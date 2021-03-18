package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T) {
	tr, err := NewRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}
	rate, err := tr.GetRate("EUR", "INR")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rate)
}

func TestError(t *testing.T) {
	tr, err := NewRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tr.GetRate("jabfj", "kanfkafn")
	if err == nil {
		t.Log("Should error out")
	}
}
