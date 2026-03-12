// Package task
package task

type State string

const (
	StateOpen   State = "open"
	StateClosed State = "closed"
)

type Task struct {
	ID     string
	Title  string
	State  State
	Labels []string
	URL    string
}
