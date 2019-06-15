package main

import (
	"ER-Network-Calculator/ipUtils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type IndexResponse struct {
	CAddr   ipUtils.IpAddr
	BinAddr string
	Network string
	Prefix  int
	IsValid bool
	IsCidr  bool
	IsIpv6  bool
	Error   string
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

	var ip ipUtils.IpAddr
	var err error
	if strings.Count(addr, ":") >= 2 { // IPv6
		ip, err = ipUtils.ParseIpv6(addr)
		if err != nil {
			return IndexResponse{
				Error:   err.Error(),
				IsValid: false,
			}
		}
	} else { // IPv4
		ip, err = ipUtils.ParseIpv4(addr)
		if err != nil {
			return IndexResponse{
				Error:   err.Error(),
				IsValid: false,
			}
		}
	}

	return IndexResponse{
		CAddr:   ip,
		IsValid: true,
		IsCidr:  ip.IsCidrFormatted(),
		BinAddr: ip.PrintBinary(),
		Prefix:  ip.GetPrefix(),
		Network: ip.PrintNetworkAddress(),
	}
}