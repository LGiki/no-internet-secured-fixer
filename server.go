package main

import (
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Name                    string `json:"name"`
	ActiveWebProbeHost      string `json:"activeWebProbeHost"`
	ActiveWebProbePath      string `json:"activeWebProbePath"`
	ActiveWebProbeContent   string `json:"activeWebProbeContent"`
	ActiveWebProbeHostV6    string `json:"activeWebProbeHostV6"`
	ActiveWebProbePathV6    string `json:"activeWebProbePathV6"`
	ActiveWebProbeContentV6 string `json:"activeWebProbeContentV6"`
	ActiveDnsProbeHost      string `json:"activeDnsProbeHost"`
	ActiveDnsProbeContent   string `json:"activeDnsProbeContent"`
	ActiveDnsProbeHostV6    string `json:"activeDnsProbeHostV6"`
	ActiveDnsProbeContentV6 string `json:"activeDnsProbeContentV6"`
	WebProbeLatency         time.Duration
	DnsProbeLatency         time.Duration
	WebProbeV6Latency       time.Duration
	DnsProbeV6Latency       time.Duration
	AverageLatency          time.Duration
}

func (s *Server) ToNCSIReg() *NCSIReg {
	return &NCSIReg{
		ActiveWebProbeHost:      &s.ActiveWebProbeHost,
		ActiveWebProbePath:      &s.ActiveWebProbePath,
		ActiveWebProbeContent:   &s.ActiveWebProbeContent,
		ActiveWebProbeHostV6:    &s.ActiveWebProbeHostV6,
		ActiveWebProbePathV6:    &s.ActiveWebProbePathV6,
		ActiveWebProbeContentV6: &s.ActiveWebProbeContentV6,
		ActiveDnsProbeHost:      &s.ActiveDnsProbeHost,
		ActiveDnsProbeContent:   &s.ActiveDnsProbeContent,
		ActiveDnsProbeHostV6:    &s.ActiveDnsProbeHostV6,
		ActiveDnsProbeContentV6: &s.ActiveDnsProbeContentV6,
	}
}

func testDnsProbe(host string, dnsType uint16, content string) error {
	ipType := "IPv4"
	if dnsType == dns.TypeAAAA {
		ipType = "IPv6"
	}

	dnsTypeStr := "A"
	if dnsType == dns.TypeAAAA {
		dnsTypeStr = "AAAA"
	}

	localc := new(dns.Client)
	localc.ReadTimeout = 5 * 1e9

	localm := new(dns.Msg)
	localm.RecursionDesired = true
	localm.SetQuestion(dns.Fqdn(host), dnsType)

	ra, _, err := localc.Exchange(localm, net.JoinHostPort(host, "53"))
	if ra == nil {
		return fmt.Errorf("error getting the %s address of %s: %s", ipType, host, err.Error())
	}
	if ra.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("invalid answer name %s after %s query: %s", host, dnsTypeStr, dns.RcodeToString[ra.Rcode])
	}
	if dnsType == dns.TypeA && ra.Answer[0].(*dns.A).A.String() != content {
		return fmt.Errorf("invalid content '%s', expected '%s'", ra.Answer[0].(*dns.A).A.String(), content)
	}
	if dnsType == dns.TypeAAAA && ra.Answer[0].(*dns.AAAA).AAAA.String() != content {
		return fmt.Errorf("invalid content '%s', expected '%s'", ra.Answer[0].(*dns.AAAA).AAAA.String(), content)
	}

	return nil
}

func testWebProbe(url string, content string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("HTTP status code %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	bodyString := string(bodyBytes)
	if bodyString != content {
		return fmt.Errorf("invalid response content '%s', expected '%s'", bodyString, content)
	}

	return nil
}

func (s *Server) TestWebProbeV4() error {
	err := testWebProbe(fmt.Sprintf("http://%s/%s", s.ActiveWebProbeHost, s.ActiveWebProbePath), s.ActiveWebProbeContent)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) TestWebProbeV6() error {
	err := testWebProbe(fmt.Sprintf("http://%s/%s", s.ActiveWebProbeHostV6, s.ActiveWebProbePathV6), s.ActiveWebProbeContentV6)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) TestDnsProbeV4() error {
	err := testDnsProbe(s.ActiveDnsProbeHost, dns.TypeA, s.ActiveDnsProbeContent)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) TestDnsProbeV6() error {
	err := testDnsProbe(s.ActiveDnsProbeHostV6, dns.TypeAAAA, s.ActiveDnsProbeContentV6)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Test() {
	fmt.Println("-------------------------------------------------")
	fmt.Println(fmt.Sprintf("        Testing NCSI server [%s]", s.Name))
	fmt.Println("-------------------------------------------------")

	beginTime := time.Now()
	err := s.TestWebProbeV4()
	s.WebProbeLatency = time.Since(beginTime)
	if err != nil {
		fmt.Println("Error testing web probe:", err)
		s.WebProbeLatency = -1
	} else {
		fmt.Println(fmt.Sprintf("Web probe latency: %s", s.WebProbeLatency))
	}

	beginTime = time.Now()
	err = s.TestWebProbeV6()
	s.WebProbeV6Latency = time.Since(beginTime)
	if err != nil {
		fmt.Println("Error testing web probe v6:", err)
		s.WebProbeV6Latency = -1
	} else {
		fmt.Println(fmt.Sprintf("Web probe v6 latency: %s", s.WebProbeV6Latency))
	}

	beginTime = time.Now()
	err = s.TestDnsProbeV4()
	s.DnsProbeLatency = time.Since(beginTime)
	if err != nil {
		fmt.Println("Error testing dns probe:", err)
		s.DnsProbeLatency = -1
	} else {
		fmt.Println(fmt.Sprintf("DNS probe latency: %s", s.DnsProbeLatency))
	}

	beginTime = time.Now()
	err = s.TestDnsProbeV6()
	s.DnsProbeV6Latency = time.Since(beginTime)
	if err != nil {
		fmt.Println("Error testing dns probe v6:", err)
		s.DnsProbeV6Latency = -1
	} else {
		fmt.Println(fmt.Sprintf("DNS probe v6 latency: %s", s.DnsProbeV6Latency))
	}

	if s.WebProbeLatency != -1 && s.WebProbeV6Latency != -1 && s.DnsProbeLatency != -1 && s.DnsProbeV6Latency != -1 {
		s.AverageLatency = (s.WebProbeLatency + s.WebProbeV6Latency + s.DnsProbeLatency + s.DnsProbeV6Latency) / 4
	} else {
		s.AverageLatency = -1
	}

	fmt.Println("-------------------------------------------------")
}
