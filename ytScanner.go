package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	// ErrPostTxtNotFound indicates no post text was found in a given post
	ErrPostTxtNotFound = fmt.Errorf("No post text found in page")
	// ErrPostDateNotFound indicates no post text was found in a given post
	ErrPostDateNotFound = fmt.Errorf("No post date found in page")
)

type ytTextJSON struct {
	Text string
}

// YTPostDetails contains details about a post required to make an action on it
type YTPostDetails struct {
	postTime string
	postText string
}

func cleanHTMLString(postText string) string {
	postText = strings.ReplaceAll(postText, "\\r\\n", "")
	return strings.ReplaceAll(postText, "\\u0026", "and")
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

func postText(page string) (string, error) {
	reMatches := postTextRe.FindStringSubmatch(page)
	if len(reMatches) < 2 {
		return "", ErrPostTxtNotFound
	}
	postTextMatch := reMatches[1]
	var unmarshalled []ytTextJSON
	err := json.Unmarshal([]byte(postTextMatch), &unmarshalled)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal %s: %v", postTextMatch, err)
	}
	postText := ""
	for _, postSection := range unmarshalled {
		postText += postSection.Text
	}
	return postText, nil
}

// YTPost extracts information about the latest post from YT
func YTPost(cl *http.Client) (*YTPostDetails, error) {
	page, err := communityPage(cl, ytTarget)
	if err != nil {
		return nil, err
	}
	postText, err := postText(page)
	if err != nil {
		return nil, err
	}
	postTime, err := substrPrefSuf(page, postTimePrefix, postTimeSuffix)
	if err != nil {
		return nil, ErrPostDateNotFound
	}
	return &YTPostDetails{
		postText: postText,
		postTime: postTime,
	}, nil
}
