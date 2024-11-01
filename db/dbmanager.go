package db

import (
	"github.com/yhhaiua/engine/util"
	"gorm.io/gorm"
)

var storage map[string]*util.AtomicLong

// Init 初始化数据库数据
func Init(config *MysqlConfig, args ...interface{}) {
	//连接数据库
	if !accessor.LoadConfig(config) {
		return
	}
	//自动映射表
	accessor.AutoMigrate(args...)
	//开启数据保存缓存
	globalDbServer.InitDBCache()
	//初始化自增数据
	initStorage(args...)
	//显示自动创建的数据表
	logShowTable(args...)
}

func logShowTable(args ...interface{}) {
	var table []string
	for _, v := range args {
		ie := v.(IEntity)
		table = append(table, ie.TableName())
	}
	if len(table) != 0 {
		logger.Infof("自动创建的数据表:%v", table)
	}
}

// Create 创建数据库表结构 主键id存在
func Create(entity IEntity) {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s", entity.TableName())
		return
	}
	globalDbServer.AddDBCache(save, entity)
}

// CreateElement 创建数据库表结构 主键id存在
func CreateElement(entity CElement) {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s", entity.TableName())
		return
	}
	if isDelay, ok := entity.(IDelayBefore); ok && isDelay.IsDelayBefore() {
		globalDbServer.AddDBCacheElement(save, entity)
	} else {
		entity.Before()
		globalDbServer.AddDBCache(save, entity.GetEntity())
	}
}

// CreateAuto 创建数据库表结构 主键id自动生成
func CreateAuto(entity IEntity, operateId, serverId int) {
	if entity.GetId() == 0 {
		at := storage[entity.TableName()]
		id := util.BuildPrimaryKey(operateId, serverId, at.IncrementAndGet())
		entity.SetId(id)
	}
	globalDbServer.AddDBCache(save, entity)
}

// GetDefaultId 获取自增id
func GetDefaultId(entity IEntity, operateId, serverId int) int64 {
	at := storage[entity.TableName()]
	id := util.BuildPrimaryKey(operateId, serverId, at.IncrementAndGet())
	return id
}

// Update 更新表数据
func Update(entity IEntity) {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s", entity.TableName())
		return
	}
	globalDbServer.AddDBCache(update, entity)
}

// UpdateElement 更新表数据
func UpdateElement(entity CElement) {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s", entity.TableName())
		return
	}
	if isDelay, ok := entity.(IDelayBefore); ok && isDelay.IsDelayBefore() {
		globalDbServer.AddDBCacheElement(update, entity)
	} else {
		entity.Before()
		globalDbServer.AddDBCache(update, entity.GetEntity())
	}
}

// RepairSave 修复保存
func RepairSave(entity IEntity) {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s", entity.TableName())
		return
	}
	accessor.Save(entity)
}

// Delete 删除表数据
func Delete(entity IEntity) {
	globalDbServer.AddDBCache(remove, entity)
}

// First 查找单条数据
func First(id interface{}, entity interface{}) bool {
	return accessor.First(id, entity)
}

func GetDB() *gorm.DB {
	return accessor.db
}

// FindAll 查找所有数据 entity 切片
func FindAll(entity interface{}) {
	accessor.FindAll(entity)
}

// FindCond 条件查询 dest 切片  //accessor.db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
func FindCond(dest interface{}, query interface{}, args ...interface{}) {
	accessor.FindCond(dest, query, args...)
}

// Close 关闭数据库连接
func Close() {
	accessor.Close()
}

// ShutDown 缓存数据关闭
func ShutDown() {
	globalDbServer.ShutDown()
	accessor.Close()
}
func initStorage(args ...interface{}) {
	storage = make(map[string]*util.AtomicLong)
	for _, v := range args {
		ie := v.(IEntity)
		if is, ok := ie.(ISince); !ok || is.GetSince() == false {
			continue
		}
		at := &util.AtomicLong{}
		maxId := GetMaxId(ie)
		if maxId > 0 {
			at.Set(util.ParseAutoincrement(maxId))
		} else {
			at.Set(0)
		}

		storage[ie.TableName()] = at
		logger.Infof("自增id的表:%s,maxId:%d", ie.TableName(), maxId)
	}
}

func GetMaxId(entity IEntity) int64 {
	hql := "select max(id) from " + entity.TableName()
	var maxId int64
	var count int64
	result := accessor.db.Table(entity.TableName()).Count(&count)
	if result.Error != nil {
		logger.Errorf("Count error:%s", result.Error.Error())
	}
	if count == 0 {
		return 0
	}
	result = result.Raw(hql).Scan(&maxId)
	if result.Error != nil {
		logger.Errorf("getMaxId error:%s", result.Error.Error())
	}
	return maxId
}
