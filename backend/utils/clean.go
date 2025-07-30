package utils

import (
	"html"
	"regexp"
	"strings"
)

// CleanHtmlTags는 문자열에서 HTML 태그, 불필요한 문구, 기자 정보 등을 제거하여 텍스트를 정제합니다.
func CleanHtmlTags(s string) string {
	// 1. HTML 엔티티 디코딩 (e.g., &quot; -> ")
	s = html.UnescapeString(s)

	// 2. HTML 태그 제거
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")

	// 3. 불필요한 번역 관련 문구 제거
	s = strings.ReplaceAll(s, "It is assumed that there may be errors in the English translation.", "")
	s = strings.ReplaceAll(s, "This is a machine translation and may contain errors.", "")

	// 4. 기자 정보 및 불필요한 문구 제거 (정규식 사용)
	// 예: [이데일리 최선미 기자] blah blah, (서울=뉴스1) blah blah, 아주경제=송윤서 기자 sys0303@ajunews.com, /연합뉴스
	re = regexp.MustCompile(`\[[^\]]+\]|【[^】]+】|\([^\)]+\)|[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}|[가-힣]+ 기자 ?=|/뉴시스|/연합뉴스`)
	s = re.ReplaceAllString(s, "")

	// 5. 연속된 공백을 하나로 줄이고, 앞뒤 공백 제거
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	return s
}