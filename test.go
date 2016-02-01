package main
import (
	"path/filepath"
	"os"
)

func main() {
	filename := "/tmp/a.txt"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		println("aaaaaa")
		return
	}
	println("exist")
}

func testFilewalk() {


	filepath.Walk("/tmp/", func(path string, info os.FileInfo, err error) error{
		if err != nil {
			println(err)
		}
		println(path)
		return nil
	})
}