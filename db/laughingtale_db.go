package db

import (
	"context"
	"net"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/shashank-mugiwara/laughingtale/conf"
	"github.com/shashank-mugiwara/laughingtale/logger"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var laughingtaleDb *gorm.DB

func GetlaughingtaleDb() *gorm.DB {
	return laughingtaleDb
}

func InitGormPool() {

	laughingtale_logger := logger.GetLaughingTaleLogger()

	var sqlxMasterDb *sqlx.DB
	var err error
	var driver database.Driver

	// First create master node connection
	if conf.ApplicationSetting.RunType == "local" {
		laughingtale_logger.Info("Using Local Database Connection [MASTER_NODE].")

		// Information for local connection
		localAuth := &LocalAuth{
			DatabaseUser:     conf.PostgresSetting.Username,
			DatabaseHost:     conf.PostgresSetting.MasterNodeHost,
			DatabasePort:     conf.PostgresSetting.Port,
			DatabaseName:     conf.PostgresSetting.Database,
			DatabasePassword: conf.PostgresSetting.Password,
			DatabaseSchema:   conf.PostgresSetting.TablePrefix,
		}

		sqlxMasterDb, err = localAuth.Connect(context.TODO())
		if err != nil {
			laughingtale_logger.Error("Unable to create connection to database. Exiting ...")
			panic(err)
		}

	} else {
		// If the RunType is cloud or any other, we will get a connection to given RDS instance
		// using IAM user auth approach
		laughingtale_logger.Info("Using RDS Database Connection [MASTER_NODE].")

		// Information for IAM Authentication RDS connection
		iamAuth := &IAMAuth{
			DatabaseUser:   conf.PostgresSetting.Username,
			DatabaseHost:   conf.PostgresSetting.MasterNodeHost,
			DatabasePort:   conf.PostgresSetting.Port,
			DatabaseName:   conf.PostgresSetting.Database,
			DatabaseSchema: conf.PostgresSetting.TablePrefix,
		}

		sqlxMasterDb, err = iamAuth.Connect(context.TODO())
		if err != nil {
			laughingtale_logger.Error("Unable to create connection to database. Exiting ...")
			panic(err)
		}
	}

	// Perform Migrations
	driver, err = postgres.WithInstance(sqlxMasterDb.DB, &postgres.Config{})
	if err != nil {
		laughingtale_logger.Error("Failed to configure driver for migration ...")
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://resources/db/migrations",
		"postgres", driver)

	if err != nil {
		laughingtale_logger.Error("Failed to perform database migration. Please check if the migration files are present under resources/migrations directory and are error free.")
		panic(err)
	}

	laughingtale_logger.Info("Performing Datbase migrations with migrate.")
	migrationErr := m.Up()
	if migrationErr == nil {
		laughingtale_logger.Info("Changes were detected for migrating.")
		laughingtale_logger.Info("Database migrations performed successfully.")
	} else {
		laughingtale_logger.Info("No changes detected for migrating.")
	}

	gormDB, err := gorm.Open(gormPostgres.New(gormPostgres.Config{
		Conn: sqlxMasterDb,
	}), &gorm.Config{})

	// If the RunType is local, then we will connect the application to local database or the
	// tunneled database
	if conf.ApplicationSetting.RunType == "local" {
		laughingtale_logger.Info("Using Local Database Connection [READER_NODE].")
		readerLocalAuth := &LocalAuth{
			DatabaseUser:     conf.PostgresSetting.Username,
			DatabaseHost:     conf.PostgresSetting.ReaderNodeHost,
			DatabasePort:     conf.PostgresSetting.Port,
			DatabaseName:     conf.PostgresSetting.Database,
			DatabaseSchema:   conf.PostgresSetting.TablePrefix,
			DatabasePassword: conf.PostgresSetting.Password,
		}

		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				gormPostgres.New(gormPostgres.Config{
					DSN:                  readerLocalAuth.GetLocalReaderNodeConnectionString(context.TODO(), net.LookupCNAME),
					PreferSimpleProtocol: true,
				}),
			},
			TraceResolverMode: true,
		}))

	} else {
		laughingtale_logger.Info("Using RDS Database Connection [READER_NODE].")
		readerIamAuth := &IAMAuth{
			DatabaseUser:   conf.PostgresSetting.Username,
			DatabaseHost:   conf.PostgresSetting.ReaderNodeHost,
			DatabasePort:   conf.PostgresSetting.Port,
			DatabaseName:   conf.PostgresSetting.Database,
			DatabaseSchema: conf.PostgresSetting.TablePrefix,
		}
		readerDbConnection, err := readerIamAuth.Connect(context.TODO())
		if err != nil {
			laughingtale_logger.Error("Unable to create connection to database. Exiting ...")
			panic(err)
		}

		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				gormPostgres.New(gormPostgres.Config{
					Conn:                 readerDbConnection,
					PreferSimpleProtocol: true,
				}),
			},
			TraceResolverMode: true,
		}).SetConnMaxIdleTime(1 * time.Minute).
			SetConnMaxLifetime(13 * time.Minute).
			SetMaxIdleConns(2).SetMaxOpenConns(4))
	}

	if err != nil {
		laughingtale_logger.Error("Unable to connect to the database. Exiting ...")
		panic(err)
	}

	laughingtale_logger.Info("Application successfully connected to RDS laughingtale database.")

	// Initialize laughingtaleDb for application
	laughingtaleDb = gormDB
}
