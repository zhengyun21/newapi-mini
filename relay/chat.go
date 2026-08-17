package relay

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"newapi-mini/model"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func RelayChat(ctx *gin.Context) {
	// 获取请求体内容
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": "读取/解析请求体失败"})
		return
	}
	var reqBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &reqBody)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": "请求出现错误"})
		return
	}
	// 拿到用户请求的模型
	modelName, ok := reqBody["model"].(string)
	if !ok {
		ctx.JSON(400, gin.H{"msg": "请输入正确的模型格式"})
		return
	}
	// 查找可用渠道
	var channel model.Channel
	result := model.DB.Where("status = ? AND FIND_IN_SET(?,models)", 1, modelName).First(&channel)
	if result.Error != nil {
		ctx.JSON(404, gin.H{"msg": "该渠道下暂无该模型"})
		return
	}
	// 发送请求给上游
	req, _ := http.NewRequest("POST", channel.BaseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channel.APIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(502, gin.H{"msg": "上游请求失败"})
		return
	}
	defer resp.Body.Close() // 若err != nil，resp为空，也不会泄露。
	if resp.StatusCode >= 400 {
		// 读取上游错误响应体
		errBody, _ := io.ReadAll(resp.Body)
		// 记录日志，方便排查
		// 日志内容：状态码、上游错误内容
		log.Printf("上游错误：响应码=%v 响应体=%v", resp.StatusCode, string(errBody))
		// 读完 body 后，原样返回给客户端
		ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), errBody)
		return
	}
	// 获取响应内容
	ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)

}
