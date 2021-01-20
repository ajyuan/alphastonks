package main

import (
	"testing"
)

func TestDiscoveredWithinBounds(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedOutput bool
	}{
		"1 second ago": {
			input:          "1 second ago",
			expectedOutput: true,
		}, "2 seconds ago": {
			input:          "2 seconds ago",
			expectedOutput: true,
		}, "single digit too late": {
			input:          "3 seconds ago",
			expectedOutput: false,
		}, "double digits too late": {
			input:          "10 seconds ago",
			expectedOutput: false,
		}, "undefined string": {
			input:          "Just now",
			expectedOutput: true,
		},
	}
	for name, test := range tests {
		out := discoveredWithinBounds(test.input)
		if out != test.expectedOutput {
			t.Fatalf("test \"%s\" failed expected %t, got %t", name, test.expectedOutput, out)
		}
	}
}

func TestRecommendation(t *testing.T) {
	tests := map[string]struct {
		input              string
		expectedAction     uint
		expectedMultiplier float32
	}{
		"CRNT Buy": {
			input:              "New penny stock that I am currently interested in (& invested in) is CRNT (Ceragon Networks).  Ark Invest recently purchased over 150k more shares of this penny stock from 1/11 - 1/14 (i took a picture for proof and will compare Cathie Wood's portfolios in the video).  This is a 5G play which is 1 of the biggest growth sectors but what made me really excited about this stock is that Huawei is getting banned in many countries + companies don't want to work with them because of their reputation of stealing.  CRNT said they feel like they're in a good position to take market share from Huawei (I will show all of their statements in the video that I plan on releasing sometime next week). Needham upgraded CRNT and gave it a buy rating after their presentation at its growth conference yesterday giving me more confidence in this company. If you want to help me with my DD in this company, please feel free to add notes in this comment section or my subreddit (trakstocks).  I will try to go over this briefly on Sunday and then do a full video sometime next week (maybe Monday but not positive)",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		},
		"CPSH Buy in Update": {
			input:              "A lot of people are DMing me asking for an \"exclusive option within the exclusive option\" so you can see the stock picks faster.  I will think of an option but in the meantime I have a penny stock that I am going to announce in here in a couple hours ( around 9:30am - 9:45am ET).  After this announcement, then we can figure out the other pricing options.  Please always do your research before investing because these cheap stocks are cheap for a reason (they could be bad stocks).  Remember:  CPSH has the potential to do very well in a lot of fast growing industries so pls dont compare CPSH stock performance with this new stock pick because they're different types of companies.  With that being said, I am not just talking about penny stocks just to talk about them... there always has to be at least 3 big reasons for me to want to invest in these risky stocks and there is 1 reason that does not make this your ordinary penny stock.  Stay tuned.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "CPSH Buy": {
			input:              "Hey, I just invested in the penny stock CPSH.  I will be doing a video tomorrow on this stock.  I believe this company can do well because their products are revolving around hybrid / electric vehicles, aerospace, wind turbines / clean energy.  The reason why I think this won't be a penny stock for long is because of the role they play in the clean energy sector.  Also the space industry is growing and experts believe the space industry will triple to 1.4 trillion within a decade.  Yesterday on Reddit's home page, they were showing Nasa's new Mars rover (Perseverance Rover) and how it's about to land on mars on Feb. 18th.  Nasa is using CPSH's product in that rover.  Last quarter they were net profitable for their operations and with them playing an important role in the green energy / ev sector, I think this might be the beginning of something great for this company.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "MP Buy": {
			input:              "hey, planning on doing a video on MP tomorrow.  Might be Wednesday latest but shooting for Tuesday.  I think they're going to do well in upcoming earnings and I feel like I can come pretty close to proving it. I could always be wrong so when you watch the video, if you think I am \"seeing what I want to see\" just ignore it but with Biden + what they're already doing, i think they're going to do well.  PLEASE tell me what you don't like with MP and i will try to address it but watch the video I attached where I addressed common concerns with the company.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "STPKK Buy": {
			input:              "Hey, STPK I am hearing great things about.  I am working on a video but don't know when it will be done.  I am hearing a lot about how this might do really well under Biden.  If you have any concerns about this company please leave it down below so I can try to address them in the video.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
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
			expectedMultiplier: lowBuyMult,
		}, "FRSX High Buy": {
			input:              "Just purchased the penny stock FRSX.  I covered them before but after apple \"icar\" news I have invested heavily.  Will try to post video in a few hours.  Not saying they will work with apple, just noticing growing interest in this sector and this stock seems to be the last \"autonomous driving related\" related stock left and it's under $3 a share.  They are getting more and more opportunities, expecting a grant of approximately one million USD from the European Commission, won Edison award (something Elon Musk won a while back), a COVID-19 Symptom Detection play, and will cover more in video.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "IMMR Low/High Buy": {
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
		}, "NNOX Buy/NoOp": {
			input:              "Just bought more NNOX (like.. a lot more).  If it fails at live demonstration (which is tomorrow), stock will drop VERY hard so big risk.  Big catalyst tomorrow (maybe positive maybe negative).  This isn't like some drug trial btw lol.  This is like, turning on a machine... im sure they've been doing this every day this month to be prepared.",
			expectedAction:     actionBuy,
			expectedMultiplier: lowBuyMult,
		}, "AQB Sell, MP Buy": {
			input:          "Hey, I usually tell you about when I buy more than I sell so I wanted to let you know I have currently sold all of my AQB position (since it has gone up 101% this month).  There is a big catalyst right around the corner according to the CEO about them distributing their salmon for the first time ever and they said they will announce this by the end of the year (so don't sell on account of me).  I have just had an incredible run with this stock because I got in at $2.83 and I wanted to invest heavily into MP (MP Materials).  Still working on the video for MP but it's A LOT of info.  If AQB significantly drops, I will get back in but yea, just wanted to bring some balance to this message board with my buys & sells.",
			expectedAction: actionBuy,
		}, "MP Buy 2": {
			input:              "check MP.  Working on long video that might be released tomorrow on it.  Its currently negative 5%.  It has had huge run up lately but at least put it on your radar (maybe will fall more but I am buying). Most rare earth separation and metal processing plants have sold out of production until early 2021, tightening spot supplies and pushing up domestic prices further.  China controls 80% of the rare earth elements. You need rare earth for all of these electric (ev) related products (and hydrogen).  MP is the 2nd largest supplier and only supplier in USA (& entire western hemisphere). USA is trying to be less dependent on China ( i will explain more in detail when i do video on why MP looks like this will directly benefit them). China just did export law and people speculate they will limit rare earth elements to counteract to our tarrifs (directly benefiting MP).  Still doing research but as of right now, seems like this has a lot of room to run.  This is a LONG hold imo, not short term.  Also LGVW is down again, I think this will recover nicely.  Ark purchased more shares again in LGVW (but is not announcing it in her newsletter, i took pics of her portfolio so i saw it) and have a feeling they will buy again today. I highly do not recommend selling LGVW for a loss currently.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "KNDI Sell": {
			input:          "If you're still bullish on KNDI (Kandi) after recent short seller report can you tell me why?  Currently looking into it to see if it's the perfect time to buy.  Please only respond if you're bullish on KNDI so we can filter out best reasons to still believe in this company.",
			expectedAction: actionNoOp,
		}, "FEAC Buy 2": {
			input:              "It might be better to wait for a red day on this BUT I just invested in FEAC.  Currently can't see why this isn't a $20 stock (or \"unit\") minimum.  Whenever I am this bullish on a company, I am scared I am missing something but I think the reason the stock is not higher $ is because it's under the radar. FEAC is merging with Skillz. Skillz is an online mobile multiplayer competition platform that is integrated into a lot of iOS and Android games. Players use it to compete in competitions against other players across the world (and it can be for money).  Skillz is one of the fastest growing companies (was ranked #1 fastest private company in USA in 2017 by Inc). IMO, Skillz has even greater potential than sports betting because eSports can be played 24/7.  Sports are only specific days/times.  Will release video very soon.",
			expectedAction:     actionBuy,
			expectedMultiplier: lowBuyMult,
		}, "LGVW Buy": {
			input:              "I just bought more LGVW.  IMO it seems like a no brainer that this -10 drop will recover soon.  I would not be surprised to see Cathie Wood at Ark Invest buy more today.  Ark invest said on Friday that LGVW (butterfly iQ) is the only handheld ultrasound device to abandon legacy piezoelectric technology and to offer higher quality, lower costs, and form-factor advantages over incumbents.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "Metromile NoOp": {
			input:          "After closer inspection I personally dont see the appeal in Metromile.  According to the National Association of Insurance Commissioners, Metromile has a complaint ratio much higher than the national average.  Customers most often complain about its customer service practices and unexpected rate increases.  IMO people tend to leave reviews when they're upset rather than satisfied with something so keep that in mind.  This could definitely change and this company could do very well but this \"pay per mile\" business model they have seems easy to recreate & also in this video I will show you several companies that are already doing this.",
			expectedAction: actionNoOp,
		}, "INAQ Buy": {
			input:              "Spac INAQ ( Metromile ) Mark Cuban & Chamath Palihapitiya backed company.  I am currently not invested but it could do well.  Learning more about it now.  They're a car insurance startup that offers pay-per-mile insurance",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "PSTH Buy": {
			input:              "Ok here's a stock that could do very well.  There is rumors that Bill Ackman's spac PSTH might merge with Stripe.  I heard Bill Ackman say in an interview in Sep. that stripe wasn’t ready to go public and they’re looking for companies that are ready to go public.  So basically he is saying they aren't going to choose Stripe.  But there was a top exec at Stripe who recently started following bill on twitter so people are speculating.  If PSTH merges with Stripe, it will do well imo but does anyone have any new reasons to believe he will choose them (other than stripe being a good company... I'm talking clues they will merge)?  Any NEW news on this (news after the Ackman interview in Sep.) would be much appreciated.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "LGVW Buy 2": {
			input:              "check LGVW asap!  They just announced merger with Butterfly Network.  Butterfly iQ is the only ultrasound transducer that can perform \"whole-body imaging\" with a single handheld probe using semiconductor technology. Connected to a mobile phone or tablet, it is powered by Butterfly's proprietary Ultrasound-on-Chip™ technology and harnesses the advantages of AI to deliver advanced imaging that we believe is easy-to-use, improves patient outcomes and lowers cost of care.  will try to do a video today",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "BTC NoOp": {
			input:          "Important! I am not referring anyone to BTC or telling you to use any service in my comments.  Scammers are making accounts with my pic and same name.  Pls dont fall for it.",
			expectedAction: actionNoOp,
		}, "GDRX NoOp": {
			input:          "Ceo confirmed this... GoodRX has lower prices with GoodRX Gold than Amazon Prime over 90% of the time.  Plus only 5% of the industry is mail order (yes could change) but again you can do mail orders through GDRX with your pharmacy.  The Earning Reports will be the catalysts.",
			expectedAction: actionNoOp,
		}, "GDRX Buy": {
			input:              "I bought more GDRX!  down -18% from amazon pharmacy announcement.  Wife thinks its over reaction (she's a pharmacist).  Video coming soon of our conversation on why she thinks GoodRx is still good.",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		}, "MRNA Buy": {
			input:              "MRNA 94% effective.  I got in before market hours.  I think this will do well because refrigeration process is also easier than pfizer.  Also this is why I like WeBull so you can buy before/after hours.  Here is my link if you want to download -> Two free stocks valued up to $1400 using this link: https://bit.ly/2Eieppx",
			expectedAction:     actionBuy,
			expectedMultiplier: highBuyMult,
		},
	}

	for name, test := range tests {
		testPost := &YTPostDetails{
			postText: test.input,
			postTime: "1 second ago",
		}
		actionProfile, err := Recommendation(testPost)
		if err != nil {
			t.Fatalf("test \"%s\" failed: expected no error, got %v", name, err)
		} else if actionProfile.action != test.expectedAction {
			t.Errorf("test \"%s\" failed: expected action %d, got profile %v", name, test.expectedAction, actionProfile)
		}
		/*
			else if actionProfile.multiplier != test.expectedMultiplier {
				t.Errorf("test \"%s\" failed: expected multiplier %f, got profile %v", name, test.expectedMultiplier, actionProfile)
			}
		*/
	}
}
