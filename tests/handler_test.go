package tests

import (
	"github.com/stretchr/testify/assert"
	"go-mobiclix/app"
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

	controllers.Jobs = make(chan *controllers.Job, 100)
	go app.Booking()

	handler := http.HandlerFunc(controllers.BookingHandlerV2)

	requests := 1000
	wg := sync.WaitGroup{}
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		}()
	}
	wg.Wait()
}
