package models

// Todo list schema of the todo list table
type TodoList struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}
