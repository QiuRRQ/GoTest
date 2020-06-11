package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var router *chi.Mux
var db *sql.DB

const (
	dbName = "mygo"
	dbPass = ""
	dbHost = "localhost"
	dbPort = "3306"
)

func main() {
	routers()
	log.Fatal(http.ListenAndServe(":8005", routers()))
}

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8", dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dbSource)

	if err != nil {
		log.Println("error : ", err)
	}
}

func routers() *chi.Mux {
	router.Get("/posts", returnAllUsers)
	router.Get("/getUser", DetailPost)
	router.Post("/posts", CreatePost)
	router.Put("/posts/{id}", UpdatePost)
	router.Delete("/posts/{id}", DeletePost)

	return router
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	var users Users
	var arr_user []Users
	var response Response

	db := connect()
	defer db.Close()

	rows, err := db.Query("Select id,first_name,last_name from person")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&users.Id, &users.FirstName, &users.LastName); err != nil {
			log.Fatal(err.Error())

		} else {
			arr_user = append(arr_user, users)
		}
	}

	response.Status = 1
	response.Message = "Success"
	response.Data = arr_user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var user Users
	var response Response

	json.NewDecoder(r.Body).Decode(&user)
	if r.Body != nil {
		log.Println("ada isinya")
		fmt.Println(r)
	}
	query, err := db.Prepare("Insert person SET first_name=?, last_name=?")
	if err != nil {
		log.Println("error : ", err)
	}

	_, er := query.Exec(user.FirstName, user.LastName)
	if er != nil {
		log.Println("error : ", er)
	}
	defer query.Close()

	w.Header().Set("Content-Type", "application/json")
	response.Status = 1
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var user Users
	var response Response
	// id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&user)

	query, err := db.Prepare("Update person set first_name=?, last_name=? where id=?")
	if err != nil {
		log.Println("error : ", err)
	}
	_, er := query.Exec(user.FirstName, user.LastName, user.Id)

	if er != nil {
		log.Println("error : ", er)
	}
	defer query.Close()

	w.Header().Set("Content-Type", "application/json")
	response.Status = 1
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)

}

func DetailPost(w http.ResponseWriter, r *http.Request) {
	var user Users
	var arr_user []Users
	var response Response
	// id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&user)

	query, err := db.Prepare("Select id,first_name,last_name from person where id=?")
	if err != nil {
		log.Println("error : ", err)
	}
	s, er := query.Exec(user.Id)

	if er != nil {
		log.Println("error : ", er)
	}
	defer query.Close()

	rows, err := s.RowsAffected()
	log.Println(rows)
	// scanErr := rows.Scan(&user.Id, &user.FirstName, &user.LastName)
	if scanErr != nil {
		log.Fatal(err.Error())

	} else {
		arr_user = append(arr_user, user)
	}

	response.Status = 1
	response.Message = "Success"
	response.Data = arr_user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// DeletePost remove a spesific post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var response Response

	query, err := db.Prepare("delete from person where id=?")
	if err != nil {
		log.Println("error : ", err)
	}
	_, er := query.Exec(id)
	if er != nil {
		log.Println("error : ", er)
	}
	query.Close()

	w.Header().Set("Content-Type", "application/json")
	response.Status = 1
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)

}
