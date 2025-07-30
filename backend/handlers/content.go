package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tipo-backend/models"
	"tipo-backend/utils"
)

// Naver API 응답 구조체들
type NaverNewsItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PubDate     string `json:"pubDate"`
	Originallink string `json:"originallink"` // Add originallink field
}

type NaverNewsResponse struct {
	Items []NaverNewsItem `json:"items"`
}

// extractDomain extracts the domain from a URL
func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// 문장 분리 유틸리티 함수
func splitToSentences(text string) []string {
	// 1. 다중 공백을 단일 공백으로 변경
	text = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(text), " ")

	// 2. 문장 분리 정규식
	re := regexp.MustCompile(`([^.!?]+(?:[.!?](?:\s|$)|\z))`)
	sentences := re.FindAllString(text, -1)

	if len(sentences) == 0 {
		return []string{}
	}

	var result []string
	var buffer strings.Builder

	for _, s := range sentences {
		trimmedSentence := strings.TrimSpace(s)
		if trimmedSentence == "" {
			continue
		}

		// 현재 문장이 짧으면 버퍼에 추가
		if len([]rune(trimmedSentence)) < 10 {
			if buffer.Len() > 0 {
				buffer.WriteString(" ")
			}
			buffer.WriteString(trimmedSentence)
		} else { // 현재 문장이 길면
			// 버퍼에 있던 짧은 문장들을 현재 문장 앞에 붙여서 결과에 추가
			if buffer.Len() > 0 {
				buffer.WriteString(" ")
			}
			buffer.WriteString(trimmedSentence)
			result = append(result, buffer.String())
			buffer.Reset() // 버퍼 비우기
		}
	}

	// 루프가 끝난 후 버퍼에 남은 내용이 있다면 (마지막 문장들이 짧은 경우)
	// 그냥 결과에 추가한다.
	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}

	// 분리된 문장이 하나도 없는 경우 (모든 문장이 짧아서 버퍼에만 있다가 마지막에 추가된 경우 등)
	if len(result) == 0 && len([]rune(text)) > 0 {
		return []string{text} // 원본 텍스트 반환
	}

	return result
}

func ContentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	startStr := r.URL.Query().Get("start")
	displayStr := r.URL.Query().Get("display")

	if query == "" {
		http.Error(w, "'query' parameter is required", http.StatusBadRequest)
		return
	}

	start, err := strconv.Atoi(startStr)
	if err != nil || start < 1 {
		start = 1 // 유효하지 않거나 없는 경우 1로 설정
	}

	display, err := strconv.Atoi(displayStr)
	if err != nil || display < 1 || display > 100 {
		display = 10 // 유효하지 않거나 없는 경우 10으로 설정 (네이버 API 최대 100)
	}

	// Naver API 설정
	clientID := os.Getenv("NAVER_CLIENT_ID")
	clientSecret := os.Getenv("NAVER_CLIENT_SECRET")

	apiURL := "https://openapi.naver.com/v1/search/news.json"

	// Naver API는 최대 1000개 결과까지만 제공. 1000개 넘어가면 다시 처음부터 시작
	if start > 1000 {
		start = 1
	}

	// 다중 키워드 처리
	keywords := strings.Split(query, ",")
	numKeywords := len(keywords)
	if numKeywords == 0 {
		http.Error(w, "'query' parameter cannot be empty", http.StatusBadRequest)
		return
	}

	// 각 키워드별로 가져올 기사 수 계산
	baseArticlesPerKeyword := display / numKeywords
	remainder := display % numKeywords

	var articlesToFetch []int
	for i := 0; i < numKeywords; i++ {
		count := baseArticlesPerKeyword
		if i < remainder {
			count++
		}
		articlesToFetch = append(articlesToFetch, count)
	}

	var allItems [][]NaverNewsItem // 각 키워드별로 가져온 아이템들을 저장
	maxItems := 0 // 가장 많은 아이템을 가져온 키워드의 아이템 수

	for idx, kw := range keywords {
		currentArticlesToFetch := articlesToFetch[idx]
		if currentArticlesToFetch == 0 {
			allItems = append(allItems, []NaverNewsItem{})
			continue
		}
		fullAPIURL := apiURL + "?query=" + url.QueryEscape(strings.TrimSpace(kw)) + "&start=" + strconv.Itoa(start) + "&display=" + strconv.Itoa(currentArticlesToFetch)

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

		var data NaverNewsResponse
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allItems = append(allItems, data.Items)
		if len(data.Items) > maxItems {
			maxItems = len(data.Items)
		}
	}

	// 최종 응답 데이터 (슬라이스)
	var contentResponses []models.ContentResponse

	// 아이템들을 번갈아가면서 섞기
	for i := 0; i < maxItems; i++ {
		for j := 0; j < numKeywords; j++ {
			if i < len(allItems[j]) {
				item := allItems[j][i]
				cleanedDesc := utils.CleanHtmlTags(item.Description)
				contentResponses = append(contentResponses, models.ContentResponse{
					Title:     utils.CleanHtmlTags(item.Title),
					Source:    extractDomain(item.Originallink), // extract domain from originallink
					Keyword:   strings.TrimSpace(keywords[j]), // 해당 기사를 가져온 개별 키워드 사용
					PubDate:   item.PubDate,
					Sentences: splitToSentences(cleanedDesc),
				})
				if len(contentResponses) >= display {
					break // display 개수만큼만 가져오기
				}
			}
		}
		if len(contentResponses) >= display {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contentResponses)
}


