package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"text/template"

	_ "embed"
)

const (
	SqlExt  = ".sql"
	TmplExt = ".tmpl"
)

func genSql(_ *cobra.Command, args []string) error {
	dir := args[0]
	if !isAbs(dir) {
		dir = join(CurrentDir, dir)
	}

	stat, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("os.Stat(%s): %w", quote(dir), err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("expects directory, got file from %s", quote(dir))
	}

	dirEntries, err := list(dir)
	if err != nil {
		return fmt.Errorf("os.ReadDir(%s): %w", quote(dir), err)
	}

	var (
		sqlMap  = make(map[string]string, len(dirEntries))
		tmplMap = make(map[string]string, len(dirEntries))
	)

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		var (
			targetMap map[string]string
			targetExt string

			entryName = entry.Name()
		)

		switch ext := ext(entryName); ext {
		case SqlExt:
			targetMap = sqlMap
			targetExt = SqlExt
		case TmplExt:
			targetMap = tmplMap
			targetExt = TmplExt
		default:
			continue
		}

		content, err := read(join(dir, entryName))
		if err != nil {
			return fmt.Errorf("os.ReadFile(%s): %w", quote(join(dir, entryName)), err)
		}

		targetMap[trimSuffix(base(entryName), targetExt)] = string(content)
	}

	inspectCtx, err := inspectSql(join(CurrentDir, CurrentFile), LineNum+1)
	if err != nil {
		return fmt.Errorf("inspectSql(%s, %d): %w", quote(join(CurrentDir, CurrentFile)), LineNum, err)
	}

	code, err := genSqlCode(PackageName, inspectCtx.Ident.Name, pointer, sqlMap, tmplMap)
	if err != nil {
		return fmt.Errorf("genSqlCode(%s, %s): %w", quote(PackageName), quote(inspectCtx.Ident.Name), err)
	}

	fmtCode, err := format.Source(code)
	if err != nil {
		return fmt.Errorf("format.Source: \n\n%s\n\n%w", code, err)
	}

	if output == "" {
		output = "sql.go"
	}

	if err := write(join(CurrentDir, output), fmtCode, FileMode); err != nil {
		return fmt.Errorf("os.WriteFile(%s, %0x)", join(CurrentDir, output), FileMode)
	}

	return nil
}

type SqlContext struct {
	Ident *ast.Ident
}

func inspectSql(file string, line int) (*SqlContext, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return nil, err
	}

	var (
		genDecl   *ast.GenDecl
		valueSpec *ast.ValueSpec
	)

inspectDecl:
	for _, declIface := range f.Decls {
		if hit(fset, declIface, line) {
			if decl := declIface.(*ast.GenDecl); decl.Tok == token.VAR {
				genDecl = decl
				break inspectDecl
			}
		}
	}

	if genDecl == nil {
		return nil, fmt.Errorf(
			"no available 'SqlLoader' variable declaration (*ast.GenDecl) found, "+
				"available *ast.GenDecl are: \n\n"+
				"%s\n\n", concat(nodeMap(f.Decls, fmtNode), "\n"))
	}

inspectSpec:
	for _, specIface := range genDecl.Specs {
		if hit(fset, specIface, line) {
			spec := specIface.(*ast.ValueSpec)
			if n := len(spec.Names); n > 1 {
				return nil, fmt.Errorf(
					"the variable declaration of 'SqlLoader' "+
						"should have only one variable name, "+
						"got %d", n)
			}
			valueSpec = spec
			break inspectSpec
		}
	}

	if valueSpec == nil {
		return nil, fmt.Errorf(
			"no available 'SqlLoader' variable declaration (*ast.ValueSpec) found, "+
				"available *ast.ValueSpec are: \n\n"+
				"%s\n\n", concat(nodeMap(genDecl.Specs, fmtNode), "\n"))
	}

	return &SqlContext{Ident: valueSpec.Names[0]}, nil
}

//go:embed templates/sql.tmpl
var SqlTemplate string

type TmplArgs struct {
	Package string
	Ident   string
	Pointer bool
	SqlMap  map[string]string
	TmplMap map[string]string
}

func genSqlCode(
	pkg string,
	ident string,
	pointer bool,
	sqlMap map[string]string,
	tmplMap map[string]string,
) ([]byte, error) {
	tmpl, err := template.
		New("loadc(sql)").
		Funcs(template.FuncMap{
			"quote":    quote,
			"camelize": camelize,
		}).
		Parse(SqlTemplate)

	if err != nil {
		return nil, err
	}

	args := TmplArgs{
		Package: pkg,
		Ident:   ident,
		Pointer: pointer,
		SqlMap:  sqlMap,
		TmplMap: tmplMap,
	}

	var dst bytes.Buffer
	if err = tmpl.Execute(&dst, &args); err != nil {
		return nil, err
	}

	return dst.Bytes(), nil
}
