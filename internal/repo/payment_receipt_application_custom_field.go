package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
)

type PaymentReceiptApplicationCustomFieldDBData struct {
	repository.BaseProcessCustomFieldDBData
	PaymentReceiptId int64
}

func (*PaymentReceiptApplicationCustomFieldDBData) TableName() string {
	return "payment_receipt_application_custom_field"
}

func (param *PaymentReceiptApplicationCustomFieldDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *PaymentReceiptApplicationCustomFieldDBData) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	if param.PaymentReceiptId > 0 {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt_id = ?", param.PaymentReceiptId))
	}
	return conditions
}

type PaymentReceiptApplicationCustomFieldRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *PaymentReceiptApplicationCustomFieldDBData) (*PaymentReceiptApplicationCustomFieldDBData, error)
	Count(context.Context, *PaymentReceiptApplicationCustomFieldDBData) (int64, error)
	List(context.Context, string, int32, int32, *PaymentReceiptApplicationCustomFieldDBData) (*[]PaymentReceiptApplicationCustomFieldDBData, int64, error)
}

type paymentReceiptApplicationCustomFieldRepo struct {
	repository.BaseRepo
}

func (r *paymentReceiptApplicationCustomFieldRepo) Get(ctx context.Context, param *PaymentReceiptApplicationCustomFieldDBData) (*PaymentReceiptApplicationCustomFieldDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*PaymentReceiptApplicationCustomFieldDBData), handler.HandleError(err)
}

func (r *paymentReceiptApplicationCustomFieldRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *PaymentReceiptApplicationCustomFieldDBData) (*[]PaymentReceiptApplicationCustomFieldDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]PaymentReceiptApplicationCustomFieldDBData), count, handler.HandleError(err)
}

func (r *paymentReceiptApplicationCustomFieldRepo) Count(ctx context.Context, param *PaymentReceiptApplicationCustomFieldDBData) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}
