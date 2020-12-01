package run

import (
	"os"

	"github.com/cavaliercoder/grab"
)

//Megabox struct is setting for Megabox Movie List json
type Megabox struct {
	Host        string
	ImgPath     string
	PayLoad     string
	ImgSvrURL   string
	SQLFileName string
	SQLFile     *os.File
	Client      *grab.Client
	MovieList   []MovieList `json:"movieList"`
}

//MovieList is saved json data for movie list
type MovieList struct {
	MovRank int    `json:"boxoRank"` //예매율
	MovieNo string `json:"movieNo"`  // 영화 넘버
	MovMst
	MovDT
	MovGenre
	MovType
	MovImg
}

//MovMst is DB struct(table column)
type MovMst struct {
	Title    string `json:"movieNm"`           // 영화이름
	Cnt      int    `json:"boxoKofTotAdncCnt"` // 누적관객
	OpenDate string `json:"rfilmDe"`           // 개봉일
	Stat     string `json:"onairYn"`           // Status
	Age      string `json:"admisClassNm"`      // 등급
}

//MovDT is DB struct(table column)
type MovDT struct {
	Direct  string // 감독
	Actor   string // 출연진
	Runtime string `json:"playTime"`     // 러닝타임
	Summary string `json:"movieSynopCn"` // 설명
}

//MovGenre is DB struct(table column)
type MovGenre struct {
	GenreName string //장르이름
}

//MovType is DB struct(table column)
type MovType struct {
	TypeName string // 2D, 3D등
}

//MovImg is DB struct(table column)
type MovImg struct {
	ImgPath string `json:"imgPathNm"` //포스터 저장 위치
	ImgName string
}

type TypeMst struct {
	TypeName string
}
