package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestValidConfig(t *testing.T) {
	infile, err := os.Open("./test/fcos-config.yaml")
	if err != nil {
		t.Fatalf("failed to open fcos config: %v\n", err)
	}
	defer infile.Close()

	dataIn, err := ioutil.ReadAll(infile)
	if err != nil {
		t.Fatalf("failed to read %s: %v\n", infile.Name(), err)
	}

	req, err := http.NewRequest("POST", "/api/v1/transpile", bytes.NewBuffer(dataIn))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/x-yaml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transpile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"ignition":{"config":{"replace":{"source":null,"verification":{}}},"security":{"tls":{}},"timeouts":{},"version":"3.0.0"},"passwd":{},"storage":{"files":[{"group":{},"path":"/etc/test.cfg","user":{},"contents":{"source":"data:,test-test-test","verification":{}},"mode":384}]},"systemd":{"units":[{"contents":"[Unit]\nDescription=Test\nAfter=network-online.target\nWants=network.target\n\n[Service]\nExecStart=/usr/bin/test\n\n[Install]\nWantedBy=multi-user.target\n","enabled":true,"name":"test.service"},{"mask":true,"name":"docker.service"}]}}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestInvalidConfig(t *testing.T) {
	infile, err := os.Open("./test/fcos-error.yaml")
	if err != nil {
		t.Fatalf("failed to open fcos config: %v\n", err)
	}
	defer infile.Close()

	dataIn, err := ioutil.ReadAll(infile)
	if err != nil {
		t.Fatalf("failed to read %s: %v\n", infile.Name(), err)
	}

	req, err := http.NewRequest("POST", "/api/v1/transpile?strict=true", bytes.NewBuffer(dataIn))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/x-yaml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transpile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	expected := `{"error":"can't translate config"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
