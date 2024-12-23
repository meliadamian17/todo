package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

type Todo struct {
	Task string `json:"task"`
}

type CLI struct {
	L  ListCommand   `cmd:"" help:"List all todos."`
	A  AddCommand    `cmd:"" help:"Add a new todo."`
	RM RemoveCommand `cmd:"" help:"Remove a todo by its index."`
	MV MoveCommand   `cmd:"" help:"Move TODO to different spot in list by its index. I.e MV 3 1 will shift todo from place 3 to place 1, shifting all remaing downwards"`
}

type ListCommand struct{}

type AddCommand struct {
	Task string `arg:"" required:"" help:"Task to add to the todo list."`
}

type RemoveCommand struct {
	Index int `arg:"" required:"" help:"Index of the todo to remove."`
}

type MoveCommand struct {
	Start  int `arg:"" required:"" help:"Index of todo to move."`
	Target int `arg:"" required:"" help:"Target index of moving todo."`
}

var todosFile string

func validIndex(l []Todo, indices ...int) error {
	n := len(l)
	for _, idx := range indices {
		if idx < 0 || idx >= n {
			return fmt.Errorf("index %d is out of bounds (1 to %d)", idx, n)
		}
	}
	return nil
}

func remove(slice []Todo, idx int) []Todo {
	return append(slice[:idx], slice[idx+1:]...)
}

func pop(slice []Todo, idx int) ([]Todo, Todo, error) {
	if idx < 0 || idx >= len(slice) {
		return nil, Todo{}, fmt.Errorf("index %d is out of bounds", idx)
	}
	popped := slice[idx]
	slice = remove(slice, idx)
	return slice, popped, nil
}

func MoveTodo(current []Todo, s, t int) ([]Todo, error) {

	if err := validIndex(current, s, t); err != nil {
		return nil, err
	}

	current, moving, err := pop(current, s)
	if err != nil {
		return nil, err
	}

	if t > len(current) {
		t = len(current)
	}

	current = append(current[:t], append([]Todo{moving}, current[t:]...)...)

	return current, nil
}

func LoadTodos() ([]Todo, error) {
	file, err := os.ReadFile(todosFile)
	if os.IsNotExist(err) {
		return []Todo{}, nil
	} else if err != nil {
		return nil, err
	}

	var todos []Todo
	err = json.Unmarshal(file, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func SaveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(todosFile, data, 0644)
}

func listTodos(l []Todo) {

	if len(l) == 0 {
		fmt.Println("No todos found. Add some!")
		return
	}

	fmt.Println("Your Todos:")
	for i, todo := range l {
		fmt.Printf("%d. %s\n", i+1, todo.Task)
	}

	return
}

func (l *ListCommand) Run() error {
	todos, err := LoadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %v", err)
	}

	listTodos(todos)

	return nil
}

func (a *AddCommand) Run() error {
	todos, err := LoadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %v", err)
	}

	todos = append(todos, Todo{Task: a.Task})

	err = SaveTodos(todos)
	if err != nil {
		return fmt.Errorf("failed to save todo: %v", err)
	}

	fmt.Printf("Added: %s\n", a.Task)
	listTodos(todos)
	return nil
}

func (r *RemoveCommand) Run() error {
	todos, err := LoadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %v", err)
	}

	if r.Index <= 0 || r.Index > len(todos) {
		return fmt.Errorf("invalid index: %d", r.Index)
	}

	task := todos[r.Index-1].Task
	todos = append(todos[:r.Index-1], todos[r.Index:]...)

	err = SaveTodos(todos)
	if err != nil {
		return fmt.Errorf("failed to save todos: %v", err)
	}

	fmt.Printf("Removed: %s\n", task)
	listTodos(todos)
	return nil
}

func (mv *MoveCommand) Run() error {
	current, err := LoadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %v", err)
	}
	todos, err := MoveTodo(current, mv.Start-1, mv.Target-1)
	if err != nil {
		return fmt.Errorf("failed to move todos: %v", err)
	}

	if err := SaveTodos(todos); err != nil {
		return fmt.Errorf("failed to save todos: %v", err)
	}

	fmt.Printf("Moved %v from %v to %v\n", current[mv.Start-1].Task, mv.Start, mv.Target)
	listTodos(todos)
	return nil
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine home directory: %v\n", err)
		os.Exit(1)
	}
	todosFile = filepath.Join(homeDir, ".config", "todo", "todos.json")

	var cli CLI
	ctx := kong.Parse(&cli)
	err = ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
