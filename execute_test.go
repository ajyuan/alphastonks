package main

import (
	"os"
	"testing"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
)

var (
	testAlpacaID        = "PKV5AKRYU0GCJQHMZR0B"
	testAlpacaSecret    = "Njlg16OvFzrW4SATteS4oPi42EJySP9xvlyaiIjf"
	testAlpacaMarketURL = "https://paper-api.alpaca.markets"
)

func alpacaCl() *alpaca.Client {
	log.Debug("Setting up Alpaca Client")
	if common.Credentials().ID == "" {
		os.Setenv(common.EnvApiKeyID, alpacaID)
	}
	if common.Credentials().Secret == "" {
		os.Setenv(common.EnvApiSecretKey, alpacaSecret)
	}
	alpaca.SetBaseUrl(alpacaMarketURL)
	return alpaca.NewClient(common.Credentials())
}

func TestExecute(t *testing.T) {
	testCl := alpacaCl()
	tests := map[string]struct {
		input *ActionProfile
	}{
		"CRNT Execute": {
			input: &ActionProfile{
				ticker:     "CRNT",
				action:     actionBuy,
				multiplier: highBuyMult,
			},
		},
	}
	for name, test := range tests {
		err := Execute(test.input, testCl)
		if err != nil {
			log.Fatalf("test \"%s\" failed, got error %v", name, err)
		}
	}
}
