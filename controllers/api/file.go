package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/beegoexample/models/oss"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

type FileController struct {
	beego.Controller
}

// @router /api/file/upload [post]
func (c *FileController) Upload() {

	result := map[string]interface{}{}

	defer func() {
		c.Data["json"] = &result
		c.ServeJSON()
	}()

	en := context.ParseEncoding(c.Ctx.Request)
	if en != "" {
		c.Ctx.ResponseWriter.Header().Set("Content-Encoding", en)
	}
	// todo: ?既然c.ServeJSON 已经明确设置该头了，为什么没有生效？
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	f, fh, err := c.GetFile("file")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		result["retcode"] = 4000001
		result["msg"] = "no file to upload"
		return
	}
	defer f.Close()
	//ReadForm 返回的时候回删除可能使用到的临时文件的，GetFile返回的是文件的句柄，所以不需要显式删除临时文件，文件小于10MB不会用到临时文件

	bucket, err := oss.NewBucket()
	if err != nil {
		msg := err.Error()
		http.Error(c.Ctx.ResponseWriter, msg, 500)
		result["retcode"] = 5000001
		result["msg"] = msg
		return
	}

	err = bucket.PutObject(fh.Filename, f)
	if err != nil {
		msg := "upload file fail: " + err.Error()
		http.Error(c.Ctx.ResponseWriter, msg, 500)
		result["retcode"] = 5000002
		result["msg"] = msg
		return
	}

	endpoint := beego.AppConfig.String("oss.endpoint")
	url := "http://" + bucket.BucketName + endpoint + "/" + fh.Filename
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	result["retcode"] = 2000000
	result["msg"] = "created"
	result["body"] = map[string]string{"path": url}

}

// @router /api/file/download [get]
func (c *FileController) Download() {

	bucket, err := oss.NewBucket()
	if err != nil {
		msg := err.Error()
		http.Error(c.Ctx.ResponseWriter, msg, 500)
		return
	}

	filename := c.GetString("name")
	if filename == "" {
		http.Error(c.Ctx.ResponseWriter, "need param name", 400)
		return
	}

	//	object, err := url.ParseRequestURI(file)
	//	if err != nil {
	//		http.Error(c.Ctx.ResponseWriter, "param name invalid", 400)
	//		return
	//	}
	//	filename := filepath.Base(file)
	r, err := bucket.GetObject(filename)
	if err != nil {
		beego.Debug(err)
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

	// todo: 这里需要优化
	content, err := ioutil.ReadAll(r)
	if err != nil {
		beego.Debug(err)
		http.Error(c.Ctx.ResponseWriter, "very said", 500)
		return
	}
	// todo: 这里应该提供一个支持ioReader的方法
	context.WriteBody(en, c.Ctx.ResponseWriter, content)
}

// @router /api/file/list [get]
func (c *FileController) List() {
	bucket, err := oss.NewBucket()
	if err != nil {
		msg := err.Error()
		http.Error(c.Ctx.ResponseWriter, msg, 500)
		return
	}

	lsRes, err := bucket.ListObjects()
	if err != nil {
		msg := err.Error()
		http.Error(c.Ctx.ResponseWriter, msg, 500)
		return
	}

	paths := map[string]string{}

	for _, object := range lsRes.Objects {
		paths[object.Key] = object.Key
	}

	c.Data["paths"] = paths
	c.TplName = "api/file/list.tpl"
}

type slice []string

func (this *slice) has(e string) bool {
	for _, item := range *this {
		if e == item {
			return true
		}
	}
	return false
}
