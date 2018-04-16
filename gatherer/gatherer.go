package gatherer

import "github.com/gorilla/websocket"

// Gatherer should implement ...
type Gatherer interface {
	Set(...interface{})
	Run(conn *websocket.Conn)
	Report() interface{}
}

// Init Gatherers.
func Init() map[string]Gatherer {
	return map[string]Gatherer{
		"basicInfo": NewBasicInfo(),
		"whois":     NewWhois(),
		"tracert":   NewTracer(),
		"cms":       NewCMSDetector(),
		"port":      NewPortScanner(),
		"crawl":     NewCrawler(),
	}
}
