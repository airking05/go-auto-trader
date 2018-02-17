package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/airking05/go-auto-trader/config"
	"github.com/airking05/go-auto-trader/logger"
	"github.com/airking05/go-auto-trader/models"
	"github.com/airking05/go-auto-trader/repositories"
	"github.com/airking05/go-auto-trader/services"
	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Server struct {
	TradersService services.TraderServiceImpl `inject:"inline"`
}

func help_and_exit() {
	fmt.Fprintf(os.Stderr, "%s config.yml\n", os.Args[0])
	os.Exit(1)
}

func migrateDB(db *gorm.DB) {
	models.Migrate(db)
}

func main() {
	if len(os.Args) != 2 {
		help_and_exit()
	}
	confPath := os.Args[1]
	conf := config.ReadConfig(confPath)

	fmt.Printf(strconv.FormatBool(conf.Debug))

	// open db
	db, err := gorm.Open("mysql", conf.DBConnection)
	db.LogMode(false)
	if err != nil {
		panic(errors.Wrap(err, "failed to open db"))
	}
	migrateDB(db)
	if err := models.Migrate(db); err != nil {
		panic(err)
	}

	botmap := services.NewBotMap()
	server := Server{}
	// init DI
	var g inject.Graph
	err = g.Provide(
		&inject.Object{Value: db},

		// Repositories
		&inject.Object{Value: &repositories.TraderGormRepository{}},
		&inject.Object{Value: &repositories.ChartRepositoryGorm{}},
		&inject.Object{Value: &repositories.OrderDataStorage{}},
		&inject.Object{Value: &repositories.PositionStorage{}},

		// Services
		&inject.Object{Value: &services.TraderServiceImpl{}},

		// Server
		&inject.Object{Value: &server},
	)
	if err != nil {
		panic(err)
	}
	if err := g.Populate(); err != nil {
		panic(err)
	}
	server.TradersService.BotMap = botmap
	logger.Get().Info("starting job_server...")

	traderConfig := &models.TraderGorm{
		ExchangeID:               conf.TraderConfig.ExchangeID,
		APIKey:                   conf.TraderConfig.APIKey,
		SecretKey:                conf.TraderConfig.SecretKey,
		Trading:                  conf.TraderConfig.Trading,
		Settlement:               conf.TraderConfig.Settlement,
		LossCutRate:              conf.TraderConfig.LossCutRate,
		ProfitTakeRate:           conf.TraderConfig.ProfitTakeRate,
		AssetDistributionRate:    conf.TraderConfig.AssetDistributionRate,
		PositionType:             conf.TraderConfig.PositionType,
		Duration:                 conf.TraderConfig.Duration,
		WaitLimitSecond:          conf.TraderConfig.WaitLimitSecond,
		MakePositionLogicsString: conf.TraderConfig.MakePositionLogicsString,
	}
	server.TradersService.Run(traderConfig)
	for {
		time.Sleep(time.Minute)
	}

}
