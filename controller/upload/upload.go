package upload

import (
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const BaseURL = "/api/upload"

type Controller struct {}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e echo.Echo) (err error) {
	e.POST(BaseURL, upload)
	return nil
}

func upload(c echo.Context) (err error) {
	fmt.Println(c)
	fmt.Println(c.FormFile("file"))
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	t := time.Now()
	nameList := strings.Split(file.Filename, ".")
	if len(nameList) != 2 {
		return c.String(http.StatusOK, "请检查文件名是否有非法字符，上传失败")
	} else {
		fileName := string(nameList[0] + "-" + t.Format("20060102150405") + "." + nameList[1])
		fmt.Println(fileName)
		dst, err := os.Create("/root/go/src/BlogBackend/upload/" + fileName)
		fmt.Println(dst)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		return c.String(http.StatusOK, fileName)
	}

}
