package main

import (
	"bytes"
	"encoding/csv"
	"log"
	"sync"
)

type taskStore struct {
	tasks map[int]*task
	m     *sync.RWMutex
}

func (tStore taskStore) addNewTask(uris []string) *task {
	tStore.m.Lock()
	defer tStore.m.Unlock()

	id := len(tStore.tasks)
	newTask := &task{
		id:    id,
		uris:  uris,
		state: stateProcessing,
		m:     &sync.RWMutex{},
	}
	tStore.tasks[id] = newTask
	return newTask
}

func (tStore taskStore) GetTaskContent(taskID int, deleteTask bool) (*bytes.Buffer, int) {
	tStore.m.Lock()

	if _, ok := tStore.tasks[taskID]; !ok {
		tStore.m.Unlock()
		return nil, -1
	}

	task, _ := tStore.tasks[taskID]
	if deleteTask {
		delete(tStore.tasks, taskID)
	}

	tStore.m.Unlock()
	return task.getTaskContent()
}

const (
	stateNotExist   = -1
	stateProcessing = iota
	stateDone
)

type task struct {
	resData *bytes.Buffer
	id      int
	uris    []string

	state int
	m     *sync.RWMutex
}

func (t *task) process() {
	crawlerResCh := make(chan []string, 0)
	resDataProcessing := make([][]string, 0)

	for _, uri := range t.uris {
		go func(uri string) {
			getMeta(crawlerResCh, uri)
		}(uri)
	}

	for i := 0; i < len(t.uris); i++ {
		select {
		case res := <-crawlerResCh:
			if res != nil {
				//t.resData = append(t.resData, res...)
				resDataProcessing = append(resDataProcessing, res)
			}
		}
	}

	t.resData = &bytes.Buffer{}
	w := csv.NewWriter(t.resData)

	if err := w.WriteAll(resDataProcessing); err != nil {
		log.Printf("[ERROR] writing to csv: %s\n", err.Error())
	}

	log.Println("[DEBUG] processing done")
	t.m.Lock()
	defer t.m.Unlock()
	t.state = stateDone
}

func (t *task) getTaskID() int {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.id
}

func (t *task) getTaskContent() (*bytes.Buffer, int) {
	t.m.RLock()
	defer t.m.RUnlock()

	return t.resData, t.state
}
