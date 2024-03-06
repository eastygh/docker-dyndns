package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"main/config"
	"main/dns"
	"main/netutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

var dyndnsConfig = &config.Config{}

func updateZone(zone string, domain string, recordType string, ttl string, ip string) string {

	f, err := ioutil.TempFile("/tmp", "dyndns")
	if err != nil {
		return err.Error()
	}

	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf("server localhost\n"))
	w.WriteString(fmt.Sprintf("zone %s.\n", zone))
	w.WriteString(fmt.Sprintf("update delete %s.%s %s\n", domain, zone, recordType))
	w.WriteString(fmt.Sprintf("update add %s.%s %s %s %s\n", domain, zone, ttl, recordType, ip))
	w.WriteString("send\n")

	w.Flush()
	f.Close()

	cmd := exec.Command("/usr/bin/nsupdate", f.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err.Error() + ": " + stderr.String()
	}
	return out.String()
}

func main() {

	log.Info().Msg("Starting...")

	gin.SetMode(gin.ReleaseMode)

	dyndnsConfig.ParseConfig("./dyndnsConfig.json")

	router := gin.Default()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		dyndnsConfig.User: dyndnsConfig.Password,
	}))

	authorized.GET("/update", func(c *gin.Context) {

		domain := c.Query("domain")

		ip := c.Query("ip")
		if len(ip) < 1 {
			ip = c.ClientIP()
		}

		if netutil.IsDomainValid(domain, dyndnsConfig.Domains) {
			if netutil.ValidateIpV4(ip) {
				err := updateZone(dyndnsConfig.Zone, domain, "A", dyndnsConfig.TTL, ip)
				if err == "" {
					c.String(http.StatusOK, "Set record for domain: %s to ip: %s", domain, ip)
				} else {
					c.String(http.StatusBadRequest, "%s", err)
				}
			} else if netutil.ValidateIpV6(ip) {
				err := updateZone(dyndnsConfig.Zone, domain, "AAAA", dyndnsConfig.TTL, ip)
				if err == "" {
					c.String(http.StatusOK, "Set record for domain: %s to ip: %s", domain, ip)
				} else {
					c.String(http.StatusBadRequest, "%s", err)
				}
			} else {
				c.String(http.StatusBadRequest, "ip: %s ist not in a valid format", ip)
			}
		} else {
			c.String(http.StatusBadRequest, "subdomain: %s not allowed", domain)
		}
	})

	handler := new(dns.Handler)
	handler.StartDNSResolver()
	err := router.Run(":8080")
	if err != nil {
		return
	}

}
