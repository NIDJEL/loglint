package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "check log messages for style and sensitive data",
	Run:  run,
}

var sensitiveWords = []string{
	"password",
	"passwd",
	"pwd",
	"token",
	"api_key",
	"apikey",
	"secret",
	"access_key",
	"private_key",
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

			if len(call.Args) > 0 && containsSensitiveData(call.Args[0]) {
				pass.Reportf(call.Pos(), "log message should not contain sensitive data")
			}

			msg, ok := firstStingArg(call)
			if !ok {
				return true
			}

			if startsWithUpper(msg) {
				pass.Reportf(call.Pos(), "log message should start with lowercase letter")
			}

			if hasNonEnglishLetters(msg) {
				pass.Reportf(call.Pos(), "log message should be written in English")
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

func hasNonEnglishLetters(msg string) bool {
	for _, r := range msg {
		if !unicode.IsLetter(r) {
			continue
		}

		if r >= 'a' && r <= 'z' {
			continue
		}

		if r >= 'A' && r <= 'Z' {
			continue
		}

		return true
	}

	return false
}

func containsSensitiveData(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return false
		}

		value, err := strconv.Unquote(e.Value)
		if err != nil {
			return false
		}

		return hasSensitiveWord(value)

	case *ast.Ident:
		return hasSensitiveWord(e.Name)

	case *ast.BinaryExpr:
		return containsSensitiveData(e.X) || containsSensitiveData(e.Y)

	default:
		return false
	}
}

func hasSensitiveWord(text string) bool {
	text = strings.ToLower(text)

	for _, word := range sensitiveWords {
		if strings.Contains(text, word) {
			return true
		}
	}

	return false
}
