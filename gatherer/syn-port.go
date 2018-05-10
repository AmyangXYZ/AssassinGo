package gatherer

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"../logger"
	"../utils"
	"github.com/gorilla/websocket"
)

// PortScanner scans common used ports.
// WebSocket API.
type PortScanner struct {
	mconn       *utils.MuxConn
	handler     net.PacketConn
	target      string
	targetIP    net.IP
	srcIP       net.IP
	srcPort     layers.TCPPort
	OpenedPorts map[string]string
}

// NewPortScanner returns a PortScanner.
func NewPortScanner() *PortScanner {
	h, _ := net.ListenPacket("ip4:tcp", "0.0.0.0")
	return &PortScanner{
		mconn:       &utils.MuxConn{},
		handler:     h,
		OpenedPorts: make(map[string]string, 0),
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target}
func (ps *PortScanner) Set(v ...interface{}) {
	ps.mconn.Conn = v[0].(*websocket.Conn)
	ps.target = v[1].(string)
}

// Report implements Gatherer interface
func (ps *PortScanner) Report() map[string]interface{} {
	return map[string]interface{}{
		"ports": ps.OpenedPorts,
	}
}

// Run implements the Gatherer interface.
func (ps *PortScanner) Run() {
	logger.Green.Println("Ports Scanning...")
	dstaddrs, err := net.LookupIP(ps.target)
	if err != nil {
		logger.Red.Println(err)
		return
	}

	ps.targetIP = dstaddrs[0].To4()
	err = ps.localIPPort(ps.targetIP)

	if err != nil {
		logger.Red.Println(err)
		return
	}

	go func() {
		for p := 0; p < 65536; p++ {
			ps.send(p)
		}
	}()

	ps.handle()
}

func (ps *PortScanner) handle() {
	ps.handler, _ = net.ListenPacket("ip4:tcp", "0.0.0.0")

	ps.handler.SetDeadline(time.Now().Add(10 * time.Second))

	for {
		b := make([]byte, 4096)
		n, addr, err := ps.handler.ReadFrom(b)

		if err != nil {
			// logger.Red.Println(err)
			return
		} else if addr.String() == ps.targetIP.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == ps.srcPort {
					if tcp.SYN && tcp.ACK {
						ps.onResult(tcp.SrcPort)
					}
				}
			}
		}
	}
}

func (ps *PortScanner) onResult(openedPort layers.TCPPort) {
	var service string
	port := fmt.Sprintf("%d", openedPort)
	if _, ok := ps.OpenedPorts[port]; ok {
		return
	}
	if len(openedPort.String()) > len(port) {
		service = openedPort.String()[len(port)+1 : len(openedPort.String())-1]
	} else {
		service = "unknown"
	}
	logger.Blue.Printf("%-5s -  %s \n", port, openedPort)
	ret := map[string]interface{}{
		"port":    port,
		"service": service,
	}
	ps.mconn.Send(ret)
	ps.OpenedPorts[port] = service
}

func (ps *PortScanner) send(port int) {
	dstport := layers.TCPPort(port)

	ip := &layers.IPv4{
		SrcIP:    ps.srcIP,
		DstIP:    ps.targetIP,
		Protocol: layers.IPProtocolTCP,
	}
	// Our TCP header
	tcp := &layers.TCP{
		SrcPort: ps.srcPort,
		DstPort: dstport,
		Seq:     1105024978,
		SYN:     true,
		Window:  14600,
	}
	tcp.SetNetworkLayerForChecksum(ip)

	// Serialize.  Note:  we only serialize the TCP layer, because the
	// socket we get with net.ListenPacket wraps our data in IPv4 packets
	// already.  We do still need the IP layer to compute checksums
	// correctly, though.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		logger.Red.Println(err)
		return
	}

	if _, err := ps.handler.WriteTo(buf.Bytes(), &net.IPAddr{IP: ps.targetIP}); err != nil {
		logger.Red.Println(err)
	}
	return
}

// get the local ip and port based on our destination ip.
func (ps *PortScanner) localIPPort(dstip net.IP) error {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":12345")
	if err != nil {
		logger.Red.Fatal(err)
	}

	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			ps.srcIP = udpaddr.IP
			ps.srcPort = layers.TCPPort(udpaddr.Port)
			return nil
		}
	}
	logger.Red.Println("could not get local ip: " + err.Error())
	return err
}
