

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var students = []Student{
	{ID: 1, Name: "Ali", Age: 18, Email: "ali@gmail.com"},
	{ID: 2, Name: "Aruzhan", Age: 19, Email: "aruzhan@gmail.com"},
}

var nextID = 3
var mu sync.Mutex

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Student Management System!")
	fmt.Fprintln(w, "Go Web Application is running successfully.")
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var s Student

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid data", 400)
		return
	}

	mu.Lock()
	s.ID = nextID
	nextID++
	students = append(students, s)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var id int
	fmt.Sscanf(r.URL.Query().Get("id"), "%d", &id)

	mu.Lock()
	defer mu.Unlock()

	for i, s := range students {
		if s.ID == id {
			students = append(students[:i], students[i+1:]...)
			fmt.Fprintln(w, "Student deleted successfully!")
			return
		}
	}

	http.Error(w, "Student not found", 404)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/students", getStudents)
	http.HandleFunc("/add", addStudent)
	http.HandleFunc("/delete", deleteStudent)

	fmt.Println("Server started on port 3333")

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
