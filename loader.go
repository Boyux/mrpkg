package mrpkg

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type Data map[string]any

// SqlLoader load sql/tmpl
type SqlLoader struct {
	sqlMap  ConcurrentMap[string, string]
	tmplMap ConcurrentMap[string, *template.Template]
}

func (loader *SqlLoader) AddSql(id string, sql string) {
	loader.sqlMap.Set(id, sql)
}

func (loader *SqlLoader) AddTmpl(id string, tmpl *template.Template) {
	loader.tmplMap.Set(id, tmpl)
}

func (loader *SqlLoader) LoadSql(id string) string {
	sql, err := loader.LoadSqlWithErr(id)
	if err != nil {
		panic(err)
	}
	return sql
}

func (loader *SqlLoader) LoadSqlWithErr(id string) (string, error) {
	sql, ok := loader.sqlMap.Get(id)
	if !ok {
		return "", fmt.Errorf("SqlLoader.LoadSqlWithErr: no sql found for id %s", strconv.Quote(id))
	}
	return sql, nil
}

func (loader *SqlLoader) LoadTmpl(id string, data any) string {
	sql, err := loader.LoadTmplWithErr(id, data)
	if err != nil {
		panic(err)
	}
	return sql
}

func (loader *SqlLoader) LoadTmplWithErr(id string, data any) (string, error) {
	tmpl, ok := loader.tmplMap.Get(id)
	if !ok {
		return "", fmt.Errorf("SqlLoader.LoadTmplWithErr: no sql template found for id %s", strconv.Quote(id))
	}

	var dst bytes.Buffer
	if err := tmpl.Execute(&dst, data); err != nil {
		return "", fmt.Errorf("SqlLoader.LoadTmplWithErr: %w", err)
	}

	return dst.String(), nil
}

var byteType = reflect.TypeOf([]byte{})

type ToArgs interface {
	ToArgs() []any
}

func MergeArgs(args ...any) []any {
	dst := make([]any, 0, len(args))
	for _, arg := range args {
		rv := reflect.ValueOf(arg)
		if toArgs, ok := arg.(ToArgs); ok {
			dst = append(dst, MergeArgs(toArgs.ToArgs()...)...)
		} else if rv.Kind() == reflect.Slice && rv.Type() != byteType {
			for i := 0; i < rv.Len(); i++ {
				dst = append(dst, rv.Index(i).Interface())
			}
		} else {
			dst = append(dst, arg)
		}
	}
	return dst
}

func GenBindVars(data any) string {
	var n int
	switch rv := reflect.ValueOf(data); rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n = int(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = int(rv.Uint())
	case reflect.Slice:
		if rv.Type() == byteType {
			n = 1
		} else {
			n = rv.Len()
		}
	}

	bindVars := make([]string, n)
	for i := 0; i < n; i++ {
		bindVars[i] = "?"
	}

	return strings.Join(bindVars, ", ")
}
