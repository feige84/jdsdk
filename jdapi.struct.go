package jdsdk

//商品类目查询
type CategoryResp struct {
	Id       int64  `json:"id"`        //类目Id
	Name     string `json:"name"`      //类目名称
	Grade    int64  `json:"grade"`     //类目级别(类目级别 0，1，2 代表一、二、三级类目)
	ParentId int64  `json:"parent_id"` //父类目Id
}

//订单查询接口
type OrderResp struct {
	FinishTime int64     //订单完成时间(时间戳，毫秒)
	OrderEmt   int64     //下单设备(1:PC, 2:无线)
	OrderId    int64     //订单ID
	OrderTime  int64     //下单时间(时间戳，毫秒)
	ParentId   int64     //父单的订单ID，仅当发生订单拆分时返回， 0：未拆分，有值则表示此订单为子订单
	PayMonth   string    //订单维度预估结算时间（格式：yyyyMMdd），0：未结算，订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
	Plus       int64     //下单用户是否为PLUS会员 0：否，1：是
	PopId      int64     //商家ID
	SkuList    []SkuInfo //订单包含的商品信息列表
	UnionId    int64     //推客的联盟ID
	Ext1       string    //推客生成推广链接时传入的扩展字段，订单维度（需要联系运营开放白名单才能拿到数据）
	ValidCode  int64     //订单维度的有效码（-1：未知, 2.无效-拆单, 3.无效-取消, 4.无效-京东帮帮主订单, 5.无效-账号异常, 6.无效-赠品类目不返佣, 7.无效-校园订单, 8.无效-企业订单, 9.无效-团购订单, 10.无效-开增值税专用发票订单, 11.无效-乡村推广员下单, 12.无效-自己推广自己下单, 13.无效-违规订单, 14.无效-来源与备案网址不符, 15.待付款, 16.已付款, 17.已完成, 18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
	HasMore    bool      //是否还有更多, true：还有数据；false:已查询完毕，没有数据
}

type SkuInfo struct {
	ActualCosPrice    float64 //实际计算佣金的金额。订单完成后，会将误扣除的运费券金额更正。如订单完成后发生退款，此金额会更新。
	ActualFee         float64 //推客获得的实际佣金（实际计佣金额*佣金比例*最终比例）。如订单完成后发生退款，此金额会更新。
	CommissionRate    float64 //佣金比例
	EstimateCosPrice  float64 //预估计佣金额，即用户下单的金额(已扣除优惠券、白条、支付优惠、进口税，未扣除红包和京豆)，有时会误扣除运费券金额，完成结算时会在实际计佣金额中更正。如订单完成前发生退款，此金额不会更新。
	EstimateFee       float64 //推客的预估佣金（预估计佣金额*佣金比例*最终比例），如订单完成前发生退款，此金额不会更新。
	FinalRate         float64 //最终比例（分成比例+补贴比例）
	Cid1              int64   //一级类目ID
	Cid2              int64   //二级类目ID
	Cid3              int64   //三级类目ID
	FrozenSkuNum      int64   //商品售后中数量
	Pid               string  //联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID
	PositionId        int64   //推广位ID,0代表无推广位
	Price             float64 //商品单价
	SiteId            int64   //网站ID，0：无网站
	SkuId             int64   //商品ID
	SkuName           string  //商品名称
	SkuNum            int64   //商品数量
	SkuReturnNum      int64   //商品已退货数量
	SubSideRate       float64 //分成比例
	SubsidyRate       float64 //补贴比例
	UnionAlias        string  //PID所属母账号平台名称（原第三方服务商来源）
	UnionTag          string  //联盟标签数据（整型的二进制字符串(32位)，目前只返回8位：00000001。数据从右向左进行，每一位为1表示符合联盟的标签特征，第1位：京喜红包，第2位：组合推广订单，第3位：拼购订单，第5位：有效首次购订单（00011XXX表示有效首购，最终奖励活动结算金额会结合订单状态判断，以联盟后台对应活动效果数据报表https://union.jd.com/active为准）。例如：00000001:京喜红包订单，00000010:组合推广订单，00000100:拼购订单，00011000:有效首购，00000111：京喜红包+组合推广+拼购等）
	UnionTrafficGroup int64   //渠道组 1：1号店，其他：京东
	ValidCode         int64   //sku维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
	SubUnionId        string  //子联盟ID(需要联系运营开放白名单才能拿到数据)
	TraceType         int64   //2：同店；3：跨店
	PayMonth          string  //订单行维度预估结算时间（格式：yyyyMMdd） ，0：未结算。订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
	PopId             int64   //商家ID，订单行维度
	Ext1              string  //推客生成推广链接时传入的扩展字段（需要联系运营开放白名单才能拿到数据）。&lt;订单行维度&gt;
}

//关键词商品查询接口【申请】
type GoodsResp struct {
	CategoryInfo       CategoryInfo   //类目信息
	Comments           int64          //评论数
	CommissionInfo     CommissionInfo //佣金信息
	CouponInfo         CouponInfo     //优惠券信息，返回内容为空说明该SKU无可用优惠券
	GoodCommentsShare  float64        //商品好评率
	ImageInfo          ImageInfo      //图片信息
	InOrderCount30Days int64          //30天引单数量
	MaterialUrl        string         //商品落地页
	PriceInfo          PriceInfo      //价格信息
	ShopInfo           ShopInfo       //店铺信息
	SkuId              int64          //商品ID
	SkuName            string         //商品名称
	SkuImage           string         //自己加的。直接取ImageInfo[0]
	IsHot              int64          //是否爆款，1：是，0：否
	Spuid              int64          //其值为同款商品的主skuid
	BrandCode          string         //品牌code
	BrandName          string         //品牌名
	Owner              string         //g=自营，p=pop
	PinGouInfo         PinGouInfo     //拼购信息
	TotalCount         int64          //有效商品总数量
}

type CategoryInfo struct {
	Cid1     int64  //一级类目ID
	Cid1Name string //一级类目名称
	Cid2     int64  //二级类目ID
	Cid2Name string //二级类目名称
	Cid3     int64  //三级类目ID
	Cid3Name string //三级类目名称
}

type CommissionInfo struct {
	Commission      float64 //佣金
	CommissionShare float64 //佣金比例
}

type CouponInfo struct {
	CouponList []Coupon //优惠券合集
}

type Coupon struct {
	BindType     int64   //券种类 (优惠券种类：0 - 全品类，1 - 限品类（自营商品），2 - 限店铺，3 - 店铺限商品券)
	Discount     float64 //券面额
	Link         string  //券链接
	PlatformType int64   //券使用平台 (平台类型：0 - 全平台券，1 - 限平台券)
	Quota        float64 //券消费限额
	GetStartTime int64   //领取开始时间(时间戳，毫秒)
	GetEndTime   int64   //券领取结束时间(时间戳，毫秒)
	UseStartTime int64   //券有效使用开始时间(时间戳，毫秒)
	UseEndTime   int64   //券有效使用结束时间(时间戳，毫秒)
	IsBest       int64   //最优优惠券，1：是；0：否
}

type ImageInfo struct {
	imageList []string //图片合集 这里没有用官方的urlinfo结构。直接组合成数组。
}

type UrlInfo struct {
	Url string //图片链接地址，第一个图片链接为主图链接
}

type PriceInfo struct {
	Price           float64 //无线价格
	LowestPrice     float64 //最低价格
	LowestPriceType int64   //最低价格类型，1：无线价格；2：拼购价格； 3：秒杀价格
}

type ShopInfo struct {
	ShopName string //店铺名称（或供应商名称）
	ShopId   int64  //店铺Id
}

type PinGouInfo struct {
	PingouPrice     float64 //拼购价格
	PingouTmCount   int64   //拼购成团所需人数
	PingouUrl       string  //拼购落地页url
	PingouStartTime int64   //拼购开始时间(时间戳，毫秒)
	PingouEndTime   int64   //拼购结束时间(时间戳，毫秒)
}

type ResourceInfo struct {
	EliteId   int64  //频道id
	EliteName string //频道名称
}

//京粉精选商品查询接口
type JFGoodsResp struct {
	CategoryInfo          CategoryInfo   //类目信息
	Comments              int64          //评论数
	CommissionInfo        CommissionInfo //佣金信息
	CouponInfo            CouponInfo     //优惠券信息，返回内容为空说明该SKU无可用优惠券
	GoodCommentsShare     float64        //商品好评率
	ImageInfo             ImageInfo      //图片信息
	InOrderCount30Days    int64          //30天引单数量
	MaterialUrl           string         //商品落地页
	PriceInfo             PriceInfo      //价格信息
	ShopInfo              ShopInfo       //店铺信息
	SkuId                 int64          //商品ID
	SkuName               string         //商品名称
	SkuImage              string         //自己加的。直接取ImageInfo[0]
	IsHot                 int64          //是否爆款，1：是，0：否
	Spuid                 int64          //其值为同款商品的主skuid
	BrandCode             string         //品牌code
	BrandName             string         //品牌名
	Owner                 string         //g=自营，p=pop
	PinGouInfo            PinGouInfo     //拼购信息
	ResourceInfo          ResourceInfo   //资源信息
	InOrderCount30DaysSku int64          //30天引单数量(sku维度)
	TotalCount            int64          //有效商品总数量
}

//jd.union.open.goods.promotiongoodsinfo.query 获取推广商品信息接口
type PromotionGoodsResp struct {
	SkuId             int64   //商品ID
	UnitPrice         float64 //商品单价即京东价
	MaterialUrl       string  //商品落地页
	EndDate           int64   //推广结束日期(时间戳，毫秒)
	IsFreeFreightRisk int64   //是否支持运费险(1:是,0:否)
	IsFreeShipping    int64   //是否包邮(1:是,0:否,2:自营商品遵从主站包邮规则)
	CommisionRatioWl  float64 //无线佣金比例
	CommisionRatioPc  float64 //PC佣金比例
	ImgUrl            string  //图片地址
	Vid               int64   //商家ID
	Cid               int64   //一级类目ID
	CidName           string  //一级类目名称
	Cid2              int64   //二级类目ID
	Cid2Name          string  //二级类目名称
	Cid3              int64   //三级类目ID
	Cid3Name          string  //三级类目名称
	WlUnitPrice       float64 //商品无线京东价（单价为-1表示未查询到该商品单价）
	IsSeckill         int64   //是否秒杀(1:是,0:否)
	InOrderCount      int64   //30天引单数量
	ShopId            int64   //店铺ID
	IsJdSale          int64   //是否自营(1:是,0:否)
	GoodsName         string  //商品名称
	StartDate         int64   //推广开始日期（时间戳，毫秒）
}

//jd.union.open.promotion.common.get 获取通用推广链接
type PromotionCodeResp struct {
	ShortURL string //生成的推广目标链接，以短链接形式，有效期60天
	ClickURL string //生成的目标推广链接，长期有效
}

//jd.union.open.goods.link.query 链接商品查询接口【申请】
type LinkGoodsResp struct {
	SkuId     int64   //skuId
	ProductId int64   //productId
	Images    string  //图片集，逗号','分割，首张为主图
	SkuImage  string  //自己加的。直接取ImageInfo[0]
	SkuName   string  //商品名称
	Price     float64 //京东价，单位：元
	CosRatio  float64 //佣金比例，单位：%
	ShortUrl  string  //短链接
	ShopId    string  //店铺id
	ShopName  string  //店铺名称
	Sales     int64   //30天引单量
	IsSelf    string  //是否自营，g：自营，p：pop
}

//jd.union.open.goods.bigfield.query 大字段商品查询接口（内测版）【申请】
type BigFieldGoodsResp struct {
	SkuId             int64             //skuId
	SkuName           string            //商品名称
	CategoryInfo      CategoryInfo      //分类信息
	ImageInfo         ImageInfo         //图片信息
	BaseBigFieldInfo  BaseBigFieldInfo  //基础大字段信息
	BookBigFieldInfo  BookBigFieldInfo  //图书大字段信息
	VideoBigFieldInfo VideoBigFieldInfo //影音大字段信息
}

type BaseBigFieldInfo struct {
	Wdis     string //商品介绍
	PropCode string //规格参数
	WareQD   string //包装清单 (仅自营商品)
}

type BookBigFieldInfo struct {
	Comments        string //媒体评论
	Image           string //精彩文摘与插图(插图)
	ContentDesc     string //内容摘要(内容简介)
	RelatedProducts string //产品描述(相关商品)
	EditerDesc      string //编辑推荐
	Catalogue       string //目录
	BookAbstract    string //精彩摘要(精彩书摘)
	AuthorDesc      string //作者简介
	Introduction    string //前言(前言/序言)
	ProductFeatures string //产品特色
}

type VideoBigFieldInfo struct {
	Comments            string //评论
	Image               string //商品描述(精彩剧照)
	ContentDesc         string //内容摘要(内容简介)
	EditerDesc          string //编辑推荐
	Catalogue           string //目录
	BoxContents         string //包装清单
	MaterialDescription string //特殊说明
	Manual              string //说明书
	ProductFeatures     string //产品特色
}

//jd.union.open.coupon.query 优惠券领取情况查询接口【申请】
type CouponResp struct {
	TakeEndTime   int64   //券领取结束时间(时间戳，毫秒)
	TakeBeginTime int64   //券领取开始时间(时间戳，毫秒)
	RemainNum     int64   //券剩余张数
	Yn            string  //券有效状态
	Num           int64   //券总张数
	Quota         float64 //券消费限额
	Link          string  //券链接
	Discount      float64 //券面额
	BeginTime     int64   //券有效使用开始时间(时间戳，毫秒)
	EndTime       int64   //券有效使用结束时间(时间戳，毫秒)
	Platform      string  //券使用平台
}

//jd.union.open.goods.stuprice.query 学生价商品查询接口【申请】
type StuPriceGoodsResp struct {
	SkuName            string  //商品名称
	SkuId              int64   //商品id
	ImageUrl           string  //图片url
	IsStuPrice         int64   //是否学生价商品。 1：是学生价商品。 0：不是学生价商品。
	JdPrice            float64 //京东价
	StudentPrice       float64 //学生专享价
	StuPriceStartTime  int64   //专享价促销开始时间（时间戳：毫秒）
	StuPriceEndTime    int64   //专享价促销结束时间（时间戳：毫秒）
	Cid1Id             int64   //一级类目id
	Cid2Id             int64   //二级类目id
	Cid3Id             int64   //三级类目id
	Cid1Name           string  //一级类目名称
	Cid2Name           string  //二级类目名称
	Cid3Name           string  //三级类目名称
	CommissionShare    float64 //通用佣金比例，百分比
	Commission         float64 //通用佣金
	Owner              string  //是否自营。g=自营，p=pop
	InOrderCount30Days int64   //30天引入订单量（spu）
	InOrderComm30Days  float64 //30天支出佣金（spu）
	TotalCount         int64   //总数量
}

//jd.union.open.goods.seckill.query 秒杀商品查询接口【申请】
type SecKillGoodsResp struct {
	SkuName            string  //商品名称
	SkuId              int64   //商品id
	ImageUrl           string  //图片url
	IsSecKill          int64   //是秒杀。1：是商品 0：非秒杀商品
	OriPrice           float64 //原价
	SecKillPrice       float64 //秒杀价
	SecKillStartTime   int64   //秒杀开始展示时间（时间戳：毫秒）
	SecKillEndTime     int64   //秒杀结束时间（时间戳：毫秒）
	Cid1Id             int64   //一级类目id
	Cid2Id             int64   //二级类目id
	Cid3Id             int64   //三级类目id
	Cid1Name           string  //一级类目名称
	Cid2Name           string  //二级类目名称
	Cid3Name           string  //三级类目名称
	CommissionShare    float64 //通用佣金比例，百分比
	Commission         float64 //通用佣金
	Owner              string  //是否自营。g=自营，p=pop
	InOrderCount30Days int64   //30天引入订单量（spu）
	InOrderComm30Days  float64 //30天支出佣金（spu）
	JdPrice            float64 //京东价
	TotalCount         int64   //总数量
}
