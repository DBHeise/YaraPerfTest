package main

import (
	"fmt"
	"flag"
	"encoding/json"
	"YaraPerfTest"
	log "github.com/sirupsen/logrus"
)

var (
	rulefile	string
	numtimes    int
	testFolder  string
)	

func init() {
	flag.IntVar(&numtimes, "times", 10, "number of times to test each rule/file")
	flag.StringVar(&testFolder, "test", "", "folder containing test files")
	flag.StringVar(&rulefile, "rule", "", "yara rule file")
}

func main() {
	flag.Parse()

	r, err := YaraPerfTest.RunYara(rulefile, numtimes, testFolder)
	if err != nil {
		log.Error(err)
	} else {
		strB, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			log.Error(err)
		} else {
			fmt.Println(string(strB))
		}
		
	}
}