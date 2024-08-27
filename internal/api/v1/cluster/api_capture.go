package cluster

import (
	"errors"
	"fmt"
	"sync"
	"time"

	v1Capture "github.com/KubeOperator/kubepi/internal/model/v1/capture"
	kubeClient "github.com/KubeOperator/kubepi/pkg/kubernetes"

	"github.com/KubeOperator/kubepi/pkg/util/podtool"
	"k8s.io/client-go/kubernetes"

	// v1SvcCapture "github.com/KubeOperator/kubepi/internal/service/v1/capture"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var (
	count int
	lock  sync.Mutex
)

func init() {

	go func() {
		time.Sleep(time.Second * 5) // 5秒后循环检查task任务
		tmp1()
		for {
			lock.Lock()
			if count == 1 {
				// 重新加载配置
			}
			// fmt.Println("time loop check task")
			time.Sleep(time.Second * 2)
			count = 0
			lock.Unlock()
		}

	}()
}

func tmp1() (*podtool.PodTool, error) {
	var pt podtool.PodTool
	svc := cluster.NewService()
	clu, err := svc.Get("aaa", common.DBOptions{})
	if err != nil {
		return nil, errors.New("error1")
	}
	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return nil, errors.New("error2")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.New("error3")
	}
	pt = podtool.PodTool{
		Namespace:     "hwx1166232",
		PodName:       "testapp5-0",
		ContainerName: "",
		K8sClient:     clientset,
		RestClient:    config,
		ExecConfig:    podtool.ExecConfig{
			// Stdin: request.Stdin,
		},
	}
	bts, err := pt.ExecCommand([]string{
		"bash",
		"-c",
		"ls /util-linux/nsenter;",
	})
	fmt.Println(string(bts))
	return &pt, nil
}

func (h *Handler) ConfigGet() iris.Handler {
	return func(ctx *context.Context) {
		lock.Lock()
		count = 1
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
		count = 1
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
		count = 1
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
		count = 1
		defer lock.Unlock()
		var post v1Capture.Task
		if err := ctx.ReadJSON(&post); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if post.Id == "" {
			now := time.Now()
			// 转换为Unix时间戳（秒）
			timestamp := now.Unix()
			post.Id = fmt.Sprintf("%v", timestamp)
		}
		if post.TimeDuring == 0 {
			post.TimeDuring = 120
		}
		if post.Status == "" {
			post.Status = "running"
		}
		fmt.Println(post)
		if err := h.captureService.PostTask(&post, common.DBOptions{}); err != nil {
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
