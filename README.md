# httpin

[![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin)

HTTP Input for Go

## Quick View

<table>
<tr>
  <th>BEFORE (use functions provided by net/http lib)</th>
  <th>AFTER (use httpin to extract input parameters from HTTP request)</th>
</tr>
<tr>
  <td>

```go
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		// ... invalid page
		return
	}

	perPage, err := strconv.ParseInt(r.FormValue("per_page"), 10, 64)
	if err != nil {
		// ... invalid per_page
		return
	}

	isVip, _ := strconv.ParseBool(r.FormValue("is_vip"))

	// do sth.
}
```

  </td>
  <td>

```go
type ListUsersInput struct {
	Page    int  `in:"form=page"`
	PerPage int  `in:"form=per_page"`
	IsVip   bool `in:"form=is_vip"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	interfaceInput, err := httpin.New(ListUsersInput{}).ReadRequest(r)
	if err != nil {
		// err can be *httpin.InvalidField
		return
	}

	input := interfaceInput.(*ListUsersInput)
	// do sth.
}
```

  </td>
</tr>
</table>

## Features

- [x] Decode from HTTP query, i.e. `http.Request.Form`
- [x] Decode from HTTP headers, e.g. `http.Request.Header`
- [x] Builtin decoders for basic types, e.g. `bool`, `int`, `int64`, `float32`, `time.Time`, ... [full list](./decoders.go)
- [x] Decode one field by inspecting multiple keys one by one in the same source
- [x] Decode one field from multiple sources, e.g. both query and headers
- [ ] Customize decoders for user defined types
- [x] Define input struct with embedded struct fields
- [x] Tag one field as **required**
- [ ] Builtin encoders for basic types
- [ ] Customize encoders for user defined types

## Sample User Defined Input Structs

```go
type Authorization struct {
	// Decode from multiple sources, the former with higher priority
	Token string `in:"form=access_token,header=x-api-token"`
}

type Pagination struct {
	Page int `in:"form=page"`

	// Decode from multiple keys in the same source, the former with higher priority
	PerPage int `in:"form=per_page,page_size"`
}

type ListUsersInput struct {
	Gender   string `in:"form=gender"`
	AgeRange []int  `in:"form=age_range"`
	IsMember bool   `in:"form=is_member"`

	Pagination    // Embedded field works
	Authorization // Embedded field works
}
```

## Advanced - Use Middleware

First, set up the middleware for your handlers.

```go
func init() {
	mux.Handle("/users", alice.New(
		httpin.NewInput(ListUsersInput{}),
	).ThenFunc(ListUsers)).Methods("GET")
}
```

Second, fetch your input with **only one line** of code.

```go
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*UserQuery)
	// do sth.
}
```
