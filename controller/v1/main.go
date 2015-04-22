package v1

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func Router() http.Handler {
	n := negroni.New()
	r := mux.NewRouter()

	// Users collection
	users := r.PathPrefix("/users").Subrouter()
	users.Methods("POST").HandlerFunc(UserCreate)

	// Users singular
	user := r.PathPrefix("/users/{id}").Subrouter()
	user.Methods("GET").HandlerFunc(UserShow)
	user.Methods("PUT", "PATCH").HandlerFunc(UserUpdate)
	user.Methods("DELETE").HandlerFunc(UserDestroy)
	user.Methods("GET").Path("/projects").HandlerFunc(ProjectList)
	user.Methods("POST").Path("/projects").HandlerFunc(ProjectCreate)

	// Tokens collection
	tokens := r.PathPrefix("/tokens").Subrouter()
	tokens.Methods("POST").HandlerFunc(TokenCreate)

	// Tokens singular
	token := r.PathPrefix("/tokens/{key}").Subrouter()
	token.Methods("DELETE").HandlerFunc(TokenDestroy)

	// Projects singular
	project := r.PathPrefix("/projects/{id}").Subrouter()
	project.Methods("GET").HandlerFunc(ProjectShow)
	project.Methods("PUT", "PATCH").HandlerFunc(ProjectUpdate)
	project.Methods("DELETE").HandlerFunc(ProjectDestroy)
	project.Methods("GET").Path("/full").HandlerFunc(ProjectFull)

	n.UseHandler(r)

	return n
}
