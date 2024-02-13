package osexitanalyzer_test

import (
	"testing"

	"github.com/benderr/metrics/pkg/osexitanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexitanalyzer.Analyzer, "./...")
}
