package middleware

import (
	"newapi-mini/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			/*
				ctx.Abort() 只是设置了
				c.index = math.MaxInt8，
				让后续 handler 不再执行，
				但当前函数本身会继续往下走。
			*/
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "Token格式有误"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "Token格式有误"})
			return
		}
		var tokenRecord model.Token
		// key 要反义，因为是关键字
		result := model.DB.Where("`key` = ? and status = ?", token, 1).First(&tokenRecord)
		err := result.Error
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "Token不存在"})
			return
		}
		ctx.Set("token_record", tokenRecord)
		ctx.Next()
	}
}
