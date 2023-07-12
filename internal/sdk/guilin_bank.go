package sdk

import (
	"context"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-common/enum"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"strings"
)

type GuilinBankSDK interface {
	IntrabankTransfer(ctx context.Context, serialNo string, req stru.IntrabankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error)
	OutofbankTransfer(ctx context.Context, serialNo string, req stru.OutofbankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error)
	LedgerManagementTransfer(ctx context.Context, serialNo string, req stru.OutofbankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error)

	ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.TransactionDetailResponseItem, error)
	GetTransactionDetailElectronicReceipt(ctx context.Context, orderFlowNo string, host, sighHost, bankCustomerId, bankUserId string) ([]byte, error)
	QueryTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo string, isIntrabankTransfer bool, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.TransferResultResponseItem, error)
	QueryAccountBalance(ctx context.Context, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.QueryAccountBalanceResponse, error)
	UploadBatchTransferPayInfo(ctx context.Context, req stru.UploadBatchTransferPayInfoRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferPayInfoResponse, error)
	UploadBatchTransferRecInfo(ctx context.Context, batchNo string, req []stru.UploadBatchTransferRecInfoList, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferRecInfoResponse, error)
	QueryBatchTransferResult(ctx context.Context, req stru.QueryBatchTransferResultRequestBodyData, host, sighHost, bankCustomerId, bankUserId string) ([]stru.QueryBatchTransferResultResponseRow, error)

	QueryLedgerManagementTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo string, isIntrabankTransfer bool, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.LedgerManagementTransferResultResponseItem, error)
}

type guilinBankSDK struct{}

func sign(ctx context.Context, head stru.HeadData, body interface{}, sighHost string) string {
	result, err := util.PostXmlHttpStringResult(ctx, sighHost,
		"INFOSEC_SIGN/1.0", false, head, body)
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}
	return result[strings.Index(result, "<sign>")+6 : strings.Index(result, "</sign>")]
}

func (s *guilinBankSDK) IntrabankTransfer(ctx context.Context, serialNo string, req stru.IntrabankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error) {
	head := stru.NewHeadData(enum.GuilinBankIntrabankTransferServiceCode, serialNo, bankCustomerId, bankUserId)
	if req.PayRem == "" {
		req.PayRem = enum.DefaultRecommend
	}
	if req.ConfirmCode == "" {
		req.ConfirmCode = "9999"
	}
	body := stru.IntrabankTransferRequestBody{
		IntrabankTransferRequest: req,
		BusinessCode:             enum.GuilinBankIntrabankTransferBusinessCode,
	}
	var bankTransferResponse stru.BankTransferResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&bankTransferResponse, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &bankTransferResponse, handler.HandleError(err)
}

func (s *guilinBankSDK) OutofbankTransfer(ctx context.Context, serialNo string, req stru.OutofbankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error) {
	head := stru.NewHeadData(enum.GuilinBankOutofbankTransferServiceCode, serialNo, bankCustomerId, bankUserId)
	if req.PayRem == "" {
		req.PayRem = enum.DefaultRecommend
	}
	if req.ConfirmCode == "" {
		req.ConfirmCode = "9999"
	}
	body := stru.OutofbankTransferRequestBody{
		OutofbankTransferRequest: req,
		BusinessCode:             enum.GuilinBankOutofbankTransferBusinessCode,
	}
	var bankTransferResponse stru.BankTransferResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&bankTransferResponse, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &bankTransferResponse, handler.HandleError(err)
}

func (s *guilinBankSDK) LedgerManagementTransfer(ctx context.Context, serialNo string, req stru.OutofbankTransferRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.BankTransferResponse, error) {
	head := stru.NewHeadData(enum.GuilinBankLedgerManagementTransferServiceCode, serialNo, bankCustomerId, bankUserId)
	if req.PayRem == "" {
		req.PayRem = enum.DefaultRecommend
	}
	if req.ConfirmCode == "" {
		req.ConfirmCode = "9999"
	}
	body := stru.OutofbankTransferRequestBody{
		OutofbankTransferRequest: req,
		BusinessCode:             enum.GuilinBankLedgerManagementTransferBusinessCode,
	}
	var bankTransferResponse stru.BankTransferResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&bankTransferResponse, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &bankTransferResponse, handler.HandleError(err)
}

func (s *guilinBankSDK) ListTransactionDetail(ctx context.Context, accountNo, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.TransactionDetailResponseItem, error) {
	turnPageNum := 1
	turnPageShowNum := 100
	var result []stru.TransactionDetailResponseItem
	for {
		id, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		head := stru.NewHeadData(enum.GuilinBankTransactionDetailServiceCode, id, bankCustomerId, bankUserId)
		body := stru.TransactionDetailRequestBody{
			AccountNo:       strings.ReplaceAll(accountNo, " ", ""),
			BeginDate:       beginDate,
			EndDate:         endDate,
			TurnPageNum:     turnPageNum,
			TurnPageShowNum: turnPageShowNum,
			QueryFlag:       "01",
		}
		var response stru.TransactionDetailResponse
		err := util.PostXmlHttpResult(ctx,
			fmt.Sprintf("%s/corporbank/httpAccess", host),
			"application/xml;charset=UTF-8",
			&response, stru.RequestData{
				Head: *head,
				Body: body,
				Sign: sign(ctx, *head, body, sighHost),
			})
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if response.Body.Items != nil && len(response.Body.Items) > 0 {
			result = append(result, response.Body.Items...)
		} else {
			break
		}
		turnPageNum += 1
	}
	return result, nil
}

func (s *guilinBankSDK) GetTransactionDetailElectronicReceipt(ctx context.Context, orderFlowNo string, host, sighHost, bankCustomerId, bankUserId string) ([]byte, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	head := stru.NewHeadData(enum.GuilinBankTransactionDetailElectronicReceiptServiceCode, id, bankCustomerId, bankUserId)
	body := stru.TransactionDetailElectronicReceiptRequestBody{
		OrderFlowNo: orderFlowNo,
		FileType:    "pdf",
		IsDownload:  true,
	}
	result, err := util.PostXmlHttpBytesResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8", true,
		stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return result, nil
}

func (s *guilinBankSDK) QueryTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo string, isIntrabankTransfer bool, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.TransferResultResponseItem, error) {

	turnPageBeginPos := 1
	turnPageShowNum := 100
	var result []stru.TransferResultResponseItem
	businessCode := "020102"
	if isIntrabankTransfer {
		businessCode = "020100"
	}

	for {
		id, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		head := stru.NewHeadData(enum.GuilinBankTransferResultServiceCode, id, bankCustomerId, bankUserId)
		body := stru.BatchQueryTransferResultRequestBody{
			SearchPayAccount: payAccount,
			SearchRecAccount: recAccount,
			OrderFlowNo:      orderFlowNo,
			BeginDate:        beginDate,
			EndDate:          endDate,
			TurnPageBeginPos: turnPageBeginPos,
			TurnPageShowNum:  turnPageShowNum,
			BusinessCode:     businessCode,
		}
		var response stru.TransferListResponse
		err := util.PostXmlHttpResult(ctx,
			fmt.Sprintf("%s/corporbank/httpAccess", host),
			"application/xml;charset=UTF-8",
			&response, stru.RequestData{
				Head: *head,
				Body: body,
				Sign: sign(ctx, *head, body, sighHost),
			})
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if response.Body.Items != nil && len(response.Body.Items) > 0 {
			result = append(result, response.Body.Items...)
		} else {
			break
		}
		turnPageBeginPos += 1
	}
	return result, nil

}

func (s *guilinBankSDK) QueryAccountBalance(ctx context.Context, accountNo, host, sighHost, bankCustomerId, bankUserId string) (*stru.QueryAccountBalanceResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	head := stru.NewHeadData(enum.GuilinBankQueryAccountBalanceServiceCode, id, bankCustomerId, bankUserId)
	body := stru.QueryAccountBalanceRequestBody{
		QueryAccountBalanceRequest: stru.QueryAccountBalanceRequest{
			AccountNo: accountNo,
		},
	}
	var queryAccountBalance stru.QueryAccountBalanceResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&queryAccountBalance, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &queryAccountBalance, handler.HandleError(err)
}

func (s *guilinBankSDK) UploadBatchTransferPayInfo(ctx context.Context, req stru.UploadBatchTransferPayInfoRequest, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferPayInfoResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	head := stru.NewHeadData(enum.GuilinBankBatchTransferServiceCode, id, bankCustomerId, bankUserId)
	body := stru.UploadBatchTransferPayInfoRequestBody{
		UploadBatchTransferPayInfoRequest: req,
	}
	var uploadBatchTransferPayInfoResponse stru.UploadBatchTransferPayInfoResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&uploadBatchTransferPayInfoResponse, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &uploadBatchTransferPayInfoResponse, nil
}

func (s *guilinBankSDK) UploadBatchTransferRecInfo(ctx context.Context, batchNo string, req []stru.UploadBatchTransferRecInfoList, host, sighHost, bankCustomerId, bankUserId string) (*stru.UploadBatchTransferRecInfoResponse, error) {
	id, sfErr := util.SonyflakeID()
	if sfErr != nil {
		return nil, handler.HandleError(sfErr)
	}
	head := stru.NewHeadData(enum.GuilinBankBatchTransferDetailServiceCode, id, bankCustomerId, bankUserId)
	body := stru.UploadBatchTransferRecInfoRequestBody{
		UploadBatchTransferRecInfoRequestBodyData: stru.UploadBatchTransferRecInfoRequestBodyData{
			BatchNo:      batchNo,
			AccountLists: req,
		},
	}
	var uploadBatchTransferRecInfoResponse stru.UploadBatchTransferRecInfoResponse
	err := util.PostXmlHttpResult(ctx,
		fmt.Sprintf("%s/corporbank/httpAccess", host),
		"application/xml;charset=UTF-8",
		&uploadBatchTransferRecInfoResponse, stru.RequestData{
			Head: *head,
			Body: body,
			Sign: sign(ctx, *head, body, sighHost),
		})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &uploadBatchTransferRecInfoResponse, nil
}

func (s *guilinBankSDK) QueryBatchTransferResult(ctx context.Context, req stru.QueryBatchTransferResultRequestBodyData, host, sighHost, bankCustomerId, bankUserId string) ([]stru.QueryBatchTransferResultResponseRow, error) {
	turnPageBeginPos := 1
	turnPageShowNum := 100
	var result []stru.QueryBatchTransferResultResponseRow
	businessCode := "020104"

	for {
		id, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		head := stru.NewHeadData(enum.GuilinBankBatchTransferResultServiceCode, id, bankCustomerId, bankUserId)
		body := stru.QueryBatchTransferResultRequestBody{
			SearchPayAccount: req.SearchPayAccount,
			BeginDate:        req.BeginDate,
			EndDate:          req.EndDate,
			BatchNo:          req.BatchNo,
			OrderFlowNo:      req.OrderFlowNo,
			TurnPageBeginPos: turnPageBeginPos,
			TurnPageShowNum:  turnPageShowNum,
			OrderState:       req.OrderState,
			BusinessCode:     businessCode,
		}
		var response stru.QueryBatchTransferResultResponse
		err := util.PostXmlHttpResult(ctx,
			fmt.Sprintf("%s/corporbank/httpAccess", host),
			"application/xml;charset=UTF-8",
			&response, stru.RequestData{
				Head: *head,
				Body: body,
				Sign: sign(ctx, *head, body, sighHost),
			})
		if err != nil {
			return nil, err
		}
		if response.Body.Rows != nil && len(response.Body.Rows) > 0 {
			result = append(result, response.Body.Rows...)
		} else {
			break
		}
		turnPageBeginPos += 1
	}
	return result, nil
}

func (s *guilinBankSDK) QueryLedgerManagementTransferResult(ctx context.Context, payAccount string, recAccount string, orderFlowNo string, isIntrabankTransfer bool, beginDate, endDate string, host, sighHost, bankCustomerId, bankUserId string) ([]stru.LedgerManagementTransferResultResponseItem, error) {

	turnPageBeginPos := 1
	turnPageShowNum := 100
	var result []stru.LedgerManagementTransferResultResponseItem
	recBankType := "1"
	if isIntrabankTransfer {
		recBankType = "0"
	}

	for {
		id, sfErr := util.SonyflakeID()
		if sfErr != nil {
			return nil, handler.HandleError(sfErr)
		}
		head := stru.NewHeadData(enum.GuilinBankTransferPushResultServiceCode, id, bankCustomerId, bankUserId)
		body := stru.BatchQueryLedgerManagementTransferResultRequestBody{
			PayAccount:       payAccount,
			SearchRecAccount: recAccount,
			OrderFlowNo:      orderFlowNo,
			BeginDate:        beginDate,
			EndDate:          endDate,
			TurnPageBeginPos: turnPageBeginPos,
			TurnPageShowNum:  turnPageShowNum,
			RecBankType:      recBankType,
		}
		var response stru.LedgerManagementTransferListResponse
		err := util.PostXmlHttpResult(ctx,
			fmt.Sprintf("%s/corporbank/httpAccess", host),
			"application/xml;charset=UTF-8",
			&response, stru.RequestData{
				Head: *head,
				Body: body,
				Sign: sign(ctx, *head, body, sighHost),
			})
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if response.Body.Items != nil && len(response.Body.Items) > 0 {
			result = append(result, response.Body.Items...)
		} else {
			break
		}
		turnPageBeginPos += 1
	}
	return result, nil

}
