package database

import (
	"fmt"
	"os"

	"github.com/48Club/ip-waf-helper/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type server struct {
	*gorm.DB
}

var Server = server{}

func init() {
	engine, err := gorm.Open(mysql.Open(gethDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = engine.AutoMigrate(&types.IPWaf{})
	if err != nil {
		panic(err)
	}
	Server.DB = engine
}

func gethDSN() string {
	return fmt.Sprintf("root:%s@tcp(mysql:3306)/ip_waf_helper", os.Getenv("MARIADB_ROOT_PASSWORD"))
}
