package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
	"strconv"
	"transfer-breakpoint/common/request"
	"transfer-breakpoint/common/response"
	"transfer-breakpoint/model"
	"transfer-breakpoint/services"
	"transfer-breakpoint/utils"
)

type FileUploadApi struct{}

var ServiceGroupApp = new(services.BreakpointContinueService)

// @Tags fileUpload
// @Summary 断点续传到服务器
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "an example for breakpoint resume, 断点续传示例"
// @Success 200 {object} response.Response{msg=string} "断点续传到服务器"
// @Router /fileUpload/breakpointContinue [post]
func (u *FileUploadApi) BreakpointContinue(c *gin.Context) {
	fileMd5 := c.Request.FormValue("fileMd5")
	fileName := c.Request.FormValue("fileName")
	chunkMd5 := c.Request.FormValue("chunkMd5")
	chunkNumber, _ := strconv.Atoi(c.Request.FormValue("chunkNumber"))
	chunkTotal, _ := strconv.Atoi(c.Request.FormValue("chunkTotal"))
	_, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("接收文件失败!", err)
		response.FailWithMessage("接收文件失败", c)
		return
	}
	f, err := FileHeader.Open()
	if err != nil {
		log.Println("文件读取失败!", err)
		response.FailWithMessage("文件读取失败", c)
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cen, _ := ioutil.ReadAll(f)
	if !utils.CheckMd5(cen, chunkMd5) {
		log.Println("检查md5失败!", err)
		response.FailWithMessage("检查md5失败", c)
		return
	}
	err, file := ServiceGroupApp.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		log.Println("查找或创建记录失败!", err)
		response.FailWithMessage("查找或创建记录失败", c)
		return
	}
	err, pathc := utils.BreakPointContinue(cen, fileName, chunkNumber, chunkTotal, fileMd5)
	if err != nil {
		log.Println("断点续传失败!", err)
		response.FailWithMessage("断点续传失败", c)
		return
	}

	if err = ServiceGroupApp.CreateFileChunk(file.ID, pathc, chunkNumber); err != nil {
		log.Println("创建文件记录失败!", err)
		response.FailWithMessage("创建文件记录失败", c)
		return
	}
	response.OkWithMessage("切片创建成功", c)
}

// @Tags fileUpload
// @Summary 查找文件
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "Find the file, 查找文件"
// @Success 200 {object} response.Response{data=response.FileResponse,msg=string} "查找文件,返回包括文件详情"
// @Router /fileUpload/findFile [post]
func (u *FileUploadApi) FindFile(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	chunkTotal, _ := strconv.Atoi(c.Query("chunkTotal"))
	err, file := ServiceGroupApp.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		log.Println("查找失败!", err)
		response.FailWithMessage("查找失败", c)
	} else {
		response.OkWithDetailed(response.FileResponse{File: file}, "查找成功", c)
	}
}

// @Tags fileUpload
// @Summary 创建文件
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件完成"
// @Success 200 {object} response.Response{data=response.FilePathResponse,msg=string} "创建文件,返回包括文件路径"
// @Router /fileUpload/findFile [post]
func (b *FileUploadApi) BreakpointContinueFinish(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	err, filePath := utils.MakeFile(fileName, fileMd5)
	if err != nil {
		log.Println("文件创建失败!", err)
		response.FailWithDetailed(response.FilePathResponse{FilePath: filePath}, "文件创建失败", c)
	} else {
		response.OkWithDetailed(response.FilePathResponse{FilePath: filePath}, "文件创建成功", c)
	}
}

// @Tags fileUpload
// @Summary 删除切片
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "删除缓存切片"
// @Success 200 {object} response.Response{msg=string} "删除切片"
// @Router /fileUpload/removeChunk [post]
func (u *FileUploadApi) RemoveChunk(c *gin.Context) {
	var file model.File
	c.ShouldBindJSON(&file)
	err := utils.RemoveChunk(file.FileMd5)
	if err != nil {
		log.Println("缓存切片删除失败!", err)
		return
	}
	err = ServiceGroupApp.DeleteFileChunk(file.FileMd5, file.FileName, file.FilePath)
	if err != nil {
		log.Println(err.Error(), err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("缓存切片删除成功", c)
	}
}

// @Tags fileUpload
// @Summary 上传文件示例
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件示例"
// @Success 200 {object} response.Response{data=response.FileUploadResponse,msg=string} "上传文件示例,返回包括文件详情"
// @Router /fileUpload/upload [post]
func (u *FileUploadApi) UploadFile(c *gin.Context) {
	var file model.FileUpload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("接收文件失败!", err)
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err, file = ServiceGroupApp.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		log.Println("修改数据库链接失败!", err)
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(response.FileUploadResponse{File: file}, "上传成功", c)
}

// @Tags fileUpload
// @Summary 删除文件
// @Produce  application/json
// @Param data body model.FileUpload true "传入文件里面id即可"
// @Success 200 {object} response.Response{msg=string} "删除文件"
// @Router /fileUpload/deleteFile [post]
func (u *FileUploadApi) DeleteFile(c *gin.Context) {
	var file model.FileUpload
	_ = c.ShouldBindJSON(&file)
	if err := ServiceGroupApp.DeleteFile(file); err != nil {
		log.Println("删除失败!", err)
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// @Tags fileUpload
// @Summary 分页文件列表
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页文件列表,返回包括列表,总数,页码,每页数量"
// @Router /fileUpload/getFileList [post]
func (u *FileUploadApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	err, list, total := ServiceGroupApp.GetFileRecordInfoList(pageInfo)
	if err != nil {
		log.Println("获取失败!", err)
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
