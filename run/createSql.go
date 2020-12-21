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
		"\n\nINSERT INTO mov_mst(kr_title, en_title, rank, cnt, op_date, stat, age)\n" +
			"VALUES('" + items.Title + "', '" + items.TitleEng + "', " + rank + ", " + cnt + ", '" + items.OpenDate + "', '" + items.Stat + "', '" + items.Age + "')\n" +
			"ON DUPLICATE KEY UPDATE rank = " + rank + ", cnt =" + cnt + ", stat = '" + items.Stat + "', id = LAST_INSERT_ID(id);\n" +
			"SET @movid = LAST_INSERT_ID();\n" +
			"INSERT INTO mov_dt " +
			"VALUES (@movid, '" +
			items.Direct + "', '" + items.Actor + "', '" + items.Runtime + "', '" + items.Summary + "')\n" +
			"ON DUPLICATE KEY UPDATE\n director= '" +
			items.Direct + "', actor ='" + items.Actor + "', runtime ='" + items.Runtime + "', " +
			"summary = '" + items.Summary + "';\n" +

			func() string { // 장르는 따로 구분하여 넣기 때문에 슬라이싱하여 '장르의 수 만큼 구문을 생성'한다.
				s := strings.Split(items.GenreName, ",")
				gnr := ""
				for _, s := range s {
					s = strings.TrimSpace(s)
					gnr += "INSERT INTO mov_genre(mov_id, name) SELECT @movid, '" + s +
						"'\nFROM DUAL WHERE NOT EXISTS(SELECT mov_id, name FROM mov_genre\n" +
						"WHERE mov_id = @movid AND name='" + s + "');\n"
				}
				return gnr
			}() +

			"INSERT INTO mov_img(mov_id,img_path)\n" +
			"SELECT @movid, '" + items.ImgName[1:] +
			"'\nFROM DUAL WHERE NOT EXISTS(SELECT mov_id, img_path FROM mov_img\n" +
			"WHERE mov_id = @movid AND img_path='" + items.ImgName[1:] + "');\n" +

			func() string { // 상영 타입 또한 여러가지 타입이 동시에 존재할 수 있기에 익명함수를 통해 '타입갯수만큼 생성' 후 리턴
				s := strings.Split(items.TypeName, ",")
				typ := ""

				for _, s := range s {
					s = strings.TrimSpace(s)
					typ += "INSERT INTO type_mst(name) " +
						"VALUES('" + s + "')\n" +
						"ON DUPLICATE KEY UPDATE\nid = (@typeid := LAST_INSERT_ID(id));\n" +
						"INSERT INTO mov_type(mov_id,type_id)\n" +
						"SELECT @movid, @typeid\n" +
						"FROM DUAL WHERE NOT EXISTS(SELECT mov_id, type_id FROM mov_type\n" +
						"WHERE mov_id = @movid AND type_id=@typeid);\n"
				}
				return typ
			}(),
	)

	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.Title)
	return
}
