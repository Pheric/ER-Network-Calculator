package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type IndexResponse struct {
	Addr    net.IP
	BinAddr string
	Network *net.IPNet
	Prefix  int
	Netmask string
	IsValid bool
	IsCidr  bool
	IsIpv6 bool
}
func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./www/index.html", "./www/templates/header.html")
	if err != nil {
		e := fmt.Sprintf("error parsing index page: %v", err)
		if _, err = w.Write([]byte("503 Internal Server Error")); err != nil {
			e = fmt.Sprintf("error writing error to index page:\ninitial error:%s\ncurrent error: %v", e, err)
		}

		log.Println(e)
		return
	}

	if err := t.Execute(w, getIpInfo(r.URL.Query().Get("addr"))); err != nil {
		log.Printf("error executing index template: %v\n", err)
	}
}

func getIpInfo(addr string) IndexResponse {
	if addr == "" {
		return IndexResponse{}
	}

	ip, ipNet, err := net.ParseCIDR(addr)
	if err == nil {
		var netmask string
		for i, octet := range ipNet.Mask {
			netmask += strconv.FormatInt(int64(octet), 10)

			if i < len(ipNet.Mask) - 1 {
				netmask += "."
			}
		}

		isV6 := !strings.Contains(ip.String(), ".")
		binAddr := ""
		if !isV6 {
			binAddr = ipToBinary(addr)
		}
		return IndexResponse {
			Addr:    ip,
			BinAddr: binAddr,
			Network: ipNet,
			Netmask: netmask,
			IsValid: true,
			IsCidr:  true,
			IsIpv6: isV6,
		}
	}

	ip = net.ParseIP(addr)
	if ip == nil {
		return IndexResponse {
			IsValid: false,
		}
	}

	isV6 := !strings.Contains(ip.String(), ".")
	binAddr := ""
	if !isV6 {
		binAddr = ipToBinary(addr)
	}
	return IndexResponse {
		Addr:    ip,
		BinAddr: binAddr,
		IsValid: true,
		IsCidr:  false,
		IsIpv6: isV6,
	}
}

func ipToBinary(addr string) string {
	var s string

		split := strings.Split(addr, ".")
		for i, octet := range split {
			pOctet, err := strconv.Atoi(octet)
			if err != nil {
				return ""
			}

			formatted := strconv.FormatInt(int64(pOctet), 2)
			for ; 8-len(formatted) > 0; {
				formatted = "0" + formatted
			}
			s += formatted
			if i < len(split)-1 {
				s += "."
			}
		}

	return s
}