package tasks

import "copilotlens/internal/github"

var taskList []ITask

type ITask interface {
	Run()
	Stop()
}

func Init(cache *github.CacheManager, client *github.Client) *CacheWarmTask {
	cleanTask := NewCacheCleanTask(cache)
	warmTask := NewCacheWarmTask(client)
	taskList = append(taskList, cleanTask, warmTask)
	for _, t := range taskList {
		t.Run()
	}
	return warmTask
}

func Stop() {
	for _, t := range taskList {
		t.Stop()
	}
}
