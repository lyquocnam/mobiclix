package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-mobiclix/app/controllers"
	"go-mobiclix/app/models"
	"go-mobiclix/lib"
	"log"
	"net/http"
	"time"
)

func Run() {
	lib.ConnectDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/booking", controllers.BookingHandlerV2)

	controllers.Jobs = make(chan *controllers.Job, 200)

	go Booking()

	server := http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}

	fmt.Println("-> Listening on port localhost:8080")
	log.Fatalln(server.ListenAndServe())
}

func Booking() {
	for {
		select {
		case job := <-controllers.Jobs:
			var ticket models.Ticket
			tx := lib.DB.Begin()
			tx.Raw(`SELECT * FROM tickets WHERE is_booked = ? LIMIT 1 for update`, false).Scan(&ticket)

			if ticket.ID == 0 {
				tx.Rollback()
				job.Pool <- false
				return
			}

			if err := tx.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Update("is_booked", true).Error; err != nil {
				tx.Rollback()
				job.Pool <- false
				return
			}

			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				job.Pool <- false
				return
			}

			job.Pool <- true
		}
	}
}
