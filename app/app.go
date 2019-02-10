package app

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"go-mobiclix/app/controllers"
	"go-mobiclix/app/models"
	"go-mobiclix/lib"
	"log"
	"net/http"
	"time"
)

var (
	maxWorkers   = flag.Int("max_workers", 5, "The number of workers to resolve the ticket")
	maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
	port         = flag.Int("port", 8080, "The server port")
)

func Run() {
	lib.ConnectDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/booking", controllers.BookingHandlerV2)

	controllers.Jobs = make(chan *controllers.Job, 500)

	go Booking()

	server := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(`:%v`, port),
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}

	fmt.Printf(`-> Listening on port localhost:%v`, port)
	log.Fatalln(server.ListenAndServe())
}

func Booking() {
	for {
		select {
		case job := <-controllers.Jobs:
			go func(j *controllers.Job) {
				var ticket models.Ticket
				tx := lib.DB.Begin()
				tx.Raw(`SELECT * FROM tickets WHERE is_booked = ? LIMIT 1 for update`, false).Scan(&ticket)

				if ticket.ID == 0 {
					tx.Rollback()
					j.Pool <- false
					return
				}

				if err := tx.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Update("is_booked", true).Error; err != nil {
					tx.Rollback()
					j.Pool <- false
					return
				}

				if err := tx.Commit().Error; err != nil {
					tx.Rollback()
					j.Pool <- false
					return
				}

				fmt.Println(ticket.ID)

				j.Pool <- true
			}(job)
		}
	}
}
