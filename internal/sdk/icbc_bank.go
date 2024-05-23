package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type IcbcBankSDK interface {
	QueryAgreeNo(ctx context.Context, zuId, account string) (*stru.Agreement, error)
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, accCompNo string, agreeNo string) ([]stru.IcbcAccDetailItem, error)
	IcbcUserAcctSignatureApply(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (string, error)
	IcbcUserAcctSignatureAgreePush(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (*stru.IcbcSignConfirmResponse, error)
	IcbcUserAcctSignatureQuery(ctx context.Context, accountNo string, accCompNo string) (*stru.IcbcSignatureQueryResponse, error)
	IcbcReceiptNoQuery(ctx context.Context, accountNo, accCompNo, agreeNo, serialNo string) (*stru.IcbcReceiptNoQueryResponse, error)
}

type icbcBankSDK struct {
}

func (i *icbcBankSDK) IcbcReceiptNoQuery(ctx context.Context, accountNo, accCompNo, agreeNo, serialNo string) (*stru.IcbcReceiptNoQueryResponse, error) {
	request := stru.NewIcbcGlobalRequest()
	serl := []string{serialNo}
	cond := stru.IcbcReceiptNoQueryCond{
		SeqList: serl,
	}
	fseqNo, _ := util.SonyflakeID()
	request.BizContent = &stru.IcbcReceiptNoQueryRequest{
		FseqNo:    fseqNo,
		CorpNo:    config.GetString(enum.IcbcCorpNo, ""),
		AccCompNo: accCompNo,
		AgreeNo:   agreeNo,
		Account:   accountNo,
		CurrType:  "",
		QryType:   "2",
		QryCond:   cond,
	}
	var result stru.IcbcReceiptNoQueryResponse
	resultInterface, err := stru.ICBCPostHttpResult(enum.IcbcAdsReceiptAryURL, *request)
	if err != nil {
		return nil, err
	}
	jsonString, _ := json.Marshal(resultInterface)
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil, err
	}

	if result.RetCode != "9008100" {
		return nil, errors.New(result.RetMsg)
	}
	return &result, nil
}
func (i *icbcBankSDK) QueryAgreeNo(ctx context.Context, zuId, account string) (*stru.Agreement, error) {
	//IcbcAdsAgreementGryURL
	request := stru.NewIcbcGlobalRequest()
	inqwork := stru.Inqwork{
		BegNum: "0",
		FetNum: "10",
	}
	cond := stru.Cond{
		QryType:   "1",
		AccCompNo: zuId,
		Account:   account,
		CurrType:  "",
		AgrList:   nil,
	}
	request.BizContent = &stru.QueryAgreeNoRequest{
		Inqwork: inqwork,
		Corpno:  config.GetString(enum.IcbcCorpNo, ""),
		Cond:    cond,
	}
	var result stru.QueryAgreeNoResponse
	resultInterface, err := stru.ICBCPostHttpResult(enum.IcbcAdsAgreementGryURL, *request)
	if err != nil {
		return nil, err
	}
	// 检查 response.ResponseBizContent 是否可以转换为 QueryAgreeNoResponse

	jsonString, _ := json.Marshal(resultInterface)
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil, err
	}

	if result.RetCode != "9008100" {
		return nil, errors.New(result.RetMsg)
	}
	if len(result.AgrList) == 0 {
		return nil, errors.New("根据相关编号未能查询到协议信息[" + "zuId]")
	}
	agreement := result.AgrList[0]
	return &agreement, nil
}

func (i *icbcBankSDK) IcbcUserAcctSignatureQuery(ctx context.Context, accountNo string, accCompNo string) (*stru.IcbcSignatureQueryResponse, error) {
	url := enum.IcbcAdsPartNerGryURL
	request := stru.NewIcbcGlobalRequest()
	corpNo, _ := strconv.ParseInt(config.GetString(enum.IcbcCorpNo, ""), 10, 64)
	request.BizContent = &stru.IcbcSignatureQueryRequest{
		StartIndex:  "0",
		QrySize:     "10",
		CorpNo:      corpNo,
		AccCompNo:   accCompNo,
		AccCompName: "",
	}
	var result stru.IcbcSignatureQueryResponse
	resultInterface, err := stru.ICBCPostHttpResult(url, *request)
	if err != nil {
		return nil, err
	}
	// 检查 response.ResponseBizContent 是否可以转换为 IcbcSignatureQueryResponse
	if responseValue, ok := resultInterface.(stru.IcbcSignatureQueryResponse); ok {
		result = responseValue
	} else {
		zap.L().Info(fmt.Sprintf("ICBCPostResult: ResponseBizContent 不是预期的 stru.IcbcSignatureQueryResponse resultInterface=%+v", resultInterface))
	}

	if result.RetCode != "9008100" {
		return nil, errors.New(result.RetMsg)
	}
	return &result, nil
}

func (i *icbcBankSDK) IcbcUserAcctSignatureApply(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (string, error) {
	appId := config.GetString(enum.IcbcAppId, "")
	corpNo := config.GetString(enum.IcbcCorpNo, "")
	coMode := "1"
	request := stru.NewIcbcGlobalRequest()
	acclist := make([]*stru.AccListItem, 0)
	acclist = append(acclist, &stru.AccListItem{
		Account:       accountNo,
		CurrType:      "1",
		AccFlag:       "1",
		CnTioFlag:     "1",
		IsMainAcc:     "1",
		ReceiptFlag:   "1",
		StatementFlag: "1",
	})
	request.BizContent = &stru.IcbcSignRequest{
		AppID:      appId,
		ApiName:    "ADSSIGN",
		ApiVersion: "001.001.001.001",
		CorpNo:     corpNo,
		CoMode:     coMode,
		AccCompNo:  accCompNo,
		Account:    accountNo,
		CurrType:   "1",
		AccFlag:    "1",
		CnTioFlag:  "1",
		Phone:      phone,
		EpType:     "1",
		EpTimes:    "12",
		Remark:     remark,
		AccList:    acclist,
		PayAccName: "",
		PayAccNo:   "",
		PayBegDate: "",
		PayCurr:    "",
		PayLimit:   "",
	}
	resu := stru.ICBCPostHttpUIResult(*request)
	zap.L().Info(fmt.Sprintf("==ICBCPostHttpUIResult%v", resu))

	return resu, nil
}
func (i *icbcBankSDK) IcbcUserAcctSignatureAgreePush(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (*stru.IcbcSignConfirmResponse, error) {
	appId := config.GetString(enum.IcbcAppId, "")
	corpNo := config.GetString(enum.IcbcCorpNo, "")
	coMode := "1"
	//申请签约之后把待确认信息同步到工行,
	confirmRequest := stru.NewIcbcGlobalRequest()
	confirmList := make([]*stru.ConfirmListItem, 0)
	confirmList = append(confirmList, &stru.ConfirmListItem{
		AccCompNo: accCompNo,
		AcCount:   accountNo,
		CurrType:  "1",
		CHANNEL:   "2",
	})
	confirmRequest.BizContent = &stru.IcbcSignConfirmRequest{
		AppID:       appId,
		CorpNo:      corpNo,
		CoMode:      coMode,
		ConfirmList: confirmList,
	}
	var confirmResu stru.IcbcSignConfirmResponse
	resultInterface, err := stru.ICBCPostHttpResult(enum.IcbcAdsAgrConfirmSynURL, *confirmRequest)
	if err != nil {
		return nil, err
	}
	// 检查 response.ResponseBizContent 是否可以转换为 IcbcSignConfirmResponse
	if responseValue, ok := resultInterface.(stru.IcbcSignConfirmResponse); ok {
		confirmResu = responseValue
	} else {
		zap.L().Info(fmt.Sprintf("ICBCPostResult: ResponseBizContent 不是预期的 stru.IcbcSignConfirmResponse resultInterface=%+v", resultInterface))
	}
	return &confirmResu, nil
}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, accCompNo string, agreeNo string) ([]stru.IcbcAccDetailItem, error) {

	begin, _ := time.Parse("20060102", beginDate)
	beginD := begin.Format("2006-01-02")

	end, _ := time.Parse("20060102", endDate)

	endD := end.Format("2006-01-02")
	serialNo := ""

	hasNext := true
	var resultDetails []stru.IcbcAccDetailItem
	for hasNext {
		request := stru.NewIcbcGlobalRequest()
		seq, _ := util.SonyflakeID()
		accDetailRequest := &stru.AccDetailRequest{
			FSeqNo:    seq,
			Account:   account,
			CurrType:  "1",
			StartDate: beginD,
			EndDate:   endD,
			SerialNo:  serialNo,
			CorpNo:    config.GetString(enum.IcbcCorpNo, ""),
			AccCompNo: accCompNo,
			AgreeNo:   agreeNo,
		}
		request.BizContent = accDetailRequest
		var result stru.AccDetailResponse
		resultInterface, err := stru.ICBCPostHttpResult(enum.IcbcAccDetailURL, *request)
		if err != nil {
			return nil, err
		}
		jsonString, _ := json.Marshal(resultInterface)
		err = json.Unmarshal(jsonString, &result)
		if err != nil {
			return nil, err
		}
		resultDetails = append(resultDetails, result.DtlList...)

		if result.NextPage == "1" {
			hasNext = true
			serialNo = result.SerialNo
		} else {
			hasNext = false
		}
	}
	return resultDetails, nil
}
