package run

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//DetailRequest is Request for detail movie's info
func (m *Megabox) DetailRequest(list *MovieList) {
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
	req.Header.Add("Cookie", "WMONID=zcT2Xq9q57s; SESSION=MWQyNjkwNGItM2UzNy00ODYyLTk5NjMtYTAxMDg0ODc4OTA4; JSESSIONID=jB8RGqjtL1gk798CwmTE1cu2QFDZqPePnXeN6lOpkUtl4qEOVbG6X0wNY9YQceL9.b25fbWVnYWJveF9kb21haW4vbWVnYS1vbi1zZXJ2ZXIy; SameSite=None")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	list.DetailParse(res.Body)

	return
}

func (list *MovieList) DetailParse(body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Create Doc Error")
	}

	info := doc.Find("div.inner-wrap").Find("div.movie-info.infoContent")
	info.Find("p").Each(func(i int, s *goquery.Selection) {
		s.SetText(strings.Replace(s.Text(), ":", "", -1))
		if strings.Contains(s.Text(), "상영타입") {
			s.SetText(strings.Replace(s.Text(), "상영타입", "", -1))
			list.MtypeName = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "감독") {
			s.SetText(strings.Replace(s.Text(), "감독", "", -1))
			list.MdtDrt = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "장르") {
			s.SetText(strings.Replace(s.Text(), "장르", "", -1))
			list.MgnrName = strings.TrimSpace(strings.Split(s.Text(), "/")[0])
		} else if strings.Contains(s.Text(), "출연진") {
			s.SetText(strings.Replace(s.Text(), "출연진", "", -1))
			list.MdtAct = strings.TrimSpace(s.Text())
		}
	})
	return
}
