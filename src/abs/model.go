package abs

import (
	"fmt"
	app "src/config"
	"src/library/tool"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dialect             = "charset=utf8mb4&loc=PRC&parseTime=true&maxAllowedPacket=0&timeout=10s"
	retryConnDbMaxTimes = 3
)

type Model struct {
}

type dbCollector struct {
	once    sync.Once
	pointer *gorm.DB
}

var et, bbs dbCollector

func (m *Model) getClient(dbc *dbCollector, db string) *gorm.DB {
	var (
		err   error
		fatal bool
	)
	for i := 1; i <= retryConnDbMaxTimes; i++ {
		dbc.once.Do(func() {
			if dbc.pointer, err = m.openMysql(db); err == nil {
				m.setDefaultConfig(dbc.pointer)
			}
		})
		if err != nil {
			if i < retryConnDbMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				dbc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		if err = dbc.pointer.DB().Ping(); err != nil {
			if i < retryConnDbMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				dbc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		break
	}
	if fatal {
		defer func() {
			errMsg := fmt.Sprintf("get db client of %s occur error: %s", db, err.Error())
			app.Logger.Error(errMsg)
			tool.PushSimpleMessage(errMsg, true)
		}()
		return nil
	}
	return dbc.pointer
}

func (m *Model) openMysql(db string) (*gorm.DB, error) {
	dt, ok := app.DB[db]
	if !ok {
		return nil, fmt.Errorf("dialect of %s is not exist", db)
	}
	tcp := dt.User + ":" + dt.Pwd + "@tcp(" + dt.Host + ":" + dt.Port + ")/" + dt.Db + "?" + dialect
	return gorm.Open("mysql", tcp)
}

func (m *Model) setDefaultConfig(db *gorm.DB) {
	db.DB().SetMaxOpenConns(8)
	db.DB().SetMaxIdleConns(2)
	db.DB().SetConnMaxLifetime(time.Second * 30)
}

func (m *Model) SetMaxOpenConns(db *gorm.DB, num int) {
	if num > 0 && num <= 1000 {
		db.DB().SetMaxOpenConns(num)
	}
}

func (m *Model) GetEnjoyThin() *gorm.DB {
	return m.getClient(&et, "et")
}

func (m *Model) GetBbs() *gorm.DB {
	return m.getClient(&bbs, "bbs")
}
