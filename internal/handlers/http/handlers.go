package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Tasha-kyb/my-telegram-bot/internal/app"
	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

type Handler struct {
	usecase *app.Service
}

func NewHandler(usecase *app.Service) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	profile := model.Profile{
		ID:        request.UserID,
		Username:  request.Username,
		CreatedAt: time.Now(),
	}

	response, err := h.usecase.CreateProfile(r.Context(), profile)
	if err != nil {
		log.Printf("Ошибка создания профиля, %v", err)
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}

	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) AddCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := h.decodeJSON(r, &category)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()

	response, err := h.usecase.AddCategory(r.Context(), category)
	if err != nil {
		h.handleError(w, err, "Ошибка при создании категории")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	categories, err := h.usecase.GetAllCategories(r.Context(), request.UserID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, categories)
}
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
		ID     int   `json:"id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	response, err := h.usecase.DeleteCategory(r.Context(), request.UserID, request.ID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense
	err := h.decodeJSON(r, &expense)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()

	log.Printf("📥 Получен expense: UserID=%d, Amount=%.2f, Category=%s, Description=%s",
		expense.UserID, expense.Amount, expense.Category, expense.Description)

	response, err := h.usecase.AddExpense(r.Context(), &expense)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) TodayExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	response, err := h.usecase.TodayExpense(r.Context(), request.UserID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) WeekExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	response, err := h.usecase.WeekExpense(r.Context(), request.UserID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) MonthExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	response, err := h.usecase.MonthExpense(r.Context(), request.UserID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) StatsExpense(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int64 `json:"user_id"`
	}
	err := h.decodeJSON(r, &request)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	defer r.Body.Close()
	response, err := h.usecase.StatsExpense(r.Context(), request.UserID)
	if err != nil {
		h.handleError(w, err, "Внутренняя ошибка сервера")
		return
	}
	h.respondJSON(w, http.StatusCreated, response)
}
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Ошибка кодирования JSON-ответа %v", err)
		}
	}
}
func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
func (h *Handler) decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
	log.Printf("⚠️ Ошибка: %v", err)
	switch {
	case strings.Contains(err.Error(), "не найдена"):
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	case strings.Contains(err.Error(), "уже существует"):
		h.respondError(w, http.StatusConflict, err.Error())
		return
	case strings.Contains(err.Error(), "не хватает параметров"):
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	default:
		h.respondError(w, http.StatusInternalServerError, defaultMessage)
		return
	}
}
