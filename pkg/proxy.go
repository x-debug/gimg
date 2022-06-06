package pkg

import (
	"context"
	"gimg/config"
	lg "gimg/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const HTTP_SCHEMA = "http://"
const HTTPS_SCHEMA = "https://"

//Request is req info
type Request struct {
	Url      *url.URL
	OrignReq *http.Request
}

//HashVal return hash value of the request.url, IF two request is equal, their hash value also is equal
func (req *Request) HashVal() string {
	return CalcMd5Str(req.Url.String())
}

//NewRequest build request
func newRequest(orignRequest *http.Request, baseUri string) *Request {
	urlQuery := orignRequest.URL.Query()
	remoteUrl := urlQuery.Get("remote")
	if remoteUrl == "" {
		return nil
	}
	if !strings.Contains(remoteUrl, HTTP_SCHEMA) {
		if !strings.Contains(remoteUrl, HTTPS_SCHEMA) {
			remoteUrl = baseUri + remoteUrl
		}
	}
	url, err := url.Parse(remoteUrl)
	if err != nil {
		return nil
	}
	return &Request{OrignReq: orignRequest, Url: url}
}

//RemoteProxy is remote proxy for fetch images
type FileProxy struct {
	client    *http.Client
	basePath  string
	cachePath string
	timeout   time.Duration
	logger    lg.Logger
}

//Default setting
const clientMaxIdleConns = 30
const IdleConnTimeout = 60

//NewProxy build remote proxy
func NewProxy(cacheSavePath string, conf *config.ProxyConf, logger lg.Logger) *FileProxy {
	tr := &http.Transport{
		MaxIdleConns:    clientMaxIdleConns,
		IdleConnTimeout: 60 * time.Second,
	}

	proxy := &FileProxy{client: &http.Client{Transport: tr}, basePath: conf.BaseRemotePath, cachePath: cacheSavePath, logger: logger, timeout: time.Duration(conf.RequestTimeout) * time.Second}
	return proxy
}

//CloneRequest return cloned request, client put the parameter into http query's remote field, for example:
//https://you_domain_server?remote=https://chenxf.org/hello.jpg&w=100&h=100
//the remote url recommended that URLs not contain query strings
func (fp *FileProxy) CloneRequest(orignReq *http.Request) *Request {
	return newRequest(orignReq, fp.basePath)
}

//Do handle a client request, filehash is identity of file, IF the file already exists, the func return immediately.
//throw the TimeoutError IF client request timed out, OR throw the IOError IF func write the file.
func (fp *FileProxy) Do(req *Request, fileHash string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), fp.timeout)
	defer cancel()

	actualReq, err := http.NewRequestWithContext(ctx, "GET", req.Url.String(), nil)
	if err != nil {
		return err
	}
	fp.logger.Info("Query request for file proxy", lg.String("URL", req.Url.String()), lg.String("FileHash", fileHash))
	resp, err := fp.client.Do(actualReq)
	if err != nil {
		fp.logger.Error("Proxy request remote url error", lg.String("URL", req.Url.String()), lg.Error(err))
		return err
	}
	fileName := fp.cachePath + "/" + fileHash

	//func return IF the image file already exists,so the expired data maybe cached, but caller can delete the old images manually
	if _, err := os.Stat(fileName); err == nil {
		fp.logger.Info("Found image file is exist, so just return immediately", lg.String("FileName", fileName))
		return nil
	}

	fp.logger.Info("Create image file", lg.String("FileName", fileName))
	var file *os.File
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	defer file.Close()

	var rBytes []byte
	if rBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	file.Write(rBytes)
	return nil
}
