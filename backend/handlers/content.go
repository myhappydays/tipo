package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tipo-backend/models"
	"tipo-backend/utils"
)

// Naver API 응답 구조체들
type NaverBookResponse struct {
	Items []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Author      string `json:"author"`
	} `json:"items"`
}



type NaverNewsResponse struct {
	Items []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"items"`
}

// 문장 분리 유틸리티 함수
func splitToSentences(text string) []string {
	text = strings.ReplaceAll(text, ". ", ".")
	text = strings.ReplaceAll(text, "? ", "?")
	text = strings.ReplaceAll(text, "! ", "!")
	sentences := strings.Split(text, ".")
	var result []string
	for _, s := range sentences {
		if len(strings.TrimSpace(s)) > 0 {
			result = append(result, strings.TrimSpace(s)+".")
		}
	}
	return result
}

func ContentHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	query := r.URL.Query().Get("query")

	if category == "" || query == "" {
		http.Error(w, "'category' and 'query' parameters are required", http.StatusBadRequest)
		return
	}

	// Naver API 설정
	clientID := os.Getenv("NAVER_CLIENT_ID")
	clientSecret := os.Getenv("NAVER_CLIENT_SECRET")

	

	apiURL := ""

	switch category {
	case "book":
		apiURL = "https://openapi.naver.com/v1/search/book.json"
	
	case "news":
		apiURL = "https://openapi.naver.com/v1/search/news.json"
	default:
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	fullAPIURL := apiURL + "?query=" + url.QueryEscape(query)
	

	// Naver API 요청
	req, err := http.NewRequest("GET", fullAPIURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Naver-Client-Id", clientID)
	req.Header.Set("X-Naver-Client-Secret", clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	// 최종 응답 데이터
	var contentResponse models.ContentResponse

	// 카테고리별 데이터 파싱 및 가공
	switch category {
	case "book":
		var data NaverBookResponse
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(data.Items) > 0 {
			item := data.Items[0]
			cleanedDesc := utils.CleanHtmlTags(item.Description)
			contentResponse = models.ContentResponse{
				Title:     utils.CleanHtmlTags(item.Title),
				Source:    utils.CleanHtmlTags(item.Author),
				Category:  "도서",
				Image:     item.Image,
				Sentences: splitToSentences(cleanedDesc),
			}
		}
	
	case "news":
		var data NaverNewsResponse
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(data.Items) > 0 {
			item := data.Items[0]
			cleanedDesc := utils.CleanHtmlTags(item.Description)
			contentResponse = models.ContentResponse{
				Title:     utils.CleanHtmlTags(item.Title),
				Source:    "뉴스",
				Category:  "뉴스",
				Image:     "", // 뉴스 API는 이미지를 제공하지 않음
				Sentences: splitToSentences(cleanedDesc),
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contentResponse)
}
