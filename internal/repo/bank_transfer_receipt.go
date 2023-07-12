package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
)

type BankTransferReceiptDBData struct {
	repository.BaseDBData
	Title                     string
	ProcessInstanceId         int64
	OriginatorUserName        string
	SerialNo                  string
	PayAccount                string
	PayAccountName            string
	RecAccount                string
	RecAccountName            string
	PayAmount                 float64
	CurrencyType              string
	PayRem                    string
	PubPriFlag                string
	RecBankType               string
	RecAccountOpenBank        string
	UnionBankNo               string
	ClearBankNo               string
	RmtType                   string
	TransferFlag              string
	ChargeFee                 float64
	OrderState                string
	RetCode                   string
	RetMessage                string
	OrderFlowNo               string
	ProcessBusinessId         string
	ProcessComment            string
	DetailHostFlowNo          string
	CommentUserId             int64
	CommentUserName           string
	ProcessStatus             string
	RecAccountOpenBankFilling string
	PayAccountOpenBankFilling string
	PayAccountOpenBank        string
	PayAccountType            string
}

type BankTransferReceiptDBDataParam struct {
	BankTransferReceiptDBData
	ProcessInstanceIds    []int64
	ExcludeOrderStates    []string
	IsOrderFlowNoNotEmpty bool
	SerialNo              string
	BusinessId            string
	PayAccount            string
	RecAccount            string
	PayAmount             []float64
	OriginatorUser        string
	CommentUser           string
	CreateTime            []string
	TotalStatus           string
	OrderState            string
	PayAmountMin          float64
	PayAmountMax          float64
	PayAccountType        string
}

func (*BankTransferReceiptDBData) TableName() string {
	return "bank_transfer_receipt"
}

func (param *BankTransferReceiptDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankTransferReceiptDBDataParam) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	if param.ProcessInstanceId != 0 {
		conditions = append(conditions, repository.NewAndCondition("process_instance_id = ?", param.ProcessInstanceId))
	}
	if len(param.ProcessInstanceIds) > 0 {
		conditions = append(conditions, repository.NewAndCondition("process_instance_id in ?", param.ProcessInstanceIds))
	}
	if len(param.ExcludeOrderStates) > 0 {
		conditions = append(conditions, repository.NewAndCondition("order_state not in ?", param.ExcludeOrderStates))
	}
	if param.IsOrderFlowNoNotEmpty {
		conditions = append(conditions, repository.NewAndCondition("order_flow_no <> '' and order_flow_no is not null "))
	}
	if param.ProcessStatus != "" {
		conditions = append(conditions, repository.NewAndCondition("process_status = ?", param.ProcessStatus))
	}
	if param.SerialNo != "" {
		conditions = append(conditions, repository.NewAndCondition("serial_no like ?", util.GetLikeString(param.SerialNo)))
	}
	if param.BusinessId != "" {
		conditions = append(conditions, repository.NewAndCondition("process_business_id like ?", util.GetLikeString(param.BusinessId)))
	}
	if param.PayAccount != "" {
		conditions = append(conditions, repository.NewAndCondition("pay_account_name like ?", util.GetLikeString(param.PayAccount)))
	}
	if param.RecAccount != "" {
		conditions = append(conditions, repository.NewAndCondition("rec_account_name like ?", util.GetLikeString(param.RecAccount)))
	}
	if len(param.PayAmount) > 0 {
		conditions = append(conditions, repository.NewAndCondition("pay_amount >= ? and pay_amount <= ?", param.PayAmount[0], param.PayAmount[1]))
	}
	if param.OriginatorUser != "" {
		conditions = append(conditions, repository.NewAndCondition("originator_user_name like ?", util.GetLikeString(param.OriginatorUser)))
	}
	if param.CommentUser != "" {
		conditions = append(conditions, repository.NewAndCondition("comment_user_name like ?", util.GetLikeString(param.CommentUser)))
	}
	if param.TotalStatus != "" {
		conditions = append(conditions, repository.NewAndCondition("process_status = ?", param.TotalStatus))
	}
	if param.OrderState != "" {
		conditions = append(conditions, repository.NewAndCondition("order_state = ?", param.OrderState))
	}
	if len(param.CreateTime) > 1 {
		if param.CreateTime[0] != "" {
			conditions = append(conditions, repository.NewAndCondition("created_at >= ?", param.CreateTime[0]))
		}
		if param.CreateTime[1] != "" {
			conditions = append(conditions, repository.NewAndCondition("created_at < date_add(?, interval 1 day)", param.CreateTime[1]))
		}
	}
	if param.PayAmountMin > 0 {
		conditions = append(conditions, repository.NewAndCondition("pay_amount >= ?", param.PayAmountMin))
	}
	if param.PayAmountMax > 0 {
		conditions = append(conditions, repository.NewAndCondition("pay_amount <= ?", param.PayAmountMax))
	}

	if param.Title != "" {
		conditions = append(conditions, repository.NewAndCondition("title like ?", util.GetLikeString(param.Title)))
	}
	if param.OrganizationId != 0 {
		conditions = append(conditions, repository.NewAndCondition("organization_id = ?", param.OrganizationId))
	}
	if param.PayAccountType != "" {
		conditions = append(conditions, repository.NewAndCondition("pay_account_type = ?", param.PayAccountType))
	}
	return conditions
}

type BankTransferReceiptRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankTransferReceiptDBData) (*BankTransferReceiptDBData, error)
	Count(context.Context, *BankTransferReceiptDBDataParam) (int64, error)
	List(context.Context, string, int32, int32, *BankTransferReceiptDBDataParam) (*[]BankTransferReceiptDBData, int64, error)
}

type bankTransferReceiptRepo struct {
	repository.BaseRepo
}

func (r *bankTransferReceiptRepo) Get(ctx context.Context, param *BankTransferReceiptDBData) (*BankTransferReceiptDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankTransferReceiptDBData), handler.HandleError(err)
}

func (r *bankTransferReceiptRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankTransferReceiptDBDataParam) (*[]BankTransferReceiptDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankTransferReceiptDBData), count, handler.HandleError(err)
}

func (r *bankTransferReceiptRepo) Count(ctx context.Context, param *BankTransferReceiptDBDataParam) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
