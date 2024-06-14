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
	"time"
)

type IcbcBankSDK interface {
	QueryAgreeNo(ctx context.Context, zuId, account string) (*stru.Agreement, error)
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, agreeNo string) ([]stru.IcbcAccDetailItem, error)
	IcbcReceiptNoQuery(ctx context.Context, accountNo, accCompNo, agreeNo, serialNo string) (*stru.IcbcReceiptNoQueryResponse, error)
	IcbcReceiptFileDownload(ctx context.Context, accountNo, accCompNo, agreeNo, serialNo string) (string, error)
}

type icbcBankSDK struct {
}

func (i *icbcBankSDK) IcbcReceiptFileDownload(ctx context.Context, accountNo, accCompNo, agreeNo, serialNo string) (string, error) {
	////生成privateKey
	//privateKeyBytes, err := base64.StdEncoding.DecodeString(config.GetString(enum.IcbcPrivateKey, ""))
	//if err != nil {
	//	fmt.Errorf("failed to decode Base64 encoded private key: %w", err)
	//}
	//
	//err = ioutil.WriteFile("./ssh_key.pem", privateKeyBytes, 0600)
	//
	//fmt.Println("SSH密钥文件已成功生成")
	//return "", nil
	client := stru.SftpClient()
	dir, err := client.ReadDir("/")
	if err != nil {
		return "", err
	}
	for _, v := range dir {
		zap.L().Info(fmt.Sprintf("sftp%s", v.Name()))
	}
	return "", err
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
		Account:   "",
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

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, agreeNo string) ([]stru.IcbcAccDetailItem, error) {
	//i.IcbcReceiptFileDownload(ctx, account, agreeNo, "")
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
			AccCompNo: config.GetString(enum.IcbcAccCompNo, ""),
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
