package testutil

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
)

var dbServer *server.Server

func InitDB() {
	if dbServer != nil {
		return
	}
	db := memory.NewDatabase("ideashare")
	provider := memory.NewDBProvider(db)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	s, err := server.NewServer(
		server.Config{Protocol: "tcp", Address: ":3319"},
		sqle.NewDefault(provider),
		memory.NewSessionBuilder(provider),
		nil,
	)
	if err != nil {
		panic(err)
	}
	dbServer = s
	go func() {
		if err = dbServer.Start(); err != nil {
			panic(err)
		}
	}()
}

func TearDownDB() {
	_ = dbServer.Close()
}
