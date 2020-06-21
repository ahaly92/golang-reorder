package main

import (
	. "github.com/ahaly92/golang-reorder/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()
	ginEngine.GET("/todos", Todos)

	ginEngine.Run(":4000")
}
