package main

import (
	"ggapi/db"
	"ggapi/routers"
)

func main() {

	db.Init()
    routers.Init()
}