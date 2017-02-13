# commiter

1.commiter简介
commiter是一款由go语言开发的，基于kubernetes集群自动打包docker镜像的工具。该工具在打包镜像的功能基础上添加了：将容器ip和镜像名信息保存到etcd中功能，实现了生产环境下应用的可迁移，同时自动为kubernetes RC文件pod副本中的每一个pod的生成对应的接管RC文件。
2.commiter使用
  Usage: commiter COMMAND [args...]
  Version: 0.2.0
  Author:liuchjlu
  Email:liucaihong@iie.ac.cn
  Commands:
      manage  [apiserver address] [etcd path]           Create the config in etcd.
      import  [localyaml path]    [etcd path]  [0/1/2]  Import the status of pod from yaml. 0:don't commit; 1:commit & push to registry; 2:commit & don't push to registry.
      commit  [etcd path]                               Image commit.
      help
  Use case
      commiter manage 192.168.11.52:8080 $ETCDPATH
      commiter import test.yaml $ETCDPATH 1
      commiter commit $ETCDPATH 

1)import命令：
在K8s主节点执行，将系统名、组件名、物理ip、镜像仓库、打包状态信息导入到etcd中/commiter/images/ip/app/component/state和/commiter/images/ip/app/component/reporsitory中。
import.yaml样例如下：
App: 306
Component: mango
Repository: 192.168.11.51:5000
Ips:
  - Ip: 192.168.11.52
导入import.yaml后的etcd目录结构：
/commiter/images/192.168.11.52/306/mango/state
/commiter/images/192.168.11.52/306/mango/repository
其中state的value保存镜像打包状态信息，目前仅支持0（不可打包）和1（打包并上传仓库）状态，repository的value保存的为mango组件镜像仓库地址

2）commiter命令：
在需要打包镜像的节点执行，首先获取本机物理ip，然后根据ip去etcd中查询相应目录下的系统名(以下用app表示)和组件名（以下用component表示）（可能有多对app,component，用切片保存），然后从docker接口获取本机运行的容器信息（其中包括容器id，容器ip，启动容器的镜像名，容器名字等，其中容器ip获取先从etcd总/commiter/ips/app-component-ip:tag下查询，如果没有找到在去/assigner/ips下查询容器ip），如果容器名字中含有apprc-component字段则对容器打包镜像，镜像名字命名为registory/app-component-ip:tag，如:192.168.11.51:5000/201-databus-192.168.11.234:20161228，然后将键值对<镜像名:容器ip>写到etcd中/commiter/ips目录下，然后将镜像上传到镜像仓库。如果以后需要迁移的话直接用新打的镜像启动容器的时候会先从etcd中/commiter/ips目录下获取ip，解决了通过assigner保存ip信息的不可迁移的问题。在完成每一个pod镜像打包之后，根据FTP服务器相对应/app/component目录下的apprc-component*.yaml模板文件（如果对应rc是分层的，则component*写成相对应层信息，如：204rc-databus-1.yaml或306rc-mysql.yaml），为该pod生成一个接管的rc文件，接管rc文件生成后命名为：apprc-component-ip.yaml，并将该接管rc文件上传到ftpserver中/app/componenet目录下。

3）manage命令：
将k8s apiserver地址和ftpserver地址导入到etcd /commiter/config/下

3.操作步骤：
1) ftp服务端和客户端安装,ftp用户名为commiter，密码为woaiiie@iiecloud，ftp服务工作目录为/home/commiter
2) 将kubectl,commiter拷贝到node节点/usr/bin，并添加可执行权限，将download.template和upload.template拷贝到node节点/etc/.commiter/template目录下
3) 从etcd中获取rc文件保存到本地，保存到主节点/home/commiter/app/component目录下，文件命名为apprc-component.yaml（如果为分层则需要带上层的信息，如：204rc-databus-1.yaml，不分层如：204rc-mysql.yaml）
4) 主节点执行commiter manage命令（该步骤为集群初始化，仅执行一次）
4) 编写import.yaml文件，并在主节点执行commiter import import.yaml $ETCDPATH 1，完成组件打包状态更新
5) 各从节点执行commiter commit $ETCDATH命令，进行打包、上传、K8s添加label、生成接管pod的rc等操作
6) 检查主节点相应/home/commiter/app/componnet目录下是否有pod的rc文件生成、pod的ip label是否添加、pod rc是否create

4.使用注意事项
1）为了确保模板文件跟线上系统完全一致，模板rc文件直接从etcd中获取，然后在改成模板。
2）生成模板需要改动地方：
	a）rc.memadata.name改为："name":"{{.App}}rc-{{.Component}}-{{.ContainerIp}}"；
	b）rc.metadata.labels：增加一项:"ip":"{{.ContainerIp}}"；
	c）rc.spec.replicas中副本数改为1；
	d）rc.spec.template.medata.labels：增加一项："ip":"{{.ContainerIp}}"；
	e）rc.spec.template.spec.containers.image改为："image":"{{.Image}}"；
	f) rc.spec.selector：增加一项："ip":"{{.ContainerIp}}"。
3）在commiter commit运行之前模板rc文件需要放到ftp服务器/app/component目录对应app/component下
4）K8s每一个node需要添加如下文件：
	a)将download.template和upload.template文件拷贝到/etc/.commiter/template目录下
	b)将kubectl、assigner、commiter文件拷贝到/ust/bin目录下。




