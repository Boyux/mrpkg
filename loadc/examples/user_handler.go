// Code generated by loadc, DO NOT EDIT

package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/Boyux/mrpkg"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func NewUserHandler(drv string, dsn string) UserHandler {
	return &implUserHandler{
		Core: sqlx.MustOpen(drv, dsn),
	}
}

func NewUserHandlerFromDB(core *sqlx.DB) UserHandler {
	return &implUserHandler{
		Core: core,
	}
}

func NewUserHandlerFromCore(core interface {
	Rebind(query string) string
	Beginx() (*sqlx.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}) UserHandler {
	return &implUserHandler{
		Core: core,
	}
}

type implUserHandler struct {
	withTx bool
	Core   interface {
		Rebind(query string) string
		Beginx() (*sqlx.Tx, error)
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
		PrepareNamed(query string) (*sqlx.NamedStmt, error)
		PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
		Exec(query string, args ...interface{}) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		Get(dest interface{}, query string, args ...interface{}) error
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		Select(dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	}
}

func (imp *implUserHandler) Get(ctx context.Context, id int64) (*User, error) {
	var (
		v0Get  = new(User)
		errGet error
	)

	sqlTmplGet := template.Must(
		template.
			New("Get").
			Funcs(template.FuncMap{
				"bindvars": mrpkg.GenBindVars,
			}).
			Parse("SELECT *\nFROM user\nWHERE id = ?;\r\n\r\n"),
	)

	sqlGet := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(sqlGet)
	defer sqlGet.Reset()

	if errGet = sqlTmplGet.Execute(sqlGet, map[string]any{
		"ctx": ctx,
		"id":  id,
	}); errGet != nil {
		return v0Get, fmt.Errorf("error executing %s template: %w", strconv.Quote("Get"), errGet)
	}

	sqlQueryGet := strings.TrimSpace(sqlGet.String())
	sqlQueryGet = imp.Core.Rebind(sqlQueryGet)

	argsGet := mrpkg.MergeArgs(
		id,
	)

	startGet := time.Now()

	errGet = imp.Core.GetContext(ctx, v0Get, sqlQueryGet, argsGet...)

	if logGet, okGet := imp.Core.(interface {
		Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
	}); okGet {
		logGet.Log(ctx, "Get", sqlQueryGet, argsGet, time.Since(startGet))
	}

	if errGet != nil {
		return v0Get, fmt.Errorf("error executing %s sql: \n\n%s\n\n%w", strconv.Quote("Get"), sqlQueryGet, errGet)
	}

	return v0Get, nil
}

func (imp *implUserHandler) QueryByName(name string) ([]User, error) {
	var (
		v0QueryByName  []User
		errQueryByName error
	)

	sqlTmplQueryByName := template.Must(
		template.
			New("QueryByName").
			Funcs(template.FuncMap{
				"bindvars": mrpkg.GenBindVars,
			}).
			Parse("SELECT\r\nid,\r\nname\r\nFROM user\r\nWHERE\r\nname = :name\r\n\r\n"),
	)

	sqlQueryByName := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(sqlQueryByName)
	defer sqlQueryByName.Reset()

	if errQueryByName = sqlTmplQueryByName.Execute(sqlQueryByName, map[string]any{
		"name": name,
	}); errQueryByName != nil {
		return v0QueryByName, fmt.Errorf("error executing %s template: %w", strconv.Quote("QueryByName"), errQueryByName)
	}

	sqlQueryQueryByName := strings.TrimSpace(sqlQueryByName.String())
	sqlQueryQueryByName = imp.Core.Rebind(sqlQueryQueryByName)

	argsQueryByName := mrpkg.MergeNamedArgs(map[string]any{
		"name": name,
	})

	startQueryByName := time.Now()

	stmtQueryByName, errQueryByName := imp.Core.PrepareNamed(sqlQueryQueryByName)
	if errQueryByName != nil {
		return v0QueryByName, fmt.Errorf("error creating %s prepare statement: %w", strconv.Quote("QueryByName"), errQueryByName)
	}
	errQueryByName = stmtQueryByName.Select(&v0QueryByName, argsQueryByName)

	if logQueryByName, okQueryByName := imp.Core.(interface {
		Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
	}); okQueryByName {
		logQueryByName.Log(context.Background(), "QueryByName", sqlQueryQueryByName, argsQueryByName, time.Since(startQueryByName))
	}

	if errQueryByName != nil {
		return v0QueryByName, fmt.Errorf("error executing %s sql: \n\n%s\n\n%w", strconv.Quote("QueryByName"), sqlQueryQueryByName, errQueryByName)
	}

	return v0QueryByName, nil
}

func (imp *implUserHandler) Update(ctx context.Context, user *UserUpdate) error {
	var (
		errUpdate error
	)

	sqlTmplUpdate := template.Must(
		template.
			New("Update").
			Funcs(template.FuncMap{
				"bindvars": mrpkg.GenBindVars,
			}).
			Parse("UPDATE user SET name = ? WHERE id = ?;\r\n\r\n"),
	)

	sqlUpdate := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(sqlUpdate)
	defer sqlUpdate.Reset()

	if errUpdate = sqlTmplUpdate.Execute(sqlUpdate, map[string]any{
		"ctx":  ctx,
		"user": user,
	}); errUpdate != nil {
		return fmt.Errorf("error executing %s template: %w", strconv.Quote("Update"), errUpdate)
	}

	txUpdate, errUpdate := imp.Core.BeginTxx(ctx, nil)
	if errUpdate != nil {
		return fmt.Errorf("error creating %s transaction: %w", strconv.Quote("Update"), errUpdate)
	}
	if !imp.withTx {
		defer txUpdate.Rollback()
	}

	offsetUpdate := 0
	argsUpdate := mrpkg.MergeArgs(
		user,
	)

	for _, splitSqlUpdate := range strings.Split(sqlUpdate.String(), ";") {
		splitSqlUpdate = strings.TrimSpace(splitSqlUpdate)
		if splitSqlUpdate == "" {
			continue
		}

		countUpdate := strings.Count(splitSqlUpdate, "?")
		splitSqlUpdate = imp.Core.Rebind(splitSqlUpdate)

		startUpdate := time.Now()

		_, errUpdate = txUpdate.ExecContext(ctx, splitSqlUpdate, argsUpdate[offsetUpdate:offsetUpdate+countUpdate]...)

		if logUpdate, okUpdate := imp.Core.(interface {
			Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
		}); okUpdate {
			logUpdate.Log(ctx, "Update", splitSqlUpdate, argsUpdate, time.Since(startUpdate))
		}

		if errUpdate != nil {
			return fmt.Errorf("error executing %s sql: \n\n%s\n\n%w", strconv.Quote("Update"), splitSqlUpdate, errUpdate)
		}

		offsetUpdate += countUpdate
	}

	if !imp.withTx {
		if errUpdate := txUpdate.Commit(); errUpdate != nil {
			return fmt.Errorf("error committing %s transaction: %w", strconv.Quote("Update"), errUpdate)
		}
	}

	return nil
}

func (imp *implUserHandler) UpdateName(ctx context.Context, id int64, name string) (sql.Result, error) {
	var (
		v0UpdateName  sql.Result
		errUpdateName error
	)

	sqlTmplUpdateName := template.Must(
		template.
			New("UpdateName").
			Funcs(template.FuncMap{
				"bindvars": mrpkg.GenBindVars,
			}).
			Parse("UPDATE user SET name = :name WHERE id = :id;\r\n\r\n"),
	)

	sqlUpdateName := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(sqlUpdateName)
	defer sqlUpdateName.Reset()

	if errUpdateName = sqlTmplUpdateName.Execute(sqlUpdateName, map[string]any{
		"ctx":  ctx,
		"id":   id,
		"name": name,
	}); errUpdateName != nil {
		return v0UpdateName, fmt.Errorf("error executing %s template: %w", strconv.Quote("UpdateName"), errUpdateName)
	}

	txUpdateName, errUpdateName := imp.Core.BeginTxx(ctx, nil)
	if errUpdateName != nil {
		return v0UpdateName, fmt.Errorf("error creating %s transaction: %w", strconv.Quote("UpdateName"), errUpdateName)
	}
	if !imp.withTx {
		defer txUpdateName.Rollback()
	}

	argsUpdateName := mrpkg.MergeNamedArgs(map[string]any{
		"id":   id,
		"name": name,
	})

	for _, splitSqlUpdateName := range strings.Split(sqlUpdateName.String(), ";") {
		splitSqlUpdateName = strings.TrimSpace(splitSqlUpdateName)
		if splitSqlUpdateName == "" {
			continue
		}
		splitSqlUpdateName = imp.Core.Rebind(splitSqlUpdateName)

		startUpdateName := time.Now()

		stmtUpdateName, errUpdateName := txUpdateName.PrepareNamedContext(ctx, splitSqlUpdateName)
		if errUpdateName != nil {
			return v0UpdateName, fmt.Errorf("error creating %s prepare statement: %w", strconv.Quote("UpdateName"), errUpdateName)
		}

		v0UpdateName, errUpdateName = stmtUpdateName.ExecContext(ctx, argsUpdateName)

		if logUpdateName, okUpdateName := imp.Core.(interface {
			Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
		}); okUpdateName {
			logUpdateName.Log(ctx, "UpdateName", splitSqlUpdateName, argsUpdateName, time.Since(startUpdateName))
		}

		if errUpdateName != nil {
			return v0UpdateName, fmt.Errorf("error executing %s sql: \n\n%s\n\n%w", strconv.Quote("UpdateName"), splitSqlUpdateName, errUpdateName)
		}

	}

	if !imp.withTx {
		if errUpdateName := txUpdateName.Commit(); errUpdateName != nil {
			return v0UpdateName, fmt.Errorf("error committing %s transaction: %w", strconv.Quote("UpdateName"), errUpdateName)
		}
	}

	return v0UpdateName, nil
}

func NewUserHandlerFromTxAndLog(core *sqlx.Tx, log interface {
	Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
}) UserHandler {
	return &implUserHandler{
		withTx: true,
		Core: &txUserHandler{
			Tx:  core,
			log: log,
		},
	}
}

type txUserHandler struct {
	*sqlx.Tx
	log interface {
		Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
	}
}

func (tx txUserHandler) Beginx() (*sqlx.Tx, error) {
	return tx.Tx, nil
}

func (tx txUserHandler) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return tx.Tx, nil
}

func (tx txUserHandler) Log(ctx context.Context, caller string, query string, args any, elapse time.Duration) {
	if tx.log != nil {
		tx.log.Log(ctx, caller, query, args, elapse)
	}
}

func (imp *implUserHandler) WithTx(ctx context.Context, f func(UserHandler) error) error {
	inner, err := imp.Core.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction in %s: %w", strconv.Quote("WithTx"), err)
	}

	defer inner.Rollback()

	core := &txUserHandler{
		Tx: inner,
	}

	if log, ok := imp.Core.(interface {
		Log(ctx context.Context, caller string, query string, args any, elapse time.Duration)
	}); ok {
		core.log = log
	}

	tx := &implUserHandler{
		withTx: true,
		Core:   core,
	}

	if err = f(tx); err != nil {
		return err
	}

	if err = inner.Commit(); err != nil {
		return fmt.Errorf("error committing transaction in %s: %w", strconv.Quote("WithTx"), err)
	}

	return nil
}
