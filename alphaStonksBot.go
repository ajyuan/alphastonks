package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/grassmudhorses/vader-go"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	// Community page parsing vars
	ytTarget       = "https://www.youtube.com/c/Deadnsyde/community"
	cookieData     = "LOGIN_INFO=AFmmF2swRQIgHMWsQHQ-90wybzpWtLWiT2ZBzVTEpBFKRTY8uZdr2KcCIQCLbYFdhRci3b09nS5XwhMIcJSyTPYCcj06VInnihPG4g:QUQ3MjNmdy12NW1Hc0l5d1lwRkNJTVoyYVhwR1d1MzJKa2Q2VEc4RDRvM0xwdlE4R3FuS2hoVDRadXdPaWtZeEk1TExWdjRVVG1WcGJnMmIxdVJ4Q0JTcWt4RlhoeU16R09YVU9XX2E5Zk1RT3ZQSnpUdzFrMmI5M0Zhb2RCMTBfMjdPMG4tNjhNdWo0dGw4MWZnZFkzcXdDVGg4U0tFa1QzQTVGRm9hMHNpN3BBdWZ6Tnk2MnE0; __Secure-3PSID=5gdH74PTb42Ro1o50WJwfDou628f_muSSJ2NXUIWDT0ksniTlIQ9jnM90C7zEyoUopRlrg."
	postTextPrefix = "\"contentText\":{\"runs\":[{\"text\":\""
	postTextSuffix = "\"}]},"
	postTimePrefix = "\"publishedTimeText\":{\"runs\":[{\"text\":\""
	postTimeSuffix = "\","

	// Modifiers
	buyLowConfidence   = 0.15
	buyMedConfidence   = 0.6
	buyHighConfidence  = 0.85
	sellHighConfidence = -0.8
	lowBuyMult         = 10
	medBuyMult         = 5
	highBuyMult        = 1
	highSellMult       = 100

	// Actions
	actionNoOp = 0
	actionBuy  = 1
	actionSell = 2
)

var (
	// Alpaca Config
	alpacaID        = os.Getenv("ALPACA_ID")
	alpacaSecret    = os.Getenv("ALPACA_SECRET")
	alpacaMarketURL = os.Getenv("ALPACA_MARKET_URL")

	// ErrInsufficientFunds indicates not enough funds to buy stock at calculated limit price
	ErrInsufficientFunds = fmt.Errorf("No money")

	// Ticker Identification Config
	tickerFalsePositives = []string{"I", "A", "ET", "DD", "CEO", "NOT", "USD", "VERY", "SUPER"}

	// Time Config
	nyTimezone *time.Location

	log = logrus.New()
)

// ActionProfile contains info to execute a market operation on a stock
type ActionProfile struct {
	ticker     string
	action     uint
	multiplier int32
}

func timer() func() {
	start := time.Now()
	return func() {
		log.Infof("Iteration completed in %v", time.Since(start))
	}
}

func stringIn(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func roundPriceDown(value float32) float64 {
	return math.Floor(float64(value*100)) / 100
}

func substr(page, prefix, suffix string) (string, error) {
	si, n, ei := strings.Index(page, prefix)+len(prefix), len(suffix), -1
	if si == -1 {
		return "", fmt.Errorf("extractPosts failed to find data prefix \"%s\"", postTextPrefix)
	}
	for i := si + 1; i < len(page)-n; i++ {
		if page[i:i+n] == suffix {
			ei = i
			break
		}
	}
	if ei == -1 {
		return "", fmt.Errorf("extractPosts failed to find data suffix \"%s\"", suffix)
	}
	return page[si:ei], nil
}

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

// PastAH indicates the time is past after-hours trading
func PastAH() bool {
	return time.Now().In(nyTimezone).Hour() == 18
}

func communityPage(cl *http.Client, target string) (string, error) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return "", fmt.Errorf("fetchPage failed to form request: %v", err)
	}
	req.Header.Set("Cookie", cookieData)
	resp, err := cl.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetchPage failed to request page: %v", err)
	}
	page, err := ioutil.ReadAll(resp.Body)
	f, _ := os.Create("com-page.html")
	defer f.Close()
	f.Write(page)
	return string(page), err
}

// resolveTickers resolves tickers if multiple candidates are found
func resolveTickers(tickers []string) (string, error) {
	//TODO: Implement ticker resolution
	return "", fmt.Errorf("resolveTickers found multiple possible tickers: %v", tickers)
}

// Ticker retrieves the ticker being discussed in a current post
func Ticker(postText string) (string, error) {
	words := strings.Fields(postText)
	tickers := []string{}
	var skipWord bool
	for _, word := range words {
		if string(word[0]) == "$" {
			word = word[1:]
		}
		if !unicode.IsLetter(rune(word[len(word)-1])) {
			word = word[:len(word)-1]
		}
		if len(word) > 5 {
			continue
		}
		if stringIn(tickerFalsePositives, word) {
			continue
		}
		for _, char := range word {
			if !unicode.IsUpper(char) || !unicode.IsLetter(char) {
				skipWord = true
				continue
			}
		}
		if !skipWord && len(word) != 0 && !stringIn(tickers, word) {
			tickers = append(tickers, word)
			continue
		}
		skipWord = false
	}
	if len(tickers) == 1 {
		return tickers[0], nil
	}
	return resolveTickers(tickers)
}

// Recommendation populates a stock profile with a recommendation
func Recommendation(profile *ActionProfile, postText string) {
	score := vader.GetSentiment(postText)
	//fmt.Printf("Ticker: %s\nText: %s\nPositivity: %f\nNegativity: %f\nNeutral: %f\nCompound: %f\n\n", profile.ticker, postText, score.Positive, score.Negative, score.Neutral, score.Compound)
	if score.Compound > buyHighConfidence {
		profile.action = actionBuy
		profile.multiplier = highBuyMult
	} else if score.Compound > buyMedConfidence {
		profile.action = actionBuy
		profile.multiplier = medBuyMult
	} else if score.Compound > buyLowConfidence {
		profile.action = actionBuy
		profile.multiplier = lowBuyMult
	} else if score.Compound < sellHighConfidence {
		profile.action = actionSell
		profile.multiplier = highSellMult
	}
}

// Action checks the YT feed, analyze new posts, recommend action
func Action() (*ActionProfile, error) {
	cl := &http.Client{}
	page, err := communityPage(cl, ytTarget)
	if err != nil {
		return nil, err
	}
	postText, err := substr(page, postTextPrefix, postTextSuffix)
	if err != nil {
		return nil, err
	}
	postTime, err := substr(page, postTimePrefix, postTimeSuffix)
	if postTime != "Just now" || postText == "" {
		log.Infof("Latest post time was %s, skipping iteration", postTime)
		return &ActionProfile{}, nil
	}
	ticker, err := Ticker(postText)
	if err != nil {
		// TODO: Better ticker error handling
		return &ActionProfile{}, nil
	}
	profile := &ActionProfile{ticker: ticker}
	Recommendation(profile, postText)
	return profile, err
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
	maxBuyableShares := acct.BuyingPower.Div(orderLimitPrice).Sub(decimal.NewFromFloat32(0.5)).Round(0)
	req.Qty = maxBuyableShares.Div(decimal.NewFromInt32(action.multiplier))
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

// Execute executes an action profile
func Execute(action *ActionProfile) error {
	if action.action == actionNoOp {
		return nil
	}
	alpacaCl := alpaca.NewClient(common.Credentials())
	req, err := orderRequest(alpacaCl, action)
	if err != nil {
		return err
	}
	if req.Qty.Equal(decimal.Zero) {
		log.Infof("Insufficient funds to %s %s %v@%v", req.Side, *req.AssetKey, req.Qty, req.LimitPrice)
		return ErrInsufficientFunds
	}
	order, err := alpacaCl.PlaceOrder(*req)
	if err != nil {
		return fmt.Errorf("Execute failed to execute order %v: %v", req, err)
	}
	fmt.Printf("Order Placed: %v", order)
	return nil
}

// Tick performs one check and potential buy
func Tick() error {
	defer timer()()
	action, err := Action()
	if err != nil {
		return err
	}
	err = Execute(action)
	if err != nil {
		return err
	}
	return nil
}

func setup() {
	log.SetLevel(logrus.DebugLevel)
	log.Debug("Establishing NY Time Offset")
	var err error
	nyTimezone, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("setup failed to establish NY time: %v", err)
	}
	log.Debug("Setting up Alpaca Client")
	if common.Credentials().ID == "" {
		os.Setenv(common.EnvApiKeyID, alpacaID)
	}
	if common.Credentials().Secret == "" {
		os.Setenv(common.EnvApiSecretKey, alpacaSecret)
	}
	alpaca.SetBaseUrl(alpacaMarketURL)
}

func main() {
	setup()
	for true {
		err := Tick()
		if err != nil {
			if err == ErrInsufficientFunds {
			} else {
				log.Error(err)
				panic(err)
			}
		}
		time.Sleep(time.Millisecond * 150)
		if IsAH() {
			log.Infof("The time is %v and markets are closed, shutting down", time.Now().In(nyTimezone))
			os.Exit(0)
		}
	}
}
