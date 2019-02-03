package tests

import (
	"github.com/stretchr/testify/assert"
	"go-mobiclix/app/controllers"
	"go-mobiclix/lib"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func Test1kRequest(t *testing.T) {
	lib.ConnectDatabase()

	req, _ := http.NewRequest("POST", "/booking", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.BookingHandler)

	requests := 1000
	wg := sync.WaitGroup{}
	wg.Add(requests)
	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		}()
	}
	wg.Wait()
}
