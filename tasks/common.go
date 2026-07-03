package tasks

import "copilotlens/internal/github"

var taskList []ITask

type ITask interface {
	Run()
	Stop()
}

func Init(cache *github.CacheManager) {
	taskList = append(taskList, NewCacheCleanTask(cache))
	for _, t := range taskList {
		t.Run()
	}
}

func Stop() {
	for _, t := range taskList {
		t.Stop()
	}
}
