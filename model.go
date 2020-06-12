package main

type UsersRequest struct { //for taking json post request
	FirstName string `form:"firstname" json:"first_name"`
	LastName  string `form:"lastname" json:"last_name"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []UsersViewModel
}

type UsersViewModel struct { //for returning fesponse format
	Id        string `form:"id" json:"id"`
	FirstName string `form:"firstname" json:"first_name"`
	LastName  string `form:"lastname" json:"last_name"`
}

type UsersEntity struct { //for database operational
	Id        int    `form:"id" json:"id"`
	FirstName string `form:"firstname" json:"first_name"`
	LastName  string `form:"lastname" json:"last_name"`
}
