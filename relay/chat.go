package relay

import (
	"bytes"
	"io"
	"net/http"
	"newapi-mini/model"
	"strings"

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
	modelName := reqBody["model"].(string)
	// 查找可用渠道
	var channel model.Channel
	model.DB.Where("status = ?", 1).First(&channel)
	// 判断渠道是否有 modelName
	contains := strings.Contains(channel.Models, modelName)
	if !contains {
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
	defer resp.Body.Close()
	// 获取响应内容
	ctx.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
}
