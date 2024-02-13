package analyzers

import (
	"encoding/json"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

type configData struct {
	StaticcheckAll bool     // Include all staticcheck [SA] rules
	Staticcheck    []string // Select only custom staticcheck [SA] rules
	SimpleAll      bool     // Include all simple [S] rules
	Simple         []string // Include only custom simple [S] rules
	StylecheckAll  bool     // Include all styles [ST] rules
	Stylecheck     []string // Include only custom style [ST] rules
}

type analyzerList struct {
	config configData
}

func New(config string) (*analyzerList, error) {
	conf, err := loadConfig(config)
	if err != nil {
		return nil, err
	}

	return &analyzerList{config: *conf}, nil
}

func loadConfig(config string) (*configData, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	var cfg configData
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (a *analyzerList) loadStaticCheckAnalyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0)
	staticChecks := mapKeys(a.config.Staticcheck)
	simpleChecks := mapKeys(a.config.Simple)
	styleChecks := mapKeys(a.config.Stylecheck)

	for _, v := range staticcheck.Analyzers {
		if staticChecks[v.Analyzer.Name] || a.config.StaticcheckAll {
			checks = append(checks, v.Analyzer)
		}
	}

	for _, v := range simple.Analyzers {
		if simpleChecks[v.Analyzer.Name] || a.config.SimpleAll {
			checks = append(checks, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if styleChecks[v.Analyzer.Name] || a.config.StylecheckAll {
			checks = append(checks, v.Analyzer)
		}
	}
	return checks
}

func (a *analyzerList) loadPassesAnaluzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		shadow.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}
}

func (a *analyzerList) GetAnalyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0)

	passes := a.loadPassesAnaluzers()
	static := a.loadStaticCheckAnalyzers()

	checks = append(checks, passes...)
	checks = append(checks, static...)
	return checks
}

func mapKeys(keys []string) map[string]bool {
	staticChecks := make(map[string]bool)

	for _, v := range keys {
		staticChecks[v] = true
	}
	return staticChecks
}
