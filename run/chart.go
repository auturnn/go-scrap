package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

//GetList is parsing of Megabox Movie List json. And after running go-rootine for other Function.
func (m *Megabox) GetList() {
	method := "POST"
	client := &http.Client{}
	req, err := http.NewRequest(method, m.Host, strings.NewReader(m.PayLoad))
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

	m.SQLFile, _ = os.OpenFile(
		m.SQLFileName,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644),
	)
	if m.SQLFileName == "./sql/movie-ing.sql" {
		_, err = m.SQLFile.WriteString(
			"DELETE FROM mov_type;\n" +
				"DELETE FROM mov_genre;\n" +
				"DELETE FROM mov_img;\n" +
				"DELETE FROM mov_dt;\n" +
				"DELETE FROM mov_mst;\n" +
				"ALTER TABLE mov_mst AUTO_INCREMENT = 1;\n",
		)
	}

	var wait sync.WaitGroup
	wait.Add(len(m.MovieList))

	for _, items := range m.MovieList {
		go func(items MovieList) {
			defer wait.Done()

			if strings.ContainsAny(items.Title, ":") {
				items.Title = strings.ReplaceAll(items.Title, ":", "")
			} else if strings.ContainsAny(items.Title, "/") {
				items.Title = strings.ReplaceAll(items.Title, "/", "")
			}

			if items.Runtime == "MSC02" {
				items.Runtime = "N"
			}

			//정상적인 이미지 경로의 경우 모두 길이가 58로 고정.
			if len(items.ImgPath) == 58 {
				ext := items.ImgPath[len(items.ImgPath)-4:]
				items.ImgName = m.ImgPath + items.Title + ext
				m.PosterDown(items)
				m.DetailRequest(&items)
				m.CreateSQL(&items)
			} else {
				log.Println("Pass :", items.Title)
			}
		}(items)
	}

	wait.Wait()

	return
}
