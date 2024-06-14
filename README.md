## harbor rest api v2.x img-clear



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
    ➜  harbor-clean git:(master) ✗ ./harbor-img-clear
I0614 18:27:48.086809   23146 harbor.go:121] http://10.200.82.51/api/v2.0/projects/isap/repositories/deepglint%252Fisap/artifacts?page=1&page_size=100
当前tag: 29  ，保留tag: 30   of isap/deepglint/isap                      ,无需删除!
```

### 3. 在harbor后台页面立即 清理或者设置定时任务 清理
    说明： 每个project 其实可以自己设置策略；比如：保留最近推送的40个 artifacts基于条件tags匹配**基于条件 无 Tag
    但是项目多了以后麻烦；所以可以设置all统一脚本跑！