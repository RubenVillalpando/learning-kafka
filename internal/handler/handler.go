package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RubenVillalpando/learning-kafka/internal/db"
	"github.com/RubenVillalpando/learning-kafka/internal/kafka"
	"github.com/RubenVillalpando/learning-kafka/internal/model"
)

type Handler struct {
	db          *db.DB
	kafkaClient *kafka.Client
}

func New(kc *kafka.Client) *Handler {

	return &Handler{
		db:          db.New(),
		kafkaClient: kc,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	default:
		fmt.Fprintf(w, "Unsupported method: %s", r.Method)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	if id == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Missing Query Parameter id"))
		return
	}

	customer, err := h.db.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was a problem with the server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was a problem with the server"))
		return
	}

	id, err := h.db.StoreOne(&customer)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("user created with id: %s", id)))
	h.kafkaClient.ProduceMessage(&customer)

}
