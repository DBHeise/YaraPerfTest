//+build linux cgo

package YaraPerfTest

import (
	"io/ioutil"
	"os"
	"time"
	"path"

	"github.com/hillu/go-yara"
	log "github.com/sirupsen/logrus"
)

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
	files, err := ioutil.ReadDir(testFolder)
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
			hits, err := r.ScanFile(filename, 0, 0)
			if firstHit == nil {
				firstHit = hits
			}
			elapsed := time.Since(start)
			if err != nil {
				log.Warn("Error scanning file [%s]: %s", filename, err)
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
