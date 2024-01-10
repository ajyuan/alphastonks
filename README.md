## While this stock generated marginal returns during it's inception, it is no longer an effective strategy. The effect of stock market influencers has been diminished since the pandemic, and similar bots seem to have been started by others.

Alphastonks is an experimental stock trading bot.

# Motivation

During the COVID pandemic, some stock influencers popped up on YouTube. They would call our small-cap stocks, and viewers would immediately buy them up, causing the share price to dramatically rise. This bot aims to monitor the YouTube social feed for new stocks being announced, and purchase some shares before others buy them.

# Strategy

### 1. Scan the social feed for YouTuber DeadNSYEd

If a stock ticker is detected, perform sentiment analysis on the message. Messages may promote or shun a stock, so we want to make sure we only buy stocks that are being promoted.

### 2. Buy the stock on Alpaca

Alpaca is a brokerage that has APIs.

### 3. Sell the stock after X minutes of holding

The stock price often collapses after a critial number of users see and buy the stock. Sell it before then.
