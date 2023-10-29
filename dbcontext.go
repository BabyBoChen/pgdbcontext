package main

import (
	"database/sql"
	"errors"
	"pgsql/data"
)

type DbContext struct {
	connStr     string
	Conn        *sql.DB
	Transaction *sql.Tx
}

func New(connStr string) (*DbContext, error) {
	var ptrDb *DbContext
	conn, err := sql.Open("postgres", connStr)
	if err == nil {
		var db DbContext
		db.connStr = connStr
		db.Conn = conn
		ptrDb = &db
	} else {
		ptrDb = nil
	}
	return ptrDb, err
}

func (db *DbContext) beginTransction() error {
	var err error
	if db.Conn != nil {
		if db.Transaction == nil {
			db.Transaction, err = db.Conn.Begin()
		}
	}
	return err
}

func (db *DbContext) Query(cmdTxt string, args ...interface{}) (*data.DataTable, error) {

	var dt *data.DataTable
	var err error

	db.beginTransction()
	var rows *sql.Rows
	rows, err = db.Transaction.Query(cmdTxt, args...)
	if err == nil {
		var dataTable data.DataTable
		dataRows := make([]data.DataRow, 0)
		dataTable.Rows = &dataRows

		for rows.Next() {
			colTypes, _ := rows.ColumnTypes()
			colNames, _ := rows.Columns()

			var dataRow data.DataRow
			cells := make([]data.DataCell, len(colTypes))
			cellValuePtrs := make([]interface{}, len(colTypes))
			dataRow.Cells = &cells
			dataRow.RowState = data.Unchanged

			for i := 0; i < len(cells); i++ {
				cell := data.CreateEmptyCell(colNames[i], colTypes[i].DatabaseTypeName())
				cells[i] = *cell
				cellValuePtrs[i] = cell.GetCellValuePtr()
			}
			err = rows.Scan(cellValuePtrs...)

			if err == nil {
				for i := 0; i < len(cells); i++ {
					cells[i].Row = &dataRow
					cells[i].DerefValue()
				}
				r := *dataTable.Rows
				r = append(r, dataRow)
				dataTable.Rows = &r
			}
		}
		dt = &dataTable
	}

	return dt, err
}

func (db *DbContext) Commit() error {
	if db.Transaction != nil {
		err := db.Transaction.Commit()
		if err == nil {
			db.Transaction = nil
		}
		return err
	} else {
		return errors.New("transaction is nil")
	}
}

func (db *DbContext) Dispose() {
	if db.Transaction != nil {
		db.Transaction.Rollback()
		db.Transaction = nil
	}
	if db.Conn != nil {
		db.Conn.Close()
	}
}
