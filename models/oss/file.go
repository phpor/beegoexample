package oss
import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"fmt"
)

func NewBucket() (bucket *oss.Bucket, err error) {
	endpoint := beego.AppConfig.String("oss.endpoint")
	accesskeyid := beego.AppConfig.String("oss.accesskeyid")
	accesskeysecret := beego.AppConfig.String("oss.accesskeysecret")

	client, err := oss.New(endpoint, accesskeyid, accesskeysecret)
	if err != nil {
		err = fmt.Errorf("init storage fail: " + err.Error())
		return
	}
	bucketname := beego.AppConfig.String("oss.bucket")

	bucket, err = client.Bucket(bucketname)
	if err != nil {
		err = fmt.Errorf("init bucket fail: " + err.Error())
	}
	return
}