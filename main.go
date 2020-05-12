package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/coreos/fcct/config"
	"github.com/coreos/fcct/config/common"
	"github.com/gorilla/mux"
)

func transpile(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strict, err := strconv.ParseBool(query.Get("strict"))
	if err != nil {
		strict = false
	}
	pretty, err := strconv.ParseBool(query.Get("pretty"))
	if err != nil {
		pretty = false
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't read body"}`))
		return
	}
	options := common.TranslateOptions{}
	options.Pretty = pretty
	options.Strict = strict
	dataOut, _, err := config.Translate(body, options)
	if err != nil {
		log.Printf("Error translating config: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't translate config"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(append(dataOut, '\n')); err != nil {
		log.Printf("Failed to write config: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/transpile", transpile).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
