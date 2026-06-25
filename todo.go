package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const tasksFile = "todos.json"

func main() {
	fmt.Print("Todolist\n\n")

	if len(os.Args) < 2 {
		fmt.Println("veuillez renseigner une commande (add, list, done, delete)")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("veuillez renseigner un titre pour votre tâche")
			return
		}

		title := strings.TrimSpace(strings.Join(os.Args[2:], " "))
		if title == "" {
			fmt.Println("veuillez renseigner un titre pour votre tâche")
			return
		}
		addTask(title)

	case "list":
		listTasks()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("veuillez renseigner l'ID de la tâche à marquer complétée")
			return
		}

		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		markTaskCompleted(id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("veuillez renseigner l'ID de la tâche à supprimer")
			return
		}

		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		deleteTask(id)

	default:
		fmt.Println("commande inconnue: ", command)
		fmt.Println("commandes disponibles: add, list, done, delete")
	}
}

func addTask(title string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	nextID := 1

	if len(tasks) > 0 {
		maxID := 0
		for i := range tasks {
			if tasks[i].ID > maxID {
				maxID = tasks[i].ID
			}
		}

		nextID = maxID + 1
	}

	tasks = append(tasks, Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
	})

	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	listTasks()
	fmt.Println("\ntâche ajoutée avec succès")
}

func listTasks() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("aucune tâche pour l'instant")
		return
	}

	for i := range tasks {
		if tasks[i].Completed {
			fmt.Print("[x] ")
		} else {
			fmt.Print("[ ] ")
		}
		fmt.Printf("%d - %s\n", tasks[i].ID, tasks[i].Title)
	}
}

func markTaskCompleted(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		fmt.Println("pas d'ID correspondant")
		return
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	listTasks()
	fmt.Printf("\nla tâche %d est marquée comme complétée\n", id)
}

func loadTasks() ([]Task, error) {
	content, err := os.ReadFile(tasksFile)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(string(content)) == "" {
		return []Task{}, nil
	}

	var tasks []Task

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	fileContent, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(tasksFile, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("impossible de charger les tâches:", err)
		return
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("pas d'ID correspondant")
		return
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("impossible de sauvegarder les tâches:", err)
		return
	}

	listTasks()
	fmt.Printf("\nla tâche %d a été supprimée\n", id)
}

func parseID(value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if id < 0 {
		return 0, fmt.Errorf("ID invalide, veuillez renseigner un nombre positif")
	}

	return id, nil
}
