package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "assem254673"
	dbname   = "assem"
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	// Connect to the PostgreSQL database
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new task
	taskName := "Task 2"
	err = createTask(db, taskName)
	if err != nil {
		panic(err)
	}

	// List all tasks
	// tasks, err := listTasks(db)
	// if err != nil {
	//  panic(err)
	// }
	// fmt.Println("Tasks:")
	// for _, task := range tasks {
	//  fmt.Printf("%d. %s (Completed: %t)\n", task.ID, task.Name, task.Completed)
	// }

	// // Update a task (mark as completed)
	// taskID := 2
	// err = updateTask(db, taskID)
	// if err != nil {
	//  panic(err)
	// }

	// // Delete a task
	// taskID := 2
	// err = deleteTask(db, taskID)
	// if err != nil {
	// 	panic(err)
	// }
}

func createTask(db *sql.DB, name string) error {
	// Check if task with the same name already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE name = $1", name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("task with name %s already exists", name)
	}

	// Insert the new task into the database
	_, err = db.Exec("INSERT INTO tasks(name) VALUES($1)", name)
	return err
}

func listTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func updateTask(db *sql.DB, id int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the task's completion status
	_, err = tx.Exec("UPDATE tasks SET completed = TRUE WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	return err
}

func deleteTask(db *sql.DB, id int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the task from the database
	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	return err
}
