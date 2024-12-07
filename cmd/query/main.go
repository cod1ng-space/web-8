package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgre"
	dbname   = "sandbox"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	msg, err := h.dbProvider.SelectQuery()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else if msg == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неинициализированное значение в базе данных!"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, " + msg + "!"))
	}
}
func (h *Handlers) PostQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	input := struct {
		Msg string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else if input.Msg == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Пустая строка!"))
	} else {
		err = h.dbProvider.UpdateQuery(input.Msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

// Методы для работы с базой данных
func (dp *DatabaseProvider) SelectQuery() (string, error) {
	var msg string

	row := dp.db.QueryRow("SELECT name FROM query")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}

	return msg, nil
}
func (dp *DatabaseProvider) UpdateQuery(msg string) error {
	_, err := dp.db.Exec("UPDATE query SET name = $1", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	// Регистрируем обработчики
	http.HandleFunc("/get", h.GetQuery)
	http.HandleFunc("/post", h.PostQuery)

	fmt.Println("Сервер запущен")
	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
