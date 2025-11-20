package storages

import (
	"19_11_2026_go/internal/models"
	"encoding/json"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type Storage interface {
	CreateTask(urls []string) *models.Task
	UpdateTask(taskID int, results map[string]string)
	GetTasksByIDs(ids []int) []*models.Task
	Save() error
	Load() error
}

type FileStorage struct {
	mu       sync.RWMutex
	tasks    map[int]*models.Task
	nextID   int32
	filePath string
}

func NewStorage(filePath string) (*FileStorage, error) {
	s := &FileStorage{
		tasks:    make(map[int]*models.Task),
		filePath: filePath,
		nextID:   0,
	}
	err := s.Load()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

func (s *FileStorage) CreateTask(urls []string) *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID := atomic.AddInt32(&s.nextID, 1)

	results := make(map[string]string, len(urls))
	for _, url := range urls {
		results[url] = "pending"
	}

	task := &models.Task{
		ID:      int(newID),
		URLs:    urls,
		Results: results,
	}
	s.tasks[task.ID] = task

	return task
}

func (s *FileStorage) UpdateTask(taskID int, results map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[taskID]; ok {
		task.Results = results
	}
}

func (s *FileStorage) GetTasksByIDs(ids []int) []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	foundTasks := make([]*models.Task, 0, len(ids))
	for _, id := range ids {
		if task, ok := s.tasks[id]; ok {
			foundTasks = append(foundTasks, task)
		}
	}
	return foundTasks
}

func (s *FileStorage) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *FileStorage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Storage file not found")
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &s.tasks); err != nil {
		log.Printf("Could not unmarshal storage data: %v", err)
		s.tasks = make(map[int]*models.Task)
		return nil
	}

	var maxID int32 = 0
	for id := range s.tasks {
		if int32(id) > maxID {
			maxID = int32(id)
		}
	}
	s.nextID = maxID
	log.Println("Successfully loaded tasks from storage")

	return nil
}
