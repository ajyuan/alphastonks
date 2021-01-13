package main

import (
	"testing"
)

func TestRecommendation(t *testing.T) {
	tests := map[string]struct {
		text               string
		expectedAction     uint
		expectedMultiplier int32
	}{
		"CPSH Buy": {
			text:               "Hey, I just invested in the penny stock CPSH.  I will be doing a video tomorrow on this stock.  I believe this company can do well because their products are revolving around hybrid / electric vehicles, aerospace, wind turbines / clean energy.  The reason why I think this won't be a penny stock for long is because of the role they play in the clean energy sector.  Also the space industry is growing and experts believe the space industry will triple to 1.4 trillion within a decade.  Yesterday on Reddit's home page, they were showing Nasa's new Mars rover (Perseverance Rover) and how it's about to land on mars on Feb. 18th.  Nasa is using CPSH's product in that rover.  Last quarter they were net profitable for their operations and with them playing an important role in the green energy / ev sector, I think this might be the beginning of something great for this company.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "MP Buy/NoOp": {
			text:           "hey, planning on doing a video on MP tomorrow.  Might be Wednesday latest but shooting for Tuesday.  I think they're going to do well in upcoming earnings and I feel like I can come pretty close to proving it. I could always be wrong so when you watch the video, if you think I am \"seeing what I want to see\" just ignore it but with Biden + what they're already doing, i think they're going to do well.  PLEASE tell me what you don't like with MP and i will try to address it but watch the video I attached where I addressed common concerns with the company.",
			expectedAction: actionNoOp,
		}, "STPKK Buy": {
			text:               "Hey, STPK I am hearing great things about.  I am working on a video but don't know when it will be done.  I am hearing a lot about how this might do really well under Biden.  If you have any concerns about this company please leave it down below so I can try to address them in the video.",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		}, "BNGO NoOp": {
			text:           "PLEASE do your research because I dont want this BNGO to hit $20 because I am missing something!  With that being said, I just got off the phone and am less bullish with BNGO.  VERY hard subject to tackle so please do your research carefully, Jan 11th - 15th could still be a catalyst.  From my understanding there is 2 ways to look at this stock, from a \"clinical market\" side and as a \"research market\" side.  Clinical market has more of a TAM then Research Market (it's a night and day difference apparently).  I am under the impression BNGO currently can only address a very small part of \"clinical market\" and more so a \"research market\" company.",
			expectedAction: actionNoOp,
		}, "BABA NoOp": {
			text:           "Hey, I did a video on BABA a few days ago but after reading about what the Chinese government does to these \"outspoken\" billionaires, I would think waiting before you get in, might be best (in case there are any \"forced confessions\" from Jack Ma like what has happened in the past with other companies).  Will do a video on this soon.  This stock has one of the most \"buy\" recommendations from analysts I have seen so that's why this interests me so much.",
			expectedAction: actionNoOp,
		}, "BNGO High Buy": {
			text:               "Ok after closer look I have invested in BNGO.  Will try to post a video this week (maybe wed or thursday).  Main reason why is because of the study published 6 days ago where they show ways BNGO ($1 a share) is better than PACB ($24 a share) and the day after the study was published, Ark Invest was trying to learn more about BNGO from their CEO.  If the BNGO share price performs even 1/4th as good as PACB I will be VERY happy lol.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "BNGO Med Buy": {
			text:               "I talked about PACB on youtube when it was $6 (now $24) because of Ark putting them on my radar.  While BNGO is still at $1, pls read this study to determine if you think BNGO is better than PACB.  Trying my best to not influence your decision lol.  Maybe I am reading this wrong, but it seems like BNGO is better in a couple ways. https://apnews.com/press-release/glob...",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		}, "FRSX Buy": {
			text:               "Just purchased the penny stock FRSX.  I covered them before but after apple \"icar\" news I have invested heavily.  Will try to post video in a few hours.  Not saying they will work with apple, just noticing growing interest in this sector and this stock seems to be the last \"autonomous driving related\" related stock left and it's under $3 a share.  They are getting more and more opportunities, expecting a grant of approximately one million USD from the European Commission, won Edison award (something Elon Musk won a while back), a COVID-19 Symptom Detection play, and will cover more in video.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "IMMR Buy": {
			text:               "Hey, currently I am not invested in this but working on a video about IMMR.  They could do well soon.  Gathering up a lot of interesting information about this company on why it might be a good stock to invest in but will not get the video done until maybe next week.  I think it's def worth it for you to look closely at this company.  If you look at it closely and do NOT like it, please leave reasons why so we can see : )",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "SAMA Buy": {
			text:               "Hey, take a look at SAMA.  They will be merging with Clever Leaves on the 17th.  Clever Leaves is a \"Global Company and leading vertically integrated producer of  medical cannabis and hemp extracts\".  They recently signed a deal with Canopy Growth and export their product from Columbia.  They are GMP certified which apparently gives them a big advantage to sell medical grade cannabis and sell at a higher cost.  Check it out.  It's down -5% currently.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "BFT Buy": {
			text:               "Very briefly looked into BFT (Paysafe) and from what I can see, it looks like a great buy.  I just now bought shares, quite a bit because I think this will have a lot of momentum.  I can see this hitting at least $15 sometime this week or next week.  Doing more research now but this is looking really good (from what I can currently see).  I would not worry about the recent run up because I think it's small when looking at the bigger picture.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "FEAC Buy": {
			text:               "Sorry for the repetitive videos about FEAC but after looking even closer at the company, I am even more bullish on this one.  99% positive this stock will double (not sure when but confident it will eventually,.  It has ran up recently so could be a pull back).  I will do a video today about more reasons why because of new things I discovered.  I have been trying to identify every negative thing with the company in the past 2 videos and the biggest negative catalyst I can imagine is some regulation law comes into place banning this \"skills based\" gambling.  Also I believe App store is going from 30% fee to 15% fee on in app purchases (which is better than before but I dont like how apple has so much control).  Not 100% sure this effects skillz (i am pretending that it does) but if apple raises the price again, could be another negative catalyst potentially.  Other than that, I am struggling to find more problems with the company.  Bullish and this is a long term play.  So far REALLY liking this ceo too.  Seems focused and what he says makes sense to me.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "NNOX Buy": {
			text:               "Just bought more NNOX (like.. a lot more).  If it fails at live demonstration (which is tomorrow), stock will drop VERY hard so big risk.  Big catalyst tomorrow (maybe positive maybe negative).  This isn't like some drug trial btw lol.  This is like, turning on a machine... im sure they've been doing this every day this month to be prepared.",
			expectedAction:     actionBuy,
			expectedMultiplier: medBuyMult,
		},
	}

	for name, test := range tests {
		testActionProfile := &ActionProfile{}
		Recommendation(testActionProfile, test.text)
		if testActionProfile.action != test.expectedAction {
			t.Errorf("test \"%s\" failed: expected action %d, got %d", name, test.expectedAction, testActionProfile.action)
		} else if testActionProfile.multiplier != test.expectedMultiplier {
			t.Errorf("test \"%s\" failed: expected multiplier %d, got %d", name, test.expectedMultiplier, testActionProfile.multiplier)
		}
	}
}
