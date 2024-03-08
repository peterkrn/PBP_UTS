package main

import (
	"fmt"
	"log"
	"net/http"
	"uts/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rooms/get_rooms", controllers.GetAllRooms).Methods("GET")
	router.HandleFunc("/rooms/get_rooms_details", controllers.GetRoomDetails).Methods("GET")
	router.HandleFunc("/rooms/insert_room", controllers.InsertRoom).Methods("POST")
	router.HandleFunc("/rooms/leave_room", controllers.LeaveRoom).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected on port :8888")
	log.Println("Connected on port :8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
