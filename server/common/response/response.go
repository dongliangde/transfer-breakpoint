package response

import (
	"net/http"
	"transfer-breakpoint/model"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, data, message, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, data, message, c)
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type FilePathResponse struct {
	FilePath string `json:"filePath"`
}

type FileUploadResponse struct {
	File model.FileUpload `json:"file"`
}

type FileResponse struct {
	File model.File `json:"file"`
}
