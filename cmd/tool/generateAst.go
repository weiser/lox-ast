package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage generateAst <output_directory>")
		os.Exit(1)
	}

	outputDir := args[0]

	defineAst(outputDir, "Expr", []string{
		"Binary : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal : Object value",
		"Unary : Token operator, Expr right",
	})
}

// TODO: consider making a type that embeds *os.File that implements this, but also lets us format the golang code before we write it out
func writeWithNewline(f *os.File, s string) {
	f.WriteString(fmt.Sprintln(s))
}

func defineAst(outputDir string, baseName string, types []string) {
	path := fmt.Sprintf("%v/%v.go", outputDir, baseName)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	writeWithNewline(f, "package expr")
	writeWithNewline(f, "// DO NOT MODIFY. GENERATED VIA `go run cmd/tool/generateAst.go expr`")
	writeWithNewline(f, `import . "github.com/weiser/lox/token"`)
	writeWithNewline(f, fmt.Sprintf("type %v struct {", baseName))
	writeWithNewline(f, "}")
	writeWithNewline(f, fmt.Sprintf(`type %vInterface interface {
		Accept(evi ExprVisitorInterface)
	}`, baseName))

	defineVisitor(f, baseName, types)

	writeWithNewline(f, fmt.Sprintf("func (o *%v) Accept(evi ExprVisitorInterface) {", baseName))
	writeWithNewline(f, "evi.VisitExpr(o)")
	writeWithNewline(f, "}")

	// define AST types
	for _, typ := range types {
		splits := strings.Split(typ, ":")

		exprType := strings.TrimSpace(splits[0])
		fields := strings.TrimSpace(splits[1])

		defineType(f, baseName, exprType, fields)
	}

}

func defineVisitor(f *os.File, baseName string, types []string) {
	writeWithNewline(f, "type ExprVisitorInterface interface {")
	writeWithNewline(f, fmt.Sprintf("Visit%v(e *%v)", baseName, baseName))
	// each `type` has format `Type : ....`
	for _, typ := range types {
		splits := strings.Split(typ, ":")

		exprType := strings.TrimSpace(splits[0])
		writeWithNewline(f, fmt.Sprintf("Visit%v(e *%v)", exprType, exprType))
	}
	writeWithNewline(f, "}")

}

func defineType(f *os.File, baseName string, exprType string, fields string) {
	writeWithNewline(f, fmt.Sprintf("type %v struct {", exprType))
	writeWithNewline(f, fmt.Sprintf("*%v", baseName))
	for _, field := range strings.Split(strings.TrimSpace(fields), ",") {
		// field is 'Token operator'. needs to be "operator Token" in struct
		fs := strings.Split(strings.TrimSpace(field), " ")
		fieldt := fs[0]
		// "Object" is the java equivalent to "interface{}"
		if fieldt == "Object" {
			fieldt = "interface{}"
		} else if fieldt == baseName {
			fieldt += "Interface"
		}
		writeWithNewline(f, fmt.Sprintf("%v %v", strings.Title(fs[1]), fieldt))
	}
	writeWithNewline(f, "}")

	writeWithNewline(f, fmt.Sprintf("func (o *%v) Accept(evi ExprVisitorInterface) {", exprType))
	writeWithNewline(f, fmt.Sprintf("evi.Visit%v(o)", exprType))
	writeWithNewline(f, "}")

}
