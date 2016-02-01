package api

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"os"
	"net/http"
	"github.com/astaxie/beego/context"
	"mime"
	"strings"
)

type TestController struct {
	beego.Controller
}

var fileDir string = beego.AppConfig.DefaultString("tmpdir", "/tmp/")
var allowExt = slice(beego.AppConfig.DefaultStrings("allowext", slice{".txt", ".jpg"}))

// @router /test/file/upload [post]
func (c *TestController) Upload() {
	f, fh, err := c.GetFile("file")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	f.Close()
	c.SaveToFile("file", fileDir + fh.Filename)
	c.Redirect("/test/file/list", 302)
}


// @router /test/file/download [get]
func (c *TestController) Download() {
	filename := c.GetString("name")
	if filename == "" {
		http.Error(c.Ctx.ResponseWriter, "Not Found", 404)
		return
	}
	filename = fileDir + filename
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		http.Error(c.Ctx.ResponseWriter, "Not Found", 404)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(filename))
	if ctype == "" {
		ctype = "application/oct-stream"
	}
	en := context.ParseEncoding(c.Ctx.Request)
	if en != "" {
		c.Ctx.ResponseWriter.Header().Set("Content-Encoding", en)
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", ctype)

	f, err := os.Open(filename)
	context.WriteFile(en, c.Ctx.ResponseWriter, f)
}

// @router /test/file/list [get]
func (c *TestController) List() {
	paths := map[string]string{}
	filepath.Walk(fileDir, func(path string, finfo os.FileInfo, err error) error{
		if err != nil {
			return err
		}
		if finfo.IsDir() {
			return nil
		}
		if !allowExt.has(filepath.Ext(path)) {
			return nil
		}
		paths[filepath.Base(path)] = strings.TrimPrefix(path, fileDir)
		return  nil
	})
	c.Data["paths"] = paths
	c.TplName = "file/list.tpl"
}

type slice []string

func (this *slice) has(e string) bool{
	for _, item := range *this {
		if e == item {
			return true
		}
	}
	return false
}