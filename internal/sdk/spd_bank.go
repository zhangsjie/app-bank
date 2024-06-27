package sdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-common/enum"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SPDBankSDK interface {
	/*IntrabankTransfer(ctx context.Context, serialNo string, req stru.IntrabankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error)
	OutofbankTransfer(ctx context.Context, serialNo string, req stru.OutofbankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error)
	UploadBatchTransferPayInfo(ctx context.Context, req stru.UploadBatchTransferPayInfoRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferPayInfoResponse, error)
	UploadBatchTransferRecInfo(ctx context.Context, batchNo string, req []stru.UploadBatchTransferRecInfoList, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferRecInfoResponse, error)
	QueryBatchTransferResult(ctx context.Context, req stru.QueryBatchTransferResultRequestBodyData, host, sighHost, bankCustomerId, bankUserId string) ([]stru.QueryBatchTransferResultResponseRow, error)*/
	BankTransfer(ctx context.Context, serialNo string, req stru.SPDBankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferResponseData, error)
	ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.SPDBankTransferDetailResponseItem, error)
	RequestTransactionDetailElectronicReceipt(ctx context.Context, accountNo, beginDate, endDate, backhostGyno, subpoenaSeqNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferDetailElectronicReceiptResponseItem, error)
	DownloadTransactionDetailElectronicReceipt(ctx context.Context, accountNo, beginDate, endDate, backhostGyno, subpoenaSeqNo, host, sighHost, fileHost, bankCustomerId, bankUserId string) ([]byte, error)
	QueryTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo, elecChequeNo string, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferResultResponseData, error)
	QueryAccountBalance(ctx context.Context, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankQueryAccountBalanceResponseData, error)
	CreateVirtualAccount(ctx context.Context, req []stru.SPDCreateVirtualAccountRequestItem, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDCreateVirtualAccountResponseData, error)
	QueryVirtualAccountBalance(ctx context.Context, accountNo string, virtualAccountNo []string, offset int, host, sighHost, bankCustomerId, bankUserId string) ([]stru.SPDQueryVirtualAccountBalanceResponseItem, error)
	BankVirtualAccountTransfer(ctx context.Context, serialNo string, req stru.SPDBankVirtualAccountTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankVirtualAccountTransferResponseData, error)
}

type spdBankSDK struct{}

// spdSign
//  @Description: 	获取签名
//  @param ctx
//  @param head		请求头
//  @param body		请求body
//  @param sighHost	验签地址
//  @return string
//
func spdSign(ctx context.Context, body interface{}, sighHost string) string {
	result, err := util.PostGBKXmlHttpStringResult(ctx, sighHost,
		"INFOSEC_SIGN/1.0", false, body)
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}
	signature := result[strings.Index(result, "<sign>")+6 : strings.Index(result, "</sign>")]
	signature = strings.ReplaceAll(signature, "\r", "")
	signature = strings.ReplaceAll(signature, "\n", "")
	return signature
}

// spdVerifySign
//  @Description: 验证签名接口 解析签名
//  @param ctx
//  @param sign
//  @param sighHost
//  @return string
//
func spdVerifySign(ctx context.Context, sign string, sighHost string) string {
	result, err := util.PostStringHttpStringResult(ctx, sighHost,
		"INFOSEC_VERIFY_SIGN/1.0", sign)
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}
	return result[strings.Index(result, "<sic>")+5 : strings.Index(result, "</sic>")]
}

// BankTransfer
//  @Description:			转账
//  @receiver s
//  @param ctx
//  @param serialNo 		对应企业发送报文的报文号 当天内唯一
//  @param req				转账请求体
//  @param host				nc_http 服务地址
//  @param sighHost			nc_sign 验签地址
//  @param bankCustomerId	浦发企业客户号
//  @param bankUserId		浦发付款名称
//  @return *stru.BankTransferResponse
//  @return error
//
func (s *spdBankSDK) BankTransfer(ctx context.Context, serialNo string, req stru.SPDBankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferResponseData, error) {
	if req.AcctName == "" {
		req.AcctName = bankUserId
	}
	// 根据body获取sign
	body := &stru.SPDBankTransferRequestBody{
		SPDBankTransferRequest: req,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankTransferServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	// 返回的sign对象
	var responseSignData stru.SPDResponseSignData
	err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
		zap.L().Info(fmt.Sprintf("浦发服务: %s 转账验签异常code: %s msg:%s\n", enum.SPDBankTransferServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
		return nil, nil
	}
	// 最终返回的对象
	var bankTransferResponseData stru.SPDBankTransferResponseData
	// 设置签名请求后的返回信息
	bankTransferResponseData.ReturnCode = responseSignData.Body.ReturnCode
	bankTransferResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
	// 拿到返回的sign, 在解析作为返参
	responseSign := responseSignData.Body.Sign
	if responseSign != "" {
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &bankTransferResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
	}
	return &bankTransferResponseData, nil
}

//
// ListTransactionDetail
//  @Description: 			查询交易明细
//  @receiver s
//  @param ctx
//  @param accountNo		账号
//  @param beginDate		开始日期
//  @param endDate			结束日期
//  @param host				nc_http 服务地址
//  @param sighHost			nc_sign 验签地址
//  @param bankCustomerId	浦发企业客户号
//  @param bankUserId		浦发付款名称
//  @return []stru.TransactionDetailResponseItem
//  @return error
//
func (s *spdBankSDK) ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.SPDBankTransferDetailResponseItem, error) {
	pageNum := 1
	pageSize := 5
	var result []stru.SPDBankTransferDetailResponseItem
	for {
		serialNo, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		beginNumber := (pageNum-1)*pageSize + 1
		// 根据body获取sign
		body := &stru.SPDBankTransferDetailRequestBody{
			SPDBankTransferDetailRequest: stru.SPDBankTransferDetailRequest{
				AcctNo:        accountNo,
				DateBeginDate: beginDate,
				DateEndDate:   endDate,
				BeginNumber:   beginNumber,
				QueryNumber:   pageSize,
			},
		}
		signature := spdSign(ctx, body, sighHost)
		sign := &stru.SPDRequestSignBody{
			Sign: signature,
		}
		// 拿到sign 作为入参再次调用接口
		head := stru.NewSPDHeadData(enum.SPDBankListTransferDetailServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
		request := &stru.SPDRequestData{
			Head: *head,
			Body: sign,
		}
		// 返回的sign对象
		var responseSignData stru.SPDResponseSignData
		err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
			zap.L().Info(fmt.Sprintf("浦发服务: %s 查询交易明细验签异常code: %s msg:%s\n", enum.SPDBankListTransferDetailServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
			break
		}
		// 最终返回的对象
		var bankTransferDetailResponseData stru.SPDBankTransferDetailResponseData
		// 设置签名请求后的返回信息
		bankTransferDetailResponseData.ReturnCode = responseSignData.Body.ReturnCode
		bankTransferDetailResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
		// 拿到返回的sign, 在解析作为返参
		responseSign := responseSignData.Body.Sign
		if responseSign == "" {
			break
		}
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &bankTransferDetailResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if bankTransferDetailResponseData.Items != nil && len(bankTransferDetailResponseData.Items) > 0 {
			result = append(result, bankTransferDetailResponseData.Items...)
		} else {
			break
		}
		pageNum += 1
	}

	return result, nil
}

// RequestTransactionDetailElectronicReceipt
//  @Description:			电子回单下载申请
//  @receiver s
//  @param ctx
//  @param backhostGyno
//  @param subpoenaSeqNo
//  @param host
//  @param sighHost
//  @param bankCustomerId
//  @param bankUserId
//  @return []byte
//  @return error
//
func (s *spdBankSDK) RequestTransactionDetailElectronicReceipt(ctx context.Context, accountNo, beginDate, endDate, backhostGyno, subpoenaSeqNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferDetailElectronicReceiptResponseItem, error) {
	var result stru.SPDBankTransferDetailElectronicReceiptResponseItem
	req := stru.SPDBankTransferDetailElectronicReceiptRequest{
		BillDownloadChanel: "1",
		AcctNo:             accountNo,
		SingleOrBatchFlag:  "0",
		BackhostGyno:       backhostGyno,
		SubpoenaSeqNo:      subpoenaSeqNo,
		BeginDate:          beginDate,
		EndDate:            endDate,
		BeginNumber:        1,
		QueryNumber:        1,
	}
	// 根据body获取sign
	body := &stru.SPDBankTransferDetailElectronicReceiptRequestBody{
		SPDBankTransferDetailElectronicReceiptRequest: req,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	serialNo, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankTransferDetailElectronicReceiptDownloadRequestServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	// 返回的sign对象
	var responseSignData stru.SPDResponseSignData
	err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
		zap.L().Info(fmt.Sprintf("浦发服务: %s 文件下载验签异常code: %s msg:%s\n", enum.SPDBankTransferDetailElectronicReceiptDownloadRequestServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
		return nil, nil
	}
	// 最终返回的对象
	var receiptResponseData stru.SPDBankTransferDetailElectronicReceiptResponseData
	// 设置签名请求后的返回信息
	receiptResponseData.ReturnCode = responseSignData.Body.ReturnCode
	receiptResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
	// 拿到返回的sign, 在解析作为返参
	responseSign := responseSignData.Body.Sign
	if responseSign == "" {
		return nil, nil
	}
	responseData := spdVerifySign(ctx, responseSign, sighHost)
	err = xml.Unmarshal([]byte(responseData), &receiptResponseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 取第一个
	if &receiptResponseData != nil && len(receiptResponseData.Items) > 0 {
		result = receiptResponseData.Items[0]
	}
	return &result, nil
}

// DownloadTransactionDetailElectronicReceipt
//  @Description: 			获取交易明细电子回单
//  @receiver s
//  @param ctx
//  @param backhostGyno		柜员流水号: 对应交易明细的 tellerJnlNo
//  @param subpoenaSeqNo	传票组内序号: 对应交易明细的 summonsNumber
//  @param host
//  @param sighHost
//  @param bankCustomerId
//  @param bankUserId
//  @return []byte
//  @return error
//
func (s *spdBankSDK) DownloadTransactionDetailElectronicReceipt(ctx context.Context, accountNo, beginDate, endDate, backhostGyno, subpoenaSeqNo, host, sighHost, fileHost, bankCustomerId, bankUserId string) ([]byte, error) {
	receipt, err := s.RequestTransactionDetailElectronicReceipt(ctx, accountNo, beginDate, endDate, backhostGyno, subpoenaSeqNo, host, sighHost, bankCustomerId, bankUserId)
	//银行建议发起请求1分钟后再开启回单下载申请,生成回单需要时间
	time.Sleep(time.Second * 30)
	// todo 测试写死
	//receipt, err := s.RequestTransactionDetailElectronicReceipt(ctx, accountNo, "20210101", "20210131", "999709310001", "1", host, sighHost, bankCustomerId, bankUserId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if receipt == nil {
		return nil, nil
	}
	//receipt.AcceptNo = "20221010000010520003"
	//accountNo = "952A9997220008092"
	req := stru.SPDBankTransferDetailElectronicReceiptDownloadRequest{
		AcctNo:           accountNo,
		FileDownloadFlag: "1",
		FileDownloadPar:  receipt.AcceptNo,
	}
	body := &stru.SPDBankTransferDetailElectronicReceiptDownloadRequestBody{
		SPDBankTransferDetailElectronicReceiptDownloadRequest: req,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	serialNo, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankTransferDetailElectronicReceiptDownloadServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	resultBytes, err := util.PostSPDXmlHttpBytesResult(ctx, fileHost, "application/xml;charset=UTF-8", request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	responseStr := string(resultBytes)
	fileBinaryStr := responseStr[strings.Index(responseStr, "</packet>")+9:]
	fileBinary := []byte(fileBinaryStr)
	return fileBinary, nil
}

func (s *spdBankSDK) QueryTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo, elecChequeNo string, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankTransferResultResponseData, error) {
	pageNum := 1
	pageSize := 20
	// 最终返回的对象
	var bankTransferResultResponseData stru.SPDBankTransferResultResponseData
	for {
		serialNo, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		beginNumber := (pageNum-1)*pageSize + 1
		// 根据body获取sign
		body := &stru.SPDBankTransferResultRequestBody{
			SPDBankTransferResultRequest: stru.SPDBankTransferResultRequest{
				AcctNo:            payAccount,
				PayeeAcctNo:       recAccount,
				BeginDate:         beginDate,
				EndDate:           endDate,
				AcceptNo:          orderFlowNo,
				ElecChequeNo:      elecChequeNo,
				QueryNumber:       pageSize,
				BeginNumber:       beginNumber,
				SingleOrBatchFlag: "0",
			},
		}
		signature := spdSign(ctx, body, sighHost)
		sign := &stru.SPDRequestSignBody{
			Sign: signature,
		}
		// 拿到sign 作为入参再次调用接口
		head := stru.NewSPDHeadData(enum.SPDBankTransferResultServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
		request := &stru.SPDRequestData{
			Head: *head,
			Body: sign,
		}
		// 返回的sign对象
		var responseSignData stru.SPDResponseSignData
		err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
			zap.L().Info(fmt.Sprintf("浦发服务: %s 查询转账结果验签异常code: %s msg:%s\n", enum.SPDBankTransferResultServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
			break
		}
		// 设置签名请求后的返回信息
		bankTransferResultResponseData.ReturnCode = responseSignData.Body.ReturnCode
		bankTransferResultResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
		// 拿到返回的sign, 在解析作为返参
		responseSign := responseSignData.Body.Sign
		if responseSign == "" {
			break
		}
		var tempBankTransferResultResponseData stru.SPDBankTransferResultResponseData
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &tempBankTransferResultResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if tempBankTransferResultResponseData.Items != nil && len(tempBankTransferResultResponseData.Items) > 0 {
			bankTransferResultResponseData.Items = append(bankTransferResultResponseData.Items, tempBankTransferResultResponseData.Items...)
		} else {
			break
		}
		pageNum += 1
	}
	return &bankTransferResultResponseData, nil
}

func (s *spdBankSDK) QueryAccountBalance(ctx context.Context, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankQueryAccountBalanceResponseData, error) {
	// 根据body获取sign
	var requestItems []stru.SPDBankQueryAccountBalanceRequestItem
	requestItem := stru.SPDBankQueryAccountBalanceRequestItem{
		AcctNo: accountNo,
	}
	requestItems = append(requestItems, requestItem)

	body := &stru.SPDBankQueryAccountBalanceRequestBody{
		Items: requestItems,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	serialNo, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankQueryAccountBalanceServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	// 返回的sign对象
	var responseSignData stru.SPDResponseSignData
	err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData,
		request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
		zap.L().Info(fmt.Sprintf("浦发服务: %s 查询账户余额验签异常code: %s msg:%s\n", enum.SPDBankQueryAccountBalanceServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
		return nil, nil
	}
	// 最终返回的对象
	var bankQueryAccountBalanceResponseData stru.SPDBankQueryAccountBalanceResponseData
	// 设置签名请求后的返回信息
	bankQueryAccountBalanceResponseData.ReturnCode = responseSignData.Body.ReturnCode
	bankQueryAccountBalanceResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 拿到返回的sign, 在解析作为返参
	responseSign := responseSignData.Body.Sign
	if responseSign != "" {
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &bankQueryAccountBalanceResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
	}
	return &bankQueryAccountBalanceResponseData, nil
}

// CreateVirtualAccount
//  @Description: 创建虚账户
//  @receiver s
//  @param ctx
//  @param items
//  @param accountNo
//  @param host
//  @param sighHost
//  @param bankCustomerId
//  @param bankUserId
//  @return *stru.SPDCreateVirtualAccountResponseData
//  @return error
//
func (s *spdBankSDK) CreateVirtualAccount(ctx context.Context, items []stru.SPDCreateVirtualAccountRequestItem, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDCreateVirtualAccountResponseData, error) {
	// 根据body获取sign
	body := &stru.SPDCreateVirtualAccountRequestBody{
		AcctNo: accountNo,
		Items:  items,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	serialNo, err := util.SonyflakeID()
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankCreateVirtualAccountServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	// 返回的sign对象
	var responseSignData stru.SPDResponseSignData
	err = util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if responseSignData.Body.ReturnCode != "" && responseSignData.Body.ReturnMsg != "" {
		zap.L().Info(fmt.Sprintf("浦发服务: %s 创建虚账户验签异常code: %s msg:%s\n", enum.SPDBankCreateVirtualAccountServiceCode, responseSignData.Body.ReturnCode, responseSignData.Body.ReturnMsg))
		return nil, nil
	}
	// 最终返回的对象
	var createVirtualAccountResponseData stru.SPDCreateVirtualAccountResponseData
	// 拿到返回的sign, 在解析作为返参
	responseSign := responseSignData.Body.Sign
	if responseSign != "" {
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &createVirtualAccountResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
	}
	return &createVirtualAccountResponseData, nil
}

// QueryVirtualAccountBalance
//  @Description: 查询实账户下虚账户的即时金额
//  @receiver s
//  @param ctx
//  @param accountNo
//  @param virtualAccountNo
//  @param offset
//  @param host
//  @param sighHost
//  @param bankCustomerId
//  @param bankUserId
//  @return []stru.SPDQueryVirtualAccountBalanceResponseItem
//  @return error
//
func (s *spdBankSDK) QueryVirtualAccountBalance(ctx context.Context, accountNo string, virtualAccountNo []string, offset int, host, sighHost, bankCustomerId, bankUserId string) ([]stru.SPDQueryVirtualAccountBalanceResponseItem, error) {
	pageNum := 1
	pageSize := 20
	// 最终返回的对象
	var result []stru.SPDQueryVirtualAccountBalanceResponseItem
	// virtualAccountNo to map
	virtualAccountNoMap := make(map[string]string, len(virtualAccountNo))
	for _, v := range virtualAccountNo {
		virtualAccountNoMap[v] = v
	}
	for {
		serialNo, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		beginNumber := (pageNum-1)*pageSize + 1
		if offset > 0 {
			beginNumber = offset
		}
		// 根据body获取sign
		body := &stru.SPDQueryVirtualAccountBalanceRequestBody{
			AcctNo:      accountNo,
			BeginNumber: beginNumber,
			QueryNumber: pageSize,
		}
		signature := spdSign(ctx, body, sighHost)
		sign := &stru.SPDRequestSignBody{
			Sign: signature,
		}
		// 拿到sign 作为入参再次调用接口
		head := stru.NewSPDHeadData(enum.SPDBankQueryVirtualAccountBalanceServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
		request := &stru.SPDRequestData{
			Head: *head,
			Body: sign,
		}
		// 返回的sign对象
		var responseSignData stru.SPDResponseSignData
		err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		// 拿到返回的sign, 在解析作为返参
		responseSign := responseSignData.Body.Sign
		if responseSign == "" {
			break
		}
		var tempQueryVirtualAccountBalanceResponseData stru.SPDQueryVirtualAccountBalanceResponseData
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &tempQueryVirtualAccountBalanceResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if tempQueryVirtualAccountBalanceResponseData.Items != nil && len(tempQueryVirtualAccountBalanceResponseData.Items) > 0 {
			for _, v := range tempQueryVirtualAccountBalanceResponseData.Items {
				// 查询到了虚账号
				if _, exists := virtualAccountNoMap[v.VirtualAcctNo]; exists == true {
					result = append(result, v)
				}
			}
		} else {
			break
		}
		if offset > 0 {
			offset = offset + pageSize
		} else {
			pageNum += 1
		}
	}
	return result, nil
}

// BankVirtualAccountTransfer
//  @Description: 虚账户转账
//  @receiver s
//  @param ctx
//  @param serialNo
//  @param req
//  @param host
//  @param sighHost
//  @param bankCustomerId
//  @param bankUserId
//  @return *stru.SPDBankTransferResponseData
//  @return error
//
func (s *spdBankSDK) BankVirtualAccountTransfer(ctx context.Context, serialNo string, req stru.SPDBankVirtualAccountTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.SPDBankVirtualAccountTransferResponseData, error) {
	// 根据body获取sign
	body := &stru.SPDBankVirtualAccountTransferRequestBody{
		SPDBankVirtualAccountTransferRequest: req,
	}
	signature := spdSign(ctx, body, sighHost)
	sign := &stru.SPDRequestSignBody{
		Sign: signature,
	}
	// 拿到sign 作为入参再次调用接口
	head := stru.NewSPDHeadData(enum.SPDBankVirtualAccountTransferServiceCode, enum.SPDNeedSign, bankCustomerId, "", serialNo)
	request := &stru.SPDRequestData{
		Head: *head,
		Body: sign,
	}
	// 返回的sign对象
	var responseSignData stru.SPDResponseSignData
	err := util.PostSPDXmlHttpResult(ctx, host, "application/xml;charset=UTF-8", &responseSignData, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 最终返回的对象
	var bankTransferResponseData stru.SPDBankVirtualAccountTransferResponseData
	// 设置签名请求后的返回信息
	bankTransferResponseData.ReturnCode = responseSignData.Body.ReturnCode
	bankTransferResponseData.ReturnMsg = responseSignData.Body.ReturnMsg
	// 拿到返回的sign, 在解析作为返参
	responseSign := responseSignData.Body.Sign
	if responseSign != "" {
		responseData := spdVerifySign(ctx, responseSign, sighHost)
		err = xml.Unmarshal([]byte(responseData), &bankTransferResponseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
	}
	return &bankTransferResponseData, nil
}
