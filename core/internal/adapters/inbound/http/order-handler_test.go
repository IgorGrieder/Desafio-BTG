package http_test

import (
	http_internal "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/inbound/http"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
	"net/http"
	"testing"
)

func TestOrderHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		service ports.OrderService
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http_internal.NewOrderHandler(tt.service)
			h.CreateOrder(tt.w, tt.r)
		})
	}
}
