package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"ER-Network-Calculator/ipUtils"
)

type IndexResponse struct {
	Addr    net.IP
	CAddr   ipUtils.Ipv4Addr
	BinAddr string
	Network *net.IPNet
	Prefix  int
	Netmask string
	IsValid bool
	IsCidr  bool
	IsIpv6  bool
	Error string
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

	ip, err := ipUtils.ParseIpv4(addr)
	if err != nil {
		return IndexResponse{
			Error: err.Error(),
			IsValid: false,
		}
	}

	return IndexResponse{
		CAddr: ip,
		IsValid: true,
		IsCidr: ip.IsCidrFormatted(),
		BinAddr: ip.PrintBinary(),
		Prefix: ip[4],
	}
}
