package capture

// capture.StatusRunning
const (
	StatusRunning      = "StatusRunning"      // 任务正在运行
	StatusWrongPrepare = "StatusWrongPrepare" // 集群被删除，节点被删除，pod已经不存在了
	StatusStopped      = "StatusStopped"      // 任务已经终止
)

type Capture struct {
	// v1.Metadata `storm:"inline"`
	Id string `json:"id" storm:"id,index,unique" yaml:"id"` // taskid

	Config Config `json:"config" storm:"inline" yaml:"config"`
}

type Config struct {
	DefaultImg string       `json:"defaultImg" storm:"inline" yaml:"defaultImg"` // 默认使用的镜像 proxy
	ArrImg     []ArrImgItem `json:"arrImg" storm:"inline" yaml:"arrImg"`
}

type ArrImgItem struct {
	ClusterName string `json:"clusterName" storm:"inline" yaml:"clusterName"`
	Img         string `json:"img" storm:"inline" yaml:"img"`
}

type TaskYaml struct {
	Task Task `json:"task" storm:"taskId,index,unique" yaml:"task"`
}

type Task struct {
	// Name        string `json:"name" storm:"unique" `
	// Description string `json:"description"`
	// UUID        string `json:"uuid" storm:"id,index,unique"`

	Id         string `json:"id" storm:"id,index,unique" yaml:"id"`        // taskid
	TaskStatus string `json:"taskStatus" storm:"inline" yaml:"taskStatus"` // StatusWrongPrepare , StatusRunning,StatusStopped

	StartTime     int64      `json:"startTime" storm:"inline" yaml:"startTime"`
	MaxTimeDuring int        `json:"maxTimeDuring" storm:"inline" yaml:"maxTimeDuring"`
	Items         []TaskItem `json:"items" storm:"inline" yaml:"items"` // 为什么要是数组？  一个任务需要多个监听（跨集群，跨节点）
}

type TaskItem struct {
	// ItemId   string    `json:"itemId" storm:"inline" yaml:"itemId"` //
	ClusterName string `json:"clusterName" storm:"inline" yaml:"clusterName"`

	NodeName string `json:"nodename" storm:"inline" yaml:"nodename"` // 与 （ PodNamespace ，PodName ） 二选一

	PodNamespace string `json:"podNamespace" storm:"inline" yaml:"podNamespace"` // 与 NodeName 二选一
	PodName      string `json:"podName" storm:"inline" yaml:"podName"`           // 与 NodeName 二选一

	ItemStatus string    `json:"itemStatus" storm:"inline" yaml:"itemStatus"` // 所有命令执行完成，方便统计，或者说pod或者节点已经不存在时直接结束
	Commands   []Command `json:"commands" storm:"inline" yaml:"commands"`     // 为什么要是数组？  假设一个任务 tcpdump 【即要在页面显示又需要保存为文件的方式】 , 同时抓取ipip或者vxlan这样好确认使用的协议
}

type Command struct {
	CommandId   string `json:"commandId" storm:"inline" yaml:"commandId"` // 用于强制退出时
	Commandline string `json:"commandLine" storm:"inline"  yaml:"commandLine"`
	ToFile      string `json:"toFile" storm:"inline" yaml:"toFile"`
	CmdStatus   string `json:"cmdStatus" storm:"inline" yaml:"cmdStatus"`

	StartTime int64 `json:"startTime" storm:"inline" yaml:"startTime"`
	EndTime   int64 `json:"endTime" storm:"inline" yaml:"endTime"`
}
