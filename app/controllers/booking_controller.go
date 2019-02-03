package controllers

import (
	"go-mobiclix/app/models"
	"go-mobiclix/lib"
	"log"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	ids   = make([]uint, 0)
)

func BookingHandler(w http.ResponseWriter, _ *http.Request) {

	tx := lib.DB.Begin()
	var ticket models.Ticket

	if err := tx.Raw(`SELECT * FROM tickets WHERE tickets.is_booked = FALSE LIMIT 1 for update`).Scan(&ticket).Error; err != nil {
		tx.Rollback()
		log.Fatal(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	//tx.Not(ids).First(&ticket, "is_booked = ?", false)

	if ticket.ID == 0 {
		tx.Rollback()
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	ticket.IsBooked = true

	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
