package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	"go-url-shortener/internal/analyzers/mainosexit"
)

// main function of the staticlint tool performs the static code analysis via predefined and custom analyzers.
// It combines standard, 3rd-party and custom analyzers into a single slices.
// The analysis is performed via the multichecker tool.
func main() {
	mychecks := append(getStandardChecks(), append(getStaticChecks(), getCustomChecks()...)...)

	multichecker.Main(
		mychecks...,
	)
}

// getCustomChecks enables the custom static analyzers.
// mainosexit analyzer prevents the main function of the main package to explicitly use the os.Exit function.
func getCustomChecks() []*analysis.Analyzer {
	return []*analysis.Analyzer{mainosexit.Analyzer}
}

// getStandardChecks enables the set of the standard static analyzers.
// atomic analyzer checks for common mistakes using the sync/atomic package.
// httpresponse analyzer checks for mistakes using HTTP responses.
// loopclosure checks references to loop variables from within nested functions.
// nilfunc checks for useless comparisons between functions and nil.
// shadow checks for possible unintended shadowing of variables.
// stringintconv checks for string(int) conversions.
// structtag checks that struct field tags conform to reflect.StructTag.Get.
// unmarshal reports passing non-pointer or non-interface values to unmarshal.
// usesgenerics detects whether a package uses generics features.
func getStandardChecks() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		atomic.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		shadow.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		unmarshal.Analyzer,
		usesgenerics.Analyzer,
	}
}

// getStaticChecks enables all available analyzers from the staticcheck package.
func getStaticChecks() []*analysis.Analyzer {
	var mychecks []*analysis.Analyzer
	for i := 0; i < len(staticcheck.Analyzers); i++ {
		mychecks = append(mychecks, staticcheck.Analyzers[i].Analyzer)
	}
	return mychecks
}
