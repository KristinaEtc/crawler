package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Crawler struct {
	storeTasks taskStore
	Router     *gin.Engine
}

func (cr *Crawler) PostTasks(c *gin.Context) {
	var bodyRaw []byte
	n, err := c.Request.Body.Read(bodyRaw)
	statusCode := http.StatusBadRequest
	message := "Bad request"

	switch {
	case err != nil:
		log.Printf("[postTasks] error body processing: %s\n", err.Error())
	case n <= 0:
		log.Printf("[postTasks] bad boby length: %d\n", n)
	default:
		defer c.Request.Body.Close()
		createdTask := cr.storeTasks.addNewTask()
		//log.Printf("[DEBUG] got request body=[%s] from ID=[%d]\n", bodyRaw, createdTask.id)
		statusCode = http.StatusOK
		go createdTask.process(byteToSlice(bodyRaw))
	}

	c.String(statusCode, message)
}
