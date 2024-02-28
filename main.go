package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't read body"}`))
		return
	}
	options := common.TranslateBytesOptions{}
	options.Pretty = pretty
	dataOut, report, err := config.TranslateBytes(body, options)
	if err != nil {
		log.Printf("Error translating config: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"can't translate config"}`))
		return
	}

	if strict && len(report.Entries) > 0 {
		log.Print("Config produced warnings and strict was specified\n")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"config produced warnings and strict was specified"}`))
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
	body, err := io.ReadAll(r.Body)
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
