package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jaskeerat789/gRPC-tutorial/data"
	"github.com/jaskeerat789/gRPC-tutorial/protos/currency"
)

type Currency struct {
	log   hclog.Logger
	rates *data.ExchangeRates
	currency.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{log: l, rates: r}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase(), rr.GetDestination())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {

	go func() {

		for {
			rr, err := src.Recv()

			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}

			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				break
			}
			c.log.Info("Handle client request", "request", rr)
		}
	}()

	for {
		err := src.Send(&currency.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}

}
