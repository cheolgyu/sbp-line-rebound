package dao

import (
	"context"
	"log"

	"github.com/cheolgyu/base/db"
	"github.com/cheolgyu/graph/bound/c"
	"github.com/cheolgyu/model"
	"github.com/gchaincl/dotsql"
	"github.com/swithek/dotsqlx"
)

const q_last_rebound = `SELECT  dt, op, hp, lp, cp FROM hist.price WHERE CODE_ID =$1  AND  dt >=  ` +
	`(SELECT  COALESCE(MAX(X1),0) AS X1  FROM hist.rebound WHERE CODE_ID =$1  AND PRICE_TYPE=$2 ) order by dt asc `

const q_insert_hist_rebound = `INSERT INTO hist.rebound ` +
	`(code_id, price_type, x1, y1, x2, y2, x_tick, y_minus, y_percent) ` +
	`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ` +
	`ON CONFLICT (code_id,price_type,x1) DO UPDATE SET ` +
	`x1=$3, y1=$4, x2=$5, y2=$6, x_tick=$7, y_minus=$8, y_percent=$9 `

const q_insert_public_rebound = `INSERT INTO  public.tb_rebound (` +
	`code_id, cp_x1, cp_y1, cp_x2, cp_y2, cp_x_tick, cp_y_minus, cp_y_percent, op_x1, op_y1, op_x2, op_y2, op_x_tick, op_y_minus, op_y_percent, lp_x1, lp_y1, lp_x2, lp_y2, lp_x_tick, lp_y_minus, lp_y_percent, hp_x1, hp_y1, hp_x2, hp_y2, hp_x_tick, hp_y_minus, hp_y_percent)` +
	`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)` +
	` ON CONFLICT ("code_id") DO UPDATE SET ` +
	`cp_x1=$2, cp_y1=$3, cp_x2=$4, cp_y2=$5, cp_x_tick=$6, cp_y_minus=$7, cp_y_percent=$8, op_x1=$9, op_y1=$10,` +
	`op_x2=$11, op_y2=$12, op_x_tick=$13, op_y_minus=$14, op_y_percent=$15, lp_x1=$16, lp_y1=$17, lp_x2=$18, lp_y2=$19, lp_x_tick=$20,` +
	`lp_y_minus=$21, lp_y_percent=$22, hp_x1=$23, hp_y1=$24, hp_x2=$25, hp_y2=$26, hp_x_tick=$27, hp_y_minus=$28, hp_y_percent=$29`

// 바운스 마지막 시작일자 보다 큰 가격목록 조회.

// obj price 도 model.PriceMarket 로  담아서 보내고 사용시 model.PriceMarket.ToPriceStock() 이용하기.
func GetPriceByLastReBound(code_id int, price_type_id int) ([]model.PriceMarket, error) {
	res_market := []model.PriceMarket{}

	rows, err := db.Conn.QueryContext(context.Background(), q_last_rebound, code_id, price_type_id)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := model.PriceMarket{}

		if err := rows.Scan(&item.Dt, &item.OpenPrice, &item.HighPrice, &item.LowPrice, &item.ClosePrice); err != nil {
			log.Fatal(err)
			panic(err)
		}
		res_market = append(res_market, item)

	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return res_market, err
}

func InsertHistReBound(fnm string) error {
	client := db.Conn
	var dot *dotsql.DotSql
	var err error

	if dot, err = dotsql.LoadFromFile(fnm); err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		dotx := dotsqlx.Wrap(dot)
		_, err = dotx.Exec(client, c.DOTSQL_NAME_REBOUND)
	}

	return err
}

/*
 2021/07/21 19:39:43 main.go:21: [걸린시간] Elipsed Time: 16.8630616s

*/
func InsertPublicReBound(item model.TbReBound, upsert bool) error {

	client := db.Conn
	stmt, err := client.Prepare(q_insert_public_rebound)
	if err != nil {
		log.Println("쿼리:Prepare 오류: ", item)
		log.Fatal(err)
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Code_id,
		item.Cp_X1, item.Cp_Y1, item.Cp_X2, item.Cp_Y2, item.Cp_X_tick, item.Cp_Y_minus, item.Cp_Y_Percent,
		item.Op_X1, item.Op_Y1, item.Op_X2, item.Op_Y2, item.Op_X_tick, item.Op_Y_minus, item.Op_Y_Percent,
		item.Lp_X1, item.Lp_Y1, item.Lp_X2, item.Lp_Y2, item.Lp_X_tick, item.Lp_Y_minus, item.Lp_Y_Percent,
		item.Hp_X1, item.Hp_Y1, item.Hp_X2, item.Hp_Y2, item.Hp_X_tick, item.Hp_Y_minus, item.Hp_Y_Percent,
	)
	if err != nil {
		log.Println("쿼리:stmt.Exec 오류: ", item)
		log.Fatal(err)
		panic(err)
	}
	return err
}
