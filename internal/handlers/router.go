package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouterT struct {
	router *mux.Router
}

func NewRouter(handler *HandlerT) *RouterT {
	router := mux.NewRouter()

	router.HandleFunc("/start", handler.CreateProfile).Methods("POST")
	router.HandleFunc("/addCategory", handler.AddCategory).Methods("POST")
	router.HandleFunc("/categories", handler.GetAllCategories).Methods("POST")
	router.HandleFunc("/category/delete", handler.DeleteCategory).Methods("POST")
	router.HandleFunc("/add", handler.AddExpense).Methods("POST")
	router.HandleFunc("/today", handler.TodayExpense).Methods("POST")
	router.HandleFunc("/week", handler.WeekExpense).Methods("POST")
	router.HandleFunc("/month", handler.MonthExpense).Methods("POST")
	router.HandleFunc("/stats", handler.StatsExpense).Methods("POST")

	return &RouterT{router: router}
}

func (s *RouterT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
