package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}
	err := er.getRates()

	return er, err
}

func (er *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := er.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}
	dr, ok := er.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}

	return dr / br, nil
}

func (er *ExchangeRates) MonitorRate(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for k, v := range er.rates {
					change := rand.Float64() / 10
					direction := rand.Intn(1)

					if direction == 0 {
						change = 1 - change
					} else {
						change = 1 + change
					}

					er.rates[k] = v * change
				}
				ret <- struct{}{}
			}
		}
	}()
	return ret
}

func (er *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	// Parse string rates into Float54
	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		er.rates[c.Currency] = r
	}
	er.rates["EUR"] = 1
	return nil
}
