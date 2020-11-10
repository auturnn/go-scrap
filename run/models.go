package run

import (
	"os"

	"github.com/cavaliercoder/grab"
)

type Megabox struct {
	Host        string
	ImgPath     string
	PayLoad     string
	ImgSvrUrl   string `json:"imgSvrUrl"`
	SqlFileName string
	SqlFile     *os.File
	Client      *grab.Client
	MovieList   []MovieList `json:"movieList"`
}

type MovieList struct {
	MovRank  int    `json:"boxoRank"`          //예매율
	MovieNo  string `json:"movieNo"`           // 영화 넘버
	MovTitle string `json:"movieNm"`           // 영화이름
	MovCnt   int    `json:"boxoKofTotAdncCnt"` // 누적관객
	MovStat  string `json:"onairYn"`           // Status
	MovType  string // 2D, 3D
	MovGnr   string // 장르
	MovAge   string `json:"admisClassNm"` // 등급
	MovPst   string `json:"imgPathNm"`    // 포스터 이미지
	MovDrt   string // 감독
	MovRnt   string `json:"playTime"` // 러닝타임
	MovOpd   string `json:"rfilmDe"`  // 개봉일
	MovPfm   string // 출연진
	MovSmr   string `json:"movieSynopCn"` // 설명
}
