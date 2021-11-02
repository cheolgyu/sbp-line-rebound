package handler

import (
	"fmt"
	"log"
	"sync"

	"github.com/cheolgyu/stock-write-model/model"
	"github.com/cheolgyu/stock-write-project-rebound/src/dao"
)

var upsert_bound bool
var price_type_arr []model.Config
var price_type_config map[string]int

var wg_price sync.WaitGroup = sync.WaitGroup{}

func init() {

	ch_price_rebound = make(chan Calculate)
	ch_price_sql_write = make(chan ReboundSqlWrite)

	wg_price = sync.WaitGroup{}
	wg_price_rebound = sync.WaitGroup{}
	wg_price_insert = sync.WaitGroup{}

	go ChannelReboundCalculate(ch_price_rebound)
	go ChannelReboundSqlWrite(ch_price_sql_write)

	upsert_bound = true

	_price_type_arr, err := dao.GetConfig_Upper_Code()
	ChkErr(err)

	price_type_arr = _price_type_arr

	price_type_config = make(map[string]int)
	for i := range price_type_arr {
		price_type_config[price_type_arr[i].Code] = price_type_arr[i].Id
	}
}

func ReBoundHandler() {

	code_list, err := dao.GetCodeAll()
	ChkErr(err)

	log.Println(" ReBoundHandler  start")
	bp := ReBound{}
	bp.Save(code_list)

	log.Println(" ReBoundHandler  end")
}

type ReBound struct {
}

// BOUND_POINT 저장.
func (o *ReBound) Save(list []model.Code) {

	for i := range list {
		item := fmt.Sprintf("index=%v ,  %+v\n", i, list[i])
		log.Println("item:", item)

		cc := list[i]
		bc := code_rebound{Code: cc}

		wg_price.Add(1)
		wg_price_insert.Add(1)
		wg_price_rebound.Add(1)
		go bc.get_price()

	}
	wg_price.Wait()
	wg_price_rebound.Wait()
	wg_price_insert.Wait()

}

type code_rebound struct {
	model.Code
}

// CODE에 해당하는 가격목록 조회.
func (o *code_rebound) get_price() {
	defer wg_price.Done()

	var item Calculate
	item.Code = o.Code

	for i := range price_type_arr {
		gcg := rebound_price_type{
			Code:       o.Code,
			price_type: price_type_arr[i],
		}
		gcg.get_price()
		item.list = append(item.list, gcg)
		// txt := fmt.Sprintf(" gcg  %+v", gcg)
		// log.Println(txt)
	}
	ch_price_rebound <- item
}

func ChkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
