package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (m *Megabox) Pss() {
	// var p ParseHTML
	url := "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieInfo.do"
	method := "POST"

	payload := strings.NewReader(`{rpstMovieNo:	"`)

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
	log.Println(body)
	err = json.Unmarshal(body, m)
	if err != nil {
		log.Fatal(err)
	}

	return
}
