package main

import "github.com/gin-gonic/gin"

func main() {
	c := &Crawler{
		Router: gin.Default(),
	}
	c.Router.POST("/tasks/", c.PostTasks)

	// listen and serve on 0.0.0.0:8080
	c.Router.Run()
}
