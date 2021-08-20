package db

import (
	"fmt"
	"log"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type MysqlConn struct {
	Conn *gorm.DB
}

func New() *MysqlConn {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}

	return &MysqlConn{
		Conn: dbConn,
	}
}
