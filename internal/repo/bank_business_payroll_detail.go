package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
)

type BankBusinessPayrollDetailDBData struct {
	repository.BaseDBData
	BatchId        int64
	SerialNo       string  //序列号
	RecAccountName string  //收款人名称
	RecAccountNo   string  //收款人账号
	Amount         float64 //金额
	Month          string  //代发月份
	Remark         string  //用途
	Num            string  //工号
	OrderState     string  //交易状态
	ErrorCode      string  //返回错误码
	ErrorMessage   string  //返回错误信息
}

func (*BankBusinessPayrollDetailDBData) TableName() string {
	return "bank_business_payroll_detail"
}

func (param *BankBusinessPayrollDetailDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankBusinessPayrollDetailDBData) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	conditions = append(conditions, repository.NewAndCondition("deleted = ?", "0"))
	if param.RecAccountNo != "" {
		conditions = append(conditions, repository.NewAndCondition("rec_account_no like ?", util.GetLikeString(param.RecAccountNo)))
	}
	if param.RecAccountName != "" {
		conditions = append(conditions, repository.NewAndCondition("rec_account_name like ?", util.GetLikeString(param.RecAccountName)))
	}
	if param.Month != "" {
		conditions = append(conditions, repository.NewAndCondition("month = ?", param.Month))
	}
	if param.BatchId != 0 {
		conditions = append(conditions, repository.NewAndCondition("batch_id = ?", param.BatchId))
	}
	if param.Num != "" {
		conditions = append(conditions, repository.NewAndCondition("num like ?", util.GetLikeString(param.Num)))
	}
	return conditions
}

type BankBusinessPayrollDetailRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankBusinessPayrollDetailDBData) (*BankBusinessPayrollDetailDBData, error)
	Count(context.Context, *BankBusinessPayrollDetailDBData) (int64, error)
	List(context.Context, string, int32, int32, *BankBusinessPayrollDetailDBData) (*[]BankBusinessPayrollDetailDBData, int64, error)
}

type bankBusinessPayrollDetailRepo struct {
	repository.BaseRepo
}

func (r *bankBusinessPayrollDetailRepo) Get(ctx context.Context, param *BankBusinessPayrollDetailDBData) (*BankBusinessPayrollDetailDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankBusinessPayrollDetailDBData), handler.HandleError(err)
}

func (r *bankBusinessPayrollDetailRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankBusinessPayrollDetailDBData) (*[]BankBusinessPayrollDetailDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankBusinessPayrollDetailDBData), count, handler.HandleError(err)
}

func (r *bankBusinessPayrollDetailRepo) Count(ctx context.Context, param *BankBusinessPayrollDetailDBData) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
