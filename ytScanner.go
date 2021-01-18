package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

// YTPost extracts information about the latest post from YT
func YTPost(cl *http.Client) (*YTPostDetails, error) {
	page, err := communityPage(cl, ytTarget)
	if err != nil {
		return nil, err
	}
	postText, err := substrPrefSuf(page, postTextPrefix, postTextSuffix)
	if err != nil {
		return nil, err
	}
	postText = cleanHTMLString(postText)
	postTime, err := substrPrefSuf(page, postTimePrefix, postTimeSuffix)
	if err != nil {
		return nil, err
	}
	return &YTPostDetails{
		postText: postText,
		postTime: postTime,
	}, nil
}
