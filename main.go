package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cavaliercoder/grab"
)

type Megabox struct {
	Url       string
	ImgPath   string
	ImgSvrUrl string      `json:"imgSvrUrl"`
	MovieList []MovieList `json:"movieList"`
}

type MovieList struct {
	BoxoRank          int    `json:"boxoRank"`
	BoxoKofTotAdncCnt int    `json:"boxoKofTotAdncCnt"` // 누적관객
	MovieNo           string `json:"movieNo"`           // 영화 넘버
	MovieNm           string `json:"movieNm"`           // 영화이름
	ImgPathNm         string `json:"imgPathNm"`         // 포스터 이미지
	AdmisClassNm      string `json:"admisClassNm"`      // 등급
	MovieSynopCn      string `json:"movieSynopCn"`      // 설명
	PlayTime          string `json:"playTime"`          // 러닝타임
	RfilmDe           string `json:"rfilmDe"`           // 개봉일
	OnairYn           string `json:"onairYn"`           // Status
}

type Spec struct {
	Star  string
	Shap  string
	Excla string
	Alpa  string
	Enper string
}

func main() {
	log.Println("main log...")

	m := &Megabox{
		Url:     "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieList.do",
		ImgPath: "./img/poster/ing/",
	}
	m.GetList()

	file, err := os.OpenFile(
		"./sql/ing.sql",
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644),
	)
	if err != nil {
		log.Fatal("file Create Error: ", err)
	}

	var wait sync.WaitGroup
	wait.Add(len(m.MovieList))
	for _, items := range m.MovieList {
		go func(items MovieList) {
			defer wait.Done()
			if items.AdmisClassNm == "전체관람가" {
				items.AdmisClassNm = "all"
			} else {
				items.AdmisClassNm = strings.TrimRight(items.AdmisClassNm, "세이상관람가")
			}

			items.ImgPathNm = "/img/poster/ing" + items.MovieNm + ".jpg"
			rank := strconv.Itoa(items.BoxoRank)
			cnt := strconv.Itoa(items.BoxoKofTotAdncCnt)
			_, err = file.WriteString("--예매율랭킹" + rank + "\nINSERT INTO MOV_TEST(MOV_TITLE, MOV_CNT, MOV_OPD, MOV_CLD, MOV_AGE, MOV_PST, MOV_RNT, MOV_SMR)\n" +
				"VALUES('" + items.MovieNm + "', " + cnt + ", '" + items.RfilmDe + "', '" + items.OnairYn + "', '" + items.AdmisClassNm + "', '" +
				items.ImgPathNm + "', '" + items.PlayTime + "', '" + items.MovieSynopCn + "')\n",
			)
		}(items)

	}
	wait.Wait()
}

func (m *Megabox) GetList() {
	method := "POST"

	payload := strings.NewReader(`{
	"currentPage": "1",
	"ibxMovieNmSearch": "",
	"onairYn": "Y",
	"pageType": "ticketing",
	"recordCountPerPage": "40",
	"specialType": ""
	}`)

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
				m.PosterDown(items)
			}(items)
		}

		wait.Wait()
	}(m)
	return
}

func (m *Megabox) PosterDown(list MovieList) {
	client := grab.NewClient()
	request, _ := grab.NewRequest(m.ImgPath, m.ImgSvrUrl+list.ImgPathNm)
	response := client.Do(request)
	filename := response.Filename

	if strings.ContainsAny(list.MovieNm, ":") {
		list.MovieNm = strings.ReplaceAll(list.MovieNm, ":", "")
	} else if strings.ContainsAny(list.MovieNm, "/") {
		list.MovieNm = strings.ReplaceAll(list.MovieNm, "/", "")
	}

	if list.OnairYn == "MSC02" {
		list.OnairYn = "N"
	}

	for {
		select {
		case <-response.Done:
			err := response.Err()
			if err != nil {
				// ...
			}
			err = os.Rename(filename, m.ImgPath+list.MovieNm+".jpg")
			if err != nil {
				log.Println("Err:" + list.MovieNm)
			}
			log.Println("Done!"+m.ImgPath+list.MovieNm+".jpg", response.HTTPResponse.StatusCode)
			return
		}
	}
}
