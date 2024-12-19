package stru

type MinShengBaseResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
}
type MinShengTransferRequest struct {
	ReqSeq            string `json:"reqSeq"`            // 请求流水号
	AcctNo            string `json:"acctNo"`            // 账号
	PayType           string `json:"payType"`           // 支付方式 1直接支付，2 附条件支付 暂时支持1
	IsCross           string `json:"isCross"`           // 0-他行 1-本行
	Currency          string `json:"currency"`          // 仅仅支持人民币 CNY
	TransAmt          string `json:"transAmt"`          // 转账金额
	BankRoute         string `json:"bankRoute"`         // 汇路
	BankCode          string `json:"BankCode"`          //开户行号
	BankName          string `json:"BankName"`          // 开户行名
	OpenId            string `json:"openId"`            // 授权码
	Usage             string `json:"usage"`             // 用途
	CertNo            string `json:"certNo"`            //企业自制凭证号,可用于hostflow字段的填写
	PayeeAcctNo       string `json:"payeeAcctNo"`       //收款账号
	PayeeAcctName     string `json:"payeeAcctName"`     //收款账户名称
	PublicPrivateFlag string `json:"publicPrivateFlag"` // 对公对私 0-对公 1-对私
}

type MinShengBankTransferResponse struct {
	MinShengBaseResponse
	ReqSeq     string `json:"reqSeq"`     // 业务请求流水号
	RespDate   string `json:"respDate"`   // 交易日期
	RespTime   string `json:"respTime"`   // 交易时间
	RespStatus string `json:"respStatus"` // 执行结果 S-成功 E- 失败
	RespCode   string `json:"respCode"`   // 交易状态码
	RespDesc   string `json:"respDesc"`   // 返回信息
}

type MinShengSingleTransferQueryResponse struct {
	MinShengBaseResponse
	ReqSeq      string `json:"req_seq"`       //业务请求流水号
	TransSeqNo  string `json:"trans_seq_no"`  //交易流水号
	RespDate    string `json:"resp_date"`     //交易日期
	RespTime    string `json:"resp_time"`     //交易时间
	RespStatus  string `json:"resp_status"`   //执行结果 S-成功 E- 失败
	RespCode    string `json:"resp_code"`     //交易状态码 AAAAAAA 受理成功
	RespDesc    string `json:"resp_desc"`     //返回信息
	TransStatus string `json:"trans_status"`  //交易状态 1-审批中 2-审批拒绝/失败 3-交易成功（跨行交易时仅代表银行提交清算机构成功，不代表实际入账成功）4-交易支付中 5-交易失败
	TransDate   string `json:"trans_date"`    // 交易执行时间 用以下载回单使用
	TransErrMsg string `json:"trans_err_msg"` // 失败原因
}

type MinShengSingleTransferDetailResponseItem struct {
	TransSeqNo    string `json:"trans_seq_no"`    //交易流水号
	RecSeq        string `json:"rec_seq"`         //记账流水号
	AcctNo        string `json:"acct_no"`         //账号
	AcctName      string `json:"acct_name"`       //户名
	DcFlag        string `json:"dc_flag"`         //1-借方（支出） 2-贷方（收入）
	EnterAcctDate string `json:"enter_acct_date"` //入账日期
	IntrDate      string `json:"intr_date"`       //起息日期
	Amount        string `json:"amount"`          //金额
	CpAcctNo      string `json:"cp_acct_no"`      //对方账号
	CpAcctName    string `json:"cp_acct_name"`    //对方户名
	CpBankName    string `json:"cp_bank_name"`    //对方开户行名称
	CpBankAddr    string `json:"cp_bank_addr"`    //对方开户行地址
	Explain       string `json:"explain"`         // 用途/摘要
	Balance       string `json:"balance"`         // 余额
	Timestamp     string `json:"timestamp"`       // 时间戳
	Currency      string `json:"currency"`        // 币种
}

type MinShengResponse struct {
	ReturnCode   string      `json:"return_code"`
	ReturnMsg    string      `json:"return_msg"`
	ResponseBusi interface{} `json:"response_busi"` // 响应内容
}

type MinShengSingleTransferDetailResponse struct {
	MinShengBaseResponse
	ReqSeq     string                                      `json:"req_seq"`     // 业务请求流水号
	RespDate   string                                      `json:"resp_date"`   // 交易日期
	RespTime   string                                      `json:"resp_time"`   // 交易时间
	RespStatus string                                      `json:"resp_status"` // 执行结果 S-成功 E- 失败
	RespCode   string                                      `json:"resp_code"`   // 交易状态码 AAAAAAA 受理成功
	RespDesc   string                                      `json:"resp_desc"`   // 返回信息
	ItemNum    string                                      `json:"item_num"`    // 本次记录数
	TotalNum   string                                      `json:"total_num"`   // 总记录数
	Items      []*MinShengSingleTransferDetailResponseItem `json:"result_list"` // 明细列表
}

type MinShengElectronicReceiptResponse struct {
	MinShengBaseResponse
	ReqSeq      string `json:"req_seq"`      // 业务请求流水号
	RespDate    string `json:"resp_date"`    // 交易日期
	RespTime    string `json:"resp_time"`    // 交易时间
	RespStatus  string `json:"resp_status"`  // 执行结果 S-成功 E- 失败
	RespCode    string `json:"resp_code"`    // 交易状态码 AAAAAAA 受理成功
	RespDesc    string `json:"resp_desc"`    // 返回信息
	FileContent []byte `json:"file_content"` // 电子回单文件内容 经 Base64加密之后的字符窜
	FileName    string `json:"file_name"`    // 电子回单文件名
	FileType    string `json:"file_type"`    // 文件类型
}

type MinShengAuthResponse struct {
	MinShengBaseResponse
	ReqSeq     string `json:"req_seq"`     // 业务请求流水号
	RespDate   string `json:"resp_date"`   // 交易日期
	RespTime   string `json:"resp_time"`   // 交易时间
	RespStatus string `json:"resp_status"` // 执行结果 S-成功 E- 失败
	RespCode   string `json:"resp_code"`   // 交易状态码 AAAAAAA 受理成功
	RespDesc   string `json:"resp_desc"`   // 返回信息
	AuthUrl    string `json:"auth_url"`    // 授权URL
}

type MinShengAuthStatusResponse struct {
	MinShengBaseResponse
	ReqSeq     string `json:"req_seq"`     // 业务请求流水号
	RespDate   string `json:"resp_date"`   // 交易日期
	RespTime   string `json:"resp_time"`   // 交易时间
	RespStatus string `json:"resp_status"` // 执行结果 S-成功 E- 失败
	RespCode   string `json:"resp_code"`   // 交易状态码 AAAAAAA 受理成功
	RespDesc   string `json:"resp_desc"`   // 返回信息
	Status     string `json:"status"`      // 授权状态 0-已解约 1-成功 2-待授权 3-已过期 4-待审批
	OpenId     string `json:"open_id"`     // 授权码
	StartTime  string `json:"start_time"`  // 授权开始时间
	EndTime    string `json:"end_time"`    // 授权结束时间
	AuthCode   string `json:"auth_code"`   // 授权标识
}
