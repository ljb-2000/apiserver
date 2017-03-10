package app

import (
	"github.com/gorilla/mux"
)

func Run(addr string) error {
	root := mux.NewRouter()
	api := root.PathPrefix("/api/v1").Subrouter()
	installApiGroup(api)
	http.Handle("/", root)
	return http.ListenAndServe(addr, nil)
}

func installApiGroup(router *mux.Router) {

}

func installNodeApi(router *mux.Router) {

}

func installAppApi(router *mux.Router) {

}

func installContainerApi(router *mux.Router) {

}

func installDeploymentApi(router *mux.Router) {

}
