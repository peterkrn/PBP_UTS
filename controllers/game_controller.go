package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	models "uts/models"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode data as JSON and write to response body
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log the error
		log.Println("Error encoding JSON response:", err)
	}
}

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id, room_name FROM rooms"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var room models.Room
	var rooms []models.Room
	for rows.Next() {
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			log.Println(err)
		} else {
			rooms = append(rooms, room)
		}
	}

	response := struct {
		Status string                   `json:"status"`
		Data   map[string][]models.Room `json:"data"`
	}{
		Status: "success",
		Data: map[string][]models.Room{
			"rooms": rooms,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRoomDetails(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("id")

	db := connect()
	defer db.Close()

	query := `
        SELECT r.id, r.room_name, p.id, p.id_account, a.username
        FROM rooms r
        LEFT JOIN participants p ON r.id = p.id_room
        LEFT JOIN accounts a ON p.id_account = a.id
        WHERE r.id = ?
    `

	// Execute the query
	rows, err := db.Query(query, roomID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var room models.Room
	var participants []models.Participant
	for rows.Next() {
		var participant models.Participant
		var account models.Account
		if err := rows.Scan(&room.ID, &room.Name, &participant.ID, &participant.AccountID, &account.Username); err != nil {
			log.Println(err)
		} else {
			participants = append(participants, participant)
		}
	}

	// Prepare response
	var response models.RoomDetailResponse
	response.Status = "success"
	response.Data.Room.ID = room.ID
	response.Data.Room.RoomName = room.Name
	response.Data.Room.Participants = participants

	// Set response headers and encode response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request) {
	var participant models.Participant

	db := connect()
	defer db.Close()

	// Decode JSON data from request body
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		if err == io.EOF {
			log.Println("Empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Query to count participants in the room
	query := "SELECT COUNT(*) FROM participants WHERE id_room = ?"

	var count int
	err := db.QueryRow(query, participant.RoomID).Scan(&count)
	if err != nil {
		log.Println("Error counting participants:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Query to get max player from game
	var maxPlayer int
	err = db.QueryRow("SELECT max_player FROM games WHERE id = (SELECT id_game FROM rooms WHERE id = ?)", participant.RoomID).Scan(&maxPlayer)
	if err != nil {
		log.Println("Error getting max player:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If room is full, send fail response
	if count >= maxPlayer {
		response := models.InsertRoomResponse{
			Status:  "failed",
			Message: "Failed to insert room. Max player limit reached for the game.",
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	// If room is not full, insert participant
	insertQuery := "INSERT INTO participants (id_account, id_room) VALUES (?, ?) "
	_, err = db.Exec(insertQuery, participant.AccountID, participant.RoomID)
	if err != nil {
		log.Println("Error inserting participant:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := models.InsertRoomResponse{
		Status:  "success",
		Message: "Participant inserted successfully.",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	var participant models.Participant

	// Decode JSON data from request body
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		if err == io.EOF {
			log.Println("Empty request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := connect()
	defer db.Close()

	// Query to delete participant from room
	deleteQuery := "DELETE FROM participants WHERE id_account = ? AND id_room = ?"
	_, err := db.Exec(deleteQuery, participant.AccountID, participant.RoomID)
	if err != nil {
		log.Println("Error deleting participant from room:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := models.LeaveRoomResponse{
		Status:  "success",
		Message: "Participant left the room successfully.",
	}
	respondWithJSON(w, http.StatusOK, response)
}
