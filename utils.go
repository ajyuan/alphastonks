package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func timer() func() {
	start := time.Now()
	return func() {
		log.Infof("Iteration completed in %v", time.Since(start))
	}
}

func stringIn(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func roundPriceDown(value float32) float64 {
	return math.Floor(float64(value*100)) / 100
}

func substrPrefSuf(page, prefix, suffix string) (string, error) {
	si, n, ei := strings.Index(page, prefix)+len(prefix), len(suffix), -1
	if si == -1 {
		return "", fmt.Errorf("extractPosts failed to find data prefix \"%s\"", postTextPrefix)
	}
	for i := si + 1; i < len(page)-n; i++ {
		if page[i:i+n] == suffix {
			ei = i
			break
		}
	}
	if ei == -1 {
		return "", fmt.Errorf("extractPosts failed to find data suffix \"%s\"", suffix)
	}
	return page[si:ei], nil
}

func minZero(a int) int {
	if a < 0 {
		return 0
	}
	return a
}
