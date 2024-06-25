package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	czip "github.com/dablelv/cyan/zip"
	"github.com/pkg/errors"
	"gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type IcbcBankSDK interface {
	QueryAgreeNo(ctx context.Context, zuId, account string) (*stru.Agreement, error)
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, agreeNo string) ([]stru.IcbcAccDetailItem, error)
	IcbcReceiptNoQuery(ctx context.Context, accountNo, agreeNo, serialNo string) (*stru.IcbcReceiptNoQueryResponse, error)
	IcbcReceiptFileDownload(ctx context.Context) (string, error)
}

type icbcBankSDK struct {
}

func (i *icbcBankSDK) IcbcReceiptFileDownload(ctx context.Context) (string, error) {
	client := stru.SftpClient()
	defer client.Close()
	remotePath := "/Ebillsend/download"
	downloadDir, err := client.ReadDir(remotePath)
	if err != nil {
		return "", err
	}
	if downloadDir == nil || len(downloadDir) == 0 {
		zap.L().Info(fmt.Sprintf("IcbcReceiptFileDownload当前文件夹为空"))
		return "", nil
	}
	//下载.删除所有文件,然后重新从文件服务器上下载
	filePathDir, _ := util.Mkdir("tempFile/icbc")
	err = os.RemoveAll("filePathDir")
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(filePathDir); os.IsNotExist(err) {
		os.MkdirAll(filePathDir, 0755)
	}
	var wg sync.WaitGroup
	//下载zip包到本地临时路径并且解压缩
	for _, v := range downloadDir {
		wg.Add(1)
		fileName := v.Name()
		remoteFilePath := path.Join(remotePath, fileName)
		localFilePath := path.Join(filePathDir, fileName)
		if !strings.HasSuffix(fileName, ".zip") {
			continue
		}
		//创建本地文件
		localFile, err := os.Create(localFilePath)
		if err != nil {
			zap.L().Info(fmt.Sprintf("IcbcReceiptFileDownload创建本地文件失败： %s", err))
			continue
		}

		//打开远程文件
		remoteFile, err := client.Open(remoteFilePath)
		if err != nil {
			zap.L().Info(fmt.Sprintf("IcbcReceiptFileDownload打开远程文件失败： %s", err))
			continue
		}
		//复制远程文件到本地
		_, err = io.Copy(localFile, remoteFile)
		if err != nil {
			zap.L().Info(fmt.Sprintf("IcbcReceiptFileDownload复制文件内容失败： %s", err))
			continue
		}

		localFile.Close()
		remoteFile.Close()
		go func() {
			err = czip.Unzip(localFilePath, filePathDir)
			if err != nil {
				zap.L().Info(fmt.Sprintf("IcbcReceiptFileDownload解压zip包失败： %s", err))
			}
		}()
	}
	wg.Wait()
	return "", err
}

func (i *icbcBankSDK) IcbcReceiptNoQuery(ctx context.Context, accountNo, agreeNo, serialNo string) (*stru.IcbcReceiptNoQueryResponse, error) {
	request := stru.NewIcbcGlobalRequest()
	serl := []string{serialNo}
	cond := stru.IcbcReceiptNoQueryCond{
		SeqList: serl,
	}
	fseqNo, _ := util.SonyflakeID()
	request.BizContent = &stru.IcbcReceiptNoQueryRequest{
		FseqNo:    fseqNo,
		CorpNo:    config.GetString(enum.IcbcCorpNo, ""),
		AccCompNo: config.GetString(enum.IcbcAccCompNo, ""),
		AgreeNo:   agreeNo,
		Account:   accountNo,
		CurrType:  "001",
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
		AccCompNo: config.GetString(enum.IcbcAccCompNo, ""),
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
	jsonString, _ := json.Marshal(resultInterface)
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil, err
	}
	if result.RetCode != "9008100" {
		return nil, errors.New(result.RetMsg)
	}
	if len(result.AgrList) == 0 {
		return nil, errors.New("未能查询到账号协议信息[" + "zuId]")
	}
	agreement := result.AgrList[0]
	return &agreement, nil
}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string, agreeNo string) ([]stru.IcbcAccDetailItem, error) {
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
