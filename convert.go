package main

import (
	"bytes"
	evy "evylang.dev/evy/pkg/parser"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ... (other imports and your translateNode function remain)

func main() {
	if len(os.Args) < 2 { // Check for minimum number of arguments
		fmt.Println("Usage: go run main.go <directory_or_file_path>")
		os.Exit(1)
	}

	testPath := os.Args[1]

	fileInfo, err := os.Stat(testPath)
	if err != nil {
		fmt.Println("Invalid path:", err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		files, _ := ioutil.ReadDir(testPath)
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
				processFile(testPath, file.Name())
			}
		}
	} else {
		processFile(testPath, "")
	}
}

func processFile(testPath, fileName string) {
	var filePath string
	if fileName == "" {
		filePath = testPath
	} else {
		filePath = testPath + "/" + fileName
	}

	sourceCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceCode, 0)
	if err != nil {
		fmt.Println("Error parsing Go code:", err)
		return
	}
	conf := types.Config{Importer: importer.Default()}

	// types.TypeOf() requires all three maps are populated
	info := &types.Info{
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	_, err = conf.Check(filePath, fset, []*ast.File{file}, info)
	if err != nil {
		log.Fatalln(err)
	}
	evyCode := translateNode(info, file)
	if evyCode == "" {
		fmt.Println(filePath, "translation empty")
		return
	}
	evyFilePath := strings.Replace(filePath, ".go", ".evy", 1)
	err = ioutil.WriteFile(evyFilePath, []byte(evyCode), 0644)
	if err != nil {
		fmt.Println("Error writing Evy file:", err)
		return
	}

	// Execute the "evy test" command
	cmd := exec.Command("evy", "test", evyFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing 'evy test':", err)
		fmt.Println(string(output)) // Print the output for debugging
		return
	}

	fmt.Println("Exit code for", evyFilePath, "is", cmd.ProcessState.ExitCode())
}

// translateNode now converts a Go AST node to an Evy AST node
func translateNode(info *types.Info, goNode ast.Node) string {
	switch node := goNode.(type) {
	case *ast.ArrayType:
		return translateArrayType(info, node)
	case *ast.AssignStmt:
		return translateAssignStmt(info, node)
	case *ast.BadDecl:
		panic("not implemented")
	case *ast.BadExpr:
		panic("not implemented")
	case *ast.BadStmt:
		panic("not implemented")
	case *ast.BasicLit:
		return translateBasicLit(info, node)
	case *ast.BinaryExpr:
		return translateBinaryExpr(info, node)
	case *ast.BlockStmt:
		return i(translateBlockStmt(info, node))
	case *ast.BranchStmt:
		panic(node)
	case *ast.CallExpr:
		panic(node)
	case *ast.CaseClause:
		panic(node)
	case *ast.ChanType:
		panic(node)
	case *ast.CommClause:
		panic(node)
	case *ast.CompositeLit:
		panic(node)
	case *ast.DeclStmt:
		panic(node)
	case *ast.DeferStmt:
		panic(node)
	case *ast.Ellipsis:
		panic(node)
	case *ast.EmptyStmt:
		panic(node)
	case *ast.ExprStmt:
		return translateExprStmt(info, node)
	case *ast.Field:
		return translateField(info, node)
	case *ast.FieldList:
		return translateFieldList(info, node)
	case *ast.File:
		return translateFile(info, node)
	case *ast.ForStmt:
		return translateForStmt(info, node)
	case *ast.FuncDecl:
		return translateFuncDecl(info, node)
	case *ast.FuncLit:
		return translateFuncLit(info, node)
	case *ast.FuncType:
		return translateFuncType(info, node)
	case *ast.GenDecl:
		return translateGenDecl(info, node)
	case *ast.GoStmt:
		panic(node)
	case *ast.Ident:
		return translateIdent(info, node)
	case *ast.IfStmt:
		return translateIfStmt(info, node)
	case *ast.ImportSpec:
		return translateImportSpec(info, node)
	case *ast.IncDecStmt:
		return translateIncDecStmt(info, node)
	case *ast.IndexExpr:
		return translateIndexExpr(info, node)
	case *ast.InterfaceType:
		return translateInterfaceType(info, node)
	case *ast.KeyValueExpr:
		return translateKeyValueExpr(info, node)
	case *ast.LabeledStmt:
		return translateLabeledStmt(info, node)
	case *ast.MapType:
		return translateMapType(info, node)
	case *ast.Package:
		return translatePackage(info, node)
	case *ast.ParenExpr:
		return translateParenExpr(info, node)
	case *ast.RangeStmt:
		return translateRangeStmt(info, node)
	case *ast.ReturnStmt:
		return translateReturnStmt(info, node)
	case *ast.SelectorExpr:
		return translateSelectorExpr(info, node)
	case *ast.SendStmt:
		return translateSendStmt(info, node)
	case *ast.SliceExpr:
		return translateSliceExpr(info, node)
	case *ast.StarExpr:
		return translateStarExpr(info, node)
	case *ast.StructType:
		return translateStructType(info, node)
	case *ast.SwitchStmt:
		return translateSwitchStmt(info, node)
	case *ast.TypeAssertExpr:
		return translateTypeAssertExpr(info, node)
	case *ast.TypeSpec:
		return translateTypeSpec(info, node)
	case *ast.UnaryExpr:
		return translateUnaryExpr(info, node)
	case *ast.ValueSpec:
		return translateValueSpec(info, node)
	default:
		panic("")
	}
	panic("")
}

func translateExprStmt(info *types.Info, node *ast.ExprStmt) string {
	return translateExpr(info, node.X)
}

func translateBasicLit(info *types.Info, node *ast.BasicLit) string {
	switch node.Kind {
	case token.INT, token.FLOAT, token.IMAG:
		return node.Value
	case token.STRING:
		return node.Value
	case token.CHAR:
		return node.Value
	default:
		panic("Unsupported BasicLit kind")
	}
}

func translateExpr(info *types.Info, expr ast.Expr) string {
	var buf bytes.Buffer
	switch e := expr.(type) {
	case *ast.Ident:
		buf.WriteString(translateIdent(info, e))
	case *ast.BasicLit:
		buf.WriteString(e.Value)
	case *ast.BinaryExpr:
		return "(" + translateBinaryExpr(info, e) + ")"
	case *ast.UnaryExpr:
		return translateUnaryExpr(info, e)
	case *ast.ParenExpr:
		return translateParenExpr(info, e)
	case *ast.CallExpr:
		buf.WriteString(translateExpr(info, e.Fun))
		buf.WriteString(" ")
		for i, arg := range e.Args {
			if i > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString(translateExpr(info, arg))
		}
		buf.WriteString("")
	case *ast.SelectorExpr:
		buf.WriteString(translateIdent(info, e.Sel))
	case *ast.MapType:
		buf.WriteString("{}") // Use curly braces for maps
		buf.WriteString(translateExpr(info, e.Value))
	case *ast.ArrayType:
		buf.WriteString("[]")
		buf.WriteString(translateExpr(info, e.Elt))
	case *ast.CompositeLit:
		switch info.TypeOf(e).(type) { // Determine the type of the literal
		case *types.Slice:
			// Slice literal
			buf.WriteString("[")
			for i, elt := range e.Elts {
				if i > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(translateExpr(info, elt))
			}
			buf.WriteString("]")

		case *types.Map:
			// Map literal
			buf.WriteString("{")
			for i, elt := range e.Elts {
				if i > 0 {
					buf.WriteString(" ")
				}
				kvExpr := elt.(*ast.KeyValueExpr)
				buf.WriteString(translateExpr(info, kvExpr.Key)) // Key (string)
				buf.WriteString(": ")
				buf.WriteString(translateExpr(info, kvExpr.Value)) // Value
			}
			buf.WriteString("}")

		case *types.Array:
			// Array literal (translate as Evy slice)
			buf.WriteString("[")
			for i, elt := range e.Elts {
				if i > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(translateExpr(info, elt))
			}
			buf.WriteString("]")

		case *types.Struct:
			// Struct literal (translate as Evy map)
			buf.WriteString("{")
			for i, elt := range e.Elts {
				if i > 0 {
					buf.WriteString(" ")
				}
				kvExpr := elt.(*ast.KeyValueExpr)
				buf.WriteString(translateIdent(info, kvExpr.Key.(*ast.Ident))) // Field name (string)
				buf.WriteString(": ")
				buf.WriteString(translateExpr(info, kvExpr.Value)) // Value
			}
			buf.WriteString("}")
		}
		// fill this out
	// Add cases for other expression types as needed (e.g., *ast.StarExpr,
	// *ast.FuncLit, *ast.IndexExpr, etc.)

	default:
		// Handle unknown expression types by returning a placeholder or error message
		return fmt.Sprintf("/* unsupported expression type: %T */", e)
	}

	return buf.String()
}

func translateField(info *types.Info, node *ast.Field) string {
	var buf strings.Builder
	// Translate names (e.g., "x, y int")
	for i, name := range node.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(translateIdent(info, name))
	}
	// Translate the field's type (e.g., "int", "string", etc.)
	buf.WriteString(" ")
	buf.WriteString(translateExpr(info, node.Type))
	// Translate struct tags (if present)
	if node.Tag != nil {
		buf.WriteString(" ")
		buf.WriteString(node.Tag.Value) // Preserve raw tag string
	}
	return buf.String()
}

func translateFieldList(info *types.Info, node *ast.FieldList) string {
	var buf strings.Builder
	// Iterate over fields and separate them with commas
	for i, field := range node.List {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(translateField(info, field))
	}
	return buf.String()
}

func i(s string) string {
	s = strings.TrimSpace(s)        // Remove leading/trailing whitespace
	lines := strings.Split(s, "\n") // Break into lines
	for i := range lines {
		lines[i] = "    " + lines[i] // Add indent to each line
	}
	return strings.Join(lines, "\n") // Rejoin lines into a single string
}

func n(s string) string {
	return strings.TrimSpace(s) + "\n" // Trim whitespace and add newline
}

func translateForStmt(info *types.Info, node *ast.ForStmt) string {
	var buf strings.Builder
	//rangeInit, ok := node.Init.(*ast.AssignStmt)
	//if ok && len(rangeInit.Lhs) == 1 && len(rangeInit.Rhs) == 1 && rangeInit.Tok == token.DEFINE {
	//	buf.WriteString("for ")
	//	if rangeInit.Lhs[0].(*ast.Ident).Name != "_" {
	//		buf.WriteString(translateExpr(info, rangeInit.Lhs[0]))
	//		buf.WriteString(" := ")
	//	}
	//	buf.WriteString("range ")
	//	buf.WriteString(translateExpr(info, rangeInit.Rhs[0]))
	//	buf.WriteString("")
	//	buf.WriteString("\n")
	//	buf.WriteString(translateBlockStmt(info, node.Body))
	//	buf.WriteString("\nend\n")
	//	return buf.String()
	//}

	// Check if the for loop can be simplified to a while loop
	if node.Init == nil && node.Post == nil {
		buf.WriteString("while ")
		buf.WriteString(translateExpr(info, node.Cond))
		buf.WriteString("\n")
		buf.WriteString(translateBlockStmt(info, node.Body))
		buf.WriteString("\nend\n")
		return buf.String()
	}
	// Translate standard for loop
	buf.WriteString("for ")
	// Initialize
	var assignStmt *ast.AssignStmt // Declare assignStmt here
	assignStmt, ok := node.Init.(*ast.AssignStmt)
	if assignStmt != nil {
		//assignStmt, ok := node.Init.(*ast.AssignStmt) // Assign assignStmt within the if statement
		if ok {
			buf.WriteString(translateIdent(info, assignStmt.Lhs[0].(*ast.Ident))) // Assuming the LHS is a single identifier
			buf.WriteString(" := range")
		} else {
			buf.WriteString(translateStmt(info, node.Init)) // Fallback if the Init isn't a simple assignment
			buf.WriteString("; ")
		}
	}

	// Condition
	if node.Cond != nil {
		binExpr, ok := node.Cond.(*ast.BinaryExpr)
		if ok && (binExpr.Op == token.LSS || binExpr.Op == token.LEQ) { // Assuming the condition is a simple comparison
			stopExpr := translateExpr(info, binExpr.Y)
			if node.Init != nil {
				startExpr := translateExpr(info, assignStmt.Rhs[0]) // Assuming the RHS of the Init assignment is the start
				buf.WriteString(" ")
				buf.WriteString(startExpr)
				buf.WriteString(" ")
			}

			buf.WriteString(stopExpr)

			// Step
			if node.Post != nil {
				incDecStmt, ok := node.Post.(*ast.IncDecStmt)
				if ok && incDecStmt.Tok == token.INC {
					buf.WriteString(" 1") // Add step = 1 if it's a simple ++
				}
			}
			buf.WriteString(" ")
		} else {
			buf.WriteString(translateExpr(info, node.Cond)) // Fallback to standard condition translation
			buf.WriteString(" ")
		}
	}

	//// Post (should be empty for your syntax)
	//if node.Post != nil {
	//	buf.WriteString(translateStmt(info, node.Post))
	//}

	// Body of the loop
	buf.WriteString("\n")
	buf.WriteString(translateBlockStmt(info, node.Body))
	buf.WriteString("\nend\n")

	return buf.String()
}

func translateStmt(info *types.Info, stmt ast.Stmt) string {
	var buf strings.Builder
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		for i, lhs := range s.Lhs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(translateExpr(info, lhs))
		}
		buf.WriteString(" ")
		buf.WriteString(s.Tok.String())
		buf.WriteString(" ")
		for i, rhs := range s.Rhs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(translateExpr(info, rhs))
		}

	case *ast.BlockStmt:
		buf.WriteString(translateBlockStmt(info, s))

	case *ast.DeclStmt:
		buf.WriteString(translateDecl(info, s.Decl))

	case *ast.ExprStmt:
		buf.WriteString(translateExpr(info, s.X))

	case *ast.IncDecStmt:
		buf.WriteString(translateIncDecStmt(info, s))

	case *ast.IfStmt:
		buf.WriteString(translateIfStmt(info, s))

	case *ast.ForStmt:
		buf.WriteString(translateForStmt(info, s))

	case *ast.ReturnStmt:
		buf.WriteString(translateReturnStmt(info, s))

	case *ast.SwitchStmt:
		buf.WriteString(translateSwitchStmt(info, s))

	// Add cases for other statement types (e.g., *ast.BranchStmt,
	// *ast.GoStmt, *ast.DeferStmt, etc.) as needed

	default:
		return fmt.Sprintf("/* unsupported statement type: %T */", s)
	}

	return buf.String()
}

func translateDecl(info *types.Info, decl ast.Decl) string {
	var buf strings.Builder
	switch d := decl.(type) {
	case *ast.GenDecl: // General declaration (var, const, type, import)
		buf.WriteString(d.Tok.String()) // var, const, type, or import
		buf.WriteString(" ")
		for i, spec := range d.Specs {
			if i > 0 {
				buf.WriteString("\n")
			}
			switch s := spec.(type) {
			case *ast.ValueSpec: // Variable or constant declaration
				buf.WriteString(translateValueSpec(info, s))
			case *ast.TypeSpec: // Type declaration
				buf.WriteString(translateTypeSpec(info, s))
			case *ast.ImportSpec: // Import declaration
				buf.WriteString(translateImportSpec(info, s))
			default:
				return fmt.Sprintf("/* unsupported spec type: %T */", s)
			}
		}
	case *ast.FuncDecl: // Function declaration
		buf.WriteString(translateFuncDecl(info, d))

	default:
		return fmt.Sprintf("/* unsupported decl type: %T */", d)
	}

	return buf.String()
}

func translateValueSpec(info *types.Info, node *ast.ValueSpec) string {
	var buf strings.Builder
	// Names of the values being declared (e.g., "x", "y", "z")
	for i, name := range node.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(translateIdent(info, name))
	}
	// Optional type for the values
	if node.Type != nil {
		buf.WriteString(" ")
		buf.WriteString(translateExpr(info, node.Type))
	}
	// Optional initial values
	if node.Values != nil {
		buf.WriteString(" = ")
		for i, value := range node.Values {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(translateExpr(info, value))
		}
	}

	return buf.String()
}

func translateSwitchStmt(info *types.Info, node *ast.SwitchStmt) string {
	var buf strings.Builder

	buf.WriteString("switch ")

	// Optional initialization statement
	if node.Init != nil {
		buf.WriteString(translateStmt(info, node.Init))
		buf.WriteString("; ")
	}

	// Check if it's a type switch
	if node.Tag == nil {
		// No tag means it's a type switch
		buf.WriteString("true") // This is how type switches are internally represented in Go
	} else {
		// Expression-based switch
		buf.WriteString(translateExpr(info, node.Tag))
	}

	buf.WriteString(" {\n")

	for _, stmt := range node.Body.List {
		caseClause, ok := stmt.(*ast.CaseClause)
		if ok {
			// Case expressions/types
			buf.WriteString("case ")
			for i, expr := range caseClause.List {
				if i > 0 {
					buf.WriteString(", ")
				}
				if node.Tag == nil {
					// Type switch: expressions are type assertions
					buf.WriteString(translateExpr(info, expr))
				} else {
					// Expression switch: expressions are values to compare
					buf.WriteString(translateExpr(info, expr))
				}
			}
			buf.WriteString(":\n")

			// Case body
			for _, caseStmt := range caseClause.Body {
				buf.WriteString(translateStmt(info, caseStmt))
				buf.WriteString("\n")
			}
		}
	}

	buf.WriteString("}\n") // Close the switch statement

	return buf.String()
}

func translateFuncLit(info *types.Info, node *ast.FuncLit) string {
	var buf strings.Builder

	// Translate function type (parameters and results)
	buf.WriteString(translateFuncType(info, node.Type))

	// Translate the function body
	buf.WriteString(" ")
	buf.WriteString(translateBlockStmt(info, node.Body))

	return buf.String()
}

func translateFuncType(info *types.Info, node *ast.FuncType) string {
	var buf strings.Builder
	buf.WriteString("func")
	buf.WriteString("")
	buf.WriteString(translateFieldList(info, node.Params))
	buf.WriteString("")
	if node.Results != nil {
		buf.WriteString(" ")
		if len(node.Results.List) > 1 {
			buf.WriteString("")
			buf.WriteString(translateFieldList(info, node.Results))
			buf.WriteString("")
		} else {
			buf.WriteString(translateFieldList(info, node.Results))
		}
	}
	return buf.String()
}

func translateIfStmt(info *types.Info, node *ast.IfStmt) string {
	var buf strings.Builder
	buf.WriteString("if ")
	if node.Init != nil {
		buf.WriteString(translateStmt(info, node.Init))
	}
	buf.WriteString(translateExpr(info, node.Cond))
	buf.WriteString("\n")
	buf.WriteString(translateBlockStmt(info, node.Body))
	buf.WriteString("\n")
	if node.Else != nil {
		buf.WriteString(" else ")
		if elseIf, ok := node.Else.(*ast.IfStmt); ok {
			buf.WriteString(translateIfStmt(info, elseIf)) // Recursive call
		} else {
			buf.WriteString("")
			buf.WriteString(translateStmt(info, node.Else))
			buf.WriteString("\n")
		}
	}
	buf.WriteString("end\n")
	return buf.String()
}

func translateImportSpec(info *types.Info, node *ast.ImportSpec) string {
	var buf strings.Builder
	buf.WriteString("import ")

	// Optional package name/alias
	if node.Name != nil {
		buf.WriteString(node.Name.Name)
		buf.WriteString(" ")
	}

	// Import path
	buf.WriteString(node.Path.Value)

	return buf.String()
}

func translateIncDecStmt(info *types.Info, node *ast.IncDecStmt) string {
	v := translateExpr(info, node.X)
	if node.Tok == token.INC {
		return fmt.Sprintf("%v = %v + 1", v, v)
	}
	return fmt.Sprintf("%v = %v - 1", v, v)
}

func translateIndexExpr(info *types.Info, node *ast.IndexExpr) string {
	var buf strings.Builder

	// The expression being indexed (e.g., an array or slice)
	buf.WriteString(translateExpr(info, node.X))

	// Open bracket
	buf.WriteString("[")

	// The index expression
	buf.WriteString(translateExpr(info, node.Index))

	// Closing bracket
	buf.WriteString("]")

	return buf.String()
}

func translateInterfaceType(info *types.Info, node *ast.InterfaceType) string {
	var buf strings.Builder
	buf.WriteString("interface {")

	// Translate each method signature in the interface
	buf.WriteString(translateFieldList(info, node.Methods))

	buf.WriteString("}")

	return buf.String()
}
func translateKeyValueExpr(info *types.Info, node *ast.KeyValueExpr) string {
	var buf strings.Builder

	// Key expression (e.g., the "foo" in "foo: bar")
	buf.WriteString(translateExpr(info, node.Key))

	// Colon separator
	buf.WriteString(": ")

	// Value expression (e.g., the "bar" in "foo: bar")
	buf.WriteString(translateExpr(info, node.Value))

	return buf.String()
}

func translateLabeledStmt(info *types.Info, node *ast.LabeledStmt) string {
	var buf strings.Builder

	// Label name
	buf.WriteString(node.Label.Name)

	// Colon separator
	buf.WriteString(": ")

	// The statement being labeled
	buf.WriteString(translateStmt(info, node.Stmt))

	return buf.String()
}

func translateMapType(info *types.Info, node *ast.MapType) string {
	var buf strings.Builder
	buf.WriteString("map[")

	// Key type
	buf.WriteString(translateExpr(info, node.Key))

	buf.WriteString("]")

	// Value type
	buf.WriteString(translateExpr(info, node.Value))

	return buf.String()
}

func translatePackage(info *types.Info, node *ast.Package) string {
	// The package name is typically the only thing to translate here
	return node.Name
}

func translateParenExpr(info *types.Info, node *ast.ParenExpr) string {
	var buf strings.Builder

	// Opening parenthesis
	buf.WriteString("")

	// The expression inside the parentheses
	buf.WriteString(translateExpr(info, node.X))

	// Closing parenthesis
	buf.WriteString("")

	return buf.String()
}

func translateRangeStmt(info *types.Info, node *ast.RangeStmt) string {
	var buf strings.Builder
	buf.WriteString("for ")

	// Optional key/value assignments
	if node.Key != nil {
		if node.Tok == token.DEFINE { // Using := for short variable declaration
			buf.WriteString(translateExpr(info, node.Key))
			if node.Value != nil {
				buf.WriteString(", ")
				buf.WriteString(translateExpr(info, node.Value))
			}
		} else { // Using = for assignment
			buf.WriteString(translateExpr(info, node.Key))
			if node.Value != nil {
				buf.WriteString(" = ")
				buf.WriteString(translateExpr(info, node.Value))
			}
		}
		buf.WriteString(" := ")
	}

	// The expression being ranged over
	buf.WriteString("range ")
	buf.WriteString(translateExpr(info, node.X))

	// Body of the loop
	buf.WriteString(" ")
	buf.WriteString(translateBlockStmt(info, node.Body))

	return buf.String()
}
func translateReturnStmt(info *types.Info, node *ast.ReturnStmt) string {
	var buf strings.Builder
	buf.WriteString("return")
	if len(node.Results) > 0 { // Check if there are values to return
		buf.WriteString(" ")
		for i, result := range node.Results {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(translateExpr(info, result))
		}
	}
	return buf.String()
}

func translateSelectorExpr(info *types.Info, node *ast.SelectorExpr) string {
	var buf strings.Builder
	buf.WriteString(translateExpr(info, node.X))
	buf.WriteString(".")
	buf.WriteString(translateIdent(info, node.Sel))
	return buf.String()
}

func translateSendStmt(info *types.Info, node *ast.SendStmt) string {
	var buf strings.Builder
	buf.WriteString(translateExpr(info, node.Chan))
	buf.WriteString(" <- ")
	buf.WriteString(translateExpr(info, node.Value))
	return buf.String()
}

func translateSliceExpr(info *types.Info, node *ast.SliceExpr) string {
	var buf strings.Builder
	buf.WriteString(translateExpr(info, node.X))
	buf.WriteString("[")
	if node.Low != nil {
		buf.WriteString(translateExpr(info, node.Low))
	}
	buf.WriteString(":")
	if node.High != nil {
		buf.WriteString(translateExpr(info, node.High))
	}
	if node.Max != nil {
		buf.WriteString(":")
		buf.WriteString(translateExpr(info, node.Max))
	}
	buf.WriteString("]")
	return buf.String()
}

func translateStarExpr(info *types.Info, node *ast.StarExpr) string {
	var buf strings.Builder
	buf.WriteString("*")
	buf.WriteString(translateExpr(info, node.X))
	return buf.String()
}

func translateStructType(info *types.Info, node *ast.StructType) string {
	var buf strings.Builder
	buf.WriteString("struct {")
	buf.WriteString(translateFieldList(info, node.Fields))
	buf.WriteString("}")
	return buf.String()
}

func translateTypeAssertExpr(info *types.Info, node *ast.TypeAssertExpr) string {
	var buf strings.Builder
	buf.WriteString(translateExpr(info, node.X))
	buf.WriteString(".(")
	if node.Type != nil {
		buf.WriteString(translateExpr(info, node.Type))
	}
	buf.WriteString(")")
	return buf.String()
}

func translateTypeSpec(info *types.Info, node *ast.TypeSpec) string {
	var buf strings.Builder
	buf.WriteString("type ")
	buf.WriteString(node.Name.Name)
	buf.WriteString(" ")
	buf.WriteString(translateExpr(info, node.Type))
	return buf.String()
}

func translateUnaryExpr(info *types.Info, node *ast.UnaryExpr) string {
	str := "("
	str += node.Op.String()
	str += "("
	str += translateExpr(info, node.X)
	str += ")"
	str += ")"
	return str
}

// ... other parts of your translation code ...

func translateBinaryExpr(info *types.Info, node *ast.BinaryExpr) string {
	x := translateNode(info, node.X)
	y := translateNode(info, node.Y)
	return x + " " + translateOperator(node.Op) + " " + y
}

func translateArrayType(info *types.Info, node *ast.ArrayType) string {
	return "arr"
}

func translateGenDecl(info *types.Info, node *ast.GenDecl) string {
	// Handle different types of declarations within a GenDecl
	switch node.Tok {
	case token.VAR, token.CONST:
		// Handle variable declarations
		varSpecs := make([]string, len(node.Specs))
		for i, spec := range node.Specs {
			varSpecs[i] = translateNode(info, spec)
		}
		return strings.Join(varSpecs, "\n")
	case token.IMPORT:
		return ""
	default:
		panic("Unsupported GenDecl token") // Handle unsupported declaration types
	}
}

func translateFile(info *types.Info, file *ast.File) string {
	var statements []string
	for _, decl := range file.Decls {
		stmt := translateNode(info, decl)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return strings.Join(statements, "\n")
}

func translateFuncDecl(info *types.Info, funcDecl *ast.FuncDecl) string {
	var buf bytes.Buffer
	// Function signature
	name := translateIdent(info, funcDecl.Name)
	fmt.Fprintf(&buf, "func %s", name)

	// Parameters
	for i, field := range funcDecl.Type.Params.List {
		for j, name := range field.Names {
			if j > 0 {
				buf.WriteString(" ")
			}
			fmt.Fprintf(&buf, "%s ", translateIdent(info, name))
		}
		if len(field.Names) > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(translateExpr(info, field.Type))
		if i < len(funcDecl.Type.Params.List)-1 {
			buf.WriteString(" ")
		}
	}
	buf.WriteString(" ")
	//Results
	if funcDecl.Type.Results != nil {
		buf.WriteString("")
		for i, field := range funcDecl.Type.Results.List {
			if i > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString(translateExpr(info, field.Type))
		}
		buf.WriteString(" ")
	}
	buf.WriteString("\n")
	// Function body
	if funcDecl.Body != nil {
		for _, decl := range funcDecl.Body.List {
			buf.WriteString(i(translateStmt(info, decl)))
			buf.WriteString("\n")
		}
		buf.WriteString("end\n")
	} else {
		buf.WriteString("\n\t//Empty Function \n")
	}
	if name == "main" {
		buf.WriteString("main\n")
	}
	return buf.String()
}

func translateIdent(info *types.Info, ident *ast.Ident) string {
	str := ident.String()
	switch str {
	case "Println":
		return "print"
	case "Printf":
		return "printf"
	default:
		return str
	}
}

func translateBlockStmt(info *types.Info, blockStmt *ast.BlockStmt) string {
	var statements []string
	for _, goStmt := range blockStmt.List {
		evyStmt := translateNode(info, goStmt) // Recursively translate each statement in the block
		if evyStmt != "" {                     // Ignore unsupported statements (if any)
			statements = append(statements, evyStmt)
		}
	}
	return i(strings.Join(statements, "\n"))
}

func translateAssignStmt(into *types.Info, assignStmt *ast.AssignStmt) string {
	var lhs, rhs string
	for i, lhsExpr := range assignStmt.Lhs {
		lhs += translateNode(into, lhsExpr)
		rhs += translateNode(into, assignStmt.Rhs[i])
	}
	return lhs + " = " + rhs
}

func toEvyType(in string) *evy.Type {
	switch {
	case strings.Contains(in, "float"), strings.Contains(in, "int"):
		return evy.NUM_TYPE
	case in == "string":
		return evy.STRING_TYPE
	case in == "bool":
		return evy.BOOL_TYPE
	case in == "any", in == "interface{}":
		return evy.ANY_TYPE
	case strings.HasPrefix(in, "[]"):
		return &evy.Type{
			Name: evy.ARRAY,
			Sub:  toEvyType(in[2:]),
		}
	case strings.HasPrefix(in, "map[string]"): // other map types not supported in evy
		return &evy.Type{
			Name: evy.MAP,
			Sub:  toEvyType(in[len("map[string]"):]),
		}
	}
	panic("")
}

func translateOperator(op token.Token) string {
	switch op {
	case token.ADD:
		return evy.OP_PLUS.String()
	case token.SUB:
		return evy.OP_MINUS.String()
	case token.MUL:
		return evy.OP_ASTERISK.String()
	case token.QUO:
		return evy.OP_SLASH.String()
	case token.REM:
		return evy.OP_PERCENT.String()
	case token.EQL:
		return evy.OP_EQ.String()
	case token.NEQ:
		return evy.OP_NOT_EQ.String()
	case token.LSS:
		return evy.OP_LT.String()
	case token.GTR:
		return evy.OP_GT.String()
	case token.LEQ:
		return evy.OP_LTEQ.String()
	case token.GEQ:
		return evy.OP_GTEQ.String()
	case token.LAND:
		return evy.OP_AND.String()
	case token.LOR:
		return evy.OP_OR.String()
	default:
		fmt.Printf("Unsupported operator: %s\n", op) // For debugging
		return evy.OP_ILLEGAL.String()
	}
}
