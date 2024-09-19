package cluster

import (
	"log"

	v1Capture "github.com/KubeOperator/kubepi/internal/model/v1/capture"
	v1SvcCapture "github.com/KubeOperator/kubepi/internal/service/v1/capture"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
)

func prePareDB() {
	// var capturelist []v1Capture.Capture

	db := v1SvcCapture.NewService().GetDB(common.DBOptions{})
	var capturelist []v1Capture.Capture
	if err := db.All(&capturelist); err != nil {
		return
	}

	if len(capturelist) == 0 {
		err := v1SvcCapture.NewService().PostConfig(&v1Capture.Capture{
			Id: "default",
			Config: v1Capture.Config{
				DefaultImg: "swr.cn-north-4.myhuaweicloud.com/hfbbg4/proxy-prod:v0.1",
			},
		}, common.DBOptions{})
		if err != nil {
			log.Println("111dsfd: %v", err)

		}
	}

}

func testAddTaskRunBak() {
	// yamlFile := []byte(fmt.Sprintf(varTestTask, "1724392364", "running"))

	/*
		// 解析YAML数据到map
		var data capture.TaskYaml
		fmt.Println(fmt.Sprintf(varTestTask, "1724392364", "running", time.Now().Unix(), 60, "192.168.113.249", "192.168.113.249"))
		if err := yaml.Unmarshal([]byte(fmt.Sprintf(varTestTask, "1724392364", "running", time.Now().Unix(), 60, "192.168.113.249", "192.168.113.249")), &data); err != nil {
			log.Println("error: %v", err)
		}

		// 将map转换为JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("error: %v", err)
		}

		// 打印JSON数据
		fmt.Println(string(jsonData))
	*/

}

var varTestTask = `


`

// stdbuf -o0 -e0
var varTestTaskNoAdd = `
task:
  id: "" #"1724392364"
  taskStatus: running #running
    #状态 running
    #状态 stop
    #状态 finish 文件已经拷贝集中到管理pod中
    #状态 stop-type2 到时间强制终止
  startTime: 
  maxTimeDuring: 60000000 # 60 # 任务最大时间 单位 s 秒
  items:
  - clusterName: "aaa" 
    nodename: "" 
    itemStatus: "" 
    podNamespace: "hwx1166232" 
    podName: "nginx-upsteam-dfcf776f4-7sms9"
    commands:
    - commandId: ""
      commandLine: "date;sleep 1;date;sleep 1;" 
    # date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;date;sleep 1;
      toFile: ""
      cmdStatus: ""
      testStorm: "2222"
#  - clusterName: "aaa" 
#    nodename: "n56-249" 
#    podNamespace: "" 
#    podName: ""
#    commands:
#    - commandId: ""
#      commandLine: "ip a"
#      toFile: ""
#      cmdStatus: ""
#      testStorm: "1111"
`
