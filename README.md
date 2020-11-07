## Go-Scrap (Scraper)

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
### 필요 폴더
- /sql : 작성된 SQL 저장

- /img
    - /poster
        - /ing : 상영작 포스터 저장
        - /pre : 상영예정작 포스터 저장
