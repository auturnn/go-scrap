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

	items.MovPst = "/img/poster/ing/" + items.MovTitle + ".jpg"
	rank := strconv.Itoa(items.MovRank)
	cnt := strconv.Itoa(items.MovCnt)
	_, err := m.SqlFile.WriteString("--예매율랭킹 " + rank + "위" + "\nINSERT INTO MOV_TEST(MOV_TITLE, MOV_CNT, MOV_OPD, MOV_STAT, MOV_TYPE, MOV_GNR, MOV_AGE, MOV_PST, MOV_RNT, MOV_PFM, MOV_SMR)\n" +
		"VALUES('" + items.MovTitle + "', " + cnt + ", '" + items.MovOpd + "', '" + items.MovStat + "', '" + items.MovType + "', '" + items.MovGnr + "', '" + items.MovAge + "', '" +
		items.MovPst + "', '" + items.MovRnt + "', '" + items.MovPfm + "', '" + items.MovSmr + "');\n",
	)
	if err != nil {
		log.Println("SQL Error! ", err)
	}

	log.Println("SQL Done : ", items.MovTitle)
	return
}
