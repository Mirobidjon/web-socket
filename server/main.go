package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader  websocket.Upgrader
	userTodos map[string]*Client
	todoID    int
)

func main() {
	userTodos = make(map[string]*Client)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8082", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		clientReq := &ClientRequest{}
		err := conn.ReadJSON(clientReq)
		if err != nil {
			log.Println(err)
			return
		}
		clientResp := &ClientResponse{}
		log.Printf("Message from client : %#v \n\n", clientReq)
		if len(clientReq.Username) == 0 {
			return
		}

		var todos Todos
		switch clientReq.Type {
		case "hello":
			doLogin(clientReq.Username, conn)
			todos = getTodos(clientReq.Username)
		case "add":
			todos = addTodos(clientReq.Username, clientReq.Todo)
		case "delete":
			todos = removeTodo(clientReq.Username, clientReq.ID)
		case "toggle.done":
			todos = toggleDone(clientReq.Username, clientReq.ID)
		}

		clientResp.Todos = todos
		connections := userTodos[clientReq.Username].Connextions
		fmt.Printf("Updating %v Clients for user %v\n\n", len(connections), clientReq.Username)
		for _, c := range connections {
			err := c.WriteJSON(clientResp)
			if err != nil {
				doLogOut(clientReq.Username, c)
			}
		}
	}
}

func doLogOut(username string, c *websocket.Conn) {
	var tmp Connextions
	conn := userTodos[username].Connextions
	for _, v := range conn {
		if v != c {
			tmp = append(tmp, v)
		}
	}
	userTodos[username].Connextions = tmp
}

func doLogin(username string, c *websocket.Conn) {
	if userTodos[username] == nil {
		userTodos[username] = &Client{}
	}

	userTodos[username].Connextions = append(userTodos[username].Connextions, c)
}

func toggleDone(username string, id int) Todos {
	for index, todo := range userTodos[username].Todos {
		if todo.ID == id {
			userTodos[username].Todos[index].Done = !todo.Done
		}
	}
	return userTodos[username].Todos
}

func getTodos(username string) Todos {
	return userTodos[username].Todos
}

func addTodos(username string, todo Todo) Todos {
	todoID++
	todo.ID = todoID
	userTodos[username].Todos = append(userTodos[username].Todos, todo)
	return userTodos[username].Todos
}

func removeTodo(username string, id int) Todos {
	var tmp Todos
	todos := userTodos[username].Todos
	for _, todo := range todos {
		if todo.ID != id {
			tmp = append(tmp, todo)
		}
	}
	userTodos[username].Todos = tmp
	return userTodos[username].Todos
}
