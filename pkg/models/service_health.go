package models

import (
	"sync"
	"time"
)

type ServiceHealth struct {
	Status      string          `json:"status"`
	LastSuccess *time.Time      `json:"last_success,omitempty"`
	Message     string          `json:"message,omitempty"`
	StateLock   sync.Mutex      `json:"-"`
	Send        map[string]bool `json:"-"`
}

func NewServiceHealth() ServiceHealth {
	return ServiceHealth{
		Status: "undefined",
		Send:   make(map[string]bool),
	}
}

func (h *ServiceHealth) UpdateSuccess(serial string) {
	h.StateLock.Lock()
	defer h.StateLock.Unlock()

	tm := time.Now().Round(time.Second)
	h.Send[serial] = h.Status != "ok"
	h.Status = "ok"
	h.LastSuccess = &tm
	h.Message = ""
}

func (h *ServiceHealth) UpdateError(serial string, err error) {
	h.StateLock.Lock()
	defer h.StateLock.Unlock()

	h.Status = "error"
	h.Message = err.Error()
	h.Send[serial] = true
}
