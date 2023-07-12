package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
)

type BankTransactionDetailProcessInstanceDBData struct {
	repository.BaseDBData
	BankTransactionDetailId int64
	ExternalId              string
}

func (*BankTransactionDetailProcessInstanceDBData) TableName() string {
	return "bank_transaction_detail_process_instance"
}

func (param *BankTransactionDetailProcessInstanceDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankTransactionDetailProcessInstanceDBData) listConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

type BankTransactionDetailProcessInstanceRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankTransactionDetailProcessInstanceDBData) (*BankTransactionDetailProcessInstanceDBData, error)
	Count(context.Context, *BankTransactionDetailProcessInstanceDBData) (int64, error)
	List(context.Context, string, int32, int32, *BankTransactionDetailProcessInstanceDBData) (*[]BankTransactionDetailProcessInstanceDBData, int64, error)
}

type bankTransactionDetailProcessInstanceRepo struct {
	repository.BaseRepo
}

func (r *bankTransactionDetailProcessInstanceRepo) Get(ctx context.Context, param *BankTransactionDetailProcessInstanceDBData) (*BankTransactionDetailProcessInstanceDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankTransactionDetailProcessInstanceDBData), handler.HandleError(err)
}

func (r *bankTransactionDetailProcessInstanceRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankTransactionDetailProcessInstanceDBData) (*[]BankTransactionDetailProcessInstanceDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankTransactionDetailProcessInstanceDBData), count, handler.HandleError(err)
}

func (r *bankTransactionDetailProcessInstanceRepo) Count(ctx context.Context, param *BankTransactionDetailProcessInstanceDBData) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
