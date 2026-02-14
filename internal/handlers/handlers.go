package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/internal/model"
)

type HandlerT struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) *HandlerT {
	return &HandlerT{usecase: usecase}
}

func (p *HandlerT) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile model.Profile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := p.usecase.CreateProfile(r.Context(), profile)
	if err != nil {
		log.Printf("Ошибка создания профиля, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func (p *HandlerT) AddCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := p.usecase.AddCategory(r.Context(), category)
	if err != nil {
		log.Printf("Ошибка создания категории, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func (p *HandlerT) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	categories, err := p.usecase.GetAllCategories(r.Context(), request.UserID)
	if err != nil {
		log.Printf("Ошибка получения категорий, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
func (p *HandlerT) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
		ID     int   `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	response, err := p.usecase.DeleteCategory(r.Context(), request.UserID, request.ID)
	if err != nil {
		log.Printf("Ошибка при удалении категории")
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (p *HandlerT) AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := p.usecase.AddExpense(r.Context(), &expense)
	if err != nil {
		log.Printf("Ошибка создания расхода, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func (p *HandlerT) TodayExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	expenses, err := p.usecase.TodayExpense(r.Context(), request.UserID)
	if err != nil {
		log.Printf("Ошибка получения расходов, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}
func (p *HandlerT) WeekExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	expenses, err := p.usecase.WeekExpense(r.Context(), request.UserID)
	if err != nil {
		log.Printf("Ошибка получения расходов, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}
func (p *HandlerT) MonthExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	expenses, err := p.usecase.MonthExpense(r.Context(), request.UserID)
	if err != nil {
		log.Printf("Ошибка получения расходов, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}
func (p *HandlerT) StatsExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Некорректный JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	expenses, err := p.usecase.StatsExpense(r.Context(), request.UserID)
	if err != nil {
		log.Printf("Ошибка получения расходов, %v", err)
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}
