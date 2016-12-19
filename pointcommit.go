package main
import (
	"io"
	"os"
	//"fmt"

	// "github.com/docker/docker/api/types"
	// dockerclient "github.com/docker/docker/client"
	// "golang.org/x/net/context"
	// "github.com/liuchjlu/pointcommit/etcdclient"
	log "github.com/Sirupsen/logrus"
    "github.com/liuchjlu/pointcommit/cli"
)
func main() {
	logFilename := "/tmp/pointcommit.log"
	logFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetLevel(log.DebugLevel)

	log.Infoln("main.main():Start Pointcommit Main")

	//cli.Getname("http://192.168.11.52:2379")
	//cli.Commit(os.Args[1])
	cli.Run()
	//cli.Commit("http://192.168.11.52:2379","192.168.11.51:5000/")


}



