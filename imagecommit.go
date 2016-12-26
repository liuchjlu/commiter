package main

import (
	"github.com/liuchjlu/commiter/cli"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/commiter/ip"
)

func main() {
	logFilename := "/tmp/imgcommit.log"
	logFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetLevel(log.DebugLevel)

	log.Infoln("main.main():Start Imgcommit Main")
	ip, _ := ip.Getcontaineripfromcommiter("204-databus-19216811123:161221", "http://192.168.11.52:2379", "204", "databus")
	log.Infof("container ip:%+v\n", ip)
	cli.Run()
}
