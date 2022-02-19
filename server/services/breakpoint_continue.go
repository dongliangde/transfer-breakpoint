package services

import (
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
	"transfer-breakpoint/common/request"
	"transfer-breakpoint/config"
	"transfer-breakpoint/model"
	"transfer-breakpoint/utils/upload"
)

type BreakpointContinueService struct{}

//@function: FindOrCreateFile
//@description: 上传文件时检测当前文件属性，如果没有文件则创建，有则返回文件的当前切片
//@param: fileMd5 string, fileName string, chunkTotal int
//@return: err error, file model.ExaFile
func (e *BreakpointContinueService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (err error, file model.File) {
	var cfile model.File
	cfile.FileMd5 = fileMd5
	cfile.FileName = fileName
	cfile.ChunkTotal = chunkTotal

	if errors.Is(config.GVA_DB.Where("file_md5 = ? AND is_finish = ?", fileMd5, true).First(&file).Error, gorm.ErrRecordNotFound) {
		err = config.GVA_DB.Where("file_md5 = ? AND file_name = ?", fileMd5, fileName).Preload("ExaFileChunk").FirstOrCreate(&file, cfile).Error
		return err, file
	}
	cfile.IsFinish = true
	cfile.FilePath = file.FilePath
	err = config.GVA_DB.Create(&cfile).Error
	return err, cfile
}

//@function: CreateFileChunk
//@description: 创建文件切片记录
//@param: id uint, fileChunkPath string, fileChunkNumber int
//@return: error
func (e *BreakpointContinueService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	var chunk model.FileChunk
	chunk.FileChunkPath = fileChunkPath
	chunk.ExaFileID = id
	chunk.FileChunkNumber = fileChunkNumber
	err := config.GVA_DB.Create(&chunk).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除文件切片记录
//@param: fileMd5 string, fileName string, filePath string
//@return: error
func (e *BreakpointContinueService) DeleteFileChunk(fileMd5 string, fileName string, filePath string) error {
	var chunks []model.FileChunk
	var file model.File
	err := config.GVA_DB.Where("file_md5 = ? ", fileMd5).First(&file).Update("IsFinish", true).Update("file_path", filePath).Error
	if err != nil {
		return err
	}
	err = config.GVA_DB.Where("exa_file_id = ?", file.ID).Delete(&chunks).Unscoped().Error
	return err
}

//@function: Upload
//@description: 创建文件上传记录
//@param: file model.FileUpload
//@return: error
func (e *BreakpointContinueService) Upload(file model.FileUpload) error {
	return config.GVA_DB.Create(&file).Error
}

//@function: FindFile
//@description: 删除文件切片记录
//@param: id uint
//@return: error, model.FileUpload
func (e *BreakpointContinueService) FindFile(id uint) (error, model.FileUpload) {
	var file model.FileUpload
	err := config.GVA_DB.Where("id = ?", id).First(&file).Error
	return err, file
}

//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.FileUpload
//@return: err error
func (e *BreakpointContinueService) DeleteFile(file model.FileUpload) (err error) {
	var fileFromDb model.FileUpload
	err, fileFromDb = e.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = config.GVA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64
func (e *BreakpointContinueService) GetFileRecordInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := config.GVA_DB.Model(&model.FileUpload{})
	var fileLists []model.FileUpload
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return err, fileLists, total
}

//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地
//@param: header *multipart.FileHeader, noSave string
//@return: err error, file model.FileUpload
func (e *BreakpointContinueService) UploadFile(header *multipart.FileHeader, noSave string) (err error, file model.FileUpload) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(err)
	}
	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		f := model.FileUpload{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s[len(s)-1],
			Key:  key,
		}
		return e.Upload(f), f
	}
	return
}
