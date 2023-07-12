package stru

import (
	"encoding/xml"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"time"
)

type HeadData struct {
	XMLName    xml.Name `xml:"head"`
	ServiceId  string   `xml:"serviceId"`
	CustomerId string   `xml:"customerId"`
	UserId     string   `xml:"userId"`
	ReqTime    string   `xml:"reqTime"`
	SerialNo   string   `xml:"serialNo"`
}

func NewHeadData(serviceId, serialNo, bankCustomerId, bankUserId string) *HeadData {
	return &HeadData{
		ServiceId:  serviceId,
		CustomerId: bankCustomerId,
		UserId:     bankUserId,
		ReqTime:    util.FormatTimeyyyyMMddHHmmss(time.Now()),
		SerialNo:   serialNo,
	}
}

type ResponseHeadData struct {
	HeadData
	OrderFlowNo string `xml:"orderFlowNo"`
	RetCode     string `xml:"retCode"`
	RetMessage  string `xml:"retMessage"`
}

type RequestData struct {
	XMLName xml.Name    `xml:"ebank"`
	Head    HeadData    `xml:"head"`
	Body    interface{} `xml:"body"`
	Sign    string      `xml:"sign"`
}

type BankTransferResponse struct {
	XMLName xml.Name                 `xml:"ebank"`
	Head    ResponseHeadData         `xml:"head"`
	Body    BankTransferResponseBody `xml:"body"`
}

type IntrabankTransferRequest struct {
	PayAccount     string  `xml:"payAccount"`
	PayAccountName string  `xml:"payAccountName"`
	RecAccount     string  `xml:"recAccount"`
	RecAccountName string  `xml:"recAccountName"`
	PayAmount      float64 `xml:"payAmount"`
	PayRem         string  `xml:"payRem"`
	PubPriFlag     string  `xml:"pubPriFlag"`
	CurrencyType   string  `xml:"currencyType"`
	RecBankType    string  `xml:"recBankType"`
	ConfirmCode    string  `xml:"confirmCode"`
}

type IntrabankTransferRequestBody struct {
	IntrabankTransferRequest
	XMLName      xml.Name `xml:"body"`
	BusinessCode string   `xml:"businessCode"`
}

type OutofbankTransferRequest struct {
	PayAccount         string  `xml:"payAccount"`
	PayAccountName     string  `xml:"payAccountName"`
	RecAccount         string  `xml:"recAccount"`
	RecAccountName     string  `xml:"recAccountName"`
	PayAmount          float64 `xml:"payAmount"`
	PayRem             string  `xml:"payRem"`
	PubPriFlag         string  `xml:"pubPriFlag"`
	CurrencyType       string  `xml:"currencyType"`
	RecBankType        string  `xml:"recBankType"`
	ConfirmCode        string  `xml:"confirmCode"`
	TransferFlag       string  `xml:"transferFlag"`
	UnionBankNo        string  `xml:"unionBankNo"`
	ClearBankNo        string  `xml:"clearBankNo"`
	RecAccountOpenBank string  `xml:"recAccountOpenBank"`
	RmtType            string  `xml:"rmtType"`
}

type OutofbankTransferRequestBody struct {
	OutofbankTransferRequest
	XMLName      xml.Name `xml:"body"`
	BusinessCode string   `xml:"businessCode"`
}

type BankTransferResponseBody struct {
	ChargeFee  float64 `xml:"chargeFee"`
	OrderState string  `xml:"orderState"`
}

type SignResult struct {
	XMLName xml.Name       `xml:"html"`
	Head    SignHeadResult `xml:"head"`
	Body    SignBodyResult `xml:"body"`
}

type SignHeadResult struct {
	Title  string `xml:"title"`
	Result string `xml:"result"`
}

type SignBodyResult struct {
	Sign string `xml:"sign"`
}

type TransactionDetailRequestBody struct {
	XMLName         xml.Name `xml:"body"`
	AccountNo       string   `xml:"accountNo"`
	BeginDate       string   `xml:"beginDate"`
	EndDate         string   `xml:"endDate"`
	TurnPageNum     int      `xml:"turnPageNum"`
	TurnPageShowNum int      `xml:"turnPageShowNum"`
	QueryFlag       string   `xml:"queryFlag"`
}

type TransactionDetailResponse struct {
	XMLName xml.Name                      `xml:"ebank"`
	Head    ResponseHeadData              `xml:"head"`
	Body    TransactionDetailResponseBody `xml:"body"`
}

type TransactionDetailResponseBody struct {
	TurnPageTotalNum int64                           `xml:"turnPageTotalNum"`
	Items            []TransactionDetailResponseItem `xml:"list>row"`
}

type TransactionDetailResponseItem struct {
	CashFlag        string  `xml:"cashFlag"`
	PayAmount       float64 `xml:"payAmount"`
	RecAmount       float64 `xml:"recAmount"`
	BsnType         string  `xml:"bsnType"`
	TransferDate    string  `xml:"transferDate"`
	TransferTime    string  `xml:"transferTime"`
	TranChannel     string  `xml:"TranChannel"`
	CurrencyType    string  `xml:"currencyType"`
	Balance         float64 `xml:"balance"`
	OrderFlowNo     string  `xml:"orderFlowNo"`
	HostFlowNo      string  `xml:"hostFlowNo"`
	VouchersType    string  `xml:"vouchersType"`
	VouchersNo      string  `xml:"vouchersNo"`
	SummaryNo       string  `xml:"summaryNo"`
	Summary         string  `xml:"summary"`
	AcctNo          string  `xml:"acctNo"`
	AccountName     string  `xml:"accountName"`
	AccountOpenNode string  `xml:"accountOpenNode"`
}

type TransactionDetailElectronicReceiptRequestBody struct {
	XMLName     xml.Name `xml:"body"`
	OrderFlowNo string   `xml:"orderFlowNo"`
	FileType    string   `xml:"fileType"`
	IsDownload  bool     `xml:"isDownload"`
}

type BatchQueryTransferResultRequestBody struct {
	XMLName          xml.Name `xml:"body"`
	SearchPayAccount string   `xml:"searchPayAccount"`
	SearchRecAccount string   `xml:"searchRecAccount"`
	BeginDate        string   `xml:"beginDate"`
	EndDate          string   `xml:"endDate"`
	BatchNo          string   `xml:"batchNo"`
	OrderFlowNo      string   `xml:"orderFlowNo"`
	TurnPageBeginPos int      `xml:"turnPageBeginPos"`
	TurnPageShowNum  int      `xml:"turnPageShowNum"`
	OrderState       string   `xml:"orderState"`
	BusinessCode     string   `xml:"businessCode"`
}

type TransferListResponse struct {
	Head ResponseHeadData           `xml:"head"`
	Body TransferResultResponseBody `xml:"body"`
}

type TransferResultResponseBody struct {
	TurnPageTotalNum int64                        `xml:"turnPageTotalNum"`
	Items            []TransferResultResponseItem `xml:"list>row"`
}

type TransferResultResponseItem struct {
	BatchNo         string  `xml:"batchNo"`
	OrderFlowNo     string  `xml:"orderFlowNo"`
	BusinessCode    string  `xml:"businessCode"`
	PayAccount      string  `xml:"payAccount"`
	RecAccount      string  `xml:"recAccount"`
	PayAmount       float64 `xml:"payAmount"`
	OrderSubmitTime string  `xml:"orderSubmitTime"`
	OrderState      string  `xml:"orderState"`
	ErrorCode       string  `xml:"errorCode"`
	ErrorMessage    string  `xml:"errorMessage"`
	HostFlowNo      string  `xml:"hostFlowNo"`
	UserId          string  `xml:"userId"`
}

type UploadBatchTransferPayInfoRequest struct {
	TotalNumber    string  `xml:"totalNumber"`
	TotalAmount    float64 `xml:"totalAmount"`
	PayAccount     string  `xml:"payAccount"`
	PayAccountName string  `xml:"payAccountName"`
}

type UploadBatchTransferPayInfoRequestBody struct {
	XMLName xml.Name `xml:"body"`
	UploadBatchTransferPayInfoRequest
}

type UploadBatchTransferPayInfoResponse struct {
	XMLName xml.Name                           `xml:"ebank"`
	Head    ResponseHeadData                   `xml:"head"`
	Body    UploadBatchTransferPayInfoBodyData `xml:"body"`
}

type UploadBatchTransferPayInfoBodyData struct {
	BatchNo string `xml:"batchNo"`
}

type UploadBatchTransferRecInfoList struct {
	SerialNo           string  `xml:"serialNo"`
	RecAccount         string  `xml:"recAccount"`
	RecAccountName     string  `xml:"recAccountName"`
	RecAccountOpenBank string  `xml:"recAccountOpenBank"`
	UnionBankNo        string  `xml:"unionBankNo"`
	PayAmount          float64 `xml:"payAmount"`
	PayUse             string  `xml:"payUse"`
	PayRem             string  `xml:"payRem"`
	BusinessCode       string  `xml:"businessCode"`
	CurrencyType       string  `xml:"currencyType"`
	RmtType            string  `xml:"rmtType"`
	RecBankType        string  `xml:"recBankType"`
	PubPriFlag         string  `xml:"pubPriFlag"`
}

type UploadBatchTransferRecInfoRequestBody struct {
	XMLName xml.Name `xml:"body"`
	UploadBatchTransferRecInfoRequestBodyData
}

type UploadBatchTransferRecInfoRequestBodyData struct {
	BatchNo      string                           `xml:"batchNo"`
	AccountLists []UploadBatchTransferRecInfoList `xml:"list>row"`
}

type UploadBatchTransferRecInfoResponse struct {
	Head ResponseHeadData                           `xml:"head"`
	Body UploadBatchTransferRecInfoResponseBodyData `xml:"body"`
}

type UploadBatchTransferRecInfoResponseBodyData struct {
	BatchNo         string `xml:"batchNo"`
	OrderState      string `xml:"orderState"`
	OrderSubmitTime string `xml:"orderSubmitTime"`
	ErrorCode       string `xml:"errorCode"`
	ErrorMessage    string `xml:"errorMessage"`
}

type QueryAccountBalanceRequestBody struct {
	XMLName xml.Name `xml:"body"`
	QueryAccountBalanceRequest
}

type QueryAccountBalanceRequest struct {
	AccountNo string `xml:"accountNo"`
}

type QueryAccountBalanceResponse struct {
	XMLName xml.Name                        `xml:"ebank"`
	Body    QueryAccountBalanceResponseBody `xml:"body"`
}

type QueryAccountBalanceResponseBody struct {
	XMLName xml.Name `xml:"body"`
	QueryAccountBalanceBodyData
}

type QueryAccountBalanceBodyData struct {
	AccountNo        string  `xml:"accountNo"`
	AccountName      string  `xml:"accountName"`
	Balance          float64 `xml:"balance"`
	BalanceAvailable float64 `xml:"balanceAvailable"`
	BalanceFreeze    float64 `xml:"balanceFreeze"`
	AccStatus        string  `xml:"accStatus"`
	CurrencyType     string  `xml:"currencyType"`
	LastDayBal       string  `xml:"lastDayBal"`
	CshFlg           string  `xml:"cshFlg"`
	QueryDate        string  `xml:"queryDate"`
}

type QueryBatchTransferResultRequestBody struct {
	XMLName          xml.Name `xml:"body"`
	SearchPayAccount string   `xml:"searchPayAccount"`
	BeginDate        string   `xml:"beginDate"`
	EndDate          string   `xml:"endDate"`
	BatchNo          string   `xml:"batchNo"`
	OrderFlowNo      string   `xml:"orderFlowNo"`
	TurnPageBeginPos int      `xml:"turnPageBeginPos"`
	TurnPageShowNum  int      `xml:"turnPageShowNum"`
	OrderState       string   `xml:"orderState"`
	BusinessCode     string   `xml:"businessCode"`
}

type QueryBatchTransferResultRequestBodyData struct {
	SearchPayAccount string `xml:"searchPayAccount"`
	BeginDate        string `xml:"beginDate"`
	EndDate          string `xml:"endDate"`
	BatchNo          string `xml:"batchNo"`
	OrderFlowNo      string `xml:"orderFlowNo"`
	TurnPageBeginPos int    `xml:"turnPageBeginPos"`
	TurnPageShowNum  int    `xml:"turnPageShowNum"`
	OrderState       string `xml:"orderState"`
	BusinessCode     string `xml:"businessCode"`
}

type QueryBatchTransferResultResponse struct {
	Head ResponseHeadData                     `xml:"head"`
	Body QueryBatchTransferResultResponseBody `xml:"body"`
}

type QueryBatchTransferResultResponseBody struct {
	TurnPageTotalNum int64                                 `xml:"turnPageTotalNum"`
	Rows             []QueryBatchTransferResultResponseRow `xml:"list>row"`
}

type QueryBatchTransferResultResponseRow struct {
	BatchNo         string  `xml:"batchNo"`
	OrderFlowNo     string  `xml:"orderFlowNo"`
	BusinessCode    string  `xml:"businessCode"`
	PayAccount      string  `xml:"payAccount"`
	RecAccount      string  `xml:"recAccount"`
	PayAmount       float64 `xml:"payAmount"`
	OrderSubmitTime string  `xml:"orderSubmitTime"`
	OrderState      string  `xml:"orderState"`
	ErrorCode       string  `xml:"errorCode"`
	ErrorMessage    string  `xml:"errorMessage"`
	UserId          string  `xml:"userId"`
}

type BatchQueryLedgerManagementTransferResultRequestBody struct {
	XMLName          xml.Name `xml:"body"`
	PayAccount       string   `xml:"payAccount"`
	SearchRecAccount string   `xml:"searchRecAccount"`
	BeginDate        string   `xml:"beginDate"`
	EndDate          string   `xml:"endDate"`
	BatchNo          string   `xml:"batchNo"`
	OrderFlowNo      string   `xml:"orderFlowNo"`
	TurnPageBeginPos int      `xml:"turnPageBeginPos"`
	TurnPageShowNum  int      `xml:"turnPageShowNum"`
	OrderState       string   `xml:"orderState"`
	RecBankType      string   `xml:"recBankType"`
}

type LedgerManagementTransferListResponse struct {
	Head ResponseHeadData                           `xml:"head"`
	Body LedgerManagementTransferResultResponseBody `xml:"body"`
}

type LedgerManagementTransferResultResponseBody struct {
	TurnPageTotalNum int64                                        `xml:"turnPageTotalNum"`
	Items            []LedgerManagementTransferResultResponseItem `xml:"list>row"`
}

type LedgerManagementTransferResultResponseItem struct {
	BatchNo          string  `xml:"batchNo"`
	OrderFlowNo      string  `xml:"orderFlowNo"`
	BusinessCode     string  `xml:"businessCode"`
	PayAccount       string  `xml:"payAccount"`
	RecAccount       string  `xml:"recAccount"`
	PayAmount        float64 `xml:"payAmount"`
	OrderSubmitTime  string  `xml:"orderSubmitTime"`
	OrderSendTime    string  `xml:"orderSendTime"`
	OrderState       string  `xml:"orderState"`
	RecBankType      string  `xml:"recBankType"`
	PayAccountName   string  `xml:"payAccountName"`
	RecAccountName   string  `xml:"recAccountName"`
	AuthRejectReason string  `xml:"authRejectReason"`
}
