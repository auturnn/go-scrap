package run

import (
	"log"
	"strconv"
	"strings"
)

func (m *Megabox) CreateSQL(items *MovieList) {
	if items.MovAge == "전체관람가" {
		items.MovAge = "all"
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
	_, err := m.SqlFile.WriteString("\n\nINSERT INTO MOV_MST(MOV_TITLE, MOV_RANK, MOV_CNT, MOV_OPD, MOV_STAT, MOV_AGE)\n" +
		"VALUES('" + items.MovTitle + "', " + rank + ", " + cnt + ", '" + items.MovOpd + "', '" + items.MovStat + "', '" + items.MovAge + "');\n" +
		"INSERT INTO MOV_DT " +
		"VALUES((SELECT LAST_INSERT_ID()), '" +
		items.MovDrt + "', '" + items.MovAct + "', '" + items.MovLen + "', '" + items.MovSmr + "');\n" +
		"INSERT INTO MOV_GENRE " +
		"VALUES((SELECT LAST_INSERT_ID()), '" + items.MgnrName + "');\n" +
		"INSERT INTO MOV_TYPE " +
		"VALUES((SELECT LAST_INSERT_ID()), '" + items.MtypeName + "');\n" +
		"INSERT INTO MOV_IMG " +
		"VALUES((SELECT LAST_INSERT_ID()), '" + items.MimgName[1:] + "');\n")
	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.MovTitle)
	return
}
