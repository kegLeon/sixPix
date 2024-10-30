package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура !",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное  задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый  раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

type ResponseAll struct {
	Tasks []Task `json:"tasks"`
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
		return
	}
	var responseElements []Task
	for _, v := range tasks {
		responseElements = append(responseElements, v)
	}
	response := ResponseAll{Tasks: responseElements}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Allarm", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandlerInServ(w http.ResponseWriter, r *http.Request) {
	var task Task
	if r.Method != http.MethodPost {
		http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, "allarmInServ", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.WriteHeader(http.StatusCreated)
}
func HandlerId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
		return
	}
	idsh := r.URL.Query().Get("ID")
	v, exists := tasks[idsh]
	if exists {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
			http.Error(w, "Ошибка декода", http.StatusBadRequest)
			return
		}
		return
	}
	http.Error(w, "Ошибка в айlишной обработке", http.StatusBadRequest)

}
func HandlerDelID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
		return
	}
	idsh := r.URL.Query().Get("ID")
	_, exists := tasks[idsh]
	if exists {
		delete(tasks, idsh)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Такой задачи нет", http.StatusBadRequest)
}
func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	http.HandleFunc("/tasks", allHandler)
	http.HandleFunc("/tasks", HandlerInServ)
	http.HandleFunc("/tasks{id}", HandlerId)
	http.HandleFunc("/tasks{id}", HandlerDelID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
