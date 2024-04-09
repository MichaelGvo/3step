package main

import (
	"bytes"
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
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже написаны обработчики для каждого эндпоинта
// func getTaskFromMap - обработчик для получения всех задач, возвращает все tasks, которые хранятся в map[string]Task
// func getTaskFromMap принимает в качестве параметров запрос r *http.Request и ответ w http.ResponseWriter
func getTaskFromMap(w http.ResponseWriter, r *http.Request) {
	// err = nil - мы успешно сериализуем данные из tasks (var tasks = map[string]Task)
	// err != nil - вызываем ошибку (http.Error) с http.StatusInternalServerError - при ошибке сервер должен вернуть статус 500 Internal Server Error
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента, формат данных - Content-Type — application/json
	w.Header().Set("Content-Type", "application/json")
	// все успешно - http.StatusOK
	w.WriteHeader(http.StatusOK)
	// далее записанные в JSON данные записываем в тело ответа
	w.Write(resp)
}

// func postTaskOnServer - обработчик для отправки задачи на сервер, принимает задачу в теле запроса и сохраняет ее в мапе
// func postTaskOnServer принимает в качестве параметров запрос r *http.Request и ответ w http.ResponseWriter
func postTaskOnServer(w http.ResponseWriter, r *http.Request) {
	var task Task
	// используем тип bytes.Buffer для работы с байтовыми данными, в данном случае с содержимым тела запроса (r.Body)
	var buf bytes.Buffer

	// err = nil - чтение прошло успешно
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = nil - успешно десериализуем JSON из буфера в структуру Task
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	// в заголовок записываем тип контента, формат данных - Content-Type — application/json
	w.Header().Set("Content-Type", "application/json")
	// при успешном запросе сервер возвращает 201 - Created
	w.WriteHeader(http.StatusCreated)

}

// func getTaskFromMapById - обработчик для получения задачи по ID, должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе
// func getTaskFromMapById принимает в качестве параметров запрос r *http.Request и ответ w http.ResponseWriter
func getTaskFromMapById(w http.ResponseWriter, r *http.Request) {
	//chi.URLParam() принимает r, представляющий запрос, и строку с именем параметра (“id”)
	id := chi.URLParam(r, "id")

	// проверяю существует ли задача с таким ID
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задание не найдено", http.StatusBadRequest)
		return
	}

	// err = nil - мы успешно сериализуем данные из tasks (var tasks = map[string]Task)
	// err != nil - вызываем ошибку (http.Error) с http.StatusBadRequest - при ошибке сервер должен вернуть статус 400 Bad Request
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// в заголовок записываем тип контента, формат данных - Content-Type — application/json
	w.Header().Set("Content-Type", "application/json")
	// все успешно - http.StatusOK
	w.WriteHeader(http.StatusOK)
	// далее записанные в JSON данные записываем в тело ответа
	w.Write(resp)
}

// func deleteTaskFromMapById - обработчик удаления задачи по ID
// func deleteTaskFromMapById принимает в качестве параметров запрос r *http.Request и ответ w http.ResponseWriter
func deleteTaskFromMapById(w http.ResponseWriter, r *http.Request) {
	var task Task
	// используем тип bytes.Buffer для работы с байтовыми данными, в данном случае с содержимым тела запроса (r.Body)
	var buf bytes.Buffer

	// err = nil - чтение прошло успешно
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = nil - успешно десериализуем JSON из буфера в структуру Task
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delete(tasks, task.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	// здесь отрегистрированы мои обработчики
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTaskFromMap`
	r.Get("/tasks", getTaskFromMap)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTaskOnServer`
	r.Post("/tasks", postTaskOnServer)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTaskFromMapById`
	r.Get("/tasks/{id}", getTaskFromMapById)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTaskFromMapById`
	r.Delete("/tasks/{id}", deleteTaskFromMapById)

	// запускаем сервер

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
