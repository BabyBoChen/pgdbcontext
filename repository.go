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

func (repo *DbRepository) Insert(cellValues map[string]interface{}) (map[string]interface{}, error) {

	var lastInsertedId map[string]interface{}
	var err error

	cols := ""
	paramNames := ""
	vals := make([]interface{}, 0)
	colCnt := 0
	for k, v := range cellValues {
		colCnt += 1

		if colCnt > 1 {
			cols += ", "
		}
		cols += "\"" + k + "\""

		if colCnt > 1 {
			paramNames += ", "
		}
		paramNames += fmt.Sprintf("$%d", colCnt)
		vals = append(vals, v)
	}

	idCols := ""
	for i := 0; i < len(*repo.tbModel.Rows); i++ {
		row := (*repo.tbModel.Rows)[i]
		cell, cellErr := row.GetCell("isidentity")
		if cellErr == nil {
			if cell.Value(nil) == true {
				if len(idCols) > 0 {
					idCols += ", "
				}
				idColCell, _ := row.GetCell("fieldname")
				idCols += fmt.Sprintf("\"%s\"", idColCell.Value(nil))
			}
		} else {
			err = cellErr
			break
		}
	}

	var dt *DataTable
	if err == nil {
		cmdTxt := "INSERT INTO \"%s\".\"%s\"(%s) VALUES(%s) RETURNING %s ;"
		cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName, cols, paramNames, idCols)
		dt, err = repo.db.Query(cmdTxt, vals...)
	}

	if err == nil {
		if len(*dt.Rows) > 0 {
			lastInsertedId = make(map[string]interface{})
			row := (*dt.Rows)[0]
			for i := range *row.Cells {
				cell := (*row.Cells)[i]
				lastInsertedId[cell.Column.ColumnName] = cell.Value(nil)
			}
		}
	}

	return lastInsertedId, err
}
