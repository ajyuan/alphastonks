package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

const (
	// Decision Thresholds
	lowBuyScoreThresh    = 0.0
	highBuyScoreThresh   = 0.1
	lowBuyMagThresh      = 0.1
	highBuyMagThresh     = 0.1
	salienceThresh       = 0.01
	fbLowBuyScoreThresh  = 0.2
	fbHighBuyScoreThresh = 0.32

	lowBuyMult       = 0.64
	highBuyMult      = 1.0
	lowFallbackMult  = 0.2
	highFallbackMult = 0.32

	// Actions
	actionNoOp = 0
	actionBuy  = 1
	actionSell = 2
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
			log.Errorf("Error processing time %s: %v", ytTimeString, err)
			return false
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
	return false
}

// isTickerBasic returns if a string is a ticker
func isTickerBasic(word string) bool {
	if len(word) > 5 || len(word) == 0 || !isUpper(word) {
		return false
	}
	if _, ok := tickerFalsePositives[word]; ok {
		return false
	}
	return true
}

func analyzeEntitySentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitySentimentResponse, error) {
	return client.AnalyzeEntitySentiment(ctx, &languagepb.AnalyzeEntitySentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

func actionableEntity(postText string) (*languagepb.Entity, error) {
	ctx := context.Background()
	langCl, err := language.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GCP NLP Client: %v", err)
	}
	sentiment, err := analyzeEntitySentiment(ctx, langCl, postText)
	if err != nil {
		log.Errorf("Failed to analyze text: %v", err)
		return nil, nil
	}

	for _, entity := range sentiment.Entities {
		if !isTickerBasic(entity.GetName()) || entity.GetSentiment().Score < lowBuyScoreThresh {
			continue
		}
		if entity.GetSalience() < salienceThresh {
			return nil, nil
		}
		return entity, nil
	}
	return nil, nil
}

// cloudActionProfile populates a stock profile with a recommendation
func cloudActionProfile(entity *languagepb.Entity) *ActionProfile {
	//fmt.Printf("Ticker: %s\nText: %s\nPositivity: %f\nNegativity: %f\nNeutral: %f\nCompound: %f\n\n", profile.ticker, postText, score.Positive, score.Negative, score.Neutral, tickerProfile.sentiment)
	if entity == nil {
		return nil
	}
	out := &ActionProfile{ticker: entity.Name}
	if entity.GetSentiment().Score >= highBuyScoreThresh && entity.Sentiment.GetMagnitude() >= highBuyMagThresh {
		out.action = actionBuy
		out.multiplier = highBuyMult
	} else if entity.GetSentiment().Score >= lowBuyScoreThresh && entity.Sentiment.GetMagnitude() >= lowBuyMagThresh {
		out.action = actionBuy
		out.multiplier = lowBuyMult
	} else {
		return nil
	}
	return out
}

// actionProfile returns a list of actions to execute based on post content
func actionProfile(postText string) (*ActionProfile, error) {
	start := time.Now()
	entity, err := actionableEntity(postText)
	if err != nil {
		return nil, err
	}
	out := cloudActionProfile(entity)
	if out == nil {
		log.Infof("Could not generate action using cloud action profiling, using fallback action profiling.")
		out = fallbackActionProfile(postText)
	}
	log.Infof("Generated action %v in %v", out, time.Since(start))
	return out, nil
}

// Recommendation analyzes a post recommend action
func Recommendation(post *YTPostDetails) (*ActionProfile, error) {
	if (!discoveredWithinBounds(post.postTime) || post.postText == "") && !ignorePostAge {
		log.Debugf("Recommendation detected last post was created at time %s, too late to be actionable. Skipping.", post.postTime)
		return nil, nil
	}
	profile, err := actionProfile(post.postText)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
