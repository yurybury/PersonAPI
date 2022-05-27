package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/yurybury/UserManagement/rest/go/data"
)

// Products is a http.Handler
type Users struct {
	l *log.Logger
}

// NewUsers creates a users handler with the given logger
func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

// getUsers returns the products from the data store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")

	// fetch the products from the datastore
	lu := data.GetUsers()

	// serialize the list to JSON
	err := lu.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// getUser returns the user from the data store
func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	u.l.Println("Handle GET User", id)

	// fetch the user from the datastore
	usr := data.GetUser(id)

	// serialize the record to JSON
	err = usr.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")

	usr := r.Context().Value(KeyUser{}).(data.UserAdd)
	data.AddUser(&usr)
}

func (u Users) UpdateUsers(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	u.l.Println("Handle PUT User", id)
	usr := r.Context().Value(KeyUser{}).(data.UserUpdate)

	err = data.UpdateUser(id, &usr)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

func (u Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	u.l.Println("Handle DELETE User", id)

	err = data.DeleteUser(id)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

type KeyUser struct{}

func (u Users) MiddlewareValidateUserAdd(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := data.UserAdd{}

		err := usr.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, usr)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

func (u Users) MiddlewareValidateUserUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := data.UserUpdate{}

		err := usr.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, usr)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
