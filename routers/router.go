package routers

import (
	"github.com/astaxie/beego"
	"github.com/beegoexample/controllers/api"
)

func init() {
	//    beego.Router("/", &controllers.MainController{})
	beego.Include(&api.FileController{})
}
