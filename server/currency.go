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
	log           hclog.Logger
	rates         *data.ExchangeRates
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
	currency.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{log: l, rates: r, subscriptions: make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)}
	go c.handleUpdates()
	return c
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Rate: rate, Base: rr.GetBase(), Destination: rr.GetDestination()}, nil
}

func (c *Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			c.log.Debug("Map structure", c.subscriptions)
			return nil
		}

		if err != nil {
			c.log.Error("Unable to read from client", "error", err)

			return err
		}
		c.log.Info("Handle client request", "request", rr)

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*currency.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs

	}
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRate(5 * time.Second)
	for range ru {
		c.log.Info("Got Updated rates")
		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("unable to get updated rates", "base", rr.GetBase(), "destination", rr.GetDestination())
				}
				err = k.Send(&currency.RateResponse{Rate: r, Base: rr.GetBase(), Destination: rr.GetDestination()})
				if err != nil {
					c.log.Error("Unable to send updated rates", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}

func (c *Currency) deleteClientFromSubscription() {

}
