package models

type Task struct {
	ID      int               `json:"links_num"`
	URLs    []string          `json:"-"`
	Results map[string]string `json:"links"`
}
