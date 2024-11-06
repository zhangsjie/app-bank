package repo

import (
	"context"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"strings"
)

type PaymentReceiptDBData struct {
	repository.BaseProcessDBData
	Code                   string
	PayAmount              float64
	PublicPrivateFlag      string
	PayAccount             string
	PayAccountName         string
	PayAccountBankName     string
	PayAccountType         string
	ReceiveAccount         string
	ReceiveAccountName     string
	ReceiveAccountBankName string
	Purpose                string
	UnionBankNo            string
	ClearBankNo            string
	InsideOutsideBankType  string
	ChargeFee              float64
	OrderStatus            string
	RetCode                string
	RetMessage             string
	OrderFlowNo            string
	Type                   string
	PaymentModeType        string
	ProcessStatus          string
	ProcessCodes           string
	ProcessName            string
	ProcessCurrentUserId   int64
	RefundSuccess          string
	ProcessCurrentUserName string
	BusType                string
	BusOrderNo             string
	ReceiptOrderNo         string
	ApplicantId            int64
	ApplicantName          string
	FillingDt              string
	DepartmentId           int64
	DepartmentName         string
	Attachments            string
	ElectronicDocument     string
	ElectronicDocumentPng  string
	PaymentReason          string
	PaymentId              int64
	TransDate              string
}

type PaymentReceiptDBParam struct {
	PaymentReceiptDBData
	ExcludeOrderStates    []string
	IsOrderFlowNoNotEmpty bool
	CreateTime            []string
	UpdateTime            []string
	CreateTimeStart       string
	CreateTimeEnd         string
	BeginTime             string
	EndTime               string
}

func (*PaymentReceiptDBData) TableName() string {
	return "payment_receipt"
}

func (PaymentReceiptDBData) ProcessCode() string {
	return "paymentApplication"
}

func (PaymentReceiptDBData) ProcessNodeStep() int32 {
	return 2
}

func (param *PaymentReceiptDBData) getConditions() []*repository.Condition {
	return []*repository.Condition{
		repository.NewAndCondition(param),
	}
}

func (param *PaymentReceiptDBParam) listConditions() []*repository.Condition {
	var conditions []*repository.Condition
	if len(param.ExcludeOrderStates) > 0 {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_status not in ?", param.ExcludeOrderStates))
	}
	if param.IsOrderFlowNoNotEmpty {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_flow_no <> '' and order_flow_no is not null "))
	}
	if param.PayAccount != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.pay_account_name like ?", util.GetLikeString(param.PayAccount)))
	}
	if len(param.CreateTime) > 1 {
		if param.CreateTime[0] != "" {
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.created_at >= ?", param.CreateTime[0]))
		}
		if param.CreateTime[1] != "" {
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.created_at < date_add(?, interval 1 day)", param.CreateTime[1]))
		}
	}
	if len(param.UpdateTime) > 1 {
		if param.UpdateTime[0] != "" {
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.updated_at >= ?", param.UpdateTime[0]))
		}
		if param.UpdateTime[1] != "" {
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.updated_at < date_add(?, interval 1 day)", param.UpdateTime[1]))
		}
	}
	if param.OrganizationId != 0 {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.organization_id = ?", param.OrganizationId))
	}
	if param.PayAccountType != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.pay_account_type = ?", param.PayAccountType))
	}
	if param.Code != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.code like ?", util.GetLikeString(param.Code)))
	}
	if param.ReceiveAccount != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.receive_account like ?", util.GetLikeString(param.ReceiveAccount)))
	}
	if param.ReceiveAccountName != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.receive_account_name like ?", util.GetLikeString(param.ReceiveAccountName)))
	}
	if param.ReceiveAccountBankName != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.receive_account_bank_name like ?", util.GetLikeString(param.ReceiveAccountBankName)))
	}
	if param.Purpose != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.purpose like ?", util.GetLikeString(param.Purpose)))
	}
	if param.OrderStatus != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_status = ?", param.OrderStatus))
	}
	if param.OrderFlowNo != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_flow_no like ?", util.GetLikeString(param.OrderFlowNo)))
	}
	if param.CreateTimeStart != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.created_at >= ?", param.CreateTimeStart))
	}
	if param.CreateTimeEnd != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.created_at <= ?", param.CreateTimeEnd))
	}
	if param.Type != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.type = ?", param.Type))
	}
	if param.PaymentModeType != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.payment_mode_type = ?", param.PaymentModeType))
	}
	if param.ProcessName != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.process_name like ?", util.GetLikeString(param.ProcessName)))
	}
	if param.ProcessCodes != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.process_codes like ?", util.GetLikeString(param.ProcessCodes)))
	}
	if param.ProcessStatus != "" {
		if param.ProcessStatus == "1" {
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_status = '90'"))
		}
		if param.ProcessStatus == "5" {
			param.ProcessStatus = "1"
			conditions = append(conditions, repository.NewAndCondition("payment_receipt.order_status != '90'"))
		}
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.process_instance_id in (select id from base.process_instance where base.process_instance.status = ? )", param.ProcessStatus))
	}
	if param.RefundSuccess != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.refund_success = ?", param.RefundSuccess))
	}
	if param.ProcessCurrentUserId > 0 {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.process_current_user_id = ?", param.ProcessCurrentUserId))
	}
	if param.ProcessCurrentUserName != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.process_instance_id in (select id from base.process_instance where base.process_instance.current_process_user_id > 0 and base.process_instance.current_process_user_nickname like ? )", util.GetLikeString(param.ProcessCurrentUserName)))
		conditions = append(conditions, repository.NewOrCondition("payment_receipt.process_instance_id in (SELECT DISTINCT s.id FROM base.process_instance s,base.process_instance_node n,base.process_instance_item i,base.process_instance_item_multier m WHERE s.id = n.process_instance_id AND n.id = i.process_instance_node_id AND i.id = m.process_instance_item_id AND m.process_user_nickname LIKE ? AND m.`status` = '0' AND i.`status` = '0' AND s.`status` = '0')", util.GetLikeString(param.ProcessCurrentUserName)))
	}
	if param.BusType != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.bus_type = ?", param.BusType))
	}
	if param.BusOrderNo != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.bus_order_no like ?", util.GetLikeString(param.BusOrderNo)))
	}
	if param.ReceiptOrderNo != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.receipt_order_no like ?", util.GetLikeString(param.ReceiptOrderNo)))
	}
	if param.ApplicantName != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.applicant_name like ?", util.GetLikeString(param.ApplicantName)))
	}
	if param.BeginTime != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.filling_dt >= ?", param.BeginTime))
	}
	if param.EndTime != "" {
		conditions = append(conditions, repository.NewAndCondition("payment_receipt.filling_dt <= ?", param.EndTime))
	}
	return conditions
}

type PaymentReceiptRepo interface {
	repository.BaseCommonRepo
	Get(context.Context, *PaymentReceiptDBData) (*PaymentReceiptDBData, error)
	SimpleGet(context.Context, *PaymentReceiptDBData) (*PaymentReceiptDBData, error)
	GetWithoutPermission(context.Context, *PaymentReceiptDBData) (*PaymentReceiptDBData, error)
	Count(context.Context, *PaymentReceiptDBParam) (int64, error)
	List(context.Context, string, int32, int32, *PaymentReceiptDBParam) (*[]PaymentReceiptDBData, *[]PaymentReceiptDBData, int64, error)
	ListAll(context.Context, string, int32, int32, *PaymentReceiptDBParam) (*[]PaymentReceiptDBData, int64, error)
}

type paymentReceiptRepo struct {
	repository.BaseRepo
}

func (r *paymentReceiptRepo) Get(ctx context.Context, param *PaymentReceiptDBData) (*PaymentReceiptDBData, error) {
	data, err := r.BaseGet(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*PaymentReceiptDBData), handler.HandleError(err)
}

func (r *paymentReceiptRepo) SimpleGet(ctx context.Context, param *PaymentReceiptDBData) (*PaymentReceiptDBData, error) {
	conditions := []string{"deleted = ?"}
	params := []interface{}{0}
	if param.Id > 0 {
		conditions = append(conditions, "id = ?")
		params = append(params, param.Id)
	}
	if param.ProcessInstanceId > 0 {
		conditions = append(conditions, "process_instance_id = ?")
		params = append(params, param.ProcessInstanceId)
	}
	var result PaymentReceiptDBData
	db := r.Db.WithContext(ctx).Where(strings.Join(conditions, " and "), params...)
	err := db.Find(&result).Error
	return &result, handler.HandleError(err)
}

func (r *paymentReceiptRepo) GetWithoutPermission(ctx context.Context, param *PaymentReceiptDBData) (*PaymentReceiptDBData, error) {
	data, err := r.BaseGetWithoutPermission(ctx, repository.NewQueryBuilder().Where(param.getConditions()))
	if data == nil {
		return nil, handler.HandleError(err)
	}
	return data.(*PaymentReceiptDBData), handler.HandleError(err)
}

func (r *paymentReceiptRepo) List(ctx context.Context, order string, pageNum, pageSize int32, param *PaymentReceiptDBParam) (*[]PaymentReceiptDBData, *[]PaymentReceiptDBData, int64, error) {
	data, count, err := r.BaseList(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, nil, count, handler.HandleError(err)
	}
	result := data.(*[]PaymentReceiptDBData)
	var canWriteData *[]PaymentReceiptDBData
	if result != nil && len(*result) > 0 {
		ids := make([]int64, len(*result))
		for i, v := range *result {
			ids[i] = v.Id
		}
		canWriteRes, err := r.ListWritePermission(ctx, ids)
		if err != nil {
			return nil, nil, 0, handler.HandleError(err)
		}
		canWriteData = canWriteRes.(*[]PaymentReceiptDBData)
	}
	return data.(*[]PaymentReceiptDBData), canWriteData, count, handler.HandleError(err)
}

func (r *paymentReceiptRepo) Count(ctx context.Context, param *PaymentReceiptDBParam) (int64, error) {
	count, err := r.BaseCount(ctx, repository.NewQueryBuilder().Where(param.listConditions()))
	return count, handler.HandleError(err)
}

func (r *paymentReceiptRepo) ListAll(ctx context.Context, order string, pageNum, pageSize int32, param *PaymentReceiptDBParam) (*[]PaymentReceiptDBData, int64, error) {
	data, count, err := r.BaseListWithoutPermission(ctx, repository.NewListQueryBuilder(order, pageNum, pageSize).Where(param.listConditions()))
	if data == nil {
		return nil, count, handler.HandleError(err)
	}
	return data.(*[]PaymentReceiptDBData), count, handler.HandleError(err)
}
