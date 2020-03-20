package main

import (
	"BlogBackend"
	"BlogBackend/db"
	"log"
)

func main() {
	err := db.InitGlobalDB("127.0.0.1", "blog")
	if err != nil{
		log.Panic(err)
	}
	s := BlogBackend.NewServer(":1323")
	err = s.Init()
	if err!=nil{
		log.Panic(err)
	}
	s.StartServer()
}