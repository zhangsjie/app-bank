package stru

//通用接口
type IcbcGlobalRequest struct {
	AppID      string `json:"app_id"`      // APP的编号,应用在API开放平台注册时生成
	MsgID      string `json:"msg_id"`      // 消息通讯唯一编号，每次调用独立生成，APP级唯一
	SignType   string `json:"sign_type"`   // 签名类型，CA-工行颁发的证书认证，RSA-RSAWithSha1，RSA2-RSAWithSha256，缺省为RSA
	Sign       string `json:"sign"`        // 报文签名
	Timestamp  string `json:"timestamp"`   // 交易发生时间戳，yyyy-MM-dd HH:mm:ss格式
	BizContent string `json:"biz_content"` // 请求参数的集合
}

//通用接口
type IcbcGlobalResponse struct {
	Sign               string `json:"sign"`                 // 针对返回参数集合的签名
	ResponseBizContent string `json:"response_biz_content"` // 响应参数集合，包含公共和业务参数
}

//签约参数
type SignRequestParams struct {
	AppID      string        `json:"app_id"`               // APP的编号,应用在API开放平台注册时生成
	ApiName    string        `json:"api_name"`             // api接口名称
	ApiVersion string        `json:"api_version"`          // api接口版本
	CorpNo     string        `json:"corpNo"`               // 一级合作方编号
	CoMode     string        `json:"coMode"`               // 合作模式，1：代理记账；2：自主记账
	AccCompNo  string        `json:"accCompNo,omitempty"`  // 二级合作方编号，合作模式为1时，必输 (可选字段)
	Account    string        `json:"account"`              // 主账户
	CurrType   string        `json:"currType"`             // 主账户币种，1：人民币
	AccFlag    string        `json:"accFlag"`              // 本他行标志，1：本行；2：他行
	CntioFlag  string        `json:"cntioFlag"`            // 境内外标志，1：境内；2：境外
	Phone      string        `json:"phone"`                // 手机号
	EpType     string        `json:"epType"`               // 是否自动展期，0：否；1：是
	EpTimes    string        `json:"epTimes,omitempty"`    // 展期期数，epType为1时，必输 (可选字段)
	PayAccNo   string        `json:"payAccNo,omitempty"`   // 收费账号 (可选字段)
	PayCurr    string        `json:"payCurr,omitempty"`    // 收费账号币种 (可选字段)
	PayAccName string        `json:"payAccName,omitempty"` // 收费账号户名，默认为空 (可选字段)
	PayLimit   string        `json:"payLimit,omitempty"`   // 收费周期，默认为空 (可选字段)
	PayBegDate string        `json:"payBegDate,omitempty"` // 收费开始日期，默认为空 (可选字段)
	PayEndDate string        `json:"payEndDate,omitempty"` // 收费结束日期，默认为空 (可选字段)
	Remark     string        `json:"remark,omitempty"`     // 备注，默认为空 (可选字段)
	AccList    []AccListItem `json:"accList"`              // 下挂账号列表
}

type AccListItem struct {
	Account       string `json:"account"`       // 账号
	CurrType      string `json:"currType"`      // 币种
	AccFlag       string `json:"accFlag"`       // 本他行标志，1：本行；2：他行
	CntioFlag     string `json:"cntioFlag"`     // 境内外标志，1：境内；2：境外
	IsMainAcc     string `json:"isMainAcc"`     // 主帐户标志，0：否；1：是
	ReceiptFlag   string `json:"receiptFlag"`   // 开通电子单据标示，0：关；1：开
	StatementFlag string `json:"statementFlag"` // 开通账务明细查询标示，0：关；1：开
}
