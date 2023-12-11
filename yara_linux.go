//go:build linux || cgo
// +build linux cgo

package YaraPerfTest

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hillu/go-yara/v4"
	log "github.com/sirupsen/logrus"
)

func printMatches(item string, m []yara.MatchRule, err error) {
	if err != nil {
		log.Printf("%s: error: %s", item, err)
		return
	}
	if len(m) == 0 {
		log.Printf("%s: no matches", item)
		return
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s: [", item)
	for i, match := range m {
		if i > 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, "%s:%s", match.Namespace, match.Rule)
	}
	fmt.Fprint(buf, "]")
	log.Print(buf.String())
}

// RunYara - tests the yara rule against the files the given number of times
func RunYara(yaraRuleFile string, numTimes int, testFolder string) ([]YaraResult, error) {

	//Create Yara Compiler
	c, err := yara.NewCompiler()
	if err != nil {
		log.Warn("Failed to initialize YARA compiler")
		return nil, err
	}

	//Load Rule files
	f, err := os.Open(yaraRuleFile)
	if err != nil {
		log.Warn("Could not open rule file")
		return nil, err
	}
	err = c.AddFile(f, "")
	f.Close()
	if err != nil {
		log.Warn("Could not parse rule file")
		return nil, err
	}

	//Compile Rules
	start := time.Now()
	r, err := c.GetRules()
	elapsed := time.Since(start)
	if err != nil {
		log.Warn("Failed to compile rules")
		return nil, err
	}
	log.Printf("Rule Compliation Time: %s", elapsed)

	//Get files
	files, err := os.ReadDir(testFolder)
	if err != nil {
		log.Warn("Failed to get files in folder")
		return nil, err
	}

	results := make([]YaraResult, 0)

	//Scan Files
	for _, f := range files {
		filename := path.Join(testFolder, f.Name())
		log.Printf("Scanning file %s... ", filename)
		times := make([]float64, numTimes)
		var firstHit []yara.MatchRule

		for i := 0; i < numTimes; i++ {
			start := time.Now()
			s, _ := yara.NewScanner(r)
			var m yara.MatchRules
			err := s.SetCallback(&m).ScanFile(filename)
			printMatches(filename, m, err)

			elapsed := time.Since(start)
			if err != nil {
				fmt.Printf("Error scanning file [%s]: %s\n", filename, err)
				return nil, err
			}
			times[i] = float64(elapsed)
		}

		result := YaraResult{
			File:      filename,
			Stats:     Statistics{},
			FirstHits: firstHit,
		}
		result.Stats.Calculate(times)
		results = append(results, result)
	}

	return results, nil
}
