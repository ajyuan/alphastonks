package main

import (
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/sirupsen/logrus"
)

const (
	// Base clock speed: 1 tick per tickDuration in milliseconds
	tickDuration   = 512
	sleepRandRange = 512

	// Community page parsing vars
	ytTarget       = "https://www.youtube.com/c/Deadnsyde/community"
	cookieData     = "LOGIN_INFO=AFmmF2swRQIgHMWsQHQ-90wybzpWtLWiT2ZBzVTEpBFKRTY8uZdr2KcCIQCLbYFdhRci3b09nS5XwhMIcJSyTPYCcj06VInnihPG4g:QUQ3MjNmdy12NW1Hc0l5d1lwRkNJTVoyYVhwR1d1MzJKa2Q2VEc4RDRvM0xwdlE4R3FuS2hoVDRadXdPaWtZeEk1TExWdjRVVG1WcGJnMmIxdVJ4Q0JTcWt4RlhoeU16R09YVU9XX2E5Zk1RT3ZQSnpUdzFrMmI5M0Zhb2RCMTBfMjdPMG4tNjhNdWo0dGw4MWZnZFkzcXdDVGg4U0tFa1QzQTVGRm9hMHNpN3BBdWZ6Tnk2MnE0; __Secure-3PSID=5gdH74PTb42Ro1o50WJwfDou628f_muSSJ2NXUIWDT0ksniTlIQ9jnM90C7zEyoUopRlrg."
	postTextPrefix = "\"contentText\":{\"runs\":[{\"text\":\""
	postTextSuffix = "},"
	postTimePrefix = "\"publishedTimeText\":{\"runs\":[{\"text\":\""
	postTimeSuffix = "\","
)

var (
	// Testing vars
	// SET TO FALSE BEFORE PUSHING
	ignoreMarketHours = false
	ignorePostAge     = false
	benchmark         = false

	// Alpaca Config
	alpacaID        = os.Getenv("ALPACA_ID")
	alpacaSecret    = os.Getenv("ALPACA_SECRET")
	alpacaMarketURL = os.Getenv("ALPACA_MARKET_URL")

	// YT Parse Regex
	postTextRe = regexp.MustCompile("\"contentText\":{\"runs\":(\\[{.*?\\])")

	// String filters
	actionExecutableTimeFilter = []string{"hour", "day", "minute", "week", "month", "year"}
	tickerFalsePositives       = map[string]struct{}{"I": {}, "A": {}, "AI": {}, "ET": {}, "DD": {}, "DM": {}, "ML": {}, "ARK": {}, "BTC": {}, "BUT": {}, "CEO": {}, "ETF": {}, "LOT": {}, "IMO": {}, "NOT": {}, "USA": {}, "USD": {}, "LONG": {}, "VERY": {}, "COVID": {}, "SHORT": {}, "SUPER": {}, "REALLY": {}}
	abortKeywords              = map[string]struct{}{"bots": {}, "botting": {}}

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

// PastTrainingHours indicates the time is past bot trading hours
func PastTrainingHours() bool {
	return time.Now().In(nyTimezone).Hour() == 16
}

func marketHoliday(alpacaCl *alpaca.Client) bool {
	today := time.Now().In(nyTimezone).Format("2006-01-02")
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

func processArgs() {
	for _, arg := range os.Args[1:] {
		if arg == "imh" {
			ignoreMarketHours = true
		} else if arg == "ipa" {
			ignorePostAge = true
		} else if arg == "d" {
			log.SetLevel(logrus.DebugLevel)
		} else if arg == "b" {
			benchmark = true
		}
	}
}

func setup() setupOutput {
	log.SetLevel(logrus.InfoLevel)
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
	if benchmark {
		defer timer()()
	}
	post, err := YTPost(cl)
	if err != nil {
		log.Error(err)
		return nil
	}
	actions, err := Recommendation(post)
	if err != nil {
		return err
	}
	err = Execute(actions, alpacaCl)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.Infof("AlphaStonks v.%s", "2.0.1")
	setupOutput := setup()
	processArgs()
	if marketHoliday(setupOutput.alpacaCl) && !ignoreMarketHours {
		log.Infof("Today is a market holiday, shutting down")
		os.Exit(0)
	}
	log.Info("Setup complete")
	for true {
		if PastTrainingHours() && !ignoreMarketHours {
			log.Infof("Markets are closed, shutting down")
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
		sleepDuration := time.Millisecond * time.Duration(
			minZero(int(tickDuration-time.Since(tickStart).Milliseconds())+rand.Intn(sleepRandRange)))
		time.Sleep(sleepDuration)
	}
}
