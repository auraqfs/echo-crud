package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"echo-crud/entity"
	"echo-crud/internal/config"
)

var db *gorm.DB

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config, err := config.NewConfig(".env")
	checkError(err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.Name)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	checkError(err)

	sqlDB, err := db.DB()
	defer sqlDB.Close()

	executePendingMigrations(db)

	// Migrate rest of the models
	log.Info().Msg("AutoMigrate Model [table_name]")

	log.Info().Msg("  SupplierModel [" + (&entity.Supplier{}).TableName() + "]")
	db.AutoMigrate(&entity.Supplier{})

	log.Info().Msg("  ProdukModel [" + (&entity.Produk{}).TableName() + "]")
	db.AutoMigrate(&entity.Produk{})

	log.Info().Msg("  PelangganModel [" + (&entity.Pelanggan{}).TableName() + "]")
	db.AutoMigrate(&entity.Pelanggan{})

	log.Info().Msg("  TransaksiModel [" + (&entity.Transaksi{}).TableName() + "]")
	db.AutoMigrate(&entity.Transaksi{})

	log.Info().Msg("  TransaksiDetailModel [" + (&entity.TransaksiDetail{}).TableName() + "]")
	db.AutoMigrate(&entity.TransaksiDetail{})

	log.Info().Msg("  PembayaranModel [" + (&entity.Pembayaran{}).TableName() + "]")
	db.AutoMigrate(&entity.Pembayaran{})

}

func executePendingMigrations(db *gorm.DB) {
	db.AutoMigrate(&MigrationHistoryModel{})
	lastMigration := MigrationHistoryModel{}
	skipMigration := db.Order("migration_id desc").Limit(1).Find(&lastMigration).RowsAffected > 0

	// skip to last migration
	keys := make([]string, 0, len(migrations))
	for k := range migrations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// run all migrations in one transaction
	if len(migrations) == 0 {
		log.Info().Msg("No pending migrations")
	} else {
		db.Transaction(func(tx *gorm.DB) error {
			for _, k := range keys {
				if skipMigration {
					if k == lastMigration.MigrationID {
						skipMigration = false
					}
				} else {
					log.Info().Msg("  " + k)
					tx.Transaction(func(subTx *gorm.DB) error {
						// run migration update
						checkError(migrations[k](subTx))
						// insert migration id into history
						checkError(subTx.Create(MigrationHistoryModel{MigrationID: k}).Error)
						return nil
					})
				}
			}
			return nil
		})
	}
}

type mFunc func(tx *gorm.DB) error

var migrations = make(map[string]mFunc)

// MigrationHistoryModel model migration
type MigrationHistoryModel struct {
	MigrationID string `gorm:"type:text;primaryKey"`
}

// TableName name of migration table
func (model *MigrationHistoryModel) TableName() string {
	return "migration_history"
}

func checkError(err error) {
	if err != nil {
		log.Fatal().Err(err)
		panic(err)
	}
}

func registerMigration(id string, fm mFunc) {
	migrations[id] = fm
}
