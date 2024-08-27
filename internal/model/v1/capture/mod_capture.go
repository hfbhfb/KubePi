package capture

import (
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
)

type Capture struct {
	v1.Metadata `storm:"inline"`
	Config      Config `json:"config" storm:"inline" `
}

type Config struct {
	DefaultImg string       `json:"default_img" storm:"inline" `
	ArrImg     []ArrImgItem `json:"arr_img" storm:"inline"`
}

type ArrImgItem struct {
	ClusterName string `json:"cluster_name" storm:"inline"`
	Img         string `json:"img" storm:"inline"`
}

type Task struct {
	Id         string     `json:"id" storm:"id,index,unique"`
	Status     string     `json:"status" storm:"inline"`
	TimeDuring int        `json:"time_during" storm:"inline"`
	Items      []TaskItem `json:"items" storm:"inline"`
}

type TaskItem struct {
	Podip    string    `json:"podip" storm:"inline"`
	Nodeip   string    `json:"nodeip" storm:"inline"`
	Commands []Command `json:"commands" storm:"inline"`
}

type Command struct {
	Commandline string `json:"command_line" storm:"inline"`
	ToFile      string `json:"to_file" storm:"inline"`
}
