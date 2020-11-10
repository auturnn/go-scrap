package run

import (
	"sync"

	"github.com/cavaliercoder/grab"
)

func Run() {
	list := map[string]map[string]string{
		"ing": map[string]string{
			"ImgPath": "./img/poster/ing/",
			"PayLoad": `{
				"currentPage": "1",
				"ibxMovieNmSearch": "",
				"onairYn": "Y",
				"pageType": "ticketing",
				"recordCountPerPage": "40",
				"specialType": ""
			}`,
			"SqlFileName": "./sql/movie-ing.sql",
		},
		"pre": map[string]string{
			"ImgPath": "./img/poster/pre/",
			"PayLoad": `{
				"currentPage": "1",
		 	"ibxMovieNmSearch": "",
		 	"onairYn": "MSC02",
		 	"pageType": "rfilmDe",
		 	"recordCountPerPage": "40",
		 	"specialType": ""
			}`,
			"SqlFileName": "./sql/movie-pre.sql",
		},
	}
	var wait sync.WaitGroup
	wait.Add(2)
	for _, items := range list {
		go func(items map[string]string) {
			defer wait.Done()
			m := Megabox{
				Host:        "https://www.megabox.co.kr/on/oh/oha/Movie/selectMovieList.do",
				ImgPath:     items["ImgPath"],
				PayLoad:     items["PayLoad"],
				SqlFileName: items["SqlFileName"],
				Client:      grab.NewClient(),
			}
			m.GetList()
		}(items)
	}
	wait.Wait()
}
