package myanalyzer_test

import (
	"testing"
    "golang.org/x/tools/go/analysis/analysistest"
    "my"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, myanalyzer.Analyzer, "a")
}