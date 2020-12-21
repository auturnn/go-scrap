package db

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type MovMst struct {
	MovID       int
	MovTitle    string
	MovEngTitle string
	MovRank     string
	MovCnt      int
	MovOpd      time.Time
	MovStat     string
	MovAge      string
}

type MovDt struct {
	MovID   int
	MdtDrct string
	MdtAct  string
	MdtLen  string
	MdtSmr  string
}

type MovGenre struct {
	MovId    int
	MgnrName string
}

type MovType struct {
}

type MovImg struct {
}

type TheaMst struct {
}

type TheaRoom struct {
}

type City struct {
}
