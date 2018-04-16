package gatherer

import (
	"strings"

	"../logger"
	"github.com/gorilla/websocket"
	whois "github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
)

// Whois queries the domain information.
type Whois struct {
	mconn  *muxConn
	domain string
	raw    string
	info   map[string]string
}

// NewWhois returns a new Whois.
func NewWhois() *Whois {
	return &Whois{}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target string}.
func (w *Whois) Set(v ...interface{}) {
	w.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	if strings.Count(v[1].(string), ".") == 2 {
		d := strings.Split(v[0].(string), ".")
		w.domain = d[1] + "." + d[2]
		return
	}
	w.domain = v[1].(string)
}

// Report implements Gatherer interface.
func (w *Whois) Report() interface{} {
	return w.info
}

// Run implements Gatherer interface.
func (w *Whois) Run() {
	logger.Green.Println("Whois Information")
	whoisRaw, err := whois.Whois(w.domain)
	if err != nil {
		logger.Red.Println(err)
	}
	w.raw = whoisRaw
	result, _ := whois_parser.Parse(w.raw)

	ret := map[string]string{
		"domain":          w.domain,
		"registrar_name":  result.Registrar.RegistrarName,
		"admin_name":      result.Admin.Name,
		"admin_email":     result.Admin.Email,
		"admin_phone":     result.Admin.Phone,
		"created_date":    result.Registrar.CreatedDate,
		"expiration_date": result.Registrar.ExpirationDate,
		"ns":              result.Registrar.NameServers,
		"state":           strings.Split(result.Registrar.DomainStatus, " ")[0],
	}
	w.info = ret
	w.mconn.send(ret)
}
