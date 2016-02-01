package main

import (
	"github.com/astaxie/beego"
	_ "github.com/beegoexample/filters"
	_ "github.com/beegoexample/routers"
)

func main() {
	beego.Run()
}
