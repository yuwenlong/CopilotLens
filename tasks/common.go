package tasks

import "copilotlens/internal/github"

var taskList []ITask

type ITask interface {
	Run()
	Stop()
}

func Init(cache *github.CacheManager, client *github.Client) {
	taskList = append(taskList, NewCacheCleanTask(cache))
	taskList = append(taskList, NewCacheWarmTask(client))
	for _, t := range taskList {
		t.Run()
	}
}

func Stop() {
	for _, t := range taskList {
		t.Stop()
	}
}
