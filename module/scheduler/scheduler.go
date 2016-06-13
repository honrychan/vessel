package scheduler

import (
	"log"
	"fmt"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
	"github.com/containerops/vessel/utils/timer"
)

// StartStage start stage on scheduler
func StartStage(executorMap map[string]*models.Executor, hourglass *timer.Hourglass) []*models.ExecutedResult {
	readyMap := map[string]bool{"":true}

	count := len(executorMap)
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)

	running := true
	handler := stage.StartStage
	for running {
		go startProgress(executorMap, readyMap, finishChan, hourglass, handler)
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			readyMap[result.Name] = true
			running = len(resultList) != count
			log.Println(fmt.Sprintf("scheduler StartStage name = %v and finish num = %d", result.Name, len(resultList)))
		}
	}
	return resultList
}

// StopStage stop stage on scheduler
func StopStage(executorMap map[string]*models.Executor, hourglass *timer.Hourglass) []*models.ExecutedResult {
	readyMap := map[string]bool{"":true}

	count := len(executorMap)
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)

	running := true
	handler := stage.StopStage
	for running {
		go startProgress(executorMap, readyMap, finishChan, hourglass, handler)
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			readyMap[result.Name] = true
			running = len(resultList) != count
			log.Println(fmt.Sprintf("scheduler StopStage name = %v and finish num = %d", result.Name, len(resultList)))
		}
	}
	return resultList
}

func startProgress(executorMap map[string]*models.Executor, readyMap map[string]bool, finishChan chan *models.ExecutedResult, hourglass *timer.Hourglass,
	handler func(info interface{}, finishChan chan *models.ExecutedResult, hourglass *timer.Hourglass)) {
	for name, executor := range executorMap {
		if _, ok := readyMap[name]; ok {
			continue
		}
		isReady := true
		for _, from := range executor.From {
			if isReady, _ = readyMap[from]; !isReady {
				break
			}
		}
		if !isReady {
			continue
		}
		readyMap[name] = false
		go handler(executor.Info, finishChan, hourglass)
	}
}