package main

import (
	"errors"
	"fmt"
	"pgsql/utils"
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

func (repo *DbRepository) NewRow() (map[string]interface{}, error) {

	newRow := make(map[string]interface{})
	var err error

	selectSql := "SELECT "
	needComma := false
	for _, fld := range repo.tbModel.Rows {
		var fnCell *DataCell
		var dfCell *DataCell
		fnCell, err = fld.GetCell("fieldname")
		if err == nil {
			dfCell, err = fld.GetCell("defaultvalue")
		}
		if err == nil {
			if needComma {
				selectSql += ", "
			}
			selectSql += fmt.Sprintf("%s AS \"%s\"", dfCell.GetValue(), fnCell.GetValue())
			needComma = true
		} else {
			break
		}
	}

	var dt *DataTable
	if err == nil {
		dt, err = repo.db.Query(selectSql)
	}

	if err == nil {
		if len(dt.Rows) != 1 {
			err = errors.New("unable to create new row")
		}
	}

	if err == nil {
		row := dt.Rows[0]
		newRow = row.ToMap()
	}

	return newRow, err
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
	for i := 0; i < len(repo.tbModel.Rows); i++ {
		row := repo.tbModel.Rows[i]
		cell, cellErr := row.GetCell("isidentity")
		if cellErr == nil {
			if cell.GetValue() == true {
				if len(idCols) > 0 {
					idCols += ", "
				}
				idColCell, _ := row.GetCell("fieldname")
				idCols += fmt.Sprintf("\"%s\"", idColCell.GetValue())
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
		if len(dt.Rows) > 0 {
			lastInsertedId = make(map[string]interface{})
			row := dt.Rows[0]
			for i := range row.Cells {
				cell := row.Cells[i]
				lastInsertedId[cell.Column.ColumnName] = cell.GetValue()
			}
		}
	}

	return lastInsertedId, err
}

func (repo *DbRepository) Update(cellValues map[string]interface{}) error {
	var err error

	var pkCol DataColumn
	pkCol, err = repo.getPrimaryKeyColumn()

	if err == nil {
		if !utils.HasKey(cellValues, pkCol.ColumnName) {
			err = errors.New("primary key column was not found in provided cellvalues")
		}
	}

	var idCols []DataColumn
	if err == nil {
		idCols, err = repo.getIdentityColumns()
	}

	setters := ""
	colCnt := 0
	vals := make([]interface{}, 0)
	if err == nil {
		for colName, val := range cellValues {
			if colName != pkCol.ColumnName && !ContainsColumn(idCols, colName) {
				colCnt += 1
				if colCnt > 1 {
					setters += ", "
				}
				setters += fmt.Sprintf("\"%s\"=$%d", colName, colCnt)
				vals = append(vals, val)
			}
		}
		if len(setters) == 0 {
			err = errors.New("no set clause in this operation")
		}
	}

	whereSql := ""
	if err == nil {
		colCnt += 1
		whereSql = fmt.Sprintf("\"%s\"=$%d", pkCol.ColumnName, colCnt)
		vals = append(vals, cellValues[pkCol.ColumnName])

		cmdTxt := "UPDATE \"%s\".\"%s\" SET %s WHERE %s"
		cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName, setters, whereSql)
		_, err = repo.db.Query(cmdTxt, vals...)
	}

	return err
}

func (repo *DbRepository) UpdateRow(row DataRow) error {
	var err error

	var pkCol DataColumn
	pkCol, err = repo.getPrimaryKeyColumn()

	var idCols []DataColumn
	if err == nil {
		idCols, err = repo.getIdentityColumns()
	}

	setters := ""
	colCnt := 0
	vals := make([]interface{}, 0)
	if err == nil {
		for _, cell := range row.Cells {
			if cell.Column.ColumnName != pkCol.ColumnName && !ContainsColumn(idCols, cell.Column.ColumnName) {
				colCnt += 1
				if colCnt > 1 {
					setters += ", "
				}
				setters += fmt.Sprintf("\"%s\"=$%d", cell.Column.ColumnName, colCnt)
				vals = append(vals, cell.GetValue())
			}
		}
		if len(setters) == 0 {
			err = errors.New("no set clause in this operation")
		}
	}

	var pkCell *DataCell
	if err == nil {
		pkCell, err = row.GetCell(pkCol.ColumnName)
	}

	whereSql := ""
	if err == nil {
		colCnt += 1
		whereSql = fmt.Sprintf("\"%s\"=$%d", pkCol.ColumnName, colCnt)
		vals = append(vals, pkCell.GetValue())
		cmdTxt := "UPDATE \"%s\".\"%s\" SET %s WHERE %s"
		cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName, setters, whereSql)
		_, err = repo.db.Query(cmdTxt, vals...)
	}

	return err
}

func (repo *DbRepository) getIdentityColumns() ([]DataColumn, error) {

	idCols := make([]DataColumn, 0)
	var err error

	for _, row := range repo.tbModel.Rows {
		var idCell *DataCell
		idCell, err = row.GetCell("isidentity")

		if err == nil && idCell.GetValue() == true {
			var idCol DataColumn

			var fieldNameCell *DataCell
			var dataTypeCell *DataCell

			fieldNameCell, err = row.GetCell("fieldname")
			if err == nil {
				dataTypeCell, err = row.GetCell("datatype")
			}
			if err == nil {
				idCol.ColumnName = fieldNameCell.GetValue().(string)
				idCol.DataType = dataTypeCell.GetValue().(string)
				idCols = append(idCols, idCol)
			}
		}
	}

	return idCols, err
}

func (repo *DbRepository) getPrimaryKeyColumn() (DataColumn, error) {

	var pkCol DataColumn
	var err error

	hasPkCol := false

	for i := 0; i < len(repo.tbModel.Rows); i++ {
		row := repo.tbModel.Rows[i]
		cell, err := row.GetCell("isprimarykey")
		if err == nil && cell.GetValue() == true {
			hasPkCol = true
			var cellFieldName *DataCell
			var cellDataType *DataCell
			cellFieldName, err = row.GetCell("fieldname")
			if err == nil {
				cellDataType, err = row.GetCell("datatype")
			}
			if err == nil {
				pkCol.ColumnName = cellFieldName.GetValue().(string)
				pkCol.DataType = cellDataType.GetValue().(string)
			}
			break
		}
	}
	if !hasPkCol {
		err = fmt.Errorf("table \"%s\" does not have a primary key constraint", repo.TableName)
	}
	return pkCol, err
}

func (repo *DbRepository) Delete(cellValues map[string]interface{}) error {
	var err error

	var pkCol DataColumn
	pkCol, err = repo.getPrimaryKeyColumn()

	if err == nil {
		if !utils.HasKey(cellValues, pkCol.ColumnName) {
			err = errors.New("primary key column was not found in provided cellvalues")
		}
	}

	if err == nil {
		cmdTxt := "DELETE FROM \"%s\".\"%s\" WHERE \"%s\"=$1 "
		cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName, pkCol.ColumnName)
		_, err = repo.db.Query(cmdTxt, cellValues[pkCol.ColumnName])
	}

	return err
}

func (repo *DbRepository) DeleteWhere(whereSql string, args ...interface{}) error {
	cmdTxt := "DELETE FROM \"%s\".\"%s\" WHERE (%s) "
	cmdTxt = fmt.Sprintf(cmdTxt, repo.db.DefaultSchema, repo.TableName, whereSql)
	_, err := repo.db.Query(cmdTxt, args...)
	return err
}
