package stock_data

import (
	"astock/utils"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type StockData struct {
	TsCode   string
	ListDate string
}

func GetAllStockData(requestUrl string, requestToken string)  {
	response := utils.GetAllAStock(requestUrl, requestToken)

	// 连接数据库
	db, err := sql.Open("mysql", "root:931102@tcp(127.0.0.1:3306)/astock")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//不使用协程插入
	for _, items := range response.Data.Items {
		_, err = db.Exec("INSERT INTO astock_daily (ts_code, trade_date, open, high, low, close, pre_close, change, pct_chg, vol, amount,created_on, modified_on) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], int(time.Now().Unix()), int(time.Now().Unix()))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func GetStockDailyData(requestUrl string, requestToken string)  {
	// 连接数据库
	db, err := sql.Open("mysql", "root:931102@tcp(127.0.0.1:3306)/astock")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 取出所有URL记录
	rows, err := db.Query("SELECT ts_code, list_date FROM astock order by id asc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	StockDatas := make([]StockData, 0)
	for rows.Next() {
		stockData := StockData{}
		err := rows.Scan(&stockData.TsCode, &stockData.ListDate)
		if err != nil {
			fmt.Println("Failed to read row:", err)
			return
		}

		StockDatas = append(StockDatas, stockData)
	}

	currentTime := time.Now()
	dateString := currentTime.Format("20060102")
	fmt.Println("Current date in YYYYMMDD format:", dateString)

	// 使用10个协程进行请求并保存记录
	var wg sync.WaitGroup
	ch := make(chan StockData)

	const maxCount = 450
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			for stockdata := range ch {

				count++
				if count%maxCount == 0 {
					time.Sleep(1 * time.Minute)
				}

				response := utils.GetDailAStock(requestUrl, requestToken, stockdata.TsCode,stockdata.ListDate,dateString)
				for _, items := range response.Data.Items {
					_, err = db.Exec("INSERT INTO astock_daily (ts_code, trade_date, open, high, low, close, pre_close, `change`, pct_chg, vol, amount,created_on, modified_on) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], int(time.Now().Unix()), int(time.Now().Unix()))
					if err != nil {
						fmt.Println(err)
						continue
					}
				}

			}
			wg.Done()
		}()
	}

	wg.Add(10)
	for _, stockdata := range StockDatas {
		ch <- stockdata
	}
	close(ch)

	wg.Wait()
}
