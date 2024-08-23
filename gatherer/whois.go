package gatherer

import (
	"strings"

	"assassingo/logger"

	"github.com/bobesa/go-domain-util/domainutil"
	whois "github.com/likexian/whois"
	"github.com/likexian/whois-parser"
)

// Whois queries the domain information.
type Whois struct {
	target string
	raw    string
	info   map[string]interface{}
}

// NewWhois returns a new Whois.
func NewWhois() *Whois {
	return &Whois{}
}

// Set implements Gatherer interface.
// Params should be {target string}.
func (w *Whois) Set(v ...interface{}) {
	w.target = domainutil.Domain(v[0].(string))
}

// Report implements Gatherer interface.
func (w *Whois) Report() map[string]interface{} {
	return w.info
}

// Run implements Gatherer interface.
func (w *Whois) Run() {
	logger.Green.Println("Whois Information")

	whoisRaw, err := whois.Whois(w.target)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	w.raw = whoisRaw
	result, _ := whois_parser.Parse(w.raw)

	w.info = map[string]interface{}{
		"domain":          w.target,
		"registrar_name":  result.Registrar.RegistrarName,
		"admin_name":      result.Admin.Name,
		"admin_email":     result.Admin.Email,
		"admin_phone":     result.Admin.Phone,
		"created_date":    result.Registrar.CreatedDate,
		"expiration_date": result.Registrar.ExpirationDate,
		"ns":              result.Registrar.NameServers,
		"state":           strings.Split(result.Registrar.DomainStatus, " ")[0],
	}
	for k, v := range w.info {
		logger.Blue.Println(k + ": " + v.(string))
	}
}
