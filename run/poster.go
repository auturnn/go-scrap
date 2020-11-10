package run

import (
	"log"

	"github.com/hashicorp/go-getter"
)

func (m *Megabox) PosterDown(list MovieList) {
	err := getter.GetFile(m.ImgPath+list.MovTitle+".jpg", m.ImgSvrUrl+list.MovPst)
	if err != nil {
		log.Println("img Error", m.ImgPath+list.MovTitle+".jpg")
	}
	// grab의 request는 무조건적으로 실행하지만, go-getter의 getFile은 이미 동일한 이름의 파일이 존재한다면 무시하고 넘어간다.
	// request, _ := grab.NewRequest(m.ImgPath, m.ImgSvrUrl+list.MovPst)
	// response := m.Client.Do(request)
	// filename := response.Filename

	// for {
	// 	select {
	// 	case <-response.Done:
	// 		err := response.Err()
	// 		if err != nil {
	// 			log.Println("Err:" + filename)
	// 		}
	// 		err = os.Rename(filename, m.ImgPath+"/"+list.MovTitle+".jpg")
	// 		if err != nil {
	// 			log.Println("Err:" + m.ImgPath + "/" + list.MovTitle + ".jpg")
	// 		}
	// 		// log.Println("IMG Done!"+m.ImgPath+list.MovTitle+".jpg", response.HTTPResponse.StatusCode)
	// 		return
	// 	}
	// }
}
