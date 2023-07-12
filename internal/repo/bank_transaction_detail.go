package repo

import (
	"context"
	"fmt"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"strconv"
)

type BankTransactionDetailDBData struct {
	repository.BaseDBData
	MerchantAccountId     int64
	MerchantAccountName   string
	Type                  string
	CashFlag              string
	PayAmount             float64
	RecAmount             float64
	BsnType               string
	TransferDate          string
	TransferTime          string
	TranChannel           string
	CurrencyType          string
	Balance               float64
	OrderFlowNo           string
	HostFlowNo            string
	VouchersType          string
	VouchersNo            string
	SummaryNo             string
	Summary               string
	AcctNo                string
	AccountName           string
	AccountOpenNode       string
	ElectronicReceiptFile string
	ElectronicReceiptPng  string
	ProcessBusinessId     string
	ProcessInstanceId     int64
	OriginatorUserId      string
	OriginatorUserName    string
	OperationUserId       string
	OperationUserName     string
	OperationComment      string
	ProcessTotalStatus    string
	PayAccountType        string
	ExtField1             string
	ExtField2             string
	ExtField3             string
	ExtField4             string
	ExtField5             string
}

type BankTransactionDetailDBDataParam struct {
	BankTransactionDetailDBData
	PayAmountMin                float64
	PayAmountMax                float64
	RecAmountMin                float64
	RecAmountMax                float64
	TransferTimeArray           []string
	IsElectronicReceiptFileNull bool
	IsAccountNoNull             bool
}

type TransactionDetailTimeParam struct {
	StartTime string
	EndTime   string
}

func (*BankTransactionDetailDBData) TableName() string {
	return "bank_transaction_detail"
}

func (param *BankTransactionDetailDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankTransactionDetailDBDataParam) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	if param.Type != "" {
		conditions = append(conditions, repository.NewAndCondition("type = ?", param.Type))
	}
	if param.ProcessBusinessId != "" {
		conditions = append(conditions, repository.NewAndCondition("process_business_id like ?", util.GetLikeString(param.ProcessBusinessId)))
	}
	if param.OriginatorUserName != "" {
		conditions = append(conditions, repository.NewAndCondition("originator_user_name like ?", util.GetLikeString(param.OriginatorUserName)))
	}
	if param.OperationUserName != "" {
		conditions = append(conditions, repository.NewAndCondition("operation_user_name like ?", util.GetLikeString(param.OperationUserName)))
	}
	if param.MerchantAccountName != "" {
		conditions = append(conditions, repository.NewAndCondition("merchant_account_name like ?", util.GetLikeString(param.MerchantAccountName)))
	}
	if param.AccountName != "" {
		conditions = append(conditions, repository.NewAndCondition("account_name like ?", util.GetLikeString(param.AccountName)))
	}
	if param.HostFlowNo != "" {
		conditions = append(conditions, repository.NewAndCondition("host_flow_no like ?", util.GetLikeString(param.HostFlowNo)))
	}
	if param.ProcessTotalStatus != "" {
		conditions = append(conditions, repository.NewAndCondition("process_total_status = ?", param.ProcessTotalStatus))
	}
	if len(param.TransferTimeArray) > 1 {
		if param.TransferTimeArray[0] != "" {
			conditions = append(conditions, repository.NewAndCondition("transfer_time >= DATE_FORMAT(?,'%Y%m%d%H:%i:%s')", param.TransferTimeArray[0]))
		}
		if param.TransferTimeArray[1] != "" {
			conditions = append(conditions, repository.NewAndCondition("transfer_time < DATE_FORMAT(date_add(?, interval 1 day),'%Y%m%d%H:%i:%s')", param.TransferTimeArray[1]))
		}
	}
	if param.BsnType != "" {
		conditions = append(conditions, repository.NewAndCondition("bsn_type = ?", param.BsnType))
	}
	if param.PayAmountMin > 0 {
		conditions = append(conditions, repository.NewAndCondition("pay_amount >= ?", param.PayAmountMin))
	}
	if param.PayAmountMax > 0 {
		conditions = append(conditions, repository.NewAndCondition("pay_amount <= ?", param.PayAmountMax))
	}
	if param.RecAmountMin > 0 {
		conditions = append(conditions, repository.NewAndCondition("rec_amount >= ?", param.RecAmountMin))
	}
	if param.RecAmountMax > 0 {
		conditions = append(conditions, repository.NewAndCondition("rec_amount <= ?", param.RecAmountMax))
	}
	if param.PayAccountType != "" {
		conditions = append(conditions, repository.NewAndCondition("pay_account_type = ?", param.PayAccountType))
	}
	if param.IsElectronicReceiptFileNull {
		conditions = append(conditions, repository.NewAndCondition("electronic_receipt_file = '' or electronic_receipt_file is null "))
	}
	if param.ExtField3 != "" {
		conditions = append(conditions, repository.NewAndCondition("ext_field3 = ?", param.ExtField3))
	}
	if param.IsAccountNoNull {
		conditions = append(conditions, repository.NewAndCondition("acct_no is not || acct_no = ''"))
	} else {
		conditions = append(conditions, repository.NewAndCondition("acct_no != '' and acct_no is not null "))
	}
	return conditions
}

type BankTransactionDetailRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankTransactionDetailDBData) (*BankTransactionDetailDBData, error)
	Count(context.Context, *BankTransactionDetailDBDataParam) (int64, error)
	List(context.Context, string, int32, int32, *BankTransactionDetailDBDataParam) (*[]BankTransactionDetailDBData, int64, error)
	CashFlowCount(*TransactionDetailTimeParam, int64) (*[]string, *[]float64, *[]float64, error)
	BalanceFlowCountMap(*TransactionDetailTimeParam, int64) (map[string]float64, error)
}

type bankTransactionDetailRepo struct {
	repository.BaseRepo
}

func (r *bankTransactionDetailRepo) Get(ctx context.Context, param *BankTransactionDetailDBData) (*BankTransactionDetailDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankTransactionDetailDBData), handler.HandleError(err)
}

func (r *bankTransactionDetailRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankTransactionDetailDBDataParam) (*[]BankTransactionDetailDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankTransactionDetailDBData), count, handler.HandleError(err)
}

func (r *bankTransactionDetailRepo) Count(ctx context.Context, param *BankTransactionDetailDBDataParam) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}

func (r *bankTransactionDetailRepo) CashFlowCount(param *TransactionDetailTimeParam, organizationId int64) (*[]string, *[]float64, *[]float64, error) {
	simpleSql := fmt.Sprintf("select DATE_FORMAT(created_at,'%%m%%d') as days, sum(pay_amount) as payAmount, sum(rec_amount) as recAmount from bank_transaction_detail where created_at >= '%s' and created_at <= '%s' and organization_id = %d group by days ", param.StartTime, param.EndTime, organizationId)
	var result []map[string]interface{}
	err := r.Db.Raw(simpleSql).Find(&result).Error
	if err != nil {
		return nil, nil, nil, handler.HandleError(err)
	}
	dayArray := make([]string, len(result))
	payAmountArray := make([]float64, len(result))
	recAmountArray := make([]float64, len(result))
	for i, row := range result {
		for key, value := range row {
			switch key {
			case "days":
				dayArray[i] = value.(string)
			case "payAmount":
				amount, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					handler.HandleError(err)
				}
				payAmountArray[i] = amount

			case "recAmount":
				amount, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					handler.HandleError(err)
				}
				recAmountArray[i] = amount
			}
		}
	}
	return &dayArray, &payAmountArray, &recAmountArray, nil
}

func (r *bankTransactionDetailRepo) BalanceFlowCountMap(param *TransactionDetailTimeParam, organizationId int64) (map[string]float64, error) {
	simpleSql := fmt.Sprintf("select d.balance as balance , a.days as days from (select DATE_FORMAT(created_at,'%%m%%d') as days, max(`id`) as serial from bank_transaction_detail where created_at >= '%s' and created_at <= '%s' group by days) as a inner join `bank_transaction_detail` as d on a.serial = d.id where d.organization_id = %d ", param.StartTime, param.EndTime, organizationId)
	var result []map[string]interface{}
	err := r.Db.Raw(simpleSql).Find(&result).Error
	if err != nil {
		return nil, handler.HandleError(err)
	}

	balanceArray := make(map[string]float64, len(result))
	for _, row := range result {
		dayVal := ""
		balanceVal := 0.00
		for key, value := range row {
			switch key {
			case "days":
				dayVal = value.(string)
			case "balance":
				balance, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					handler.HandleError(err)
				}
				balanceVal = balance
			}
		}
		balanceArray[dayVal] = balanceVal
	}
	return balanceArray, nil
}
