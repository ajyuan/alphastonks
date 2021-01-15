package main

import (
	"testing"
)

func TestTicker(t *testing.T) {
	tests := map[string]struct {
		input          string
		shouldErr      bool
		expectedTicker string
	}{
		"Basic CRNT": {
			input:          "New penny stock that I am currently interested in (& invested in) is CRNT (Ceragon Networks).  Ark Invest recently purchased over 150k more shares of this penny stock from 1/11 - 1/14 (i took a picture for proof and will compare Cathie Wood's portfolios in the video).  This is a 5G play which is 1 of the biggest growth sectors but what made me really excited about this stock is that Huawei is getting banned in many countries + companies don't want to work with them because of their reputation of stealing.  CRNT said they feel like they're in a good position to take market share from Huawei (I will show all of their statements in the video that I plan on releasing sometime next week). Needham upgraded CRNT and gave it a buy rating after their presentation at its growth conference yesterday giving me more confidence in this company. If you want to help me with my DD in this company, please feel free to add notes in this comment section or my subreddit (trakstocks).  I will try to go over this briefly on Sunday and then do a full video sometime next week (maybe Monday but not positive)",
			expectedTicker: "CRNT",
		},
		"Basic CPSH": {
			input:          "Hey, I just invested in the penny stock CPSH.  I will be doing a video tomorrow on this stock.  I believe this company can do well because their products are revolving around hybrid / electric vehicles, aerospace, wind turbines / clean energy.  The reason why I think this won't be a penny stock for long is because of the role they play in the clean energy sector.  Also the space industry is growing and experts believe the space industry will triple to 1.4 trillion within a decade.  Yesterday on Reddit's home page, they were showing Nasa's new Mars rover (Perseverance Rover) and how it's about to land on mars on Feb. 18th.  Nasa is using CPSH's product in that rover.  Last quarter they were net profitable for their operations and with them playing an important role in the green energy / ev sector, I think this might be the beginning of something great for this company.",
			expectedTicker: "CPSH",
		}, "Basic update, mentions CPSH": {
			input:          "A lot of people are DMing me asking for an \"exclusive option within the exclusive option\" so you can see the stock picks faster.  I will think of an option but in the meantime I have a penny stock that I am going to announce in here in a couple hours ( around 9:30am - 9:45am ET).  After this announcement, then we can figure out the other pricing options.  Please always do your research before investing because these cheap stocks are cheap for a reason (they could be bad stocks).  Remember:  CPSH has the potential to do very well in a lot of fast growing industries so pls dont compare CPSH stock performance with this new stock pick because they're different types of companies.  With that being said, I am not just talking about penny stocks just to talk about them... there always has to be at least 3 big reasons for me to want to invest in these risky stocks and there is 1 reason that does not make this your ordinary penny stock.  Stay tuned.",
			expectedTicker: "CPSH",
		}, "MP Basic": {
			input:          "hey, planning on doing a video on MP tomorrow.  Might be Wednesday latest but shooting for Tuesday.  I think they're going to do well in upcoming earnings and I feel like I can come pretty close to proving it. I could always be wrong so when you watch the video, if you think I am \"seeing what I want to see\" just ignore it but with Biden + what they're already doing, i think they're going to do well.  PLEASE tell me what you don't like with MP and i will try to address it but watch the video I attached where I addressed common concerns with the company.",
			expectedTicker: "MP",
		}, "STPK Buy": {
			input:          "Hey, STPK I am hearing great things about.  I am working on a video but don't know when it will be done.  I am hearing a lot about how this might do really well under Biden.  If you have any concerns about this company please leave it down below so I can try to address them in the video.",
			expectedTicker: "STPK",
		}, "BNGO NoOp": {
			input:          "PLEASE do your research because I dont want this BNGO to hit $20 because I am missing something!  With that being said, I just got off the phone and am less bullish with BNGO.  VERY hard subject to tackle so please do your research carefully, Jan 11th - 15th could still be a catalyst.  From my understanding there is 2 ways to look at this stock, from a \"clinical market\" side and as a \"research market\" side.  Clinical market has more of a TAM then Research Market (it's a night and day difference apparently).  I am under the impression BNGO currently can only address a very small part of \"clinical market\" and more so a \"research market\" company.",
			expectedTicker: "BNGO",
		}, "BABA NoOp": {
			input:          "Hey, I did a video on BABA a few days ago but after reading about what the Chinese government does to these \"outspoken\" billionaires, I would think waiting before you get in, might be best (in case there are any \"forced confessions\" from Jack Ma like what has happened in the past with other companies).  Will do a video on this soon.  This stock has one of the most \"buy\" recommendations from analysts I have seen so that's why this interests me so much.",
			expectedTicker: "BABA",
		}, "BNGO High Buy": {
			input:          "Ok after closer look I have invested in BNGO.  Will try to post a video this week (maybe wed or thursday).  Main reason why is because of the study published 6 days ago where they show ways BNGO ($1 a share) is better than PACB ($24 a share) and the day after the study was published, Ark Invest was trying to learn more about BNGO from their CEO.  If the BNGO share price performs even 1/4th as good as PACB I will be VERY happy lol.",
			expectedTicker: "BNGO",
		}, "BNGO/PACB Conflict": {
			input:          "I talked about PACB on youtube when it was $6 (now $24) because of Ark putting them on my radar.  While BNGO is still at $1, pls read this study to determine if you think BNGO is better than PACB.  Trying my best to not influence your decision lol.  Maybe I am reading this wrong, but it seems like BNGO is better in a couple ways. https://apnews.com/press-release/glob...",
			shouldErr:      true,
			expectedTicker: "BNGO",
		}, "FRSX Buy": {
			input:          "Just purchased the penny stock FRSX.  I covered them before but after apple \"icar\" news I have invested heavily.  Will try to post video in a few hours.  Not saying they will work with apple, just noticing growing interest in this sector and this stock seems to be the last \"autonomous driving related\" related stock left and it's under $3 a share.  They are getting more and more opportunities, expecting a grant of approximately one million USD from the European Commission, won Edison award (something Elon Musk won a while back), a COVID-19 Symptom Detection play, and will cover more in video.",
			expectedTicker: "FRSX",
		}, "IMMR Buy": {
			input:          "Hey, currently I am not invested in this but working on a video about IMMR.  They could do well soon.  Gathering up a lot of interesting information about this company on why it might be a good stock to invest in but will not get the video done until maybe next week.  I think it's def worth it for you to look closely at this company.  If you look at it closely and do NOT like it, please leave reasons why so we can see : )",
			expectedTicker: "IMMR",
		}, "SAMA Buy": {
			input:          "Hey, take a look at SAMA.  They will be merging with Clever Leaves on the 17th.  Clever Leaves is a \"Global Company and leading vertically integrated producer of  medical cannabis and hemp extracts\".  They recently signed a deal with Canopy Growth and export their product from Columbia.  They are GMP certified which apparently gives them a big advantage to sell medical grade cannabis and sell at a higher cost.  Check it out.  It's down -5% currently.",
			expectedTicker: "SAMA",
		}, "BFT Buy": {
			input:          "Very briefly looked into BFT (Paysafe) and from what I can see, it looks like a great buy.  I just now bought shares, quite a bit because I think this will have a lot of momentum.  I can see this hitting at least $15 sometime this week or next week.  Doing more research now but this is looking really good (from what I can currently see).  I would not worry about the recent run up because I think it's small when looking at the bigger picture.",
			expectedTicker: "BFT",
		}, "FEAC Buy": {
			input:          "Sorry for the repetitive videos about FEAC but after looking even closer at the company, I am even more bullish on this one.  99% positive this stock will double (not sure when but confident it will eventually,.  It has ran up recently so could be a pull back).  I will do a video today about more reasons why because of new things I discovered.  I have been trying to identify every negative thing with the company in the past 2 videos and the biggest negative catalyst I can imagine is some regulation law comes into place banning this \"skills based\" gambling.  Also I believe App store is going from 30% fee to 15% fee on in app purchases (which is better than before but I dont like how apple has so much control).  Not 100% sure this effects skillz (i am pretending that it does) but if apple raises the price again, could be another negative catalyst potentially.  Other than that, I am struggling to find more problems with the company.  Bullish and this is a long term play.  So far REALLY liking this ceo too.  Seems focused and what he says makes sense to me.",
			expectedTicker: "FEAC",
		}, "NNOX Buy": {
			input:          "Just bought more NNOX (like.. a lot more).  If it fails at live demonstration (which is tomorrow), stock will drop VERY hard so big risk.  Big catalyst tomorrow (maybe positive maybe negative).  This isn't like some drug trial btw lol.  This is like, turning on a machine... im sure they've been doing this every day this month to be prepared.",
			expectedTicker: "NNOX",
		},
	}
	for name, test := range tests {
		ticker, err := Ticker(test.input)
		if (err != nil) != test.shouldErr {
			if test.shouldErr {
				t.Errorf("test \"%s\" failed: expected error, got nil", name)
			} else if !test.shouldErr {
				t.Errorf("test \"%s\" failed: expected no error, got %v", name, err)
			}
		} else if err == nil {
			if ticker != test.expectedTicker {
				t.Errorf("test \"%s\" failed: expected action %s, got %s", name, test.expectedTicker, ticker)
			}
		}
	}
}

func TestRecommendation(t *testing.T) {
	tests := map[string]struct {
		input              string
		expectedAction     uint
		expectedMultiplier int32
	}{
		"CRNT Buy": {
			input:              "New penny stock that I am currently interested in (& invested in) is CRNT (Ceragon Networks).  Ark Invest recently purchased over 150k more shares of this penny stock from 1/11 - 1/14 (i took a picture for proof and will compare Cathie Wood's portfolios in the video).  This is a 5G play which is 1 of the biggest growth sectors but what made me really excited about this stock is that Huawei is getting banned in many countries + companies don't want to work with them because of their reputation of stealing.  CRNT said they feel like they're in a good position to take market share from Huawei (I will show all of their statements in the video that I plan on releasing sometime next week). Needham upgraded CRNT and gave it a buy rating after their presentation at its growth conference yesterday giving me more confidence in this company. If you want to help me with my DD in this company, please feel free to add notes in this comment section or my subreddit (trakstocks).  I will try to go over this briefly on Sunday and then do a full video sometime next week (maybe Monday but not positive)",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		},
		"CPSH Buy": {
			input:              "Hey, I just invested in the penny stock CPSH.  I will be doing a video tomorrow on this stock.  I believe this company can do well because their products are revolving around hybrid / electric vehicles, aerospace, wind turbines / clean energy.  The reason why I think this won't be a penny stock for long is because of the role they play in the clean energy sector.  Also the space industry is growing and experts believe the space industry will triple to 1.4 trillion within a decade.  Yesterday on Reddit's home page, they were showing Nasa's new Mars rover (Perseverance Rover) and how it's about to land on mars on Feb. 18th.  Nasa is using CPSH's product in that rover.  Last quarter they were net profitable for their operations and with them playing an important role in the green energy / ev sector, I think this might be the beginning of something great for this company.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "MP Buy/NoOp": {
			input:          "hey, planning on doing a video on MP tomorrow.  Might be Wednesday latest but shooting for Tuesday.  I think they're going to do well in upcoming earnings and I feel like I can come pretty close to proving it. I could always be wrong so when you watch the video, if you think I am \"seeing what I want to see\" just ignore it but with Biden + what they're already doing, i think they're going to do well.  PLEASE tell me what you don't like with MP and i will try to address it but watch the video I attached where I addressed common concerns with the company.",
			expectedAction: actionNoOp,
		}, "STPKK Buy": {
			input:              "Hey, STPK I am hearing great things about.  I am working on a video but don't know when it will be done.  I am hearing a lot about how this might do really well under Biden.  If you have any concerns about this company please leave it down below so I can try to address them in the video.",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		}, "BNGO NoOp": {
			input:          "PLEASE do your research because I dont want this BNGO to hit $20 because I am missing something!  With that being said, I just got off the phone and am less bullish with BNGO.  VERY hard subject to tackle so please do your research carefully, Jan 11th - 15th could still be a catalyst.  From my understanding there is 2 ways to look at this stock, from a \"clinical market\" side and as a \"research market\" side.  Clinical market has more of a TAM then Research Market (it's a night and day difference apparently).  I am under the impression BNGO currently can only address a very small part of \"clinical market\" and more so a \"research market\" company.",
			expectedAction: actionNoOp,
		}, "BABA NoOp": {
			input:          "Hey, I did a video on BABA a few days ago but after reading about what the Chinese government does to these \"outspoken\" billionaires, I would think waiting before you get in, might be best (in case there are any \"forced confessions\" from Jack Ma like what has happened in the past with other companies).  Will do a video on this soon.  This stock has one of the most \"buy\" recommendations from analysts I have seen so that's why this interests me so much.",
			expectedAction: actionNoOp,
		}, "BNGO High Buy": {
			input:              "Ok after closer look I have invested in BNGO.  Will try to post a video this week (maybe wed or thursday).  Main reason why is because of the study published 6 days ago where they show ways BNGO ($1 a share) is better than PACB ($24 a share) and the day after the study was published, Ark Invest was trying to learn more about BNGO from their CEO.  If the BNGO share price performs even 1/4th as good as PACB I will be VERY happy lol.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "BNGO Med Buy": {
			input:              "I talked about PACB on youtube when it was $6 (now $24) because of Ark putting them on my radar.  While BNGO is still at $1, pls read this study to determine if you think BNGO is better than PACB.  Trying my best to not influence your decision lol.  Maybe I am reading this wrong, but it seems like BNGO is better in a couple ways. https://apnews.com/press-release/glob...",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		}, "FRSX Buy": {
			input:              "Just purchased the penny stock FRSX.  I covered them before but after apple \"icar\" news I have invested heavily.  Will try to post video in a few hours.  Not saying they will work with apple, just noticing growing interest in this sector and this stock seems to be the last \"autonomous driving related\" related stock left and it's under $3 a share.  They are getting more and more opportunities, expecting a grant of approximately one million USD from the European Commission, won Edison award (something Elon Musk won a while back), a COVID-19 Symptom Detection play, and will cover more in video.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "IMMR Buy": {
			input:              "Hey, currently I am not invested in this but working on a video about IMMR.  They could do well soon.  Gathering up a lot of interesting information about this company on why it might be a good stock to invest in but will not get the video done until maybe next week.  I think it's def worth it for you to look closely at this company.  If you look at it closely and do NOT like it, please leave reasons why so we can see : )",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "SAMA Buy": {
			input:              "Hey, take a look at SAMA.  They will be merging with Clever Leaves on the 17th.  Clever Leaves is a \"Global Company and leading vertically integrated producer of  medical cannabis and hemp extracts\".  They recently signed a deal with Canopy Growth and export their product from Columbia.  They are GMP certified which apparently gives them a big advantage to sell medical grade cannabis and sell at a higher cost.  Check it out.  It's down -5% currently.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "BFT Buy": {
			input:              "Very briefly looked into BFT (Paysafe) and from what I can see, it looks like a great buy.  I just now bought shares, quite a bit because I think this will have a lot of momentum.  I can see this hitting at least $15 sometime this week or next week.  Doing more research now but this is looking really good (from what I can currently see).  I would not worry about the recent run up because I think it's small when looking at the bigger picture.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "FEAC Buy": {
			input:              "Sorry for the repetitive videos about FEAC but after looking even closer at the company, I am even more bullish on this one.  99% positive this stock will double (not sure when but confident it will eventually,.  It has ran up recently so could be a pull back).  I will do a video today about more reasons why because of new things I discovered.  I have been trying to identify every negative thing with the company in the past 2 videos and the biggest negative catalyst I can imagine is some regulation law comes into place banning this \"skills based\" gambling.  Also I believe App store is going from 30% fee to 15% fee on in app purchases (which is better than before but I dont like how apple has so much control).  Not 100% sure this effects skillz (i am pretending that it does) but if apple raises the price again, could be another negative catalyst potentially.  Other than that, I am struggling to find more problems with the company.  Bullish and this is a long term play.  So far REALLY liking this ceo too.  Seems focused and what he says makes sense to me.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "NNOX Buy": {
			input:              "Just bought more NNOX (like.. a lot more).  If it fails at live demonstration (which is tomorrow), stock will drop VERY hard so big risk.  Big catalyst tomorrow (maybe positive maybe negative).  This isn't like some drug trial btw lol.  This is like, turning on a machine... im sure they've been doing this every day this month to be prepared.",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		},
	}

	for name, test := range tests {
		testActionProfile := &ActionProfile{}
		Recommendation(testActionProfile, test.input)
		if testActionProfile.action != test.expectedAction {
			t.Errorf("test \"%s\" failed: expected action %d, got %d", name, test.expectedAction, testActionProfile.action)
		} else if testActionProfile.multiplier != test.expectedMultiplier {
			t.Errorf("test \"%s\" failed: expected multiplier %d, got %d", name, test.expectedMultiplier, testActionProfile.multiplier)
		}
	}
}
