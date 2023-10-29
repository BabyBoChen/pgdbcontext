package data

import "github.com/google/uuid"

type DataTable struct {
	TableName string
	Columns   *[]DataColumn
	Rows      *[]DataRow
}

type DataColumn struct {
	ColumnName string
	DataType   string
}

type DataRow struct {
	Cells    *[]DataCell
	RowState DataRowState
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

func (cell *DataCell) Value(newValue interface{}) interface{} {
	if newValue != nil {
		cell.cellValue = newValue
		if cell.Row.RowState == Unchanged {
			cell.Row.RowState = Modified
		}
	}
	return cell.cellValue
}

func (cell *DataCell) GetOldValue() interface{} {
	return cell.oldValue
}
