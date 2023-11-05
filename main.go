package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Config struct {
	DomainIPMap map[string]string `json:"domain_ip_map"`
}

var domainMap map[string]string

func main() {
	parseConfig()
	startDNSServer()
}

func parseConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error parsing configuration file:", err)
		return
	}

	domainMap = config.DomainIPMap
}

func startDNSServer() {
	addr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return
	}

	fmt.Println("DNS Proxy Server started")

	for {
		buffer := make([]byte, 1024*1024*1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		go handleDNSQuery(data, conn, addr)
	}
}

func handleDNSQuery(query []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	msg := new(dns.Msg)
	err := msg.Unpack(query)
	if err != nil {
		fmt.Println("Error unpacking DNS query:", err)
		return
	}

	domain := msg.Question[0].Name
	domain = strings.TrimSuffix(domain, ".")
	//fmt.Println("Domain:", domain)
	ip, ok := domainMap[domain]
	if ok {
		response := new(dns.Msg)
		response.SetReply(msg)
		response.Answer = []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{Name: dns.Fqdn(domain), Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
				A:   net.ParseIP(ip),
			},
		}

		responseBytes, err := response.Pack()
		if err != nil {
			fmt.Println("Error packing DNS response:", err)
			return
		}

		_, err = conn.WriteToUDP(responseBytes, addr)
		if err != nil {
			fmt.Println("Error sending DNS response:", err)
			return
		}
	} else {
		dnsConfig, _ := dns.ClientConfigFromFile("resolv.conf")

		dnsConn, err := net.Dial("udp", dnsConfig.Servers[0]+":53")
		if err != nil {
			fmt.Println("Error connecting to default DNS server:", err)
			return
		}
		defer dnsConn.Close()

		_, err = dnsConn.Write(query)
		if err != nil {
			fmt.Println("Error sending DNS query to default DNS server:", err)
			return
		}

		responseBytes := make([]byte, 1024*100)
		n, err := dnsConn.Read(responseBytes)
		if err != nil {
			fmt.Println("Error receiving DNS response from default DNS server:", err)
			return
		}

		_, err = conn.WriteToUDP(responseBytes[:n], addr)
		if err != nil {
			fmt.Println("Error sending DNS response:", err)
			return
		}
	}
}
