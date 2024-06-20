package main

import (
	"context"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service"
	api "gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
)

// BankImpl implements the last service interface defined in the IDL.
type BankImpl struct {
	bankService           service.BankService
	paymentReceiptService service.PaymentReceiptService
}

// ListBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) ListBankTransferReceipt(ctx context.Context, req *api.ListBankTransferReceiptRequest) (resp *api.ListBankTransferReceiptResponse, err error) {
	return s.bankService.ListBankTransferReceipt(ctx, req)
}

// GetBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) GetBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (resp *api.BankTransferReceiptData, err error) {
	return s.bankService.GetBankTransferReceipt(ctx, req)
}

// AddBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) AddBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (resp int64, err error) {
	return s.bankService.AddBankTransferReceipt(ctx, req)
}

// DeleteBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) DeleteBankTransferReceipt(ctx context.Context, req int64) (err error) {
	return s.bankService.DeleteBankTransferReceipt(ctx, req)
}

// CountBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) CountBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (resp int64, err error) {
	return s.bankService.CountBankTransferReceipt(ctx, req)
}

// ListBankTransactionDetail implements the BankImpl interface.
func (s *BankImpl) ListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest) (resp *api.ListBankTransactionDetailResponse, err error) {
	return s.bankService.ListBankTransactionDetail(ctx, req)
}

// GetBankTransactionDetail implements the BankImpl interface.
func (s *BankImpl) GetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData) (resp *api.BankTransactionDetailData, err error) {
	return s.bankService.GetBankTransactionDetail(ctx, req)
}

// ListBankTransactionDetailProcessInstance implements the BankImpl interface.
func (s *BankImpl) ListBankTransactionDetailProcessInstance(ctx context.Context, id int64) (resp []*api.BankTransactionDetailProcessInstanceData, err error) {
	return s.bankService.ListBankTransactionDetailProcessInstance(ctx, id)
}

// GetBankCodeInfo implements the BankImpl interface.
func (s *BankImpl) GetBankCodeInfo(ctx context.Context, code string) (resp *api.BankCodeData, err error) {
	return s.bankService.GetBankCodeInfo(ctx, code)
}

// QueryBankCardInfo implements the BankImpl interface.
func (s *BankImpl) QueryBankCardInfo(ctx context.Context, cardNo string) (resp *api.QueryBankCardInfoResponse, err error) {
	return s.bankService.QueryBankCardInfo(ctx, cardNo)
}

// ListBankCode implements the BankImpl interface.
func (s *BankImpl) ListBankCode(ctx context.Context, req *api.ListBankCodeRequest) (resp *api.ListBankCodeResponse, err error) {
	return s.bankService.ListBankCode(ctx, req)
}

// AddBankCode implements the BankImpl interface.
func (s *BankImpl) AddBankCode(ctx context.Context, req *api.AddBankCodeRequest) (err error) {
	return s.bankService.AddBankCode(ctx, req)
}

// EditBankCode implements the BankImpl interface.
func (s *BankImpl) EditBankCode(ctx context.Context, req *api.BankCodeData) (err error) {
	return s.bankService.EditBankCode(ctx, req)
}

// DeleteBankCode implements the BankImpl interface.
func (s *BankImpl) DeleteBankCode(ctx context.Context, id int64) (err error) {
	return s.bankService.DeleteBankCode(ctx, id)
}

// UpdateBankTransactionRecDetail implements the BankImpl interface.
func (s *BankImpl) UpdateBankTransactionRecDetail(ctx context.Context, req *api.BankTransactionRecDetailData) (err error) {
	return s.bankService.UpdateBankTransactionRecDetail(ctx, req)
}

// GetBankCode implements the BankImpl interface.
func (s *BankImpl) GetBankCode(ctx context.Context, req *api.BankCodeData) (resp *api.BankCodeData, err error) {
	return s.bankService.GetBankCode(ctx, req)
}

// HandleTransferReceiptResult_ implements the BankImpl interface.
func (s *BankImpl) HandleTransferReceiptResult_(ctx context.Context, id int64) (err error) {
	return s.bankService.HandleTransferReceiptResult(ctx, id)
}

// EditBankTransferReceipt implements the BankImpl interface.
func (s *BankImpl) EditBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (err error) {
	return s.bankService.EditBankTransferReceipt(ctx, req)
}

// ConfirmTransaction implements the BankImpl interface.
func (s *BankImpl) ConfirmTransaction(ctx context.Context, req *api.BankTransferReceiptData) (err error) {
	return s.bankService.ConfirmTransaction(ctx, req)
}

// SyncTransferReceipt implements the BankImpl interface.
func (s *BankImpl) SyncTransferReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error) {
	return s.bankService.SyncTransferReceipt(ctx, taskId, param, organizationId)
}

// SyncTransactionDetail implements the BankImpl interface.
func (s *BankImpl) SyncTransactionDetail(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error) {
	return s.bankService.SyncTransactionDetail(ctx, taskId, param, organizationId)
}

// HandleTransactionDetail implements the BankImpl interface.
func (s *BankImpl) HandleTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) (err error) {
	return s.bankService.HandleTransactionDetail(ctx, beginDate, endDate, organizationId)
}

// HandleSyncTransferReceipt implements the BankImpl interface.
func (s *BankImpl) HandleSyncTransferReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) (err error) {
	return s.bankService.HandleSyncTransferReceipt(ctx, beginDate, endDate, organizationId)
}

// GetCashFlowMonthChartData implements the BankImpl interface.
func (s *BankImpl) GetCashFlowMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (resp *api.ChartData, err error) {
	return s.bankService.GetCashFlowMonthChartData(ctx, req)
}

// GetBalanceMonthChartData implements the BankImpl interface.
func (s *BankImpl) GetBalanceMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (resp *api.ChartData, err error) {
	return s.bankService.GetBalanceMonthChartData(ctx, req)
}

// DashboardData implements the BankImpl interface.
func (s *BankImpl) DashboardData(ctx context.Context, organizationId int64) (resp *api.DashboardData, err error) {
	return s.bankService.DashboardData(ctx, organizationId)
}

// CreateTransactionDetailProcessInstance implements the BankImpl interface.
func (s *BankImpl) CreateTransactionDetailProcessInstance(ctx context.Context, id int64) (err error) {
	return s.bankService.CreateTransactionDetailProcessInstance(ctx, id)
}

// QueryAccountBalance implements the BankImpl interface.
func (s *BankImpl) QueryAccountBalance(ctx context.Context, req *api.QueryAccountBalanceRequest) (resp *api.QueryAccountBalanceResponse, err error) {
	return s.bankService.QueryAccountBalance(ctx, req)
}

// ImportBankBusinessPayrollData implements the BankImpl interface.
func (s *BankImpl) ImportBankBusinessPayrollData(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error) {
	return s.bankService.ImportBankBusinessPayrollData(ctx, taskId, param, organizationId)
}

// ListBankBusinessPayroll implements the BankImpl interface.
func (s *BankImpl) ListBankBusinessPayroll(ctx context.Context, req *api.ListBusinessPayrollRequest) (resp *api.ListBusinessPayrollResponse, err error) {
	return s.bankService.ListBankBusinessPayroll(ctx, req)
}

// ListBankBusinessPayrollDetail implements the BankImpl interface.
func (s *BankImpl) ListBankBusinessPayrollDetail(ctx context.Context, req *api.ListBusinessPayrollDetailRequest) (resp *api.ListBusinessPayrollDetailResponse, err error) {
	return s.bankService.ListBankBusinessPayrollDetail(ctx, req)
}

// SyncBankBusinessPayrollDetail implements the BankImpl interface.
func (s *BankImpl) SyncBankBusinessPayrollDetail(ctx context.Context, req *api.SyncBusinessPayrollResultRequest) (resp *api.SyncBusinessPayrollResultResponse, err error) {
	return s.bankService.SyncBankBusinessPayrollDetail(ctx, req)
}

// HandleTransactionDetailReceipt implements the BankImpl interface.
func (s *BankImpl) HandleTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) (err error) {
	return s.bankService.HandleTransactionDetailReceipt(ctx, beginDate, endDate, organizationId)
}

// CreateVirtualAccount implements the BankImpl interface.
func (s *BankImpl) CreateVirtualAccount(ctx context.Context, req *api.CreateVirtualAccountRequest) (resp *api.CreateVirtualAccountResponse, err error) {
	return s.bankService.CreateVirtualAccount(ctx, req)
}

// SyncVirtualAccountBalance implements the BankImpl interface.
func (s *BankImpl) SyncVirtualAccountBalance(ctx context.Context) (err error) {
	s.bankService.PinganSyncVirtualAccountBalance(ctx)
	return s.bankService.SyncVirtualAccountBalance(ctx)
}

// SpdBankVirtualAccountTranscation implements the BankImpl interface.
func (s *BankImpl) SpdBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (resp *api.BankVirtualAccountTranscationResponse, err error) {
	return s.bankService.VirtualAccountTranscation(ctx, organizationId, req)
}

// ListPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) ListPaymentReceipt(ctx context.Context, req *api.ListPaymentReceiptRequest) (resp *api.ListPaymentReceiptResponse, err error) {
	return s.paymentReceiptService.ListPaymentReceipt(ctx, req)
}

// AddPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) AddPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData) (err error) {
	return s.paymentReceiptService.Create(ctx, req)
}

// ApprovePaymentReceipt implements the BankImpl interface.
func (s *BankImpl) ApprovePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData) (err error) {
	return s.paymentReceiptService.Approve(ctx, id, req)
}

// RefusePaymentReceipt implements the BankImpl interface.
func (s *BankImpl) RefusePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string) (err error) {
	return s.paymentReceiptService.Refuse(ctx, id, req, remark)
}

// PaymentReceiptRun implements the BankImpl interface.
func (s *BankImpl) PaymentReceiptRun(ctx context.Context, id int64) (err error) {
	return s.paymentReceiptService.PaymentReceiptRun(ctx, id)
}

// HandleSyncPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) HandleSyncPaymentReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) (err error) {
	return s.paymentReceiptService.HandleSyncPaymentReceipt(ctx, beginDate, endDate, organizationId)
}

// SyncPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) SyncPaymentReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error) {
	return s.paymentReceiptService.SyncPaymentReceipt(ctx, taskId, param, organizationId)
}

// GetPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) GetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	return s.paymentReceiptService.GetPaymentReceipt(ctx, id)
}

// TransmitPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) TransmitPaymentReceipt(ctx context.Context, processInstanceId int64, transmitUserId int64) (err error) {
	return s.paymentReceiptService.Transmit(ctx, processInstanceId, transmitUserId, nil)
}

// EditBankTransactionDetailExtField implements the BankImpl interface.
func (s *BankImpl) EditBankTransactionDetailExtField(ctx context.Context, req *api.BankTransactionDetailData) (err error) {
	return s.bankService.EditBankTransactionDetailExtField(ctx, req)
}

// SendBackPaymentApplication implements the BankImpl interface.
func (s *BankImpl) SendBackPaymentApplication(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string) (err error) {
	return s.paymentReceiptService.SendBack(ctx, id, req, remark)
}

// SystemRefusePaymentReceipt implements the BankImpl interface.
func (s *BankImpl) SystemRefusePaymentReceipt(ctx context.Context, id int64) (err error) {
	return s.paymentReceiptService.PaymentReceiptSystemRefuse(ctx, id)
}

// SystemApprovePaymentReceipt implements the BankImpl interface.
func (s *BankImpl) SystemApprovePaymentReceipt(ctx context.Context, id int64) (err error) {
	return s.paymentReceiptService.PaymentReceiptSystemApprove(ctx, id)
}

// QueryVirtualAccountBalance implements the BankImpl interface.
func (s *BankImpl) QueryVirtualAccountBalance(ctx context.Context, organizationId int64, bankType string) (resp *api.VirtualAccountBalanceData, err error) {
	return s.bankService.QueryVirtualAccountBalance(ctx, organizationId, bankType)
}

// WithDrawPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) WithDrawPaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData) (err error) {
	// TODO: Your code here...
	return s.paymentReceiptService.WithDraw(ctx, id, req)
}

// CommentPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) CommentPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData) (err error) {
	// TODO: Your code here...
	return s.paymentReceiptService.Comment(ctx, req)
}

// AddTagPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) AddTagPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData) (err error) {
	// TODO: Your code here...
	return s.paymentReceiptService.AddTag(ctx, req.Id, req.ProcessAddTagItemVO, req)
}

// SimpleGetBankTransactionDetail implements the BankImpl interface.
func (s *BankImpl) SimpleGetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData) (resp *api.BankTransactionDetailData, err error) {
	// TODO: Your code here...
	return s.bankService.SimpleGetBankTransactionDetail(ctx, req)
}

// SimpleGetPaymentReceipt implements the BankImpl interface.
func (s *BankImpl) SimpleGetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	// TODO: Your code here...
	return s.paymentReceiptService.SimpleGetPaymentReceipt(ctx, id)
}

// SimpleGetPaymentReceiptByProcessInstanceId implements the BankImpl interface.
func (s *BankImpl) SimpleGetPaymentReceiptByProcessInstanceId(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	// TODO: Your code here...
	return s.paymentReceiptService.SimpleGetPaymentReceiptByProcessInstanceId(ctx, id)
}

// SimpleListBankTransactionDetail implements the BankImpl interface.
func (s *BankImpl) SimpleListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest) (resp *api.ListBankTransactionDetailResponse, err error) {
	return s.bankService.SimpleListBankTransactionDetail(ctx, req)
}

// PinganBankAccountSignatureApply implements the BankImpl interface.
func (s *BankImpl) PinganBankAccountSignatureApply(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (resp *api.PinganUserAcctSignatureApplyResponse, err error) {
	return s.bankService.PinganBankAccountSignatureApply(ctx, req)
}

// PinganBankAccountSignatureQuery implements the BankImpl interface.
func (s *BankImpl) PinganBankAccountSignatureQuery(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (resp *api.PinganUserAcctSignatureApplyResponse, err error) {
	return s.bankService.PinganBankAccountSignatureQuery(ctx, req)
}

// IcbcBankAccountSignatureQuery implements the BankImpl interface.
func (s *BankImpl) IcbcBankAccountSignatureQuery(ctx context.Context, req *api.IcbcBankAccountSignatureRequest) (resp *api.IcbcBankAccountSignatureQueryResponse, err error) {
	return s.bankService.IcbcBankAccountSignatureQuery(ctx, req)
}

// IcbcBankListTransactionDetail implements the BankImpl interface.
func (s *BankImpl) IcbcBankListTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) (err error) {
	return s.bankService.IcbcBankListTransactionDetail(ctx, beginDate, endDate, organizationId)
}

// GetBankTransactionReceipt implements the BankImpl interface.
func (s *BankImpl) GetBankTransactionReceipt(ctx context.Context, id int64) (err error) {
	return s.bankService.GetBankTransactionReceipt(ctx, id)
}
