package server

import (
	"context"

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
