package model

import (
	"gorm.io/gorm"
	"time"
)

// file struct, 文件结构体
type File struct {
	ID         uint           `gorm:"primarykey"` // 主键ID
	FileName   string
	FileMd5    string
	FilePath   string
	FileChunk  []FileChunk
	ChunkTotal int
	IsFinish   bool
	CreateTime time.Time      // 创建时间
	UpdateTime time.Time      // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

// file chunk struct, 切片结构体
type FileChunk struct {
	ID              uint           `gorm:"primarykey"` // 主键ID
	ExaFileID       uint
	FileChunkNumber int
	FileChunkPath   string
	CreateTime      time.Time      // 创建时间
	UpdateTime      time.Time      // 更新时间
	DeleteTime      gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type FileUpload struct {
	ID         uint           `gorm:"primarykey"`              // 主键ID
	Name       string         `json:"name" gorm:"comment:文件名"` // 文件名
	Url        string         `json:"url" gorm:"comment:文件地址"` // 文件地址
	Tag        string         `json:"tag" gorm:"comment:文件标签"` // 文件标签
	Key        string         `json:"key" gorm:"comment:编号"`   // 编号
	CreateTime time.Time      // 创建时间
	UpdateTime time.Time      // 更新时间
	DeleteTime gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
