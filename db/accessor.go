package db

import (
	"github.com/yhhaiua/engine/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

var logger = log.GetLogger()

type Accessor struct {
	db *gorm.DB
}

func (a *Accessor) Db() *gorm.DB {
	return a.db
}

var accessor = &Accessor{}

func (a *Accessor) LoadConfig(config *MysqlConfig) bool {

	var defaultSize uint = 191
	if strings.Contains(config.Config, "utf8mb4") {
		defaultSize = 191
	} else if strings.Contains(config.Config, "utf8") {
		defaultSize = 255
	}
	mysqlConfig := mysql.Config{
		DSN:               config.Dsn(), // DSN data source name
		DefaultStringSize: defaultSize,  // string 类型字段的默认长度
	}
	db, err := gorm.Open(mysql.New(mysqlConfig),
		&gorm.Config{Logger: replaceLog})

	if err != nil {
		logger.Errorf("mysql connect error:%s", err.Error())
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("mysql connect error:%s", err.Error())
		return false
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	sqlDB.SetConnMaxLifetime(2 * time.Hour)
	a.db = db
	return true
}

func (a *Accessor) Create(entity interface{}) {
	result := a.db.Create(entity)
	if result.Error != nil {
		logger.Errorf("mysql Create error:%s", result.Error.Error())
	}
}

func (a *Accessor) Save(entity interface{}) {
	result := a.db.Save(entity)
	if result.Error != nil {
		logger.Errorf("mysql Save error:%s", result.Error.Error())
	}
}

func (a *Accessor) Delete(id interface{}, entity interface{}) {
	result := a.db.Delete(entity, id)
	if result.Error != nil {
		logger.Errorf("mysql Delete error:%s", result.Error.Error())
	}
}

func (a *Accessor) AutoMigrate(args ...interface{}) {

	if err := a.db.AutoMigrate(args...); err != nil {
		logger.Errorf("mysql AutoMigrate error:%s", err.Error())
	}
}

func (a *Accessor) First(id interface{}, entity interface{}) bool {
	result := a.db.First(entity, id)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logger.Errorf("mysql First error:%s", result.Error.Error())
			//panic(result.Error.Error())
		}
		return false
	}
	return true
}

func (a *Accessor) FindAll(entity interface{}) {
	result := a.db.Find(entity)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logger.Errorf("mysql FindAll error:%s", result.Error.Error())
			//panic(result.Error.Error())
		}
	}
}

// FindCond 条件查询 entity 切片  //accessor.db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
func (a *Accessor) FindCond(dest interface{}, query interface{}, args ...interface{}) {
	result := a.db.Where(query, args...).Find(dest)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logger.Errorf("mysql FindCond error:%s", result.Error.Error())
			//panic(result.Error.Error())
		}
	}
}

func (a *Accessor) GetTables() []string {
	tableList, _ := a.db.Migrator().GetTables()
	return tableList
}

func (a *Accessor) Close() {
	db, _ := a.db.DB()
	if db != nil {
		db.Close()
	}
}
