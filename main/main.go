package main

import (
	"BlogBackend"
	"BlogBackend/db"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func main() {
	err := db.InitGlobalDB("127.0.0.1", "blog")
	if err != nil{
		log.Panic(err)
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	s := BlogBackend.NewServer(":1323")
	err = s.Init()
	if err!=nil{
		log.Panic(err)
	}
	s.StartServer()
}