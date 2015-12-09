package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
	"github.com/mostlygeek/dns-ec2-addr/awsapi"
)

var ipdb *awsapi.IpDB

func handleResponse(w dns.ResponseWriter, r *dns.Msg) {
	var rr dns.RR
	m := new(dns.Msg)
	m.SetReply(r)

	// make sure the ipdb isn't there... meh, just do nothing
	if ipdb == nil {
		return
	}

	ip := ipdb.Get()
	if ip != nil {
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET, Ttl: 0,
		}
		rr.(*dns.A).A = net.ParseIP(*ip)
		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
	}

	// otherwise don't return anything (sadness, meh)
}

func init() {
	ipdb = awsapi.NewIpDB()
	fmt.Println("fetching host ip addresses")
	ipdb.Update("autopush", "autopush")
}

func main() {

	dns.HandleFunc(".", handleResponse)
	go func() {
		fmt.Println("Listening...")
		server := &dns.Server{Addr: ":8053", Net: "udp", TsigSecret: nil}
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Could not set up server: %v", err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

forever:
	for {
		select {
		case s := <-sig:
			fmt.Printf("Signal (%d) received, stopping\n", s)
			break forever
		}
	}

}
