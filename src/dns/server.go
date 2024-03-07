package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	_ "net"
)

type Handler struct {
	server *dns.Server
}

func (h *Handler) resolve(domain string, qType uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qType)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
	}

	for _, ans := range in.Answer {
		fmt.Println(ans)
	}
	return in.Answer
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		fmt.Printf("Received query: %s %d %d\n", question.Name, question.Qtype, question.Qclass)

		answers := h.resolve(question.Name, question.Qtype)
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
}

func (h *Handler) StartDNSResolver() {
	h.server = &dns.Server{
		Addr:      ":53",
		Net:       "udp",
		Handler:   h,
		UDPSize:   65535,
		ReusePort: true,
		//DisableBackground: true,
	}
	go func() {
		log.Info().Msgf("Starting DNS server on port %d", 53)
		err := h.server.ListenAndServe()
		fmt.Println("Stopped dns server")
		if err != nil {
			return
		}
	}()
}
