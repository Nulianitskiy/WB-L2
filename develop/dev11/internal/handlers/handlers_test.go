package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/KapDmitry/WB_L2/develop/dev11/internal/event"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/logger"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitiateHandlerUser(rep *event.MockEventRepository) *Handler {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"Nop"},
		ErrorOutputPaths: []string{"Nop"},
	}
	logger, err := logger.NewCustomLogger(zapConfig)
	if err != nil {
		panic(err.Error())
	}
	service := &Handler{
		Logger:    logger,
		Repo:      rep,
		Validator: validator.New(),
	}
	return service
}

func TestCreateEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	eventsRepo := event.NewMockEventRepository(ctrl)
	hndlr := InitiateHandlerUser(eventsRepo)

	tests := []struct {
		name         string
		method       string
		requestBody  string
		repoMock     *gomock.Call
		expected     int
		expectedBody []byte
	}{
		{
			name:        "Success",
			method:      http.MethodPost,
			requestBody: `{"user_id": 123, "data": "new_event", "date": "2024-02-19"}`,
			repoMock: eventsRepo.EXPECT().Create(event.Event{
				ID:      "",
				UserID:  123,
				RawDate: "2024-02-19",
				Date:    time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC),
				Data:    "new_event",
			}).Return(event.Event{
				ID:      "1",
				UserID:  123,
				RawDate: "2024-02-19",
				Date:    time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC),
			}, nil),
			expected: http.StatusOK,
			expectedBody: []byte(`{
				"result": {
					"event_id":"1"
					"user_id": 123,
					"data": "new_event",
					"date": "2024-02-19T00:00:00Z",

				}
			}`),
		},
		{
			name:        "Invalid Method",
			method:      http.MethodGet,
			requestBody: "",
			expected:    400,
		},
		{
			name:        "Invalid JSON",
			method:      http.MethodPost,
			requestBody: "invalid json",
			expected:    400,
		},
		{
			name:         "Repo Error",
			method:       http.MethodPost,
			requestBody:  `{"user_id": 123, "data":"sada", "date":"2024-02-19"}`,
			repoMock:     eventsRepo.EXPECT().Create(gomock.Any()).Return(event.Event{}, fmt.Errorf("bad repo")),
			expected:     503,
			expectedBody: []byte(`{"error" : "internal server error"}`),
		},
	}

	for _, tst := range tests {
		buf := bytes.NewBuffer([]byte(tst.requestBody))
		r := httptest.NewRequest(tst.method, "/create_event", buf)
		w := httptest.NewRecorder()

		hndlr.CreateEventHandler(w, r)

		if w.Code != tst.expected {
			t.Errorf("not correct code")
			if tst.expectedBody != nil {
				bd := w.Body
				bts := bd.Bytes()
				assert.Equal(t, tst.expectedBody, bts, "not equal")
			}
		}
	}
}

func TestUpdatEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	eventsRepo := event.NewMockEventRepository(ctrl)
	hndlr := InitiateHandlerUser(eventsRepo)

	tests := []struct {
		name         string
		method       string
		requestBody  string
		repoMock     *gomock.Call
		expected     int
		expectedBody []byte
	}{
		{
			name:        "Success",
			method:      http.MethodPost,
			requestBody: `{"user_id": 123, "data": "new_event", "date": "2024-02-19"}`,
			repoMock: eventsRepo.EXPECT().Update(event.Event{
				ID:      "",
				UserID:  123,
				RawDate: "2024-02-19",
				Date:    time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC),
				Data:    "new_event",
			}).Return(event.Event{
				ID:      "1",
				UserID:  123,
				RawDate: "2024-02-19",
				Date:    time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC),
			}, nil),
			expected: http.StatusOK,
			expectedBody: []byte(`{
				"result": {
					"event_id":"1"
					"user_id": 123,
					"data": "new_event",
					"date": "2024-02-19T00:00:00Z",

				}
			}`),
		},
		{
			name:        "Invalid Method",
			method:      http.MethodGet,
			requestBody: "",
			expected:    400,
		},
		{
			name:        "Invalid JSON",
			method:      http.MethodPost,
			requestBody: "invalid json",
			expected:    400,
		},
		{
			name:         "Repo Error",
			method:       http.MethodPost,
			requestBody:  `{"user_id": 123, "data":"sada", "date":"2024-02-19"}`,
			repoMock:     eventsRepo.EXPECT().Update(gomock.Any()).Return(event.Event{}, fmt.Errorf("bad repo")),
			expected:     503,
			expectedBody: []byte(`{"error" : "internal server error"}`),
		},
	}

	for _, tst := range tests {
		buf := bytes.NewBuffer([]byte(tst.requestBody))
		r := httptest.NewRequest(tst.method, "/update_event", buf)
		w := httptest.NewRecorder()

		hndlr.UpdateEventHandler(w, r)

		if w.Code != tst.expected {
			t.Errorf("not correct code")
			if tst.expectedBody != nil {
				bd := w.Body
				bts := bd.Bytes()
				assert.Equal(t, tst.expectedBody, bts, "not equal")
			}
		}
	}
}

func TestDeleteEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	eventsRepo := event.NewMockEventRepository(ctrl)
	hndlr := InitiateHandlerUser(eventsRepo)

	tests := []struct {
		name         string
		method       string
		requestBody  string
		repoMock     *gomock.Call
		expected     int
		expectedBody []byte
	}{
		{
			name:        "Success",
			method:      http.MethodPost,
			requestBody: `{"event_id":"1", "user_id": 123, "data": "new_event", "date": "2024-02-19"}`,
			repoMock:    eventsRepo.EXPECT().Delete(int64(123), "1").Return(nil).MaxTimes(1),
			expected:    http.StatusOK,
			expectedBody: []byte(`{
				"result": "OK"
			}`),
		},
		{
			name:        "Invalid Method",
			method:      http.MethodGet,
			requestBody: "",
			expected:    400,
		},
		{
			name:        "Invalid JSON",
			method:      http.MethodPost,
			requestBody: "invalid json",
			expected:    400,
		},
		{
			name:         "Repo Error",
			method:       http.MethodPost,
			requestBody:  `{"user_id": 123, "data":"sada", "date":"2024-02-19"}`,
			repoMock:     eventsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("bad repo")),
			expected:     503,
			expectedBody: []byte(`{"error" : "internal server error"}`),
		},
	}

	for _, tst := range tests {
		buf := bytes.NewBuffer([]byte(tst.requestBody))
		r := httptest.NewRequest(tst.method, "/delete_event", buf)
		w := httptest.NewRecorder()

		hndlr.DeleteEventHandler(w, r)

		if w.Code != tst.expected {
			t.Errorf("not correct code")
			if tst.expectedBody != nil {
				bd := w.Body
				bts := bd.Bytes()
				assert.Equal(t, tst.expectedBody, bts, "not equal")
			}
		}
	}
}

func TestEventsForHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	eventsRepo := event.NewMockEventRepository(ctrl)
	hndlr := InitiateHandlerUser(eventsRepo)

	baseURL := "/events_for_day"
	params := url.Values{}
	params.Set("user_id", "123")
	params.Set("date", "2024-02-19")

	tests := []struct {
		name         string
		method       string
		requestURL   string
		requestBody  string
		repoMock     *gomock.Call
		expected     int
		expectedBody []byte
	}{
		{
			name:       "Success",
			method:     http.MethodGet,
			requestURL: "/events_for_day?user_id=123&date=2024-02-19",
			repoMock: eventsRepo.EXPECT().GetEventsForRange(int64(123), time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC), time.Date(2024, time.February, 20, 0, 0, 0, 0, time.UTC)).Return(
				[]event.Event{
					{
						ID:      "1",
						UserID:  123,
						RawDate: "2024-02-19",
						Date:    time.Date(2024, time.February, 19, 0, 0, 0, 0, time.UTC),
						Data:    "new_event",
					},
				}, nil,
			),
			expected: http.StatusOK,
			expectedBody: []byte(`{
				"result": [
					{
					"event_id":"1"
					"user_id": 123,
					"data": "new_event",
					"date": "2024-02-19T00:00:00Z",
					}
				]
			}`),
		},
		{
			name:        "Invalid Method",
			method:      http.MethodPost,
			requestURL:  "/events_for_day?user_id=123&date=2024-02-19",
			requestBody: "",
			expected:    400,
		},
		{
			name:        "Invalid URL",
			method:      http.MethodGet,
			requestBody: "invalid json",
			requestURL:  "/events_for_day?",
			expected:    400,
		},
		{
			requestURL:   baseURL + "?" + params.Encode(),
			name:         "Repo Error",
			method:       http.MethodGet,
			repoMock:     eventsRepo.EXPECT().GetEventsForRange(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("bad repo")),
			expected:     503,
			expectedBody: []byte(`{"error" : "internal server error"}`),
		},
	}

	for _, tst := range tests {
		buf := bytes.NewBuffer([]byte(tst.requestBody))
		r := httptest.NewRequest(tst.method, tst.requestURL, buf)
		w := httptest.NewRecorder()

		hndlr.EventsForHandler(w, r)

		if w.Code != tst.expected {
			t.Errorf("not correct code")
			if tst.expectedBody != nil {
				bd := w.Body
				bts := bd.Bytes()
				assert.Equal(t, tst.expectedBody, bts, "not equal")
			}
		}
	}
}
