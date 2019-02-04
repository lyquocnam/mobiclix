package controllers

import (
	"fmt"
	"go-mobiclix/app/models"
	"go-mobiclix/lib"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	ids   = make([]uint, 0)
	Jobs  chan *Job
)

type Job struct {
	Pool chan bool
}

//func selectTicket(tx *gorm.DB) (*models.Ticket, error) {
//	var ticket models.Ticket
//	if err := tx.Raw(`SELECT * FROM tickets WHERE tickets.is_booked = FALSE LIMIT 1 for update`).Scan(&ticket).Error; err != nil {
//		if strings.Contains(err.Error(), "could not serialize access due to concurrent update") {
//			return selectTicket(tx)
//		} else {
//			return nil, err
//		}
//	}
//	return &ticket, nil
//}
func BookingHandlerV2(w http.ResponseWriter, _ *http.Request) {
	job := &Job{
		Pool: make(chan bool),
	}
	defer close(job.Pool)

	Jobs <- job

	for {
		select {
		case ok := <-job.Pool:
			if ok {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
	}
}

func BookingHandlerV3(w http.ResponseWriter, _ *http.Request) {
	var ticket models.Ticket
	tx := lib.DB.Begin()
	tx.Raw(`SELECT * FROM tickets WHERE is_booked = ? LIMIT 1 for update`, false).Scan(&ticket)
	// tx.First(&ticket, "is_booked = ?", false)

	if ticket.ID == 0 {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticket.IsBooked = true
	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(ticket.ID)
}

func BookingHandlerV1(w http.ResponseWriter, _ *http.Request) {
	mutex.Lock()
	var ticket models.Ticket
	lib.DB.Not(ids).First(&ticket, "is_booked = ?", false)
	if ticket.ID == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	//if err != nil {
	//	tx.Rollback()
	//	w.WriteHeader(http.StatusNotAcceptable)
	//	return
	//}
	if ticket.ID == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	ids = append(ids, ticket.ID)
	//ids = append(ids, ticket.ID)
	mutex.Unlock()
	//if err := tx.Raw(`SELECT * FROM tickets WHERE tickets.is_booked = FALSE LIMIT 1 for update`).Scan(&ticket).Error; err != nil {
	//	if strings.Contains(err.Error(), "could not serialize access due to concurrent update") {
	//
	//	} else {
	//		tx.Rollback()
	//		log.Fatal(err)
	//		w.WriteHeader(http.StatusNotAcceptable)
	//		return
	//	}
	//
	//}
	//tx.Not(ids).First(&ticket, "is_booked = ?", false)

	ticket.IsBooked = true

	if err := lib.DB.Save(&ticket).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//if err := tx.Commit().Error; err != nil {
	//	log.Fatalln(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	w.WriteHeader(http.StatusOK)
	fmt.Println(ticket.ID)
	//mutex.Lock()
	newIds := make([]uint, 0)
	for _, id := range ids {
		if id != ticket.ID {
			newIds = append(newIds, id)
		}
	}
	ids = newIds
	//mutex.Unlock()
}
