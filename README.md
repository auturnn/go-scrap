## Go-Scrap (Scraper)

[![Go Report Card](https://goreportcard.com/badge/github.com/auturnn/go-scrap)](https://goreportcard.com/report/github.com/auturnn/go-scrap)
---

이 코드는 포트폴리오용 영화 페이지를 만들기 위해\
필요한 자원의 수집을 자동화하기 위해 개발하였음을 알립니다.

해당 소스는 메가박스의 영화 리스트, 영화 정보페이지를 crawling하여\
메인 포스터 이미지를 다운하고, 세부 정보들을 SQL로 작성 및 저장한다.

---

- Language : Golang
- Framework :
    - "github.com/PuerkitoBio/goquery"
    - "github.com/cavaliercoder/grab"

---
### Folder Struct
- /run : jsonUnmashal, Crawling등의 정보수집을 위해 실행되는 소스 파일이 위치.
- /db : DB 연결 및 ORM을 이용한 DB 조작 소스
- /sql : 작성된 SQL 저장
- /img
    - /poster
        - /ing : 상영작 포스터 저장
        - /pre : 상영예정작 포스터 저장
---
### 업데이트 계획
1. SQL문 작성 이후 문서화 -> DB에 SQL문 작성
2. Spring Project에 이식
3. 서버상에서 crontab을 이용한 자동화

