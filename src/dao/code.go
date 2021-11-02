package dao

import (
	"log"

	"github.com/cheolgyu/stock-write-common/db"
	"github.com/cheolgyu/stock-write-model/model"
)

func Update_info() {
	query := `UPDATE public.info SET  updated= now() WHERE name = 'updated_rebound' `

	_, err := db.Conn.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func GetCodeAll() ([]model.Code, error) {
	var res []model.Code
	rows, err := db.Conn.Query("select id, code, code_type from meta.code   order by id  ")
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		c := model.Code{}
		if err := rows.Scan(&c.Id, &c.Code, &c.Code_type); err != nil {
			log.Fatal(err)
			panic(err)
		}
		res = append(res, c)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}
	return res, err
}

func GetConfig_Upper_Code() ([]model.Config, error) {

	var res []model.Config
	query := `select id, upper_code, upper_name, code, name from meta.config where upper_code = 'price_type' order by id ;`
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Config
		if err := rows.Scan(&item.Id, &item.Upper_code, &item.Upper_name, &item.Code, &item.Name); err != nil {
			log.Fatal(err)
			panic(err)
		}
		res = append(res, item)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return res, err
}
