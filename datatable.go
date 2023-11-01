package main

import (
	"errors"

	"github.com/google/uuid"
)

type DataTable struct {
	TableName string
	Columns   *[]DataColumn
	Rows      *[]DataRow
}

type DataColumn struct {
	ColumnName string
	DataType   string
}

func ContainsColumn(cols []DataColumn, colName string) (bool) {
	hasCol := false
	for _, col := range cols {
        if col.ColumnName == colName {
            hasCol = true
        }
    }
    return hasCol
}

type DataRow struct {
	Cells    *[]DataCell
	RowState DataRowState
}

func (row *DataRow) GetCell(colName string) (DataCell, error) {
	var cell DataCell
	var err error
	for i := 0; i < len(*row.Cells); i++ {
		c := (*row.Cells)[i]
		if c.Column.ColumnName == colName {
			cell = c
		}
	}
	if cell.Column == nil {
		err = errors.New("cell not found")
	}
	return cell, err
}

type DataRowState int

const (
	Added     DataRowState = 4
	Deleted   DataRowState = 8
	Detached  DataRowState = 1
	Modified  DataRowState = 16
	Unchanged DataRowState = 2
)

type DataCell struct {
	Column       *DataColumn
	Row          *DataRow
	ptrCellValue interface{}
	cellValue    interface{}
	oldValue     interface{}
}

func (cell *DataCell) GetCellValuePtr() interface{} {
	return cell.ptrCellValue
}

func (cell *DataCell) DerefValue() {
	var val interface{}
	if cell.Column.DataType == "UUID" {
		val = *(cell.ptrCellValue.(*uuid.UUID))
	} else if cell.Column.DataType == "NUMERIC" {
		val = *(cell.ptrCellValue.(*float64))
	} else {
		val = *(cell.ptrCellValue.(*interface{}))
	}
	cell.oldValue = val
	cell.cellValue = val
}

func (cell *DataCell) GetValue() interface{} {
	return cell.cellValue
}

func (cell *DataCell) SetValue(newValue interface{}) {
	cell.cellValue = newValue
	if cell.Row.RowState == Unchanged {
		cell.Row.RowState = Modified
	}
}

func (cell *DataCell) GetOldValue() interface{} {
	return cell.oldValue
}
