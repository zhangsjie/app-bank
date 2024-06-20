package stru

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func ICBCPostHttpResult(urlPath string, requestData IcbcGlobalRequest) (interface{}, error) {
	host := config.GetString(enum.IcbcHost, "")
	//生成privateKey
	privateKey, _ := parsePrivateKey(config.GetString(enum.IcbcPrivateKey, ""))
	// 将 BizContent 转换为 JSON 字符串
	bizContent, err := json.Marshal(requestData.BizContent)
	if err != nil {
		return nil, err
	}
	//生成签名字符串
	signString := urlPath + "?app_id=" + requestData.AppID + "&biz_content=" + string(bizContent) + "&msg_id=" + requestData.MsgID + "&sign_type=" + requestData.SignType + "&timestamp=" + requestData.Timestamp

	// 生成签名
	sign, err := generateRSASign([]byte(signString), privateKey)
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostHttpResult 生成签名字符串失败err=%+v", err))
		return nil, err
	}
	//组装请求参数

	baseUrl, _ := url.Parse(host)
	baseUrl.Path = urlPath
	//signData := s + "&sign=" + sign
	// 进行Base64编码
	params := url.Values{}
	params.Add("app_id", requestData.AppID)
	params.Add("msg_id", requestData.MsgID)
	params.Add("sign_type", requestData.SignType)
	params.Add("timestamp", requestData.Timestamp)
	params.Add("sign", sign)
	//params.Add("biz_content", string(bizContent))
	baseUrl.RawQuery = params.Encode()
	paramsContext := url.Values{}

	paramsContext.Add("biz_content", string(bizContent))

	zap.L().Info(fmt.Sprintf("ICBCPostHttp request==\n %+v,,body=\n %+v", baseUrl.String(), paramsContext))

	resp, err := http.PostForm(baseUrl.String(), paramsContext)
	if err != nil {
		zap.L().Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostHttpResult t请求工行出错bodyData=%+v,resp=%+v,err=%+v", baseUrl.String(), resp, err))
	}

	zap.L().Info(fmt.Sprintf("ICBCPostHttp response==\n %+v", string(body)))

	// 处理响应
	var response IcbcGlobalResponse
	//responseData := IcbcGlobalResponse{
	//	ResponseBizContent: result,
	//}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&response)

	if err != nil {
		zap.L().Info(fmt.Sprintf("ICBCPostResult 解析响应出错：%v", err))
		return nil, err
	}

	return response.ResponseBizContent, nil
}

func SftpClient() *sftp.Client {
	host := config.GetString("bankConfig.icbc.sftpHost", "")
	port := config.GetString("bankConfig.icbc.sftPort", "")
	userName := config.GetString("bankConfig.icbc.userName", "")
	ip := host + ":" + port
	// 读取私钥文件
	keyBytes, err := ioutil.ReadFile("id_rsa.txt")
	if err != nil {
		fmt.Errorf("failed to decode Base64 encoded private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		fmt.Println("读取私钥文件失败：", err)
		return nil
	}

	// 创建SSH配置，使用自定义的HostKeyCallback
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
		},
		//HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//	return checkHostKey(hostname, remote, knownHosts, key)
		//},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境请替换为有效的HostKeyCallback
	}

	// 连接SSH服务器
	conn, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		fmt.Println("Failed to dial SSH:", err)
	}
	// 建立SFTP会话
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		fmt.Println("Failed to create SFTP client:", err)
	}
	return sftpClient
}

func generateRSASign(data []byte, privateKey *rsa.PrivateKey) (string, error) {

	digest := sha256.Sum256(data)
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

func parsePrivateKey(privateKeyString string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Base64 encoded private key: %w", err)
	}

	//block, _ := pem.Decode(privateKeyBytes)
	//if block == nil {
	//	return nil, fmt.Errorf("failed to parse PEM block containing the key")
	//}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
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

type AccDetailRequest struct {
	FSeqNo    string `json:"fseqno"`    //合作方上送，需保证全局唯一，每次调用校验表里是否重复；建议拼接携带一级合作方编号、调用类型（0单账户流水提取；1流水批量提取；2回单下载）；
	Account   string `json:"account"`   //银行卡号
	CurrType  string `json:"currtype"`  //币种
	StartDate string `json:"startdate"` //提交起始日期 YYYY-MM-DD
	EndDate   string `json:"enddate"`   //提交结束日期 YYYY-MM-DD
	SerialNo  string `json:"serialno"`  //流水号 第一次查询送“”，分页查询下一页需要送上当前页最后一笔明细的序列号，上一页需要送上当前页第一笔明细的序列号
	CorpNo    string `json:"corpno"`    //一级合作方编号
	AccCompNo string `json:"acccompno"` //二级合作方编号 如为自主记账，可空；
	AgreeNo   string `json:"agreeno"`   //协议编号
}

type AccDetailResponse struct {
	RetCode  string              `json:"retcode"`  //返回状态码 0（成功）,-2（失败）,9008100-处理成功-999/9008101-处理失败 9008200-参数错误
	RetMsg   string              `json:"retmsg"`   //返回信息 9008100-处理成功 9008101-处理失败 9008200-参数错误
	NextPage string              `json:"nextpage"` //0无下一页；1有下一页
	SerialNo string              `json:"serialno"` //当前页最后一笔明细的流水号
	AccNo    string              `json:"accno"`    //本方账号
	DtlList  []IcbcAccDetailItem `json:"dtllist"`  //流水明细表
}

type IcbcAccDetailItem struct {
	SerialNo   int64   `json:"serialno"`   //流水号
	BusiDate   string  `json:"busidate"`   //入账日期
	BusiTime   string  `json:"busitime"`   //入账时间
	TrxCode    int64   `json:"trxcode"`    //交易代码
	DetailF    int     `json:"detailf"`    //明细性质
	DrcrF      int     `json:"drcrf"`      //借贷标记 1-借，2-贷，借是出账，贷是入账
	VouhType   int     `json:"vouhtype"`   //凭证种类
	VouhNo     int     `json:"vouhno"`     //凭证号
	Summary    string  `json:"summary"`    //摘要
	Amount     float64 `json:"amount"`     //金额
	CurrType   int     `json:"currtype"`   //币种
	Balance    float64 `json:"balance"`    //当前余额
	WorkDate   string  `json:"workdate"`   //工作日期
	ValueDay   string  `json:"valueday"`   //调整起息日期
	StatCode   int     `json:"statcode"`   //外汇统计代码
	SettMode   int     `json:"settmode"`   //外汇结算方式
	ActCode    int64   `json:"actcode"`    //账户核算机构号
	TellerNo   int64   `json:"tellerno"`   //柜员号
	AuthtlNo   int64   `json:"authtlno"`   //授权柜员号
	AuthNo     string  `json:"authno"`     //授权代号
	TermId     string  `json:"termid"`     //终端号
	RecipAcc   string  `json:"recipacc"`   //对方账号
	RecipNam   string  `json:"recipnam"`   //对方户名
	CrvouhType int     `json:"crvouhtype"` //对方凭证种类
	CrvouhNo   int     `json:"crvouhno"`   //对方凭证号
	VagenRef   string  `json:"vagen_ref"`  //业务编号
	OreF       string  `json:"oref"`       //相关业务编号
	DrbusCode  string  `json:"drbuscode"`  //借方业务代码
	CrbusCode  string  `json:"crbuscode"`  //贷方业务代码
	EnSummry   string  `json:"ensummry"`   //英文备注
	TranTel    int64   `json:"trantel"`    //经办柜员号
	Importel   int64   `json:"importel"`   //录入柜员号
	CheckTel   int64   `json:"checktel"`   //复核柜员号
	RecipcNo   int64   `json:"recipcno"`   //对方客户编号
	RecipBkn   string  `json:"recipbkn"`   //对方行号
	RecipBna   string  `json:"recipbna"`   //对方行名
	OperType   string  `json:"opertype"`   //网银业务种类
	Notes      string  `json:"notes"`      //附言
	Purpose    string  `json:"purpose"`    //用途
	ServFace   int64   `json:"servface"`   //服务界面
	EventSeq   int64   `json:"eventseq"`   //大交易序号
	PtrxSeq    int64   `json:"ptrxseq"`    //小交易序号
	UpdTranf   int     `json:"updtranf"`   //冲正标志
	RevTranf   int     `json:"revtranf"`   //正反交易标志
	// 新增的字段
	BusTypNo  int    `json:"busTypNo"`
	Cardno    string `json:"cardno"`
	ReciPna2  string `json:"reciPna2"`
	ReciPac2  string `json:"reciPac2"`
	Trxorgann string `json:"trxorgann"`
	RecipNams string `json:"recipNams"`
	OtactNo1  string `json:"otactNo1"`
	OtactNam1 string `json:"otactNam1"`
	CnlType   int    `json:"cnlType"`
	FinFo     string `json:"finFo"`
	NeagnFlg  string `json:"neagnFlg"`
	seFlag    string `json:"useFlag"`
}

type QueryAgreeNoRequest struct {
	Inqwork Inqwork `json:"inqwork"` // 翻页控制
	Corpno  string  `json:"corpno"`  // 一级合作方编号
	Cond    Cond    `json:"cond"`    // 查询条件
}
type Inqwork struct {
	BegNum string `json:"begnum"` // 查询起始位置，0开始
	FetNum string `json:"fetnum"` // 查询条数，最多支持10条，后同
}
type Cond struct {
	QryType   string        `json:"qrytype"`   // 查询类型：1按合作企业/合作企业+代账公司查询（可输入主账户），2按协议列表方式查询
	AccCompNo string        `json:"acccompno"` // 二级合作方编号，qrytype为1时选输
	Account   string        `json:"account"`   // 主账户，qrytype为1时选输
	CurrType  string        `json:"currtype"`  // 币种，qrytype为1时选输
	AgrList   []interface{} `json:"agrlist"`   // 协议列表，qrytype为2时必输，最多10条
}

type QueryAgreeNoResponse struct {
	RetCode string      `json:"retCode"` // 返回状态码 9008100-处理成功 9008101-处理失败 9008200-参数错误
	RetMsg  string      `json:"retMsg"`  // 返回信息  9008100-处理成功 9008101-处理失败 9008200-参数错误
	MsgID   string      `json:"msgid"`   // 消息通讯唯一编号，每次调用独立生成，APP级唯一
	AgrList []Agreement `json:"agrlist"` // 返回协议列表
}
type Agreement struct {
	AgreeNo         int64         `json:"agreeno"`         // 协议号
	CorpNo          string        `json:"corpno"`          // 一级合作方编号
	CoMode          string        `json:"comode"`          // 合作模式
	AccCompNo       string        `json:"acccompno"`       // 二级合作方编号
	Phone           string        `json:"phone"`           // 手机号
	Status          string        `json:"status"`          // 协议状态
	ActDate         string        `json:"actdate"`         // 生效日期
	MatDate         string        `json:"matdate"`         // 到期日期
	EpType          string        `json:"eptype"`          // 是否自动展期
	EpTimes         int64         `json:"eptimes"`         // 展期期数
	FeeFlag         string        `json:"feeflag"`         // 收费标识
	Remark          string        `json:"remark"`          // 备注
	LstModft        string        `json:"lstmodft"`        // 最后修改时间
	AccountInfoList []interface{} `json:"accountinfolist"` // 账户信息列表
}

type IcbcReceiptNoQueryRequest struct {
	FseqNo    string                 `json:"fseqno" `   // 请求序列号
	CorpNo    string                 `json:"corpno"`    // 一级合作方编号
	AccCompNo string                 `json:"acccompno"` // 二级合作方编号
	AgreeNo   string                 `json:"agreeno" `  // 协议编号
	Account   string                 `json:"account"`   // 账号/卡号
	CurrType  string                 `json:"currtype"`  // 币种
	QryType   string                 `json:"qrytype" `  // 查询类型
	QryCond   IcbcReceiptNoQueryCond `json:"qrycond"`   // 查询条件
}
type IcbcReceiptNoQueryCond struct {
	StartDate string   `json:"startDate" ` // 提交起始日期，qrytype为1时必输
	EndDate   string   `json:"endDate"`    // 提交结束日期，qrytype为1时必输
	SeqList   []string `json:"seqlist"`    // 回单编号列表，qrytype为2时必输，最多50条
}
type IcbcReceiptNoQueryResponse struct {
	RetCode string `json:"retcode" ` // 9008100-处理成功 9008101-处理失败 9008200-参数错误
	RetMsg  string `json:"retmsg"`   //
	OrderId string `json:"orderid"`  // 根据orderid在共享目录搜索回单文件，订单编号示例：“2021082711075662E0D0D48C860002E0537A196D225495”
}
