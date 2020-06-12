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
	router.Get("/getUser/{id}", DetailPost)
	router.Post("/posts", CreatePost)
	router.Put("/posts", UpdatePost)
	router.Delete("/posts/{id}", DeletePost)

	return router
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	var users UsersEntity
	var arr_user []UsersEntity
	var userVM []UsersViewModel

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

	for _, r := range arr_user {
		userVM = append(userVM,
			UsersViewModel{
				FirstName: r.FirstName,
				LastName:  r.LastName,
			},
		)
	}

	RespondWithJSON(w, 200, 200, "success", userVM, nil)

}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var user UsersRequest

	json.NewDecoder(r.Body).Decode(&user)
	if r.Body != nil {
		log.Println("ada isinya")
		fmt.Println(r)
	}

	log.Println(user.FirstName)
	query, err := db.Prepare("Insert into person(first_name,last_name)  values(?,?)")
	if err != nil {
		log.Println("error : ", err)
	}

	_, err = query.Exec(user.FirstName, user.LastName)
	if err != nil {
		log.Println("error : ", err)
	}
	defer query.Close()

	RespondWithJSON(w, 200, 200, "success", user, nil)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var user UsersEntity
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

	RespondWithJSON(w, 200, 200, "success", user, nil)

}

func DetailPost(w http.ResponseWriter, r *http.Request) {
	var user UsersEntity
	var arr_user []UsersEntity
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&user)

	fmt.Println(id)
	if r.Body != nil {
		log.Println("ada isinya")
		fmt.Println(r)
	}
	rows, err := db.Query("Select id,first_name,last_name from person where id=?", id)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName); err != nil {
			log.Fatal(err.Error())

		} else {
			arr_user = append(arr_user, user)
		}
	}

	RespondWithJSON(w, 200, 200, "success", user, nil)

}

// DeletePost remove a spesific post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var response UsersResponse

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

// RespondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, httpCode int, statCode int, message string, payload interface{}, pagination interface{}) {
	respPayload := map[string]interface{}{
		"stat_code":  statCode,
		"stat_msg":   message,
		"pagination": pagination,
		"data":       payload,
	}

	response, _ := json.Marshal(respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}
