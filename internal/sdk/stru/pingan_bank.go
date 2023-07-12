package stru

// PingAnBankTransferRequest
//
//	@Description: 转账请求结构体对象
//
// @MrchCode 企业银企直联标准代码
// @tradeDate 交易日期(yyyyMMdd)
// @tradeTime 交易时间 hhmmss
// @CnsmrSeqNo 请求方系统流水号 唯一标识一笔交易 如果某种交易要有多次请求的才能完成的，多个交易请求包流水号要保持一致
// @ThirdVoucher 转账凭证号 标示交易唯一性，同一客户上送的不可重复，建议格式：yyyymmddHHSS+8位系列。最少10位长度，要求6个月内唯一
// @CstInnerFlowNo 客户自定义凭证号 用于客户转账登记和内部识别，通过转账结果查询可以返回。银行不检查唯一性(不必填)
// @CcyCode 货币类型 RMB-人民币
// @OutAcctNo 付款人账户
// @OutAcctName 付款人姓名 (不必填)
// @OutAcctAddr 付款人地址 (不必填) 议填写付款账户的分行、网点名称议填写付款账户的分行、网点名称
// @InAcctBankNode 收款人开户行行号 跨行转账建议必输。为人行登记在册的商业银行号，若输入则长度必须在4~12位之间；
// @InAcctNo 收款人账户
// @InAcctName 收款人名；
// @InAcctBankName 收款人开户行名称,建议格式 XX银行
// @InAcctProvinceCode 收款账户银行开户省代码或省名称 (非必填) 建议跨行转账输入；对照码参考“附录-省对照表”；也可输入“附录-省对照表”中的省名称。
// @InAcctCityName 收款账户银行开户市 (非必填)
// @TranAmount 转出金额
// @UseEx 资金用途 (非必填)
// @UnionFlag 行内跨行标志 1：行内转账，0：跨行转账
// @AddrFlag “1”—同城   “2”—异地；若无法区分，可默认送1-同城。
// @MainAcctNo “1”—同城   “2”—异地；若无法区分，可默认送1-同城。
type PingAnBankTransferRequest struct {
	MrchCode           string `json:"MrchCode"`
	TradeDate          int    `json:"tradeDate"`
	TradeTime          int    `json:"tradeTime"`
	CnsmrSeqNo         string `json:"CnsmrSeqNo"`
	ThirdVoucher       string `json:"ThirdVoucher"`
	CstInnerFlowNo     string `json:"CstInnerFlowNo"`
	CcyCode            string `json:"CcyCode"`
	OutAcctNo          string `json:"OutAcctNo"`
	OutAcctName        string `json:"OutAcctName"`
	OutAcctAddr        string `json:"OutAcctAddr"`
	InAcctBankNode     string `json:"InAcctBankNode"`
	InAcctNo           string `json:"InAcctNo"`
	InAcctName         string `json:"InAcctName"`
	InAcctBankName     string `json:"InAcctBankName"`
	InAcctProvinceCode string `json:"InAcctProvinceCode"`
	InAcctCityName     string `json:"InAcctCityName"`
	TranAmount         string `json:"TranAmount"`
	UseEx              string `json:"UseEx"`
	UnionFlag          string `json:"UnionFlag"`
	AddrFlag           string `json:"AddrFlag"`
	MainAcctNo         string `json:"MainAcctNo"`
	PaymentModeType    string `json:"PaymentModeType"`
}

// PingAnBankTransferResponse
//
//	@Description: 转账返回结构体对象
type PingAnBankTransferResponse struct {
	PinganErrorResult
	ThirdVoucher   string `json:"ThirdVoucher"`
	FrontLogNo     string `json:"FrontLogNo"`
	CstInnerFlowNo string `json:"CstInnerFlowNo"`
	CcyCode        string `json:"CcyCode"`
	OutAcctName    string `json:"OutAcctName"`
	OutAcctNo      string `json:"OutAcctNo"`
	InAcctBankName string `json:"InAcctBankName"`
	InAcctNo       string `json:"InAcctNo"`
	InAcctName     string `json:"InAcctName"`
	TranAmount     string `json:"TranAmount"`
	UnionFlag      string `json:"UnionFlag"`
	Fee1           string `json:"Fee1"`
	Fee2           string `json:"Fee2"`
	HostFlowNo     string `json:"hostFlowNo"`
	HostTxDate     string `json:"hostTxDate"`
	Stt            string `json:"stt"`
}

type PingAnBankRequestBody struct {
	RequestBody   string `json:"requestBody"`
	InterfaceName string `json:"interfaceName"`
	InterfaceType int    `json:"interfaceType"`
	ZuId          string `json:"zuId"`
}
type PingAnBankResponseBody struct {
	Data string `json:"data"`
	Code int    `json:"code"`
}

//查询企业账户余额
// @Description: 请求结构体对象
// @MrchCode 企业银企直联标准代码
// @tradeDate 交易日期(yyyyMMdd)
// @tradeTime 交易时间 hhmmss
// @CnsmrSeqNo 请求方系统流水号 唯一标识一笔交易  如果某种交易要有多次请求的才能完成的，多个交易请求包流水号要保持一致(自己生成)
// @CcyType 钞汇标志 C 钞户, R汇户,默认为C
// @CcyCode 货币类型 RMB 人民币,USD 美元，HKD 港币, 默认为RMB

type PinganBankCorAcctBalanceQueryRequest struct {
	MrchCode   string `json:"MrchCode"`
	TradeDate  int    `json:"tradeDate"`
	TradeTime  int    `json:"tradeTime"`
	CnsmrSeqNo string `json:"CnsmrSeqNo"`
	Account    string `json:"Account"`
	CcyType    string `json:"CcyType"`
	CcyCode    string `json:"CcyCode"`
}
type PinganBankCorAcctBalanceQueryResponse struct {
	PinganErrorResult
	Account                  string `json:"Account"`                  //货币类型
	CcyCode                  string `json:"CcyCode"`                  //钞汇标志
	CcyType                  string `json:"CcyType"`                  //账户户名
	AccountName              string `json:"AccountName"`              //可用余额
	Balance                  string `json:"Balance"`                  //账面余额
	AccountStatus            string `json:"AccountStatus"`            //账户状态  若有多个状态， “|”分割，如：A| DGZH01	A：正常；DGZH01: 法律冻结 DGZH02: 账户止付
	HoldBalance              string `json:"HoldBalance"`              //冻结金额
	StopBalance              string `json:"StopBalance"`              //止付金额
	LastBalance              string `json:"LastBalance"`              //昨日余额
	HRate1                   string `json:"HRate1"`                   //活期存款计息执行利率
	XDRate2                  string `json:"XDRate2"`                  //协定执行利率
	AcctBalance              string `json:"AcctBalance"`              //账户余额
	AgreeDepReserveBalance   string `json:"AgreeDepReserveBalance"`   //协定额度
	BeginEffectDate          string `json:"BeginEffectDate"`          //开始生效日期
	ExpiryDate               string `json:"ExpiryDate"`               //到期日
	AgreeDepositInterest     string `json:"AgreeDepositInterest"`     //利息金额
	AgreeDepositInterestType string `json:"AgreeDepositInterestType"` //利息种类
	EnableFlagEndTime        string `json:"EnableFlagEndTime"`        //触发标识确定协议的结束时间
}

type PinganHistoryTransactionDetailsRequest struct {
	MrchCode       string `json:"MrchCode"`       //Y企业银企直联标准代码  银行提供给企业的20位唯一的标识代码
	RecvLength     int    `json:"RecvLength"`     //N接收报文长度  报文数据长度；不包括附件内容、签名内容的长度，不够左补0
	TradeDate      int    `json:"tradeDate"`      //N交易日期(yyyyMMdd)  yyyymmdd
	TradeTime      int    `json:"tradeTime"`      //N交易时间  hhmmss
	CnsmrSeqNo     string `json:"CnsmrSeqNo"`     //Y请求方系统流水号  "唯一标识一笔交易 备注：（如果某种交易要有多次请求的才能完成的，多个交易请求包流水号要保持一致）"
	AcctNo         string `json:"AcctNo"`         //Y账号
	CcyCode        string `json:"CcyCode"`        //Y币种 RMB
	BeginDate      string `json:"BeginDate"`      //Y开始日期  "若查询当日明细，开始、结束日期必须为当天；若查询历史明细，开始、结束日期必须是历史日期 格式yyyyMMdd"
	EndDate        string `json:"EndDate"`        //Y结束日期  格式yyyyMMdd
	PageNo         string `json:"PageNo"`         //Y查询页码  1：第一页，依次递增
	PageSize       string `json:"PageSize"`       //N每页明细数量  "当日明细默认每页30条记录，支持最大每页100条，若上送PageSize>100无效，等同100 历史明细默认每页30条记录，支持最大每页1000条，若上送PageSize>1000则提示输入错误且每次查询必须固定为此值，否则出现明细遗漏"
	Reserve        string `json:"Reserve"`        //N预留字段
	OrderMode      string `json:"OrderMode"`      //N记录排序标志
	BankTranFlowNo string `json:"BankTranFlowNo"` //N银行交易流水号
	OppAcctNo      string `json:"OppAcctNo"`      //N交易对手账
}
type PinganHistoryTransactionDetailsResponse struct {
	PinganErrorResult
	AcctNo       string                                `json:"AcctNo"`       //账号
	CcyCode      string                                `json:"CcyCode"`      //货币类型
	EndFlag      string                                `json:"EndFlag"`      //数据结束标志
	Reserve      string                                `json:"Reserve"`      //预留字段
	PageNo       string                                `json:"PageNo"`       //查询页码
	PageRecCount string                                `json:"PageRecCount"` //记录笔数
	List         []PinganHistoryTransactionDetailsItem `json:"list"`
}

type PinganHistoryTransactionDetailsItem struct {
	AcctDate         string `json:"AcctDate"`         //主机记账日期
	TxTime           string `json:"TxTime"`           //交易时间
	HostTrace        string `json:"HostTrace"`        //主机流水号
	BussSeqNo        string `json:"BussSeqNo"`        //业务流水号
	DetailSerialNo   string `json:"DetailSerialNo"`   //明细序号
	OutNode          string `json:"OutNode"`          //付款方网点号
	OutBankNo        string `json:"OutBankNo"`        //付款方联行号
	OutBankName      string `json:"OutBankName"`      //付款行名称
	OutAcctNo        string `json:"OutAcctNo"`        //付款方账号
	OutAcctName      string `json:"OutAcctName"`      //付款方户名
	CcyCode          string `json:"CcyCode"`          //结算币种
	TranAmount       string `json:"TranAmount"`       //交易金额
	InNode           string `json:"InNode"`           //收款方网点号
	InBankNo         string `json:"InBankNo"`         //收款方联行号
	InBankName       string `json:"InBankName"`       //收款方行名
	InAcctNo         string `json:"InAcctNo"`         //收款方账号
	InAcctName       string `json:"InAcctName"`       //收款方户名
	DcFlag           string `json:"DcFlag"`           //借贷标志 D借，出账；C贷，入账
	AbstractStr      string `json:"AbstractStr"`      //摘要，未翻译的摘要，如TRS
	VoucherNo        string `json:"VoucherNo"`        //凭证号
	TranFee          string `json:"TranFee"`          //手续费
	PostFee          string `json:"PostFee"`          //邮电费
	AcctBalance      string `json:"AcctBalance"`      //账面余额
	Purpose          string `json:"Purpose"`          //用途，附言
	AbstractStrDesc  string `json:"AbstractStr_Desc"` //中文摘要，AbstractStr的中文翻译
	ProxyPayName     string `json:"ProxyPayName"`     //代理人户名
	ProxyPayAcc      string `json:"ProxyPayAcc"`      //代理人账号
	ProxyPayBankName string `json:"ProxyPayBankName"` //代理人银行名称
	HostDate         string `json:"HostDate"`         //主机日期
	Remark1          string `json:"Remark1"`          //备注1
	Remark2          string `json:"Remark2"`          //备注2
	BeReverseFlag    string `json:"BeReverseFlag"`    //被冲正标志
	SeqTime          string `json:"SeqTime"`          //时序时间
	FeeCode          string `json:"FeeCode"`          //费用代码
}

// 该功能单笔转账指令查询_银企直联
// 交易说明
// 交易状态Stt域说明：
// 1、交易状态核实必须用stt，而不能使用Yhcljg域；
// 2、Stt=99是异常，30是失败，可以显示<Yhcljg>作为对应的银行状态描述，而不能作为成功、失败的判断依据了。
// 3、Stt=99的建议次日10点再同步一次（银行每日5点左右对账后会更新状态为20、30），若仍然为99状态，可以连续同步n日，每日同步一次即可。也可以设置报警，若第二日仍然是99就联系银行核实状态。
type PinganSignleTransferQueryRequest struct {
	MrchCode         string `json:"MrchCode"`         //Y 银行提供给企业的20位唯一的标识代码
	RecvLength       int    `json:"RecvLength"`       //N 报文数据长度；不包括附件内容、签名内容的长度，不够左补0
	TradeDate        int    `json:"tradeDate"`        //N yyyymmdd
	TradeTime        int    `json:"tradeTime"`        //N hhmmss
	CnsmrSeqNo       string `json:"CnsmrSeqNo"`       //Y "唯一标识一笔交 备注：（如果某种交易要有多次请求的才能完成的，多个交易请求包流水号要保持一致）"
	OrigThirdVoucher string `json:"OrigThirdVoucher"` //N 推荐使用；使用400409接口上送的ThirdVoucher
	OrigFrontLogNo   string `json:"OrigFrontLogNo"`   //N "不推荐使用；银行返回的转账流水号
}

type PinganSignleTransferQueryResponse struct {
	PinganErrorResult
	OrigThirdVoucher string `json:"OrigThirdVoucher"` //转账凭证号
	FrontLogNo       string `json:"FrontLogNo"`       //银行流水号
	CstInnerFlowNo   string `json:"CstInnerFlowNo"`   //客户自定义凭证号
	CcyCode          string `json:"CcyCode"`          //货币类型
	OutAcctBankName  string `json:"OutAcctBankName"`  //转出账户开户网点名
	OutAcctNo        string `json:"OutAcctNo"`        //转出账户
	InAcctBankName   string `json:"InAcctBankName"`   //转入账户网点名称
	InAcctNo         string `json:"InAcctNo"`         //转入账户
	InAcctName       string `json:"InAcctName"`       //转入账户户名
	TranAmount       string `json:"TranAmount"`       //交易金额
	UnionFlag        string `json:"UnionFlag"`        //行内跨行标志 1：行内转账，0：跨行转账
	Stt              string `json:"Stt"`              //交易状态标志 20:成功 30失败:其他为银行受理成功处理中
	IsBack           string `json:"IsBack"`           //转账退票标志 0:未退票 1:退票,默认0
	BackRem          string `json:"BackRem"`          //支付失败或退票原因描述 如果是超级网银则返回如下信息:	RJ01对方返回：账号不存在	RJ02对方返回：账号、户名不符	大小额支付则返回失败描述
	Yhcljg           string `json:"Yhcljg"`           //银行处理结果 格式为：“六位代码:中文描述”。冒号为半角。如：000000：转账成功
	//处理中的返回(以如下返回开头)：
	//MA9111:交易正在受理中
	//000000:交易受理成功待处理
	//000000:交易处理中
	//000000:交易受理成功处理中
	//成功的返回：
	//000000:转账交易成功
	//其他的返回都为失败:
	//MA9112:转账失败
	SysFlag          string `json:"SysFlag"`          //转账加急标志 Y：加急 N：普通S：特急
	Fee              string `json:"Fee"`              //转账手续费
	TransBsn         string `json:"TransBsn"`         //转账代码类型 004：单笔转账；	4014：单笔批量；	4034：汇总批量
	SubmitTime       string `json:"submitTime"`       //交易受理时间
	AccountDate      string `json:"AccountDate"`      //记账日期
	HostFlowNo       string `json:"hostFlowNo"`       //主机记账流水号
	HostErrorCode    string `json:"hostErrorCode"`    //错误码 交易失败的错误代码
	ProxyPayName     string `json:"ProxyPayName"`     //代理人户名
	ProxyPayAcc      string `json:"ProxyPayAcc"`      //代理人账号
	ProxyPayBankName string `json:"ProxyPayBankName"` //代理人银行名
}

/*
虚子账户操作编辑接口(新建删除修改恢复)
//银行实时处理清分台账编码开户请求并反馈结果。清分台账编码生成的规则是，总长度14位，组成规则有两种：302+五位公司码+六位序列号、60/61/62+2位公司码+10位顺序号。
//上送清分台账编码的6位或10位的序号，/返回完整的清分台账编码。五位、十位公司码查询参考“智能清分台账编码关系查询”接口。清分台账编码修改、删除、维护时输入完整的清分台账编码；
//清分台账编码新建、修改的同时会复制父账户的权限，//若后期父账户的权限申请变更，清分台账编码的权限不会自动同步，//需要调用B-同步清分台账编码权限功能，将清分台账编码的权限同步为最新的。
*/
type PinganCreateVirtualAccountRequest struct {
	MrchCode            string `json:"MrchCode"`            //Y 企业银企直联标准代码
	RecvLength          int    `json:"RecvLength"`          //N 接收报文长度
	tradeDate           int    `json:"tradeDate"`           //N 交易日期(yyyyMMdd)
	tradeTime           int    `json:"tradeTime"`           //N 交易时间
	CnsmrSeqNo          string `json:"CnsmrSeqNo"`          //Y 请求方系统流水号 唯一标识一笔交易
	MainAccount         string `json:"MainAccount"`         //Y 智能账号 签约账户
	CcyCode             string `json:"CcyCode"`             //N 币种
	SubAccountSeq       string `json:"SubAccountSeq"`       //N 清分台账编码序号  六位序列号、十位序号  当功能码为A-新建：必输
	SubAccount          string `json:"SubAccount"`          //N 清分台账编码  当功能码为A-新建：不允许输入
	SubAccountName      string `json:"SubAccountName"`      //Y 清分台账编码别名 长度限制：200字符
	SubAccountNameEn    string `json:"SubAccountNameEn"`    //N 清分台账编码英文别名 长度限制：200字符
	OpFlag              string `json:"OpFlag"`              //Y 功能码 A-新建 U-修改 D-删除 R-恢复
	ODFlag              string `json:"ODFlag"`              //N 清分台账编码透支标志  默认为N非透支  Y透支
	InterestFlag        string `json:"InterestFlag"`        //N 内部计息标志  默认值：N 否,Y是
	SettleInterestCycle string `json:"SettleInterestCycle"` //N 内部计息周期  D-按日	M-按月	Q-按季	Y-按年	默认值：Q，不可更改
	Rate                string `json:"Rate"`                //N 计息利率  默认值：0
	ZSZFStatus          string `json:"ZSZFStatus"`          //N 止收止付状态   1-止收，2-止付，3-止收止付，0-取消止收止付
}
type PinganCreateVirtualAccountResponse struct {
	PinganErrorResult
	MainAccount         string `json:"MainAccount"`         //智能账号
	SubAccountNo        string `json:"SubAccountNo"`        //清分台账编码  302+五位公司码+六位序列号
	SubAccountName      string `json:"SubAccountName"`      //清分台账编码别名
	Stt                 string `json:"Stt"`                 //清分台账编码状态
	LastModifyDate      string `json:"LastModifyDate"`      //最后维护日期
	InterestFlag        string `json:"InterestFlag"`        //内部计息标志
	SettleInterestCycle string `json:"SettleInterestCycle"` //内部计息周期
	Rate                string `json:"Rate"`                //计息利率
}

// 查询清分台账编码信息，包括清分台账编码的余额 。
type PinganQueryVirtualAccountBalanceRequest struct {
	MrchCode        string `json:"MrchCode"`        //企业银企直联标准代码
	CnsmrSeqNo      string `json:"CnsmrSeqNo"`      //  唯一标识一笔交易 备注：（如果某种交易要有多次请求的才能完成的，多个交易请求包流水号要保持一致）
	MainAccount     string `json:"MainAccount"`     // MainAccount
	ReqSubAccountNo string `json:"ReqSubAccountNo"` // ReqSubAccountNo
}
type PinganQueryVirtualAccountBalanceResponse struct {
	PinganErrorResult
	SubAccountNo   string `json:"SubAccountNo"`   //清分台账编码
	MainAccount    string `json:"MainAccount"`    //智能账号
	CcyCode        string `json:"CcyCode"`        //币种
	SubAccountName string `json:"SubAccountName"` //清分台账编码别名
	SubAccBalance  string `json:"SubAccBalance"`  //清分台账编码余额
	ZSZFStatus     string `json:"ZSZFStatus"`     //止收止付状态  1-止收，2-止付，3-止收止付
	Stt            string `json:"Stt"`            //清分台账编码状态 A--正常	C--销户
	LastModifyDate string `json:"LastModifyDate"` //最后维护日期
}

// 查询账户历史日期的电子回单数据 支持最近一年日期范围内的查询，单次查询的起始日期、结束日期需要在7天内
type PinganSingleDataQueryRequest struct {
	MrchCode      string  `json:"MrchCode"`      //Y 企业银企直联标准代码
	RecvLength    int     `json:"RecvLength"`    //N 接收报文长度
	TradeDate     int     `json:"tradeDate"`     //N 交易日期(yyyyMMdd)
	TradeTime     int     `json:"tradeTime"`     //N 交易时间
	CnsmrSeqNo    string  `json:"CnsmrSeqNo"`    //Y 请求方系统流水号
	AcctNo        string  `json:"AcctNo"`        //Y 账号
	ReceiptType   string  `json:"ReceiptType"`   //Y 回单类型  参照回单类型	注：查全部可送“ALL”
	SubType       string  `json:"SubType"`       //Y 子类型 参照回单类型	注：查全部可送“ALL”
	StartDate     string  `json:"StartDate"`     //Y 起始日期 格式yyyyMMdd（记账日期
	EndDate       string  `json:"EndDate"`       //Y 结束日期 格式yyyyMMdd（记账日期
	StartRecord   int     `json:"StartRecord"`   //Y 起始记录数  起始值为1，不能送0
	RecordNum     int     `json:"RecordNum"`     //Y 本批记录数
	StartAmt      float64 `json:"StartAmt"`      //N 开始金额
	EntAmt        float64 `json:"EntAmt"`        //N 结束金额
	OrderMode     string  `json:"OrderMode"`     //N 排序方式 01：交易时间从近到远	002：交易时间从远到近	003：金额升序（从小到大）	004：金额降序（从大到小）	005：回单号升序	006：回单号降序
	PayeeAcctNo   string  `json:"PayeeAcctNo"`   //N 收款人账号
	PayeeName     string  `json:"PayeeName"`     //N 收款人名称
	DrCrFlag      string  `json:"DrCrFlag"`      //N 借贷标志  D：借方交易	C：贷方交易
	Ccy           string  `json:"Ccy"`           //N 币种
	SerialNo      string  `json:"SerialNo"`      //N 顺序号
	PrintBranchId string  `json:"PrintBranchId"` //N 打印网点
	ReceiptNo     string  `json:"ReceiptNo"`     //N 回单号
	PrintFlag     string  `json:"PrintFlag"`     //N 打印标志 0：首次打印	1：补打
}
type PinganSingleDataQueryResponse struct {
	PinganErrorResult
	RecordTotalCount string                              `json:"RecordTotalCount"` //记录总数
	StartRecord      string                              `json:"StartRecord"`      //起始记录数
	ResultNum        string                              `json:"ResultNum"`        //本次返回记录数
	EndFlag          string                              `json:"EndFlag"`          //结束标志 Y:无剩余记录	N:有剩余记录
	List             []PinganSingleDataQueryResponseItem `json:"lists>list"`
}
type PinganSingleDataQueryResponseItem struct {
	ReceiptNo               string  `json:"ReceiptNo"`               //回单号
	CheckCode               string  `json:"CheckCode"`               //验证码
	ReceiptType             string  `json:"ReceiptType"`             //回单类型
	SubType                 string  `json:"SubType"`                 //回单子类
	BookingDate             string  `json:"BookingDate"`             //记账日期
	PayerName               string  `json:"PayerName"`               //付款人名称
	PayeeName               string  `json:"PayeeName"`               //收款人名称
	PayerAccNo              string  `json:"PayerAccNo"`              //付款人账号
	PayeeAccNo              string  `json:"PayeeAccNo"`              //收款人账号
	PayerAcctOpenBranchID   string  `json:"PayerAcctOpenBranchID"`   //付款人开户行
	PayeeAcctOpenBranchName string  `json:"PayeeAcctOpenBranchName"` //收款人开户行名称
	MainAcctNo              string  `json:"MainAcctNo"`              //主账号
	SubAcctNo               string  `json:"SubAcctNo"`               //子账号
	OldAcctNo               string  `json:"OldAcctNo"`               //原账号
	Ccy                     string  `json:"Ccy"`                     //币种
	TranAmt                 float64 `json:"TranAmt"`                 //交易金额
	SubBranchID             string  `json:"SubBranchID"`             //网点号
	DrCrFlag                string  `json:"DrCrFlag"`                //借贷标志
	Crpp                    string  `json:"Crpp"`                    //资金用途
	Corpus                  float64 `json:"Corpus"`                  //本金
	DepositIntRate          float64 `json:"DepositIntRate"`          //存款利率
	DepositReceiptNo        string  `json:"DepositReceiptNo"`        //存单号
	StartPeriod             string  `json:"StartPeriod"`             //起始期
	EndPeriod               string  `json:"EndPeriod"`               //结束期
	InterestTax             float64 `json:"InterestTax"`             //利息税
	IntInterest             float64 `json:"IntInterest"`             //利息
	OverdraftInterest       float64 `json:"OverdraftInterest"`       //透支利息
	TaxRate                 float64 `json:"TaxRate"`                 //税率
	LoanAcctNo              string  `json:"LoanAcctNo"`              //贷款账号
	DuebillNo               string  `json:"DuebillNo"`               //借据号
	PaidAmt                 float64 `json:"PaidAmt"`                 //还款金额
	RepayCorpus             float64 `json:"RepayCorpus"`             //还款本金
	ReplyInterest           float64 `json:"ReplyInterest"`           //还款利息
	ComInterest             float64 `json:"ComInterest"`             //复利
	CorpusBalance           float64 `json:"CorpusBalance"`           //本金余额
	DueRepayCorpus          float64 `json:"DueRepayCorpus"`          //应还本金
	RepayCount              int     `json:"RepayCount"`              //还款期数
	Commission              float64 `json:"Commission"`              //手续费金额
	MaterialFee             float64 `json:"MaterialFee"`             //工本费
	TaxedInterest           float64 `json:"TaxedInterest"`           //税后利息
	HostSeqNo               string  `json:"HostSeqNo"`               //主机流水号
	LoanIntRate             float64 `json:"LoanIntRate"`             //贷款利率
	ReceivableInterest      float64 `json:"ReceivableInterest"`      //应收利息
	TellerNo                string  `json:"TellerNo"`                //柜员号
	AuthTellerNo            string  `json:"AuthTellerNo"`            //授权柜员号
	PrintClientName         string  `json:"PrintClientName"`         //打印客户端名称
	PrintTime               string  `json:"PrintTime"`               //打印时间
	PrintTimes              int     `json:"PrintTimes"`              //打印次数
	RegionNo                string  `json:"RegionNo"`                //地区号
	TermNo                  string  `json:"TermNo"`                  //终端号
	PrintNote               string  `json:"PrintNote"`               //打印节点
	BussType                string  `json:"BussType"`                //业务类型
	IntSettleAcctNo         string  `json:"IntSettleAcctNo"`         //结息账号
	AcctOpenBranchID        string  `json:"AcctOpenBranchID"`        //账户开户行行号
	TranDate                string  `json:"TranDate"`                //交易日期
	TranTime                string  `json:"TranTime"`                //交易时间
	BranchId                string  `json:"BranchId"`                //机构号
	SerialNo                string  `json:"SerialNo"`                //顺序号
	RecordType              string  `json:"RecordType"`              //记录类型
	FrontEndCode            string  `json:"FrontEndCode"`            //前置机代码
	RemarkCode              string  `json:"RemarkCode"`              //摘要码
	Summary                 string  `json:"Summary"`                 //摘要
}

// 当日历史回单数据查询接口 。
type PinganSameDayHistoryReceiptDataQueryRequest struct {
	MrchCode         string `json:"MrchCode"`                   //Y 企业银企直联标准代码
	CnsmrSeqNo       string `json:"CnsmrSeqNo"`                 //Y 请求方系统流水号
	OutAccNo         string `json:"OutAccNo"`                   //Y 账号
	AccountBeginDate string `json:"AccountBeginDate,omitempty"` //N 记账起始日期 查询当日无需输入此字段 YYYYMMDD
	AccountEndDate   string `json:"AccountEndDate,omitempty"`   //N 记账结束日期
	HostFlow         string `json:"HostFlow"`                   //N 核心流水号 银行核心流水号、银行主机流水号。如取4005返回的HostFlowNo，4013返回的HostFlowNo.
	DcFlag           string `json:"DcFlag,omitempty"`           //N 借贷标志  D:借  C:贷
	SortType         string `json:"SortType,omitempty"`         //N 排序方式 0默认排序
	CCY              string `json:"CCY"`                        //N 币种 默认 RMB
	ReceiptType      string `json:"ReceiptType,omitempty"`      //N 回单类型
	SubReceiptType   string `json:"SubReceiptType,omitempty"`   //N 子回单类型
	RecordStartNo    string `json:"RecordStartNo"`              //N 记录起始号 记录起始号 用于分页	默认：1
	RecordNumber     string `json:"RecordNumber"`               //N 请求记录数 分页条数最大100条   默认：100
}
type PinganErrorResult struct {
	Errors  []PinganErrorResultItem `json:"Errors"`
	Message string                  `json:"Message"`
	Code    string                  `json:"Code"`
}
type PinganErrorResultItem struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

type PinganFileResult struct {
	FilePath string `json:"pngOssPath"`
}

type PinganUserAcctSignatureApplyRequest struct {
	MrchCode     string `json:"MrchCode"`
	CnsmrSeqNo   string `json:"CnsmrSeqNo"`
	ThirdVoucher string `json:"ThirdVoucher"`
	OpFlag       string `json:"OpFlag"` //“A”代表账管+签约；“D”代表账管+解约；“U”代表账管+签约修改。修改场景，只能修改开通权限BsnStr。
	ZuID         string `json:"ZuID"`
	AccountNo    string `json:"AccountNo"`
	AccountName  string `json:"AccountName"`
	BsnStr       string `json:"BsnStr"` //4001,400103,4013,ELC008,ELC002,F00101,ELC007,400409,400505,401805,401501,404710,4048
}
type PinganUserAcctSignatureApplyResponse struct {
	PinganErrorResult
	ThirdVoucher string `json:"ThirdVoucher"`
	ZuID         string `json:"ZuID"`
	OpFlag       string `json:"OpFlag"`
	Stt          string `json:"Stt"`
	AccountNo    string `json:"AccountNo"`
	Remark       string `json:"Remark"`
}
type PinganUserAcctSignatureQueryRequest struct {
	MrchCode   string `json:"MrchCode"`
	CnsmrSeqNo string `json:"CnsmrSeqNo"`
	ZuID       string `json:"ZuID"`
	AccountNo  string `json:"AccountNo"`
}

type PinganSubAcctBalanceAdjustRuest struct {
	//此接口支持相同智能账号下的智能清分台账编码之间的资金划转。此接口的付款最终状态可以使用“支付单状态查询”接口来确认状态。
	MrchCode          string `json:"MrchCode"`
	RecvLength        int    `json:"RecvLength"`
	TradeDate         int    `json:"tradeDate"`
	TradeTime         int    `json:"tradeTime"`
	CnsmrSeqNo        string `json:"CnsmrSeqNo"`
	ThirdVoucher      string `json:"ThirdVoucher"`
	CstInnerFlowNo    string `json:"CstInnerFlowNo"`
	MainAccount       string `json:"MainAccount"`
	MainAccountName   string `json:"MainAccountName"`
	CcyCode           string `json:"CcyCode"`
	OutSubAccount     string `json:"OutSubAccount"`
	OutSubAccountName string `json:"OutSubAccountName"`
	TranAmount        string `json:"TranAmount"`
	InSubAccNo        string `json:"InSubAccNo"`
	InSubAccName      string `json:"InSubAccName"`
	UseEx             string `json:"UseEx"`
}
type PinganSubAcctBalanceAdjustResponse struct {
	PinganErrorResult
	ThirdVoucher      string `json:"ThirdVoucher"`
	FrontFlowNo       string `json:"FrontFlowNo"`
	CstInnerFlowNo    string `json:"CstInnerFlowNo"`
	MainAccount       string `json:"MainAccount"`
	OutSubAccount     string `json:"OutSubAccount"`
	OutSubAccountName string `json:"OutSubAccountName"`
	OutSubAccBalance  string `json:"OutSubAccBalance"`
	CcyCode           string `json:"CcyCode"`
	TranAmount        string `json:"TranAmount"`
	InSubAccNo        string `json:"InSubAccNo"`
	InSubAccName      string `json:"InSubAccName"`
	InSubAccBalance   string `json:"InSubAccBalance"`
}
