package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var lastID int

func main() {
	fmt.Println("Welcome to do !")

	addTask("coucou")
}

func addTask(title string) {
	// load existing tasks
	content, err := os.ReadFile("tasks.json")
	if err != nil {
		fmt.Println("error reading file")
	}

	var tasks []Task

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		fmt.Println("error:", err)
	}

	if len(tasks) == 0 {
		lastID = 1
	} else {
		// last ID of the slice
		lastID = tasks[len(tasks)-1].ID
		lastID += 1
	}

	// create a task
	tasks = append(tasks, Task{
		ID:        lastID,
		Title:     title,
		Completed: false,
	})

	// save in file
	fileContent, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = os.WriteFile("tasks.json", fileContent, 0777)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("tâche ajoutée avec succès!")
	}
}
