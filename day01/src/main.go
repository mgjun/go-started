package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func main() {
	app := iris.New()

	app.Get("/get", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(path)

		_, err := context.WriteString("receive get request")
		if err != nil {
			app.Logger().Error("response error")
		}
	})
	
	app.Post("post", func(ctx context.Context) {
		path := ctx.Path()
		app.Logger().Info(path)

		_, err := ctx.WriteString("{'post':'body'}")
		if err != nil{
			app.Logger().Error("response error")
		}
	})

	if app.Run(iris.Addr(":8080")) != nil {
		app.Logger().Error("start application failed")
	}
}
