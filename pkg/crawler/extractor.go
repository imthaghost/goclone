package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Extractor(baseUrl, link, projectFilePath string, referer string, userAgent string, cookieJar *cookiejar.Jar) {
	if link == "" || strings.HasPrefix(link, "#") {
		return
	}
	fmt.Println("Extracting --> ", link)
	client := &http.Client{Jar: cookieJar}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	if referer != "" {
		req.Header.Set("Referer", referer) // 例如设置 Authorization 头
	}
	// get the html body
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// Closure
	defer resp.Body.Close()
	baseURL, err := ExtractBaseURL(baseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, err := ExtractPath(baseURL, link)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	writeFileToPath(projectFilePath+"/"+path, resp)
}

func ExtractBaseURL(baseURL string) (string, error) {
	// 解析基础URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("解析基础URL失败: %v", err)
	}
	// 获取路径部分，并去掉文件名部分
	basePath := parsedURL.Path
	// 处理路径部分，去掉文件名
	dir := path.Dir(basePath)
	// 如果路径最后没有斜杠，手动添加斜杠
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	// 构建新的基础 URL
	baseURLWithoutFile := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, dir)

	return baseURLWithoutFile, nil
}

func ExtractPath(baseURL, resourceURL string) (string, error) {
	// 解析基础URL
	baseParsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("解析基础URL失败: %v", err)
	}

	// 解析资源URL
	resourceParsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", fmt.Errorf("解析资源URL失败: %v", err)
	}

	// 获取基础URL的路径部分
	basePath := baseParsed.Path
	if strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1] // 移除末尾的斜杠
	}

	// 获取资源URL的路径部分
	resourcePath := resourceParsed.Path

	// 移除基础路径的前缀部分，获取相对路径
	if strings.HasPrefix(resourcePath, basePath) {
		relativePath := strings.TrimPrefix(resourcePath, basePath+"/")
		return relativePath, nil
	}

	return "", fmt.Errorf("资源URL不以基础URL为前缀")
}
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func writeFileToPath(projectFilePath string, resp *http.Response) {
	dir := filepath.Dir(projectFilePath)
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.OpenFile(projectFilePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		println(projectFilePath)
		panic(err)
	}
	defer f.Close()
	htmlData, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	f.Write(htmlData)
}
