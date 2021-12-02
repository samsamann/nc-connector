package producer

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const (
	mssqlProducerName = "mssql"
	mssqlScheme       = "sqlserver"
)

const (
	mssqlHostCName           = "host"
	mssqlInstanceCName       = "instance"
	mssqlDatabaseCName       = "database"
	mssqlUserCName           = "username"
	mssqlPasswordCName       = "password"
	mssqlQueryCName          = "query"
	mssqlParamsCName         = "params"
	mssqlMappingCName        = "mapping"
	mssqlMappingContentCName = "content"
	mssqlMappingPathCName    = "path"

	mssqlAppNameParam   = "app name"
	mssqlDefaultAppName = "nc-connector"
)

func initMssqlProducer(config map[string]interface{}) (stream.Producer, error) {
	c, err := processConfig(util.NewConfigMap(config))
	if err != nil {
		return nil, err
	}
	return newMssqlProducer(c), nil
}

type mssqlProducer struct {
	config mssqlConfig
}

func newMssqlProducer(c mssqlConfig) stream.Producer {
	p := new(mssqlProducer)
	p.config = c
	return p
}

func (ms mssqlProducer) Out() <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)

	//TODO: handle errors
	db, _ := mssqlConnect(mssqlScheme, ms.config.connURL)

	rows, _ := db.Query(ms.config.query)
	go func() {
		defer func() {
			close(channel)
			if db != nil {
				db.Close()
			}
		}()
		ms.process(channel, rows)
	}()

	return channel
}

func (ms mssqlProducer) process(channel chan<- stream.SyncItem, rows *sql.Rows) {
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			rows.Close()
			return
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			if str, ok := (*val).(string); ok {
				*val = strings.Trim(str, " ")
			}
			m[colName] = *val
		}
		// TODO: ckeck error
		path := m[ms.config.mapping[mssqlMappingPathCName]].(string)
		content := m[ms.config.mapping[mssqlMappingContentCName]].([]byte)
		delete(m, ms.config.mapping[mssqlMappingPathCName])
		delete(m, ms.config.mapping[mssqlMappingContentCName])

		props := make(stream.Properties)
		for k, v := range m {
			props[k] = v
		}

		channel <- stream.NewFileSyncItem(path, props, content)
	}
}

type mssqlConfig struct {
	connURL *url.URL
	query   string
	mapping map[string]string
}

func processConfig(config *util.ConfigMap) (mssqlConfig, error) {
	host := config.Get(mssqlHostCName).Required().String()
	instance := config.Get(mssqlInstanceCName).Required().String()
	database := config.Get(mssqlDatabaseCName).Required().String()
	username := config.Get(mssqlUserCName).String()
	password := config.Get(mssqlPasswordCName).String()
	query := config.Get(mssqlQueryCName).Required().String()
	connParams := config.Get(mssqlParamsCName).Map()
	mapping := config.Get(mssqlMappingCName).Required().Map()
	if err := config.Error(); err != nil {
		return mssqlConfig{}, err
	}

	p, pathExists := mapping[mssqlMappingPathCName]
	c, contentExists := mapping[mssqlMappingContentCName]
	if !pathExists || len(p) == 0 || !contentExists || len(c) == 0 {
		return mssqlConfig{},
			fmt.Errorf(
				"%q and %q do not exist or must not be empty in the mapping conf",
				mssqlMappingPathCName,
				mssqlMappingContentCName,
			)
	}

	urlQuery := make(url.Values)
	urlQuery.Add(mssqlDatabaseCName, database)
	for k, v := range connParams {
		urlQuery.Add(k, v)
	}
	if !urlQuery.Has(mssqlAppNameParam) {
		urlQuery.Add(mssqlAppNameParam, mssqlDefaultAppName)
	}
	u := &url.URL{
		Scheme:   mssqlScheme,
		User:     url.UserPassword(username, password),
		Host:     host,
		Path:     instance,
		RawQuery: urlQuery.Encode(),
	}

	return mssqlConfig{
		connURL: u,
		query:   query,
		mapping: mapping,
	}, nil
}

func mssqlConnect(driver string, u *url.URL) (db *sql.DB, err error) {
	db, err = sql.Open(driver, u.String())
	if err != nil {
		return
	}
	err = db.Ping()
	return
}
