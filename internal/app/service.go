package app

import (
	"19_11_2026_go/internal/checkers"
	"19_11_2026_go/internal/models"
	"19_11_2026_go/internal/reporter"
	"19_11_2026_go/internal/storages"
	"errors"
	"sync"
)

var ErrTasksNotFound = errors.New("no tasks found with given IDs")

type LinkServicer interface {
	CheckLinks(urls []string) (*models.Task, error)
	GenerateReport(ids []int) ([]byte, error)
}

type Service struct {
	storage storages.Storage
	checker *checkers.Checker
}

func NewService(s storages.Storage, c *checkers.Checker) *Service {
	return &Service{
		storage: s,
		checker: c,
	}
}

func (service *Service) CheckLinks(urls []string) (*models.Task, error) {
	task := service.storage.CreateTask(urls)

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			status := service.checker.Check(u)

			mu.Lock()
			results[u] = status
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	service.storage.UpdateTask(task.ID, results)

	return task, nil
}
func (service *Service) GenerateReport(ids []int) ([]byte, error) {
	tasks := service.storage.GetTasksByIDs(ids)
	if len(tasks) == 0 {
		return nil, ErrTasksNotFound
	}

	pdfBytes, err := reporter.GeneratePDF(tasks)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}
