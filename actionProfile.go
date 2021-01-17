package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/grassmudhorses/vader-go"
)

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

// Action analyzes a post recommend action
// Dependencies: YouTube
func Action(post *YTPostDetails) (*ActionProfile, error) {
	if (!discoveredWithinBounds(post.postTime) || post.postText == "") && !ignorePostAge {
		log.Debugf("Last post was created at time %s, too late to be actionable. Skipping.", post.postTime)
		return &ActionProfile{}, nil
	}
	ticker, err := Ticker(post.postText)
	if err != nil {
		// TODO: Better ticker error handling
		return &ActionProfile{}, nil
	}
	profile := &ActionProfile{ticker: ticker}
	Recommendation(profile, post.postText)
	return profile, err
}
