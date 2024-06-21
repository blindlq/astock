package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestBody struct {
	APIName string `json:"api_name"`
	Token   string `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields string `json:"fields"`
}

type ResponseData struct {
	Fields []string    `json:"fields"`
	Items  [][]interface{}  `json:"items"`
}

type Response struct {
	Data ResponseData `json:"data"`
}

//股票基础信息
func GetAllAStock(requestUrl string, requestToken string) Response {
	//请求的数据-股票基础信息-stock_basic
	//构建请求体数据
	//构建动态的params参数
	params := make(map[string]interface{})
	params["list_status"] = "L"
	//请求字段
	fields := "ts_code,symbol,name,area,industry,list_date"

	// 构建请求体数据
	requestBody := RequestBody{
		APIName: "stock_basic",
		Token:   requestToken,
		Params:  params,
		Fields:  fields,
	}

	return sendRequest(requestUrl, requestBody)
}

func GetDailAStock(requestUrl string, requestToken string, tsCode string, startDate string, endDate string) Response {
	params := make(map[string]interface{})
	params["start_date"] = startDate
	params["end_date"] = endDate
	params["ts_code"] = tsCode
	//请求字段
	fields := "ts_code,trade_date,open,high,low,close,pre_close,change,pct_chg,vol,amount"

	// 构建请求体数据
	requestBody := RequestBody{
		APIName: "daily",
		Token:   requestToken,
		Params:  params,
		Fields:  fields,
	}

	return sendRequest(requestUrl, requestBody)
}

func sendRequest(requestUrl string, requestBody RequestBody) Response {
	// 将请求体序列化为JSON字节流
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Failed to marshal request body:", err)
		return Response{}
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println("创建HTTP请求失败", err)
		return Response{}
	}

	//设置json的请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP请求失败:", err)
		return Response{}
	}
	defer resp.Body.Close()

	// 读取响应数据
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应数据失败:", err)
		return Response{}
	}

	// 解析响应数据
	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Failed to unmarshal response body:", err)
		return Response{}
	}

	return response
}