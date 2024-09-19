package capture

import (
	"fmt"

	v1Capture "github.com/KubeOperator/kubepi/internal/model/v1/capture"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
)

type Service interface {
	common.DBService
	GetConfig(capture *v1Capture.Capture, options common.DBOptions) error
	PostConfig(capture *v1Capture.Capture, options common.DBOptions) error

	PostTask(task *v1Capture.Task, options common.DBOptions) error
	GetTaskAll(options common.DBOptions) ([]v1Capture.Task, error)

	Clean(options common.DBOptions) error
}

func NewService() Service {
	return &captureSvc{
		DefaultDBService: common.DefaultDBService{},
	}
}

type captureSvc struct {
	common.DefaultDBService
}

func (c *captureSvc) GetConfig(capture *v1Capture.Capture, options common.DBOptions) error {
	db := c.GetDB(options)
	var capturelist []v1Capture.Capture
	if err := db.All(&capturelist); err != nil {
		return err
	}
	if len(capturelist) == 0 {
		tmp1 := v1Capture.Capture{}
		// tmp1.Config.DefaultImg = "aaa"
		tmp1.Id = "default"
		// tmp1.Name = "default"
		// tmp1.UUID = uuid.New().String()
		db.Save(&tmp1)
		capturelist = append(capturelist, tmp1)
	}
	if len(capturelist[0].Config.ArrImg) == 0 {
		capturelist[0].Config.ArrImg = []v1Capture.ArrImgItem{}
	}
	*capture = capturelist[0]
	return nil

}

func (c *captureSvc) PostConfig(capture *v1Capture.Capture, options common.DBOptions) error {
	tmp1 := v1Capture.Capture{}
	c.GetConfig(&tmp1, options)
	tmp1.Config = capture.Config

	db := c.GetDB(options)
	return db.Save(&tmp1)

}

func (c *captureSvc) Clean(options common.DBOptions) error {
	db := c.GetDB(options)
	// fmt.Println("clean")
	err := db.Drop(&v1Capture.Capture{})
	if err != nil {
		fmt.Println(err)
		// return err
	}
	err = db.Drop(&v1Capture.Task{}) //db.Drop("Task")
	if err != nil {
		fmt.Println(err)

		// return err
	}

	return nil

}

func (c *captureSvc) PostTask(task *v1Capture.Task, options common.DBOptions) error {
	db := c.GetDB(options)
	return db.Save(task)
}

func (c *captureSvc) GetTaskAll(options common.DBOptions) ([]v1Capture.Task, error) {
	db := c.GetDB(options)
	var tasklist []v1Capture.Task
	if err := db.All(&tasklist); err != nil {
		return []v1Capture.Task{}, err
	}
	for i, _ := range tasklist {
		if len(tasklist[i].Items) == 0 {
			tasklist[i].Items = []v1Capture.TaskItem{}
		}
	}
	return tasklist, nil
}
