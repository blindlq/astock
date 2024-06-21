package main

import (
	"astock/stock_data"
	//"astock/stock_data"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func init()  {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs/")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	stock_data.GetAllStockData(viper.GetString("RequestConfig.Url"),viper.GetString("RequestConfig.Token"))
	stock_data.GetStockDailyData(viper.GetString("RequestConfig.Url"),viper.GetString("RequestConfig.Token"))
}

