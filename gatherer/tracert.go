package gatherer

import (
	"fmt"
	"net"
	"time"

	"../logger"
	"../utils"
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
// WebSocket API.
type Tracer struct {
	mconn *utils.MuxConn
	host  string
	route []node
}

// NewTracer returns a new route tracer.
func NewTracer() *Tracer {
	return &Tracer{
		mconn: &utils.MuxConn{},
		route: make([]node, 0),
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target string}
func (t *Tracer) Set(v ...interface{}) {
	t.mconn.Conn = v[0].(*websocket.Conn)
	t.host = v[1].(string)
}

// Report implements Gatherer interface.
func (t *Tracer) Report() map[string]interface{} {
	return nil
}

// Run implements Gatherer interface.
func (t *Tracer) Run() {
	logger.Green.Println("Trace Route and GeoIP")
	t.route = []node{}
	timeout := make(chan bool)
	go func() {
		time.Sleep(20 * time.Second)
		timeout <- true
	}()

	ch := make(chan traceroute.TracerouteHop, 0)
	go func() {
		for {
			hop, ok := <-ch
			if !ok {
				return
			}
			t.printHop(hop)
		}
	}()

	_, err := traceroute.Traceroute(t.host, &traceroute.TracerouteOptions{}, ch)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	select {
	default:
		// Wait the final output.
		time.Sleep(1 * time.Second)
	case <-timeout:
		return
	}

}

func (t *Tracer) printHop(hop traceroute.TracerouteHop) {
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
		t.mconn.Send(ret)
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
	t.mconn.Send(ret)
}

func (n *node) geoip() {
	db, err := geoip2.Open("/dict/GeoLite2-City.mmdb")
	defer db.Close()
	if err != nil {
		logger.Red.Println(err)
		return
	}

	record, err := db.City(net.ParseIP(n.addr))
	if err != nil {
		logger.Red.Println(err)
		return
	}
	n.country = record.Country.Names["en"]
	n.lat = record.Location.Latitude
	n.long = record.Location.Longitude
}
