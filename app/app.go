package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-mobiclix/app/controllers"
	"go-mobiclix/lib"
	"log"
	"net/http"
	"time"
)

func Run() {
	lib.ConnectDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/booking", controllers.BookingHandler)


	server := http.Server{
		Handler: r,
		Addr: ":8080",
		WriteTimeout: 5 * time.Minute,
		ReadTimeout: 5 * time.Minute,
	}

	fmt.Println("-> Listening on port localhost:8080")
	log.Fatalln(server.ListenAndServe())
}
