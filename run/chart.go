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

	m.SqlFile, _ = os.OpenFile(
		m.SqlFileName,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644),
	)

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
			m.DetailRequest(&items)
			m.CreateSQL(&items)
		}(items)
	}

	wait.Wait()

	return
}
