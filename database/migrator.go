package database

import (
	"log"

	"github.com/DavidHuie/gomigrate"
)

func (d *Database) migrate(migrationsDir string) error {
	if migrator, err := gomigrate.NewMigratorWithLogger(d.db, gomigrate.Sqlite3{}, migrationsDir, log.Default()); err != nil {
		return err
	} else {
		return migrator.Migrate()
	}
}
