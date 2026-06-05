package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "check log messages for style and sensitive data",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !isLogCall(call) {
				return true
			}

			msg, ok := firstStingArg(call)
			if !ok {
				return true
			}

			if startsWithUpper(msg) {
				pass.Reportf(call.Pos(), "log message should start with lowercase letter")
			}

			if hasBadChars(msg) {
				pass.Reportf(call.Pos(), "log message should not contain special characters or emoji")
			}

			return true
		})
	}

	return nil, nil
}

func isLogCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	switch sel.Sel.Name {
	case "Debug", "Info", "Warn", "Error":
		return true
	default:
		return false
	}
}

func firstStingArg(call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}

	arg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return "", false
	}

	if arg.Kind != token.STRING {
		return "", false
	}

	msg, err := strconv.Unquote(arg.Value)
	if err != nil {
		return "", false
	}

	return msg, true

}

func startsWithUpper(msg string) bool {
	for _, r := range msg {
		if unicode.IsLetter(r) {
			return unicode.IsUpper(r)
		}
	}

	return false
}

func hasBadChars(msg string) bool {
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}

		return true
	}

	return false
}
