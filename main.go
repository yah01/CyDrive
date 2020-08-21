package main

import (
	"github.com/yah01/CyDrive/config"
	"github.com/yah01/cyflag"
)

var (
	dbConfig config.Config

	isOnline      bool
	serverAddress string
)

func init() {
	// Parse args
	//cyflag.BoolVar(&isServer, "--server", false, "whether run as a cdv cdpServer")
	cyflag.BoolVar(&isOnline, "--online", false, "whether is online")
	cyflag.StringVar(&serverAddress, "-h", "localhost", "set the CyDrive Server address")
	cyflag.Parse()

	// Read DB config
	//dbConfigFile, err := ioutil.ReadFile("db_config.yaml")
	//if err != nil {
	//	panic(err)
	//}
	//if err = yaml.Unmarshal(dbConfigFile, &dbConfig); err != nil {
	//	panic(err)
	//}

	// Open the log file with level INFO
	//infoLogFile, err := os.Open("info_log.txt")
	//if err != nil {
	//	panic(err)
	//}
	//infoLog = log.New(infoLogFile, "INFO", log.LstdFlags)
}

func main() {
	dbConfig.UserStoreType = "mem"
	InitServer(dbConfig)
	RunServer()
}
