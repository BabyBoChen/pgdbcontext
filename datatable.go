package main

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

const SPGetTbFldInfos string = `SELECT A.column_name::varchar as fieldname
,pg_catalog.col_description(E.oid,A.ordinal_position) as shortdesc
,coalesce(A.character_maximum_length,0)+coalesce(A.numeric_precision,0) as datalength
,coalesce(A.numeric_scale,0) as numericscale 
,CASE WHEN A.is_nullable='YES'
	THEN true
	ELSE false
	END AS isallownull
,A.column_default as defaultvalue
,CASE WHEN constraint_type='PRIMARY KEY'
	THEN true
	ELSE false
	END AS isprimarykey
,case when A.is_identity='YES'
	then true
	else false
	end as isidentity
FROM information_schema.columns AS A
LEFT JOIN information_schema.constraint_column_usage AS B ON A.table_schema=B.table_schema AND A.table_name=B.table_name AND A.column_name=B.column_name
LEFT JOIN information_schema.table_constraints AS C ON B.table_schema=C.table_schema AND B.table_name=C.table_name AND B.constraint_name=C.constraint_name
LEFT JOIN pg_catalog.pg_namespace AS D ON D.nspname=A.table_schema
LEFT JOIN pg_catalog.pg_class AS E ON E.relnamespace=D.oid AND E.relname=A.table_name
WHERE A.table_schema=$1 AND A.table_name=$2;`
