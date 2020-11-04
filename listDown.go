package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

// //HTMLDown is Get html and Create a htmlFile
func (p *ParseHTML) HTMLDown() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var res string

	err := chromedp.Run(ctx,
		chromedp.Navigate(p.url),
		// wait for footer element is visible (ie, page is loaded)

		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			checkErr(err)
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	checkErr(err)
	file, _ := os.Create(p.htmlPath + p.htmlName)
	defer file.Close()

	_, err = file.WriteString(res)
	checkErr(err)
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

//IOHtml read and save to ParseHTML.doc
func (p *ParseHTML) IOHtml() {
	p.htmlFile, _ = os.Open(p.htmlName)

	defer p.htmlFile.Close()

	f1, err := ioutil.ReadFile(p.htmlName)
	if err != nil {
		log.Fatal(err)
	}

	rhtml := strings.NewReader(string(f1))

	p.doc, err = goquery.NewDocumentFromReader(rhtml)
	if err != nil {
		panic(err)
	}

	return
}

//PostDown is megabox movie poster img download And Renaming
func (p *ParseHTML) PostDown() {
	log.Println("영화 카운트: ", p.doc.Find("div.movie-list").Find("ol").Find("li").Length())
	client := grab.NewClient()

	postFind := p.doc.Find("div.movie-list").Find("ol").Find("li")
	postFind.Each(func(i int, s *goquery.Selection) {
		imgPath, _ := s.Find("img").Attr("src")
		title := s.Find("div.tit-area").Find("p.tit").Text()

		request, _ := grab.NewRequest(p.imgPath, imgPath)
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
				p.MovTitle = append(p.MovTitle, title)
				p.MovPst = append(p.MovPst, p.imgPath+title+".jpg")
				os.Rename(filename, p.imgPath+title+".jpg")
				os.Rename(p.htmlName, p.htmlPath+title+".html")
				log.Println(imgPath, response.HTTPResponse.StatusCode)
				return
			}
		}
	})
	return
}

// //TextParse hi
// func (p *ParseHTML) TextInject() (err error) {
// 	p.sqlFile, err = os.Create(p.sqlPath + p.sqlName)
// 	if err != nil {
// 		log.Fatal("file open Error")
// 	}
// 	defer p.sqlFile.Close()
// 	for _, title := range p.MovTitle {
// 		log.Println("title:", title)
// 		_, err := p.sqlFile.WriteString("INSERT INTO MOV_TEST(MOV_TITLE) VALUES ('" + title + "')\n")
// 		if err != nil {
// 			log.Fatal("file Write err:", err)
// 		}
// 	}

// 	return err
// }
