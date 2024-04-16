package sdk

import (
	"context"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
)

// 民生银行接口对接
type MinShengSDK interface {
	BankTransfer(ctx context.Context, req stru.MinShengTransferRequest) (*stru.MinShengBankTransferResponse, error)
	QueryTransferResult(ctx context.Context, acctNo, orgReqSeq string) (*stru.MinShengSingleTransferQueryResponse, error)
	ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, startIndex, endIndex int64) (stru.MinShengSingleTransferDetailResponse, error)
	GetTransactionDetailElectronicReceipt(ctx context.Context, orderFlowNo string) ([]byte, error)
	AuthRequest(ctx context.Context, entCifno, acctNo string) (*stru.MinShengAuthResponse, error)
	QueryAuthStatus(ctx context.Context, srcReqSeq string) (*stru.MinShengAuthStatusResponse, error)
	//ReNewAuthRequest(ctx context.Context, openId, authCode string)
}

type minShengSDK struct{}

func (s *minShengSDK) BankTransfer(ctx context.Context, req stru.MinShengTransferRequest) (*stru.MinShengBankTransferResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["acct_no"] = req.AcctNo
	busiParamMap["pay_type"] = "1" // 直接支付
	busiParamMap["is_cross"] = req.IsCross
	busiParamMap["currency"] = "CNY"
	busiParamMap["trans_amt"] = req.TransAmt
	busiParamMap["bank_route"] = req.BankRoute
	busiParamMap["bank_code"] = req.BankCode // 开户行号
	busiParamMap["bank_name"] = req.BankName // 开户行名
	// 请求民生接口方法名
	method := "settlement.transfer.ent_single_order"
	// 民生接口版本号
	version := "V1.0"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengBankTransferResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) QueryTransferResult(ctx context.Context, acctNo, orgReqSeq string) (*stru.MinShengSingleTransferQueryResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["acct_no"] = acctNo
	busiParamMap["org_req_seq"] = orgReqSeq
	// 请求民生接口方法名
	method := "settlement.transfer.ent_single_order_qry"
	// 民生接口版本号
	version := "V1.1"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengSingleTransferQueryResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, startIndex, endIndex int64) (*stru.MinShengSingleTransferDetailResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}

	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["acct_no"] = accountNo
	busiParamMap["date_from"] = beginDate
	busiParamMap["date_to"] = endDate
	busiParamMap["start_index"] = startIndex
	busiParamMap["end_index"] = endIndex
	// 请求民生接口方法名
	method := "account.transinfo.detail_query"
	// 民生接口版本号
	version := "V1.1"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengSingleTransferDetailResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) GetTransactionDetailElectronicReceipt(ctx context.Context, acctNo, transSeqNo, enterAcctDate string) (*stru.MinShengElectronicReceiptResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["qry_type"] = "1" // 查询模式 1-交易流水查询（交易明细查询）
	busiParamMap["acct_no"] = acctNo
	busiParamMap["trans_seq_no"] = transSeqNo       // 交易流水号
	busiParamMap["enter_acct_date"] = enterAcctDate // 明细入账日期
	//busiParamMap["file_type"] = "PDF" // 默认就是PDF
	// 请求民生接口方法名
	method := "account.transinfo.electnote_download_new"
	// 民生接口版本号
	version := "V1.0"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengElectronicReceiptResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) AuthRequest(ctx context.Context, entCifno, acctNo string) (*stru.MinShengAuthResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["auth_type"] = "1"      // 1-企业自主授权
	busiParamMap["acct_no"] = acctNo     // 授权账号
	busiParamMap["ent_cifno"] = entCifno // 企业识别码
	//busiParamMap["file_type"] = "PDF" // 默认就是PDF
	// 请求民生接口方法名
	method := "account.openauth.apply"
	// 民生接口版本号
	version := "V1.1"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengAuthResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) QueryAuthStatus(ctx context.Context, srcReqSeq string) (*stru.MinShengAuthStatusResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	// 业务参数
	busiParamMap := make(map[string]interface{})
	busiParamMap["req_seq"] = id
	busiParamMap["src_req_seq"] = srcReqSeq // 原请求流水号
	//busiParamMap["file_type"] = "PDF" // 默认就是PDF
	// 请求民生接口方法名
	method := "account.openauth.status_qry"
	// 民生接口版本号
	version := "V1.0"
	response, err := s.invokeMinSheng(ctx, method, version, busiParamMap)
	if err != nil {
		return nil, err
	}
	minShengResponse, _ := response.(stru.MinShengResponse)
	result, _ := minShengResponse.ResponseBusi.(stru.MinShengAuthStatusResponse)
	result.ReturnCode = minShengResponse.ReturnCode
	result.ReturnMsg = minShengResponse.ReturnMsg
	return &result, nil
}

func (s *minShengSDK) invokeMinSheng(ctx context.Context, method, version string, busiParamMap map[string]interface{}) (interface{}, error) {
	var result interface{}
	host := "http://localhpst:8888/bank-web/cmbc/request"
	zap.L().Info(fmt.Sprintf("调用民生接口请求，方法名：%s,版本：%s,参数：%+v", method, version, busiParamMap))
	requestParamMap := make(map[string]interface{})
	// 民生接口名称
	requestParamMap["method"] = method
	// 民生接口版本号
	requestParamMap["version"] = version
	// 民生接口业务参数
	requestParamMap["busiParamMap"] = busiParamMap
	err := util.PostHttpResult(ctx, host, requestParamMap, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
