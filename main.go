package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/coreos/fcct/config"
	"github.com/coreos/fcct/config/common"
	iconfig "github.com/coreos/ignition/v2/config"
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

func validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't read body"}`))
		return
	}

	cfg, rpt, err := iconfig.Parse(body)
	if err != nil {
		log.Printf("Error validating config: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't validate config"}`))
		return
	}

	if len(rpt.Entries) > 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"error":"validation failed", "reason":"` + rpt.String() + `"}`))
		return
	}
	jsonConfig, jsonerr := json.MarshalIndent(cfg, "", "  ")
	if jsonerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"config json marshalling failed"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(append(jsonConfig, '\n')); err != nil {
		log.Printf("Failed to write validated config: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/transpile", transpile).Methods(http.MethodPost)
	api.HandleFunc("/validate", validate).Methods(http.MethodPost)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", r))
}
