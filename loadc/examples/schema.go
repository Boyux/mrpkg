package main

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

//go:generate go run "github.com/Boyux/mrpkg/loadc" --mode=sqlx --output=user_service.go
type UserService interface {
	// Get QUERY
	// include sql/get_user.sql
	Get(id int64) (User, error)

	// QueryByNames QUERY
	// SELECT
	//     id,
	//     name
	// FROM user
	// WHERE
	//     name IN ({{ bindvars $.names }})
	QueryByNames(names []string) ([]User, error)

	// Update EXEC
	// UPDATE user SET name = ? WHERE id = ?;
	Update(user *UserUpdate) error

	// UpdateName EXEC
	// UPDATE user SET name = ? WHERE id = ?;
	UpdateName(id int64, name string) error
}

func main() {
	userService := NewUserService("driver", "source")
	user, _ := userService.Get(1)
	userService.Update(user.Update())
	userService.UpdateName(user.Id, user.Name)
}
