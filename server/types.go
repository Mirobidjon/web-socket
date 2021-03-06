package main

import "github.com/gorilla/websocket"

type (
	Todo struct {
		ID          int    `json:"id,omitempty"`
		Description string `json:"description,omitempty"`
		Done        bool   `json:"done,omitempty"`
	}

	Todos []Todo

	ClientRequest struct {
		Username string `json:"username,omitempty"`
		Type     string `json:"type,omitempty"` // hello, add or remove
		Todo     `json:"todo,omitempty"`
		ID       int `json:"id,omitempty"`
	}

	ClientResponse struct {
		Todos `json:"todos,omitempty"`
	}

	Connextions []*websocket.Conn
	Client      struct {
		Todos
		Connextions
	}
)
