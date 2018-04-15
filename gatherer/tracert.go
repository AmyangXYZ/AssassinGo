package gatherer

import (
	"fmt"
	"net"
	"time"

	"../logger"
	"github.com/aeden/traceroute"
	"github.com/gorilla/websocket"
	geoip2 "github.com/oschwald/geoip2-golang"
)

type node struct {
	addr        string
	ttl         int
	elapsedTime time.Duration
	country     string
	lat         float64
	long        float64
}

// Tracer trace route to the target.
type Tracer struct {
	host  string
	route []node
}

// NewTracer returns a new route tracer.
func NewTracer() *Tracer {
	return &Tracer{route: make([]node, 0)}
}

// Set implements Gatherer interface.
// Params should be {}
func (t *Tracer) Set(v ...interface{}) {
	t.host = v[0].(string)
}

// Report implements Gatherer interface.
func (t *Tracer) Report() interface{} {
	return ""
}

// Run implements Gatherer interface.
func (t *Tracer) Run(conn *websocket.Conn) {
	logger.Green.Println("Trace Route and GeoIP")

	ch := make(chan traceroute.TracerouteHop, 0)
	go func() {
		for {
			hop, ok := <-ch
			if !ok {
				return
			}
			t.printHop(hop, conn)
		}
	}()

	_, err := traceroute.Traceroute(t.host, &traceroute.TracerouteOptions{}, ch)
	if err != nil {
		logger.Red.Println(err)
	}

	// Wait the final output.
	time.Sleep(1 * time.Second)
}

func (t *Tracer) printHop(hop traceroute.TracerouteHop, conn *websocket.Conn) {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	if hop.Success {
		n := node{
			addr:        addr,
			ttl:         hop.TTL,
			elapsedTime: hop.ElapsedTime,
		}
		n.geoip()
		t.route = append(t.route, n)

		logger.Blue.Printf("TTL: %d; addr: %s; ElapsedTime: %s; Country: %s; Position: (%f, %f)",
			n.ttl, n.addr, n.elapsedTime, n.country, n.lat, n.long)
		ret := map[string]interface{}{
			"ttl":          n.ttl,
			"addr":         n.addr,
			"elapsed_time": n.elapsedTime,
			"country":      n.country,
			"lat":          n.lat,
			"long":         n.long,
		}
		conn.WriteJSON(ret)
		return
	}

	n := node{
		addr: addr,
		ttl:  hop.TTL,
	}
	logger.Blue.Println(n.ttl, "no reply")
	ret := map[string]interface{}{
		"ttl":          n.ttl,
		"addr":         n.addr,
		"elapsed_time": n.elapsedTime,
		"country":      n.country,
		"lat":          n.lat,
		"long":         n.long,
	}
	conn.WriteJSON(ret)
}

func (n *node) geoip() {
	db, err := geoip2.Open("./gatherer/GeoLite2-City.mmdb")
	defer db.Close()
	if err != nil {
		logger.Red.Fatal(err)
	}

	record, err := db.City(net.ParseIP(n.addr))
	if err != nil {
		logger.Red.Println(err)
	}
	n.country = record.Country.Names["en"]
	n.lat = record.Location.Latitude
	n.long = record.Location.Longitude
}
