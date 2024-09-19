package cluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	"github.com/KubeOperator/kubepi/internal/model/v1/capture"
	v1Capture "github.com/KubeOperator/kubepi/internal/model/v1/capture"
	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/server"
	kubeClient "github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gopkg.in/yaml.v2"

	"github.com/KubeOperator/kubepi/pkg/util/podtool"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	defaultcontext "context"

	v1SvcCapture "github.com/KubeOperator/kubepi/internal/service/v1/capture"

	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// flagReloadAll bool = false
	lock sync.Mutex

	allCluster []v1Cluster.Cluster

	defaultCapture capture.Capture // 抓包功能的一些配置，默认代理镜像，每个集群单独配置不同的镜像

	allTasks []capture.Task

	tasksOutPutBuffer map[string]*bytes.Buffer // 缓存输出，全局数据
)

func init() {

	// 初始化
	tasksOutPutBuffer = make(map[string]*bytes.Buffer)
	go loopForever()
}

func getstdoutBuff(CommandId string) *bytes.Buffer {
	buf := tasksOutPutBuffer[CommandId]
	if buf == nil {
		newBuffer := bytes.Buffer{}

		tasksOutPutBuffer[CommandId] = &newBuffer
		buf = &newBuffer
	}
	return buf
}

func cleanstdoutBuff(CommandId string) *bytes.Buffer {
	newBuffer := bytes.Buffer{}

	tasksOutPutBuffer[CommandId] = &newBuffer

	return &newBuffer
}

func reloadBaseInfo() {
	allCluster = updateGetallCluster() // 获取所有集群配置
	err := v1SvcCapture.NewService().GetConfig(&defaultCapture, common.DBOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	// 重新加载配置
	tasks, err := v1SvcCapture.NewService().GetTaskAll(common.DBOptions{})
	if err != nil {
		log.Println(err)
		return
	}

	allTasks = tasks
	log.Println(allTasks)

}
func reloadallCluster() {
	allCluster = updateGetallCluster() // 获取所有集群配置
}

func loopForever() {

	time.Sleep(time.Second * 1) // 5秒后循环检查task任务
	log.SetFlags(log.Llongfile)

	log.Println("start capture")

	// 测试清空
	// v1SvcCapture.NewService().Clean(common.DBOptions{})

	// 准备默认数据 （初始化处理）
	prePareDB()

	reloadBaseInfo()

	// testAddTaskRun()     // 增加一个任务做测试
	// reloadBaseInfo() // 增加一个任务做测试

	// runCommandInPod()       // 测试在容器里运行命令（终端直接返回，而不是保存到pod文件里面）
	// getDbDir() // 测试工作目录（当代理目录执行完之后，需要清理工作用pod）
	// CreateOrDeleteProdProxyPod("aaa", "hwx1166232", "n56-249", true, false) // 测试创建pod
	// time.Sleep(time.Second * 6)
	// CreateOrDeleteProdProxyPod("aaa", "hwx1166232", "n56-249", false, true) // 测试删除pod

	// updateGetallCluster() // 测试获取所有集群

	// prepareNamespace("aaa", "hwx1166232-m") // 确认是否需要创建命名空间
	count := 1
	for {

		/*
			func() {
				lock.Lock()
				defer func() {
					flagReloadAll = false
					lock.Unlock()
				}()

				if flagReloadAll {
					reloadBaseInfo()
				}
			}()
		*/

		func() {
			reloadallCluster() // 不管之前如何处理cluster,自己功能处保留一份，并且每次执行逻辑的时候都获取一份

			for i1, _ := range allTasks {
				for i2, v2 := range allTasks[i1].Items {

					nodeName := ""
					if v2.NodeName != "" {
						nodeName = v2.NodeName
						if _, err := checkNodeExist(v2.ClusterName, nodeName); err != nil {
							allTasks[i1].Items[i2].ItemStatus = capture.StatusWrongPrepare
						}
					} else {
						nodeName, _ = getPodScheduleNodeName(v2.ClusterName, v2.PodNamespace, v2.PodName)

						if _, err := checkPodExist(v2.ClusterName, v2.PodNamespace, v2.PodName); err != nil {
							allTasks[i1].Items[i2].ItemStatus = capture.StatusWrongPrepare
						}
					}
					fixCaptureNamespace := "for-capture"
					// 先创建对应节点的pod
					// "for-capture" 是固定的代理负载的命名空间
					if err := CreateOrDeleteProdProxyPod(v2.ClusterName, fixCaptureNamespace, nodeName, true, false); err != nil {
						// allTasks[i1].Items[i2].Status = capture.StatusWrongPrepare
					}

					// if _, err := checkPodExist(v2.ClusterName, ); err != nil {
					// 	allTasks[i1].Items[i2].Status = capture.StatusWrongPrepare
					// }
					if err := CheckPodRunning(v2.ClusterName, fixCaptureNamespace, nodeName); err != nil {
					} else {
						// log.Println(len(allTasks[i1].Items))
						// log.Println(len(allTasks[i1].Items[i2].Commands))
						for i3, _ := range allTasks[i1].Items[i2].Commands {
							if allTasks[i1].Items[i2].Commands[i3].CmdStatus == "" {
								allTasks[i1].Items[i2].Commands[i3].CmdStatus = capture.StatusRunning

								targetip := "127.0.0.1" // 在节点本身运行，则使用 127.0.0.1 这个特殊地址
								if v2.NodeName != "" {

								} else {
									targetip, _ = getPodRunningIp(v2.ClusterName, v2.PodNamespace, v2.PodName)
								}
								log.Println("run command")
								// 组装代理pod的信息
								proxyPodName := fmt.Sprintf("labcapture-prox-prod-node-%v-%v", v2.ClusterName, nodeName)

								log.Println("runhere:", v2.ClusterName, fixCaptureNamespace, proxyPodName, allTasks[i1].Items[i2].Commands[i3].Commandline, targetip, allTasks[i1].Items[i2].Commands[i3].CommandId, getstdoutBuff(allTasks[i1].Items[i2].Commands[i3].CommandId))

								go func(clustername, podnamespace, podname, targetcommand, targetip, commandid string, stdoutBuff *bytes.Buffer) {
									runCommandInPod(clustername, podnamespace, podname, targetcommand, targetip, commandid, stdoutBuff)

									// 标记命令已经完成
									markCommandFinished(commandid)
								}(v2.ClusterName, fixCaptureNamespace, proxyPodName, allTasks[i1].Items[i2].Commands[i3].Commandline, targetip, allTasks[i1].Items[i2].Commands[i3].CommandId, getstdoutBuff(allTasks[i1].Items[i2].Commands[i3].CommandId))

							}
							// log.Println(allTasks[i1].Items[i2].Commands[i3].Commandline)
							// log.Println(getstdoutBuff(allTasks[i1].Items[i2].Commands[i3].CommandId).Len())

							// log.Println(count)
						}
					}

				}
			}
		}()
		count++
		time.Sleep(time.Second * 2)

	}

}

func getDbDir() error {
	cfg := server.Config()
	log.Println(cfg.Spec.DB.Path) // 获取数据存储的目录
	return nil

}

func updateGetallCluster() []v1Cluster.Cluster {
	svc := cluster.NewService()
	// clu, err := svc.Get(name, common.DBOptions{})
	var conditions commons.SearchConditions

	clusters, _, err := svc.Search(0, 1000, conditions.Conditions, common.DBOptions{})
	if err != nil && err != storm.ErrNotFound {
		return []v1Cluster.Cluster{}
	}
	// result := make([]Cluster, 0)
	// for i := range clusters {
	// 	log.Println(clusters[i].Name) // 打印所有集群名字
	// }
	return clusters
}

var templatePodYaml = `
kind: Pod
apiVersion: v1
metadata:
  name: labcapture-prox-prod-node-%v-%v
  labels:
    app: labcapture-prox-prod-node
    labcapture: labcapture-prox-prod-node
    labcapture-node: %v
spec:
  volumes:
    - name: prochost
      hostPath:
        path: /proc
        type: DirectoryOrCreate
    - name: varrunhost
      hostPath:
        path: /var/run
        type: DirectoryOrCreate
  containers:
    - name: container0
      image: swr.cn-north-4.myhuaweicloud.com/hfbbg4/proxy-prod:v0.1
      command:
        - /bin/bash
        - '-c'
        - ' echo "labcapture: prox-prod";while true; do  date; sleep 13; done'
      resources:
        limits:
          cpu: 2000m
          memory: 2000Mi
        requests:
          cpu: 12m
          memory: 10Mi
      volumeMounts:
        - name: prochost
          readOnly: true
          mountPath: /prod
        - name: varrunhost
          readOnly: true
          mountPath: /var/run
      terminationMessagePath: /dev/termination-log
      terminationMessagePolicy: File
      imagePullPolicy: IfNotPresent
      securityContext:
        privileged: true
        runAsUser: 0
  restartPolicy: Always
  terminationGracePeriodSeconds: 1
  dnsPolicy: ClusterFirst
  serviceAccountName: default
  serviceAccount: default
  nodeName: %v
  hostNetwork: true
  securityContext: {}
  hostname: prox-prod
  subdomain: prox-prod
  schedulerName: default-scheduler
  tolerations:
    - key: node.kubernetes.io/not-ready
      operator: Exists
      effect: NoExecute
      tolerationSeconds: 300
    - key: node.kubernetes.io/unreachable
      operator: Exists
      effect: NoExecute
      tolerationSeconds: 300
  priority: 0
  enableServiceLinks: true
  preemptionPolicy: PreemptLowerPriority
`

// 确认是否需要创建命名空间
func prepareNamespace(clustername, namespace string) error {
	// 获取当亲集群
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)
	if err != nil {
		return err
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	flagCreateNamespace := true
	if allNameSpace, err := clientset.CoreV1().Namespaces().List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range allNameSpace.Items {
			if namespace == v.Name {
				flagCreateNamespace = false
			}
		}

	}
	if flagCreateNamespace {
		newns, err := clientset.CoreV1().Namespaces().Create(defaultcontext.TODO(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		log.Println("create namespace: ", newns)
	}

	return nil
}

func getClusterFromBuffer(clustername string) (*v1Cluster.Cluster, error) {
	// clu, err := svc.Get(clustername, common.DBOptions{})
	for i, v := range allCluster {
		if v.Name == clustername {
			return &allCluster[i], nil
		}
	}
	return nil, errors.New("clustername not exist 22324aa")
}

func checkClusternameExist(clustername string) error {
	flageErrorNoCluster := true
	for _, v := range allCluster {
		if v.Name == clustername {
			flageErrorNoCluster = false
		}
	}
	if flageErrorNoCluster {
		return errors.New("clustername not exist 3243adasf")
	} else {
		return nil
	}

}

func checkNodeExist(clustername, nodename string) (bool, error) {
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})

	clu, err := getClusterFromBuffer(clustername)
	if err != nil {
		return false, err
	}
	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return false, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}
	flagExist := true
	if clusterAllNodes, err := clientset.CoreV1().Nodes().List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range clusterAllNodes.Items {
			if nodename == v.Name {
				flagExist = false
			}
		}

	}
	if flagExist == false {
		return true, nil
	}
	return false, errors.New("node check error ")
}

func checkPodExist(clustername, namespace, podname string) (bool, error) {
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)

	if err != nil {
		return false, err
	}
	if namespace == "" || podname == "" {
		return false, errors.New("namespace or podname is empty")
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return false, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}
	flagExist := true
	if allNameSpacePods, err := clientset.CoreV1().Pods(namespace).List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range allNameSpacePods.Items {
			if podname == v.Name {
				flagExist = false
			}
		}

	}
	if flagExist == false {
		return true, nil
	}
	return false, errors.New("check pod error")
}

func getPodScheduleNodeName(clustername, namespace, podname string) (string, error) {
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)

	if err != nil {
		return "", err
	}
	if namespace == "" || podname == "" {
		return "", errors.New("namespace or podname is empty")
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	flagNotExist := true
	nodeName := ""
	if allNameSpacePods, err := clientset.CoreV1().Pods(namespace).List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range allNameSpacePods.Items {
			if podname == v.Name {
				flagNotExist = false
				nodeName = v.Spec.NodeName
			}
		}

	}
	if flagNotExist == false {
		if nodeName != "" {
			return nodeName, nil
		}
		return "", errors.New("nodeName is empty")
	}
	return "", errors.New("check pod error")
}

func getPodRunningIp(clustername, namespace, podname string) (string, error) {
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)

	if err != nil {
		return "", err
	}
	if namespace == "" || podname == "" {
		return "", errors.New("namespace or podname is empty")
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	flagNotExist := true
	podIp := ""
	if allNameSpacePods, err := clientset.CoreV1().Pods(namespace).List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range allNameSpacePods.Items {
			if podname == v.Name {
				flagNotExist = false
				podIp = v.Status.PodIP
				break
			}
		}

	}
	if flagNotExist == false {
		if podIp != "" {
			return podIp, nil
		}
		return "", errors.New("nodeName is empty")
	}
	return "", errors.New("check pod error")
}

func CheckPodRunning(clustername, namespace, nodename string) error {
	if "" == clustername {
		return errors.New("没有选择 集群 ")
	}
	if "" == namespace {
		return errors.New("没有 命名空间 ")
	}
	if "" == nodename {
		return errors.New("没有 节点名 ")
	}
	// 获取当亲集群
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)
	if err != nil {
		return err
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	if allNameSpacePods, err := clientset.CoreV1().Pods(namespace).List(defaultcontext.TODO(), metav1.ListOptions{}); err != nil {
		log.Println(err)
	} else {
		for _, v := range allNameSpacePods.Items {
			if "Running" == v.Status.Phase {
				return nil
			}
		}

	}
	return errors.New("pod 没有准备好")
}

func markCommandFinished(commandId string) error {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()

	var targetTask *capture.Task
	for i1, _ := range allTasks {
		checkMarkTaskStatusReady := false
		for i2, _ := range allTasks[i1].Items {
			checkMarkItemStatusReady := false
			for i3, _ := range allTasks[i1].Items[i2].Commands {
				if allTasks[i1].Items[i2].Commands[i3].CommandId == commandId {
					allTasks[i1].Items[i2].Commands[i3].CmdStatus = capture.StatusStopped
					log.Println("mark command finished", allTasks[i1].Items[i2].Commands[i3].CommandId)
					log.Println(allTasks[i1].Items[i2].Commands[i3].Commandline)
					log.Println(getstdoutBuff(allTasks[i1].Items[i2].Commands[i3].CommandId).Len())

					targetTask = &allTasks[i1] // 此target将被保存

					// return nil
				} else {
					if allTasks[i1].Items[i2].Commands[i3].CmdStatus == capture.StatusStopped {

					} else {
						checkMarkItemStatusReady = true
					}
				}

			}
			if checkMarkItemStatusReady == false {
				allTasks[i1].Items[i2].ItemStatus = capture.StatusStopped
			}
			if allTasks[i1].Items[i2].ItemStatus == capture.StatusStopped || allTasks[i1].Items[i2].ItemStatus == capture.StatusWrongPrepare {

			} else {
				checkMarkTaskStatusReady = true
			}

		}
		if checkMarkTaskStatusReady == false {
			allTasks[i1].TaskStatus = capture.StatusStopped
		}
	}
	checkAndSaveTask(targetTask)

	log.Println(allTasks)

	return nil
}

func checkTask(post *v1Capture.Task) error {
	if post == nil {
		return errors.New(" post is nil ")
	}
	if post.Id == "" {
		now := time.Now()
		// 转换为Unix时间戳（秒）
		timestamp := now.Unix()
		post.Id = fmt.Sprintf("taskid-%v", timestamp)
	}
	if post.MaxTimeDuring == 0 {
		post.MaxTimeDuring = 120
	}
	if post.StartTime == 0 {
		post.StartTime = time.Now().Unix()
		log.Println(post.StartTime)
	}

	// if post.Status == "" {
	// 	post.Status = capture.StatusRunning
	// }

	for i1, v1 := range post.Items {
		if err := checkClusternameExist(v1.ClusterName); err != nil {
			return err
		}
		if v1.NodeName != "" {
			if _, err := checkNodeExist(v1.ClusterName, v1.NodeName); err != nil {
				return err
			}
		} else {
			if _, err := checkPodExist(v1.ClusterName, v1.PodNamespace, v1.PodName); err != nil {
				return err
			}
		}
		for i2, v2 := range post.Items[i1].Commands {
			if v2.CommandId == "" {
				// 初始化随机种子
				rand.Seed(time.Now().UnixNano())
				// 生成一个随机数
				randomNumber := rand.Intn(1000000000) // 生成0到9之间的随机数
				// fmt.Println("随机数是:", randomNumber)
				post.Items[i1].Commands[i2].CommandId = "commandid-" + fmt.Sprintf("%v", randomNumber) // 初始化随机种子

			}
		}

	}
	return nil
}

func checkAndSaveTask(post *v1Capture.Task) error {
	if err := checkTask(post); err != nil {
		return err
	}

	if err := v1SvcCapture.NewService().PostTask(post, common.DBOptions{}); err != nil {
		return err
	}
	// log.Println("opt save task ok")
	return nil
}

func CreateOrDeleteProdProxyPod(clustername, namespace, nodename string, forcreate, fordelete bool) error {

	if false == forcreate && false == fordelete {
		return errors.New("没有选择相应的功能")
	}
	if "" == clustername {
		return errors.New("没有选择 集群")
	}
	if "" == namespace {
		return errors.New("没有 命名空间")
	}
	if "" == nodename {
		return errors.New("没有 节点名")
	}

	if _, err := checkNodeExist(clustername, nodename); err != nil {
		return err
	}

	// 获取当亲集群
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)
	if err != nil {
		return err
	}

	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// 模板文件写死在这里
	yamlFile := []byte(fmt.Sprintf(templatePodYaml, clustername, nodename, nodename, nodename))

	// Decode YAML into a Pod object
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	decoder := codecs.UniversalDeserializer()

	pod := &v1.Pod{}
	_, _, err = decoder.Decode(yamlFile, nil, pod)
	if err != nil {
		log.Println(err)
		return err
	}

	if true == forcreate {
		prepareNamespace(clustername, namespace) // 确认是否需要创建命名空间

		exist, _ := checkPodExist(clustername, namespace, fmt.Sprintf("labcapture-prox-prod-node-%v-%v", clustername, nodename))
		// Create the Pod in Kubernetes
		// 创建Pod
		if exist != true {
			pod, err = clientset.CoreV1().Pods(namespace).Create(defaultcontext.TODO(), pod, metav1.CreateOptions{})
			if err != nil {
				log.Printf("Error creating pod 3222: %v", err)
			}

			log.Println("Pod %s created successfully!\n", pod.Name)
		}

	}

	if true == fordelete {
		log.Println("Error delete pod: %v", 2)

		// Create the Pod in Kubernetes
		// 创建Pod
		err = clientset.CoreV1().Pods(namespace).Delete(defaultcontext.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Error delete pod: %v", err)
		}

		log.Println("Pod %s delete successfully!\n", pod.Name)
	}

	return nil
}

func testAddTaskRun() { // 手动添加任务

	var data capture.TaskYaml

	if err := yaml.Unmarshal([]byte(varTestTaskNoAdd), &data); err != nil {
		log.Println(err)
	}

	// 将map转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	err = checkAndSaveTask(&data.Task)
	// err = v1SvcCapture.NewService().PostTask(&data.Task, common.DBOptions{})
	if err != nil {
		log.Println(err)

	}
	if true {
		// 打印JSON数据
		log.Println(string(jsonData))

	}
}

func runCommandInPod(clustername, podnamespace, podname, targetcommand, targetip, commandid string, stdoutBuff *bytes.Buffer) (*podtool.PodTool, error) {
	var pt podtool.PodTool
	// svc := cluster.NewService()
	// clu, err := svc.Get(clustername, common.DBOptions{})
	clu, err := getClusterFromBuffer(clustername)
	if err != nil {
		return nil, err
	}
	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	pt = podtool.PodTool{
		Namespace:     podnamespace,
		PodName:       podname,
		ContainerName: "",
		K8sClient:     clientset,
		RestClient:    config,
		ExecConfig:    podtool.ExecConfig{
			// Stdin: request.Stdin,
		},
	}
	/*
		flagFinish := false
		defer func() {
			flagFinish = true
		}()
		go func() {
			i := 0
			time.Sleep(time.Second)
			for {
				aa := pt.ExecConfig.Stdout
				if t, ok := aa.(*bytes.Buffer); ok {
					log.Println(string(t.Bytes()))
				}
				time.Sleep(time.Second)
				i++
				if i == 10 {

				}
				if flagFinish {
					break
				}
			}
		}()
	*/
	// var stdoutBuff bytes.Buffer
	pt.ExecConfig.Stdout = stdoutBuff
	// pt.ExecConfig.Stderr = stdoutBuff
	pt.ExecConfig.Command = []string{
		"bash",
		"-c",
		fmt.Sprintf("export targetip=%v; export whichnsenter=nsenter-prod; export targetcommand='%v' ; proxy-prod -t %v ; ", targetip, targetcommand, commandid),
	}
	// pt.ExecConfig.Tty = false
	var typeAction podtool.ActionType = "Exec"
	typeAction = "Exec"
	err = pt.Exec(typeAction)

	// strs, err := pt.ExecCommand(pt.ExecConfig.Command)
	// log.Println(string(strs))
	if err != nil {
		return nil, err
	}

	// time.Sleep(time.Second * 30)

	return &pt, nil
}

func (h *Handler) ConfigGet() iris.Handler {
	return func(ctx *context.Context) {
		lock.Lock()
		// flagReload = true
		defer lock.Unlock()

		capture := v1Capture.Capture{}
		err := h.captureService.GetConfig(&capture, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", capture)
		return
	}
}

func (h *Handler) ConfigPost() iris.Handler {
	return func(ctx *context.Context) {
		lock.Lock()
		// flagReload = true
		defer lock.Unlock()

		var post v1Capture.Capture
		if err := ctx.ReadJSON(&post); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if post.Config.DefaultImg == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "post.Config.DefaultImg")
			return
		}
		if err := h.captureService.PostConfig(&post, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		ctx.StatusCode(iris.StatusOK)
		return
	}
}

func (h *Handler) TaskGetAll() iris.Handler {
	return func(ctx *context.Context) {
		lock.Lock()
		// flagReload = true
		defer lock.Unlock()
		tasks, err := h.captureService.GetTaskAll(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", tasks)

	}
}

func (h *Handler) TaskPost() iris.Handler {
	return func(ctx *context.Context) {
		lock.Lock()
		// flagReload = true
		defer lock.Unlock()
		var post v1Capture.Task
		if err := ctx.ReadJSON(&post); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		if err := checkAndSaveTask(&post); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		// 返回，以备注id和状态 running,stop
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", post)

		return
	}
}
