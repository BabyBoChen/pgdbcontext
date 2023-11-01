package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db, err := NewDbContext(connStr)
	defer db.Dispose()

	var repo *DbRepository
	if err == nil {
		repo, err = db.GetRepository("dish")
	}

	if err == nil {
		toUpdate := make(map[string]interface{})
		toUpdate["dish_id"] = "e90ba433-4181-4945-b1e5-3d2235bfef9d"
		toUpdate["unit_price"] = 75
		err = repo.Update(toUpdate)
		db.Commit()
	}

	// var lastInsertedId map[string]interface{}
	// if err == nil {
	// 	cellValues := make(map[string]interface{})
	// 	cellValues["title"] = "鮭魚壽司"
	// 	cellValues["unit_price"] = 25.00
	// 	cellValues["row_order"] = 13
	// 	lastInsertedId, err = repo.Insert(cellValues)
	// }

	// if err == nil {
	// 	db.Commit()
	// 	fmt.Println(lastInsertedId)
	// }

	var dt *DataTable
	if err == nil {
		dt, err = repo.Select("",)
	}

	if err == nil {
		for i := 0; i < len(*dt.Rows); i++ {
			row := (*dt.Rows)[i]
			for j := 0; j < len(*row.Cells); j++ {
				cell := (*row.Cells)[j]
				fmt.Println(cell.GetValue())
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
