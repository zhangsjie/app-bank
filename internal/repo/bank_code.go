package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
)

// BankCodeDBData  银行编码 联行号 清算行号 关系表/*
type BankCodeDBData struct {
	repository.BaseCommonDBData
	BankName      string //银行全称
	BankAliasName string //银行别名 简称
	BankCode      string //银行编号
	UnionBankNo   string //联行号
	ClearBankNo   string //清算行号
}

func (*BankCodeDBData) TableName() string {
	return "bank_code"
}

func (param *BankCodeDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *BankCodeDBData) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	if param.BankName != "" {
		conditions = append(conditions, repository.NewAndCondition("bank_name like ?", util.GetLikeString(param.BankName)))
	}
	if param.BankAliasName != "" {
		conditions = append(conditions, repository.NewAndCondition("bank_alias_name like ?", util.GetLikeString(param.BankAliasName)))
	}
	return conditions
}

type BankCodeRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *BankCodeDBData) (*BankCodeDBData, error)
	Count(context.Context, *BankCodeDBData) (int64, error)
	List(context.Context, string, int32, int32, *BankCodeDBData) (*[]BankCodeDBData, int64, error)
}

type bankCodeRepo struct {
	repository.BaseRepo
}

func (r *bankCodeRepo) Get(ctx context.Context, param *BankCodeDBData) (*BankCodeDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*BankCodeDBData), handler.HandleError(err)
}

func (r *bankCodeRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *BankCodeDBData) (*[]BankCodeDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]BankCodeDBData), count, handler.HandleError(err)
}

func (r *bankCodeRepo) Count(ctx context.Context, param *BankCodeDBData) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
