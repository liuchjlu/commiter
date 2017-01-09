package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/commiter/cmd"
	"github.com/liuchjlu/commiter/etcdclient"
	"os"
	"strings"
	"text/template"
)

var download string = "/etc/.commiter/sh/download.sh"
var upload string = "/etc/.commiter/sh/upload.sh"

//create podrc.yaml,download.sh,upload.sh. And put the podrc.yaml to the ftpserver:/home/commiter/app/component
func Rc(app, component, containerip, containername, image, etcdpath string) {
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.RC() Error when create etcdclient:%+v\n", err)
		return
	}
	apiserver, err := client.GetAbsoluteKey("/commiter/config/apiserver")
	if err != nil {
		log.Fatalf("cli.RC() Error when get apiserver address from etcd:%+v\n", err)
	}
	ftpserver, err := client.GetAbsoluteKey("/commiter/config/ftpserver")
	if err != nil {
		log.Fatalf("cli.RC() Error when get ftpserver address from etcd:%+v\n", err)
	}
	BuildDownload(app, component, containername, ftpserver)
	BUildUpload(app, component, containerip, containername, ftpserver)
	//get component1 which include the info of level
	podname := strings.Split(containername, "_")[2]
	var component1 string
	if len(strings.Split(podname, "-")) == 3 {
		component1 = component
	} else if len(strings.Split(podname, "-")) == 4 {
		component1 = component + "-" + strings.Split(podname, "-")[2]
	} else {
		log.Fatalln("cli.BuildDownload(): There are some error when get component name.")
	}
	podyamlname := "/home/commiter/" + app + "rc-" + component1 + "-" + containerip + ".yaml"
	//todo: run download.sh,run upload.sh
	commandname1 := "bash"
	commandname2 := "kubectl"
	params1 := []string{download}
	params2 := []string{upload}
	params3 := []string{"--server=" + apiserver, "label", "pods", strings.Split(containername, "_")[2], "ip=" + containerip}
	params4 := []string{"--server=" + apiserver, "create", "-f", podyamlname}
	//get the app-componnet1.yaml from ftpserver
	if err := cmd.ExecCommand(commandname1, params1); !err {
		log.Fatalf("cli.RC() Error when run the script of download rc template:%+v\n", err)
	} else {
		log.Infoln("cli.RC() Download rc template success")
	}
	//create rc for pod
	BuildYaml(app, component, containerip, image, containername)
	//upload the rc of pod to ftpserver
	if err := cmd.ExecCommand(commandname1, params2); !err {
		log.Fatalf("cli.RC() Error when run the script of upload rc template:%+v\n", err)
	} else {
		log.Infoln("cli.RC() Upload rc template success")
	}
	//add label of ip:containerip for pod
	if err := cmd.ExecCommand(commandname2, params3); !err {
		log.Fatalf("cli.RC() Error when add label ip to the pod:%+v,err:%+v\n", strings.Split(containername, "_")[2], err)
	} else {
		log.Infof("cli.RC() Add label ip successful. Pod:%+v. IP:%+v\n", strings.Split(containername, "_")[2], containerip)
	}
	//create rc for the pod rc
	if err := cmd.ExecCommand(commandname2, params4); !err {
		log.Fatalf("cli.RC() K8s create rc for the pod rc:%+v\n", err)
	}
	return
}

//create podrc.yaml, and save to the local path:/home/commiter
func BuildYaml(app, component, containerip, image, containername string) {
	filepath := "/home/commiter/"
	os.MkdirAll(filepath, 0755)
	podname := strings.Split(containername, "_")[2]
	log.Infof("cli.BuildYaml() podname=%+v\n", podname)
	if len(strings.Split(podname, "-")) == 3 {
		component = component
		log.Debugf("cli.BUildYaml component:%+v\n", component)
	} else if len(strings.Split(podname, "-")) == 4 {
		component = component + "-" + strings.Split(podname, "-")[2]
		log.Debugf("cli.BUildYaml component:%+v\n", component)
	} else {
		log.Fatalf("cli.BuildDownload(): There are some error when get component name.")
	}
	podrcname := app + "rc-" + component + "-" + containerip + ".yaml"
	filename := filepath + podrcname
	dstFile, err := os.Create(filename)
	defer dstFile.Close()
	if err != nil {
		log.Fatalf("cli.BuildYaml() os.create(file) failed:%+v\n", err)
		return
	}
	templatename := app + "rc-" + component + ".yaml"
	log.Debugf("cli.BuildYaml =====templatename:%+v\n", templatename)
	t, err := template.New(templatename).ParseFiles("/home/commiter/" + templatename)
	if err != nil {
		log.Fatalln("cli.BuildYaml() use template failed:%+v\n", err)
		return
	}
	data := struct {
		App         string
		Component   string
		ContainerIp string
		Image       string
	}{
		App:         app,
		Component:   component,
		ContainerIp: containerip,
		Image:       image,
	}
	err = t.Execute(dstFile, data)
	if err != nil {
		log.Fatalf("cli.BuildYaml() execute template fail:%+v\n", err)
		return
	}
	dstFile.WriteString("\n")
	log.Infoln("BuildYaml complete.")
	return
}

//create download.sh
func BuildDownload(app, component, containername, ftpserver string) {
	filepath := "/etc/.commiter/sh/"
	os.MkdirAll(filepath, 0755)
	filename := filepath + "download.sh"
	dstFile, err := os.Create(filename)
	defer dstFile.Close()
	if err != nil {
		fmt.Println("os.create(file) failed.")
		return
	}
	t, err := template.New("download.template").ParseFiles("/etc/.commiter/template/download.template")
	if err != nil {
		fmt.Println("BuildDownload failed.")
		return
	}
	podname := strings.Split(containername, "_")[2]
	log.Infof("cli.BUildDownload ****Container Name: %+v\n", containername)
	log.Infof("cli.BuildDownload ****Pod Name: %+v\n", podname)
	lenth := len(strings.Split(podname, "-"))
	log.Debugf("cli.BuildDownload() ===len(strings.Split(podname, -):%+v\n", lenth)
	var component1 string
	if len(strings.Split(podname, "-")) == 3 {
		component1 = component
	} else if len(strings.Split(podname, "-")) == 4 {
		component1 = component + "-" + strings.Split(podname, "-")[2]
	} else {
		log.Fatalln("cli.BuildDownload(): There are some error when get component name.")
	}
	templatename := app + "rc-" + component1 + ".yaml"
	data := struct {
		App          string
		Component    string
		TemplateName string
		FtpServer    string
		UserName     string
		PassWord     string
	}{
		App:          app,
		Component:    component,
		TemplateName: templatename,
		FtpServer:    ftpserver,
		UserName:     "commiter",
		PassWord:     "woaiiie@iiecloud",
	}
	err = t.Execute(dstFile, data)
	if err != nil {
		fmt.Println("t.execute failed.")
		fmt.Println(err)
		return
	}
	dstFile.WriteString("\n")
	fmt.Println("BuildDownload complete.")
	return
}

//create upload.sh
func BUildUpload(app, component, containerip, containername, ftpserver string) {
	filepath := "/etc/.commiter/sh/"
	os.MkdirAll(filepath, 0755)
	filename := filepath + "upload.sh"
	dstFile, err := os.Create(filename)
	defer dstFile.Close()
	if err != nil {
		fmt.Println("os.create(file) failed.")
		return
	}
	t, err := template.New("upload.template").ParseFiles("/etc/.commiter/template/upload.template")
	if err != nil {
		fmt.Println("BUildUpload failed.")
		return
	}
	var component1 string
	podname := strings.Split(containername, "_")[2]
	if len(strings.Split(podname, "-")) == 3 {
		component1 = component
	} else if len(strings.Split(podname, "-")) == 4 {
		component1 = component + "-" + strings.Split(podname, "-")[2]
	} else {
		log.Errorln("cli.BuildDownload(): There are some error when get component name.")
	}
	data := struct {
		App       string
		Component string
		PodRc     string
		FtpServer string
		UserName  string
		PassWord  string
	}{
		App:       app,
		Component: component,
		PodRc:     app + "rc-" + component1 + "-" + containerip + ".yaml",
		FtpServer: ftpserver,
		UserName:  "commiter",
		PassWord:  "woaiiie@iiecloud",
	}
	err = t.Execute(dstFile, data)
	if err != nil {
		fmt.Println("t.execute failed.")
		fmt.Println(err)
		return
	}
	dstFile.WriteString("\n")
	fmt.Println("BUildUpload complete.")
	return
}
