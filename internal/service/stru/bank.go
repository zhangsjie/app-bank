package stru

import (
	"encoding/json"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"time"
)

func ConvertBankTransferReceiptData(dbData repo.BankTransferReceiptDBData) *api.BankTransferReceiptData {
	createdAtBytes, _ := json.Marshal(dbData.CreatedAt)
	updatedAtBytes, _ := json.Marshal(dbData.UpdatedAt)
	return &api.BankTransferReceiptData{
		Id:                        dbData.Id,
		OrganizationId:            dbData.OrganizationId,
		CreatedAt:                 createdAtBytes,
		UpdatedAt:                 updatedAtBytes,
		ProcessInstanceId:         dbData.ProcessInstanceId,
		OriginatorUserName:        dbData.OriginatorUserName,
		SerialNo:                  dbData.SerialNo,
		PayAccount:                dbData.PayAccount,
		PayAccountName:            dbData.PayAccountName,
		RecAccount:                dbData.RecAccount,
		RecAccountName:            dbData.RecAccountName,
		PayAmount:                 dbData.PayAmount,
		CurrencyType:              dbData.CurrencyType,
		PayRem:                    dbData.PayRem,
		PubPriFlag:                dbData.PubPriFlag,
		RecBankType:               dbData.RecBankType,
		RecAccountOpenBank:        dbData.RecAccountOpenBank,
		UnionBankNo:               dbData.UnionBankNo,
		ClearBankNo:               dbData.ClearBankNo,
		RmtType:                   dbData.RmtType,
		TransferFlag:              dbData.TransferFlag,
		ChargeFee:                 dbData.ChargeFee,
		OrderState:                dbData.OrderState,
		RetCode:                   dbData.RetCode,
		RetMessage:                dbData.RetMessage,
		OrderFlowNo:               dbData.OrderFlowNo,
		ProcessBusinessId:         dbData.ProcessBusinessId,
		ProcessComment:            dbData.ProcessComment,
		DetailHostFlowNo:          dbData.DetailHostFlowNo,
		CommentUserId:             dbData.CommentUserId,
		CommentUserName:           dbData.CommentUserName,
		ProcessStatus:             dbData.ProcessStatus,
		PayAccountOpenBank:        dbData.PayAccountOpenBank,
		PayAccountOpenBankFilling: dbData.PayAccountOpenBankFilling,
		RecAccountOpenBankFilling: dbData.RecAccountOpenBankFilling,
		Title:                     dbData.Title,
		PayAccountType:            dbData.PayAccountType,
		ElectronicReceiptFile:     dbData.ElectronicReceiptFile,
	}
}

func ConvertBankTransferReceiptDBData(data api.BankTransferReceiptData) *repo.BankTransferReceiptDBData {
	var createdAt time.Time
	var updatedAt time.Time
	json.Unmarshal(data.CreatedAt, &createdAt)
	json.Unmarshal(data.UpdatedAt, &updatedAt)
	return &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id:        data.Id,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			OrganizationId: data.OrganizationId,
		},
		ProcessInstanceId:         data.ProcessInstanceId,
		OriginatorUserName:        data.OriginatorUserName,
		SerialNo:                  data.SerialNo,
		PayAccount:                data.PayAccount,
		PayAccountName:            data.PayAccountName,
		RecAccount:                data.RecAccount,
		RecAccountName:            data.RecAccountName,
		PayAmount:                 data.PayAmount,
		CurrencyType:              data.CurrencyType,
		PayRem:                    data.PayRem,
		PubPriFlag:                data.PubPriFlag,
		RecBankType:               data.RecBankType,
		RecAccountOpenBank:        data.RecAccountOpenBank,
		UnionBankNo:               data.UnionBankNo,
		ClearBankNo:               data.ClearBankNo,
		RmtType:                   data.RmtType,
		TransferFlag:              data.TransferFlag,
		ChargeFee:                 data.ChargeFee,
		OrderState:                data.OrderState,
		RetCode:                   data.RetCode,
		RetMessage:                data.RetMessage,
		OrderFlowNo:               data.OrderFlowNo,
		DetailHostFlowNo:          data.DetailHostFlowNo,
		ProcessBusinessId:         data.ProcessBusinessId,
		ProcessComment:            data.ProcessComment,
		CommentUserId:             data.CommentUserId,
		CommentUserName:           data.CommentUserName,
		ProcessStatus:             data.ProcessStatus,
		RecAccountOpenBankFilling: data.RecAccountOpenBankFilling,
		PayAccountOpenBankFilling: data.PayAccountOpenBankFilling,
		PayAccountOpenBank:        data.PayAccountOpenBank,
		Title:                     data.Title,
		PayAccountType:            data.PayAccountType,
	}
}

func ConvertBankTransactionDetailData(dbData repo.BankTransactionDetailDBData) *api.BankTransactionDetailData {
	createdAtBytes, _ := json.Marshal(dbData.CreatedAt)
	updatedAtBytes, _ := json.Marshal(dbData.UpdatedAt)
	return &api.BankTransactionDetailData{
		Id:                    dbData.Id,
		OrganizationId:        dbData.OrganizationId,
		CreatedAt:             createdAtBytes,
		UpdatedAt:             updatedAtBytes,
		MerchantAccountId:     dbData.MerchantAccountId,
		MerchantAccountName:   dbData.MerchantAccountName,
		Type:                  dbData.Type,
		CashFlag:              dbData.CashFlag,
		PayAmount:             dbData.PayAmount,
		RecAmount:             dbData.RecAmount,
		BsnType:               dbData.BsnType,
		TransferDate:          dbData.TransferDate,
		TransferTime:          dbData.TransferTime,
		TranChannel:           dbData.TranChannel,
		CurrencyType:          dbData.CurrencyType,
		Balance:               dbData.Balance,
		OrderFlowNo:           dbData.OrderFlowNo,
		HostFlowNo:            dbData.HostFlowNo,
		VouchersType:          dbData.VouchersType,
		VouchersNo:            dbData.VouchersNo,
		SummaryNo:             dbData.SummaryNo,
		Summary:               dbData.Summary,
		AcctNo:                dbData.AcctNo,
		AccountName:           dbData.AccountName,
		AccountOpenNode:       dbData.AccountOpenNode,
		ElectronicReceiptFile: dbData.ElectronicReceiptFile,
		ProcessBusinessId:     dbData.ProcessBusinessId,
		ProcessInstanceId:     dbData.ProcessInstanceId,
		OriginatorUser:        dbData.OriginatorUserName,
		OperationUser:         dbData.OperationUserName,
		OperationComment:      dbData.OperationComment,
		ProcessTotalStatus:    dbData.ProcessTotalStatus,
		PayAccountType:        dbData.PayAccountType,
		ExtField3:             dbData.ExtField3,
	}
}

func ConvertBankTransactionDetailDataAndMerchantAccount(dbData repo.BankTransactionDetailDBData, merchantAccount string) *api.BankTransactionDetailData {
	createdAtBytes, _ := json.Marshal(dbData.CreatedAt)
	updatedAtBytes, _ := json.Marshal(dbData.UpdatedAt)
	return &api.BankTransactionDetailData{
		Id:                    dbData.Id,
		OrganizationId:        dbData.OrganizationId,
		CreatedAt:             createdAtBytes,
		UpdatedAt:             updatedAtBytes,
		MerchantAccountId:     dbData.MerchantAccountId,
		MerchantAccount:       merchantAccount,
		MerchantAccountName:   dbData.MerchantAccountName,
		Type:                  dbData.Type,
		CashFlag:              dbData.CashFlag,
		PayAmount:             dbData.PayAmount,
		RecAmount:             dbData.RecAmount,
		BsnType:               dbData.BsnType,
		TransferDate:          dbData.TransferDate,
		TransferTime:          dbData.TransferTime,
		TranChannel:           dbData.TranChannel,
		CurrencyType:          dbData.CurrencyType,
		Balance:               dbData.Balance,
		OrderFlowNo:           dbData.OrderFlowNo,
		HostFlowNo:            dbData.HostFlowNo,
		VouchersType:          dbData.VouchersType,
		VouchersNo:            dbData.VouchersNo,
		SummaryNo:             dbData.SummaryNo,
		Summary:               dbData.Summary,
		AcctNo:                dbData.AcctNo,
		AccountName:           dbData.AccountName,
		AccountOpenNode:       dbData.AccountOpenNode,
		ElectronicReceiptFile: dbData.ElectronicReceiptFile,
		ProcessBusinessId:     dbData.ProcessBusinessId,
		ProcessInstanceId:     dbData.ProcessInstanceId,
		OriginatorUser:        dbData.OriginatorUserName,
		OperationUser:         dbData.OperationUserName,
		OperationComment:      dbData.OperationComment,
		ProcessTotalStatus:    dbData.ProcessTotalStatus,
		PayAccountType:        dbData.PayAccountType,
	}
}

func ConvertBankCodeData(dbData repo.BankCodeDBData) *api.BankCodeData {
	return &api.BankCodeData{
		Id:            dbData.Id,
		BankName:      dbData.BankName,
		BankAliasName: dbData.BankAliasName,
		BankCode:      dbData.BankCode,
		UnionBankNo:   dbData.UnionBankNo,
		ClearBankNo:   dbData.ClearBankNo,
	}
}

func ConvertAddBankCodeDBData(data api.AddBankCodeRequest) *repo.BankCodeDBData {
	return &repo.BankCodeDBData{
		BankName:      data.BankName,
		BankAliasName: data.BankAliasName,
		UnionBankNo:   data.UnionBankNo,
		ClearBankNo:   data.ClearBankNo,
	}
}

func ConvertBankCodeDBData(data api.BankCodeData) *repo.BankCodeDBData {
	return &repo.BankCodeDBData{
		BaseCommonDBData: repository.BaseCommonDBData{
			Id: data.Id,
		},
		BankName:      data.BankName,
		BankAliasName: data.BankAliasName,
		UnionBankNo:   data.UnionBankNo,
		ClearBankNo:   data.ClearBankNo,
	}
}

type DingtalkAttachmentData struct {
	SpaceId  string `json:"spaceId"`
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
	FileType string `json:"fileType"`
	FileId   string `json:"fileId"`
}

func FormatDayTime(dateTime time.Time) string {
	timeUnix := dateTime.Unix() //已知的时间戳
	formatTimeStr := time.Unix(timeUnix, 0).Format("20060102")
	return formatTimeStr
}

func FormatDayTimeStamp(timestamp int64) string {
	formatTimeStr := time.Unix(timestamp, 0).Format("20060102")
	return formatTimeStr
}

type SyncDateRequest struct {
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

func FormatMonthChartTimeLabelData(year int, month time.Month) []string {
	var timeRange []string
	timeFormatType := "20060102"
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	_, _, lastDay := thisMonth.AddDate(0, 1, -1).Date()
	for i := 0; i < lastDay; i++ {
		dayTime := thisMonth.AddDate(0, 0, i).Format(timeFormatType)
		timeRange = append(timeRange, dayTime[4:])
	}
	return timeRange
}

func FormatMonthTimeCondition(year int, month time.Month) *repo.TransactionDetailTimeParam {
	videoUvReq := repo.TransactionDetailTimeParam{}
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	videoUvReq.StartTime = thisMonth.AddDate(0, 0, 0).Format("2006-01-02") + " 00:00:00"
	videoUvReq.EndTime = thisMonth.AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
	return &videoUvReq
}

func FormatMonthTimeConditionMore(year int, month time.Month, offset int) *repo.TransactionDetailTimeParam {
	videoUvReq := repo.TransactionDetailTimeParam{}
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	videoUvReq.StartTime = thisMonth.AddDate(0, 0, offset).Format("2006-01-02") + " 00:00:00"
	videoUvReq.EndTime = thisMonth.AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
	return &videoUvReq
}

func FormatMonthChartTimeLabelDataMore(year int, month time.Month, offset int) ([]string, int) {
	var timeRange []string
	timeFormatType := "20060102"
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	_, _, lastDay := thisMonth.AddDate(0, 1, -1).Date()
	for i := offset; i < 0; i++ {
		dayTime := thisMonth.AddDate(0, 0, i).Format(timeFormatType)
		timeRange = append(timeRange, dayTime[4:])
	}
	for i := 0; i < lastDay; i++ {
		dayTime := thisMonth.AddDate(0, 0, i).Format(timeFormatType)
		timeRange = append(timeRange, dayTime[4:])
	}
	return timeRange, lastDay
}

func FormatWeekChartTimeLabelData(offset int) []string {
	var timeRange []string
	timeFormatType := "20060102"

	year, month, day := time.Now().Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startOffset := offset + 1
	for i := startOffset; i < 0; i++ {
		dayTime := thisWeek.AddDate(0, 0, i).Format(timeFormatType)
		timeRange = append(timeRange, dayTime[4:])
	}
	timeRange = append(timeRange, thisWeek.Format(timeFormatType)[4:])
	return timeRange
}

func FormatDayTimeCondition(offset int) *repo.TransactionDetailTimeParam {
	videoUvReq := repo.TransactionDetailTimeParam{}
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	videoUvReq.StartTime = today.AddDate(0, 0, offset+1).Format("2006-01-02") + " 00:00:00"
	videoUvReq.EndTime = today.Format("2006-01-02") + " 23:59:59"
	return &videoUvReq
}

func FormatMonthTime(month int) time.Month {
	monthTime := time.January
	switch month {
	case 1:
		monthTime = time.January
	case 2:
		monthTime = time.February
	case 3:
		monthTime = time.March
	case 4:
		monthTime = time.April
	case 5:
		monthTime = time.May
	case 6:
		monthTime = time.June
	case 7:
		monthTime = time.July
	case 8:
		monthTime = time.August
	case 9:
		monthTime = time.September
	case 10:
		monthTime = time.October
	case 11:
		monthTime = time.November
	case 12:
		monthTime = time.December
	default:
		_, month, _ := time.Now().Date()
		monthTime = month
	}
	return monthTime
}

type ChartData struct {
	Labels   []string        `json:"labels"`
	Datasets []*ChartDataSet `json:"datasets"`
}

type ChartDataSet struct {
	Label                string    `json:"label"`
	Data                 []float64 `json:"data"`
	BorderColor          string    `json:"borderColor"`
	BackgroundColor      []string  `json:"backgroundColor"`
	HoverBackgroundColor []string  `json:"hoverBackgroundColor"`
}

func FormatLineChartMultipleDataSetArray(labelList []string, dayArray []string, payAmountArray []float64, recAmountArray []float64) []*api.ChartDataSet {
	var chartData []*api.ChartDataSet
	var firstData api.ChartDataSet
	firstDataList := make([]float64, len(labelList))
	for i, item := range labelList {
		for m, label := range dayArray {
			if item == label {
				firstDataList[i] = payAmountArray[m]
				break
			} else {
				firstDataList[i] = 0
			}
		}
	}
	firstData.Data = firstDataList
	firstData.Label = "支出"
	firstData.BorderColor = "#42A5F5"
	firstData.BackgroundColor = []string{"#42A5F5"}
	firstData.HoverBackgroundColor = []string{"#64B5F6"}
	chartData = append(chartData, &firstData)

	var secondData api.ChartDataSet
	secondDataList := make([]float64, len(labelList))
	for i, item := range labelList {
		for m, label := range dayArray {
			if item == label {
				secondDataList[i] = recAmountArray[m]
				break
			} else {
				secondDataList[i] = 0
			}
		}
	}
	secondData.Data = secondDataList
	secondData.Label = "收入"
	secondData.BorderColor = "#FFA726"
	secondData.BackgroundColor = []string{"#FFA726"}
	secondData.HoverBackgroundColor = []string{"#FFB74D"}
	chartData = append(chartData, &secondData)
	return chartData
}

func FormatLineChartSingleDataSetArray(labelList []string, dataLabelArray []string, dataArray []float64) []*api.ChartDataSet {
	var chartData []*api.ChartDataSet
	var data api.ChartDataSet
	dataList := make([]float64, len(labelList))
	for i, item := range labelList {
		for m, label := range dataLabelArray {
			if item == label {
				dataList[i] = dataArray[m]
				break
			} else {
				dataList[i] = 0
			}
		}
	}
	data.Label = "余额"
	data.Data = dataList
	data.BorderColor = "#66BB6A"
	data.BackgroundColor = []string{"#66BB6A"}
	data.HoverBackgroundColor = []string{"#81C784"}
	chartData = append(chartData, &data)
	return chartData
}

func HandleBalanceData(weekTimeLabels []string, balanceDataMap map[string]float64, dayCount int) (*[]string, *[]float64) {
	timeLabels := weekTimeLabels[7:]
	dataArray := make([]float64, dayCount)
	backTimeLabels := weekTimeLabels[0:7]
	for i, time := range timeLabels {
		if i == 0 {
			if _, ok := balanceDataMap[time]; ok && balanceDataMap[time] > 0 {
				dataArray[i] = balanceDataMap[time]
			} else {
				//从之前时间中获取余额数据
				for m := 6; m <= 0; m-- {
					if _, keyOK := balanceDataMap[backTimeLabels[m]]; keyOK && balanceDataMap[backTimeLabels[m]] > 0 {
						dataArray[i] = balanceDataMap[backTimeLabels[m]]
						break
					}
				}
			}
		} else {
			if _, ok := balanceDataMap[time]; ok && balanceDataMap[time] > 0 {
				dataArray[i] = balanceDataMap[time]
			} else {
				dataArray[i] = dataArray[i-1]
			}
		}
	}
	return &timeLabels, &dataArray
}

type ImportPayrollDataRequest struct {
	PayAccountName string  `json:"payAccountName" form:"payAccountName"`
	PayAccountNo   string  `json:"payAccountNo" form:"payAccountNo"`
	Remark         string  `json:"remark" form:"remark"`
	Month          string  `json:"month" form:"month"`
	FileUrl        string  `json:"fileUrl" form:"fileUrl"`
	FileData       []byte  `json:"fileData" form:"fileData"`
	TotalAmount    float64 `json:"totalAmount" form:"totalAmount"`
	TotalCount     int64   `json:"totalCount" form:"totalCount"`
}
type PdfToImageJsdkRequest struct {
	OssPath    string `json:"ossPath"`
	PdfUrl     string `json:"pdfUrl"`
	SourceType int    `json:"sourceType"`
}
type PdfToImageJsdkResponse struct {
	OssPath string `json:"ossPath"`
}

type MinShengTransactionDetailResponse struct {
	TransSeqNo    string `json:"trans_seq_no"`
	RecSeq        string `json:"rec_seq"`
	AcctNo        string `json:"acct_no"`
	AcctName      string `json:"acct_name"`
	DcFlag        string `json:"dc_flag"`
	EnterAcctDate string `json:"enter_acct_date"`
	IntrDate      string `json:"intr_date"`
	Amount        string `json:"amount"`
	CpAcctNo      string `json:"cp_acct_no"`
	CpAcctName    string `json:"cp_acct_name"`
	CpBankName    string `json:"cp_bank_name"`
	CpBankAddr    string `json:"cp_bank_addr"`
	Explain       string `json:"explain"`
	Balance       string `json:"balance"`
	Currency      string `json:"currency"`
	Timestamp     string `json:"timestamp"`
}
