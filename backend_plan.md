# Tipo 프로젝트 백엔드 개발 계획 (Golang)

- **문서 버전:** 1.0
- **작성일:** 2025-07-29
- **기술 스택:** Golang (Go)

---

## 1. MVP 목표

- Naver API 키 등 민감한 정보를 서버에서 안전하게 관리한다.
- 클라이언트(프론트엔드)의 요청에 따라 Naver 뉴스 API를 대신 호출하는 프록시(Proxy) 서버 역할을 수행한다.
- Google Trends RSS 피드를 파싱하여 트렌딩 키워드를 제공한다.
- Naver 뉴스 API로부터 받은 데이터를 타자 연습에 적합한 형태로 가공하여 일관된 JSON 형식으로 제공한다.
    - `Source` (언론사), `Keyword` (검색 키워드), `PubDate` (발행일) 정보 포함.

---

## 2. 기술 명세 (Technical Specifications)

- **언어:** Go (Golang)
- **웹 프레임워크:** 표준 라이브러리 `net/http` 사용 (외부 프레임워크 없이 구현)
    - *이유: MVP 기능은 간단한 API 엔드포인트 하나로 충분하므로, 가볍고 빠른 표준 라이브러리가 가장 효율적입니다.*
- **라우팅:** `http.HandleFunc`를 사용한 기본 라우터
- **환경 변수 관리:** `os.Getenv` (표준 라이브러리) 또는 `godotenv` 라이브러리 활용
- **RSS 파싱:** `github.com/mmcdole/gofeed` 라이브러리 활용

---

## 3. 개발 단계별 계획

### 1단계: 프로젝트 구조 및 환경 설정

1.  **프로젝트 초기화:**
    ```bash
    go mod init tipo-backend
    ```

2.  **환경 변수 설정:**
    - 프로젝트 루트에 `.env` 파일을 생성하고 Naver API 키를 저장합니다.
      ```
      NAVER_CLIENT_ID="YOUR_NAVER_CLIENT_ID"
      NAVER_CLIENT_SECRET="YOUR_NAVER_CLIENT_SECRET"
      ```
    - `.gitignore` 파일에 `.env`를 추가합니다.

3.  **기본 파일 구조:**
    ```
    /tipo-backend
    ├── go.mod
    ├── go.sum
    ├── main.go         # 애플리케이션 진입점
    ├── .env            # 환경 변수
    └── /handlers       # HTTP 요청 핸들러
        ├── content.go
        └── trending.go
    └── /models         # 데이터 구조체 정의
        └── content.go
    └── /utils          # 유틸리티 함수 (예: HTML 태그 제거)
        └── clean.go
    ```

### 2단계: 핵심 로직 구현

1.  **데이터 모델 정의 (`/models/content.go`):**
    - API 응답으로 내려줄 JSON 구조체를 정의합니다.

    ```go
    package models

    type ContentResponse struct {
        Title     string   `json:"title"`
        Source    string   `json:"source"`    // 언론사 도메인
        Keyword   string   `json:"keyword"`   // 검색 키워드
        PubDate   string   `json:"pubDate"`   // 발행일
        Sentences []string `json:"sentences"`
    }
    ```

2.  **API 핸들러 구현 (`/handlers/content.go` 및 `/handlers/trending.go`):**
    - `/api/content` 엔드포인트의 로직을 작성합니다.
        - 클라이언트로부터 `query` (쉼표로 구분된 다중 키워드), `start`, `display` 파라미터를 받습니다.
        - `os.Getenv`를 사용해 환경 변수에서 Naver API 키를 읽어옵니다.
        - `net/http` 클라이언트를 사용해 Naver 뉴스 API를 호출하고 헤더를 설정합니다.
        - 다중 키워드에 대해 `display` 개수를 균등하게 나누어 기사를 가져오고, 번갈아가면서 섞어(`interleaving`) 최종 `display` 개수만큼의 `models.ContentResponse` 객체 리스트를 JSON으로 인코딩하여 응답합니다.
        - `utils` 패키지의 함수를 호출하여 텍스트를 정제합니다.
        - `Source`는 기사의 원본 링크에서 도메인을 추출하여 사용하고, `Keyword`는 해당 기사를 가져온 개별 키워드를 사용하며, `PubDate`를 포함합니다.
    - `/api/trending_keywords` 엔드포인트의 로직을 작성합니다.
        - 클라이언트로부터 `count` (가져올 키워드 개수)와 `random` (랜덤 키워드 반환 여부) 파라미터를 받습니다.
        - Google Trends RSS 피드(`https://trends.google.com/trending/rss?geo=KR`)를 파싱하여 트렌딩 키워드를 가져옵니다.
        - 요청된 개수만큼의 키워드 리스트를 JSON으로 인코딩하여 응답하거나, `random` 파라미터가 `true`일 경우 단일 랜덤 키워드를 반환합니다.

3.  **HTML 태그 제거 유틸리티 (`/utils/clean.go`):**
    - 정규식(`regexp` 패키지)을 사용하여 문자열에서 HTML 태그를 제거하는 간단한 함수를 작성합니다.

4.  **메인 함수 작성 (`main.go`):**
    - `http.HandleFunc`를 사용하여 `/api/content` 경로와 `handlers.ContentHandler` 함수를 매핑합니다.
    - `http.HandleFunc`를 사용하여 `/api/trending_keywords` 경로와 `handlers.TrendingKeywordsHandler` 함수를 매핑합니다.
    - `http.ListenAndServe`를 호출하여 웹 서버를 시작합니다. (예: `:8080` 포트)

### 3단계: 서버 실행 및 테스트

1.  **라이브러리 설치 (필요시):**
    ```bash
    go get github.com/joho/godotenv
    go get github.com/mmcdole/gofeed
    ```

2.  **서버 실행:**
    ```bash
    go run main.go
    ```

3.  **테스트:**
    - `curl`을 사용하여 API가 정상적으로 동작하는지 확인합니다.
    ```bash
    # 뉴스 콘텐츠 API 테스트 (단일 키워드)
    curl "http://localhost:8080/api/content?query=IT&start=1&display=10"

    # 뉴스 콘텐츠 API 테스트 (다중 키워드)
    curl "http://localhost:8080/api/content?query=IT,경제&start=1&display=10"

    # 트렌딩 키워드 API 테스트 (기본 5개)
    curl "http://localhost:8080/api/trending_keywords"

    # 트렌딩 키워드 API 테스트 (랜덤 1개)
    curl "http://localhost:8080/api/trending_keywords?random=true"
    ```

---

## 4. MVP 이후 확장 계획 (v2.0)

- **라우터 교체:** 기능이 복잡해지면 `gorilla/mux` 또는 `chi`와 같은 라우팅 라이브러리를 도입하여 미들웨어, 경로 변수 처리, RESTful API 설계 등을 용이하게 할 수 있습니다.
- **데이터베이스 연동:** 사용자 계정 및 통계 기능 추가 시, `database/sql` 표준 패키지와 PostgreSQL 또는 MySQL 드라이버를 사용하여 DB 연동 로직을 구현합니다.
- **캐싱:** Redis 같은 인메모리 저장소를 활용하여 Naver API 응답을 캐싱함으로써 성능을 향상시키고 API 호출 횟수를 줄일 수 있습니다.
- **Google Trends API 연동:** RSS 파싱 대신 Google Trends의 공식 API를 사용하여 더 풍부하고 정확한 트렌딩 키워드 데이터를 가져올 수 있습니다.
- **뉴스 기사 본문 파싱:** 현재는 기사 제목과 요약만 가져오지만, 기사 본문을 직접 파싱하여 더 긴 글감을 제공할 수 있습니다.
