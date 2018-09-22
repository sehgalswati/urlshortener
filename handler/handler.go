// Package handlers provides HTTP request handlers.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/sehgalswati/urlshortener/backend"
)

// New returns an http handler for the url shortener.
func New(prefix string, backend backend.BackendService) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, backend}
	mux.HandleFunc("/encode/", responseHandler(h.encode))
	mux.HandleFunc("/", h.redirect)
	mux.HandleFunc("/urlinfo/", responseHandler(h.decode))
	//mux.HandleFunc("/removestaleentries/", h.removestale)
	return mux
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type handler struct {
	prefix  string
	backend backend.BackendService
}

func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}

func (h handler) encode(w io.Writer, r *http.Request) (interface{}, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method)
	}

	var input struct{ URL string }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}

	url := strings.TrimSpace(input.URL)
	if url == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("URL is empty")
	}

	c, err := h.backend.Save(url)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}

	return h.prefix + c, http.StatusCreated, nil
}

func (h handler) decode(w io.Writer, r *http.Request) (interface{}, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	code := r.URL.Path[len("/urlinfo/"):]
	//model, err := h.backend.Load(code)
	model, err := h.backend.LoadInfo(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return model, http.StatusOK, nil
}

func (h handler) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Path[len("/"):]

	url, err := h.backend.Load(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
		return
	}

	http.Redirect(w, r, string(url), http.StatusMovedPermanently)
}
/*
func (h handler) removestale(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete {
		fmt.Errorf("Method %s not allowed", r.Method)
	}
	c, err := h.backend.RemoveStale()
	if err != nil {
		fmt.Errorf("Couldn't find any stale entry in database: %v", err)
	}
	fmt.Printf("remove stale in handler.go uid being removed are %v", c )
	//http.Redirect(w, r, string(h.prefix+url), http.StatusOK)
	//return http.StatusOK, nil
	
	
	code := r.URL.Path[len("/deletetuple/"):]
	model, err := h.backend.Load(code)
	model, err := h.backend.DeleteTuple(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}
	http.Redirect(w, r, string(url), http.StatusMovedPermanently)
	return model, http.StatusOK, nil
	
}
*/
