package handlers

import (
	"net/http"

	httphandler "github.com/Tasha-kyb/my-telegram-bot/internal/handlers/http"
	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func NewRouter(handler *httphandler.Handler) *Router {
	router := mux.NewRouter()

	router.HandleFunc("/start", handler.CreateProfile).Methods("POST")
	router.HandleFunc("/category/add", handler.AddCategory).Methods("POST")
	router.HandleFunc("/categories", handler.GetAllCategories).Methods("POST")
	router.HandleFunc("/category/delete", handler.DeleteCategory).Methods("POST")
	router.HandleFunc("/add", handler.AddExpense).Methods("POST")
	router.HandleFunc("/today", handler.TodayExpense).Methods("POST")
	router.HandleFunc("/week", handler.WeekExpense).Methods("POST")
	router.HandleFunc("/month", handler.MonthExpense).Methods("POST")
	router.HandleFunc("/stats", handler.StatsExpense).Methods("POST")

	return &Router{router: router}
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
