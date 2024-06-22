package scripts

//抓取同花顺每日前100名热榜

import (
	"astock/global"
	utils2 "astock/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type StockHotData struct {
	Symbol string
}

type HotData struct {
	RankField []string    `json:"rank_field"`
	RankList  [][]interface{}  `json:"rank_list"`
}

type Response struct {
	Data HotData `json:"data"`
}

func GetStockDailyHotData() {
	// 连接数据库
	db, err := sql.Open("mysql", "root:931102@tcp(127.0.0.1:3306)/astock")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 取出所有URL记录
	rows, err := db.Query("SELECT symbol FROM astock order by id asc limit 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	stockHotData := StockHotData{}
	for rows.Next() {
		err := rows.Scan(&stockHotData.Symbol)
		if err != nil {
			fmt.Println("Failed to read row:", err)
		}
	}

	spiderUrl := global.ScriptsSetting.DailyHotUrl + stockHotData.Symbol + "&data_type=rank"
	//fmt.Println(spiderUrl)
	body := utils2.SendGetRequest(spiderUrl)

	// 解析响应数据
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Failed to unmarshal response body:", err)
	}

	// 提取前100名股票数据
	// 获取当前时间
	currentTime := time.Now()
	dateString := currentTime.Format("20060102")

	i := 0
	count := 100
	for _, stock := range response.Data.RankList {
		_, err = db.Exec("INSERT INTO astock_daily_hot_rank (symbol, trade_date, stock_name, rank, `change`, pct_chg, close, marketvalue, created_on, modified_on) VALUES (?,?,?,?,?,?,?,?,?,?)", stock[0], dateString, stock[1], stock[2], stock[3], stock[4], stock[5], stock[6], int(time.Now().Unix()), int(time.Now().Unix()))
		if err != nil {
			fmt.Println(err)
			continue
		}

		i++
		// 控制数量，当达到指定数量时停止提取
		if i >= count {
			break
		}
	}
}




