package gatherer

import (
	"fmt"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// Whois queries the domain information.
type Whois struct {
	target string
	raw    string
	info   map[string]interface{}
}

// NewWhois returns a new Whois.
func NewWhois() *Whois {
	return &Whois{
		info: make(map[string]interface{}),
	}
}

// Set implements Gatherer interface.
// Params should be {target string}.
func (w *Whois) Set(v ...interface{}) {
	if len(v) > 0 {
		if target, ok := v[0].(string); ok {
			w.target = domainutil.Domain(target)
		}
	}
}

// Report implements Gatherer interface.
func (w *Whois) Report() map[string]interface{} {
	return w.info
}

// Run implements Gatherer interface.
func (w *Whois) Run() {
	fmt.Println("Gathering Whois Information")

	whoisRaw, err := whois.Whois(w.target)
	if err != nil {
		fmt.Println(err)
	}
	w.raw = whoisRaw

	result, err := whoisparser.Parse(w.raw)
	if err != nil {
		fmt.Println(err)
	}

	w.info = map[string]interface{}{
		"domain":          w.target,
		"registrar_name":  result.Registrar.Name,
		"admin_name":      result.Administrative.Name,
		"admin_email":     result.Administrative.Email,
		"admin_phone":     result.Administrative.Phone,
		"created_date":    result.Domain.CreatedDate,
		"expiration_date": result.Domain.ExpirationDate,
		"ns":              strings.Join(result.Domain.NameServers, ", "),
		"state":           getFirstStatus(result.Domain.Status),
	}

	for k, v := range w.info {
		fmt.Printf("%s: %v\n", k, v)
	}
}

// Helper function to get the first value from a map of string slices
func getFirstValue(m map[string][]string, key string) string {
	if values, ok := m[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

// Helper function to get the first status and split it
func getFirstStatus(statuses []string) string {
	if len(statuses) > 0 {
		return strings.Split(statuses[0], " ")[0]
	}
	return ""
}
