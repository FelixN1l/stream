package main

import (
	"flag"
	"log"

	"github.com/aiocloud/stream/api"
	"github.com/aiocloud/stream/dns"
	"github.com/aiocloud/stream/mitm"
)

var (
	flags struct {
		Path string
	}
)

func main() {
	flag.StringVar(&flags.Path, "c", "/etc/stream.json", "Path")
	flag.Parse()

	if err := api.Load(flags.Path); err != nil {
		log.Fatalf("[Stream][main][api.Load] %v", err)
	}

	{
		info := api.Get()
		dns.StrictDNS = info.DNS.Strict
		go dns.Listen(info.DNS.Listen)

		for i := 0; i < len(info.TCP.HTP); i++ {
			go mitm.ListenHTTP(info.TCP.HTP[i])
		}

		for i := 0; i < len(info.TCP.TLS); i++ {
			go mitm.ListenTLS(info.TCP.TLS[i])
		}
	}

	api.Boot()
}
