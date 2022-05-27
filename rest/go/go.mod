module github.com/yurybury/UserManagement/rest/go

go 1.18

require (
	github.com/gorilla/mux v1.8.0
	github.com/nicholasjackson/env v0.6.0
)

replace github.com/yurybury/UserManagement/rest/go/handlers => ./handlers

replace github.com/yurybury/UserManagement/rest/go/data => ./data
