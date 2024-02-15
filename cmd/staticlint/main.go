// Package staticlint lint your code and make it great again.
//
// Include all analysis/passes linters.
//
// Include analyzers from staticcheck (with config file for include SA, ST and S rules).
//
// Include acshiicheck package, checks that all code identifiers does not
// have non-ASCII symbols in the name.
//
// Include osexitanalyzer package, for check os.Exit function in main
// Look all analyzers (run)
//
//	cmd/staticlint/staticlint help
//
// Customize lint with arguments, for example,
// if you want to run only `osexit` analyzer:
//
//	cmd/staticlint/staticlint -osexit <package_path>
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
