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

	var dt *DataTable
	if err == nil {
		dt, err = repo.Select(" dish_id=$1 ", "c9313e1b-4f03-431d-8987-5a96aa050fad")
	}

	if err == nil {
		for i := 0; i < len(*dt.Rows); i++ {
			row := (*dt.Rows)[i]
			for j := 0; j < len(*row.Cells); j++ {
				cell := (*row.Cells)[j]
				fmt.Println(cell.Value(nil))
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
