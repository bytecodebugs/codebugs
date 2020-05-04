package prepare

import "github.com/kataras/iris/v12"

type Router interface {
	Init(app *iris.Application)
}
