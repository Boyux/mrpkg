package main

import (
	"context"
	"database/sql"
)

type User struct {
	Id   int64
	Name string
}

func (u *User) Update() *UserUpdate {
	return &UserUpdate{
		Name: u.Name,
		Id:   u.Id,
	}
}

type UserUpdate struct {
	Name string
	Id   int64
}

func (update *UserUpdate) ToArgs() []any {
	return []any{
		update.Name,
		update.Id,
	}
}

//go:generate go run "github.com/Boyux/mrpkg/loadc" --mode=sqlx --features=sqlx/log,sqlx/rebind --output=user_handler.go
type UserHandler interface {
	WithTx(context.Context, func(UserHandler) error) error

	// Get QUERY
	// include sql/get_user.sql
	Get(ctx context.Context, id int64) (*User, error)

	// QueryByName QUERY NAMED
	// SELECT
	//     id,
	//     name
	// FROM user
	// WHERE
	//     name = :name
	QueryByName(name string) ([]User, error)

	// Update EXEC
	// UPDATE user SET name = ? WHERE id = ?;
	Update(ctx context.Context, user *UserUpdate) error

	// UpdateName EXEC NAMED
	// UPDATE user SET name = :name WHERE id = :id;
	UpdateName(ctx context.Context, id int64, name string) (sql.Result, error)
}

type Inner struct {
	Host string
}

type UserResponse struct{}

func (*UserResponse) Err() error                     { panic("unimplemented") }
func (*UserResponse) ScanValues(...any) error        { panic("unimplemented") }
func (*UserResponse) FromBytes(string, []byte) error { panic("unimplemented") }
func (*UserResponse) Break() bool                    { panic("unimplemented") }

//go:generate go run "github.com/Boyux/mrpkg/loadc" --mode=api --features=api/cache,api/log,api/client --output=user_service.go
type UserService interface {
	Inner() *Inner
	Response() *UserResponse

	// GetUser GET {{ $.UserService.Host }}/user/{{ $.id }}
	GetUser(ctx context.Context, id int64) (*User, error)

	// GetUsers GET {{ $.UserService.Host }}/users?{{ range $index, $id := $.ids }}{{ if gt $index 0 }}&{{ end }}{{ $id }}{{ end }}
	GetUsers(ids ...int64) ([]User, error)

	// UpdateUser PUT {{ $.UserService.Host }}/user
	// Content-Type: application/json
	//
	// {
	//     "id": {{ $.user.Id }},
	//     "name": {{ $.user.Name }}
	// }
	UpdateUser(user *User) error
}

func main() {
	var user *User

	userHandler := NewUserHandler("driver", "source")
	user, _ = userHandler.Get(context.Background(), 1)
	userHandler.Update(context.Background(), user.Update())
	userHandler.UpdateName(context.Background(), user.Id, user.Name)

	userService := NewUserService(&Inner{Host: "host"})
	user, _ = userService.GetUser(context.Background(), 1)
	userService.UpdateUser(user)
}
