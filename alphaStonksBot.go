package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/sirupsen/logrus"
)

const (
	// Testing vars
	// SET TO FALSE BEFORE PUSHING
	ignoreMarketHours = false
	ignorePostAge     = false

	// Base clock speed: 1 tick per tickDuration in milliseconds
	tickDuration = 512

	// Community page parsing vars
	ytTarget       = "https://www.youtube.com/c/Deadnsyde/community"
	cookieData     = "LOGIN_INFO=AFmmF2swRQIgHMWsQHQ-90wybzpWtLWiT2ZBzVTEpBFKRTY8uZdr2KcCIQCLbYFdhRci3b09nS5XwhMIcJSyTPYCcj06VInnihPG4g:QUQ3MjNmdy12NW1Hc0l5d1lwRkNJTVoyYVhwR1d1MzJKa2Q2VEc4RDRvM0xwdlE4R3FuS2hoVDRadXdPaWtZeEk1TExWdjRVVG1WcGJnMmIxdVJ4Q0JTcWt4RlhoeU16R09YVU9XX2E5Zk1RT3ZQSnpUdzFrMmI5M0Zhb2RCMTBfMjdPMG4tNjhNdWo0dGw4MWZnZFkzcXdDVGg4U0tFa1QzQTVGRm9hMHNpN3BBdWZ6Tnk2MnE0; __Secure-3PSID=5gdH74PTb42Ro1o50WJwfDou628f_muSSJ2NXUIWDT0ksniTlIQ9jnM90C7zEyoUopRlrg."
	postTextPrefix = "\"contentText\":{\"runs\":[{\"text\":\""
	postTextSuffix = "},"
	postTimePrefix = "\"publishedTimeText\":{\"runs\":[{\"text\":\""
	postTimeSuffix = "\","

	// Modifiers
	buyLowConfidence   = 0.13
	buyMedConfidence   = 0.6
	buyHighConfidence  = 0.15
	sellHighConfidence = -0.8
	lowBuyMult         = 0.65
	medBuyMult         = 5.0
	highBuyMult        = 1.0
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

	// String filters
	actionExecutableTimeFilter = []string{"hour", "day", "minute", "week", "month", "year"}
	tickerFalsePositives       = []string{"I", "A", "ET", "DD", "DM", "ARK", "CEO", "ETF", "NOT", "USD", "VERY", "SUPER", "REALLY"}

	// Time Config
	nyTimezone       *time.Location
	ytReqTimeout     = time.Second * 16
	orderTimeInForce = time.Second * 8

	log = logrus.New()
)

type setupOutput struct {
	httpCl   *http.Client
	alpacaCl *alpaca.Client
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

func marketHoliday(alpacaCl *alpaca.Client) bool {
	today := time.Now().Format("2006-01-02")
	calendar, err := alpacaCl.GetCalendar(&today, &today)
	if err != nil {
		log.Errorf("marketHoliday error occurred when retrieving calendar: %v", err)
		// Don't let a failed calendar check conclude today is not a market day
		return false
	}
	if len(calendar) != 1 {
		log.Errorf("marketHoliday returned an calendar list of len %d, assuming market open at standard hours", len(calendar))
		return false
	}
	return today != calendar[0].Date
}

func setup() setupOutput {
	log.SetLevel(logrus.DebugLevel)
	rand.Seed(time.Now().UnixNano())
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
	alpacaCl := alpaca.NewClient(common.Credentials())
	return setupOutput{
		httpCl: &http.Client{
			Timeout: ytReqTimeout,
		},
		alpacaCl: alpacaCl,
	}
}

// Tick performs all steps to do one iteration of the check & buy algo
func Tick(cl *http.Client, alpacaCl *alpaca.Client) error {
	// defer timer()()
	post, err := YTPost(cl)
	if err != nil {
		return err
	}
	action := Recommendation(post)
	err = Execute(action, alpacaCl)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Infof("AlphaStonks v.%s", "1.06")
	setupOutput := setup()
	if marketHoliday(setupOutput.alpacaCl) && !ignoreMarketHours {
		log.Infof("The time is %v and it is a market holiday, shutting down", time.Now())
		os.Exit(0)
	}
	log.Info("Setup complete")
	for true {
		if PastAH() && !ignoreMarketHours {
			log.Infof("The time is %v and markets are closed, shutting down", time.Now().In(nyTimezone))
			os.Exit(0)
		}

		tickStart := time.Now()
		err := Tick(setupOutput.httpCl, setupOutput.alpacaCl)
		if err != nil {
			if err == ErrInsufficientFunds {
			} else {
				log.Error(err)
				panic(err)
			}
		}
		if time.Since(tickStart) > time.Second*3 {
			log.Warnf("Thottling detected, last request took %v", time.Since(tickStart))
		}
		/* Removed random timeout
		sleepDuration := time.Millisecond * time.Duration(
			minZero(int(tickDuration-time.Since(tickStart).Milliseconds())+rand.Intn(sleepRandRange)))
		*/
		sleepDuration := time.Millisecond * time.Duration(
			minZero(int(tickDuration-time.Since(tickStart).Milliseconds())))
		time.Sleep(sleepDuration)
	}
}
