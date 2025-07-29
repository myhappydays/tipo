# Tipo 프로젝트 백엔드 개발 계획 (Golang)

- **문서 버전:** 1.0
- **작성일:** 2025-07-29
- **기술 스택:** Golang (Go)

---

## 1. MVP 목표

- Naver API 키 등 민감한 정보를 서버에서 안전하게 관리한다.
- 클라이언트(프론트엔드)의 요청에 따라 Naver API(도서, 영화, 뉴스)를 대신 호출하는 프록시(Proxy) 서버 역할을 수행한다.
- Naver API로부터 받은 데이터를 타자 연습에 적합한 형태로 가공하여 일관된 JSON 형식으로 제공한다.

---

## 2. 기술 명세 (Technical Specifications)

- **언어:** Go (Golang)
- **웹 프레임워크:** 표준 라이브러리 `net/http` 사용 (외부 프레임워크 없이 구현)
    - *이유: MVP 기능은 간단한 API 엔드포인트 하나로 충분하므로, 가볍고 빠른 표준 라이브러리가 가장 효율적입니다.*
- **라우팅:** `http.HandleFunc`를 사용한 기본 라우터
- **환경 변수 관리:** `os.Getenv` (표준 라이브러리) 또는 `godotenv` 라이브러리 활용

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
        └── content.go
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
        Source    string   `json:"source"`
        Category  string   `json:"category"`
        Image     string   `json:"image"`
        Sentences []string `json:"sentences"`
    }
    ```

2.  **API 핸들러 구현 (`/handlers/content.go`):**
    - `/api/content` 엔드포인트의 로직을 작성합니다.
    - 클라이언트로부터 `category`와 `query` 파라미터를 받습니다.
    - `os.Getenv`를 사용해 환경 변수에서 Naver API 키를 읽어옵니다.
    - `net/http` 클라이언트를 사용해 Naver API를 호출하고 헤더를 설정합니다.
    - 응답 받은 JSON 데이터를 파싱하여 필요한 정보만 추출합니다.
    - `utils` 패키지의 함수를 호출하여 텍스트를 정제합니다.
    - 최종적으로 `models.ContentResponse` 구조체에 데이터를 담아 JSON으로 인코딩하여 응답합니다.

3.  **HTML 태그 제거 유틸리티 (`/utils/clean.go`):**
    - 정규식(`regexp` 패키지)을 사용하여 문자열에서 HTML 태그를 제거하는 간단한 함수를 작성합니다.

4.  **메인 함수 작성 (`main.go`):**
    - `http.HandleFunc`를 사용하여 `/api/content` 경로와 `handlers.ContentHandler` 함수를 매핑합니다.
    - `http.ListenAndServe`를 호출하여 웹 서버를 시작합니다. (예: `:8080` 포트)

### 3단계: 서버 실행 및 테스트

1.  **라이브러리 설치 (필요시):**
    ```bash
    go get github.com/joho/godotenv # .env 파일 로드를 위해
    ```

2.  **서버 실행:**
    ```bash
    go run main.go
    ```

3.  **테스트:**
    - 웹 브라우저나 `curl`을 사용하여 API가 정상적으로 동작하는지 확인합니다.
    ```bash
    curl "http://localhost:8080/api/content?category=book&query= 어린왕자"
    ```

---

## 4. MVP 이후 확장 계획 (v2.0)

- **라우터 교체:** 기능이 복잡해지면 `gorilla/mux` 또는 `chi`와 같은 라우팅 라이브러리를 도입하여 미들웨어나 경로 변수 처리를 용이하게 할 수 있습니다.
- **데이터베이스 연동:** 사용자 계정 및 통계 기능 추가 시, `database/sql` 표준 패키지와 PostgreSQL 또는 MySQL 드라이버를 사용하여 DB 연동 로직을 구현합니다.
- **캐싱:** `sync.Map`이나 Redis 같은 인메모리 저장소를 활용하여 Naver API 응답을 캐싱함으로써 성능을 향상시키고 API 호출 횟수를 줄일 수 있습니다.
