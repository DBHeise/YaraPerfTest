package YaraPerfTest

import (
	"time"
	"github.com/montanaflynn/stats"
)

type Statistics struct {
	Min           float64 `json:"Min"`
	Max           float64 `json:"Max"`
	Mean          float64 `json:"Mean"`
	HarmonicMean  float64 `json:"HarmonicMean"`
	GeometricMean float64 `json:"GeometricMean"`
	StdDevP       float64 `json:"StdDevP"`
	StdDevS       float64 `json:"StdDevS"`
}

func (s *Statistics) Calculate(data []float64) {
	s.Min, _ = stats.Min(data)
	s.Max, _ = stats.Max(data)
	s.Mean, _ = stats.Mean(data)
	s.HarmonicMean, _ = stats.HarmonicMean(data)
	s.GeometricMean, _ = stats.GeometricMean(data)
	s.StdDevP, _ = stats.StdDevP(data)
	s.StdDevS, _ = stats.StdDevS(data)

	//Convert to seconds
	s.Min = s.Min / float64(time.Second)
	s.Max = s.Max / float64(time.Second)
	s.Mean = s.Mean / float64(time.Second)
	s.HarmonicMean = s.HarmonicMean / float64(time.Second)
	s.GeometricMean = s.GeometricMean / float64(time.Second)
	s.StdDevP = s.StdDevP / float64(time.Second)
	s.StdDevS = s.StdDevS / float64(time.Second)
}

type YaraResult struct {
	File      string
	Stats     Statistics
	FirstHits interface{}
}
