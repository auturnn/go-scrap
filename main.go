package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
)

type Megabox struct {
	Url         string
	ImgPath     string
	PayLoad     string
	ImgSvrUrl   string `json:"imgSvrUrl"`
	SqlFileName string
	SqlFile     *os.File
	MovieList   []MovieList `json:"movieList"`
}

type MovieList struct {
	MovRank  int    `json:"boxoRank"`          //예매율
	MovieNo  string `json:"movieNo"`           // 영화 넘버
	MovTitle string `json:"movieNm"`           // 영화이름
	MovCnt   int    `json:"boxoKofTotAdncCnt"` // 누적관객
	MovStat  string `json:"onairYn"`           // Status
	MovType  string // 2D, 3D
	MovGnr   string // 장르
	MovAge   string `json:"admisClassNm"` // 등급
	MovPst   string `json:"imgPathNm"`    // 포스터 이미지
	MovDrt   string // 감독
	MovRnt   string `json:"playTime"` // 러닝타임
	MovOpd   string `json:"rfilmDe"`  // 개봉일
	MovPfm   string // 출연진
	MovSmr   string `json:"movieSynopCn"` // 설명
}

func main() {
	log.Println("Start Log...")

	m := &Megabox{
		Url:     "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieList.do",
		ImgPath: "./img/poster/ing/",
		PayLoad: `{
			"currentPage": "1",
			"ibxMovieNmSearch": "",
			"onairYn": "Y",
			"pageType": "ticketing",
			"recordCountPerPage": "40",
			"specialType": ""
			}`,
		// 상영예정작 가져오기.
		// PayLoad: `{
		// 	"currentPage": "1",
		// 	"ibxMovieNmSearch": "",
		// 	"onairYn": "MSC02",
		// 	"pageType": "rfilmDe",
		// 	"recordCountPerPage": "40",
		// 	"specialType": ""
		// }`,
		SqlFileName: "./sql/movie-ing.sql",
	}

	m.GetList()

	m.SqlFile, _ = os.OpenFile(
		m.SqlFileName,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644),
	)

	var wait sync.WaitGroup
	wait.Add(len(m.MovieList))

	list := &m.MovieList
	for _, items := range *list {
		go func(items MovieList) {
			defer wait.Done()
			m.DetailRequest(&items)
			m.CreateSQL(&items)
		}(items)
	}
	wait.Wait()
}

func (m *Megabox) CreateSQL(items *MovieList) {
	if items.MovAge == "전체관람가" {
		items.MovAge = "all"
	} else {
		items.MovAge = strings.TrimRight(items.MovAge, "세이상관람가")
	}

	items.MovPst = "/img/poster/ing/" + items.MovTitle + ".jpg"
	rank := strconv.Itoa(items.MovRank)
	cnt := strconv.Itoa(items.MovCnt)
	_, err := m.SqlFile.WriteString("--예매율랭킹 " + rank + "위" + "\nINSERT INTO MOV_TEST(MOV_TITLE, MOV_CNT, MOV_OPD, MOV_STAT, MOV_TYPE, MOV_GNR, MOV_AGE, MOV_PST, MOV_RNT, MOV_PFM, MOV_SMR)\n" +
		"VALUES('" + items.MovTitle + "', " + cnt + ", '" + items.MovOpd + "', '" + items.MovStat + "', '" + items.MovType + "', '" + items.MovGnr + "', '" + items.MovAge + "', '" +
		items.MovPst + "', '" + items.MovRnt + "', '" + items.MovPfm + "', '" + items.MovSmr + "');\n",
	)
	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.MovTitle)
	return
}

func (m *Megabox) GetList() {
	method := "POST"

	payload := strings.NewReader(m.PayLoad)

	client := &http.Client{}
	req, err := http.NewRequest(method, m.Url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "WMONID=zcT2Xq9q57s; SameSite=None; SESSION=Yzc4OTJjYmYtNmU5NS00MGVhLWFmNDMtY2FkMWZkYjc3MjE1; JSESSIONID=oN4vAfK8XnaOkcCoDOZ1jxIY6CP27dVObnvHlnNx9IZomHQSVJWIWwQbaWDVIh3y.b25fbWVnYWJveF9kb21haW4vbWVnYS1vbi1zZXJ2ZXI1")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, m)
	if err != nil {
		log.Fatal(err)
	}

	func(m *Megabox) {
		var wait sync.WaitGroup
		wait.Add(len(m.MovieList))

		for _, items := range m.MovieList {
			go func(items MovieList) {
				defer wait.Done()

				if strings.ContainsAny(items.MovTitle, ":") {
					items.MovTitle = strings.ReplaceAll(items.MovTitle, ":", "")
				} else if strings.ContainsAny(items.MovTitle, "/") {
					items.MovTitle = strings.ReplaceAll(items.MovTitle, "/", "")
				}

				if items.MovRnt == "MSC02" {
					items.MovRnt = "N"
				}

				m.PosterDown(items)
			}(items)
		}

		wait.Wait()
	}(m)

	return
}

func (m *Megabox) PosterDown(list MovieList) {
	client := grab.NewClient()
	request, _ := grab.NewRequest(m.ImgPath, m.ImgSvrUrl+list.MovPst)
	response := client.Do(request)
	filename := response.Filename

	for {
		select {
		case <-response.Done:
			err := response.Err()
			if err != nil {
				// ...
			}
			err = os.Rename(filename, m.ImgPath+"/"+list.MovTitle+".jpg")
			if err != nil {
				log.Println("Err:" + list.MovTitle)
			}
			// log.Println("Done!"+m.ImgPath+list.MovTitle+".jpg", response.HTTPResponse.StatusCode)
			return
		}
	}
}

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
			list.MovType = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "감독") {
			s.SetText(strings.Replace(s.Text(), "감독", "", -1))
			list.MovDrt = strings.TrimSpace(s.Text())
		} else if strings.Contains(s.Text(), "장르") {
			s.SetText(strings.Replace(s.Text(), "장르", "", -1))
			list.MovGnr = strings.TrimSpace(strings.Split(s.Text(), "/")[0])
		} else if strings.Contains(s.Text(), "출연진") {
			s.SetText(strings.Replace(s.Text(), "출연진", "", -1))
			list.MovPfm = strings.TrimSpace(s.Text())
		}
	})
	return
}
