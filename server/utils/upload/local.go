package upload

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
	"transfer-breakpoint/config"
	"transfer-breakpoint/utils"
)

type Local struct{}

//@object: *Local
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error
func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(config.LocalPath, os.ModePerm)
	if mkdirErr != nil {
		log.Println("function os.MkdirAll() Filed", mkdirErr.Error())
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := config.LocalPath + "/" + filename

	f, openError := file.Open() // 读取文件
	if openError != nil {
		log.Println("function file.Open() Filed", openError.Error())
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		log.Println("function os.Create() Filed", createErr.Error())
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		log.Println("function io.Copy() Filed", copyErr.Error())
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}

//@object: *Local
//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error
func (*Local) DeleteFile(key string) error {
	p := config.LocalPath + "/" + key
	if strings.Contains(p, config.LocalPath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
