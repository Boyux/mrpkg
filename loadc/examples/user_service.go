// Code generated by loadc, DO NOT EDIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Boyux/mrpkg"
	"io"
	"net/http"
	"net/textproto"
	"text/template"
)

func NewUserService(inner *Inner) UserService {
	return implUserService{inner}
}

type implUserService struct {
	inner *Inner
}

func (imp implUserService) Inner() *Inner {
	return imp.inner
}

func (implUserService) Response() *UserResponse {
	return new(UserResponse)
}

func (imp implUserService) GetUser(id int64) (User, error) {
	var innerGetUser any = imp.inner

	if cacheGetUser, okGetUser := innerGetUser.(interface {
		GetCache(string, ...any) []any
	}); okGetUser {
		if cacheValuesGetUser := cacheGetUser.GetCache("GetUser", id); cacheValuesGetUser != nil {
			return cacheValuesGetUser[0].(User), nil
		}
	}

	var (
		addrTmplGetUser = template.New("AddressGetUser")
	)

	addrGetUser := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(addrGetUser)
	defer addrGetUser.Reset()

	responseBodyGetUser := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(responseBodyGetUser)
	defer responseBodyGetUser.Reset()

	var (
		v0GetUser       User
		errGetUser      error
		responseGetUser interface {
			Err() error
			ScanValues(...any) error
			FromBytes(string, []byte) error
			Break() bool
		} = imp.Response()
	)

	if errGetUser = template.Must(addrTmplGetUser.Parse("{{ $.UserService.Host }}/user/{{ $.id }}")).
		Execute(addrGetUser, map[string]any{
			"UserService": imp.inner,
			"id":          id,
		}); errGetUser != nil {
		return v0GetUser, fmt.Errorf("error building 'GetUser' url: %w", errGetUser)
	}

	urlGetUser := addrGetUser.String()
	requestGetUser, errGetUser := http.NewRequest("GET", urlGetUser, http.NoBody)
	if errGetUser != nil {
		return v0GetUser, fmt.Errorf("error building 'GetUser' request: %w", errGetUser)
	}

	httpResponseGetUser, errGetUser := http.DefaultClient.Do(requestGetUser)
	if errGetUser != nil {
		return v0GetUser, fmt.Errorf("error sending 'GetUser' request: %w", errGetUser)
	}

	if _, errGetUser = io.Copy(responseBodyGetUser, httpResponseGetUser.Body); errGetUser != nil {
		httpResponseGetUser.Body.Close()
		return v0GetUser, fmt.Errorf("error copying 'GetUser' response body: %w", errGetUser)
	} else {
		httpResponseGetUser.Body.Close()
	}

	if httpResponseGetUser.StatusCode < 200 || httpResponseGetUser.StatusCode > 299 {
		return v0GetUser, fmt.Errorf("response status code %d for 'GetUser' with body: \n\n%s\n\n", httpResponseGetUser.StatusCode, responseBodyGetUser.String())
	}

	if errGetUser = responseGetUser.FromBytes("GetUser", responseBodyGetUser.Bytes()); errGetUser != nil {
		return v0GetUser, fmt.Errorf("error converting 'GetUser' response: %w", errGetUser)
	}

	if errGetUser = responseGetUser.Err(); errGetUser != nil {
		return v0GetUser, fmt.Errorf("error returned from 'GetUser' response: %w", errGetUser)
	}

	if errGetUser = responseGetUser.ScanValues(&v0GetUser); errGetUser != nil {
		return v0GetUser, fmt.Errorf("error scanning value from 'GetUser' response: %w", errGetUser)
	}

	if cacheGetUser, okGetUser := innerGetUser.(interface {
		SetCache(string, []any, ...any)
	}); okGetUser {
		cacheGetUser.SetCache(
			"GetUser",
			[]any{id},
			v0GetUser)
	}

	return v0GetUser, nil
}

func (imp implUserService) GetUsers(ids ...int64) ([]User, error) {
	var innerGetUsers any = imp.inner

	if cacheGetUsers, okGetUsers := innerGetUsers.(interface {
		GetCache(string, ...any) []any
	}); okGetUsers {
		if cacheValuesGetUsers := cacheGetUsers.GetCache("GetUsers", ids); cacheValuesGetUsers != nil {
			return cacheValuesGetUsers[0].([]User), nil
		}
	}

	var (
		valuesGetUsers = make([]User, 0, 10)
		nGetUsers      = 0
		pageGetUsers   = func() int {
			current := nGetUsers
			nGetUsers++
			return current
		}
		addrTmplGetUsers = template.New("AddressGetUsers").Funcs(template.FuncMap{
			"page": pageGetUsers,
		})
	)

	addrGetUsers := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(addrGetUsers)
	defer addrGetUsers.Reset()

	responseBodyGetUsers := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(responseBodyGetUsers)
	defer responseBodyGetUsers.Reset()

loop:
	for {
		var (
			v0GetUsers       []User
			errGetUsers      error
			responseGetUsers interface {
				Err() error
				ScanValues(...any) error
				FromBytes(string, []byte) error
				Break() bool
			} = imp.Response()
		)

		if errGetUsers = template.Must(addrTmplGetUsers.Parse("{{ $.UserService.Host }}/users?{{ range $index, $id := $.ids }}{{ if gt $index 0 }}&{{ end }}$id{{ end }}")).
			Execute(addrGetUsers, map[string]any{
				"UserService": imp.inner,
				"ids":         ids,
			}); errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error building 'GetUsers' url: %w", errGetUsers)
		}

		urlGetUsers := addrGetUsers.String()
		requestGetUsers, errGetUsers := http.NewRequest("GET", urlGetUsers, http.NoBody)
		if errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error building 'GetUsers' request: %w", errGetUsers)
		}

		httpResponseGetUsers, errGetUsers := http.DefaultClient.Do(requestGetUsers)
		if errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error sending 'GetUsers' request: %w", errGetUsers)
		}

		if _, errGetUsers = io.Copy(responseBodyGetUsers, httpResponseGetUsers.Body); errGetUsers != nil {
			httpResponseGetUsers.Body.Close()
			return v0GetUsers, fmt.Errorf("error copying 'GetUsers' response body: %w", errGetUsers)
		} else {
			httpResponseGetUsers.Body.Close()
		}

		if httpResponseGetUsers.StatusCode < 200 || httpResponseGetUsers.StatusCode > 299 {
			return v0GetUsers, fmt.Errorf("response status code %d for 'GetUsers' with body: \n\n%s\n\n", httpResponseGetUsers.StatusCode, responseBodyGetUsers.String())
		}

		if errGetUsers = responseGetUsers.FromBytes("GetUsers", responseBodyGetUsers.Bytes()); errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error converting 'GetUsers' response: %w", errGetUsers)
		}

		if errGetUsers = responseGetUsers.Err(); errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error returned from 'GetUsers' response: %w", errGetUsers)
		}

		if errGetUsers = responseGetUsers.ScanValues(&v0GetUsers); errGetUsers != nil {
			return v0GetUsers, fmt.Errorf("error scanning value from 'GetUsers' response: %w", errGetUsers)
		}

		valuesGetUsers = append(valuesGetUsers, v0GetUsers...)
		if responseGetUsers.Break() {
			break loop
		}
	}

	if cacheGetUsers, okGetUsers := innerGetUsers.(interface {
		SetCache(string, []any, ...any)
	}); okGetUsers {
		cacheGetUsers.SetCache(
			"GetUsers",
			[]any{ids},
			valuesGetUsers)
	}

	return valuesGetUsers, nil
}

func (imp implUserService) UpdateUser(user *User) error {
	var innerUpdateUser any = imp.inner

	if cacheUpdateUser, okUpdateUser := innerUpdateUser.(interface {
		GetCache(string, ...any) []any
	}); okUpdateUser {
		if cacheValuesUpdateUser := cacheUpdateUser.GetCache("UpdateUser", user); cacheValuesUpdateUser != nil {
			return nil
		}
	}

	var (
		addrTmplUpdateUser   = template.New("AddressUpdateUser")
		headerTmplUpdateUser = template.New("HeaderUpdateUser")
	)

	addrUpdateUser := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(addrUpdateUser)
	defer addrUpdateUser.Reset()

	headerUpdateUser := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(headerUpdateUser)
	defer headerUpdateUser.Reset()

	responseBodyUpdateUser := mrpkg.GetObj[*bytes.Buffer]()
	defer mrpkg.PutObj(responseBodyUpdateUser)
	defer responseBodyUpdateUser.Reset()

	var (
		errUpdateUser      error
		responseUpdateUser interface {
			Err() error
			ScanValues(...any) error
			FromBytes(string, []byte) error
			Break() bool
		} = imp.Response()
	)

	if errUpdateUser = template.Must(addrTmplUpdateUser.Parse("{{ $.UserService.Host }}/user")).
		Execute(addrUpdateUser, map[string]any{
			"UserService": imp.inner,
			"user":        user,
		}); errUpdateUser != nil {
		return fmt.Errorf("error building 'UpdateUser' url: %w", errUpdateUser)
	}

	if errUpdateUser = template.Must(headerTmplUpdateUser.Parse("Content-Type: application/json\r\n\r\n{\r\n\"id\": {{ $.user.Id }},\r\n\"name\": {{ $.user.Name }}\r\n}\r\n\r\n")).
		Execute(headerUpdateUser, map[string]any{
			"UserService": imp.inner,
			"user":        user,
		}); errUpdateUser != nil {
		return fmt.Errorf("error building 'UpdateUser' header: %w", errUpdateUser)
	}
	bufReaderUpdateUser := bufio.NewReader(headerUpdateUser)
	mimeHeaderUpdateUser, errUpdateUser := textproto.NewReader(bufReaderUpdateUser).ReadMIMEHeader()
	if errUpdateUser != nil {
		return fmt.Errorf("error reading 'UpdateUser' header: %w", errUpdateUser)
	}

	urlUpdateUser := addrUpdateUser.String()
	requestUpdateUser, errUpdateUser := http.NewRequest("PUT", urlUpdateUser, bufReaderUpdateUser)
	if errUpdateUser != nil {
		return fmt.Errorf("error building 'UpdateUser' request: %w", errUpdateUser)
	}

	for kUpdateUser, vvUpdateUser := range mimeHeaderUpdateUser {
		for _, vUpdateUser := range vvUpdateUser {
			requestUpdateUser.Header.Add(kUpdateUser, vUpdateUser)
		}
	}

	httpResponseUpdateUser, errUpdateUser := http.DefaultClient.Do(requestUpdateUser)
	if errUpdateUser != nil {
		return fmt.Errorf("error sending 'UpdateUser' request: %w", errUpdateUser)
	}

	if _, errUpdateUser = io.Copy(responseBodyUpdateUser, httpResponseUpdateUser.Body); errUpdateUser != nil {
		httpResponseUpdateUser.Body.Close()
		return fmt.Errorf("error copying 'UpdateUser' response body: %w", errUpdateUser)
	} else {
		httpResponseUpdateUser.Body.Close()
	}

	if httpResponseUpdateUser.StatusCode < 200 || httpResponseUpdateUser.StatusCode > 299 {
		return fmt.Errorf("response status code %d for 'UpdateUser' with body: \n\n%s\n\n", httpResponseUpdateUser.StatusCode, responseBodyUpdateUser.String())
	}

	if errUpdateUser = responseUpdateUser.FromBytes("UpdateUser", responseBodyUpdateUser.Bytes()); errUpdateUser != nil {
		return fmt.Errorf("error converting 'UpdateUser' response: %w", errUpdateUser)
	}

	if errUpdateUser = responseUpdateUser.Err(); errUpdateUser != nil {
		return fmt.Errorf("error returned from 'UpdateUser' response: %w", errUpdateUser)
	}

	if errUpdateUser = responseUpdateUser.ScanValues(); errUpdateUser != nil {
		return fmt.Errorf("error scanning value from 'UpdateUser' response: %w", errUpdateUser)
	}

	if cacheUpdateUser, okUpdateUser := innerUpdateUser.(interface {
		SetCache(string, []any, ...any)
	}); okUpdateUser {
		cacheUpdateUser.SetCache(
			"UpdateUser",
			[]any{user},
		)
	}

	return nil
}
