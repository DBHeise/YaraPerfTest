package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"YaraPerfTest"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

var (
	staticFolder  string
	host          string
	port          string
	testFolder    string
	defaultTimes  int
)


func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	flag.StringVar(&staticFolder, "static", ".\\static", "the directory to serve static files from.")
	flag.StringVar(&host, "host", "", "the host to listen with/on")
	flag.StringVar(&port, "port", "1234", "the port to listen on")
	flag.StringVar(&testFolder, "test", "~/testfiles", "folder containing files to test")
	flag.IntVar(&defaultTimes, "time", 10, "number of times to test each rule/file")

}

func writeTempFile(content string) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "test.*.yara")
	if err != nil {
		return tmpfile, err
	}

	if _, err = tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return tmpfile, err
	}
	if err = tmpfile.Close(); err != nil {
		return tmpfile, err
	}
	return tmpfile, err
}


func testRule(c echo.Context) error {
	log.Debug("Testing rule")
	yaraRule := c.FormValue("yaraRule")

	//Temp file
	tmpfile, err := writeTempFile(yaraRule)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}	
	defer os.Remove(tmpfile.Name())


	r, err := YaraPerfTest.RunYara(tmpfile.Name(), defaultTimes, testFolder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, r)
	}
}

func main() {
	flag.Parse()
	log.RegisterExitHandler(func() {
		if err := recover(); err != nil {
			fmt.Println(errors.Wrap(err, 2).ErrorStack())
		}
	})

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.DisableHTTP2 = true

	// echo Logger interface friendly wrapper around logrus logger to use it for default echo logger
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
	e.Use(logrusmiddleware.Hook())

	// Other Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   staticFolder,
		Browse: false,
		Index:  "index.html",
	}))

	e.POST("/Test", testRule)

	// Create server
	srv := &http.Server{
		Addr:         host + ":" + port,
		WriteTimeout: 10 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}

	// Start server
	go func() {
		if err := e.StartServer(srv); err != nil {
			log.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
