package run

import (
	"log"
	"strconv"
	"strings"
)

//CreateSQL is saved .sql file for mysql
func (m *Megabox) CreateSQL(items *MovieList) {
	if items.Age == "전체관람가" {
		items.Age = "all"
	} else if items.Age == "청소년관람불가" {
		items.Age = "18"
	} else {
		items.Age = strings.TrimRight(items.Age, "세이상관람가")
	}

	if items.Stat == "MSC01" {
		items.Stat = "Y"
	} else if items.Stat == "MSC02" {
		items.Stat = "L"
	}

	rank := strconv.Itoa(items.MovRank)
	cnt := strconv.Itoa(items.Cnt)
	_, err := m.SQLFile.WriteString(
		"\n\nINSERT INTO mov_mst(title, rank, cnt, open_date, stat, age)\n" +
			"VALUES('" + items.Title + "', " + rank + ", " + cnt + ", '" + items.OpenDate + "', '" + items.Stat + "', '" + items.Age + "')\n" +
			"ON DUPLICATE KEY UPDATE rank = " + rank + ", cnt =" + cnt + ", stat = '" + items.Stat + "', mov_idx = LAST_INSERT_ID(mov_idx);\n" +

			"INSERT INTO mov_dt " +
			"VALUES((SELECT LAST_INSERT_ID()), '" +
			items.Direct + "', '" + items.Actor + "', '" + items.Runtime + "', '" + items.Summary + "')\n" +
			"ON DUPLICATE KEY UPDATE\nmov_idx = LAST_INSERT_ID(mov_idx), director= '" +
			items.Direct + "', actor='" + items.Actor + "', runtime='" + items.Runtime + "', " +
			"summary= '" + items.Summary + "';\n" +

			func() string { // 장르는 따로 구분하여 넣기 때문에 슬라이싱하여 '장르의 수 만큼 구문을 생성'한다.
				s := strings.Split(items.GenreName, ",")
				gnr := ""
				for _, s := range s {
					s = strings.TrimSpace(s)
					gnr += "INSERT INTO mov_genre(mov_idx,genre_name) SELECT (SELECT LAST_INSERT_ID()),'" + s +
						"'\nFROM DUAL WHERE NOT EXISTS(SELECT mov_idx, genre_name FROM mov_genre\n" +
						"WHERE mov_idx = (SELECT LAST_INSERT_ID()) AND genre_name='" + s + "');\n"
				}
				return gnr
			}() +

			func() string { // 상영 타입 또한 여러가지 타입이 동시에 존재할 수 있기에 익명함수를 통해 '타입갯수만큼 생성' 후 리턴
				s := strings.Split(items.TypeName, ",")
				typ := ""
				mtyp := ""
				for _, s := range s {
					s = strings.TrimSpace(s)

					typ += "INSERT INTO type_mst(type_name) SELECT ('" + s +
						"' FROM DUAL WHERE NOT EXISTS(SELECT type_name FROM type_mst " +
						"WHERE type_name = '" + s + "');\n"

					mtyp += "INSERT INTO mov_type(mov_idx,type_idx) SELECT (SELECT LAST_INSERT_ID()),'" + s +
						"' FROM DUAL WHERE NOT EXISTS(SELECT mov_idx, type_idx FROM mov_type " +
						"WHERE mov_idx = (SELECT LAST_INSERT_ID()) AND type_idx='" + s + "');\n"
				}
				return typ + mtyp
			}() +

			"INSERT INTO MOV_IMG(mov_idx,MIMG_PATH)\n" +
			"SELECT (SELECT LAST_INSERT_ID()), '" + items.ImgName[1:] +
			"'\nFROM DUAL WHERE NOT EXISTS(SELECT mov_idx, MIMG_PATH FROM MOV_IMG\n" +
			"WHERE mov_idx = (SELECT LAST_INSERT_ID()) AND MIMG_PATH='" + items.ImgName[1:] + "');\n",
	)

	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.Title)
	return
}
