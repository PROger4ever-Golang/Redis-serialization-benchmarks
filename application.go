package main

import (
	"fmt"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/implementations"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
)

type Application struct {
	conf             *Config
	commonConnection redis.Conn
	rejsonConnection redis.Conn

	ErrorLogger *log.Logger

	MapImplementation    implementations.Implementation
	FlatImplementation   implementations.Implementation
	JsonImplementation   implementations.Implementation
	GobImplementation    implementations.Implementation
	RejsonImplementation implementations.Implementation

	Implementations map[string]implementations.Implementation
}

func (app *Application) Finalize() {
	for _, imp := range app.Implementations {
		imp.Del()
	}
	app.commonConnection.Close()
	app.rejsonConnection.Close()
}

func NewApplication() (app *Application, err error) {
	app = &Application{
		ErrorLogger: log.New(os.Stderr, "", 0),
	}

	app.conf, err = LoadConfig("config.json", "config-custom.json")
	if err != nil {
		return nil, utils.WrapIfError(err, "NewApplication->LoadConfig()")
	}

	commonAddress := fmt.Sprintf("%s:%d", app.conf.CommonRedis.Host, app.conf.CommonRedis.Port)
	app.commonConnection, err = redis.Dial("tcp", commonAddress)
	if err != nil {
		return nil, utils.WrapIfError(err, "NewApplication->redis.Dial(commonAddress)")
	}
	rejsonAddress := fmt.Sprintf("%s:%d", app.conf.RejsonRedis.Host, app.conf.RejsonRedis.Port)
	app.rejsonConnection, err = redis.Dial("tcp", rejsonAddress)
	if err != nil {
		return nil, utils.WrapIfError(err, "NewApplication->redis.Dial(rejsonAddress)")
	}

	app.MapImplementation = implementations.NewMapImplementation(app.commonConnection)
	app.FlatImplementation = implementations.NewFlatImplementation(app.commonConnection)
	app.JsonImplementation = implementations.NewJsonImplementation(app.commonConnection)
	app.GobImplementation = implementations.NewGobImplementation(app.commonConnection)
	app.RejsonImplementation = implementations.NewRejsonImplementation(app.rejsonConnection)

	app.Implementations = map[string]implementations.Implementation{
		"Map":    app.MapImplementation,
		"Flat":   app.FlatImplementation,
		"Json":   app.JsonImplementation,
		"Gob":    app.GobImplementation,
		"Rejson": app.RejsonImplementation,
	}

	return app, nil
}
