package models

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Game struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	MaxPlayer int    `json:"max_player"`
}

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"room_name"`
}

type Participant struct {
	ID        int `json:"id"`
	AccountID int `json:"id_account"`
	RoomID    int `json:"id_room"`
}

type RoomResponse struct {
	Status string `json:"status"`
	Data   struct {
		Rooms []Room `json:"rooms"`
	} `json:"data"`
}

type RoomDetailResponse struct {
	Status string `json:"status"`
	Data   struct {
		Room struct {
			ID           int           `json:"id"`
			RoomName     string        `json:"room_name"`
			Participants []Participant `json:"participants"`
		} `json:"room"`
	} `json:"data"`
}

type InsertRoomResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LeaveRoomResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
