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
			"DELETE FROM user_test;\n" +
				"DELETE FROM reserv_mst;\n" +
				"DELETE FROM rest_cnt;\n" +
				"DELETE FROM play_mst;\n" +
				"DELETE FROM seat_mst;\n" +
				"DELETE FROM room_mst;\n" +
				"DELETE FROM thea_mst;\n" +
				"DELETE FROM city_mst;\n" +
				"DELETE FROM mov_type;\n" +
				"DELETE FROM type_mst;\n" +
				"DELETE FROM mov_genre;\n" +
				"DELETE FROM mov_img;\n" +
				"DELETE FROM mov_dt;\n" +
				"DELETE FROM mov_mst;\n" +
				"ALTER TABLE type_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE play_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE city_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE reserv_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE rest_cnt AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE seat_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE room_mst AUTO_INCREMENT = 1;\n" +
				"ALTER TABLE thea_mst AUTO_INCREMENT = 1;\n" +
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
				ext := items.ImgPath[len(items.ImgPath)-4:] //이미지마다 확장자가 다르기 때문에 ex) .jpg / .gif / .png
				items.ImgName = m.ImgPath + items.Title + ext
				m.PosterDown(items)
				m.DetailRequest(&items)
				m.CreateSQL(&items)
			} else {
				log.Println("[PASS] IMG ERROR :", items.Title)
			}
		}(items)
	}

	wait.Wait()

	return
}
