package model

import "fmt"
import "xorm.io/xorm"
import _ "github.com/lib/pq"
import _ "github.com/go-sql-driver/mysql"

type Config struct {
	Name        string           `json:"name"`
	DriverName  string           `json:"driver_name"`
	DataSources []string         `json:"data_sources"`
	GroupPolicy xorm.GroupPolicy `json:"group_policy"`
}

type Database struct {
	*xorm.EngineGroup
}

var dbs = make(map[string]*Database)

func Init(cfgs ...*Config) error {
	if len(cfgs) == 0 {
		return fmt.Errorf("db config can not be nil")
	}
	for _, cfg := range cfgs {
		engineGroup, err := xorm.NewEngineGroup(cfg.DriverName, cfg.DataSources, cfg.GroupPolicy)
		if err != nil {
			return fmt.Errorf("init engine group err: %s; cfg: %q", err, cfg)
		}
		dbs[cfg.Name] = &Database{engineGroup}
	}
	return nil
}

func Get(name string) (*Database, error) {
	db, ok := dbs[name]
	if !ok {
		return nil, fmt.Errorf("db %s not found", name)
	}
	return db, nil
}
