package main

import (
	"github.com/benderr/metrics/pkg/analyzers"
	"github.com/benderr/metrics/pkg/osexitanalyzer"
	"github.com/tdakkota/asciicheck"
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis/multichecker"
)

const Config = `./staticcheck.config`

func main() {

	a, err := analyzers.New(Config)
	if err != nil {
		panic(err)
	}

	checks := a.GetAnalyzers()

	checks = append(checks, mnd.Analyzer)

	checks = append(checks, asciicheck.NewAnalyzer())

	checks = append(checks, osexitanalyzer.Analyzer)

	multichecker.Main(
		checks...,
	)
}
