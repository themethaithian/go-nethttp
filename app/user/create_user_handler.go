package user

import "github.com/themethaithian/go-nethttp/app"

type CreateUser struct {
	Username   string  `json:"username"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
}

func (h *handler) CreateUser(ctx app.Context) {
	var user CreateUser
	if err := ctx.Bind(&user); err != nil {
		ctx.BadRequest(err)
		return
	}

	ctx.OK("success fully create user!")
}
