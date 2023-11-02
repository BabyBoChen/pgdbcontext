package main

import (
	"database/sql"
	"errors"
)

type DbContext struct {
	connStr       string
	Conn          *sql.DB
	Transaction   *sql.Tx
	DefaultSchema string
}

func NewDbContext(connStr string) (*DbContext, error) {
	var ptrDb *DbContext
	conn, err := sql.Open("postgres", connStr)
	if err == nil {
		var db DbContext
		db.connStr = connStr
		db.Conn = conn
		db.DefaultSchema = "public"
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

func (db *DbContext) Query(cmdTxt string, args ...interface{}) (*DataTable, error) {

	var dt *DataTable
	var err error

	err = db.beginTransction()

	var rows *sql.Rows
	if err == nil {
		rows, err = db.Transaction.Query(cmdTxt, args...)
	}

	var dataTable DataTable
	if err == nil {
		var dataCols []*DataColumn
		colTypes, _ := rows.ColumnTypes()
		colNames, _ := rows.Columns()
		dataCols = make([]*DataColumn, len(colNames))
		for i := 0; i < len(colNames); i++ {
			var col DataColumn
			col.ColumnName = colNames[i]
			col.DataType = colTypes[i].DatabaseTypeName()
			dataCols[i] = &col
		}
		dataTable.Columns = dataCols

		dataTable.Rows = make([]*DataRow, 0)
		for rows.Next() {
			var dataRow DataRow
			dataRow.RowState = Unchanged
			cells := make([]*DataCell, len(colTypes))
			cellValuePtrs := make([]interface{}, len(colTypes))
			for i := 0; i < len(cells); i++ {
				cell := CreateEmptyCell(colNames[i], colTypes[i].DatabaseTypeName())
				cells[i] = cell
				cellValuePtrs[i] = cell.GetCellValuePtr()
			}
			dataRow.Cells = cells
			err = rows.Scan(cellValuePtrs...)

			if err == nil {
				for i := 0; i < len(cells); i++ {
					cells[i].Row = &dataRow
					cells[i].DerefValue()
				}
				dataTable.Rows = append(dataTable.Rows, &dataRow)
			} else {
				break
			}
		}
	}

	if err == nil {
		dt = &dataTable
	}

	return dt, err
}

const SPGetTbFldInfos string = `SELECT A.column_name::varchar as fieldname
,pg_catalog.col_description(E.oid,A.ordinal_position) as shortdesc
,A.udt_name::varchar AS datatype
,coalesce(A.character_maximum_length,0)+coalesce(A.numeric_precision,0) as datalength
,coalesce(A.numeric_scale,0) as numericscale 
,CASE WHEN A.is_nullable='YES'
	THEN true
	ELSE false
	END AS isallownull
,CASE WHEN A.column_default IS NOT NULL AND A.is_identity='NO'
	THEN A.column_default
	ELSE 'NULL'
	END AS defaultvalue
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

func (db *DbContext) GetRepository(tableName string) (*DbRepository, error) {
	var repoPtr *DbRepository
	var repo DbRepository
	repo.TableName = tableName
	repo.db = db
	schema, err := db.Query(SPGetTbFldInfos, db.DefaultSchema, tableName)
	if err == nil {
		repo.tbModel = schema
		repoPtr = &repo
	}
	return repoPtr, err
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
