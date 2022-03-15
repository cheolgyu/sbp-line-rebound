package handler

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/cheolgyu/graph/bound/dao"
	"github.com/cheolgyu/model"
)

var ch_price_rebound chan Calculate
var wg_price_rebound sync.WaitGroup

type Calculate struct {
	model.Code
	list []rebound_price_type
}

func ChannelReboundCalculate(ch chan Calculate) {
	log.Println("run  ChannelCalculate")
	// loop over the data from the channel
	for v := range ch {

		txt := fmt.Sprintf("ch_price_rebound <   %+v", v.Code)
		log.Println(txt)

		v.exec()
		ri := ReboundSqlWrite{
			Code: v.Code,
			list: v.list,
		}
		ch_price_sql_write <- ri

	}
}

func (o *Calculate) exec() {
	defer wg_price_rebound.Done()

	for i := range o.list {
		o.list[i].get_rebound_point()
	}

}

type rebound_price_type struct {
	model.Code
	price_type model.Config
	PriceList  []model.PriceMarket
	PointList  []model.Point
}

// price_type별 각각의 RE_BOUND_POINT 구하기
func (o *rebound_price_type) get_rebound_point() {
	/*
		변화 구분 0: eq 1 증가 -1 감소
		변화시작점 변화종료점
		변화 기간

	*/

	count := len(o.PriceList)
	if count < 2 {
		return
	}

	chg_value := change_graph_direction(o.get_price_value(0), o.get_price_value(0+1))
	chg_start_x := uint(o.PriceList[0].Dt)
	chg_start_y := o.get_price_value(0)

	chg_tick := 0
	//log.Println("chg_start_x, chg_start_y,chg_tick", chg_start_x, chg_start_y, chg_tick)

	for i := 0; i < count; i++ {

		// log.Println("iiiiiiiiiiiiiiiiiii==================>>>>>>%s", i)
		// txt := fmt.Sprintf("loop %+v", o.PriceList[i])
		// log.Println(txt)

		chg_tick++
		n := i + 1
		last_save := false
		if n == count {
			//log.Println("마지막은 저장해야지")
			last_save = true
			n = i
		}

		x1 := uint(o.PriceList[i].Dt)
		y1 := o.get_price_value(i)

		y2 := o.get_price_value(n)

		g_way := change_graph_direction(y1, y2)
		chg := is_rebound(chg_value, g_way)
		//log.Println("g_way, chg_value, chg", g_way, chg_value, chg)

		if chg || last_save {
			var bp = model.Point{
				Code_id:    o.Code.Id,
				Price_type: o.price_type.Id,
			}
			var err = bp.Set(chg_start_x, chg_start_y, x1, y1, uint(chg_tick))
			if err != nil {
				txt := fmt.Sprintf("o.PointList= %+v \n", o.PointList)
				log.Println(txt)
				panic(err)
			}

			o.PointList = append(o.PointList, bp)

			chg_tick = 0
			chg_value = g_way
			chg_start_x = x1
			chg_start_y = y1

		}

	}

}

// price_type별 각각의 가격 목록 조회.
func (o *rebound_price_type) get_price() {

	pmlist, err := dao.GetPriceByLastReBound(o.Code.Id, o.price_type.Id)

	if err != nil {
		log.Println("오류:GetPriceByLastReBound :", o.Code.Id, o.price_type.Id)
	}
	o.PriceList = pmlist
}

func (o *rebound_price_type) get_price_value(i int) float32 {
	switch o.price_type.Id {
	case price_type_config["low"]:
		return float32(o.PriceList[i].LowPrice)
	case price_type_config["high"]:
		return float32(o.PriceList[i].HighPrice)
	case price_type_config["close"]:
		return float32(o.PriceList[i].ClosePrice)
	case price_type_config["open"]:
		return float32(o.PriceList[i].OpenPrice)
	default:
		panic(errors.New("머머머머ㅓ머지?"))
	}
}

func is_rebound(ago_g_way int, cur_g_way int) bool {

	if cur_g_way == 0 {
		return false
	}
	if ago_g_way == cur_g_way {
		return false
	} else {
		return true
	}

}

func change_graph_direction(y1 float32, y2 float32) int {
	g_way := 0
	if y1 < y2 {
		g_way = 1
	} else if y1 > y2 {
		g_way = -1
	} else if y1 == y2 {
		g_way = 0
	} else {
		panic("머냐")
	}
	return g_way
}
