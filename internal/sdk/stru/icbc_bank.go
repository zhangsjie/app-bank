package stru

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ICBCPostHttpResult(host string, requestData IcbcGlobalRequest, result *interface{}) error {
	//生成privateKey
	privateKey, err := parsePrivateKey([]byte(config.GetString(enum.IcbcPrivateKey, "")))
	if err != nil {
		// 处理错误
	}
	// 生成签名并设置到请求结构体的 "sign" 字段
	sign, err := generateRSASign(requestData, privateKey)
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult 生成签名出错requestData=%+v,privateKey=%+v", requestData, privateKey))
		return nil
	}
	requestData.Sign = sign
	// 将请求结构体转换为 URL 编码的字符串，并作为请求体发送
	bodyData, err := encodeStructToURLValues(requestData)
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult encodeStructToURLValues转换编码出错requestData=%+v", requestData))
	}
	req, err := http.NewRequest("POST", host, strings.NewReader(bodyData))
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult httpRequest生成出错bodyData=%+v,host=%+v", bodyData, host))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求并处理响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult httpRequest请求工行出错bodyData=%+v,resp=%+v,err=%+v", bodyData, resp, err))
	}

	zap.L().Info(fmt.Sprintf("ICBCPostResult httpRequest请求工行成功bodyData=%+v,resp=%+v,", bodyData, resp))
	defer resp.Body.Close()

	// 处理响应
	responseData := IcbcGlobalResponse{
		ResponseBizContent: result,
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)

	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult 解析响应出错：%v", err))
		return err
	}

	zap.L().Info(fmt.Sprintf("ICBCPostResult 请求工行成功：bodyData=%+v, resp=%+v", bodyData, resp))
	return nil
}

func encodeStructToURLValues(data interface{}) (string, error) {
	values := url.Values{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(jsonData, &values)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

func generateRSASign(data interface{}, privateKey *rsa.PrivateKey) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(jsonData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
func NewIcbcGlobalRequest() *IcbcGlobalRequest {
	msgId, _ := util.SonyflakeID()
	request := IcbcGlobalRequest{
		AppID:     config.GetString(enum.IcbcAppId, ""),
		MsgID:     msgId,
		SignType:  "RSA2",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	return &request
}

func parsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA private key")
	}

	return rsaPrivateKey, nil
}

// IcbcGlobalRequest 通用接口
type IcbcGlobalRequest struct {
	AppID      string      `json:"app_id"`      // APP的编号,应用在API开放平台注册时生成
	MsgID      string      `json:"msg_id"`      // 消息通讯唯一编号，每次调用独立生成，APP级唯一
	SignType   string      `json:"sign_type"`   // 签名类型，CA-工行颁发的证书认证，RSA-RSAWithSha1，RSA2-RSAWithSha256，缺省为RSA
	Sign       string      `json:"sign"`        // 报文签名
	Timestamp  string      `json:"timestamp"`   // 交易发生时间戳，yyyy-MM-dd HH:mm:ss格式
	BizContent interface{} `json:"biz_content"` // 请求参数的集合
}

// IcbcGlobalResponse 通用返回结构
type IcbcGlobalResponse struct {
	Sign               string      `json:"sign"`                 // APP的编号,应用在API开放平台注册时生成
	ResponseBizContent interface{} `json:"response_biz_content"` // 消息通讯唯一编号，每次调用独立生成，APP级唯一
}

// IcbcSignRequest 签约参数
type IcbcSignRequest struct {
	AppID      string        `json:"app_id"`               // APP的编号,应用在API开放平台注册时生成
	ApiName    string        `json:"api_name"`             // api接口名称
	ApiVersion string        `json:"api_version"`          // api接口版本
	CorpNo     string        `json:"corpNo"`               // 一级合作方编号
	CoMode     string        `json:"coMode"`               // 合作模式，1：代理记账；2：自主记账
	AccCompNo  string        `json:"accCompNo,omitempty"`  // 二级合作方编号，合作模式为1时，必输 (可选字段)
	Account    string        `json:"account"`              // 主账户
	CurrType   string        `json:"currType"`             // 主账户币种，1：人民币
	AccFlag    string        `json:"accFlag"`              // 本他行标志，1：本行；2：他行
	CnTioFlag  string        `json:"cntioFlag"`            // 境内外标志，1：境内；2：境外
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
	CnTioFlag     string `json:"cntioFlag"`     // 境内外标志，1：境内；2：境外
	IsMainAcc     string `json:"isMainAcc"`     // 主帐户标志，0：否；1：是
	ReceiptFlag   string `json:"receiptFlag"`   // 开通电子单据标示，0：关；1：开
	StatementFlag string `json:"statementFlag"` // 开通账务明细查询标示，0：关；1：开
}

// AccDetailRequest 单账号交易明细查询参数
type AccDetailRequest struct {
	FSeqNo    string `json:"fseqno"`    //合作方上送，需保证全局唯一，每次调用校验表里是否重复；建议拼接携带一级合作方编号、调用类型（0单账户流水提取；1流水批量提取；2回单下载）；
	Account   string `json:"account"`   //银行卡号
	CurrType  int    `json:"currtype"`  //币种
	StartDate string `json:"startdate"` //提交起始日期 YYYY-MM-DD
	EndDate   string `json:"enddate"`   //提交结束日期 YYYY-MM-DD
	SerialNo  string `json:"serialno"`  //流水号 第一次查询送“”，分页查询下一页需要送上当前页最后一笔明细的序列号，上一页需要送上当前页第一笔明细的序列号
	CorpNo    string `json:"corpno"`    //一级合作方编号
	AccCompNo string `json:"acccompno"` //二级合作方编号 如为自主记账，可空；
	AgreeNo   string `json:"agreeno"`   //协议编号
}

type AccDetailResponse struct {
	RetCode  string            `json:"retcode"`  //返回状态码 0（成功）,-2（失败）,9008100-处理成功-999/9008101-处理失败 9008200-参数错误
	RetMsg   string            `json:"retmsg"`   //返回信息 9008100-处理成功 9008101-处理失败 9008200-参数错误
	NextPage string            `json:"nextpage"` //0无下一页；1有下一页
	SerialNo string            `json:"serialno"` //当前页最后一笔明细的流水号
	AccNo    string            `json:"accno"`    //本方账号
	DtlList  IcbcAccDetailItem `json:"dtllist"`  //流水明细表
}

type IcbcAccDetailItem struct {
	SerialNo  string `json:"serialno"`  //流水号
	BusiDate  string `json:"busidate"`  //入账日期
	BusiTime  string `json:"busitime"`  //入账时间
	TrxCode   string `json:"trxcode"`   //交易代码
	DetailF   string `json:"detailf"`   //明细性质
	DrcrF     string `json:"drcrf"`     //借贷标记
	VouhType  string `json:"vouhtype"`  //凭证种类
	VouhNo    string `json:"vouhno"`    //凭证号
	Summary   string `json:"summary"`   //摘要
	Amount    string `json:"amount"`    //金额
	CurrType  string `json:"currtype"`  //币种
	Balance   string `json:"balance"`   //当前余额
	WorkDate  string `json:"workdate"`  //工作日期
	ValueDay  string `json:"valueday"`  //调整起息日期
	StatCode  string `json:"statcode"`  //外汇统计代码
	SettMode  string `json:"settmode"`  //外汇结算方式
	ActCode   string `json:"actcode"`   //账户核算机构号
	TellerNo  string `json:"tellerno"`  //柜员号
	AuthtlNo  string `json:"authtlno"`  //授权柜员号
	AuthNo    string `json:"authno"`    //授权代号
	TermId    string `json:"termid"`    //终端号
	RecipAcc  string `json:"recipacc"`  //对方账号
	RecipNam  string `json:"recipnam"`  //对方户名
	CrvouHtyp string `json:"crvouhtyp"` //对方凭证种类
	CrvouhNo  string `json:"crvouhno"`  //对方凭证号
	VagenRef  string `json:"vagen_ref"` //业务编号
	OreF      string `json:"oref"`      //相关业务编号
	DrbusCode string `json:"drbuscode"` //借方业务代码
	CrbusCode string `json:"crbuscode"` //贷方业务代码
	EnSummry  string `json:"ensummry"`  //英文备注
	TranTel   string `json:"trantel"`   //经办柜员号
	Importel  string `json:"importel"`  //录入柜员号
	CheckTel  string `json:"checktel"`  //复核柜员号
	RecipcNo  string `json:"recipcno"`  //对方客户编号
	RecipBkn  string `json:"recipbkn"`  //对方行号
	RecipBna  string `json:"recipbna"`  //对方行名
	OperType  string `json:"opertype"`  //网银业务种类
	Notes     string `json:"notes"`     //附言
	Purpose   string `json:"purpose"`   //用途
	ServFace  string `json:"servface"`  //服务界面
	EventSeq  string `json:"eventseq"`  //大交易序号
	PtrxSeq   string `json:"ptrxseq"`   //小交易序号
	UpdTranf  string `json:"updtranf"`  //冲正标志
	RevTranf  string `json:"revtranf"`  //正反交易标志
}
