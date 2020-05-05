package prepare

import "github.com/kataras/iris/v12"

func Service(router Router, address string) error {
	app := iris.New()
	router.Init(app)
	if err := app.Listen(address, iris.WithoutServerError(iris.ErrServerClosed), iris.WithConfiguration(iris.Configuration{
		DisableInterruptHandler:   false,
		DisablePathCorrection:     false,
		EnablePathEscape:          false,
		FireMethodNotAllowed:      false,
		DisableAutoFireStatusCode: false,
		TimeFormat:                "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                   "UTF-8",
	})); err != nil {
		return err
	}
	go APICheck()
	return nil
}
