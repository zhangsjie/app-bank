package sub_process

import (
	"context"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service/stru"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	api2 "gitlab.yoyiit.com/youyi/app-base/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-base/kitex_gen/api/base"
	"gitlab.yoyiit.com/youyi/app-base/process"
	"gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api/oa"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"strconv"
	"time"
)

var _ process.SubProcess = new(PaymentReceiptSubProcess)

type PaymentReceiptSubProcess struct {
	paymentReceiptRepo                       repo.PaymentReceiptRepo
	oaClient                                 oa.Client
	baseClient                               base.Client
	paymentReceiptApplicationCustomFieldRepo repo.PaymentReceiptApplicationCustomFieldRepo
}

func (p *PaymentReceiptSubProcess) SubmitBefore(ctx context.Context, id int64, param interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PaymentReceiptSubProcess) Cancel(ctx context.Context, processInstanceId int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *PaymentReceiptSubProcess) Transmit(ctx context.Context, param interface{}) error {
	return nil
}

func (p *PaymentReceiptSubProcess) WithDraw(ctx context.Context, id int64, param interface{}) error {
	req := param.(*api2.ProcessResultData)
	//跨节点删除
	if req.ExtOne == "888" {
		handler.HandleError(p.paymentReceiptRepo.DeleteById(ctx, id))
	}
	return nil
}

func (p *PaymentReceiptSubProcess) CreateAfter(ctx context.Context, id int64, param interface{}) error {
	req := param.(*api2.ProcessResultData)
	err := p.paymentReceiptRepo.UpdateSelectedFieldsByIdWithoutPermission(ctx, id, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: req.ProcessInstanceId,
		},
		ProcessName:            req.ProcessName,
		ProcessCodes:           req.ProcessCodes,
		ProcessCurrentUserName: req.ProcessCurrentUserName,
		ProcessCurrentUserId:   req.ProcessCurrentUserId,
		ProcessStatus:          req.ProcessStatus,
	}, []string{"ProcessInstanceId", "ProcessName", "ProcessCodes", "ProcessCurrentUserName", "ProcessCurrentUserId", "ProcessStatus"})
	if err != nil {
		return handler.HandleError(err)
	}
	return nil
}

func (p *PaymentReceiptSubProcess) Edit(ctx context.Context, id int64, param interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PaymentReceiptSubProcess) SyncProcessInfo(ctx context.Context, processInstanceId int64, param interface{}) error {
	result := param.(*api2.ProcessResultData)
	dbData, err := p.paymentReceiptRepo.Get(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: processInstanceId,
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if dbData == nil || dbData.Id == 0 {
		return nil
	}
	err = p.paymentReceiptRepo.UpdateSelectedFieldsByIdWithoutPermission(ctx, dbData.Id, &repo.PaymentReceiptDBData{
		ProcessStatus: result.ProcessStatus,
	}, []string{"ProcessStatus"})
	return handler.HandleError(err)
}

func (p *PaymentReceiptSubProcess) ReLaunch(ctx context.Context, id int64, remark string) error {
	//TODO implement me
	return nil
}

func (p *PaymentReceiptSubProcess) Create(ctx context.Context, param interface{}) (int64, error) {
	req := param.(*api.PaymentReceiptData)
	data := stru.ConvertPaymentReceiptDBData(*req)
	code, err := util.SonyflakeID()
	if err != nil {
		return 0, handler.HandleError(err)
	}
	data.Code = code
	data.OrderStatus = "0"
	data.RefundSuccess = "0"
	id, err := p.paymentReceiptRepo.Add(ctx, data)

	if req.CustomFields != nil && len(req.CustomFields) > 0 {
		customFields := make([]repo.PaymentReceiptApplicationCustomFieldDBData, len(req.CustomFields))
		for i, v := range req.CustomFields {
			var customField repo.PaymentReceiptApplicationCustomFieldDBData
			customField.ProcessCustomFieldId = v.Id
			customField.ProcessCustomFieldName = v.Name
			customField.ProcessCustomFieldValue = v.Value
			customField.ProcessCustomFieldSort = v.Sort
			customField.OrganizationId = v.OrganizationId
			customField.PaymentReceiptId = id
			customFields[i] = customField
		}
		_, err = p.paymentReceiptApplicationCustomFieldRepo.BatchAdd(ctx, customFields)
		if err != nil {
			return 0, handler.HandleError(err)
		}
	}
	return id, handler.HandleError(err)
}

func (p *PaymentReceiptSubProcess) UpdateProcessInstanceId(ctx context.Context, id int64, param interface{}) error {
	result := param.(*api2.ProcessResultData)
	return p.paymentReceiptRepo.UpdateByIdWithoutPermission(ctx, id, &repo.PaymentReceiptDBData{
		ProcessCurrentUserName: result.ProcessCurrentUserName,
		ProcessCurrentUserId:   result.ProcessCurrentUserId,
		ProcessStatus:          result.ProcessStatus,
	})
}

func (p *PaymentReceiptSubProcess) Approve(ctx context.Context, id int64, param interface{}) error {
	req := param.(*api.PaymentReceiptData)
	return handler.HandleError(p.paymentReceiptRepo.UpdateById(ctx, id, stru.ConvertPaymentReceiptDBData(*req)))
}

func (p *PaymentReceiptSubProcess) SystemApprove(ctx context.Context, id int64, param interface{}) error {
	req := param.(*api.PaymentReceiptData)
	return handler.HandleError(p.paymentReceiptRepo.UpdateById(ctx, id, stru.ConvertPaymentReceiptDBData(*req)))
}

func (p *PaymentReceiptSubProcess) Refuse(ctx context.Context, id int64, param interface{}) error {
	return nil
}

func (p *PaymentReceiptSubProcess) SystemRefuse(ctx context.Context, id int64, param interface{}) error {
	return nil
}

func (p *PaymentReceiptSubProcess) Delete(ctx context.Context, id int64) error {
	return handler.HandleError(p.paymentReceiptRepo.DeleteById(ctx, id))
}

func (p *PaymentReceiptSubProcess) ProcessConfigCode() string {
	return "paymentApplication"
}

func (p *PaymentReceiptSubProcess) ProcessNodeStep() int32 {
	return 2
}

func (p *PaymentReceiptSubProcess) MiniprogramUrl(ctx context.Context, id int64) string {
	paymentReceipt, err := p.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil {
		return ""
	}
	if paymentReceipt.Type == "1" {
		paymentApplication, err := p.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("page/paymentDetail/index?id=%d", paymentApplication.Id)
	}
	if paymentReceipt.Type == "2" {
		reimburseApplication, err := p.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("reimburseApplication/detail/index?id=%d", reimburseApplication.Id)
	}
	return ""
}

func (p *PaymentReceiptSubProcess) PcUrl(ctx context.Context, id int64) string {
	paymentReceipt, err := p.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil {
		return ""
	}
	if paymentReceipt.Type == "1" {
		paymentApplication, err := p.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("https://apply.%s/oa/payment-application/%d", config.GetString("dingtalk.pc.domain", ""), paymentApplication.Id)
	}
	if paymentReceipt.Type == "2" {
		reimburseApplication, err := p.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("https://apply.%s/oa/reimburse-application/%d", config.GetString("dingtalk.pc.domain", ""), reimburseApplication.Id)
	}
	return ""
}

func (p *PaymentReceiptSubProcess) TodoFields(ctx context.Context, id int64) [][]string {
	paymentReceipt, err := p.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil {
		return nil
	}

	if paymentReceipt.Type == "1" {
		paymentApplication, err := p.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return nil
		}
		purposeStr := []rune(paymentApplication.PaymentReason)
		if len(purposeStr) > 20 {
			purposeStr = append(purposeStr[:20], []rune("...")...)
		}
		result := [][]string{
			{"待办生成时间", time.Now().Format("2006-01-02 15:04:05")},
			{"付款金额", strconv.FormatFloat(paymentApplication.PayAmount, 'f', 2, 64)},
			{"收款方户名", paymentApplication.ReceiveAccountName},
			{"付款事由", string(purposeStr)},
		}
		return result
	}
	if paymentReceipt.Type == "2" {
		reimburseApplication, err := p.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, id)
		if err != nil {
			return nil
		}
		purposeStr := []rune(reimburseApplication.ReimburseInstruction)
		if len(purposeStr) > 20 {
			purposeStr = append(purposeStr[:20], []rune("...")...)
		}
		result := [][]string{
			{"待办生成时间", time.Now().Format("2006-01-02 15:04:05")},
			{"报销金额", strconv.FormatFloat(reimburseApplication.PayAmount, 'f', 2, 64)},
			{"收款方户名", reimburseApplication.ReceiveAccountName},
			{"报销说明", string(purposeStr)},
		}
		return result
	}
	return nil
}

func (p *PaymentReceiptSubProcess) SendBack(ctx context.Context, id int64, param interface{}) error {
	paymentReceipt, err := p.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			BaseDBData: repository.BaseDBData{
				BaseCommonDBData: repository.BaseCommonDBData{
					Id: id,
				},
			},
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	processInstance, err := p.baseClient.QueryProcessInstance(ctx, []int64{paymentReceipt.ProcessInstanceId})
	if err != nil {
		return handler.HandleError(err)
	}
	if processInstance == nil || len(processInstance) == 0 {
		return handler.HandleNewError("process instance not exists")
	}
	if processInstance[paymentReceipt.ProcessInstanceId].CurrentProcessNodeStep == 1 {
		return handler.HandleError(p.paymentReceiptRepo.DeleteById(ctx, id))
	}
	return nil
}
