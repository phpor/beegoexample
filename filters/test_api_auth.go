package filters
import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
	"net/http"
)


func TestFilterApiAuth(ctx *context.Context) {
	panic(ctx)
	// todo: 这里添加具体的认证
	authOk := (ctx.Input.Query("auth") != "")
	if ! authOk {
		// 过滤器中设置状态行、输出内容、重定向都会阻止继续后续逻辑，参看router.go
		http.Error(ctx.ResponseWriter, "Forbidden", 403)
	}
}

func init() {
	beego.InsertFilter("/test/file/*", beego.BeforeRouter, TestFilterApiAuth)
}
