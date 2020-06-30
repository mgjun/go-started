package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"mvc-login/datasource"
	"mvc-login/repositories"
	"mvc-login/services"
	"mvc-login/web/controllers"
	"mvc-login/web/middleware"
	"time"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.HandleDir("/public", "./web/public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().
			GetStringDefault("message", "The page your are looking for does not exist"))
		ctx.View("shared/error.html")
	})

	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}

	repo := repositories.NewUserRepository(db)
	userService := services.NewUserService(repo)

	users := mvc.New(app.Party("/users"))

	users.Router.Use(middleware.BasicAuth)
	users.Register(userService)
	users.Handle(new(controllers.UsersController))

	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})

	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
