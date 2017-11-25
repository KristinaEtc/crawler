package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	answerBagRequst      = "Bad request"
	answerTaskProcessing = "Task processing"
	answerTaskOK         = "Task finished process"
	answerTaskNotFound   = "Task not found"

	// do not delete by default
	deleteAfterSendingTask = false
)

type Crawler struct {
	storeTasks taskStore
	Router     *gin.Engine
}

func InitCrawler() *Crawler {
	c := &Crawler{
		Router: gin.Default(),
		storeTasks: taskStore{
			tasks: make(map[int]*task, 0),
			m:     &sync.RWMutex{},
		},
	}

	c.Router.POST("/tasks/", c.PostTask)
	c.Router.GET("/tasks/:taskID/", c.GetTask)
	return c
}

func (cr *Crawler) PostTask(c *gin.Context) {
	bodyRaw, err := ioutil.ReadAll(c.Request.Body)
	statusCode := http.StatusBadRequest
	message := answerBagRequst

	switch {
	case err != nil:
		log.Printf("[postTasks] error body processing: %s\n", err.Error())
	case len(bodyRaw) <= 0:
		log.Printf("[postTasks] bad boby length: %d\n", len(bodyRaw))
	default:

		createdTask := cr.storeTasks.addNewTask(getURIsFromBodyRequest(bodyRaw))
		id := createdTask.getTaskID()
		log.Printf("[DEBUG] got request\tID=%d\n", id)
		statusCode = http.StatusOK
		message = fmt.Sprintf("taskId:%d", id)

		go createdTask.process()
	}

	c.String(statusCode, message)
}

func (cr *Crawler) GetTask(c *gin.Context) {
	log.Printf("tasks=%+v\n", cr.storeTasks)
	paramTaskID := c.Param("taskID")

	taskID, err := strconv.Atoi(paramTaskID)
	if err != nil {
		log.Printf("[ERR] could not parse task id %s\n", paramTaskID)
		c.String(http.StatusBadRequest, answerBagRequst)
		return
	}

	paramDelete := c.DefaultQuery("delete", "0")
	deleteTask, err := strconv.ParseBool(paramDelete)
	if err != nil {
		log.Printf("[ERR] wrong delete param %s\tid=%d\n", paramDelete, taskID)
		deleteTask = deleteAfterSendingTask
	}

	log.Printf("got request\tID=%d\n", taskID)

	var statusCode int
	var message string
	_, state := cr.storeTasks.GetTaskContent(taskID, false)
	log.Println("state=", state)
	switch {
	case state == stateNotExist:
		log.Println("[DEBUG] case stateNotExist")
		statusCode = http.StatusNotFound
		message = answerTaskNotFound
	case state == stateProcessing:
		log.Println("[DEBUG] case stateProcessing")
		statusCode = http.StatusNoContent
		message = answerTaskProcessing
	case state == stateDone:
		log.Println("[DEBUG] case stateDone")
		statusCode = http.StatusOK
		resData, _ := cr.storeTasks.GetTaskContent(taskID, deleteTask)
		log.Printf("[DEBUG] resData = [%s]\n", resData.String())

		c.Data(http.StatusOK, "text/csv", resData.Bytes())
		return
	}
	c.String(statusCode, message)
}
