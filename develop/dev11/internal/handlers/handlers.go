package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/KapDmitry/WB_L2/develop/dev11/internal/event"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/logger"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/response"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Repo      event.EventRepository
	Logger    logger.Logger
	Validator *validator.Validate
}

func (h *Handler) validateJSON(newEvent event.Event) error {
	err := h.Validator.Struct(newEvent)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			h.Logger.LogW("Error", "valiudation errors", map[string]interface{}{"Field": err.Field(), "Msg": err.Tag()})
		}
		return fmt.Errorf("not valid")
	}
	return nil

}

func (h *Handler) setTime(rawDate string) (time.Time, error) {
	timeStr := "2006-01-02"
	date, err := time.Parse(timeStr, rawDate)
	if err != nil {
		h.Logger.Log("Error", err.Error())
		return time.Now(), err
	}
	return date, nil
}

func (h *Handler) parseJSON(r io.Reader) (event.Event, error) {
	var newEvent event.Event
	err := json.NewDecoder(r).Decode(&newEvent)
	if err != nil {
		return newEvent, fmt.Errorf("not valid JSON")
	}

	err = h.validateJSON(newEvent)
	if err != nil {
		return newEvent, fmt.Errorf("not valid JSON")
	}

	newEvent.Date, err = h.setTime(newEvent.RawDate)
	if err != nil {
		return newEvent, fmt.Errorf("not valid JSON")
	}
	return newEvent, nil
}

func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 400)
		return
	}

	newEvent, err := h.parseJSON(r.Body)
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	newEvent, err = h.Repo.Create(newEvent)
	if err != nil {
		h.Logger.Log("Error", err.Error())
		response.ServerResponseWriter(w, 503, map[string]interface{}{"error": "internal server error"})
		return
	}

	response.ServerResponseWriter(w, 200, map[string]interface{}{"result": newEvent})

}
func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 400)
		return
	}

	newEvent, err := h.parseJSON(r.Body)
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	newEvent, err = h.Repo.Update(newEvent)
	if err != nil {
		h.Logger.Log("Error", err.Error())
		response.ServerResponseWriter(w, 503, map[string]interface{}{"error": "internal server error"})
		return
	}

	response.ServerResponseWriter(w, 200, map[string]interface{}{"result": newEvent})
}

func (h *Handler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 400)
		return
	}
	newEvent, err := h.parseJSON(r.Body)
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	err = h.Repo.Delete(newEvent.UserID, newEvent.ID)
	if err != nil {
		h.Logger.Log("Error", err.Error())
		response.ServerResponseWriter(w, 503, map[string]interface{}{"error": "internal server error"})
		return
	}

	response.ServerResponseWriter(w, 200, map[string]interface{}{"result": "OK"})
}

func (h *Handler) EventsForHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 400)
		return
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Failed to parse URL", 400)
		return
	}

	endpoint := u.Path

	err = r.ParseForm()
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	rawUserID := r.Form.Get("user_id")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}
	rawDate := r.Form.Get("date")
	date, err := h.setTime(rawDate)
	if err != nil {
		response.ServerResponseWriter(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}

	var endTime time.Time
	if endpoint == "/events_for_day" {
		endTime = date.AddDate(0, 0, 1)
	}
	if endpoint == "/events_for_month" {
		fmt.Println("true")
		endTime = date.AddDate(0, 1, 0)
	}
	if endpoint == "/events_for_week" {
		endTime = date.AddDate(0, 0, 7)
	}

	events, err := h.Repo.GetEventsForRange(int64(userID), date, endTime)
	if err != nil {
		h.Logger.Log("Error", err.Error())
		response.ServerResponseWriter(w, 503, map[string]interface{}{"error": "internal server error"})
		return
	}

	response.ServerResponseWriter(w, 200, map[string]interface{}{"result": events})
}
