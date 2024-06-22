package main

import (
	"astock/global"
	"astock/pkg/setting"
	"astock/scripts"
	"flag"
	"log"
	"strings"
)

var (
	config string
)

func init()  {
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}

	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func main() {
	//fmt.Println(global.ScriptsSetting.DailyHotUrl)
	//scripts.GetAllStockData(viper.GetString("RequestConfig.Url"),viper.GetString("RequestConfig.Token"))
	//scripts.GetStockDailyData(viper.GetString("RequestConfig.Url"),viper.GetString("RequestConfig.Token"))
	scripts.GetStockDailyHotData()
}

func setupFlag() error {
	flag.StringVar(&config,"config","configs/","指定配置文件路径")
	flag.Parse()

	return nil
}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(config,",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("RequestConfig", &global.ScriptsSetting)
	return nil
}

