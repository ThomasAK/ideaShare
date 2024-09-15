package testutil

import (
	"fmt"
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

func PrintTableContents(table string) {
	fmt.Println("Printing contents of table " + table)
	db, _ := Container.Db.DB()
	rows, _ := db.Query("select * from " + table)
	cols, _ := rows.Columns()
	row := make([]interface{}, len(cols))
	rowPtr := make([]interface{}, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}
	fmt.Println(cols)
	for rows.Next() {
		rows.Scan(rowPtr...)
		for i, col := range row {
			if col == nil {
				row[i] = "NULL"
			} else {
				switch v := col.(type) {
				case []byte:
					row[i] = string(v)
				}
			}
		}
		fmt.Println(row...)
	}
}
