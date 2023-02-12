package sockets

import "fmt"

type SocketResponse struct {
	MType string      `json:"type"`
	Data  interface{} `json:"data"`
}

func GetClientByUserID(userID uint) *Client {
	for client := range wsHub.clients {
		if client.userId == userID {
			return client
		}
	}
	return nil
}

func SendMessageToClient(userID uint, messageType string, message interface{}) {
	fmt.Println("Sending message to client")
	var response SocketResponse
	response.MType = messageType
	response.Data = message

	if wsHub == nil {
		fmt.Println("wsHub is nil")
		return
	}

	client := GetClientByUserID(userID)

	fmt.Println("Client: ", client, userID)

	if client != nil {
		client.mux.Lock()
		client.conn.WriteJSON(response)
		client.mux.Unlock()
	}
}
