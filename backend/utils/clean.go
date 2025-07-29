package utils

import "regexp"

// CleanHtmlTags는 문자열에서 HTML 태그를 제거합니다.
func CleanHtmlTags(s string) string {
	// 정규식을 사용하여 모든 HTML 태그를 찾습니다. 예: <b>, </b>, <p> 등
	re := regexp.MustCompile(`<[^>]*>`)
	// 찾은 모든 태그를 빈 문자열로 대체합니다.
	return re.ReplaceAllString(s, "")
}
