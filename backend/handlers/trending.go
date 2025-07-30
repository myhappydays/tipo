package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"math/rand"
	"time"
	"strings"
	"log"

	"github.com/mmcdole/gofeed"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TrendingKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	randomStr := r.URL.Query().Get("random")

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 5 // 기본값 5개
	}

	random := false
	if randomStr == "true" {
		random = true
	}

	// Google Trends RSS 피드 URL (한국 기준)
	feedURL := "https://trends.google.com/trending/rss?geo=KR"

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		log.Printf("Error parsing Google Trends RSS feed: %v", err) // 에러 로깅 추가
		http.Error(w, "Failed to parse Google Trends RSS feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var keywords []string
	for _, item := range feed.Items {
		// Google Trends RSS의 제목은 보통 '키워드 - 트렌드' 형식
		// 여기서는 '키워드' 부분만 추출
		parts := strings.Split(item.Title, " - ")
		if len(parts) > 0 {
			keywords = append(keywords, strings.TrimSpace(parts[0]))
		}
	}

	if random && len(keywords) > 0 {
		selectedKeyword := keywords[rand.Intn(len(keywords))]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"keyword": selectedKeyword})
		return
	}

	// 요청된 개수만큼만 반환
	if len(keywords) > count {
		keywords = keywords[:count]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keywords)
}
