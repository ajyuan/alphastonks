package main

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/shopspring/decimal"
)

var (
	// ErrInsufficientFunds indicates not enough funds to buy stock at calculated limit price
	ErrInsufficientFunds = fmt.Errorf("No money")
)

// IsAH determines if the current time is within the after-hours trading window
func IsAH() bool {
	currHour := time.Now().In(nyTimezone).Hour()
	if currHour >= 16 && currHour < 18 {
		return true
	} else if currHour == 9 && (time.Now().In(nyTimezone).Minute() < 30) {
		return true
	}
	return false
}

// actionPrice estimates an upper limit price that the order should be filled by
func actionPrice(alpacaCl *alpaca.Client, action *ActionProfile) (float64, error) {
	resp, err := alpacaCl.GetLastQuote(action.ticker)
	if err != nil {
		return 0, fmt.Errorf("actionPrice failed to retrieve last quote for ticker %s: %v", action.ticker, err)
	}
	lastPrice := resp.Last.AskPrice
	if lastPrice <= 25 {
		return roundPriceDown(lastPrice * 1.1), nil
	} else if lastPrice <= 50 {
		return roundPriceDown(lastPrice * 1.05), nil
	}
	return roundPriceDown(lastPrice * 1.03), nil
}

// actionValue computes the price and quantity for an action
// DO NOT USE FOR SELL ORDERS
func actionValue(alpacaCl *alpaca.Client, req *alpaca.PlaceOrderRequest, action *ActionProfile) error {
	acct, err := alpacaCl.GetAccount()
	if err != nil {
		return fmt.Errorf("actionQty failed to retrieve account details: %v", err)
	}
	orderPriceFloat, err := actionPrice(alpacaCl, action)
	if err != nil {
		return err
	}
	orderLimitPrice := decimal.NewFromFloat(orderPriceFloat)
	req.LimitPrice = &orderLimitPrice
	maxBuyableShares := acct.BuyingPower.Div(orderLimitPrice).Sub(decimal.NewFromFloat32(0.5))
	req.Qty = maxBuyableShares.Mul(decimal.NewFromFloat32(action.multiplier)).Round(0)
	return nil
}

// orderRequest generates the order request to be executed by Alpaca
func orderRequest(alpacaCl *alpaca.Client, action *ActionProfile) (*alpaca.PlaceOrderRequest, error) {
	req := alpaca.PlaceOrderRequest{
		AssetKey:    &action.ticker,
		Type:        alpaca.Limit,
		TimeInForce: alpaca.Day,
	}
	if action.action == actionBuy {
		req.Side = alpaca.Buy
	} else {
		req.Side = alpaca.Sell
	}
	err := actionValue(alpacaCl, &req, action)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func doAction(action *ActionProfile, alpacaCl *alpaca.Client) error {
	if action.action == actionNoOp {
		return nil
	}
	req, err := orderRequest(alpacaCl, action)
	if err != nil {
		return err
	}
	if req.Qty.LessThanOrEqual(decimal.Zero) {
		log.Infof("Insufficient funds to %s %s %v@%v", req.Side, *req.AssetKey, req.Qty, req.LimitPrice)
		return ErrInsufficientFunds
	}
	order, err := alpacaCl.PlaceOrder(*req)
	if err != nil {
		return fmt.Errorf("Execute failed to execute order %v: %v", *req, err)
	}
	log.Infof("Order Placed: %v", order)
	time.Sleep(orderTimeInForce)
	alpacaCl.CancelOrder(order.ID)
	return nil
}

// Execute executes a list of action profiles
// Dependencies: AlpacaAPI
func Execute(action *ActionProfile, alpacaCl *alpaca.Client) error {
	if action == nil {
		return nil
	}
	err := doAction(action, alpacaCl)
	if err != nil {
		return err
	}
	return nil
}

// goodTransationCheck checks if the transaction was computed too late to achieve ROI
func goodTransactionCheck(alpacaCl *alpaca.Client, action *ActionProfile) bool {
	/*
		barParams := alpaca.ListBarParams{
			Timeframe: "1Min",
			EndDt:     time.Now(),
		}
		bars, err := alpacaCl.GetSymbolBars(action.ticker, barParams)
	*/
	return true
}
