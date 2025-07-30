import os
import requests
import json
import re
import html
from dotenv import load_dotenv

# .env 파일에서 환경 변수 로드
load_dotenv()

# 네이버 API 자격 증명
NAVER_CLIENT_ID = os.getenv("NAVER_CLIENT_ID")
NAVER_CLIENT_SECRET = os.getenv("NAVER_CLIENT_SECRET")

# API URL
NAVER_API_URL = "https://openapi.naver.com/v1/search/news.json"

def clean_text(text):
    """HTML 태그를 제거하고 HTML 엔티티를 변환합니다."""
    # HTML 엔티티 변환 (e.g., &quot; -> ")
    text = html.unescape(text)
    # HTML 태그 제거
    cleaner = re.compile('<.*?>')
    clean_text = re.sub(cleaner, '', text)
    return clean_text

def fetch_naver_news(query, display=10):
    """네이버 뉴스 API를 호출하여 뉴스 데이터를 가져옵니다."""
    if not NAVER_CLIENT_ID or not NAVER_CLIENT_SECRET:
        print("오류: NAVER_CLIENT_ID와 NAVER_CLIENT_SECRET 환경 변수가 필요합니다.")
        return None

    headers = {
        "X-Naver-Client-Id": NAVER_CLIENT_ID,
        "X-Naver-Client-Secret": NAVER_CLIENT_SECRET,
    }
    params = {
        "query": query,
        "display": display,
        "start": 1,
    }

    try:
        response = requests.get(NAVER_API_URL, headers=headers, params=params)
        response.raise_for_status()  # HTTP 오류 발생 시 예외 발생
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"API 요청 중 오류 발생: {e}")
        return None

def process_and_save(api_response, query, filename):
    """API 응답을 처리하고 지정된 파일에 JSON으로 저장합니다."""
    if not api_response or "items" not in api_response:
        print("API 응답에서 'items'를 찾을 수 없습니다.")
        return

    output_data = []
    for item in api_response["items"]:
        # 원본 링크에서 도메인 추출
        source_url = item.get("originallink", "")

        processed_item = {
            "title": clean_text(item["title"]),
            "source": source_url,
            "keyword": query,
            "pubDate": item["pubDate"],
            "sentences": [clean_text(item["description"])] # description을 그대로 문장으로 사용
        }
        output_data.append(processed_item)

    # 파일에 저장
    try:
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(output_data, f, indent=2, ensure_ascii=False)
        print(f"성공적으로 {len(output_data)}개의 뉴스 기사를 '{filename}'에 저장했습니다.")
    except IOError as e:
        print(f"파일 저장 중 오류 발생: {e}")

if __name__ == "__main__":
    SEARCH_QUERY = "IT"
    RESULT_COUNT = 10
    OUTPUT_FILENAME = "test/naver_news_it_results.json"

    print(f"'{SEARCH_QUERY}' 키워드로 {RESULT_COUNT}개의 뉴스 기사를 가져옵니다...")
    
    # .env 파일이 현재 스크립트와 같은 디렉토리에 있는지 확인
    dotenv_path = os.path.join(os.path.dirname(__file__), '.env')
    if os.path.exists(dotenv_path):
        load_dotenv(dotenv_path=dotenv_path)
        print(".env 파일에서 환경 변수를 로드했습니다.")
    else:
        # 루트 폴더에서도 찾아보기
        dotenv_path_root = os.path.join(os.path.dirname(__file__), '..', '.env')
        if os.path.exists(dotenv_path_root):
            load_dotenv(dotenv_path=dotenv_path_root)
            print("프로젝트 루트 폴더의 .env 파일에서 환경 변수를 로드했습니다.")
        else:
            print("경고: .env 파일을 찾을 수 없습니다. 환경 변수가 시스템에 설정되어 있는지 확인하세요.")


    news_data = fetch_naver_news(SEARCH_QUERY, RESULT_COUNT)

    if news_data:
        process_and_save(news_data, SEARCH_QUERY, OUTPUT_FILENAME)
