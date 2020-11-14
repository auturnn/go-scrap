package run

import (
	"log"
	"strconv"
	"strings"
)

func (m *Megabox) CreateSQL(items *MovieList) {
	if items.MovAge == "전체관람가" {
		items.MovAge = "all"
	} else if items.MovAge == "청소년관람불가" {
		items.MovAge = "18"
	} else {
		items.MovAge = strings.TrimRight(items.MovAge, "세이상관람가")
	}

	if items.MovStat == "MSC01" {
		items.MovStat = "Y"
	} else if items.MovStat == "MSC02" {
		items.MovStat = "L"
	}

	rank := strconv.Itoa(items.MovRank)
	cnt := strconv.Itoa(items.MovCnt)
	_, err := m.SqlFile.WriteString(
		"\n\nINSERT INTO MOV_MST(MOV_TITLE, MOV_RANK, MOV_CNT, MOV_OPD, MOV_STAT, MOV_AGE)\n" +
			"VALUES('" + items.MovTitle + "', " + rank + ", " + cnt + ", '" + items.MovOpd + "', '" + items.MovStat + "', '" + items.MovAge + "')\n" +
			"ON DUPLICATE KEY UPDATE MOV_RANK = " + rank + ", MOV_CNT =" + cnt + ", MOV_STAT = '" + items.MovStat + "', MOV_SEQ = LAST_INSERT_ID(MOV_SEQ);\n" +

			"INSERT INTO MOV_DT " +
			"VALUES((SELECT LAST_INSERT_ID()), '" +
			items.MovDrt + "', '" + items.MovAct + "', '" + items.MovLen + "', '" + items.MovSmr + "')\n" +
			"ON DUPLICATE KEY UPDATE\nMOV_SEQ = LAST_INSERT_ID(MOV_SEQ), MDT_DRCT= '" +
			items.MovDrt + "', MDT_ACT='" + items.MovAct + "', MDT_LEN='" + items.MovLen + "', " +
			"MDT_SMR= '" + items.MovSmr + "';\n" +

			func() string { // 장르는 따로 구분하여 넣기 때문에 슬라이싱하여 '장르의 수 만큼 구문을 생성'한다.
				s := strings.Split(items.MgnrName, ",")
				gnr := ""
				for _, s := range s {
					s = strings.TrimSpace(s)
					gnr += "INSERT INTO MOV_GENRE(MOV_SEQ,MGNR_NAME) SELECT (SELECT LAST_INSERT_ID()),'" + s +
						"'\nFROM DUAL WHERE NOT EXISTS(SELECT MOV_SEQ, MGNR_NAME FROM MOV_GENRE\n" +
						"WHERE MOV_SEQ = (SELECT LAST_INSERT_ID()) AND MGNR_NAME='" + s + "');\n"
				}
				return gnr
			}() +

			func() string { // 상영 타입 또한 여러가지 타입이 동시에 존재할 수 있기에 익명함수를 통해 '타입갯수만큼 생성' 후 리턴
				s := strings.Split(items.MtypeName, ",")
				typ := ""
				for _, s := range s {
					s = strings.TrimSpace(s)
					typ += "INSERT INTO MOV_TYPE(MOV_SEQ,MTYPE_NAME) SELECT (SELECT LAST_INSERT_ID()),'" + s +
						"' FROM DUAL WHERE NOT EXISTS(SELECT MOV_SEQ, MTYPE_NAME FROM MOV_TYPE " +
						"WHERE MOV_SEQ = (SELECT LAST_INSERT_ID()) AND MTYPE_NAME='" + s + "');\n"
				}
				return typ
			}() +

			"INSERT INTO MOV_IMG(MOV_SEQ,MIMG_PATH)\n" +
			"SELECT (SELECT LAST_INSERT_ID()), '" + items.MimgName[1:] +
			"'\nFROM DUAL WHERE NOT EXISTS(SELECT MOV_SEQ, MIMG_PATH FROM MOV_IMG\n" +
			"WHERE MOV_SEQ = (SELECT LAST_INSERT_ID()) AND MIMG_PATH='" + items.MimgName[1:] + "');\n",
	)

	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.MovTitle)
	return
}
