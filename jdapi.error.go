package jdsdk

import "errors"

var (
	ApiErrInfo = make(map[int64]error)
)

func init() {
	ApiErrInfo[2] = errors.New("限制时间内调用失败次数")
	ApiErrInfo[4] = errors.New("缺少版本参数")
	ApiErrInfo[5] = errors.New("不支持的版本号")
	ApiErrInfo[7] = errors.New("缺少时间戳参数")
	ApiErrInfo[8] = errors.New("非法的时间戳参数")
	ApiErrInfo[11] = errors.New("缺少签名参数")
	ApiErrInfo[12] = errors.New("无效签名")
	ApiErrInfo[14] = errors.New("缺少方法名参数")
	ApiErrInfo[15] = errors.New("不存在的方法名")
	ApiErrInfo[18] = errors.New("缺少access_token参数")
	ApiErrInfo[19] = errors.New("无效access_token")
	ApiErrInfo[20] = errors.New("缺少app_key参数")
	ApiErrInfo[21] = errors.New("无效app_key")
	ApiErrInfo[24] = errors.New("无权调用API")
	ApiErrInfo[28] = errors.New("接口已停用")
	ApiErrInfo[43] = errors.New("系统处理错误")

	ApiErrInfo[61] = errors.New("参数｛0｝值不合法，请参照帮助文档确认！")
	ApiErrInfo[62] = errors.New("json转换时错误，错误的请求参数")
	ApiErrInfo[63] = errors.New("json格式不合法")
	ApiErrInfo[64] = errors.New("此类型商家无权调用本接口")
	ApiErrInfo[65] = errors.New("平台连接后端服务超时")
	ApiErrInfo[66] = errors.New("平台连接后端服务不可用")
	ApiErrInfo[67] = errors.New("平台连接后端服务处理过程中出现未知异常信息")
	ApiErrInfo[68] = errors.New("验证可选字段异常信息")
	ApiErrInfo[69] = errors.New("获取数据失败")

	ApiErrInfo[70] = errors.New("该订单正在出库中")
	ApiErrInfo[71] = errors.New("当前的ID不属于此商家")
	ApiErrInfo[72] = errors.New("当前的用户不是此类型（如FBP, SOP等）的商家")
	ApiErrInfo[73] = errors.New("该api是增值api，请将您的app入住云鼎平台方可调用")
	ApiErrInfo[74] = errors.New("获取调用这ip失败")
	ApiErrInfo[75] = errors.New("非法的调用这ip")
	ApiErrInfo[77] = errors.New("非法的http协议")
	ApiErrInfo[78] = errors.New("HTTP 调用异常")
	ApiErrInfo[79] = errors.New("该商家没有权限通过该appKey调用该api")

	ApiErrInfo[80] = errors.New("调用的IP不在APP设置的IP白名单中")
	ApiErrInfo[83] = errors.New("该api是解密api,所传参数不包含密文参数")
	ApiErrInfo[85] = errors.New("解密失败")
	ApiErrInfo[86] = errors.New("后台没有配置接口解密参数")
	ApiErrInfo[87] = errors.New("解密字段不是由改appKey所获取的，请使用正确的appkey调用")
	ApiErrInfo[88] = errors.New("没有调用该接口权限，请到“控制中心”->“应用管理“->“接口管理”进行申请")
	ApiErrInfo[89] = errors.New("非法appkey")

	ApiErrInfo[90] = errors.New("被监控为非法商家")
	ApiErrInfo[91] = errors.New("京东pin被监测为非法")
	ApiErrInfo[92] = errors.New("App在测试状态已超过三个月，需要上线后才能继续调用")
}
