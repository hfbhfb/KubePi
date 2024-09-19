


# 运行后端
cd /d/projs/fork8sdir/code-ref/kube-pi/cmd/server
gowatch


# 编译docker镜像
make my_all

# 备份到备份目录 
#rsync -avh --progress /d/projs/fork8sdir/code-ref/kube-pi /d/tmp/forbuild/
#cd /d/tmp/forbuild/kube-pi;  make clean -f prepare-Makefile;make build_docker
cd /d/tmp/forbuild/kube-pi; code .


git log --stat --name-status
#    关注的文件
git add internal/api/v1/cluster/api_capture.go
git add internal/api/v1/cluster/api_capture_fortask.go
git add internal/api/v1/cluster/cluster.go
git add internal/model/v1/capture/mod_capture.go
git add internal/service/v1/capture/errors.go
git add internal/service/v1/capture/svc_capture.go
git add web/kubepi/src/api/capture.js
git add web/kubepi/src/business/cluster-management/tcpdump/index.vue
git add web/kubepi/src/router/modules/clusters.js
git add zzz项目333.sh



# 运行 dashboard web
# 前端项目代码2
cd /d/projs/fork8sdir/code-ref/kube-pi/web/dashboard
code-insiders .



# 运行 terminal web
# 前端项目代码3
cd /d/projs/fork8sdir/code-ref/kube-pi/web/terminal
npm run start


# docker run -v /root/tmp/aatmp1:/root/tmp/aatmp1 -it swr.cn-north-4.myhuaweicloud.com/hfbbg4/util-linux:util-linux bash
# util-linux$ git branch
#   master
# * prod
# util-linux$ git remote -v
# origin  git@github.com:hfbhfb/util-linux.git (fetch)
# origin  git@github.com:hfbhfb/util-linux.git (push)
# proc->prod重映射项目
# swr.cn-north-4.myhuaweicloud.com/hfbbg4/util-linux:util-linux    /util-linux
#	command: ["/bin/sh","-c"," export targetip=10.0.230.176; export whichnsenter=nsenter-prod; export targetcommand='tcpdump -nnSX ;date;sleep 2;date;sleep 2;ip a;' ; proxy-prod -t 25328524 ;while true; do  date; sleep 13; done"]
# export forkill=25328524 ;proxy-prod
cd /d/projs/c-cxx/util-linux
code .



# 其它技术细节 （ kv数据库 Storm is a simple and powerful toolkit for BoltDB.  ）
https://github.com/asdine/storm

