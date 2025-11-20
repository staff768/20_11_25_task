package models

type Task struct {
	ID      int               `json:"links_num"`
	URLs    []string          `json:"-"`
	Results map[string]string `json:"links"`
}
type CheckRequest struct {
	Links []string `json:"links"`
}
type GenerateRequest struct {
	IDs []int `json:"links_list"`
}
