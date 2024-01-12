package exitcheck

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestMainExitCheck(t *testing.T) {
	// TODO: Add test cases.
	analysistest.Run(t, analysistest.TestData(), MainExitCheckAnalyzer, "./good")
}
