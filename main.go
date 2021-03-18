package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/jaskeerat789/gRPC-tutorial/data"
	"github.com/jaskeerat789/gRPC-tutorial/protos/currency"
	"github.com/jaskeerat789/gRPC-tutorial/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()

	r, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to fetch Rates from the central bank", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(r, log)
	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)

}
