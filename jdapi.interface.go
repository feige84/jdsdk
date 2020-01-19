package jdsdk

import (
	"fmt"
	"strings"
)

//jd.union.open.category.goods.get 商品类目查询
func (client *ApiReq) NormalGetCategory(parentId, grade int64) ([]CategoryResp, *ApiErrorInfo) {
	client.ParamName = "req"
	resp, err := client.Execute("jd.union.open.category.goods.get", ApiParams{
		"parentId": fmt.Sprint(parentId),
		"grade":    fmt.Sprint(grade),
	})
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var categoryList = []CategoryResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	categoryResp := resp.Get("data")
	if respCode == 200 && categoryResp.Exists() && categoryResp.IsArray() {
		for _, value := range categoryResp.Array() {
			category := CategoryResp{}
			category.Id = value.Get("id").Int()
			category.ParentId = value.Get("parentId").Int()
			category.Grade = value.Get("grade").Int()
			category.Name = value.Get("name").String()
			categoryList = append(categoryList, category)
		}
		return categoryList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.order.query 订单查询接口
func (client *ApiReq) NormalGetOrder(orderTime string, pageNo, pageSize, timeType, childUnionId int64, key string) ([]OrderResp, *ApiErrorInfo) {
	client.ParamName = "orderReq"
	params := ApiParams{}
	params["pageNo"] = fmt.Sprint(pageNo) //页码，返回第几页结果
	if pageSize <= 0 {
		params["pageSize"] = "500" //每页包含条数，上限为500
	} else {
		params["pageSize"] = fmt.Sprint(pageSize)
	}
	params["type"] = fmt.Sprint(timeType)             //订单时间查询类型(1：下单时间，2：完成时间，3：更新时间)
	params["time"] = orderTime                        //查询时间，建议使用分钟级查询，格式：yyyyMMddHH、yyyyMMddHHmm或yyyyMMddHHmmss，如201811031212 的查询范围从12:12:00--12:12:59
	params["childUnionId"] = fmt.Sprint(childUnionId) //子站长ID（需要联系运营开通PID账户权限才能拿到数据），childUnionId和key不能同时传入
	params["key"] = key                               //其他推客的授权key，查询工具商订单需要填写此项，childUnionid和key不能同时传入
	resp, err := client.Execute("jd.union.open.order.query", params)
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var orderList = []OrderResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	orderResp := resp.Get("data")
	if respCode == 200 {
		if orderResp.Exists() && orderResp.IsArray() {
			for _, value := range orderResp.Array() {
				order := OrderResp{}

				//订单信息
				order.FinishTime = value.Get("finishTime").Int() //订单完成时间(时间戳，毫秒)
				order.OrderEmt = value.Get("orderEmt").Int()     //下单设备(1:PC, 2:无线)
				order.OrderId = value.Get("orderId").Int()       //订单ID
				order.OrderTime = value.Get("orderTime").Int()   //下单时间(时间戳，毫秒)
				order.ParentId = value.Get("parentId").Int()     //父单的订单ID，仅当发生订单拆分时返回， 0：未拆分，有值则表示此订单为子订单
				order.PayMonth = value.Get("payMonth").String()  //订单维度预估结算时间（格式：yyyyMMdd），0：未结算，订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
				order.Plus = value.Get("plus").Int()             //下单用户是否为PLUS会员 0：否，1：是
				order.PopId = value.Get("popId").Int()           //商家ID
				order.UnionId = value.Get("unionId").Int()       //推客的联盟ID
				order.Ext1 = value.Get("ext1").String()          //推客生成推广链接时传入的扩展字段，订单维度（需要联系运营开放白名单才能拿到数据）
				order.ValidCode = value.Get("validCode").Int()   //订单维度的有效码（-1：未知, 2.无效-拆单, 3.无效-取消, 4.无效-京东帮帮主订单, 5.无效-账号异常, 6.无效-赠品类目不返佣, 7.无效-校园订单, 8.无效-企业订单, 9.无效-团购订单, 10.无效-开增值税专用发票订单, 11.无效-乡村推广员下单, 12.无效-自己推广自己下单, 13.无效-违规订单, 14.无效-来源与备案网址不符, 15.待付款, 16.已付款, 17.已完成, 18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
				order.HasMore = value.Get("hasMore").Bool()      //是否还有更多, true：还有数据；false:已查询完毕，没有数据

				//商品信息

				skuListResp := value.Get("skuList")
				if skuListResp.Exists() && skuListResp.IsArray() {
					skuList := []SkuInfo{}
					for _, v := range skuListResp.Array() {
						sku := SkuInfo{}
						sku.ActualCosPrice = v.Get("actualCosPrice").Float()     //实际计算佣金的金额。订单完成后，会将误扣除的运费券金额更正。如订单完成后发生退款，此金额会更新。
						sku.ActualFee = v.Get("actualFee").Float()               //推客获得的实际佣金（实际计佣金额*佣金比例*最终比例）。如订单完成后发生退款，此金额会更新。
						sku.CommissionRate = v.Get("commissionRate").Float()     //佣金比例
						sku.EstimateCosPrice = v.Get("estimateCosPrice").Float() //预估计佣金额，即用户下单的金额(已扣除优惠券、白条、支付优惠、进口税，未扣除红包和京豆)，有时会误扣除运费券金额，完成结算时会在实际计佣金额中更正。如订单完成前发生退款，此金额不会更新。
						sku.EstimateFee = v.Get("estimateFee").Float()           //推客的预估佣金（预估计佣金额*佣金比例*最终比例），如订单完成前发生退款，此金额不会更新。
						sku.FinalRate = v.Get("finalRate").Float()               //最终比例（分成比例+补贴比例）
						sku.Cid1 = v.Get("cid1").Int()                           //一级类目ID
						sku.Cid2 = v.Get("cid2").Int()                           //二级类目ID
						sku.Cid3 = v.Get("cid3").Int()                           //三级类目ID
						sku.FrozenSkuNum = v.Get("frozenSkuNum").Int()           //商品售后中数量
						sku.Pid = v.Get("pid").String()                          //联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID
						sku.PositionId = v.Get("positionId").Int()               //推广位ID,0代表无推广位
						sku.Price = v.Get("price").Float()                       //商品单价
						sku.SiteId = v.Get("siteId").Int()                       //网站ID，0：无网站
						sku.SkuId = v.Get("skuId").Int()                         //商品ID
						sku.SkuName = v.Get("skuName").String()                  //商品名称
						sku.SkuNum = v.Get("skuNum").Int()                       //商品数量
						sku.SkuReturnNum = v.Get("skuReturnNum").Int()           //商品已退货数量
						sku.SubSideRate = v.Get("subSideRate").Float()           //分成比例
						sku.SubsidyRate = v.Get("subsidyRate").Float()           //补贴比例
						sku.UnionAlias = v.Get("unionAlias").String()            //PID所属母账号平台名称（原第三方服务商来源）
						sku.UnionTag = v.Get("unionTag").String()                //联盟标签数据（整型的二进制字符串(32位)，目前只返回8位：00000001。数据从右向左进行，每一位为1表示符合联盟的标签特征，第1位：京喜红包，第2位：组合推广订单，第3位：拼购订单，第5位：有效首次购订单（00011XXX表示有效首购，最终奖励活动结算金额会结合订单状态判断，以联盟后台对应活动效果数据报表https://union.jd.com/active为准）。例如：00000001:京喜红包订单，00000010:组合推广订单，00000100:拼购订单，00011000:有效首购，00000111：京喜红包+组合推广+拼购等）
						sku.UnionTrafficGroup = v.Get("unionTrafficGroup").Int() //渠道组 1：1号店，其他：京东
						sku.ValidCode = v.Get("validCode").Int()                 //sku维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
						sku.SubUnionId = v.Get("subUnionId").String()            //子联盟ID(需要联系运营开放白名单才能拿到数据)
						sku.TraceType = v.Get("traceType").Int()                 //2：同店；3：跨店
						sku.PayMonth = v.Get("payMonth").String()                //订单行维度预估结算时间（格式：yyyyMMdd） ，0：未结算。订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
						sku.PopId = v.Get("popId").Int()                         //商家ID，订单行维度
						sku.Ext1 = v.Get("popId").String()                       //推客生成推广链接时传入的扩展字段（需要联系运营开放白名单才能拿到数据）。&lt;订单行维度&gt;

						skuList = append(skuList, sku)
					}
					order.SkuList = skuList
				}

				orderList = append(orderList, order)
			}
		}
		return orderList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.promotion.common.get 获取通用推广链接
func (client *ApiReq) NormalGetPromotionUrl(materialId string, siteId, positionId int64, subUnionId, ext1, pid, couponUrl string) (*PromotionCodeResp, *ApiErrorInfo) {
	client.ParamName = "promotionCodeReq"
	params := ApiParams{}

	params["materialId"] = materialId     //推广物料落体页
	params["siteId"] = fmt.Sprint(siteId) //站点ID是指在联盟后台的推广管理中的网站Id、APPID（1、通用转链接口禁止使用社交媒体id入参；2、订单来源，即投放链接的网址或应用必须与传入的网站ID/AppID备案一致，否则订单会判“无效-来源与备案网址不符”）
	if positionId > 0 {
		params["positionId"] = fmt.Sprint(positionId) //推广位id
	}
	if subUnionId != "" {
		params["subUnionId"] = subUnionId //子联盟ID（需要联系运营开通权限才能拿到数据）
	}
	if ext1 != "" {
		params["ext1"] = ext1 //推客生成推广链接时传入的扩展字段（查看订单对应字段信息，需要联系运营开放白名单才能看到）
	}
	if pid != "" {
		params["pid"] = pid //联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID
	}
	if couponUrl != "" {
		params["couponUrl"] = couponUrl //优惠券领取链接，在使用优惠券、商品二合一功能时入参，且materialId须为商品详情页链接
	}
	resp, err := client.Execute("jd.union.open.promotion.common.get", params)
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	promotionResp := resp.Get("data")
	if respCode == 200 && promotionResp.Exists() {
		promotion := &PromotionCodeResp{}
		promotion.ClickURL = promotionResp.Get("clickURL").String()
		return promotion, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.goods.promotiongoodsinfo.query 获取通用推广链接,单个的
func (client *ApiReq) NormalGetPromotionGoodsInfoSingle(skuId int64) (*PromotionGoodsResp, *ApiErrorInfo) {
	goodsInfo, apiErr := client.NormalGetPromotionGoodsInfoMultiple(fmt.Sprint(skuId))
	if goods, exists := goodsInfo[skuId]; exists {
		return goods, apiErr
	} else {
		return nil, apiErr
	}
}

//jd.union.open.goods.promotiongoodsinfo.query 获取通用推广链接,批量的
func (client *ApiReq) NormalGetPromotionGoodsInfoMultiple(skuIds string) (map[int64]*PromotionGoodsResp, *ApiErrorInfo) {
	client.ParamName = ""
	resp, err := client.Execute("jd.union.open.goods.promotiongoodsinfo.query", ApiParams{
		"skuIds": skuIds,
	})
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		fmt.Println("xxxxxxxxxx")
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	//fmt.Println("respCode:", respCode)
	promotionGoodsResp := resp.Get("data")
	if respCode == 200 && promotionGoodsResp.Exists() {
		promotionGoodsList := make(map[int64]*PromotionGoodsResp)
		if promotionGoodsResp.IsArray() && len(promotionGoodsResp.Array()) > 0 {
			for _, vaule := range promotionGoodsResp.Array() {
				goods := PromotionGoodsResp{}
				goods.SkuId = vaule.Get("skuId").Int()                         //商品ID
				goods.UnitPrice = vaule.Get("unitPrice").Float()               //商品单价即京东价
				goods.MaterialUrl = vaule.Get("materialUrl").String()          //商品落地页
				goods.EndDate = vaule.Get("endDate").Int()                     //推广结束日期(时间戳，毫秒)
				goods.IsFreeFreightRisk = vaule.Get("isFreeFreightRisk").Int() //是否支持运费险(1:是,0:否)
				goods.IsFreeShipping = vaule.Get("isFreeShipping").Int()       //是否包邮(1:是,0:否,2:自营商品遵从主站包邮规则)
				goods.CommisionRatioWl = vaule.Get("commisionRatioWl").Float() //无线佣金比例
				goods.CommisionRatioPc = vaule.Get("commisionRatioPc").Float() //PC佣金比例
				goods.ImgUrl = vaule.Get("imgUrl").String()                    //图片地址
				goods.Vid = vaule.Get("vid").Int()                             //商家ID
				goods.Cid = vaule.Get("cid").Int()                             //一级类目ID
				goods.CidName = vaule.Get("cidName").String()                  //一级类目名称
				goods.Cid2 = vaule.Get("cid2").Int()                           //二级类目ID
				goods.Cid2Name = vaule.Get("cid2Name").String()                //二级类目名称
				goods.Cid3 = vaule.Get("cid3").Int()                           //三级类目ID
				goods.Cid3Name = vaule.Get("cid3Name").String()                //三级类目名称
				goods.WlUnitPrice = vaule.Get("wlUnitPrice").Float()           //商品无线京东价（单价为-1表示未查询到该商品单价）
				goods.IsSeckill = vaule.Get("isSeckill").Int()                 //是否秒杀(1:是,0:否)
				goods.InOrderCount = vaule.Get("inOrderCount").Int()           //30天引单数量
				goods.ShopId = vaule.Get("shopId").Int()                       //店铺ID
				goods.IsJdSale = vaule.Get("isJdSale").Int()                   //是否自营(1:是,0:否)
				goods.GoodsName = vaule.Get("goodsName").String()              //商品名称
				goods.StartDate = vaule.Get("startDate").Int()                 //推广开始日期（时间戳，毫秒）
				promotionGoodsList[goods.SkuId] = &goods
			}
			if len(promotionGoodsList) > 0 {
				return promotionGoodsList, nil
			}
		}
		return nil, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.goods.jingfen.query 京粉精选商品查询接口
func (client *ApiReq) NormalGetJFGoods(eliteId, pageIndex, pageSize int64, sortName, sort string) ([]JFGoodsResp, *ApiErrorInfo) {
	client.ParamName = "goodsReq"
	params := ApiParams{}
	params["eliteId"] = fmt.Sprint(eliteId)     //1-好券商品,2-京粉APP-jingdong.超级大卖场,3-小程序-jingdong.好券商品,4-京粉APP-jingdong.主题聚惠1-jingdong.服装运动,5-京粉APP-jingdong.主题聚惠2-jingdong.精选家电,6-京粉APP-jingdong.主题聚惠3-jingdong.超市,7-京粉APP-jingdong.主题聚惠4-jingdong.居家生活,10-9.9专区,11-品牌好货-jingdong.潮流范儿,12-品牌好货-jingdong.精致生活,13-品牌好货-jingdong.数码先锋,14-品牌好货-jingdong.品质家电,15-京仓配送,16-公众号-jingdong.好券商品,17-公众号-jingdong.9.9,18-公众号-jingdong.京东配送
	params["pageIndex"] = fmt.Sprint(pageIndex) //页码，返回第几页结果
	if pageSize <= 0 {
		params["pageSize"] = "50" //每页数量，默认20，上限50
	} else {
		params["pageSize"] = fmt.Sprint(pageSize)
	}
	params["sortName"] = sortName //排序字段(price：单价, commissionShare：佣金比例, commission：佣金， inOrderCount30DaysSku：sku维度30天引单量，comments：评论数，goodComments：好评数)
	params["sort"] = sort         //asc,desc升降序,默认降序
	resp, err := client.Execute("jd.union.open.goods.jingfen.query", params)
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var goodsList = []JFGoodsResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	goodsResp := resp.Get("data")
	if respCode == 200 && goodsResp.Exists() && goodsResp.IsArray() {
		for _, value := range goodsResp.Array() {
			goods := JFGoodsResp{}

			//商品基本信息
			goods.Comments = value.Get("comments").Int()                           //评论数
			goods.GoodCommentsShare = value.Get("goodCommentsShare").Float()       //商品好评率
			goods.InOrderCount30Days = value.Get("inOrderCount30Days").Int()       //30天引单数量
			goods.MaterialUrl = value.Get("materialUrl").String()                  //商品落地页
			goods.SkuId = value.Get("skuId").Int()                                 //商品ID
			goods.SkuName = value.Get("skuName").String()                          //商品名称
			goods.IsHot = value.Get("isHot").Int()                                 //是否爆款，1：是，0：否
			goods.Spuid = value.Get("spuid").Int()                                 //其值为同款商品的主skuid
			goods.BrandCode = value.Get("brandCode").String()                      //品牌code
			goods.BrandName = value.Get("brandName").String()                      //品牌名
			goods.Owner = value.Get("owner").String()                              //g=自营，p=pop
			goods.InOrderCount30DaysSku = value.Get("inOrderCount30DaysSku").Int() //30天引单数量(sku维度)
			goods.TotalCount = value.Get("totalCount").Int()                       //有效商品总数量

			//类目信息
			categoryInfoResp := value.Get("categoryInfo")
			if categoryInfoResp.Exists() {
				goods.CategoryInfo.Cid1 = categoryInfoResp.Get("cid1").Int()
				goods.CategoryInfo.Cid1Name = categoryInfoResp.Get("cid1Name").String()
				goods.CategoryInfo.Cid2 = categoryInfoResp.Get("cid2").Int()
				goods.CategoryInfo.Cid2Name = categoryInfoResp.Get("cid2Name").String()
				goods.CategoryInfo.Cid3 = categoryInfoResp.Get("cid3").Int()
				goods.CategoryInfo.Cid3Name = categoryInfoResp.Get("cid3Name").String()
			}

			//图片信息
			imageListResp := value.Get("imageInfo.imageList")
			if imageListResp.Exists() && imageListResp.IsArray() {
				goods.SkuImage = imageListResp.Array()[0].Get("url").String() //自己加的。取第一张作为主图
				for _, v := range imageListResp.Array() {
					goods.ImageInfo.imageList = append(goods.ImageInfo.imageList, v.Get("url").String())
				}
			}

			//佣金信息
			commissionInfoResp := value.Get("commissionInfo")
			if commissionInfoResp.Exists() {
				goods.CommissionInfo.Commission = commissionInfoResp.Get("commission").Float()
				goods.CommissionInfo.CommissionShare = commissionInfoResp.Get("commissionShare").Float()
			}

			//价格信息
			priceInfoResp := value.Get("priceInfo")
			if priceInfoResp.Exists() {
				goods.PriceInfo.Price = priceInfoResp.Get("price").Float()
			}

			//店铺信息
			shopInfoResp := value.Get("shopInfo")
			if shopInfoResp.Exists() {
				goods.ShopInfo.ShopId = shopInfoResp.Get("shopId").Int()
				goods.ShopInfo.ShopName = shopInfoResp.Get("shopName").String()
			}

			//优惠券信息，返回内容为空说明该SKU无可用优惠券
			couponListResp := value.Get("couponInfo.couponList")
			if couponListResp.Exists() && couponListResp.IsArray() {
				for _, v := range couponListResp.Array() {
					coupon := Coupon{}
					coupon.BindType = v.Get("bindType").Int()         //券种类 (优惠券种类：0 - 全品类，1 - 限品类（自营商品），2 - 限店铺，3 - 店铺限商品券)
					coupon.Discount = v.Get("discount").Float()       //券面额
					coupon.Link = v.Get("link").String()              //券链接
					coupon.PlatformType = v.Get("platformType").Int() //券使用平台 (平台类型：0 - 全平台券，1 - 限平台券)
					coupon.Quota = v.Get("quota").Float()             //券消费限额
					coupon.GetStartTime = v.Get("getStartTime").Int() //领取开始时间(时间戳，毫秒)
					coupon.GetEndTime = v.Get("getEndTime").Int()     //券领取结束时间(时间戳，毫秒)
					coupon.UseStartTime = v.Get("useStartTime").Int() //券有效使用开始时间(时间戳，毫秒)
					coupon.UseEndTime = v.Get("useEndTime").Int()     //券有效使用结束时间(时间戳，毫秒)
					goods.CouponInfo.CouponList = append(goods.CouponInfo.CouponList, coupon)
				}
			}

			//拼购信息
			pinGouInfoResp := value.Get("pinGouInfo")
			if pinGouInfoResp.Exists() {
				goods.PinGouInfo.PingouPrice = pinGouInfoResp.Get("pingouPrice").Float()
				goods.PinGouInfo.PingouTmCount = pinGouInfoResp.Get("pingouTmCount").Int()
				goods.PinGouInfo.PingouUrl = pinGouInfoResp.Get("pingouUrl").String()
			}

			//资源信息
			resourceInfoResp := value.Get("resourceInfo")
			if resourceInfoResp.Exists() {
				goods.ResourceInfo.EliteId = resourceInfoResp.Get("eliteId").Int()
				goods.ResourceInfo.EliteName = pinGouInfoResp.Get("eliteName").String()
			}

			goodsList = append(goodsList, goods)
		}
		return goodsList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.user.pid.get 获取PID
func (client *ApiReq) NormalGetPid(unionId, childUnionId, promotionType int64, positionName, mediaName string) (string, *ApiErrorInfo) {
	client.ParamName = "pidReq"
	client.V = "1.1"
	params := ApiParams{}

	params["unionId"] = fmt.Sprint(unionId)             //联盟ID
	params["childUnionId"] = fmt.Sprint(childUnionId)   //子站长ID
	params["promotionType"] = fmt.Sprint(promotionType) //推广类型,1APP推广 2聊天工具推广
	if positionName != "" {
		params["positionName"] = positionName //子站长的推广位名称，如不存在则创建，不填则由联盟根据母账号信息创建
	}
	params["mediaName"] = mediaName //媒体名称，即子站长的app应用名称，推广方式为app推广时必填，且app名称必须为已存在的app名称
	resp, err := client.Execute("jd.union.open.user.pid.get", params)
	if err != nil {
		return "", err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return "", apiErrInfo
	}
	respCode := resp.Get("code").Int()
	promotionResp := resp.Get("data")
	if respCode == 200 && promotionResp.Exists() {
		return promotionResp.String(), nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return "", &errInfo
	}
}

//jd.union.open.goods.query 关键词商品查询接口【申请】
func (client *ApiReq) AdvancedGetGoods(cid1, cid2, cid3, pageIndex, pageSize int64, skuIds interface{}, keyword string,
	priceFrom, priceTo float64, commissionShareStart, commissionShareEnd int64,
	owner, sortName, sort string, isCoupon, isPG, isHot int64, pingouPriceStart, pingouPriceEnd float64,
	brandCode string, shopId int64) ([]GoodsResp, *ApiErrorInfo) {

	client.ParamName = "goodsReqDTO"
	params := ApiParams{}
	if cid1 > 0 {
		params["cid1"] = fmt.Sprint(cid1) //一级类目id
	}
	if cid2 > 0 {
		params["cid2"] = fmt.Sprint(cid2) //二级类目id
	}
	if cid3 > 0 {
		params["cid3"] = fmt.Sprint(cid3) //一级类目id
	}
	params["pageIndex"] = fmt.Sprint(pageIndex) //页码，返回第几页结果
	if pageSize <= 0 {
		params["pageSize"] = "30" //每页数量，单页数最大30，默认20
	} else {
		params["pageSize"] = fmt.Sprint(pageSize)
	}

	//skuid集合(一次最多支持查询100个sku)，英文逗号分割，数组类型开发时记得加[]
	if skuIds != nil {
		switch result := skuIds.(type) {
		case string:
			params["skuIds"] = result
		case []int64:
			tmp, ids := "", ""
			for _, v := range result {
				ids = tmp + fmt.Sprint(v)
				tmp = ","
			}
			if ids != "" {
				params["skuIds"] = ids
			}
		case []string:
			ids := strings.Join(result, ",")
			if ids != "" {
				params["skuIds"] = ids
			}
		}
	}

	if keyword != "" {
		params["keyword"] = keyword //关键词，字数同京东商品名称一致，目前未限制
	}
	if priceFrom > 0 {
		params["pricefrom"] = fmt.Sprint(priceFrom) //商品价格下限
	}
	if priceTo > 0 {
		params["priceto"] = fmt.Sprint(priceTo) //商品价格上限
	}

	if commissionShareStart > 0 {
		params["commissionShareStart"] = fmt.Sprint(commissionShareStart) //佣金比例区间开始
	}
	if commissionShareEnd > 0 {
		params["commissionShareEnd"] = fmt.Sprint(commissionShareEnd) //佣金比例区间结束
	}
	if owner != "" {
		params["owner"] = owner //商品类型：自营[g]，POP[p]
	}
	if sortName != "" {
		params["sortName"] = sortName //排序字段(price：单价, commissionShare：佣金比例, commission：佣金， inOrderCount30Days：30天引单量， inOrderComm30Days：30天支出佣金)
	}
	if sort != "" {
		params["sort"] = sort //asc,desc升降序,默认降序
	}
	if isCoupon > 0 {
		params["isCoupon"] = "1" //是否是优惠券商品，1：有优惠券，0：无优惠券
	}
	if isPG > 0 {
		params["isPG"] = "1" //是否是拼购商品，1：拼购商品，0：非拼购商品
		if pingouPriceStart > 0 {
			params["pingouPriceStart"] = fmt.Sprint(pingouPriceStart) //拼购价格区间开始
		}
		if pingouPriceEnd > 0 {
			params["pingouPriceEnd"] = fmt.Sprint(pingouPriceEnd) //拼购价格区间结束
		}
	}
	if isHot > 0 {
		params["isHot"] = "1" //是否是爆款，1：爆款商品，0：非爆款商品
	}
	if brandCode != "" {
		params["brandCode"] = brandCode //品牌code
	}
	if shopId > 0 {
		params["shopId"] = fmt.Sprint(shopId) //店铺Id
	}

	resp, err := client.Execute("jd.union.open.goods.query", params)
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var goodsList = []GoodsResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	goodsResp := resp.Get("data")
	if respCode == 200 && goodsResp.Exists() && goodsResp.IsArray() {
		for _, value := range goodsResp.Array() {
			goods := GoodsResp{}

			//商品基本信息
			goods.Comments = value.Get("comments").Int()                     //评论数
			goods.GoodCommentsShare = value.Get("goodCommentsShare").Float() //商品好评率
			goods.InOrderCount30Days = value.Get("inOrderCount30Days").Int() //30天引单数量
			goods.MaterialUrl = value.Get("materialUrl").String()            //商品落地页
			goods.SkuId = value.Get("skuId").Int()                           //商品ID
			goods.SkuName = value.Get("skuName").String()                    //商品名称
			goods.IsHot = value.Get("isHot").Int()                           //是否爆款，1：是，0：否
			goods.Spuid = value.Get("spuid").Int()                           //其值为同款商品的主skuid
			goods.BrandCode = value.Get("brandCode").String()                //品牌code
			goods.BrandName = value.Get("brandName").String()                //品牌名
			goods.Owner = value.Get("owner").String()                        //g=自营，p=pop
			goods.TotalCount = value.Get("totalCount").Int()                 //有效商品总数量

			//类目信息
			categoryInfoResp := value.Get("categoryInfo")
			if categoryInfoResp.Exists() {
				goods.CategoryInfo.Cid1 = categoryInfoResp.Get("cid1").Int()
				goods.CategoryInfo.Cid1Name = categoryInfoResp.Get("cid1Name").String()
				goods.CategoryInfo.Cid2 = categoryInfoResp.Get("cid2").Int()
				goods.CategoryInfo.Cid2Name = categoryInfoResp.Get("cid2Name").String()
				goods.CategoryInfo.Cid3 = categoryInfoResp.Get("cid3").Int()
				goods.CategoryInfo.Cid3Name = categoryInfoResp.Get("cid3Name").String()
			}

			//图片信息
			imageListResp := value.Get("imageInfo.imageList")
			if imageListResp.Exists() && imageListResp.IsArray() {
				goods.SkuImage = imageListResp.Array()[0].Get("url").String() //自己加的。取第一张作为主图
				for _, v := range imageListResp.Array() {
					goods.ImageInfo.imageList = append(goods.ImageInfo.imageList, v.Get("url").String())
				}
			}

			//佣金信息
			commissionInfoResp := value.Get("commissionInfo")
			if commissionInfoResp.Exists() {
				goods.CommissionInfo.Commission = commissionInfoResp.Get("commission").Float()
				goods.CommissionInfo.CommissionShare = commissionInfoResp.Get("commissionShare").Float()
			}

			//价格信息
			priceInfoResp := value.Get("priceInfo")
			if priceInfoResp.Exists() {
				goods.PriceInfo.Price = priceInfoResp.Get("price").Float()
				goods.PriceInfo.LowestPrice = priceInfoResp.Get("lowestPrice").Float()
				goods.PriceInfo.LowestPriceType = priceInfoResp.Get("lowestPriceType").Int() //最低价格类型，1：无线价格；2：拼购价格； 3：秒杀价格
			}

			//店铺信息
			shopInfoResp := value.Get("shopInfo")
			if shopInfoResp.Exists() {
				goods.ShopInfo.ShopId = shopInfoResp.Get("shopId").Int()
				goods.ShopInfo.ShopName = shopInfoResp.Get("shopName").String()
			}

			//优惠券信息，返回内容为空说明该SKU无可用优惠券
			couponListResp := value.Get("couponInfo.couponList")
			if couponListResp.Exists() && couponListResp.IsArray() {
				for _, v := range couponListResp.Array() {
					coupon := Coupon{}
					coupon.BindType = v.Get("bindType").Int()         //券种类 (优惠券种类：0 - 全品类，1 - 限品类（自营商品），2 - 限店铺，3 - 店铺限商品券)
					coupon.Discount = v.Get("discount").Float()       //券面额
					coupon.Link = v.Get("link").String()              //券链接
					coupon.PlatformType = v.Get("platformType").Int() //券使用平台 (平台类型：0 - 全平台券，1 - 限平台券)
					coupon.Quota = v.Get("quota").Float()             //券消费限额
					coupon.GetStartTime = v.Get("getStartTime").Int() //领取开始时间(时间戳，毫秒)
					coupon.GetEndTime = v.Get("getEndTime").Int()     //券领取结束时间(时间戳，毫秒)
					coupon.UseStartTime = v.Get("useStartTime").Int() //券有效使用开始时间(时间戳，毫秒)
					coupon.UseEndTime = v.Get("useEndTime").Int()     //券有效使用结束时间(时间戳，毫秒)
					coupon.IsBest = v.Get("isBest").Int()             //最优优惠券，1：是；0：否
					goods.CouponInfo.CouponList = append(goods.CouponInfo.CouponList, coupon)
				}
			}

			//拼购信息
			pinGouInfoResp := value.Get("pinGouInfo")
			if pinGouInfoResp.Exists() {
				goods.PinGouInfo.PingouPrice = pinGouInfoResp.Get("pingouPrice").Float()
				goods.PinGouInfo.PingouTmCount = pinGouInfoResp.Get("pingouTmCount").Int()
				goods.PinGouInfo.PingouUrl = pinGouInfoResp.Get("pingouUrl").String()
				goods.PinGouInfo.PingouStartTime = pinGouInfoResp.Get("pingouStartTime").Int()
				goods.PinGouInfo.PingouEndTime = pinGouInfoResp.Get("pingouEndTime").Int()
			}

			goodsList = append(goodsList, goods)
		}
		return goodsList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.goods.link.query 链接商品查询接口【申请】
func (client *ApiReq) AdvancedGetLinkGoods(goodsUrl, subUnionId string) ([]LinkGoodsResp, *ApiErrorInfo) {

	client.ParamName = "goodsReq"

	resp, err := client.Execute("jd.union.open.goods.query", ApiParams{
		"goodsUrl":   goodsUrl,
		"subUnionId": subUnionId,
	})
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var goodsList = []LinkGoodsResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	linkGoodsResp := resp.Get("data")
	if respCode == 200 && linkGoodsResp.Exists() && linkGoodsResp.IsArray() {
		for _, value := range linkGoodsResp.Array() {
			goods := LinkGoodsResp{}

			goods.SkuId = value.Get("skuId").Int()         //skuId
			goods.ProductId = value.Get("productId").Int() //productId
			goods.Images = value.Get("images").String()    //图片集，逗号','分割，首张为主图
			if goods.Images != "" {
				images := strings.Split(goods.Images, ",")
				if len(images) > 0 {
					goods.SkuImage = images[0]
				}
			}
			goods.SkuName = value.Get("skuName").String()   //商品名称
			goods.Price = value.Get("price").Float()        //京东价，单位：元
			goods.CosRatio = value.Get("cosRatio").Float()  //佣金比例，单位：%
			goods.ShortUrl = value.Get("shortUrl").String() //短链接
			goods.ShopId = value.Get("shopId").String()     //店铺id
			goods.ShopName = value.Get("shopName").String() //店铺名称
			goods.Sales = value.Get("sales").Int()          //30天引单量
			goods.IsSelf = value.Get("isSelf").String()     //是否自营，g：自营，p：pop

			goodsList = append(goodsList, goods)
		}
		return goodsList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.promotion.byunionid.get 通过unionId获取推广链接【申请】
func (client *ApiReq) AdvancedGetUnionPromotionUrl(materialId string, unionId, positionId int64, pid, couponUrl, subUnionId string, chainType int64) (*PromotionCodeResp, *ApiErrorInfo) {
	client.ParamName = "promotionCodeReq"
	params := ApiParams{}

	params["materialId"] = materialId       //推广物料链接，建议链接使用微Q前缀，能较好适配微信手Q页面
	params["unionId"] = fmt.Sprint(unionId) //目标推客的联盟ID
	if positionId > 0 {
		params["positionId"] = fmt.Sprint(positionId) //新增推广位id （不填的话，为其默认生成一个唯一此接口推广位-名称：微信手Q短链接）
	}
	if pid != "" {
		params["pid"] = pid //子帐号身份标识，格式为子站长ID_子站长网站ID_子站长推广位ID
	}
	if couponUrl != "" {
		params["couponUrl"] = couponUrl //优惠券领取链接，在使用优惠券、商品二合一功能时入参，且materialId须为商品详情页链接
	}
	if subUnionId != "" {
		params["subUnionId"] = subUnionId //子联盟ID（需要联系运营开通权限才能拿到数据）
	}
	if chainType > 0 {
		params["chainType"] = fmt.Sprint(chainType) //转链类型，1：长链， 2 ：短链 ，3： 长链+短链，默认短链
	}
	resp, err := client.Execute("jd.union.open.promotion.byunionid.get", params)
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	promotionResp := resp.Get("data")
	if respCode == 200 && promotionResp.Exists() {
		promotion := &PromotionCodeResp{}
		promotion.ShortURL = promotionResp.Get("shortURL").String()
		promotion.ClickURL = promotionResp.Get("clickURL").String()
		return promotion, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.promotion.bysubunionid.get 通过subUnionId获取推广链接【申请】
func (client *ApiReq) AdvancedGetSubUnionPromotionUrl(materialId string, unionId, positionId int64, pid, couponUrl, subUnionId string, chainType int64) (*PromotionCodeResp, *ApiErrorInfo) {
	client.ParamName = "promotionCodeReq"
	params := ApiParams{}

	params["materialId"] = materialId //推广物料链接，建议链接使用微Q前缀，能较好适配微信手Q页面
	if subUnionId != "" {
		params["subUnionId"] = subUnionId //子联盟ID（需要联系运营开通权限才能拿到数据）
	}
	if positionId > 0 {
		params["positionId"] = fmt.Sprint(positionId) //新增推广位id （不填的话，为其默认生成一个唯一此接口推广位-名称：微信手Q短链接）
	}
	if pid != "" {
		params["pid"] = pid //子帐号身份标识，格式为子站长ID_子站长网站ID_子站长推广位ID
	}
	if couponUrl != "" {
		params["couponUrl"] = couponUrl //优惠券领取链接，在使用优惠券、商品二合一功能时入参，且materialId须为商品详情页链接
	}
	if chainType > 0 {
		params["chainType"] = fmt.Sprint(chainType) //转链类型，1：长链， 2 ：短链 ，3： 长链+短链，默认短链
	}
	resp, err := client.Execute("jd.union.open.promotion.bysubunionid.get", params)
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	promotionResp := resp.Get("data")
	if respCode == 200 && promotionResp.Exists() {
		promotion := &PromotionCodeResp{}
		promotion.ShortURL = promotionResp.Get("shortURL").String() //生成的推广目标链接，以短链接形式，有效期60天
		promotion.ClickURL = promotionResp.Get("clickURL").String() //生成推广目标的长链，长期有效
		return promotion, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.goods.bigfield.query 大字段商品查询接口（内测版）【申请】
func (client *ApiReq) AdvancedGetGoodsBigField(skuIds, fields interface{}) (*BigFieldGoodsResp, *ApiErrorInfo) {
	client.ParamName = "goodsReq"
	params := ApiParams{}
	apiErrInfo := &ApiErrorInfo{}

	if skuIds == nil {
		apiErrInfo.Code = 78
		apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
		return nil, apiErrInfo
	}

	//skuId集合
	switch result := skuIds.(type) {
	case string:
		params["skuIds"] = result
	case []int64:
		tmp, ids := "", ""
		for _, v := range result {
			ids = tmp + fmt.Sprint(v)
			tmp = ","
		}
		if ids != "" {
			params["skuIds"] = ids
		}
	case []string:
		ids := strings.Join(result, ",")
		if ids != "" {
			params["skuIds"] = ids
		}
	}

	if fields != nil {
		//查询域集合，不填写则查询全部
		switch result := skuIds.(type) {
		case string:
			params["fields"] = result
		case []int64:
			tmp, f := "", ""
			for _, v := range result {
				f = tmp + fmt.Sprint(v)
				tmp = ","
			}
			if f != "" {
				params["fields"] = f
			}
		case []string:
			f := strings.Join(result, ",")
			if f != "" {
				params["fields"] = f
			}
		}
	}

	resp, err := client.Execute("jd.union.open.goods.bigfield.query", params)
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	bigFieldResp := resp.Get("data")
	if respCode == 200 && bigFieldResp.Exists() {
		bigField := &BigFieldGoodsResp{}

		//类目信息
		categoryInfoResp := bigFieldResp.Get("categoryInfo")
		if categoryInfoResp.Exists() {
			bigField.CategoryInfo.Cid1 = categoryInfoResp.Get("cid1").Int()
			bigField.CategoryInfo.Cid1Name = categoryInfoResp.Get("cid1Name").String()
			bigField.CategoryInfo.Cid2 = categoryInfoResp.Get("cid2").Int()
			bigField.CategoryInfo.Cid2Name = categoryInfoResp.Get("cid2Name").String()
			bigField.CategoryInfo.Cid3 = categoryInfoResp.Get("cid3").Int()
			bigField.CategoryInfo.Cid3Name = categoryInfoResp.Get("cid3Name").String()
		}

		//图片信息
		imageListResp := bigFieldResp.Get("imageInfo.imageList")
		if imageListResp.Exists() && imageListResp.IsArray() {
			for _, v := range imageListResp.Array() {
				bigField.ImageInfo.imageList = append(bigField.ImageInfo.imageList, v.Get("url").String())
			}
		}

		//基础大字段信息
		baseBigFieldInfoResp := bigFieldResp.Get("baseBigFieldInfo")
		if baseBigFieldInfoResp.Exists() {
			bigField.BaseBigFieldInfo.Wdis = baseBigFieldInfoResp.Get("wdis").String()
			bigField.BaseBigFieldInfo.PropCode = baseBigFieldInfoResp.Get("propCode").String()
			bigField.BaseBigFieldInfo.WareQD = baseBigFieldInfoResp.Get("wareQD").String()
		}

		//图书大字段信息
		bookBigFieldInfoResp := bigFieldResp.Get("bookBigFieldInfo")
		if baseBigFieldInfoResp.Exists() {
			bigField.BookBigFieldInfo.Comments = bookBigFieldInfoResp.Get("comments").String()               //媒体评论
			bigField.BookBigFieldInfo.Image = bookBigFieldInfoResp.Get("image").String()                     //精彩文摘与插图(插图)
			bigField.BookBigFieldInfo.ContentDesc = bookBigFieldInfoResp.Get("contentDesc").String()         //内容摘要(内容简介)
			bigField.BookBigFieldInfo.RelatedProducts = bookBigFieldInfoResp.Get("relatedProducts").String() //产品描述(相关商品)
			bigField.BookBigFieldInfo.EditerDesc = bookBigFieldInfoResp.Get("editerDesc").String()           //编辑推荐
			bigField.BookBigFieldInfo.Catalogue = bookBigFieldInfoResp.Get("catalogue").String()             //目录
			bigField.BookBigFieldInfo.BookAbstract = bookBigFieldInfoResp.Get("bookAbstract").String()       //精彩摘要(精彩书摘)
			bigField.BookBigFieldInfo.AuthorDesc = bookBigFieldInfoResp.Get("authorDesc").String()           //作者简介
			bigField.BookBigFieldInfo.Introduction = bookBigFieldInfoResp.Get("introduction").String()       //前言(前言/序言)
			bigField.BookBigFieldInfo.ProductFeatures = bookBigFieldInfoResp.Get("productFeatures").String() //产品特色
		}

		//基础大字段信息
		videoBigFieldInfoResp := bigFieldResp.Get("videoBigFieldInfo")
		if videoBigFieldInfoResp.Exists() {
			bigField.VideoBigFieldInfo.Comments = videoBigFieldInfoResp.Get("comments").String()                        //评论
			bigField.VideoBigFieldInfo.Image = videoBigFieldInfoResp.Get("image").String()                              //商品描述(精彩剧照)
			bigField.VideoBigFieldInfo.ContentDesc = videoBigFieldInfoResp.Get("contentDesc").String()                  //内容摘要(内容简介)
			bigField.VideoBigFieldInfo.EditerDesc = videoBigFieldInfoResp.Get("editerDesc").String()                    //编辑推荐
			bigField.VideoBigFieldInfo.Catalogue = videoBigFieldInfoResp.Get("catalogue").String()                      //目录
			bigField.VideoBigFieldInfo.BoxContents = videoBigFieldInfoResp.Get("box_Contents").String()                 //包装清单
			bigField.VideoBigFieldInfo.MaterialDescription = videoBigFieldInfoResp.Get("material_Description").String() //特殊说明
			bigField.VideoBigFieldInfo.Manual = videoBigFieldInfoResp.Get("manual").String()                            //说明书
			bigField.VideoBigFieldInfo.ProductFeatures = videoBigFieldInfoResp.Get("productFeatures").String()          //产品特色
		}
		return bigField, nil
	}

	errInfo := ApiErrorInfo{}
	errInfo.Code = respCode
	if e, ok := ApiErrInfo[errInfo.Code]; ok {
		errInfo.Message = e.Error()
	} else {
		errInfo.Message = resp.Get("message").String()
	}
	return nil, &errInfo
}

//jd.union.open.coupon.query 优惠券领取情况查询接口【申请】
func (client *ApiReq) AdvancedCouponQuery(couponUrls interface{}) ([]CouponResp, *ApiErrorInfo) {
	client.ParamName = ""
	params := ApiParams{}
	apiErrInfo := &ApiErrorInfo{}

	if couponUrls == nil {
		apiErrInfo.Code = 78
		apiErrInfo.Message = ApiErrInfo[apiErrInfo.Code].Error()
		return nil, apiErrInfo
	}

	//couponUrls 集合
	switch result := couponUrls.(type) {
	case string:
		params["couponUrls"] = result
	case []int64:
		tmp, ids := "", ""
		for _, v := range result {
			ids = tmp + fmt.Sprint(v)
			tmp = ","
		}
		if ids != "" {
			params["couponUrls"] = ids
		}
	case []string:
		ids := strings.Join(result, ",")
		if ids != "" {
			params["couponUrls"] = ids
		}
	}

	resp, err := client.Execute("jd.union.open.coupon.query", params)
	if err != nil {
		return nil, err
	}
	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}
	respCode := resp.Get("code").Int()
	couponResp := resp.Get("data")
	if respCode == 200 && couponResp.Exists() && couponResp.IsArray() {
		couponList := []CouponResp{}
		for _, value := range couponResp.Array() {
			coupon := CouponResp{}
			coupon.TakeEndTime = value.Get("takeEndTime").Int()     //券领取结束时间(时间戳，毫秒)
			coupon.TakeBeginTime = value.Get("takeBeginTime").Int() //券领取开始时间(时间戳，毫秒)
			coupon.RemainNum = value.Get("remainNum").Int()         //券剩余张数
			coupon.Yn = value.Get("yn").String()                    //券有效状态
			coupon.Num = value.Get("num").Int()                     //券总张数
			coupon.Quota = value.Get("quota").Float()               //券消费限额
			coupon.Link = value.Get("link").String()                //券链接
			coupon.Discount = value.Get("discount").Float()         //券面额
			coupon.BeginTime = value.Get("beginTime").Int()         //券有效使用开始时间(时间戳，毫秒)
			coupon.EndTime = value.Get("endTime").Int()             //券有效使用结束时间(时间戳，毫秒)
			coupon.Platform = value.Get("platform").String()        //券使用平台
			couponList = append(couponList, coupon)
		}
		return couponList, nil
	}

	errInfo := ApiErrorInfo{}
	errInfo.Code = respCode
	if e, ok := ApiErrInfo[errInfo.Code]; ok {
		errInfo.Message = e.Error()
	} else {
		errInfo.Message = resp.Get("message").String()
	}
	return nil, &errInfo
}

//jd.union.open.goods.seckill.query 秒杀商品查询接口【申请】
func (client *ApiReq) AdvancedGetSecKillGoods(skuIds interface{}, pageIndex, pageSize, isBeginSecKill, cid1, cid2, cid3 int64,
	secKillPriceFrom, secKillPriceTo, commissionShareFrom, commissionShareTo float64, owner, sortName, sort string) ([]SecKillGoodsResp, *ApiErrorInfo) {

	client.ParamName = "goodsReq"
	params := ApiParams{}
	if cid1 > 0 {
		params["cid1"] = fmt.Sprint(cid1) //一级类目id
	}
	if cid2 > 0 {
		params["cid2"] = fmt.Sprint(cid2) //二级类目id
	}
	if cid3 > 0 {
		params["cid3"] = fmt.Sprint(cid3) //一级类目id
	}
	params["pageIndex"] = fmt.Sprint(pageIndex) //页码，返回第几页结果
	if pageSize <= 0 {
		params["pageSize"] = "30" //每页数量，单页数最大30，默认20
	} else {
		params["pageSize"] = fmt.Sprint(pageSize)
	}

	//skuid集合(一次最多支持查询100个sku)，英文逗号分割，数组类型开发时记得加[]
	if skuIds != nil {
		switch result := skuIds.(type) {
		case string:
			params["skuIds"] = result
		case []int64:
			tmp, ids := "", ""
			for _, v := range result {
				ids = tmp + fmt.Sprint(v)
				tmp = ","
			}
			if ids != "" {
				params["skuIds"] = ids
			}
		case []string:
			ids := strings.Join(result, ",")
			if ids != "" {
				params["skuIds"] = ids
			}
		}
	}

	if isBeginSecKill > 0 {
		params["isBeginSecKill"] = "1" //是否返回未开始秒杀商品。1=返回，0=不返回
	}
	if secKillPriceFrom > 0 {
		params["secKillPriceFrom"] = fmt.Sprint(secKillPriceFrom) //秒杀价区间开始（单位：元）
	}
	if secKillPriceTo > 0 {
		params["secKillPriceTo"] = fmt.Sprint(secKillPriceTo) //秒杀价区间结束
	}

	if commissionShareFrom > 0 {
		params["commissionShareFrom"] = fmt.Sprint(commissionShareFrom) //佣金比例区间开始
	}
	if commissionShareTo > 0 {
		params["commissionShareTo"] = fmt.Sprint(commissionShareTo) //佣金比例区间结束
	}
	if owner != "" {
		params["owner"] = owner //商品类型：自营[g]，POP[p]
	}
	if sortName != "" {
		params["sortName"] = sortName //排序字段(price：单价, commissionShare：佣金比例, commission：佣金， inOrderCount30Days：30天引单量， inOrderComm30Days：30天支出佣金)
	}
	if sort != "" {
		params["sort"] = sort //asc,desc升降序,默认降序
	}

	resp, err := client.Execute("jd.union.open.goods.seckill.query", params)
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var goodsList = []SecKillGoodsResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	goodsResp := resp.Get("data")
	if respCode == 200 && goodsResp.Exists() && goodsResp.IsArray() {
		for _, value := range goodsResp.Array() {
			goods := SecKillGoodsResp{}

			goods.SkuName = value.Get("skuName").String()                    //商品名称
			goods.SkuId = value.Get("skuId").Int()                           //商品id
			goods.ImageUrl = value.Get("imageUrl").String()                  //图片url
			goods.IsSecKill = value.Get("isSecKill").Int()                   //是秒杀。1：是商品 0：非秒杀商品
			goods.OriPrice = value.Get("oriPrice").Float()                   //原价
			goods.SecKillPrice = value.Get("secKillPrice").Float()           //秒杀价
			goods.SecKillStartTime = value.Get("secKillStartTime").Int()     //秒杀开始展示时间（时间戳：毫秒）
			goods.SecKillEndTime = value.Get("secKillEndTime").Int()         //秒杀结束时间（时间戳：毫秒）
			goods.Cid1Id = value.Get("cid1Id").Int()                         //一级类目id
			goods.Cid2Id = value.Get("cid2Id").Int()                         //二级类目id
			goods.Cid3Id = value.Get("cid3Id").Int()                         //三级类目id
			goods.Cid1Name = value.Get("cid1Name").String()                  //一级类目名称
			goods.Cid2Name = value.Get("cid2Name").String()                  //二级类目名称
			goods.Cid3Name = value.Get("cid3Name").String()                  //三级类目名称
			goods.CommissionShare = value.Get("commissionShare").Float()     //通用佣金比例，百分比
			goods.Commission = value.Get("commission").Float()               //通用佣金
			goods.Owner = value.Get("owner").String()                        //是否自营。g=自营，p=pop
			goods.InOrderCount30Days = value.Get("inOrderCount30Days").Int() //30天引入订单量（spu）
			goods.InOrderComm30Days = value.Get("inOrderComm30Days").Float() //30天支出佣金（spu）
			goods.JdPrice = value.Get("jdPrice").Float()                     //京东价
			goods.TotalCount = value.Get("totalCount").Int()                 //总数量

			goodsList = append(goodsList, goods)
		}
		return goodsList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.goods.stuprice.query 学生价商品查询接口【申请】
func (client *ApiReq) AdvancedGetStudentGoods(skuIds interface{}, pageIndex, pageSize, cid1, cid2, cid3 int64,
	stuPriceFrom, stuPriceTo, commissionShareFrom, commissionShareTo float64, owner, sortName, sort string) ([]StuPriceGoodsResp, *ApiErrorInfo) {

	client.ParamName = "goodsReq"
	params := ApiParams{}
	if cid1 > 0 {
		params["cid1"] = fmt.Sprint(cid1) //一级类目id
	}
	if cid2 > 0 {
		params["cid2"] = fmt.Sprint(cid2) //二级类目id
	}
	if cid3 > 0 {
		params["cid3"] = fmt.Sprint(cid3) //一级类目id
	}
	params["pageIndex"] = fmt.Sprint(pageIndex) //页码，返回第几页结果
	if pageSize <= 0 {
		params["pageSize"] = "30" //每页数量，单页数最大30，默认20
	} else {
		params["pageSize"] = fmt.Sprint(pageSize)
	}

	//skuid集合(一次最多支持查询100个sku)，英文逗号分割，数组类型开发时记得加[]
	if skuIds != nil {
		switch result := skuIds.(type) {
		case string:
			params["skuIds"] = result
		case []int64:
			tmp, ids := "", ""
			for _, v := range result {
				ids = tmp + fmt.Sprint(v)
				tmp = ","
			}
			if ids != "" {
				params["skuIds"] = ids
			}
		case []string:
			ids := strings.Join(result, ",")
			if ids != "" {
				params["skuIds"] = ids
			}
		}
	}

	if stuPriceFrom > 0 {
		params["stuPriceFrom"] = fmt.Sprint(stuPriceFrom) //学生专享价区间开始（单位：元）
	}
	if stuPriceTo > 0 {
		params["stuPriceTo"] = fmt.Sprint(stuPriceTo) //学生专享价区间结束（单位：元）
	}

	if commissionShareFrom > 0 {
		params["commissionShareFrom"] = fmt.Sprint(commissionShareFrom) //佣金比例区间开始
	}
	if commissionShareTo > 0 {
		params["commissionShareTo"] = fmt.Sprint(commissionShareTo) //佣金比例区间结束
	}
	if owner != "" {
		params["owner"] = owner //商品类型：自营[g]，POP[p]
	}
	if sortName != "" {
		params["sortName"] = sortName //排序字段(price：单价, commissionShare：佣金比例, commission：佣金， inOrderCount30Days：30天引单量， inOrderComm30Days：30天支出佣金)
	}
	if sort != "" {
		params["sort"] = sort //asc,desc升降序,默认降序
	}

	resp, err := client.Execute("jd.union.open.goods.stuprice.query", params)
	if err != nil {
		return nil, err
	}

	if apiErrInfo := client.CheckApiErr(resp); apiErrInfo != nil {
		return nil, apiErrInfo
	}

	var goodsList = []StuPriceGoodsResp{} //创建一个以goodsId为索引的字典
	respCode := resp.Get("code").Int()
	goodsResp := resp.Get("data")
	if respCode == 200 && goodsResp.Exists() && goodsResp.IsArray() {
		for _, value := range goodsResp.Array() {
			goods := StuPriceGoodsResp{}

			goods.SkuName = value.Get("skuName").String()                    //商品名称
			goods.SkuId = value.Get("skuId").Int()                           //商品id
			goods.ImageUrl = value.Get("imageUrl").String()                  //图片url
			goods.IsStuPrice = value.Get("isStuPrice").Int()                 //是否学生价商品。 1：是学生价商品。 0：不是学生价商品。
			goods.JdPrice = value.Get("jdPrice").Float()                     //京东价
			goods.StudentPrice = value.Get("studentPrice").Float()           //学生专享价
			goods.StuPriceStartTime = value.Get("stuPriceStartTime").Int()   //专享价促销开始时间（时间戳：毫秒）
			goods.StuPriceEndTime = value.Get("stuPriceEndTime").Int()       //专享价促销结束时间（时间戳：毫秒）
			goods.Cid1Id = value.Get("cid1Id").Int()                         //一级类目id
			goods.Cid2Id = value.Get("cid2Id").Int()                         //二级类目id
			goods.Cid3Id = value.Get("cid3Id").Int()                         //三级类目id
			goods.Cid1Name = value.Get("cid1Name").String()                  //一级类目名称
			goods.Cid2Name = value.Get("cid2Name").String()                  //二级类目名称
			goods.Cid3Name = value.Get("cid3Name").String()                  //三级类目名称
			goods.CommissionShare = value.Get("commissionShare").Float()     //通用佣金比例，百分比
			goods.Commission = value.Get("commission").Float()               //通用佣金
			goods.Owner = value.Get("owner").String()                        //是否自营。g=自营，p=pop
			goods.InOrderCount30Days = value.Get("inOrderCount30Days").Int() //30天引入订单量（spu）
			goods.InOrderComm30Days = value.Get("inOrderComm30Days").Float() //30天支出佣金（spu）
			goods.TotalCount = value.Get("totalCount").Int()                 //总数量

			goodsList = append(goodsList, goods)
		}
		return goodsList, nil
	} else {
		errInfo := ApiErrorInfo{}
		errInfo.Code = respCode
		if e, ok := ApiErrInfo[errInfo.Code]; ok {
			errInfo.Message = e.Error()
		} else {
			errInfo.Message = resp.Get("message").String()
		}
		return nil, &errInfo
	}
}

//jd.union.open.position.create 创建推广位【申请】
//jd.union.open.position.query 查询推广位【申请】
//jd.union.open.coupon.importation 优惠券导入【申请】
