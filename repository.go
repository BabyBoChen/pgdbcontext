package main

import (
	"fmt"
)

type DbRepository struct {
	db        *DbContext
	TableName string
	tbModel   *DataTable
}

func (repo *DbRepository) Select(whereSql string, args ...interface{}) (*DataTable, error) {
	cmdTxt := "SELECT * FROM \"%s\".\"%s\" "
	cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName)
	if len(whereSql) > 0 {
		cmdTxt = fmt.Sprintf("%s WHERE (%s)", cmdTxt, whereSql)
	}
	return repo.db.Query(cmdTxt, args...)
}
