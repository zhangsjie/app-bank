package sdk

import (
	"context"
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
	QueryAgreeNo(ctx context.Context, zuId, account string) (string, error)
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, accCompNo string) ([]stru.IcbcAccDetailItem, error)
	IcbcUserAcctSignatureApply(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (string, error)
	IcbcUserAcctSignatureAgreePush(ctx context.Context, accountNo string, phone string, remark string, accCompNo string) (*stru.IcbcSignConfirmResponse, error)
	IcbcUserAcctSignatureQuery(ctx context.Context, accountNo string, accCompNo string) (*stru.IcbcSignatureQueryResponse, error)
}

type icbcBankSDK struct {
}

func (i *icbcBankSDK) QueryAgreeNo(ctx context.Context, zuId, account string) (string, error) {
	//IcbcAdsAgreementGryURL
	request := stru.NewIcbcGlobalRequest()
	inqwork := stru.Inqwork{
		BegNum: 0,
		FetNum: 10,
	}
	cond := stru.Cond{
		QryType:   1,
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
	err := stru.ICBCPostHttpResult(enum.IcbcAdsAgreementGryURL, *request, &result)
	if err != nil {
		return "", err
	}
	if result.RetCode != "9008100" {
		return "", errors.New(result.RetMsg)
	}
	agreeNo := result.AgrList[0].AgreeNo
	return agreeNo, nil
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
	err := stru.ICBCPostHttpResult(url, *request, &result)
	if err != nil {
		return nil, err
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
	confirmResu := stru.IcbcSignConfirmResponse{}
	err := stru.ICBCPostHttpResult(enum.IcbcAdsAgrConfirmSynURL, *confirmRequest, &confirmResu)
	if err != nil {
		return nil, err
	}

	return &confirmResu, nil
}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, accCompNo string) ([]stru.IcbcAccDetailItem, error) {
	accDetailUrlPath := config.GetString(enum.IcbcAccDetailURL, "")
	request := stru.NewIcbcGlobalRequest()
	seq, _ := util.SonyflakeID()

	begin, _ := time.Parse("20060102", beginDate)
	beginD := begin.Format("2006-01-02")

	end, _ := time.Parse("20060102", endDate)

	endD := end.Format("2006-01-02")
	serialNo := ""
	accDetailRequest := &stru.AccDetailRequest{
		FSeqNo:    seq,
		Account:   account,
		CurrType:  1,
		StartDate: beginD,
		EndDate:   endD,
		SerialNo:  serialNo,
		CorpNo:    config.GetString(enum.IcbcCorpNo, ""),
		AccCompNo: config.GetString(enum.IcbcCorpNo, ""),
		AgreeNo:   "",
	}
	request.BizContent = accDetailRequest
	var result stru.AccDetailResponse
	err := stru.ICBCPostHttpResult(accDetailUrlPath, *request, &result)
	if err != nil {
		return nil, err
	}
	hasNext := false
	var resultDetail []stru.IcbcAccDetailItem
	resultDetail = append(resultDetail, result.DtlList...)
	if result.NextPage == "1" {
		hasNext = true
		serialNo = result.DtlList[len(result.DtlList)-1].SerialNo
	}
	for hasNext {
		accDetailRequest.SerialNo = serialNo
		err := stru.ICBCPostHttpResult(accDetailUrlPath, *request, &result)
		if err != nil {
			return nil, err
		}
		resultDetail = append(resultDetail, result.DtlList...)
		if result.NextPage == "1" {
			hasNext = true
			serialNo = result.DtlList[len(result.DtlList)-1].SerialNo
		}
	}
	return resultDetail, nil
}
