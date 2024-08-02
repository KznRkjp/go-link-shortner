package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := make([]*analysis.Analyzer, 0, 245)
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	multichecker.Main(mychecks...)
}
