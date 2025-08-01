package command

import (
	"log"
	"os"
	"strings"

	"erajaya-interview/constants"

	"erajaya-interview/migrations"

	"erajaya-interview/script"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func Commands(injector *do.Injector) bool {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	var scriptName string

	migrate := false
	seed := false
	run := false
	scriptFlag := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}
		if arg == "--seed" {
			seed = true
		}
		if arg == "--run" {
			run = true
		}
		if strings.HasPrefix(arg, "--script:") {
			scriptFlag = true
			scriptName = strings.TrimPrefix(arg, "--script:")
		}
	}

	if migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("error migration: %v", err)
		}
		log.Println("migration completed successfully")
	}

	if seed {
		if err := migrations.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

	if scriptFlag {
		if err := script.Script(scriptName, db); err != nil {
			log.Fatalf("error script: %v", err)
		}
		log.Println("script run successfully")
	}

	if run {
		return true
	}

	return false
}
