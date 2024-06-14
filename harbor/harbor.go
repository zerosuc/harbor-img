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

// getProjectID . 根据projectName 获取projectID

func (c *Client) GetProjectID(projectName string) (projectID int, err error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/v2.0/projects?name=" + projectName)
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

	var projects []model.Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return
	}
	for _, p := range projects { //返回的是模糊查询的结果，所以需要做个判断
		if p.Name == projectName {
			return p.ID, nil
		}
	}

	return 0, errors.New("not found")
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

// func (c *Client) GetRepo(projectId int) (repoNames []string, err error)

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
	url := fmt.Sprintf(c.BaseURL+"/api/v2.0/projects/%s/repositories/%s/artifacts", projectName, doubleEncodedRepoName)
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

	klog.Infoln("response body: ", string(body))
	// 解析 JSON 响应
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return
	}

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

//// 获取tag列表
//func (c *Client) GetRepoTags(repo string) (tags model.Tags, err error) {
//	resp, err := c.Client.Get(c.BaseURL + "/api/v2.0/repositories/" + repo + "/tags")
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return
//	}
//
//	err = json.Unmarshal(body, &tags) //json转结构体
//	if err != nil {
//		return
//	}
//	// 排序
//	sort.Sort(tags)
//	if !sort.IsSorted(tags) {
//		return nil, errors.New("tags not sorted")
//	}
//	return
//}

// DeleteRepoTag delete tags with repo name and tag.
func (c *Client) DeleteRepoTag(repo string, tag string) (err error) {
	request, err := http.NewRequest("DELETE", c.BaseURL+"/api/v2.0/repositories/"+repo+"/tags/"+tag, nil)
	if err != nil {
		return
	}
	resp, err := c.Client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("resp code=%v", resp.StatusCode)
		return
	}
	return
}
