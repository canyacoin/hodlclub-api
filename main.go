package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/ethereum/go-ethereum/ethclient"
	logging "github.com/op/go-logging"
)

var (
	serviceID = "hodlclub-api"
	router    *gin.Engine
	logger    = logging.MustGetLogger("main")
	startedAt = time.Now()
	// ethNodeURL = "https://mainnet.infura.io/v3/4feb084407684fbcafdeac4dd88a2e83"
	ethNodeURL = "https://inherently-fast-salmon.quiknode.io:443/568eb2c7-8118-47e6-a518-7549987fef57/V2Ul02LX0aGtV7TVObmL8A==/"
	// ethNodeURL       = "http://localhost:8545"
	canyCoinContract = "0x1d462414fe14cf489c7a21cac78509f4bf8cd7c0"
	startingBlock    = int64(4332959)
	// startingBlock    = int64(6664087)
	nodeConnection   *ethclient.Client
	contractInstance *CanYaCoin
	db               *gorm.DB
)

func init() {
	var err error
	logFormatter := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.10s} %{id:03x}%{color:reset} %{message}`)
	logging.SetFormatter(logFormatter)
	consoleBackend := logging.NewLogBackend(os.Stdout, "", 0)
	consoleBackend.Color = true
	logging.SetLevel(logging.DEBUG, "main")

	dsn := mustGetenv("DSN")
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("unable to connect to dsn: %s with error: %s", dsn, err.Error())
	}

	var tables []interface{}
	tables = append(tables, &BalanceRecord{}, &LastBlock{}, &BlacklistRecord{})
	for i := range tables {
		err = db.AutoMigrate(tables[i]).Error
		if err != nil {
			logger.Fatalf("unable to automigrate table with error: %s", err.Error())
		}
	}

	var count uint
	err = db.Table("last_blocks").Count(&count).Error
	if err != nil {
		logger.Fatalf("unable to count last_blocks: %s", err.Error())
	}

	if count == 0 {
		logger.Infof("loading initial CanYaCoin deployment start block: %d", startingBlock)
		err = updateLastBlockProcessed(startingBlock)
		if err != nil {
			logger.Errorf("error updating last processed block height: %s", err.Error())
		}
	}

	router = gin.Default()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://hodl-club.firebaseapp.com", "https://hodl.canya.com"}
	router.Use(cors.New(config))

	router.GET("/status", statusHandler)
	router.GET("/hodlers", hodlersHandler)
	router.GET("/process-blockchain", processBlockchainHandler)
	router.GET("/", rootHandler)
}

func main() {
	router.Run()
}
