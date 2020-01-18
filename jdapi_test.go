package jdsdk

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	client := NewClient("xxx", "xxx")
	//client.GetCache = FileGetCache
	//client.SetCache = FileSetCache
	//client.WriteErrLog = InsertErrLog
	client.CacheLife = 0

	//jd.union.open.category.goods.get 商品类目查询
	//data, apiErr := client.NormalGetCategory(0, 0)
	//fmt.Println("return:", data, apiErr)

	//jd.union.open.order.query 订单查询接口
	//data2, apiErr2 := client.NormalGetOrder("20190711", 1, 20, 1, 0, "")
	//fmt.Println("return:", data2, apiErr2)

	//jd.union.open.promotion.common.get 获取通用推广链接//这个pid一直都是这样，提示错误，用不了。
	data3, apiErr3 := client.NormalGetPromotionUrl("https://item.jd.com/23484023378.html", 1224241220, 1847851615, "", "222", "", "")
	fmt.Println("return:", data3, apiErr3)

	//jd.union.open.goods.promotiongoodsinfo.query 获取通用推广链接,批量的
	//data4, apiErr4 := client.NormalGetPromotionGoodsInfoMultiple("50070163505,37173339786")
	//fmt.Println("return:", data4, apiErr4)

	//jd.union.open.goods.promotiongoodsinfo.query 获取通用推广链接,单个的
	//data5, apiErr5 := client.NormalGetPromotionGoodsInfoSingle(50070163505)
	//fmt.Println("return:", data5, apiErr5)

	//jd.union.open.goods.jingfen.query 京粉精选商品查询接口
	//data6, apiErr6 := client.NormalGetJFGoods(1, 1, 20, "inOrderCount30DaysSku", "desc")
	//fmt.Println("return:", data6, apiErr6)

	//jd.union.open.user.pid.get 获取PID //母账号无权限
	//data7, apiErr7 := client.NormalGetPid(1000218958, 1224241220, 1, "", "花生小宝")
	//fmt.Println("return:", data7, apiErr7)
}
