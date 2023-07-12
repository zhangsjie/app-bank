package stru

import (
	"encoding/xml"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"time"
)

/**
报文格式说明:
报文＝报文头＋报文体（packet＝head＋body）
报文头＝（报文字段）n 	（head＝（field）n）
报文体＝（报文字段|报文循环字段）*  （body＝（field|lists）*）*）
报文循环字段＝（报文字段）*	（lists＝（list）*） (list=(field)*)

注：n：表示有n个字段 *：表示有一个或者多个字段

*/

// SPDHeadData
//  @Description: 浦发银行 报文头 数据
//
type SPDHeadData struct {
	XMLName   xml.Name `xml:"head"`
	TransCode string   `xml:"transCode"`
	SignFlag  string   `xml:"signFlag"`
	MasterID  string   `xml:"masterID"`
	YqzlNo    string   `xml:"yqzlNo"`
	PacketID  string   `xml:"packetID"`
	TimeStamp string   `xml:"timeStamp"`
}

// NewSPDHeadData
//  @Description: 新建请求头
//  @param transCode
//  @param signFlag
//  @param masterID
//  @param yqzlNo
//  @param packetID
//  @return *SPDHeadData
//
func NewSPDHeadData(transCode, signFlag, masterID, yqzlNo, packetID string) *SPDHeadData {
	return &SPDHeadData{
		TransCode: transCode,
		SignFlag:  signFlag,
		MasterID:  masterID,
		YqzlNo:    yqzlNo,
		PacketID:  packetID,
		TimeStamp: util.FormatDateTime(time.Now()),
	}
}

// SPDHeadResponseData
//  @Description: 请求头返回的数据
//
type SPDHeadResponseData struct {
	SPDHeadData
	ReturnCode string `xml:"returnCode"`
	ReturnMsg  string `xml:"returnMsg"`
	Signature  string `xml:"Signature"`
}

// SPDRequestData
//  @Description: 请求数据 包含请求头和请求体, 请求体是any类型
//
type SPDRequestData struct {
	XMLName xml.Name    `xml:"packet"`
	Head    SPDHeadData `xml:"head"`
	Body    interface{} `xml:"body"`
}

type SPDReturnCodeMsg struct {
	ReturnCode string `xml:"returnCode"`
	ReturnMsg  string `xml:"returnMsg"`
}

// SPDResponseSignData
//  @Description: 请求返回数据, 浦发的返回体是一个签名, 需要再验签, 解析成具体返回对象
//
type SPDResponseSignData struct {
	XMLName xml.Name            `xml:"packet"`
	Head    SPDHeadResponseData `xml:"head"`
	Body    SPDRequestSignBody  `xml:"body"`
}

// SPDRequestSignBody
//  @Description: 请求体签名
//
type SPDRequestSignBody struct {
	Sign string `xml:"signature"`
	SPDReturnCodeMsg
}

// SPDBankTransferRequest
//  @Description: 转账请求结构体对象
//
type SPDBankTransferRequest struct {
	AuthMasterID        string  `xml:"authMasterID"`
	ElecChequeNo        string  `xml:"elecChequeNo"`
	AcctNo              string  `xml:"acctNo"`
	AcctName            string  `xml:"acctName"`
	PayeeAcctNo         string  `xml:"payeeAcctNo"`
	PayeeName           string  `xml:"payeeName"`
	PayeeBankName       string  `xml:"payeeBankName"`
	Amount              float64 `xml:"amount"`
	SysFlag             string  `xml:"sysFlag"`
	RemitLocation       string  `xml:"remitLocation"`
	Note                string  `xml:"note"`
	PayeeBankSelectFlag string  `xml:"payeeBankSelectFlag"`
	PayeeBankNo         string  `xml:"payeeBankNo"`
}

// SPDBankTransferRequestBody
//  @Description: 转账请求体
//
type SPDBankTransferRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferRequest
}

// SPDBankTransferResponseData
//  @Description: 转账返回结构体对象
//
type SPDBankTransferResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferResponseBody
}

// SPDBankTransferResponseBody
//  @Description: 转账返回体
//
type SPDBankTransferResponseBody struct {
	TransStatus string `xml:"transStatus"`
	AcceptNo    string `xml:"acceptNo"`
	SPDReturnCodeMsg
}

// SPDBankTransferDetailRequestBody
//  @Description: 交易明细请求body
//
type SPDBankTransferDetailRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferDetailRequest
}

// SPDBankTransferDetailRequest
//  @Description: 交易明细请求结构体
//
type SPDBankTransferDetailRequest struct {
	AcctNo        string `xml:"acctNo"`
	DateBeginDate string `xml:"dateBeginDate"`
	DateEndDate   string `xml:"dateEndDate"`
	QueryNumber   int    `xml:"queryNumber"`
	BeginNumber   int    `xml:"beginNumber"`
}

// SPDBankTransferDetailResponseData
//  @Description: 交易明细返回数据
//
type SPDBankTransferDetailResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferDetailResponseBody
	SPDReturnCodeMsg
}

// SPDBankTransferDetailResponseBody
//  @Description: 交易明细返回主体
//
type SPDBankTransferDetailResponseBody struct {
	TotalNumber int64                               `xml:"totalNumber"`
	AcctNo      string                              `xml:"acctNo"`
	Currency    string                              `xml:"currency"`
	MasterName  string                              `xml:"masterName"`
	Reserve1    string                              `xml:"reserve1"`
	Reserve2    string                              `xml:"reserve2"`
	Reserve3    string                              `xml:"reserve3"`
	Reserve4    string                              `xml:"reserve4"`
	Reserve5    string                              `xml:"reserve5"`
	Items       []SPDBankTransferDetailResponseItem `xml:"lists>list"`
}

// SPDBankTransferDetailResponseItem
//  @Description: 交易明细单个Item
//
type SPDBankTransferDetailResponseItem struct {
	TransTime            string  `xml:"transTime"`
	OppositeBankName     string  `xml:"oppositeBankName"`
	SubAccount           string  `xml:"subAccount"`
	TellerJnlNo          string  `xml:"tellerJnlNo"`
	OppositeBankNo       string  `xml:"oppositeBankNo"`
	SubAcctName          string  `xml:"subAcctName"`
	Remark               string  `xml:"remark"`
	BusinessCode         string  `xml:"businessCode"`
	SystemDate           string  `xml:"systemDate"`
	OrgNo                string  `xml:"orgNo"`
	DebitFlag            string  `xml:"debitFlag"`
	TransAmount          float64 `xml:"transAmount"`
	AcctBalance          float64 `xml:"acctBalance"`
	SummaryCode          string  `xml:"summaryCode"`
	Remark1              string  `xml:"remark1"`
	VoucherNo            string  `xml:"voucherNo"`
	CustAcctNo           string  `xml:"custAcctNo"`
	MasterAcctNoType     string  `xml:"masterAcctNoType"`
	TransGy              string  `xml:"transGy"`
	AuthGy               string  `xml:"authGy"`
	ReserveDomain        string  `xml:"reserveDomain"`
	SummonsNumber        string  `xml:"summonsNumber"`
	FillUpMark           string  `xml:"fillUpMark"`
	TransDate            string  `xml:"transDate"`
	UniqueIdentification string  `xml:"uniqueIdentification"`
	OldChannel           string  `xml:"oldChannel"`
	OldTransCode         string  `xml:"oldTransCode"`
	OldTransDate         string  `xml:"oldTransDate"`
	OldStartDate         string  `xml:"oldStartDate"`
	OldEntrustSeqNo      string  `xml:"oldEntrustSeqNo"`
	OldElecChequeNo      string  `xml:"oldElecChequeNo"`
	IsFlag               string  `xml:"isFlag"`
	CustomerNote         string  `xml:"customerNote"`
	BackReason           string  `xml:"backReason"`
	Note1                string  `xml:"note1"`
	Note2                string  `xml:"note2"`
	Note3                string  `xml:"note3"`
}

// SPDBankTransferDetailElectronicReceiptRequestBody
//  @Description: 交易明细电子回单申请-请求体
//
type SPDBankTransferDetailElectronicReceiptRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferDetailElectronicReceiptRequest
}

// SPDBankTransferDetailElectronicReceiptRequest
//  @Description: 交易明细电子回单申请-结构体
//
type SPDBankTransferDetailElectronicReceiptRequest struct {
	BillDownloadChanel string `xml:"billDownloadChanel"`
	AcctNo             string `xml:"acctNo"`
	SingleOrBatchFlag  string `xml:"singleOrBatchFlag"`
	BackhostGyno       string `xml:"backhostGyno"`
	SubpoenaSeqNo      string `xml:"subpoenaSeqNo"`
	BeginDate          string `xml:"beginDate"`
	EndDate            string `xml:"endDate"`
	BeginNumber        int    `xml:"beginNumber"`
	QueryNumber        int    `xml:"queryNumber"`
}

// SPDBankTransferDetailElectronicReceiptResponseData
//  @Description: 交易明细电子回单申请-返回数据
//
type SPDBankTransferDetailElectronicReceiptResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferDetailElectronicReceiptResponseBody
	SPDReturnCodeMsg
}

type SPDBankTransferDetailElectronicReceiptResponseBody struct {
	TotalNumber int64                                                `xml:"totalNumber"`
	Reserve1    int64                                                `xml:"reserve1"`
	Reserve2    int64                                                `xml:"reserve2"`
	Reserve3    int64                                                `xml:"reserve3"`
	Items       []SPDBankTransferDetailElectronicReceiptResponseItem `xml:"lists>list"`
}

type SPDBankTransferDetailElectronicReceiptResponseItem struct {
	AcceptNo       string  `xml:"acceptNo"`
	Result         string  `xml:"result"`
	ResultMsg      string  `xml:"resultMsg"`
	OppositeAcctNo string  `xml:"oppositeAcctNo"`
	OppositeBankNo string  `xml:"oppositeBankNo"`
	DebitFlag      string  `xml:"debitFlag"`
	BusinessCode   string  `xml:"businessCode"`
	BackhostGyno   string  `xml:"backhostGyno"`
	SubpoenaSeqNo  string  `xml:"subpoenaSeqNo"`
	TransDate      string  `xml:"transDate"`
	TransAmount    float64 `xml:"transAmount"`
	Reserve1       string  `xml:"reserve1"`
	Reserve2       string  `xml:"reserve2"`
	Reserve3       string  `xml:"reserve3"`
}

// SPDBankTransferDetailElectronicReceiptDownloadRequestBody
//  @Description: 交易明细电子回单下载-请求体
//
type SPDBankTransferDetailElectronicReceiptDownloadRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferDetailElectronicReceiptDownloadRequest
}

type SPDBankTransferDetailElectronicReceiptDownloadRequest struct {
	AcctNo           string `xml:"acctNo"`
	FileDownloadFlag string `xml:"fileDownloadFlag"`
	FileDownloadPar  string `xml:"fileDownloadPar"`
}

// todo
// 交易明细电子回单下载-返回data 定义

// SPDBankTransferResultRequestBody
//  @Description: 转账结果查询-请求体
//
type SPDBankTransferResultRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferResultRequest
}

type SPDBankTransferResultRequest struct {
	AcctNo            string `xml:"acctNo"`
	PayeeAcctNo       string `xml:"payeeAcctNo"`
	BeginDate         string `xml:"beginDate"`
	EndDate           string `xml:"endDate"`
	ElecChequeNo      string `xml:"elecChequeNo"`
	AcceptNo          string `xml:"acceptNo"`
	QueryNumber       int    `xml:"queryNumber"`
	BeginNumber       int    `xml:"beginNumber"`
	SingleOrBatchFlag string `xml:"singleOrBatchFlag"`
}

// SPDBankTransferResultResponseData
//  @Description: 转账结果查询-返回数据
//
type SPDBankTransferResultResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankTransferResultResponseBody
	SPDReturnCodeMsg
}

type SPDBankTransferResultResponseBody struct {
	TotalCount int64                               `xml:"totalCount"`
	Items      []SPDBankTransferResultResponseItem `xml:"lists>list"`
}

type SPDBankTransferResultResponseItem struct {
	ElecChequeNo  string  `xml:"elecChequeNo"`
	AcceptNo      string  `xml:"acceptNo"`
	SerialNo      string  `xml:"serialNo"`
	TransDate     string  `xml:"transDate"`
	BespeakDate   string  `xml:"bespeakDate"`
	PromiseDate   string  `xml:"PromiseDate"`
	AcctNo        string  `xml:"acctNo"`
	AcctName      string  `xml:"acctName"`
	PayeeAcctNo   string  `xml:"payeeAcctNo"`
	PayeeName     string  `xml:"payeeName"`
	PayeeType     string  `xml:"payeeType"`
	PayeeBankName string  `xml:"payeeBankName"`
	PayeeAddress  string  `xml:"payeeAddress"`
	Amount        float64 `xml:"amount"`
	SysFlag       string  `xml:"sysFlag"`
	RemitLocation string  `xml:"remitLocation"`
	Note          string  `xml:"note"`
	TransStatus   string  `xml:"transStatus"`
	SeqNo         string  `xml:"seqNo"`
}

// SPDBankQueryAccountBalanceRequestBody
//  @Description: 查询账户余额-请求体
//
type SPDBankQueryAccountBalanceRequestBody struct {
	XMLName xml.Name                                `xml:"body"`
	Items   []SPDBankQueryAccountBalanceRequestItem `xml:"lists>list"`
}

type SPDBankQueryAccountBalanceRequestItem struct {
	AcctNo string `xml:"acctNo"`
}

// SPDBankQueryAccountBalanceResponseData
//  @Description: 查询账户余额-返回数据
//
type SPDBankQueryAccountBalanceResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankQueryAccountBalanceResponseBody
	SPDReturnCodeMsg
}

type SPDBankQueryAccountBalanceResponseBody struct {
	Items []SPDBankQueryAccountBalanceResponseItem `xml:"lists>list"`
}

type SPDBankQueryAccountBalanceResponseItem struct {
	AcctNo           string  `xml:"acctNo"`
	MasterID         string  `xml:"masterID"`
	Balance          float64 `xml:"balance"`
	ReserveBalance   float64 `xml:"reserveBalance"`
	FreezeBalance    float64 `xml:"freezeBalance"`
	CortrolBalance   float64 `xml:"cortrolBalance"`
	CanUseBalance    float64 `xml:"canUseBalance"`
	OverdraftBalance float64 `xml:"overdraftBalance"`
}

// SPDCreateVirtualAccountRequestBody
//  @Description: 创建虚账户-请求体
//
type SPDCreateVirtualAccountRequestBody struct {
	XMLName xml.Name                             `xml:"body"`
	AcctNo  string                               `xml:"acctNo"`
	Items   []SPDCreateVirtualAccountRequestItem `xml:"lists>list"`
}

type SPDCreateVirtualAccountRequestItem struct {
	VirtualAccountName string  `xml:"masterName"`
	VirtualAccountNo   string  `xml:"virtualAcctNo"`
	Rate               float64 `xml:"rate"`
}

// SPDCreateVirtualAccountResponseData
//  @Description: 创建虚账户返回数据
//
type SPDCreateVirtualAccountResponseData struct {
	XMLName        xml.Name `xml:"body"`
	BusinessStatus string   `xml:"businessStatus"`
	JnlSeqNo       string   `xml:"jnlSeqNo"`
	Explain        string   `xml:"explain"`
	SPDCreateVirtualAccountResponseBody
}

type SPDCreateVirtualAccountResponseBody struct {
	Items []SPDCreateVirtualAccountResponseItem `xml:"lists>list"`
}

type SPDCreateVirtualAccountResponseItem struct {
	VirtualAccountName string  `xml:"masterName"`
	VirtualAccountNo   string  `xml:"virtualAcctNo"`
	Rate               float64 `xml:"rate"`
}

// SPDQueryVirtualAccountBalanceRequestBody
//  @Description: 实账户下虚账户的即时账户余额 - 请求体
//
type SPDQueryVirtualAccountBalanceRequestBody struct {
	XMLName     xml.Name `xml:"body"`
	AcctNo      string   `xml:"acctNo"`
	BeginNumber int      `xml:"beginNumber"`
	QueryNumber int      `xml:"queryNumber"`
}

// SPDQueryVirtualAccountBalanceResponseData
//  @Description: 实账户下虚账户的即时账户余额 - 返回数据
//
type SPDQueryVirtualAccountBalanceResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDQueryVirtualAccountBalanceResponseBody
}

type SPDQueryVirtualAccountBalanceResponseBody struct {
	VirAcctTotalBalance float64                                     `xml:"virAcctTotalBalance"`
	TrueAcctBalance     float64                                     `xml:"trueAcctBalance"`
	Items               []SPDQueryVirtualAccountBalanceResponseItem `xml:"lists>list"`
}

type SPDQueryVirtualAccountBalanceResponseItem struct {
	AddSubtractNum  float64 `xml:"addSubtractNum"`
	AccumulateNum   float64 `xml:"accumulateNum"`
	VirtualAcctNo   string  `xml:"virtualAcctNo"`
	AccountBalance  float64 `xml:"accountBalance"`
	VirtualAcctName string  `xml:"virtualAcctName"`
}

// SPDBankVirtualAccountTransferRequestBody
//  @Description: 虚账户转账请求体
//
type SPDBankVirtualAccountTransferRequestBody struct {
	XMLName xml.Name `xml:"body"`
	SPDBankVirtualAccountTransferRequest
}

type SPDBankVirtualAccountTransferRequest struct {
	ElectronNumber   string  `xml:"electronNumber"`
	AcctNo           string  `xml:"acctNo"`
	PayerVirAcctNo   string  `xml:"payerVirAcctNo"`
	PayerName        string  `xml:"payerName"`
	PayeeAcctNo      string  `xml:"payeeAcctNo"`
	PayeeAcctName    string  `xml:"payeeAcctName"`
	PayeeBankName    string  `xml:"payeeBankName"`
	PayeeBankAddress string  `xml:"payeeBankAddress"`
	TransAmount      float64 `xml:"transAmount"`
	OwnItBankFlag    string  `xml:"ownItBankFlag"`
	RemitLocation    string  `xml:"remitLocation"`
	Note             string  `xml:"note"`
}

// SPDBankVirtualAccountTransferResponseData
//  @Description: 虚账户转账返回结构体对象
//
type SPDBankVirtualAccountTransferResponseData struct {
	XMLName xml.Name `xml:"body"`
	SPDBankVirtualAccountTransferResponseBody
}

type SPDBankVirtualAccountTransferResponseBody struct {
	BackhostStatus string `xml:"backhostStatus"`
	JnlSeqNo       string `xml:"jnlSeqNo"`
	SPDReturnCodeMsg
}
