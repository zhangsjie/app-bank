package sdk

import (
	"context"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"strconv"
	"time"
)

type IcbcBankSDK interface {
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error)
	IcbcUserAcctSignatureApply(ctx context.Context, accountNo string, phone string, remark string) (string, error)
	IcbcUserAcctSignatureQuery(ctx context.Context, accountNo string, phone string, remark string) (*stru.IcbcSignatureQueryResponse, error)
}

type icbcBankSDK struct {
}

func (i *icbcBankSDK) IcbcUserAcctSignatureQuery(ctx context.Context, accountNo string, phone string, remark string) (*stru.IcbcSignatureQueryResponse, error) {
	queryUrl := config.GetString(enum.IcbcHost, "") + enum.IcbcAdsPartNerGryURL
	request := stru.NewIcbcGlobalRequest()
	corpNo := config.GetString(enum.IcbcCorpNo, "")
	num, _ := strconv.ParseInt(corpNo, 10, 64)
	request.BizContent = &stru.IcbcSignatureQueryRequest{
		StartIndex:  "0",
		QrySize:     "10",
		CorpNo:      num,
		AccCompNo:   "",
		AccCompName: "",
	}
	var result stru.IcbcSignatureQueryResponse
	err := stru.ICBCPostHttpResult(queryUrl, *request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (i *icbcBankSDK) IcbcUserAcctSignatureApply(ctx context.Context, accountNo string, phone string, remark string) (string, error) {
	signUrl := config.GetString(enum.IcbcHost, "") + config.GetString(enum.IcbcAdsAgrSigUiURL, "")
	request := stru.NewIcbcGlobalRequest()
	acclist := make([]*stru.AccListItem, 0)
	acclist = append(acclist, &stru.AccListItem{
		Account:       accountNo,
		CurrType:      "1",
		AccFlag:       "1",
		CnTioFlag:     "1",
		IsMainAcc:     "0",
		ReceiptFlag:   "1",
		StatementFlag: "1",
	})
	request.BizContent = &stru.IcbcSignRequest{
		AppID:      config.GetString(enum.IcbcAppId, ""),
		ApiName:    "ADSSIGN",
		ApiVersion: "001.001.001.001",
		CorpNo:     config.GetString(enum.IcbcCorpNo, ""),
		CoMode:     "2",
		AccCompNo:  "",
		Account:    config.GetString(enum.IcbcAccountNo, ""),
		CurrType:   "1",
		AccFlag:    "1",
		CnTioFlag:  "1",
		Phone:      phone,
		EpType:     "1",
		EpTimes:    "12",
		Remark:     remark,
		AccList:    acclist,
	}
	resu := stru.ICBCPostHttpUIResult(signUrl, *request)

	return resu, nil
}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error) {
	accDetailUrl := config.GetString(enum.IcbcHost, "") + config.GetString(enum.IcbcAccDetailURL, "")
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
	err := stru.ICBCPostHttpResult(accDetailUrl, *request, &result)
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
		err := stru.ICBCPostHttpResult(accDetailUrl, *request, &result)
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
