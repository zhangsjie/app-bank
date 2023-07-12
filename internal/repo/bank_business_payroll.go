package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
)

type BankBusinessPayrollDBData struct {
	repository.BaseDBData
	Name            string
	PayAccountName  string
	PayAccountNo    string
	Month           string
	Remark          string
	Count           int64
	TotalMoney      float64
	Status          int64
	CreatedUserName string
	FileUrl         string
	BatchNo         string
	State           string
	Msg             string
	OrderSubmitTime string
}

type BankBusinessPayrollDBDataParam struct {
	BankBusinessPayrollDBData
	CreateTimeRange []string
}

func (*BankBusinessPayrollDBData) TableName() string {
	return "bank_business_payroll"
}

func (param *BankBusinessPayrollDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankBusinessPayrollDBDataParam) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	conditions = append(conditions, repository.NewAndCondition("deleted = ?", "0"))
	if len(param.CreateTimeRange) == 2 {
		conditions = append(conditions, repository.NewAndCondition("created_at >= ? and created_at <= ?", param.CreateTimeRange[0], param.CreateTimeRange[1]))
	}
	if param.PayAccountNo != "" {
		conditions = append(conditions, repository.NewAndCondition("pay_account_no like ?", util.GetLikeString(param.PayAccountNo)))
	}
	if param.PayAccountName != "" {
		conditions = append(conditions, repository.NewAndCondition("pay_account_name like ?", util.GetLikeString(param.PayAccountName)))
	}
	if param.Month != "" {
		conditions = append(conditions, repository.NewAndCondition("month = ?", param.Month))
	}
	if param.Status == 0 || param.Status == 1 || param.Status == 2 {
		conditions = append(conditions, repository.NewAndCondition("status = ?", param.Status))
	}
	return conditions
}

type BankBusinessPayrollRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankBusinessPayrollDBData) (*BankBusinessPayrollDBData, error)
	Count(context.Context, *BankBusinessPayrollDBDataParam) (int64, error)
	List(context.Context, string, int32, int32, *BankBusinessPayrollDBDataParam) (*[]BankBusinessPayrollDBData, int64, error)
}

type bankBusinessPayrollRepo struct {
	repository.BaseRepo
}

func (r *bankBusinessPayrollRepo) Get(ctx context.Context, param *BankBusinessPayrollDBData) (*BankBusinessPayrollDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankBusinessPayrollDBData), handler.HandleError(err)
}

func (r *bankBusinessPayrollRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankBusinessPayrollDBDataParam) (*[]BankBusinessPayrollDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankBusinessPayrollDBData), count, handler.HandleError(err)
}

func (r *bankBusinessPayrollRepo) Count(ctx context.Context, param *BankBusinessPayrollDBDataParam) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
