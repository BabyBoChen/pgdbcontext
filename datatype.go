package main

import (
	"github.com/google/uuid"
)

func CreateEmptyCell(colName string, dbType string) *DataCell {
	var cell DataCell
	var col DataColumn
	col.ColumnName = colName
	col.DataType = dbType
	cell.Column = &col
	var ptr interface{}
	if dbType == "UUID" {
		var u uuid.UUID
		ptr = &u
	} else if dbType == "NUMERIC" {
		var f float64
		ptr = &f
	} else {
		var obj interface{}
		ptr = &obj
	}
	cell.ptrCellValue = ptr
	return &cell
}
