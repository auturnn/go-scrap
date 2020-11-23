package run

import (
	"log"
	"os"

	"github.com/cavaliercoder/grab"
)

//PosterDown is download for movie main poster
func (m *Megabox) PosterDown(list MovieList) {
	//해당 이미지가 이미 있는 지 확인. false일 경우 retrun
	if _, err := os.Stat(list.MimgName); !os.IsNotExist(err) {
		log.Println("Img already :" + list.MimgName)
		return
	}

	// go-getter 사용시의 코드
	// err := getter.GetFile(list.MimgName, m.ImgSvrUrl+list.MimgPath)
	// if err != nil {
	// 	log.Println("Img DownLoad Error :", list.MimgName)
	// }

	request, _ := grab.NewRequest(m.ImgPath, m.ImgSvrURL+list.MimgPath)
	response := m.Client.Do(request)
	filename := response.Filename

	for {
		select {
		case <-response.Done:
			err := response.Err()
			if err != nil {
				log.Println("Poster download Err:" + filename)
			}

			err = os.Rename(filename, list.MimgName)
			if err != nil {
				log.Println("Poster Rename Err:" + list.MimgName)
			}
			// log.Println("IMG Done!"+m.ImgPath+list.MovTitle+".jpg", response.HTTPResponse.StatusCode)
			return
		}
	}
}
