package clean

import (
	"fmt"
	"k8s.io/klog"
	"sort"

	"harbor-img-clear/harbor"
)

var (
	url         string
	user        string
	password    string
	projectName string
	keepNum     int
)

func deleteTagByID(harborClient harbor.Client, projectName string, keepNum int) (err error) {
	repoNames, err := harborClient.GetRepoNames(projectName)
	//klog.Infoln(repoNames)
	if err != nil {
		return
	}
	var size int64
	for _, repoName := range repoNames {
		tags, _ := harborClient.GetRepoTags(projectName, repoName)
		if len(tags) > keepNum { //tag大于keepNum才执行
			fmt.Printf("当前tag: %-4d，保留的tag: %-4d of %-40s ，开始执行删除\n", len(tags), keepNum, repoName)
			sort.Sort(tags)
			//klog.Infoln(tags)
			////自定义排序，根据tag的创建时间戳正序排列
			toDeleteTags := tags[keepNum-1:] //需要删除的tag切片
			for _, tag := range toDeleteTags {
				if len(tag.Tags) > 0 {
					fmt.Printf("     删除image: %s:%s, 创建时间为: %s\n", repoName, tag.Tags[0].Name, tag.Tags[0].PushTime)
					reference := tag.Digest
					err := harborClient.DeleteRepoTag(projectName, repoName, reference)
					if err != nil {
						fmt.Printf("image: %s:%s DeleteRepoTag: %s\n", repoName, tag.Tags[0].Name, err)
						continue
					}
					size += int64(tag.Size)
				}
			}
			fmt.Printf("repo: %s共清理: %.2f MB\n", repoName, float64(size)/1024/1024)
		} else {
			fmt.Printf("当前tag: %-4d，保留tag: %-4d of %-40s ,无需删除! \n", len(tags), keepNum, repoName)
		}

	}
	return
}

func Clean(url, user, password, projectName string, keepNum int) (err error) {
	harborClient := *harbor.NewClient(user, password, url)
	if projectName == "all" {
		allid, _ := harborClient.GetAllProjectID()
		//klog.Infoln(allid)
		for _, id := range allid {
			klog.Infoln(id)
			deleteTagByID(harborClient, projectName, keepNum)
		}
	} else {
		deleteTagByID(harborClient, projectName, keepNum)
	}
	return nil

}
