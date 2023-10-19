package stru

//通用接口
type IcbcGlobalRequest struct {
	AppID      string      `json:"app_id"`      // APP的编号,应用在API开放平台注册时生成
	MsgID      string      `json:"msg_id"`      // 消息通讯唯一编号，每次调用独立生成，APP级唯一
	SignType   string      `json:"sign_type"`   // 签名类型，CA-工行颁发的证书认证，RSA-RSAWithSha1，RSA2-RSAWithSha256，缺省为RSA
	Sign       string      `json:"sign"`        // 报文签名
	Timestamp  string      `json:"timestamp"`   // 交易发生时间戳，yyyy-MM-dd HH:mm:ss格式
	BizContent interface{} `json:"biz_content"` // 请求参数的集合
}

//通用接口
type IcbcGlobalResponse struct {
	Sign               string `json:"sign"`                 // 针对返回参数集合的签名
	ResponseBizContent string `json:"response_biz_content"` // 响应参数集合，包含公共和业务参数
}

//签约参数
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

//单账号交易明细查询参数
type AccDetailRequest struct {
	Fseqno    string `json:"fseqno"`    //合作方上送，需保证全局唯一，每次调用校验表里是否重复；建议拼接携带一级合作方编号、调用类型（0单账户流水提取；1流水批量提取；2回单下载）；
	Account   string `json:"account"`   //银行卡号
	Currtype  int    `json:"currtype"`  //币种
	Startdate string `json:"startdate"` //提交起始日期 YYYY-MM-DD
	Enddate   string `json:"enddate"`   //提交结束日期 YYYY-MM-DD
	Serialno  string `json:"serialno"`  //流水号 第一次查询送“”，分页查询下一页需要送上当前页最后一笔明细的序列号，上一页需要送上当前页第一笔明细的序列号
	Corpno    string `json:"corpno"`    //一级合作方编号
	Acccompno string `json:"acccompno"` //二级合作方编号 如为自主记账，可空；
	Agreeno   string `json:"agreeno"`   //协议编号
}

type AccDetailResponse struct {
	Retcode  string            `json:"retcode"`  //返回状态码 0（成功）,-2（失败）,9008100-处理成功-999/9008101-处理失败 9008200-参数错误
	Retmsg   string            `json:"retmsg"`   //返回信息 9008100-处理成功 9008101-处理失败 9008200-参数错误
	Nextpage string            `json:"nextpage"` //0无下一页；1有下一页
	Serialno string            `json:"serialno"` //当前页最后一笔明细的流水号
	Accno    string            `json:"accno"`    //本方账号
	DtlList  IcbcAccDetailItem `json:"dtllist"`  //流水明细表
}

type IcbcAccDetailItem struct {
	Serialno  string `json:"serialno"`  //流水号
	Busidate  string `json:"busidate"`  //入账日期
	Busitime  string `json:"busitime"`  //入账时间
	Trxcode   string `json:"trxcode"`   //交易代码
	Detailf   string `json:"detailf"`   //明细性质
	Drcrf     string `json:"drcrf"`     //借贷标记
	Vouhtype  string `json:"vouhtype"`  //凭证种类
	Vouhno    string `json:"vouhno"`    //凭证号
	Summary   string `json:"summary"`   //摘要
	Amount    string `json:"amount"`    //金额
	Currtype  string `json:"currtype"`  //币种
	Balance   string `json:"balance"`   //当前余额
	Workdate  string `json:"workdate"`  //工作日期
	Valueday  string `json:"valueday"`  //调整起息日期
	Statcode  string `json:"statcode"`  //外汇统计代码
	Settmode  string `json:"settmode"`  //外汇结算方式
	Actcode   string `json:"actcode"`   //账户核算机构号
	Tellerno  string `json:"tellerno"`  //柜员号
	Authtlno  string `json:"authtlno"`  //授权柜员号
	Authno    string `json:"authno"`    //授权代号
	Termid    string `json:"termid"`    //终端号
	Recipacc  string `json:"recipacc"`  //对方账号
	Recipnam  string `json:"recipnam"`  //对方户名
	Crvouhtyp string `json:"crvouhtyp"` //对方凭证种类
	Crvouhno  string `json:"crvouhno"`  //对方凭证号
	VagenRef  string `json:"vagen_ref"` //业务编号
	Oref      string `json:"oref"`      //相关业务编号
	Drbuscode string `json:"drbuscode"` //借方业务代码
	Crbuscode string `json:"crbuscode"` //贷方业务代码
	Ensummry  string `json:"ensummry"`  //英文备注
	Trantel   string `json:"trantel"`   //经办柜员号
	Importel  string `json:"importel"`  //录入柜员号
	Checktel  string `json:"checktel"`  //复核柜员号
	Recipcno  string `json:"recipcno"`  //对方客户编号
	Recipbkn  string `json:"recipbkn"`  //对方行号
	Recipbna  string `json:"recipbna"`  //对方行名
	Opertype  string `json:"opertype"`  //网银业务种类
	Notes     string `json:"notes"`     //附言
	Purpose   string `json:"purpose"`   //用途
	Servface  string `json:"servface"`  //服务界面
	Eventseq  string `json:"eventseq"`  //大交易序号
	Ptrxseq   string `json:"ptrxseq"`   //小交易序号
	Updtranf  string `json:"updtranf"`  //冲正标志
	Revtranf  string `json:"revtranf"`  //正反交易标志

}
