package filters
import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
)


func FilterApiAuth(ctx *context.Context) {

}

func init() {
	beego.InsertFilter("/api/file/*", beego.BeforeRouter, FilterApiAuth)
}
