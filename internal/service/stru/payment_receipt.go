package stru

import (
	"encoding/json"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/go-core/repository"
)

func ConvertPaymentReceiptDBData(data api.PaymentReceiptData) *repo.PaymentReceiptDBData {
	result := repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: data.OrganizationId,
			},
			ProcessInstanceId: data.ProcessInstanceId,
		},
		Code:                   data.Code,
		PayAmount:              data.PayAmount,
		PublicPrivateFlag:      data.PublicPrivateFlag,
		PayAccount:             data.PayAccount,
		PayAccountName:         data.PayAccountName,
		PayAccountBankName:     data.PayAccountBankName,
		PayAccountType:         data.PayAccountType,
		ReceiveAccount:         data.ReceiveAccount,
		ReceiveAccountName:     data.ReceiveAccountName,
		ReceiveAccountBankName: data.ReceiveAccountBankName,
		Purpose:                data.Purpose,
		UnionBankNo:            data.UnionBankNo,
		ClearBankNo:            data.ClearBankNo,
		InsideOutsideBankType:  data.InsideOutsideBankType,
		ChargeFee:              data.ChargeFee,
		OrderStatus:            data.OrderStatus,
		RetCode:                data.RetCode,
		RetMessage:             data.RetMessage,
		OrderFlowNo:            data.OrderFlowNo,
		Type:                   data.Type,
		PaymentModeType:        data.PaymentModeType,
		ApplicantName:          data.ApplicantName,
		ApplicantId:            data.ApplicantId,
		FillingDt:              data.FillingDt,
		DepartmentName:         data.DepartmentName,
		DepartmentId:           data.DepartmentId,
		Attachments:            data.Attachments,
		ElectronicDocument:     data.ElectronicDocument,
		ElectronicDocumentPng:  data.ElectronicDocumentPng,
		PaymentReason:          data.PaymentReason,
	}
	return &result
}

func ConvertPaymentReceiptData(dbData repo.PaymentReceiptDBData) *api.PaymentReceiptData {
	createdAtBytes, _ := json.Marshal(dbData.CreatedAt)
	updatedAtBytes, _ := json.Marshal(dbData.UpdatedAt)
	return &api.PaymentReceiptData{
		Id:                     dbData.Id,
		OrganizationId:         dbData.OrganizationId,
		CreatedAt:              createdAtBytes,
		UpdatedAt:              updatedAtBytes,
		ProcessInstanceId:      dbData.ProcessInstanceId,
		Code:                   dbData.Code,
		PayAmount:              dbData.PayAmount,
		PublicPrivateFlag:      dbData.PublicPrivateFlag,
		PayAccount:             dbData.PayAccount,
		PayAccountName:         dbData.PayAccountName,
		PayAccountBankName:     dbData.PayAccountBankName,
		PayAccountType:         dbData.PayAccountType,
		ReceiveAccount:         dbData.ReceiveAccount,
		ReceiveAccountName:     dbData.ReceiveAccountName,
		ReceiveAccountBankName: dbData.ReceiveAccountBankName,
		Purpose:                dbData.Purpose,
		UnionBankNo:            dbData.UnionBankNo,
		ClearBankNo:            dbData.ClearBankNo,
		InsideOutsideBankType:  dbData.InsideOutsideBankType,
		ChargeFee:              dbData.ChargeFee,
		OrderStatus:            dbData.OrderStatus,
		RetCode:                dbData.RetCode,
		RetMessage:             dbData.RetMessage,
		OrderFlowNo:            dbData.OrderFlowNo,
		Type:                   dbData.Type,
		PaymentModeType:        dbData.PaymentModeType,
		ApplicantName:          dbData.ApplicantName,
		ApplicantId:            dbData.ApplicantId,
		FillingDt:              dbData.FillingDt,
		DepartmentName:         dbData.DepartmentName,
		DepartmentId:           dbData.DepartmentId,
		Attachments:            dbData.Attachments,
		BusOrderNo:             dbData.BusOrderNo,
		BusType:                dbData.BusType,
		RefundSuccess:          dbData.RefundSuccess,
		ReceiptOrderNo:         dbData.ReceiptOrderNo,
		ElectronicDocument:     dbData.ElectronicDocument,
		ElectronicDocumentPng:  dbData.ElectronicDocumentPng,
		PaymentReason:          dbData.PaymentReason,
	}
}
