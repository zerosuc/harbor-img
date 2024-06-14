package model

import (
	"time"
)

// Project /api/projects?name=cloud
type Project struct {
	Name string `json:"name"`
	ID   int    `json:"project_id"`
}

// Repo /api/repositories?project_id=2
type Repo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Tag /api/repositories/cloud/demojava/tags  test-bcd2e9d
//type Tag struct {
//	Size    int64     `json:"size"`
//	Name    string    `json:"name"`
//	Created time.Time `json:"created"`
//}

type Tag struct {
	Accessories       interface{} `json:"accessories"`
	Digest            string      `json:"digest"`
	Icon              string      `json:"icon"`
	ID                int         `json:"id"`
	Labels            interface{} `json:"labels"`
	ManifestMediaType string      `json:"manifest_media_type"`
	MediaType         string      `json:"media_type"`
	ProjectID         int         `json:"project_id"`
	PullTime          time.Time   `json:"pull_time"`
	PushTime          time.Time   `json:"push_time"`
	References        []struct {
		ChildDigest string `json:"child_digest"`
		ChildID     int    `json:"child_id"`
		ParentID    int    `json:"parent_id"`
		Platform    struct {
			OsFeatures   interface{} `json:"OsFeatures"`
			Architecture string      `json:"architecture"`
			Os           string      `json:"os"`
		} `json:"platform"`
		Urls        interface{} `json:"urls"`
		Annotations struct {
			VndDockerReferenceDigest string `json:"vnd.docker.reference.digest"`
			VndDockerReferenceType   string `json:"vnd.docker.reference.type"`
		} `json:"annotations,omitempty"`
	} `json:"references"`
	RepositoryID int `json:"repository_id"`
	Size         int `json:"size"`
	Tags         []struct {
		ArtifactID   int       `json:"artifact_id"`
		ID           int       `json:"id"`
		Immutable    bool      `json:"immutable"`
		Name         string    `json:"name"`
		PullTime     time.Time `json:"pull_time"`
		PushTime     time.Time `json:"push_time"`
		RepositoryID int       `json:"repository_id"`
		Signed       bool      `json:"signed"`
	} `json:"tags"`
	Type string `json:"type"`
}

// Tags implement the sort interface
type Tags []Tag //结构体类型的切片

// 以下三个函数用于自定义排序，自定义排序必须实现这3个方法。表示根据创建时间
func (t Tags) Len() int           { return len(t) }
func (t Tags) Less(i, j int) bool { return t[i].PullTime.After(t[j].PullTime) } //时间比较，i在j之后，则交互i，j。可理解为按照时间戳大小正序排序
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
