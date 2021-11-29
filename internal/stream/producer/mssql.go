package producer

import (
	"database/sql"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const (
	mssqlProducerName = "mssql"
	mssqlScheme       = "sqlserver"
)

const (
	mssqlHostCName     = "host"
	mssqlInstanceCName = "instance"
	mssqlDatabaseCName = "database"
	mssqlUserCName     = "username"
	mssqlPasswordCName = "password"
	mssqlQueryCName    = "query"
	mssqlParamsCName   = "params"

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
	for rows.Next() {

	}
}

type mssqlConfig struct {
	connURL *url.URL
	query   string
	mapping interface{}
}

func processConfig(config *util.ConfigMap) (mssqlConfig, error) {
	host := config.Get(mssqlHostCName).Required().String()
	instance := config.Get(mssqlInstanceCName).Required().String()
	database := config.Get(mssqlDatabaseCName).Required().String()
	username := config.Get(mssqlUserCName).String()
	password := config.Get(mssqlPasswordCName).String()
	query := config.Get(mssqlQueryCName).Required().String()
	connParams := config.Get(mssqlParamsCName).Map()
	if err := config.Error(); err != nil {
		return mssqlConfig{}, err
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
