package main

import (
	"astock/stock_data"
	"astock/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
	//"net/http"
	//"time"
)


func main() {
	stock_data.GetAllStockData()
	stock_data.GetStockDailyData()
}

