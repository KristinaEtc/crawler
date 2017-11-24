package main

import (
	"sync"
)

type taskStore struct {
	tasks map[int]task
	m     *sync.RWMutex
}

func (tStore taskStore) addNewTask() task {
	tStore.m.Lock()
	defer tStore.m.Unlock()

	id := len(tStore.tasks)
	newTask := task{
		id:    id,
		state: stateProcessing,
	}
	tStore.tasks[id] = newTask
	return newTask
}

const (
	stateProcessing = 0
	stateDone       = iota
)

type task struct {
	data []byte
	id   int

	state int
	m     *sync.RWMutex
}

func (t *task) process(uris []string) {
	crawlerResCh := make(chan []byte, len(uris))

	for _, uri := range uris {
		go getMeta(crawlerResCh, uri)
	}

	for i := 0; i < len(uris); i++ {
		select {
		case res := <-crawlerResCh:
			t.data = append(t.data, res...)
		}
	}

	t.m.Lock()
	defer t.m.Unlock()
	t.state = stateDone
}
