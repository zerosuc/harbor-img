## harbor rest api v2.x img-clear



### 0. 前置
```
➜  harbor-clean git:(master) ✗ make build
go build -ldflags "-X 'harbor-img/version.GIT_TAG=v0.1' -X 'harbor-img/version.GIT_COMMIT=8d3157dbbb3dde6121aeb47437aa579e335c5852' -X 'harbor-img/version.GIT_BRANCH=master' -X 'harbor-img/version.BUILD_TIME=2024-06-17T12:00:07+0800' -X 'harbor-img/version.GO_VERSION=go version go1.21.4 darwin/amd64'" -o harbor-img main.go
➜  harbor-clean git:(master) ✗ ./harbor-img -v

Version   : v0.1
Build Time: 2024-06-17T12:00:07+0800
Git Branch: master
Git Commit: 8d3157dbbb3dde6121aeb47437aa579e335c5852
Go Version: go version go1.21.4 darwin/amd64
➜  harbor-clean git:(master) ✗ ./harbor-img -h

harbor-img 用于清理harbor的仓库中的tag,以释放存储资源

Usage:
  harbor-img [flags]
  harbor-img [command]

Examples:
./harbor-img clear --address http://10.200.82.51  --user admin --password Harbor12345 --project appsvc  --keepNum 30

Available Commands:
  clear       clear 仓库镜像清理
  help        Help about any command

Flags:
  -h, --help      help for harbor-img
  -v, --version    harbor-img 当前版本

Use "harbor-img [command] --help" for more information about a command.
```

### 1. 配置 config.yaml文件
```shell
harbor:
  url: "http://10.200.82.51"
  username: "admin"
  password: "Harbor12345"
  project: "isap" 
  num: 30 # 需要保留最新的tag数
```
### 2. 执行
```
    
```

### 3. 在harbor后台页面立即 清理或者设置定时任务 清理
    说明： 每个project 其实可以自己设置策略；比如：保留最近推送的40个 artifacts基于条件tags匹配**基于条件 无 Tag
    但是项目多了以后麻烦；所以可以设置all统一脚本跑！