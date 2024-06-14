package harbor

import (
	"encoding/json"
	"errors"
	"fmt"
	"harbor-img-clear/model"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// 最大一页100
var pageSize = 100

// Client
type Client struct {
	Client  *http.Client //http.Clinet类型，结构体嵌套
	BaseURL string
}

// NewClient
func NewClient(username, password, baseURL string) *Client {
	client := &http.Client{ //定义client变量，是一个结构体
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { //匿名函数
				req.SetBasicAuth(username, password)
				return nil, nil
			},
		},
	}
	return &Client{
		Client:  client,
		BaseURL: baseURL,
	}
}
func (c *Client) GetAllProjectID() (allid []model.Project, err error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/v2.0/projects?page=1&page_size=100")
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code is:%v", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var allProjects []model.Project

	err = json.Unmarshal(body, &allProjects) //解析json
	if err != nil {
		return
	}
	for _, a := range allProjects {
		allid = append(allid, a)
	}

	return allid, nil
}

func (c *Client) GetRepoNames(projectName string) (repoNames []string, err error) {
	url := fmt.Sprintf(c.BaseURL+"/api/v2.0/projects/%s/repositories", projectName)
	resp, err := c.Client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var repos []model.Repo
	err = json.Unmarshal(body, &repos) //json转结构体
	if err != nil {
		return
	}
	for _, repo := range repos { //结构体转切片
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

// Client 结构体定义
type HClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
}

func (c *Client) GetRepoTags(projectName, repo string) (tags model.Tags, err error) {

	// 去掉仓库名称中的项目名称和斜杠
	trimmedRepoName := strings.TrimPrefix(repo, projectName+"/")
	// 对处理后的仓库名称进行URL编码
	encodedRepoName := url.PathEscape(trimmedRepoName)
	// 对路径中的%符号进行二次编码
	doubleEncodedRepoName := strings.ReplaceAll(encodedRepoName, "%", "%25")

	client := &HClient{
		BaseURL:  c.BaseURL,
		Username: "admin",
		Password: "Harbor12345",
		Client:   &http.Client{},
	}
	// 获取 Token
	token, err := client.getToken()
	if err != nil {
		return
	}

	// 创建请求
	url := fmt.Sprintf(c.BaseURL+"/api/v2.0/projects/%s/repositories/%s/artifacts?page=1&page_size=%d", projectName, doubleEncodedRepoName, pageSize)
	klog.Infoln(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//klog.Infoln("response body: ", string(body))
	// 解析 JSON 响应
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return
	}
	//klog.Infoln(tags)

	// 排序标签
	sort.Sort(tags)
	if !sort.IsSorted(tags) {
		return nil, errors.New("tags not sorted")
	}
	return
}

// 获取 Token 的函数
func (c *HClient) getToken() (string, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/service/token?service=harbor-registry", nil)
	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResponse struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &tokenResponse)

	return tokenResponse.Token, nil
}

// DeleteRepoTag 在harbor 后台页面抓包可以得到
// http://10.200.82.51/api/v2.0/projects/appsvc/repositories/vsap%252Ffront/artifacts/sha256:4d101e7450831885d1319f6687036e3529a9086785535bd547be84c50bfbb167
func (c *Client) DeleteRepoTag(projectName, repo, reference string) (err error) {
	requestURL := ""
	if strings.Contains(repo, "/") {
		// 去掉仓库名称中的项目名称和斜杠
		trimmedRepoName := strings.TrimPrefix(repo, projectName+"/")
		// 对处理后的仓库名称进行URL编码
		encodedRepoName := url.PathEscape(trimmedRepoName)
		// 对路径中的%符号进行二次编码
		doubleEncodedRepoName := strings.ReplaceAll(encodedRepoName, "%", "%25")
		klog.Infoln(doubleEncodedRepoName)
		requestURL = fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s", c.BaseURL, projectName, doubleEncodedRepoName, reference)
	} else {
		requestURL = fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s", c.BaseURL, projectName, repo, reference)
	}
	// 构建请求URL
	request, err := http.NewRequest("DELETE", requestURL, nil)
	klog.Infoln(requestURL)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	client := &HClient{
		BaseURL:  c.BaseURL,
		Username: "admin",
		Password: "Harbor12345",
		Client:   &http.Client{},
	}
	// 获取 Token
	token, err := client.getToken()
	if err != nil {
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)
	// 设置请求头
	request.Header.Set("accept", "application/json")
	//request.Header.Set("authorization", "Bearer "+c.Token)

	// 执行请求
	resp, err := c.Client.Do(request)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code: %v", resp.StatusCode)
	}

	return nil
}
