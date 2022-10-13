package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"net/http"
	"sort"
	"strings"
	"text/template"

	_ "embed"
)

const (
	ExprErrorIdent = "error"

	MethodInner    = "Inner"
	MethodResponse = "Response"

	FeatureCache = "api/cache"
)

var (
	fileContent, _ = read(join(CurrentDir, CurrentFile))
)

func genApi(_ *cobra.Command, _ []string) error {
	inspectCtx, err := inspectApi(join(CurrentDir, CurrentFile), LineNum+1)
	if err != nil {
		return fmt.Errorf("inspectApi(%s, %d): %w", quote(join(CurrentDir, CurrentFile)), LineNum, err)
	}

	if !checkResponse(inspectCtx.Methods) {
		return fmt.Errorf("checkResponse: no '%s() T' method found in Interface", MethodResponse)
	}

	for _, method := range inspectCtx.Methods {
		if method.Ident != MethodResponse && method.Ident != MethodInner {
			if l := len(method.Out); l == 0 || !checkErrorType(method.Out[l-1]) {
				return fmt.Errorf("checkErrorType: no 'error' found in method %s returned value",
					quote(method.Ident))
			}
		}

		if (method.Ident == MethodResponse || method.Ident == MethodInner) &&
			(len(method.In) != 0 || len(method.Out) != 1) {
			return fmt.Errorf(
				"%s method can only have no income params "+
					"and 1 returned value", quote(method.Ident))
		}

		if method.Ident == MethodResponse {
			if !checkResponseType(method) {
				return fmt.Errorf(
					"checkResponseType: returned type of %s "+
						"should be kind of *ast.Ident or *ast.StarExpr",
					quote(MethodResponse))
			}
		}

		if len(method.Out) > 2 {
			return fmt.Errorf("%s method expects 2 returned value at most, got %d",
				quote(method.Ident),
				len(method.Out))
		}
	}

	code, err := genApiCode(inspectCtx)
	if err != nil {
		return fmt.Errorf("genApiCode: \n\n%#v\n\n%s", inspectCtx, err.Error())
	}

	fmtCode, err := format.Source(code)
	if err != nil {
		return fmt.Errorf("format.Source: \n\n%s\n\n%w", code, err)
	}

	if output == "" {
		output = "api.go"
	}

	if err := write(join(CurrentDir, output), fmtCode, FileMode); err != nil {
		return fmt.Errorf("os.WriteFile(%s, %0x)", join(CurrentDir, output), FileMode)
	}

	return nil
}

type ApiContext struct {
	Package  string
	Ident    string
	Generics map[string]ast.Expr
	Methods  []*Method
	Features []string
}

func (ctx *ApiContext) SortGenerics() []string {
	idents := make([]string, 0, len(ctx.Generics))
	for k := range ctx.Generics {
		idents = append(idents, k)
	}
	sort.Slice(idents, func(i, j int) bool {
		return ctx.Generics[idents[i]].Pos() < ctx.Generics[idents[j]].Pos()
	})
	return idents
}

func (ctx *ApiContext) GenericsRepr(withType bool) string {
	if len(ctx.Generics) == 0 {
		return ""
	}

	var dst bytes.Buffer
	dst.WriteByte('[')
	for index, name := range ctx.SortGenerics() {
		expr := ctx.Generics[name]
		dst.WriteString(name)
		if withType {
			dst.WriteByte(' ')
			dst.WriteString(getRepr(expr, fileContent))
		}
		if index < len(ctx.Generics)-1 {
			dst.WriteString(", ")
		}
	}
	dst.WriteByte(']')

	return dst.String()
}

func (ctx *ApiContext) HasFeature(feature string) bool {
	for _, current := range ctx.Features {
		if current == feature {
			return true
		}
	}
	return false
}

func (ctx *ApiContext) HasHeader() bool {
	for _, method := range ctx.Methods {
		if method.Header != "" {
			return true
		}
	}
	return false
}

func (ctx *ApiContext) HasInner() bool {
	return hasInner(ctx.Methods)
}

func (ctx *ApiContext) InnerType() ast.Node {
	for _, method := range ctx.Methods {
		if method.Ident == MethodInner {
			return method.Out[0]
		}
	}
	return nil
}

func inspectApi(file string, line int) (*ApiContext, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, fileContent, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var (
		genDecl   *ast.GenDecl
		typeSpec  *ast.TypeSpec
		ifaceType *ast.InterfaceType
	)

inspectDecl:
	for _, declIface := range f.Decls {
		if hit(fset, declIface, line) {
			if decl := declIface.(*ast.GenDecl); decl.Tok == token.TYPE {
				genDecl = decl
				break inspectDecl
			}
		}
	}

	if genDecl == nil {
		return nil, fmt.Errorf(
			"no available 'Interface' type declaration (*ast.GenDecl) found, "+
				"available *ast.GenDecl are: \n\n"+
				"%s\n\n", concat(nodeMap(f.Decls, fmtNode), "\n"))
	}

inspectType:
	for _, specIface := range genDecl.Specs {
		if hit(fset, specIface, line) {
			spec := specIface.(*ast.TypeSpec)
			if iface, ok := spec.Type.(*ast.InterfaceType); ok && hit(fset, iface, line) {
				typeSpec = spec
				ifaceType = iface
				break inspectType
			}
		}
	}

	if ifaceType == nil {
		return nil, fmt.Errorf(
			"no available 'Interface' type declaration (*ast.InterfaceType) found, "+
				"available *ast.GenDecl are: \n\n"+
				"%s\n\n", concat(nodeMap(f.Decls, fmtNode), "\n"))
	}

	for _, method := range ifaceType.Methods.List {
		if !checkInput(method.Type.(*ast.FuncType)) {
			return nil, fmt.Errorf(""+
				"input params for method %s should "+
				"contain 'Name' and 'Type' both",
				quote(method.Names[0].Name))
		}
	}

	apiFeatures := make([]string, 0, len(features))
	for _, feature := range features {
		if hasPrefix(feature, "api") {
			apiFeatures = append(apiFeatures, feature)
		}
	}

	generics := make(map[string]ast.Expr, 10)
	if typeSpec.TypeParams != nil {
		for _, param := range typeSpec.TypeParams.List {
			for _, name := range param.Names {
				generics[name.Name] = param.Type
			}
		}
	}

	return &ApiContext{
		Package:  PackageName,
		Ident:    typeSpec.Name.Name,
		Generics: generics,
		Methods:  nodeMap(ifaceType.Methods.List, inspectMethod),
		Features: apiFeatures,
	}, nil
}

type Method struct {
	Meta   string
	Header string
	Ident  string
	In     map[string]ast.Expr
	Out    []ast.Expr
}

func (method *Method) SortIn() []string {
	idents := make([]string, 0, len(method.In))
	for k := range method.In {
		idents = append(idents, k)
	}
	sort.Slice(idents, func(i, j int) bool {
		return method.In[idents[i]].Pos() < method.In[idents[j]].Pos()
	})
	return idents
}

func (method *Method) MetaArgs() []string {
	rawArgs := splitArgs(method.Meta)
	args := make([]string, 0, len(rawArgs))
	for i := 0; i < len(rawArgs); i++ {
		if rawArgs[i] != "" && rawArgs[i] != " " {
			args = append(args, rawArgs[i])
		}
	}
	return args
}

func (method *Method) TmplURL() string {
	args := method.MetaArgs()
	return args[len(args)-1]
}

var availableMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
}

func (method *Method) MethodHTTP() string {
	args := method.MetaArgs()
	for _, httpMethod := range availableMethods {
		if toUpper(args[len(args)-2]) == httpMethod {
			return httpMethod
		}
	}
	return ""
}

func (method *Method) ReturnSlice() bool {
	return len(method.Out) > 1 && isSlice(method.Out[0])
}

func inspectMethod(node ast.Node) (method *Method) {
	field := node.(*ast.Field)
	method = new(Method)
	if field.Doc != nil {
		method.Meta = trimSpace(
			trimPrefix(field.Doc.List[0].Text,
				"//"),
		)
		var buffer bytes.Buffer
		for _, header := range field.Doc.List[1:] {
			buffer.WriteString(trimSpace(trimPrefix(header.Text, "//")))
			buffer.WriteString("\r\n")
		}
		method.Header = buffer.String()
	}
	method.Ident = field.Names[0].Name
	funcType := field.Type.(*ast.FuncType)
	inParams := funcType.Params.List
	method.In = make(map[string]ast.Expr, len(inParams))
	for _, param := range inParams {
		for _, name := range param.Names {
			method.In[name.Name] = param.Type
		}
	}
	outParams := funcType.Results.List
	method.Out = make([]ast.Expr, 0, len(outParams))
	for _, param := range outParams {
		method.Out = append(method.Out, param.Type)
	}
	return method
}

func checkResponse(methods []*Method) bool {
	for _, method := range methods {
		if method.Ident == MethodResponse {
			return true
		}
	}
	return false
}

func checkResponseType(method *Method) bool {
	switch method.Out[0].(type) {
	case *ast.Ident, *ast.StarExpr:
		return true
	default:
		return false
	}
}

func hasInner(methods []*Method) bool {
	for _, method := range methods {
		if method.Ident == MethodInner {
			return true
		}
	}
	return false
}

func checkInput(method *ast.FuncType) bool {
	for _, param := range method.Params.List {
		if len(param.Names) == 0 {
			return false
		}
	}
	return true
}

func checkErrorType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == ExprErrorIdent
}

//go:embed templates/api.tmpl
var ApiTemplate string

func genApiCode(ctx *ApiContext) ([]byte, error) {
	tmpl, err := template.
		New("loadc(api)").
		Funcs(template.FuncMap{
			"quote": quote,
			"sub": func(x, y int) int {
				return x - y
			},
			"getRepr": func(node ast.Node) string {
				return getRepr(node, fileContent)
			},
			"methodResp": func() string {
				return MethodResponse
			},
			"isResponse": func(ident string) bool {
				return ident == MethodResponse
			},
			"isInner": func(ident string) bool {
				return ident == MethodInner
			},
			"newType": func(expr ast.Expr) string {
				return newType(expr, fileContent)
			},
			"httpMethodHasBody": func(method string) bool {
				switch method {
				case http.MethodGet:
					return false
				case http.MethodPost, http.MethodPut:
					return true
				default:
					return false
				}
			},
			"headerHasBody": func(header string) bool {
				if index := strings.Index(header, "\r\n\r\n"); index != -1 {
					return len(header[index+4:]) > 0
				}
				return false
			},
		}).
		Parse(ApiTemplate)

	if err != nil {
		return nil, err
	}

	var dst bytes.Buffer
	if err = tmpl.Execute(&dst, ctx); err != nil {
		return nil, err
	}

	return dst.Bytes(), nil
}
