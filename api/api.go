package api

import (
	"errors"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"tail-project/rest"
)

// data path
var (
	maxUploadSize   int64 = 3 * 1024 * 1024
	DataPath              = "data"
	HostNameAndPort       = ""
	api                   = new(ProjectServerApi)
)

func renderError(response *restful.Response, message string, statusCode int) {
	rest.WriteEntity(nil, errors.New(message), response)
}

// upload project api
type ProjectServerApi struct {
}

func (*ProjectServerApi) upload(request *restful.Request, response *restful.Response) {
	suffix := request.PathParameter("suffix")
	r := request.Request
	w := response.ResponseWriter
	// validate file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(response, "文件太大了", http.StatusBadRequest)
		return
	}

	// parse and validate file and post parameters
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		renderError(response, "索引文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(response, "读取文件失败", http.StatusBadRequest)
		return
	}

	newPath := GetFilePath(suffix)

	i := strings.LastIndex(newPath, "/")
	if i ==-1{
		i = strings.LastIndex(newPath, "\\")
	}
	dir := newPath[:i]
	_, err = os.Lstat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		renderError(response, "目录创建失败", http.StatusInternalServerError)
		return
	}
	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(response, "创建文件失败", http.StatusInternalServerError)
		return
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(response, "写入文件失败", http.StatusInternalServerError)
		return
	}
	//去掉 data path描述
	path :=strings.ReplaceAll(newPath,DataPath,"")
	if HostNameAndPort != "" {
		path = "http://" + HostNameAndPort +path
	}
	rest.WriteEntity(path, nil, response)
}

func init() {
	binder, webService := rest.NewJsonWebServiceBinder("/project")
	webService.Route(webService.POST("upload/{suffix}").To(api.upload))
	binder.BindAdd()
}
