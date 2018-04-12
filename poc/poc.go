package poc

import "github.com/gorilla/websocket"

// POC just need to implements Run().
type POC interface {
	Set(...interface{})
	Run(conn *websocket.Conn)
	Report() interface{}
}

// POCMap is a poc map.
var POCMap = map[string]POC{
	"SeaCMSv654": NewSeaCMSv654(),
}
