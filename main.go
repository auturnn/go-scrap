package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
)

type ParseHTML struct {
	url      string // https://www.megabox.co.kr/movie-detail?rpstMovieNo=
	respName string // 20046500
	htmlPath string
	htmlName string
	htmlFile *os.File
	imgPath  string
	imgName  string
	sqlPath  string
	sqlName  string
	sqlFile  *os.File

	doc *goquery.Document

	MovTest
}

//MovTest table structs
type MovTest struct {
	MovTitle []string
	MovOpd   []string
	MovType  []string
	MovGnr   []string
	MovAge   []string
	MovPst   []string
	MovDrt   []string
	MovRnt   []string
	MovPfm   []string
	MovSmr   []string
}

type Megabox struct {
	ImgPath   string
	ImgSvrUrl string      `json:"imgSvrUrl"`
	MovieList []MovieList `json:"movieList"`
}

type MovieList struct {
	MovieNo   string `json:"movieNo"`
	MovieNm   string `json:"movieNm"`
	ImgPathNm string `json:"imgPathNm"`
}

func main() {
	log.Println("main log...")

	m := &Megabox{
		ImgPath: "./img/poster/ing/",
	}
	m.ListData()

	var wait sync.WaitGroup
	wait.Add(len(m.MovieList))

	for i, items := range m.MovieList {
		log.Println(items)
		go func(i int, items MovieList) {
			defer wait.Done()
			m.PosterDown(i, items)
		}(i, items)
	}
	wait.Wait()
}

func (m *Megabox) ListData() {
	url := "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieList.do"
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
	req, err := http.NewRequest(method, url, payload)

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

	return
}

func (m *Megabox) PosterDown(i int, list MovieList) {
	client := grab.NewClient()
	request, _ := grab.NewRequest(m.ImgPath, m.ImgSvrUrl+list.ImgPathNm)
	response := client.Do(request)
	filename := response.Filename

	t := time.NewTicker(time.Second)

	for {
		select {
		case <-t.C:
			fmt.Printf("%.02f%% complete\n", response.Progress())
		case <-response.Done:
			err := response.Err()
			if err != nil {
				// ...
			}
			os.Rename(filename, m.ImgPath+list.MovieNm+".jpg")
			log.Println(i, "번째 Done!"+m.ImgPath+list.MovieNm+".jpg", response.HTTPResponse.StatusCode)
			return
		}
	}
}
