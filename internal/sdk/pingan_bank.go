package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type PinganBankSDK interface {
	BankTransfer(ctx context.Context, req stru.PingAnBankTransferRequest, zuId string) (*stru.PingAnBankTransferResponse, error)
	BankVirtualTransfer(ctx context.Context, req stru.PingAnBankTransferRequest) (*stru.PingAnBankTransferResponse, error)
	SubAcctBalanceAdjust(ctx context.Context, req stru.PinganSubAcctBalanceAdjustRuest) (*stru.PinganSubAcctBalanceAdjustResponse, error)
	QueryAccountBalance(ctx context.Context, accountNo string, host string, mrchCode string, zuId string) (*stru.PinganBankCorAcctBalanceQueryResponse, error)
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, zuId string) ([]stru.PinganHistoryTransactionDetailsItem, error)
	ListVirtualTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.PinganHistoryTransactionDetailsItem, error)
	QueryTransferResult(ctx context.Context, ThirdVoucher, zuId, accountNo string) (*stru.PinganSignleTransferQueryResponse, error)
	CreateVirtualAccount(ctx context.Context, req stru.PinganCreateVirtualAccountRequest) (*stru.PinganCreateVirtualAccountResponse, error)
	QueryVirtualAccount(ctx context.Context, req stru.PinganQueryVirtualAccountBalanceRequest) (*stru.PinganQueryVirtualAccountBalanceResponse, error)
	UploadTransactionDetailElectronic(ctx context.Context, req stru.PinganSameDayHistoryReceiptDataQueryRequest, zuId string) (string, error)
	UploadVirtualTransactionDetailElectronic(ctx context.Context, req stru.PinganSameDayHistoryReceiptDataQueryRequest) (string, error)

	UserAcctSignatureApply(ctx context.Context, accountNo string, accountName string, host string, mrchCode string) (*stru.PinganUserAcctSignatureApplyResponse, error)
	UserAcctSignatureQuery(ctx context.Context, accountNo string, accountName string, host string, BankCustomerId string, zuId string) (*stru.PinganUserAcctSignatureApplyResponse, error)
}

type pinganBankSDK struct{}

func (s *pinganBankSDK) ListVirtualTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.PinganHistoryTransactionDetailsItem, error) {
	pageNum := 1
	pageSize := 30
	var result []stru.PinganHistoryTransactionDetailsItem
	for {
		serialNo, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		pageNumStr := strconv.Itoa(pageNum)
		pageSizeStr := strconv.Itoa(pageSize)
		requestBody := &stru.PinganHistoryTransactionDetailsRequest{
			MrchCode:   config.GetString(enum.PinganIntelligenceMrchCode, ""),
			CnsmrSeqNo: serialNo,
			AcctNo:     account,
			CcyCode:    "RMB",
			BeginDate:  beginDate,
			EndDate:    endDate,
			PageNo:     pageNumStr,
			PageSize:   pageSizeStr,
		}
		requestBodyJson, _ := json.Marshal(requestBody)
		request := &stru.PingAnBankRequestBody{
			RequestBody:   string(requestBodyJson),
			InterfaceName: "/bedl/InquiryAccountDayHistoryTransactionDetails",
			InterfaceType: 2,
		}
		// 返回的对象
		var responseData stru.PinganHistoryTransactionDetailsResponse
		err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkUrl, ""), request, &responseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		err = handelPinganResult(responseData.PinganErrorResult)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		result = append(result, responseData.List...)
		if responseData.EndFlag == "N" {
			pageNum += 1
		} else {
			break
		}

	}
	return result, nil
}

// BankTransfer
//
//	@Description:			转账
//	@receiver s
//	@param ctx
//	@param serialNo 		对应企业发送报文的报文号 当天内唯一
//	@param req				转账请求体
//	@return *
//	@return error
func (s *pinganBankSDK) BankTransfer(ctx context.Context, req stru.PingAnBankTransferRequest, zuId string) (*stru.PingAnBankTransferResponse, error) {
	requestBodyJson, _ := json.Marshal(req)
	var request *stru.PingAnBankRequestBody
	if req.PaymentModeType == "" {
		if config.GetString(enum.PinganPlatformAccount, "") == req.OutAcctNo {
			request = &stru.PingAnBankRequestBody{
				RequestBody:   string(requestBodyJson),
				InterfaceName: "/bedl/CorSingleTransfer",
			}
		} else {
			request = &stru.PingAnBankRequestBody{
				RequestBody:   string(requestBodyJson),
				InterfaceName: "/bedl/SignleTransferPrepaidExpenses",
				ZuId:          zuId,
			}
		}
	} else if req.PaymentModeType == "1" {
		request = &stru.PingAnBankRequestBody{
			RequestBody:   string(requestBodyJson),
			InterfaceName: "/bedl/CorSingleTransfer",
		}
	} else {
		request = &stru.PingAnBankRequestBody{
			RequestBody:   string(requestBodyJson),
			InterfaceName: "/bedl/SignleTransferPrepaidExpenses",
			ZuId:          zuId,
		}
	}
	if req.OutAcctNo == config.GetString(enum.PinganIntelligenceAccountNo, "") {
		request.InterfaceType = 2
	}
	var responseData stru.PingAnBankTransferResponse
	err := util.PostHttpResult(ctx,
		config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)

	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}
func (s *pinganBankSDK) BankVirtualTransfer(ctx context.Context, req stru.PingAnBankTransferRequest) (*stru.PingAnBankTransferResponse, error) {
	requestBodyJson, _ := json.Marshal(req)
	var request *stru.PingAnBankRequestBody

	request = &stru.PingAnBankRequestBody{
		RequestBody:   string(requestBodyJson),
		InterfaceName: "/bedl/CorSingleTransfer",
		InterfaceType: 2,
	}

	var responseData stru.PingAnBankTransferResponse
	err := util.PostHttpResult(ctx,
		config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)

	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}
func (s *pinganBankSDK) SubAcctBalanceAdjust(ctx context.Context, req stru.PinganSubAcctBalanceAdjustRuest) (*stru.PinganSubAcctBalanceAdjustResponse, error) {
	requestBodyJson, _ := json.Marshal(req)
	var request *stru.PingAnBankRequestBody

	request = &stru.PingAnBankRequestBody{
		RequestBody:   string(requestBodyJson),
		InterfaceName: "/bedl/SubAcctBalanceAdjust",
		InterfaceType: 2,
	}

	var responseData stru.PinganSubAcctBalanceAdjustResponse
	err := util.PostHttpResult(ctx,
		config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}

// 查询账户详情
func (s *pinganBankSDK) QueryAccountBalance(ctx context.Context, accountNo string, host string, mrchCode string, zuId string) (*stru.PinganBankCorAcctBalanceQueryResponse, error) {
	seqNoId, err := util.SonyflakeID()
	if err != nil {
		return nil, handler.HandleError(err)
	}

	requestBody := stru.PinganBankCorAcctBalanceQueryRequest{
		MrchCode:   mrchCode,
		CnsmrSeqNo: seqNoId,
		Account:    accountNo,
	}
	var request *stru.PingAnBankRequestBody
	requestBodyJson, _ := json.Marshal(requestBody)
	if config.GetString(enum.PinganPlatformAccount, "") == accountNo {
		request = &stru.PingAnBankRequestBody{
			RequestBody:   string(requestBodyJson),
			InterfaceName: "/bedl/CorAcctBalanceQuery",
		}
	} else {
		request = &stru.PingAnBankRequestBody{
			RequestBody:   string(requestBodyJson),
			InterfaceName: "/bedl/CorAcctBalanceQuery",
			ZuId:          zuId,
		}
	}

	var responseData stru.PinganBankCorAcctBalanceQueryResponse
	err = util.PostHttpResult(ctx, host, &request, &responseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}

// 查询历史交易明细
func (s *pinganBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, zuId string) ([]stru.PinganHistoryTransactionDetailsItem, error) {
	pageNum := 1
	pageSize := 30
	var result []stru.PinganHistoryTransactionDetailsItem
	for {
		serialNo, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		pageNumStr := strconv.Itoa(pageNum)
		pageSizeStr := strconv.Itoa(pageSize)
		requestBody := &stru.PinganHistoryTransactionDetailsRequest{
			MrchCode:   config.GetString("bankConfig.pingAn.accountKeeper.mrchCode", ""),
			CnsmrSeqNo: serialNo,
			AcctNo:     account,
			CcyCode:    "RMB",
			BeginDate:  beginDate,
			EndDate:    endDate,
			PageNo:     pageNumStr,
			PageSize:   pageSizeStr,
		}
		requestBodyJson, _ := json.Marshal(requestBody)
		var request *stru.PingAnBankRequestBody
		if config.GetString(enum.PinganPlatformAccount, "") == account {
			request = &stru.PingAnBankRequestBody{
				RequestBody:   string(requestBodyJson),
				InterfaceName: "/bedl/InquiryAccountDayHistoryTransactionDetails",
			}
		} else {
			request = &stru.PingAnBankRequestBody{
				RequestBody:   string(requestBodyJson),
				InterfaceName: "/bedl/InquiryAccountDayHistoryTransactionDetails",
				ZuId:          zuId,
			}
		}
		// 返回的对象
		var responseData stru.PinganHistoryTransactionDetailsResponse
		err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		err = handelPinganResult(responseData.PinganErrorResult)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		result = append(result, responseData.List...)
		if responseData.EndFlag == "N" {
			pageNum += 1
		} else {
			break
		}

	}
	return result, nil
}
func (s *pinganBankSDK) QueryTransferResult(ctx context.Context, ThirdVoucher, zuId, accountNo string) (*stru.PinganSignleTransferQueryResponse, error) {
	CnsmrSeqNo, _ := util.SonyflakeID()
	mrchCode := ""
	interfaceName := ""
	interfaceType := 0

	if strings.HasPrefix(ThirdVoucher, enum.PinganFlexPrefix) {
		//灵活用工的转账
		mrchCode = config.GetString(enum.PinganIntelligenceMrchCode, "")
		interfaceName = "/bedl/CorSingleTransferQuery"
		interfaceType = 2
		zuId = ""
	} else {
		//银企转账
		mrchCode = config.GetString(enum.PinganIntelligenceMrchCode, "")
		if config.GetString(enum.PinganPlatformAccount, "") != accountNo {
			//账管+模式
			interfaceName = "/bedl/SignleTransferPrepaidExpensesProgressQuery"
		} else {
			//网银模式
			interfaceName = "/bedl/CorSingleTransferQuery"
		}
	}
	requestBody := &stru.PinganSignleTransferQueryRequest{
		MrchCode:         mrchCode,
		CnsmrSeqNo:       CnsmrSeqNo,
		OrigThirdVoucher: ThirdVoucher,
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(requestBodyJson),
		InterfaceName: interfaceName,
		InterfaceType: interfaceType,
		ZuId:          zuId,
	}

	var responseData stru.PinganSignleTransferQueryResponse
	err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}
func (s *pinganBankSDK) CreateVirtualAccount(ctx context.Context, req stru.PinganCreateVirtualAccountRequest) (*stru.PinganCreateVirtualAccountResponse, error) {
	reqJson, _ := json.Marshal(req)
	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(reqJson),
		InterfaceName: "/bedl/SubAcctMaintenance",
		InterfaceType: 2,
	}
	var responseData stru.PinganCreateVirtualAccountResponse

	err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}

func (s *pinganBankSDK) QueryVirtualAccount(ctx context.Context, req stru.PinganQueryVirtualAccountBalanceRequest) (*stru.PinganQueryVirtualAccountBalanceResponse, error) {
	reqJson, _ := json.Marshal(req)
	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(reqJson),
		InterfaceName: "/bedl/SubAccountBalanceQuery",
		InterfaceType: 2,
	}
	var responseData stru.PinganQueryVirtualAccountBalanceResponse
	err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkUrl, ""), &request, &responseData)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseData.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseData, nil
}
func (s *pinganBankSDK) UploadTransactionDetailElectronic(ctx context.Context, req stru.PinganSameDayHistoryReceiptDataQueryRequest, zuId string) (string, error) {
	//平安的文件下载功能,对于当天的下载不需要填写日期,
	todayTime := util.FormatTimeyyyyMMdd(time.Now())
	if req.AccountBeginDate == todayTime || req.AccountEndDate == todayTime {
		req.AccountBeginDate = ""
		req.AccountEndDate = ""
	}
	reqJson, _ := json.Marshal(req)
	var request *stru.PingAnBankRequestBody
	if config.GetString(enum.PinganPlatformAccount, "") == req.OutAccNo {
		request = &stru.PingAnBankRequestBody{
			RequestBody: string(reqJson),
		}
	} else {
		request = &stru.PingAnBankRequestBody{
			RequestBody: string(reqJson),
			ZuId:        zuId,
		}
	}
	var fileResult stru.PinganFileResult
	zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证请求参数%+v", request))
	err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkFileUrl, ""), &request, &fileResult)
	if err != nil {
		return "", err
	}
	return fileResult.FilePath, nil
}

func (s *pinganBankSDK) UploadVirtualTransactionDetailElectronic(ctx context.Context, req stru.PinganSameDayHistoryReceiptDataQueryRequest) (string, error) {
	//平安的文件下载功能,对于当天的下载不需要填写日期,
	todayTime := util.FormatTimeyyyyMMdd(time.Now())
	if req.AccountBeginDate == todayTime || req.AccountEndDate == todayTime {
		req.AccountBeginDate = ""
		req.AccountEndDate = ""
	}
	reqJson, _ := json.Marshal(req)

	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(reqJson),
		InterfaceType: 2,
	}
	var fileResult stru.PinganFileResult
	zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载平安智能清分电子凭证请求参数%+v", request))
	err := util.PostHttpResult(ctx, config.GetString(enum.PinganJsdkFileUrl, ""), &request, &fileResult)
	if err != nil {
		return "", err
	}

	return fileResult.FilePath, nil
}

func (s *pinganBankSDK) UserAcctSignatureApply(ctx context.Context, accountNo string, accountName string, host string, mrchCode string) (*stru.PinganUserAcctSignatureApplyResponse, error) {
	serialNo, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	requestBody := &stru.PinganUserAcctSignatureApplyRequest{
		MrchCode:     mrchCode,
		CnsmrSeqNo:   serialNo,
		ThirdVoucher: fmt.Sprintf("%s-%s", serialNo, util.FormatTimeyyyyMMdd(time.Now())),
		OpFlag:       "A",
		ZuID:         fmt.Sprintf("ZuId-%s", serialNo),
		AccountNo:    accountNo,
		AccountName:  accountName,
		BsnStr:       config.GetString("bankConfig.pingAn.accountKeeper.BsnStr", "4001,4013,400409,401805,404710,4048,404711,C001,C00202,C00203"),
	}

	requestBodyJson, _ := json.Marshal(requestBody)

	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(requestBodyJson),
		InterfaceName: "/bedl/UserAcctSignatureApply",
	}

	var resultBody stru.PinganUserAcctSignatureApplyResponse

	err := util.PostHttpResult(ctx, host, &request, &resultBody)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(resultBody.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &resultBody, nil
}
func (s *pinganBankSDK) UserAcctSignatureQuery(ctx context.Context, accountNo string, accountName string, host string, BankCustomerId string, zuId string) (*stru.PinganUserAcctSignatureApplyResponse, error) {
	serialNo, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	requestBody := &stru.PinganUserAcctSignatureQueryRequest{
		MrchCode:   BankCustomerId,
		CnsmrSeqNo: serialNo,
		ZuID:       zuId,
		AccountNo:  accountNo,
	}

	requestBodyJson, _ := json.Marshal(requestBody)

	request := &stru.PingAnBankRequestBody{
		RequestBody:   string(requestBodyJson),
		InterfaceName: "/bedl/SignatureResultQuery",
	}

	var responseBody stru.PinganUserAcctSignatureApplyResponse

	err := util.PostHttpResult(ctx, host, &request, &responseBody)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	err = handelPinganResult(responseBody.PinganErrorResult)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &responseBody, nil
}
func handelPinganResult(result stru.PinganErrorResult) error {
	if result.Code != "" {
		return handler.HandleNewError(fmt.Sprintf("==平安服务请求异常%v", result))
	}
	return nil
}
