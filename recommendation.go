package main

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/grassmudhorses/vader-go"
	"github.com/jdkato/prose/v2"
)

// TickerProfile contains the sentiment for a ticker
type TickerProfile struct {
	sentiment float64
	count     int
}

// ActionProfile contains info to execute a market operation on a stock
type ActionProfile struct {
	ticker     string
	action     uint
	multiplier float32
}

// discoveredWithinBounds determines if the current time is within execution time parameters
func discoveredWithinBounds(ytTimeString string) bool {
	if len(ytTimeString) < 9 {
		log.Errorf("discoveredWitinBounds error: Post discovery time %s is less than 9 characters", ytTimeString)
		return true
	}
	if ytTimeString[1:8] == " second" {
		age, err := strconv.Atoi(string(ytTimeString[0]))
		if err != nil {
			log.Errorf("Unknown time %s", ytTimeString)
		}
		if age <= 2 {
			return true
		}
	}
	for _, filterWord := range actionExecutableTimeFilter {
		if strings.Contains(ytTimeString, filterWord) {
			return false
		}
	}
	log.Errorf("Unknown time %s", ytTimeString)
	return false
}

// isTickerBasic returns if a string is a ticker
func isTickerBasic(word string) bool {
	if string(word[0]) == "$" {
		word = word[1:]
	}
	if !unicode.IsLetter(rune(word[len(word)-1])) {
		word = word[:len(word)-1]
	}
	if len(word) > 5 || len(word) == 0 || !isUpper(word) {
		return false
	}
	if _, ok := tickerFalsePositives[word]; ok {
		return false
	}
	return true
}

func nerTickersIdx(postText string, tickerIdxs map[string][]int) (map[string][]int, error) {
	doc, err := prose.NewDocument(postText)
	if err != nil {
		return nil, err
	}
	nerEntities := doc.Entities()
	if len(nerEntities) == 0 {
		return tickerIdxs, nil
	}
	nerTickerIdxs := map[string][]int{}
	gpeFound := false
	for _, ent := range doc.Entities() {
		if ent.Label == "GPE" {
			if idxs, ok := tickerIdxs[ent.Text]; ok {
				nerTickerIdxs[ent.Text] = idxs
				gpeFound = true
			}
		}
	}
	if gpeFound {
		return nerTickerIdxs, nil
	}
	for _, ent := range doc.Entities() {
		if idxs, ok := tickerIdxs[ent.Text]; ok {
			nerTickerIdxs[ent.Text] = idxs
		}
	}
	return nerTickerIdxs, nil
}

// tickerIdxs returns a mapping of ticker to indexes where it occurs
func tickerIdxs(postText string, lines []string) map[string][]int {
	tickerIdxs := map[string][]int{}
	for i, line := range lines {
		var currLineTicker string
		words := strings.Split(line, " ")
		for _, word := range words {
			if len(word) == 0 {
				continue
			}
			if isTickerBasic(word) {
				// Throw out line if ticker conflict
				// TODO: Analyze which ticker is more liked
				if currLineTicker != "" {
					log.Warnf("Ticker conflict for sentence \"line\", candidates: %s and %s", currLineTicker, word)
					continue
				}
				currLineTicker = word
			}
		}
		if currLineTicker != "" {
			tickerIdxs[currLineTicker] = append(tickerIdxs[currLineTicker], i)
		}
	}
	if len(tickerIdxs) > 1 {
		nerTickers, err := nerTickersIdx(postText, tickerIdxs)
		if err != nil {
			log.Error(err)
		} else {
			tickerIdxs = nerTickers
		}
	}
	return tickerIdxs
}

// actionProfile returns a mapping of possible tickers to their rated sentiment
func actionProfile(postText string) *ActionProfile {
	lines := strings.Split(postText, ". ")
	tickerIdxs := tickerIdxs(postText, lines)
	out := ActionProfile{}

	var topScore float64
	for ticker, idxs := range tickerIdxs {
		var currScore float64
		extraLinesRead := 0.0
		for _, idx := range idxs {
			currScore += (vader.GetSentiment(lines[idx]).Compound)
			j := idx + 1
			for j < len(lines) && !intIn(idxs, j) {
				extraSent := vader.GetSentiment(lines[j]).Compound
				currScore += extraSent
				extraLinesRead += 1.0
				if math.Abs(extraSent) >= 0.2 {
					break
				}
				j++
			}
		}
		currScore = currScore / (float64(len(idxs)) + extraLinesRead)

		if currScore > topScore {
			// TODO: Mkt cap based tiebreaking for equal score (ex. tickers mentioned once, in same sentence)
			out.ticker = ticker
		}
		actionWeight(&out, currScore)
	}
	return &out
}

// actionWeight populates a stock profile with a recommendation
func actionWeight(profile *ActionProfile, sentiment float64) {
	//fmt.Printf("Ticker: %s\nText: %s\nPositivity: %f\nNegativity: %f\nNeutral: %f\nCompound: %f\n\n", profile.ticker, postText, score.Positive, score.Negative, score.Neutral, tickerProfile.sentiment)
	if sentiment > buyHighConfidence {
		profile.action = actionBuy
		profile.multiplier = highBuyMult
	} else if sentiment > buyLowConfidence {
		profile.action = actionBuy
		profile.multiplier = lowBuyMult
	}
}

// Recommendation analyzes a post recommend action
func Recommendation(post *YTPostDetails) *ActionProfile {
	if (!discoveredWithinBounds(post.postTime) || post.postText == "") && !ignorePostAge {
		log.Debugf("Last post was created at time %s, too late to be actionable. Skipping.", post.postTime)
		return &ActionProfile{}
	}
	profile := actionProfile(post.postText)
	return profile
}
