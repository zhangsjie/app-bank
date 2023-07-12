// Code generated by Kitex v0.4.4. DO NOT EDIT.

package bank

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	api "gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	ListBankTransferReceipt(ctx context.Context, req *api.ListBankTransferReceiptRequest, callOptions ...callopt.Option) (r *api.ListBankTransferReceiptResponse, err error)
	GetBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r *api.BankTransferReceiptData, err error)
	AddBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r int64, err error)
	EditBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (err error)
	DeleteBankTransferReceipt(ctx context.Context, req int64, callOptions ...callopt.Option) (err error)
	CountBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r int64, err error)
	ConfirmTransaction(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (err error)
	HandleTransferReceiptResult_(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
	ListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest, callOptions ...callopt.Option) (r *api.ListBankTransactionDetailResponse, err error)
	GetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData, callOptions ...callopt.Option) (r *api.BankTransactionDetailData, err error)
	HandleTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error)
	CreateTransactionDetailProcessInstance(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
	EditBankTransactionDetailExtField(ctx context.Context, req *api.BankTransactionDetailData, callOptions ...callopt.Option) (err error)
	ListBankTransactionDetailProcessInstance(ctx context.Context, id int64, callOptions ...callopt.Option) (r []*api.BankTransactionDetailProcessInstanceData, err error)
	GetBankCodeInfo(ctx context.Context, code string, callOptions ...callopt.Option) (r *api.BankCodeData, err error)
	QueryBankCardInfo(ctx context.Context, cardNo string, callOptions ...callopt.Option) (r *api.QueryBankCardInfoResponse, err error)
	ListBankCode(ctx context.Context, req *api.ListBankCodeRequest, callOptions ...callopt.Option) (r *api.ListBankCodeResponse, err error)
	GetBankCode(ctx context.Context, req *api.BankCodeData, callOptions ...callopt.Option) (r *api.BankCodeData, err error)
	AddBankCode(ctx context.Context, req *api.AddBankCodeRequest, callOptions ...callopt.Option) (err error)
	EditBankCode(ctx context.Context, req *api.BankCodeData, callOptions ...callopt.Option) (err error)
	DeleteBankCode(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
	HandleSyncTransferReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error)
	UpdateBankTransactionRecDetail(ctx context.Context, req *api.BankTransactionRecDetailData, callOptions ...callopt.Option) (err error)
	SyncTransferReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error)
	SyncTransactionDetail(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error)
	DashboardData(ctx context.Context, organizationId int64, callOptions ...callopt.Option) (r *api.DashboardData, err error)
	GetCashFlowMonthChartData(ctx context.Context, req *api.MonthChartDataRequest, callOptions ...callopt.Option) (r *api.ChartData, err error)
	GetBalanceMonthChartData(ctx context.Context, req *api.MonthChartDataRequest, callOptions ...callopt.Option) (r *api.ChartData, err error)
	QueryAccountBalance(ctx context.Context, req *api.QueryAccountBalanceRequest, callOptions ...callopt.Option) (r *api.QueryAccountBalanceResponse, err error)
	ImportBankBusinessPayrollData(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error)
	ListBankBusinessPayroll(ctx context.Context, req *api.ListBusinessPayrollRequest, callOptions ...callopt.Option) (r *api.ListBusinessPayrollResponse, err error)
	ListBankBusinessPayrollDetail(ctx context.Context, req *api.ListBusinessPayrollDetailRequest, callOptions ...callopt.Option) (r *api.ListBusinessPayrollDetailResponse, err error)
	SyncBankBusinessPayrollDetail(ctx context.Context, req *api.SyncBusinessPayrollResultRequest, callOptions ...callopt.Option) (r *api.SyncBusinessPayrollResultResponse, err error)
	HandleTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error)
	CreateVirtualAccount(ctx context.Context, req *api.CreateVirtualAccountRequest, callOptions ...callopt.Option) (r *api.CreateVirtualAccountResponse, err error)
	SyncVirtualAccountBalance(ctx context.Context, callOptions ...callopt.Option) (err error)
	QueryVirtualAccountBalance(ctx context.Context, organizationId int64, bankType string, callOptions ...callopt.Option) (r *api.VirtualAccountBalanceData, err error)
	SpdBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r *api.BankVirtualAccountTranscationResponse, err error)
	ListPaymentReceipt(ctx context.Context, req *api.ListPaymentReceiptRequest, callOptions ...callopt.Option) (r *api.ListPaymentReceiptResponse, err error)
	GetPaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (r *api.PaymentReceiptData, err error)
	AddPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error)
	ApprovePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error)
	RefusePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string, callOptions ...callopt.Option) (err error)
	PaymentReceiptRun(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
	TransmitPaymentReceipt(ctx context.Context, processInstanceId int64, transmitUserId int64, callOptions ...callopt.Option) (err error)
	SendBackPaymentApplication(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string, callOptions ...callopt.Option) (err error)
	WithDrawPaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error)
	CommentPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error)
	HandleSyncPaymentReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error)
	SyncPaymentReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error)
	PinganBankAccountSignatureApply(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest, callOptions ...callopt.Option) (r *api.PinganUserAcctSignatureApplyResponse, err error)
	PinganBankAccountSignatureQuery(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest, callOptions ...callopt.Option) (r *api.PinganUserAcctSignatureApplyResponse, err error)
	SystemRefusePaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
	SystemApprovePaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kBankClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kBankClient struct {
	*kClient
}

func (p *kBankClient) ListBankTransferReceipt(ctx context.Context, req *api.ListBankTransferReceiptRequest, callOptions ...callopt.Option) (r *api.ListBankTransferReceiptResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankTransferReceipt(ctx, req)
}

func (p *kBankClient) GetBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r *api.BankTransferReceiptData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetBankTransferReceipt(ctx, req)
}

func (p *kBankClient) AddBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r int64, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddBankTransferReceipt(ctx, req)
}

func (p *kBankClient) EditBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.EditBankTransferReceipt(ctx, req)
}

func (p *kBankClient) DeleteBankTransferReceipt(ctx context.Context, req int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteBankTransferReceipt(ctx, req)
}

func (p *kBankClient) CountBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r int64, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CountBankTransferReceipt(ctx, req)
}

func (p *kBankClient) ConfirmTransaction(ctx context.Context, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ConfirmTransaction(ctx, req)
}

func (p *kBankClient) HandleTransferReceiptResult_(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleTransferReceiptResult_(ctx, id)
}

func (p *kBankClient) ListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest, callOptions ...callopt.Option) (r *api.ListBankTransactionDetailResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankTransactionDetail(ctx, req)
}

func (p *kBankClient) GetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData, callOptions ...callopt.Option) (r *api.BankTransactionDetailData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetBankTransactionDetail(ctx, req)
}

func (p *kBankClient) HandleTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleTransactionDetail(ctx, beginDate, endDate, organizationId)
}

func (p *kBankClient) CreateTransactionDetailProcessInstance(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateTransactionDetailProcessInstance(ctx, id)
}

func (p *kBankClient) EditBankTransactionDetailExtField(ctx context.Context, req *api.BankTransactionDetailData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.EditBankTransactionDetailExtField(ctx, req)
}

func (p *kBankClient) ListBankTransactionDetailProcessInstance(ctx context.Context, id int64, callOptions ...callopt.Option) (r []*api.BankTransactionDetailProcessInstanceData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankTransactionDetailProcessInstance(ctx, id)
}

func (p *kBankClient) GetBankCodeInfo(ctx context.Context, code string, callOptions ...callopt.Option) (r *api.BankCodeData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetBankCodeInfo(ctx, code)
}

func (p *kBankClient) QueryBankCardInfo(ctx context.Context, cardNo string, callOptions ...callopt.Option) (r *api.QueryBankCardInfoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.QueryBankCardInfo(ctx, cardNo)
}

func (p *kBankClient) ListBankCode(ctx context.Context, req *api.ListBankCodeRequest, callOptions ...callopt.Option) (r *api.ListBankCodeResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankCode(ctx, req)
}

func (p *kBankClient) GetBankCode(ctx context.Context, req *api.BankCodeData, callOptions ...callopt.Option) (r *api.BankCodeData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetBankCode(ctx, req)
}

func (p *kBankClient) AddBankCode(ctx context.Context, req *api.AddBankCodeRequest, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddBankCode(ctx, req)
}

func (p *kBankClient) EditBankCode(ctx context.Context, req *api.BankCodeData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.EditBankCode(ctx, req)
}

func (p *kBankClient) DeleteBankCode(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteBankCode(ctx, id)
}

func (p *kBankClient) HandleSyncTransferReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleSyncTransferReceipt(ctx, beginDate, endDate, organizationId)
}

func (p *kBankClient) UpdateBankTransactionRecDetail(ctx context.Context, req *api.BankTransactionRecDetailData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateBankTransactionRecDetail(ctx, req)
}

func (p *kBankClient) SyncTransferReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SyncTransferReceipt(ctx, taskId, param, organizationId)
}

func (p *kBankClient) SyncTransactionDetail(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SyncTransactionDetail(ctx, taskId, param, organizationId)
}

func (p *kBankClient) DashboardData(ctx context.Context, organizationId int64, callOptions ...callopt.Option) (r *api.DashboardData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DashboardData(ctx, organizationId)
}

func (p *kBankClient) GetCashFlowMonthChartData(ctx context.Context, req *api.MonthChartDataRequest, callOptions ...callopt.Option) (r *api.ChartData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetCashFlowMonthChartData(ctx, req)
}

func (p *kBankClient) GetBalanceMonthChartData(ctx context.Context, req *api.MonthChartDataRequest, callOptions ...callopt.Option) (r *api.ChartData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetBalanceMonthChartData(ctx, req)
}

func (p *kBankClient) QueryAccountBalance(ctx context.Context, req *api.QueryAccountBalanceRequest, callOptions ...callopt.Option) (r *api.QueryAccountBalanceResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.QueryAccountBalance(ctx, req)
}

func (p *kBankClient) ImportBankBusinessPayrollData(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ImportBankBusinessPayrollData(ctx, taskId, param, organizationId)
}

func (p *kBankClient) ListBankBusinessPayroll(ctx context.Context, req *api.ListBusinessPayrollRequest, callOptions ...callopt.Option) (r *api.ListBusinessPayrollResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankBusinessPayroll(ctx, req)
}

func (p *kBankClient) ListBankBusinessPayrollDetail(ctx context.Context, req *api.ListBusinessPayrollDetailRequest, callOptions ...callopt.Option) (r *api.ListBusinessPayrollDetailResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListBankBusinessPayrollDetail(ctx, req)
}

func (p *kBankClient) SyncBankBusinessPayrollDetail(ctx context.Context, req *api.SyncBusinessPayrollResultRequest, callOptions ...callopt.Option) (r *api.SyncBusinessPayrollResultResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SyncBankBusinessPayrollDetail(ctx, req)
}

func (p *kBankClient) HandleTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleTransactionDetailReceipt(ctx, beginDate, endDate, organizationId)
}

func (p *kBankClient) CreateVirtualAccount(ctx context.Context, req *api.CreateVirtualAccountRequest, callOptions ...callopt.Option) (r *api.CreateVirtualAccountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateVirtualAccount(ctx, req)
}

func (p *kBankClient) SyncVirtualAccountBalance(ctx context.Context, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SyncVirtualAccountBalance(ctx)
}

func (p *kBankClient) QueryVirtualAccountBalance(ctx context.Context, organizationId int64, bankType string, callOptions ...callopt.Option) (r *api.VirtualAccountBalanceData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.QueryVirtualAccountBalance(ctx, organizationId, bankType)
}

func (p *kBankClient) SpdBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData, callOptions ...callopt.Option) (r *api.BankVirtualAccountTranscationResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SpdBankVirtualAccountTranscation(ctx, organizationId, req)
}

func (p *kBankClient) ListPaymentReceipt(ctx context.Context, req *api.ListPaymentReceiptRequest, callOptions ...callopt.Option) (r *api.ListPaymentReceiptResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListPaymentReceipt(ctx, req)
}

func (p *kBankClient) GetPaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (r *api.PaymentReceiptData, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetPaymentReceipt(ctx, id)
}

func (p *kBankClient) AddPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddPaymentReceipt(ctx, req)
}

func (p *kBankClient) ApprovePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ApprovePaymentReceipt(ctx, id, req)
}

func (p *kBankClient) RefusePaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RefusePaymentReceipt(ctx, id, req, remark)
}

func (p *kBankClient) PaymentReceiptRun(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PaymentReceiptRun(ctx, id)
}

func (p *kBankClient) TransmitPaymentReceipt(ctx context.Context, processInstanceId int64, transmitUserId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.TransmitPaymentReceipt(ctx, processInstanceId, transmitUserId)
}

func (p *kBankClient) SendBackPaymentApplication(ctx context.Context, id int64, req *api.PaymentReceiptData, remark string, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SendBackPaymentApplication(ctx, id, req, remark)
}

func (p *kBankClient) WithDrawPaymentReceipt(ctx context.Context, id int64, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.WithDrawPaymentReceipt(ctx, id, req)
}

func (p *kBankClient) CommentPaymentReceipt(ctx context.Context, req *api.PaymentReceiptData, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CommentPaymentReceipt(ctx, req)
}

func (p *kBankClient) HandleSyncPaymentReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleSyncPaymentReceipt(ctx, beginDate, endDate, organizationId)
}

func (p *kBankClient) SyncPaymentReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SyncPaymentReceipt(ctx, taskId, param, organizationId)
}

func (p *kBankClient) PinganBankAccountSignatureApply(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest, callOptions ...callopt.Option) (r *api.PinganUserAcctSignatureApplyResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PinganBankAccountSignatureApply(ctx, req)
}

func (p *kBankClient) PinganBankAccountSignatureQuery(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest, callOptions ...callopt.Option) (r *api.PinganUserAcctSignatureApplyResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PinganBankAccountSignatureQuery(ctx, req)
}

func (p *kBankClient) SystemRefusePaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SystemRefusePaymentReceipt(ctx, id)
}

func (p *kBankClient) SystemApprovePaymentReceipt(ctx context.Context, id int64, callOptions ...callopt.Option) (err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SystemApprovePaymentReceipt(ctx, id)
}
