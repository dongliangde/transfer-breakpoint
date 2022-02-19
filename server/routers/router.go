package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"transfer-breakpoint/apis"
	_ "transfer-breakpoint/docs"
	"transfer-breakpoint/middleware/cors"
)

/**
 * @Function: InitRouter
 * @Description: 初始化路由
 * @return *gin.Engine
 */
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(cors.HandlerCors())
	r.Use(gin.Recovery())
	gin.SetMode("debug")
	//无需认证的路由
	noCheckRoleRouter(r)
	return r
}

/**
 * @Function: noCheckRoleRouter
 * @Description:
 * @param r
 */
func noCheckRoleRouter(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fileUploadRouter := r.Group("fileUpload")
	fileUploadApi := apis.FileUploadApi{}
	{
		fileUploadRouter.POST("upload", fileUploadApi.UploadFile)                                 // 上传文件
		fileUploadRouter.POST("getFileList", fileUploadApi.GetFileList)                           // 获取上传文件列表
		fileUploadRouter.POST("deleteFile", fileUploadApi.DeleteFile)                             // 删除指定文件
		fileUploadRouter.POST("breakpointContinue", fileUploadApi.BreakpointContinue)             // 断点续传
		fileUploadRouter.GET("findFile", fileUploadApi.FindFile)                                  // 查询当前文件成功的切片
		fileUploadRouter.POST("breakpointContinueFinish", fileUploadApi.BreakpointContinueFinish) // 查询当前文件成功的切片
		fileUploadRouter.POST("removeChunk", fileUploadApi.RemoveChunk)                           // 查询当前文件成功的切片
	}
}
