package run

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//DetailRequest is Request for detail movie's info
func (m *Megabox) DetailRequest(list *MovieList) {
	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()

		url := "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieInfo.do"
		method := "POST"

		payload := strings.NewReader(`{
			"rpstMovieNo":	"` + list.MovieNo + `"
		}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// req.Header.Add("Cookie", "WMONID=zcT2Xq9q57s; SESSION=MWQyNjkwNGItM2UzNy00ODYyLTk5NjMtYTAxMDg0ODc4OTA4; JSESSIONID=jB8RGqjtL1gk798CwmTE1cu2QFDZqPePnXeN6lOpkUtl4qEOVbG6X0wNY9YQceL9.b25fbWVnYWJveF9kb21haW4vbWVnYS1vbi1zZXJ2ZXIy; SameSite=None")

		pRes, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer pRes.Body.Close()

		list.PostParse(pRes.Body)
	}()

	go func() {
		defer wait.Done()
		gRes, _ := http.Get("https://www.megabox.co.kr/movie-detail?rpstMovieNo=" + list.MovieNo)
		list.GetParse(gRes.Body)
		defer gRes.Body.Close()
	}()

	wait.Wait()

	return
}

//GetParse is Parse to movie's Eng-Title
func (list *MovieList) GetParse(body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Create getDoc Error")
	}

	TitleEng := doc.Find("div#contents").Find("div.movie-detail-page").Find("div.movie-detail-cont").Find("p.title-eng").Text()
	cnt := strings.Count(TitleEng, "'")
	if cnt == 1 { // "'" 하나있을 경우 ''으로 고쳐 SQL Error 제거
		list.TitleEng = strings.Replace(TitleEng, "'", "''", 1)
	} else if cnt%2 != 0 { // cnt가 1이 아닌 홀수인 경우 해당 위치를 정확하게 알수없기에 공백처리(에러방지)
		list.TitleEng = ""
	} else if cnt == 0 { // 이외의 정상적으로 처리될 수 있는 문자열들
		list.TitleEng = TitleEng
	}
}

//PostParse is Parse to movie's detail info
func (list *MovieList) PostParse(body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Create Doc Error")
	}
	info := doc.Find("div.inner-wrap").Find("div.movie-info.infoContent")
	// log.Println(info.Html())
	info.Find("p").Each(func(_ int, s *goquery.Selection) {
		s.SetText(strings.Replace(s.Text(), ":", "", -1))
		if strings.Contains(s.Text(), "상영타입") {
			s.SetText(strings.Replace(s.Text(), "상영타입", "", -1))
			list.TypeName = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "감독") {
			s.SetText(strings.Replace(s.Text(), "감독", "", -1))
			list.Direct = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "장르") {
			s.SetText(strings.Replace(s.Text(), "장르", "", -1))
			list.GenreName = strings.TrimSpace(strings.Split(s.Text(), "/")[0])
		} else if strings.Contains(s.Text(), "출연진") {
			s.SetText(strings.Replace(s.Text(), "출연진", "", -1))
			list.Actor = strings.TrimSpace(s.Text())
		}
	})
	return
}
