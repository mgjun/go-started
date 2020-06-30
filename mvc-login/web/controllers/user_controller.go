package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"mvc-login/datamodels"
	"mvc-login/services"
)

type UserController struct {
	Ctx iris.Context

	Service services.UserService

	Session *sessions.Session
}

const userIDKey = "UserID"

func (c *UserController) getCurrentUserID() int64 {
	userId := c.Session.GetInt64Default(userIDKey, 0)
	return userId
}

func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *UserController) logout() {
	c.Session.Destroy()
}

var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	return registerStaticView
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	var (
		firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)

	u, err := c.Service.Create(password, datamodels.User{
		Username:  username,
		Firstname: firstname,
	})

	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		Err:  err,
		Path: "/user/me",
	}
}

var loginStaticView = mvc.View{
	Name: "user/login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	return loginStaticView
}

// PostLogin handles POST: http://localhost:8080/user/register.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	u, found := c.Service.GetByUsernameAndPassword(username, password)

	if !found {
		return mvc.Response{
			Path: "/user/register",
		}
	}
	
	c.Session.Set(userIDKey, u.ID)
	
	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	
	u,found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		c.logout()
		return mvc.Response{
			Path: "/user/login",
		}
	}
	
	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User": u,
		},
	}
	
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout()  {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/user/login")
}