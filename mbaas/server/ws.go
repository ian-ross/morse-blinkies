package server

import "github.com/gorilla/websocket"

// WSUpgrader is a protocol upgrader for making WebSockets.
var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
