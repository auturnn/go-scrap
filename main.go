package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
)

type ParseHTML struct {
	url      string
	fileName string
	downPath string
}

func main() {
	log.Println("main log...")
	p := &ParseHTML{
		url:      "https://www.megabox.co.kr/movie/comingsoon",
		fileName: "./ing.html",
		downPath: "./img/ing/",
	}
	// p.HTMLDown()
	p.PostDown()
}

func (p *ParseHTML) HTMLDown() {
	res, err := http.Get(p.url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code: %d %s", res.StatusCode, res.Status)
	}

	f1, err := os.Create(p.fileName)
	if err != nil {
		f1.Close()
		log.Fatal("파일 생성 실패", err)
	}

	doc := htmlParse(res.Body)
	html, _ := doc.Html()

	_, err = f1.WriteString(html)
	if err != nil {
		log.Fatal(err)
	}

	f1.Close()

}

func (p *ParseHTML) PostDown() {
	f1, err := ioutil.ReadFile(p.fileName)
	if err != nil {
		log.Fatal(err)
	}
	rhtml := strings.NewReader(string(f1))
	doc := htmlParse(rhtml)

	log.Println("영화 카운트: ", doc.Find("div.movie-list").Find("ol").Find("li").Length())
	client := grab.NewClient()

	postFind := doc.Find("div.movie-list").Find("ol").Find("li")
	postFind.Each(func(i int, s *goquery.Selection) {
		imgPath, _ := s.Find("img").Attr("src")
		title := s.Find("div.tit-area").Find("p.tit").Text()
		log.Println("title: " + title)

		request, _ := grab.NewRequest(p.downPath, imgPath)
		response := client.Do(request)
		filename := response.Filename
		log.Println(filename)

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
				err = os.Rename(filename, p.downPath+title+".jpg")
				if err != nil {
					log.Fatal(err)
				}

				log.Println(imgPath, response.HTTPResponse.StatusCode)
				return
			}
		}
	})
}

func htmlParse(rhtml io.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(rhtml)
	if err != nil {
		panic(err)
	}
	return doc
}
