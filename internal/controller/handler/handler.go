package handler

import (
	"blog_service2/internal/model"
	"blog_service2/internal/repository"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	log  *log.Logger
	repo repository.Storage
}

func New() (*Handler, error) {
	l := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime)
	r, err := repository.New()
	if err != nil {
		return nil, err
	}

	return &Handler{
		log:  l,
		repo: r,
	}, nil
}

func (h *Handler) getID(w http.ResponseWriter, ps httprouter.Params) (int, bool) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		h.log.Printf("handler: %s", err.Error())
		return 0, false
	}
	return id, true
}

func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var str string
	if len(r.URL.RawQuery) > 0 {
		str = r.URL.Query().Get("title")
		if str == "" {
			w.WriteHeader(400)
			h.log.Printf("handler: Invalid argument")
			return
		}
	}

	recs, err := h.repo.Read(str)
	if err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(recs); err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
	}
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := h.getID(w, ps)
	if !ok {
		return
	}
	rec, err := h.repo.ReadOne(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			h.log.Printf("handler: %s", err.Error())
			return
		}
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(rec); err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
	}
}

func (h *Handler) AddRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var rec model.Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil || rec.Title == "" || rec.Text == "" {
		w.WriteHeader(400)
		h.log.Printf("handler: Incorrect input data")
		return
	}
	if err := h.repo.Insert(rec); err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
		return
	}
	w.WriteHeader(201)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := h.getID(w, ps)
	if !ok {
		return
	}
	var rec model.Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil || rec.Title == "" || rec.Text == "" {
		w.WriteHeader(400)
		h.log.Printf("handler: Incorrect input data")
		return
	}
	rec.Id = id

	res, err := h.repo.Update(rec)
	if err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
		return
	}

	if res == 0 {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := h.getID(w, ps)
	if !ok {
		return
	}
	if err := h.repo.Remove(id); err != nil {
		w.WriteHeader(500)
		h.log.Printf("handler: %s", err.Error())
	}
	w.WriteHeader(204)
}
