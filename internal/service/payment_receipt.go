package service

import (
	"context"
	"encoding/json"
	"fmt"
	bankEnum "gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk"
	sdkStru "gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service/stru"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	baseApi "gitlab.yoyiit.com/youyi/app-base/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-base/kitex_gen/api/base"
	"gitlab.yoyiit.com/youyi/app-base/process"
	"gitlab.yoyiit.com/youyi/app-dingtalk/kitex_gen/api/dingtalk"
	invoiceApi "gitlab.yoyiit.com/youyi/app-invoice/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-invoice/kitex_gen/api/invoice"
	oaApi "gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api/oa"
	"gitlab.yoyiit.com/youyi/app-soms/kitex_gen/api/soms"
	"gitlab.yoyiit.com/youyi/go-common/enum"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type PaymentReceiptService interface {
	process.ProcessInterface
	ListPaymentReceipt(ctx context.Context, req *api.ListPaymentReceiptRequest) (resp *api.ListPaymentReceiptResponse, err error)
	GetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error)
	SimpleGetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error)
	SimpleGetPaymentReceiptByProcessInstanceId(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error)
	PaymentReceiptRun(ctx context.Context, id int64) (err error)
	HandleSyncPaymentReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	SyncPaymentReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error)

	handleGuilinBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	handleSPDBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	handlePinganBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error

	updatePaymentReceiptResult(ctx context.Context, req repo.PaymentReceiptDBData, chargeFee float64, orderStatus,
		retCode, retMessage, orderFlowNo string) error
	guilinBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error)
	spdBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error)
	pinganBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error)
	PaymentReceiptSystemApprove(ctx context.Context, id int64) (err error)
	PaymentReceiptSystemRefuse(ctx context.Context, id int64) (err error)
}

type paymentReceiptService struct {
	process.Process
	paymentReceiptRepo                       repo.PaymentReceiptRepo
	baseClient                               base.Client
	bankCodeRepo                             repo.BankCodeRepo
	guilinBankSDK                            sdk.GuilinBankSDK
	spdBankSDK                               sdk.SPDBankSDK
	pinganBankSDK                            sdk.PinganBankSDK
	oaClient                                 oa.Client
	dingtalkClient                           dingtalk.Client
	somsClient                               soms.Client
	paymentReceiptApplicationCustomFieldRepo repo.PaymentReceiptApplicationCustomFieldRepo
	invoiceClient                            invoice.Client
}

func (s *paymentReceiptService) ListPaymentReceipt(ctx context.Context, req *api.ListPaymentReceiptRequest) (resp *api.ListPaymentReceiptResponse, err error) {
	dbData, canWriteData, count, err := s.paymentReceiptRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.PaymentReceiptDBParam{
		PaymentReceiptDBData: repo.PaymentReceiptDBData{
			Code:                   req.Code,
			ReceiveAccount:         req.ReceiveAccount,
			ReceiveAccountBankName: req.ReceiveAccountBankName,
			ReceiveAccountName:     req.ReceiveAccountName,
			Purpose:                req.Purpose,
			OrderStatus:            req.OrderStatus,
			OrderFlowNo:            req.OrderFlowNo,
			UnionBankNo:            req.UnionBankNo,
			ClearBankNo:            req.ClearBankNo,
			Type:                   req.Type,
			ProcessStatus:          req.ProcessStatus,
			ProcessCurrentUserName: req.ProcessCurrentUserName,
			ProcessCurrentUserId:   req.ProcessCurrentUserId,
			ProcessCodes:           req.ProcessCodes,
			ProcessName:            req.ProcessName,
			RefundSuccess:          req.RefundSuccess,
			ReceiptOrderNo:         req.ReceiptOrderNo,
			PaymentModeType:        req.PaymentModeType,
			ApplicantName:          req.ApplicantName,
		},
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
		BeginTime:       req.BeginTime,
		EndTime:         req.EndTime,
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	resultData := make([]*api.PaymentReceiptData, len(*dbData))
	canWriteMap := make(map[int64]bool)
	for _, v := range *canWriteData {
		canWriteMap[v.Id] = true
	}
	for i, v := range *dbData {
		result := stru.ConvertPaymentReceiptData(v)
		if canWriteMap[v.Id] == true {
			result.CanWrite = true
		}
		resultData[i] = result
	}
	return &api.ListPaymentReceiptResponse{
		Data:  resultData,
		Count: count,
	}, nil
}

func (s *paymentReceiptService) GetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	dbData, err := s.paymentReceiptRepo.Get(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			BaseDBData: repository.BaseDBData{
				BaseCommonDBData: repository.BaseCommonDBData{
					Id: id,
				},
			},
		},
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}

	customFields, _, err := s.paymentReceiptApplicationCustomFieldRepo.List(ctx, "created_at", 0, 0, &repo.PaymentReceiptApplicationCustomFieldDBData{
		PaymentReceiptId: dbData.Id,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	receipt := stru.ConvertPaymentReceiptData(*dbData)
	var resultCustomFields []*api.CustomField
	if customFields != nil && len(*customFields) > 0 {
		resultCustomFields = make([]*api.CustomField, len(*customFields))
		for i, v := range *customFields {
			resultCustomFields[i] = &api.CustomField{
				Id:             v.Id,
				Name:           v.ProcessCustomFieldName,
				Value:          v.ProcessCustomFieldValue,
				OrganizationId: v.OrganizationId,
				Sort:           v.ProcessCustomFieldSort,
				FieldId:        v.ProcessCustomFieldId,
			}
		}
		receipt.CustomFields = resultCustomFields
	}
	return receipt, nil
}

func (s *paymentReceiptService) SimpleGetPaymentReceipt(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	dbData, err := s.paymentReceiptRepo.SimpleGet(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			BaseDBData: repository.BaseDBData{
				BaseCommonDBData: repository.BaseCommonDBData{
					Id: id,
				},
			},
		},
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}

	customFields, _, err := s.paymentReceiptApplicationCustomFieldRepo.List(ctx, "created_at", 0, 0, &repo.PaymentReceiptApplicationCustomFieldDBData{
		PaymentReceiptId: dbData.Id,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	receipt := stru.ConvertPaymentReceiptData(*dbData)
	var resultCustomFields []*api.CustomField
	if customFields != nil && len(*customFields) > 0 {
		resultCustomFields = make([]*api.CustomField, len(*customFields))
		for i, v := range *customFields {
			resultCustomFields[i] = &api.CustomField{
				Id:             v.Id,
				Name:           v.ProcessCustomFieldName,
				Value:          v.ProcessCustomFieldValue,
				OrganizationId: v.OrganizationId,
				Sort:           v.ProcessCustomFieldSort,
				FieldId:        v.ProcessCustomFieldId,
			}
		}
		receipt.CustomFields = resultCustomFields
	}
	return receipt, nil
}

func (s *paymentReceiptService) SimpleGetPaymentReceiptByProcessInstanceId(ctx context.Context, id int64) (resp *api.PaymentReceiptData, err error) {
	dbData, err := s.paymentReceiptRepo.SimpleGet(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	return stru.ConvertPaymentReceiptData(*dbData), nil
}

func (s *paymentReceiptService) PaymentReceiptRun(ctx context.Context, id int64) (err error) {
	paymentReceipt, err := s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
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
	if paymentReceipt == nil || paymentReceipt.Id == 0 {
		return handler.HandleNewError("payment receipt not exists")
	}
	if paymentReceipt.OrderStatus != "" {
		if paymentReceipt.OrderStatus == "0" {
			var orderStatus string
			if paymentReceipt.PayAmount > 0 {
				if paymentReceipt.PayAccount == "" {
					return handler.HandleNewError("PayAccount must not empty")
				}
				if paymentReceipt.PayAccountName == "" {
					return handler.HandleNewError("PayAccountName must not empty")
				}
				if paymentReceipt.ReceiveAccount == "" {
					return handler.HandleNewError("RecAccount must not empty")
				}
				if paymentReceipt.ReceiveAccountName == "" {
					return handler.HandleNewError("RecAccountName must not empty")
				}
				if paymentReceipt.PayAmount <= 0 {
					return handler.HandleNewError("PayAmount must greater than 0")
				}
				if paymentReceipt.PublicPrivateFlag == "" {
					return handler.HandleNewError("PubPriFlag must not empty")
				}
				//处理行内行外
				insideOutsideBankType := "0"
				if paymentReceipt.PayAccountBankName != paymentReceipt.ReceiveAccountBankName {
					insideOutsideBankType = "1"
				}
				paymentReceipt.InsideOutsideBankType = insideOutsideBankType
				//联行号清算行号
				bankCode, err := s.bankCodeRepo.Get(ctx, &repo.BankCodeDBData{
					BankName: paymentReceipt.ReceiveAccountBankName,
				})
				if err != nil {
					return handler.HandleError(err)
				}
				paymentReceipt.UnionBankNo = bankCode.UnionBankNo
				paymentReceipt.ClearBankNo = bankCode.ClearBankNo
				//处理备注
				payRemark := fmt.Sprintf("%s[%s]", paymentReceipt.Purpose, paymentReceipt.Code)
				//处理账户空格
				paymentReceipt.ReceiveAccount = strings.ReplaceAll(paymentReceipt.ReceiveAccount, " ", "")
				paymentReceipt.PayAccount = strings.ReplaceAll(paymentReceipt.PayAccount, " ", "")

				switch paymentReceipt.PayAccountType {
				case "0":
					orderStatus, err = s.guilinBankPayment(ctx, paymentReceipt, payRemark)
				case "1":
					// 转账前校验浦发银行接口是否已经存在这个流水数据, 存在则获取交易状态, 否则重新转账
					organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
						OrganizationId: paymentReceipt.OrganizationId,
						Type:           enum.SPDBankType,
					})
					if err != nil {
						return handler.HandleError(err)
					}
					transferDate := stru.FormatDayTime(paymentReceipt.CreatedAt)
					bankTransferResultResponseData, err := s.spdBankSDK.QueryTransferResult(ctx, paymentReceipt.PayAccount, paymentReceipt.ReceiveAccount, paymentReceipt.OrderFlowNo, paymentReceipt.Code, transferDate, transferDate, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
					zap.L().Info(fmt.Sprintf("PaymentReceiptRun spdBankSDK.QueryTransferResult查询转账结果:%v", bankTransferResultResponseData))
					// 更新单据状态
					if bankTransferResultResponseData != nil && len(bankTransferResultResponseData.Items) > 0 {
						bankTransferResponse := bankTransferResultResponseData.Items[0]
						orderStatus = enum.GetOrderState(bankTransferResponse.TransStatus, bankTransferResponse.TransStatus)
						returnCode := bankTransferResponse.TransStatus
						returnMsg := enum.GetOrderMsg(bankTransferResponse.TransStatus, "")
						if err := s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, &repo.PaymentReceiptDBData{
							OrderStatus: orderStatus,
							RetCode:     returnCode,
							RetMessage:  returnMsg,
							OrderFlowNo: bankTransferResponse.AcceptNo,
						}); err != nil {
							return handler.HandleError(err)
						}
					} else {
						orderStatus, err = s.spdBankPayment(ctx, paymentReceipt, payRemark)
						if err != nil {
							return handler.HandleError(err)
						}
					}
				case "2":
					orderStatus, err = s.pinganBankPayment(ctx, paymentReceipt, payRemark)
					if err != nil {
						return handler.HandleError(err)
					}
				}
			} else {
				// 交易金额为0 直接变更为成功
				orderStatus = enum.GuilinBankTransferSuccessResult
				if err := s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, &repo.PaymentReceiptDBData{
					OrderStatus: orderStatus,
					RetCode:     "000000",
				}); err != nil {
					return handler.HandleError(err)
				}
			}
			if orderStatus != "" {
				//同步付款申请的银行状态
				if err = s.updateApplicationOrderStatus(ctx, paymentReceipt, orderStatus); err != nil {
					return handler.HandleError(err)
				}
				return handler.HandleError(s.handleProcessResult(ctx, orderStatus, paymentReceipt))
			}
		} else {
			//同步付款申请的银行状态
			if err = s.updateApplicationOrderStatus(ctx, paymentReceipt, paymentReceipt.OrderStatus); err != nil {
				return handler.HandleError(err)
			}
			return handler.HandleError(s.handleProcessResult(ctx, paymentReceipt.OrderStatus, paymentReceipt))
		}
	}
	return nil
}

func (s *paymentReceiptService) updateApplicationOrderStatus(ctx context.Context, paymentReceipt *repo.PaymentReceiptDBData, orderStatus string) error {
	if paymentReceipt.Type == "1" {
		paymentApplication, err := s.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, paymentReceipt.ProcessInstanceId)
		if err != nil {
			return handler.HandleError(err)
		}
		if paymentApplication != nil && paymentApplication.Id != 0 && paymentApplication.OrderStatus != orderStatus {
			if err = s.oaClient.EditPaymentApplicationWithoutPermission(ctx, &oaApi.PaymentApplicationData{
				Id:          paymentApplication.Id,
				OrderStatus: orderStatus,
			}); err != nil {
				return handler.HandleError(err)
			}
		}
	} else {
		reimburseApplication, err := s.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, paymentReceipt.ProcessInstanceId)
		if err != nil {
			return handler.HandleError(err)
		}
		if reimburseApplication != nil && reimburseApplication.Id != 0 && reimburseApplication.OrderStatus != orderStatus {
			if err = s.oaClient.EditReimburseApplicationWithoutPermission(ctx, &oaApi.ReimburseApplicationData{
				Id:          reimburseApplication.Id,
				OrderStatus: orderStatus,
			}); err != nil {
				return handler.HandleError(err)
			}
		}
	}
	/*//调用进销存模块
	if orderStatus == enum.GuilinBankTransferSuccessResult && paymentApplication.BusType != "" && paymentApplication.BusType != "0" && paymentApplication.BusOrderNo != "" {
		s.somsClient.EditSimsOrderAmount(ctx, &api2.SimsOrderAmountData{
			ReturnFlag: "0",
			PayType:    "1",
			Amount:     paymentApplication.PayAmount,
			BusType:    paymentApplication.BusType,
			BillCode:   paymentApplication.BusOrderNo,
		})
	}*/
	return nil
}

func (s *paymentReceiptService) guilinBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error) {
	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: req.OrganizationId,
		Type:           enum.GuilinBankType,
	})
	if err != nil {
		return "", handler.HandleError(err)
	}
	var bankTransferResponse sdkStru.BankTransferResponse
	switch req.PaymentModeType {
	case "0":
		res, err := s.guilinBankSDK.LedgerManagementTransfer(ctx, req.Code, sdkStru.OutofbankTransferRequest{
			PayAccount:         req.PayAccount,
			PayAccountName:     req.PayAccountName,
			RecAccount:         req.ReceiveAccount,
			RecAccountName:     req.ReceiveAccountName,
			PayAmount:          req.PayAmount,
			PayRem:             payRemark,
			PubPriFlag:         req.PublicPrivateFlag,
			CurrencyType:       "CNY",
			RecBankType:        req.InsideOutsideBankType,
			TransferFlag:       "1",
			UnionBankNo:        req.UnionBankNo,
			ClearBankNo:        req.ClearBankNo,
			RecAccountOpenBank: req.ReceiveAccountBankName,
			RmtType:            "0",
		}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			return "", handler.HandleError(err)
		}
		if res != nil {
			bankTransferResponse = *res
		}
	case "1":
		switch req.InsideOutsideBankType {
		case "0":
			res, err := s.guilinBankSDK.IntrabankTransfer(ctx, req.Code, sdkStru.IntrabankTransferRequest{
				PayAccount:     req.PayAccount,
				PayAccountName: req.PayAccountName,
				RecAccount:     req.ReceiveAccount,
				RecAccountName: req.ReceiveAccountName,
				PayAmount:      req.PayAmount,
				PayRem:         payRemark,
				PubPriFlag:     req.PublicPrivateFlag,
				CurrencyType:   "CNY",
				RecBankType:    req.InsideOutsideBankType,
			}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			if err != nil {
				return "", handler.HandleError(err)
			}
			if res != nil {
				bankTransferResponse = *res
			}
		case "1":
			res, err := s.guilinBankSDK.OutofbankTransfer(ctx, req.Code, sdkStru.OutofbankTransferRequest{
				PayAccount:         req.PayAccount,
				PayAccountName:     req.PayAccountName,
				RecAccount:         req.ReceiveAccount,
				RecAccountName:     req.ReceiveAccountName,
				PayAmount:          req.PayAmount,
				PayRem:             payRemark,
				PubPriFlag:         req.PublicPrivateFlag,
				CurrencyType:       "CNY",
				RecBankType:        req.InsideOutsideBankType,
				TransferFlag:       "1",
				UnionBankNo:        req.UnionBankNo,
				ClearBankNo:        req.ClearBankNo,
				RecAccountOpenBank: req.ReceiveAccountBankName,
				RmtType:            "0",
			}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			if err != nil {
				return "", handler.HandleError(err)
			}
			if res != nil {
				bankTransferResponse = *res
			}
		}
	}

	if bankTransferResponse.Head.RetCode != enum.GuilinBankSuccessRetCode {
		bankTransferResponse.Body.OrderState = enum.GuilinBankTransferFailResult
	} else if req.PaymentModeType == "0" {
		bankTransferResponse.Body.OrderState = "51"
	}
	if err = s.updatePaymentReceiptResult(ctx, *req, bankTransferResponse.Body.ChargeFee, bankTransferResponse.Body.OrderState,
		bankTransferResponse.Head.RetCode, bankTransferResponse.Head.RetMessage, bankTransferResponse.Head.OrderFlowNo); err != nil {
		return "", handler.HandleError(err)
	}
	return bankTransferResponse.Body.OrderState, nil
}

func (s *paymentReceiptService) updatePaymentReceiptResult(ctx context.Context, req repo.PaymentReceiptDBData, chargeFee float64, orderStatus, retCode, retMessage, orderFlowNo string) error {
	return s.paymentReceiptRepo.UpdateById(ctx, req.Id, &repo.PaymentReceiptDBData{
		InsideOutsideBankType: req.InsideOutsideBankType,
		UnionBankNo:           req.UnionBankNo,
		ClearBankNo:           req.ClearBankNo,
		ChargeFee:             chargeFee,
		OrderStatus:           orderStatus,
		RetCode:               retCode,
		RetMessage:            retMessage,
		OrderFlowNo:           orderFlowNo,
	})
}

// 返回0:未结束 1:成功 2:失败
func judgeBankPaymentEnd(orderStatus string) int {
	if orderStatus == enum.GuilinBankTransferSuccessResult {
		return 1
	} else if orderStatus == enum.GuilinBankTransferRevokeResult || orderStatus == enum.GuilinBankTransferDeleteResult ||
		orderStatus == enum.GuilinBankTransferFailResult || orderStatus == enum.GuilinBankTransferRefuseResult {
		return 2
	} else {
		return 0
	}
}

func (s *paymentReceiptService) sendMessage(ctx context.Context, paymentReceiptDBData *repo.PaymentReceiptDBData, result int) error {
	resultStr := "失败"
	statusBg := "EF4444"
	if result == 1 {
		resultStr = "成功"
		statusBg = "67c23a"
	}
	organization, err := s.baseClient.GetOrganization(ctx, &baseApi.OrganizationData{
		Id: paymentReceiptDBData.OrganizationId,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	processInstanceMap, err := s.baseClient.QueryProcessInstance(ctx, []int64{paymentReceiptDBData.ProcessInstanceId})
	if err != nil {
		return handler.HandleError(err)
	}
	processInstance := processInstanceMap[paymentReceiptDBData.ProcessInstanceId]
	if processInstance == nil || processInstance.Id == 0 {
		return handler.HandleNewError("process instance not exists")
	}
	userOrganization, err := s.baseClient.GetUserOrganization(ctx, &baseApi.UserOrganizationData{
		OrganizationId: organization.Id,
		UserId:         processInstance.OriginatorUserId,
	})
	if err != nil {
		return handler.HandleError(err)
	}

	var paymentId int64 = 0
	if paymentReceiptDBData.Type == "1" || paymentReceiptDBData.Type == "" {
		paymentApplication, err := s.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, paymentReceiptDBData.ProcessInstanceId)
		if err != nil {
			return handler.HandleError(err)
		}
		err = s.dingtalkClient.SendCorpOACorpConversation(ctx,
			map[string]string{
				"statusBg":             statusBg,
				"resultStr":            resultStr,
				"name":                 processInstance.Name,
				"code":                 processInstance.Code,
				"paymentReceiptCode":   paymentReceiptDBData.Code,
				"payAccount":           paymentReceiptDBData.PayAccount,
				"payAccountName":       paymentReceiptDBData.PayAccountName,
				"receiveAccount":       paymentReceiptDBData.ReceiveAccount,
				"receiveAccountName":   paymentReceiptDBData.ReceiveAccountName,
				"payAmount":            strconv.FormatFloat(paymentReceiptDBData.PayAmount, 'f', 2, 64),
				"paymentApplicationId": strconv.FormatInt(paymentApplication.Id, 10),
			}, organization.DingtalkCorpId, userOrganization.DingtalkUserId, config.GetString("dingtalk.paymentApplication.template", ""))
		if err != nil {
			errStr := fmt.Sprintf("付款申请单=发送消息给发起人userId=%d,nickName=%s 失败,原因如下", userOrganization.UserId, userOrganization.Nickname)
			println(errStr)
			println(err.Error())
		}
		paymentId = paymentApplication.Id
	} else if paymentReceiptDBData.Type == "2" {
		reimburseApplication, err := s.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, paymentReceiptDBData.ProcessInstanceId)
		if err != nil {
			return handler.HandleError(err)
		}
		err = s.dingtalkClient.SendCorpOACorpConversation(ctx,
			map[string]string{
				"statusBg":               statusBg,
				"resultStr":              resultStr,
				"name":                   processInstance.Name,
				"code":                   processInstance.Code,
				"paymentReceiptCode":     paymentReceiptDBData.Code,
				"payAccount":             paymentReceiptDBData.PayAccount,
				"payAccountName":         paymentReceiptDBData.PayAccountName,
				"receiveAccount":         paymentReceiptDBData.ReceiveAccount,
				"receiveAccountName":     paymentReceiptDBData.ReceiveAccountName,
				"reimbursementMount":     strconv.FormatFloat(reimburseApplication.ReimbursementMount, 'f', 2, 64),
				"amountOnCredit":         strconv.FormatFloat(reimburseApplication.AmountOnCredit, 'f', 2, 64),
				"offsetAmount":           strconv.FormatFloat(reimburseApplication.OffsetAmount, 'f', 2, 64),
				"payAmount":              strconv.FormatFloat(paymentReceiptDBData.PayAmount, 'f', 2, 64),
				"reimburseApplicationId": strconv.FormatInt(reimburseApplication.Id, 10),
			}, organization.DingtalkCorpId, userOrganization.DingtalkUserId, config.GetString("dingtalk.reimburseApplication.template", ""))
		if err != nil {
			errStr := fmt.Sprintf("报销申请单=发送消息给发起人userId=%d,nickName=%s 失败,原因如下", userOrganization.UserId, userOrganization.Nickname)
			println(errStr)
			println(err.Error())
		}
		paymentId = reimburseApplication.Id
	}

	//发送给抄送人
	ccopyUsersOrganization, err := s.baseClient.ListCcopyUserOrganizationByProcessInstanceId(ctx, &baseApi.UserOrganizationData{
		OrganizationId:    organization.Id,
		ProcessInstanceId: processInstance.Id,
	})
	if err != nil || ccopyUsersOrganization == nil {
		return handler.HandleError(err)
	}

	for _, cc := range ccopyUsersOrganization {
		orderType := "付款申请"
		if paymentReceiptDBData.Type == "2" {
			orderType = "报销申请"
		}
		err := s.send(ctx, processInstance, paymentReceiptDBData, cc, organization, statusBg, resultStr, orderType, paymentId)
		if err != nil {
			errStr := fmt.Sprintf("%s=发送消息给抄送人userId=%d,nickName=%s 失败,原因如下", orderType, cc.UserId, cc.Nickname)
			println(errStr)
			println(err.Error())
		}
	}
	return nil
}

func (s *paymentReceiptService) send(ctx context.Context, processInstance *baseApi.ProcessInstanceData, receiptConfirmOrderDBData *repo.PaymentReceiptDBData,
	userOrganization *baseApi.UserOrganizationData, organization *baseApi.OrganizationData, statusBg, resultStr, resultType string,
	paymentId int64) error {
	if userOrganization.DingtalkUserId == "" {
		return handler.HandleNewError("DingtalkUserId为空,不是钉钉用户")
	}
	url := fmt.Sprintf("page/paymentDetail/index?id=%d", paymentId)
	if receiptConfirmOrderDBData.Type == "2" {
		url = fmt.Sprintf("reimburseApplication/detail/index?id=%d", paymentId)
	}
	err := s.dingtalkClient.SendCorpOACorpConversation(ctx,
		map[string]string{
			"statusBg":           statusBg,
			"resultStr":          resultStr,
			"resultType":         resultType,
			"name":               processInstance.Name,
			"code":               processInstance.Code,
			"paymentReceiptCode": receiptConfirmOrderDBData.Code,
			"payAccount":         receiptConfirmOrderDBData.PayAccount,
			"payAccountName":     receiptConfirmOrderDBData.PayAccountName,
			"receiveAccount":     receiptConfirmOrderDBData.ReceiveAccount,
			"receiveAccountName": receiptConfirmOrderDBData.ReceiveAccountName,
			"payAmount":          strconv.FormatFloat(receiptConfirmOrderDBData.PayAmount, 'f', 2, 64),
			"url":                url,
		}, organization.DingtalkCorpId, userOrganization.DingtalkUserId, config.GetString("dingtalk.Ccopy.template", ""))
	if err != nil {
		return handler.HandleError(err)
	}
	return err
}

func (s *paymentReceiptService) handleProcessResult(ctx context.Context, orderStatus string, paymentReceiptDBData *repo.PaymentReceiptDBData) error {
	result := judgeBankPaymentEnd(orderStatus)
	if result == 1 {
		processResult, err := s.baseClient.SuccessProcessInstance(ctx, paymentReceiptDBData.ProcessInstanceId)
		if err != nil {
			return err
		}
		err = s.paymentReceiptRepo.UpdateSelectedFieldsByIdWithoutPermission(ctx, paymentReceiptDBData.Id, &repo.PaymentReceiptDBData{
			ProcessStatus: processResult.ProcessStatus,
		}, []string{"ProcessStatus"})
		if err != nil {
			return handler.HandleError(err)
		}

		// 完成发票
		if paymentReceiptDBData.Type == "2" && paymentReceiptDBData.PaymentId > 0 {
			response, err := s.oaClient.ListReimburseInvoiceApplication(ctx, paymentReceiptDBData.PaymentId)
			if err != nil {
				return handler.HandleError(err)
			}
			if response != nil && response.Count > 0 {
				var invoiceIds []int64
				for _, v := range response.Data {
					invoiceIds = append(invoiceIds, v.InvoiceId)
				}
				s.invoiceClient.CompleteReimburseInvoice(ctx, &invoiceApi.CompleteReimburseInvoiceRequest{
					InvoiceIds: invoiceIds,
					ReceiptId:  paymentReceiptDBData.PaymentId,
				})
			}
		}
		return s.sendMessage(ctx, paymentReceiptDBData, result)
	} else if result == 2 {
		processResult, err := s.baseClient.RefuseProcessInstance(ctx, paymentReceiptDBData.ProcessInstanceId, s.SubProcess.ProcessNodeStep(), "")
		if err != nil {
			return err
		}
		err = s.paymentReceiptRepo.UpdateSelectedFieldsByIdWithoutPermission(ctx, paymentReceiptDBData.Id, &repo.PaymentReceiptDBData{
			ProcessStatus: processResult.ProcessStatus,
		}, []string{"ProcessStatus"})
		if err != nil {
			return handler.HandleError(err)
		}

		// 取消发票
		if paymentReceiptDBData.Type == "2" && paymentReceiptDBData.PaymentId > 0 {
			response, err := s.oaClient.ListReimburseInvoiceApplication(ctx, paymentReceiptDBData.PaymentId)
			if err != nil {
				return handler.HandleError(err)
			}
			if response != nil && response.Count > 0 {
				var invoiceIds []int64
				for _, v := range response.Data {
					invoiceIds = append(invoiceIds, v.InvoiceId)
				}
				s.invoiceClient.CancelReimburseInvoice(ctx, &invoiceApi.CancelReimburseInvoiceRequest{
					InvoiceIds: invoiceIds,
					ReceiptId:  paymentReceiptDBData.PaymentId,
				})
			}
		}
		return s.sendMessage(ctx, paymentReceiptDBData, result)
	}
	return nil
}

func (s *paymentReceiptService) spdBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error) {
	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: req.OrganizationId,
		Type:           enum.SPDBankType,
	})
	if err != nil {
		return "", handler.HandleError(err)
	}
	// 收款行名称 (当SysFlag=1即跨行转帐时必须输入)
	payeeBankName := ""
	// 同城异地标志
	remitLocation := "0"
	// 收款行速选标志 (1-速选, 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效。如果希望跨行汇款自动处理，请务必填写此项。
	payeeBankSelectFlag := ""
	// 收款行行号（人行现代支付系统行号）如果速选标志为1，请务必填写此项。 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效
	payeeBankNo := ""
	if req.InsideOutsideBankType == "1" {
		// 由于分行同城已经取消，跨行支付时，同城异地标志建议固定送1异地
		remitLocation = "1"
		payeeBankName = req.ReceiveAccountBankName
		payeeBankSelectFlag = "1"
		payeeBankNo = req.UnionBankNo
	}
	request := sdkStru.SPDBankTransferRequest{
		AuthMasterID:        organizationBankConfig.BankCustomerId,
		ElecChequeNo:        req.Code,
		AcctNo:              req.PayAccount,
		AcctName:            req.PayAccountName,
		PayeeAcctNo:         req.ReceiveAccount,
		PayeeName:           req.ReceiveAccountName,
		PayeeBankName:       payeeBankName,
		Amount:              req.PayAmount,
		SysFlag:             req.InsideOutsideBankType,
		RemitLocation:       remitLocation,
		Note:                payRemark,
		PayeeBankSelectFlag: payeeBankSelectFlag,
		PayeeBankNo:         payeeBankNo,
	}
	bankTransferResponse, err := s.spdBankSDK.BankTransfer(ctx, req.Code, request,
		organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return "", handler.HandleError(err)
	}
	// 处理订单状态
	orderStatus := bankTransferResponse.TransStatus
	if bankTransferResponse.TransStatus != "" {
		orderStatus = enum.GetOrderState(bankTransferResponse.TransStatus, orderStatus)
		bankTransferResponse.ReturnMsg = enum.GetOrderMsg(bankTransferResponse.TransStatus, bankTransferResponse.ReturnMsg)
	}
	if err = s.updatePaymentReceiptResult(ctx, *req, 0, orderStatus,
		bankTransferResponse.TransStatus, bankTransferResponse.ReturnMsg, bankTransferResponse.AcceptNo); err != nil {
		return "", handler.HandleError(err)
	}
	return orderStatus, nil
}

func (s *paymentReceiptService) pinganBankPayment(ctx context.Context, req *repo.PaymentReceiptDBData, payRemark string) (string, error) {

	// 收款行名称 (当SysFlag=1即跨行转帐时必须输入)
	payeeBankName := ""
	// 同城异地标志
	remitLocation := "0"

	// 收款行行号（人行现代支付系统行号）如果速选标志为1，请务必填写此项。 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效
	payeeBankNo := ""
	remitLocation = "1"
	//1：行内转账，0：跨行转账
	unionFlag := "1"
	if req.InsideOutsideBankType == "0" {
		unionFlag = "1"
	} else {
		unionFlag = "0"
	}
	if req.InsideOutsideBankType == "1" {
		// 由于分行同城已经取消，跨行支付时，同城异地标志建议固定送1异地
		payeeBankName = req.ReceiveAccountBankName
		payeeBankNo = req.UnionBankNo
	}
	uuid, _ := util.SonyflakeID()
	request := sdkStru.PingAnBankTransferRequest{
		MrchCode:        config.GetString(bankEnum.PinganMrchCode, ""),
		CnsmrSeqNo:      uuid,
		OutAcctNo:       req.PayAccount,
		OutAcctName:     req.PayAccountName,
		ThirdVoucher:    req.Code,
		CcyCode:         "RMB",
		InAcctNo:        req.ReceiveAccount,
		InAcctName:      req.ReceiveAccountName,
		InAcctBankName:  payeeBankName,
		TranAmount:      strconv.FormatFloat(req.PayAmount, 'f', -1, 64),
		UnionFlag:       unionFlag,
		AddrFlag:        remitLocation,
		UseEx:           payRemark,
		InAcctBankNode:  payeeBankNo,
		PaymentModeType: req.PaymentModeType,
	}
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Type:    "2",
		Account: req.PayAccount,
	})
	if err != nil {
		return "", handler.HandleError(err)
	}
	transfer, err := s.pinganBankSDK.BankTransfer(ctx, request, bankAccount.ZuId)
	if err != nil {
		return "", handler.HandleError(err)
	}
	chargeF, err := strconv.ParseFloat(transfer.Fee1, 64)
	if err != nil {
		return "", handler.HandleError(err)
	}
	orderStatus := enum.GetOrderState(transfer.Stt, transfer.Stt)
	if err = s.updatePaymentReceiptResult(ctx, *req, chargeF, orderStatus,
		transfer.Stt, transfer.HostTxDate, transfer.HostFlowNo); err != nil {
		return "", handler.HandleError(err)
	}
	return orderStatus, nil
}

func (s *paymentReceiptService) HandleSyncPaymentReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	s.handleGuilinBankSyncPaymentReceipt(ctx, enum.GuilinBankType, beginDate, endDate, organizationId)
	s.handleSPDBankSyncPaymentReceipt(ctx, enum.SPDBankType, beginDate, endDate, organizationId)
	s.handlePinganBankSyncPaymentReceipt(ctx, enum.PinganBankType, beginDate, endDate, organizationId)
	return nil
}

func (s *paymentReceiptService) handleGuilinBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("HandleGuilinBankSyncTransferReceipt 查询桂林所有未完成的付款单据")
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRefuseResult}
	var createTimeParam []string
	if beginDate != "" && beginDate != "" {
		begin, err := util.ParseDate(beginDate)
		if err != nil {
			return handler.HandleError(err)
		}
		end, err := util.ParseDate(endDate)
		if err != nil {
			return handler.HandleError(err)
		}
		createTimeParam = append(createTimeParam, util.FormatTimeyyyyMMddBar(begin), util.FormatTimeyyyyMMddBar(end))
	}
	paymentReceiptList, _, err := s.paymentReceiptRepo.ListAll(ctx, "", 0, 0, &repo.PaymentReceiptDBParam{
		ExcludeOrderStates:    excludeOrderStates,
		IsOrderFlowNoNotEmpty: true,
		PaymentReceiptDBData: repo.PaymentReceiptDBData{
			BaseProcessDBData: repository.BaseProcessDBData{
				BaseDBData: repository.BaseDBData{
					OrganizationId: organizationId,
				},
			},
			PayAccountType: bankType,
		},
		CreateTime: createTimeParam,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if paymentReceiptList != nil && len(*paymentReceiptList) > 0 {
		for _, paymentReceipt := range *paymentReceiptList {
			//查询转账结果
			isIntrabankTransfer := paymentReceipt.InsideOutsideBankType == "0"
			beginDate2 := stru.FormatDayTime(paymentReceipt.CreatedAt)
			endDate2 := stru.FormatDayTime(time.Now())
			if beginDate2 == endDate {
				endDate2 = stru.FormatDayTimeStamp(time.Now().Unix() + 86400)
			}
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: paymentReceipt.OrganizationId,
				Type:           bankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			//处理账户空格
			paymentReceipt.ReceiveAccount = strings.ReplaceAll(paymentReceipt.ReceiveAccount, " ", "")
			paymentReceipt.PayAccount = strings.ReplaceAll(paymentReceipt.PayAccount, " ", "")

			orderFlowNo := paymentReceipt.OrderFlowNo
			if paymentReceipt.PaymentModeType == "0" {
				orderFlowNo = paymentReceipt.Code
			}

			transferResults, err := s.guilinBankSDK.QueryTransferResult(ctx, paymentReceipt.PayAccount, paymentReceipt.ReceiveAccount, orderFlowNo, isIntrabankTransfer, beginDate2, endDate2, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			zap.L().Info(fmt.Sprintf("guilinBankSDK.QueryTransferResult查询转账结果:%v", transferResults))
			if err != nil {
				return handler.HandleError(err)
			}
			if transferResults != nil && len(transferResults) > 0 {
				for _, result := range transferResults {
					//比较订单状态
					if result.OrderState == paymentReceipt.OrderStatus && result.ErrorCode == "" {
						continue
					}
					updateReceipt := &repo.PaymentReceiptDBData{
						OrderStatus: result.OrderState,
					}
					if result.ErrorCode != "" && result.ErrorCode != "000000" {
						updateReceipt.RetCode = result.ErrorCode
						updateReceipt.RetMessage = result.ErrorMessage
					}
					if err = s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, updateReceipt); err != nil {
						return handler.HandleError(err)
					}
				}
			} else if paymentReceipt.PaymentModeType == "0" {
				//账管云的失败需要用这个接口
				ledgerManagementTransferResults, err := s.guilinBankSDK.QueryLedgerManagementTransferResult(ctx, paymentReceipt.PayAccount, paymentReceipt.ReceiveAccount, orderFlowNo, isIntrabankTransfer, beginDate2, endDate2, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
				zap.L().Info(fmt.Sprintf("guilinBankSDK.QueryLedgerManagementTransferResult查询转账结果:%v", transferResults))
				if err != nil {
					return handler.HandleError(err)
				}
				if ledgerManagementTransferResults != nil && len(ledgerManagementTransferResults) > 0 {
					for _, result := range ledgerManagementTransferResults {
						if result.OrderState == "99" {
							updateReceipt := &repo.PaymentReceiptDBData{
								OrderStatus: result.OrderState,
								RetMessage:  result.AuthRejectReason,
							}
							if err = s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, updateReceipt); err != nil {
								return handler.HandleError(err)
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func (s *paymentReceiptService) handleSPDBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("HandleSPDBankSyncTransferReceipt 查询浦发所有未完成的付款单据")
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRefuseResult}
	var createTimeParam []string
	if beginDate != "" && beginDate != "" {
		begin, err := util.ParseDate(beginDate)
		if err != nil {
			return handler.HandleError(err)
		}
		end, err := util.ParseDate(endDate)
		if err != nil {
			return handler.HandleError(err)
		}
		createTimeParam = append(createTimeParam, util.FormatTimeyyyyMMddBar(begin), util.FormatTimeyyyyMMddBar(end))
	}
	paymentReceiptList, _, err := s.paymentReceiptRepo.ListAll(ctx, "", 0, 0, &repo.PaymentReceiptDBParam{
		ExcludeOrderStates:    excludeOrderStates,
		IsOrderFlowNoNotEmpty: true,
		PaymentReceiptDBData: repo.PaymentReceiptDBData{
			BaseProcessDBData: repository.BaseProcessDBData{
				BaseDBData: repository.BaseDBData{
					OrganizationId: organizationId,
				},
			},
			PayAccountType: bankType,
		},
		CreateTime: createTimeParam,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if paymentReceiptList != nil && len(*paymentReceiptList) > 0 {
		for _, paymentReceipt := range *paymentReceiptList {
			beginDate2 := stru.FormatDayTime(paymentReceipt.CreatedAt)
			endDate2 := stru.FormatDayTime(time.Now())
			if beginDate2 == endDate {
				endDate2 = stru.FormatDayTimeStamp(time.Now().Unix() + 86400)
			}
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: paymentReceipt.OrganizationId,
				Type:           bankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			//处理账户空格
			paymentReceipt.ReceiveAccount = strings.ReplaceAll(paymentReceipt.ReceiveAccount, " ", "")
			paymentReceipt.PayAccount = strings.ReplaceAll(paymentReceipt.PayAccount, " ", "")
			bankTransferResultResponseData, err := s.spdBankSDK.QueryTransferResult(ctx, paymentReceipt.PayAccount, paymentReceipt.ReceiveAccount, paymentReceipt.OrderFlowNo, paymentReceipt.Code, beginDate2, endDate2, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			zap.L().Info(fmt.Sprintf("spdBankSDK.QueryTransferResult查询转账结果:%v", bankTransferResultResponseData))
			if err != nil {
				return handler.HandleError(err)
			}
			if bankTransferResultResponseData != nil && len(bankTransferResultResponseData.Items) > 0 {
				for _, result := range bankTransferResultResponseData.Items {
					// 处理订单状态(默认失败)
					orderState := enum.GuilinBankTransferFailResult
					// 获取映射的订单状态
					result.TransStatus = enum.GetOrderState(result.TransStatus, orderState)

					//比较订单状态
					if result.TransStatus == paymentReceipt.OrderStatus {
						continue
					}
					updateReceipt := &repo.PaymentReceiptDBData{
						OrderStatus: result.TransStatus,
					}
					if bankTransferResultResponseData.ReturnCode != "" && bankTransferResultResponseData.ReturnCode != enum.SPDBankSuccessRetCode {
						updateReceipt.RetCode = bankTransferResultResponseData.ReturnCode
						updateReceipt.RetMessage = bankTransferResultResponseData.ReturnMsg
					}
					if err = s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, updateReceipt); err != nil {
						return handler.HandleError(err)
					}
				}
			}
		}
	}
	return nil
}
func (s *paymentReceiptService) handlePinganBankSyncPaymentReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("handlePinganBankSyncPaymentReceipt 查询平安所有未完成的付款单据")
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRefuseResult}
	var createTimeParam []string
	if beginDate != "" && beginDate != "" {
		begin, err := util.ParseDate(beginDate)
		if err != nil {
			return handler.HandleError(err)
		}
		end, err := util.ParseDate(endDate)
		if err != nil {
			return handler.HandleError(err)
		}
		createTimeParam = append(createTimeParam, util.FormatTimeyyyyMMddBar(begin), util.FormatTimeyyyyMMddBar(end))
	}
	paymentReceiptList, _, err := s.paymentReceiptRepo.ListAll(ctx, "", 0, 0, &repo.PaymentReceiptDBParam{
		ExcludeOrderStates:    excludeOrderStates,
		IsOrderFlowNoNotEmpty: true,
		PaymentReceiptDBData: repo.PaymentReceiptDBData{
			BaseProcessDBData: repository.BaseProcessDBData{
				BaseDBData: repository.BaseDBData{
					OrganizationId: organizationId,
				},
			},
			PayAccountType: bankType,
		},
		CreateTime: createTimeParam,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if paymentReceiptList != nil && len(*paymentReceiptList) > 0 {
		for _, paymentReceipt := range *paymentReceiptList {

			//处理账户空格
			paymentReceipt.ReceiveAccount = strings.ReplaceAll(paymentReceipt.ReceiveAccount, " ", "")
			paymentReceipt.PayAccount = strings.ReplaceAll(paymentReceipt.PayAccount, " ", "")

			bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
				OrganizationId: paymentReceipt.OrganizationId,
				Type:           bankType,
				Account:        paymentReceipt.PayAccount,
			})
			if err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt查询账号详情失败:%v", err))
				continue
			}
			result, err := s.pinganBankSDK.QueryTransferResult(ctx, paymentReceipt.Code, bankAccount.ZuId, bankAccount.Account)
			zap.L().Info(fmt.Sprintf("pinganBankSDK.QueryTransferResult查询转账结果:%v", result))
			if err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt查询转账结果失败:%v", err))
				continue
			}

			// 处理订单状态(默认失败)
			//orderState := enum.GuilinBankTransferFailResult
			// 获取映射的订单状态
			result.Stt = enum.GetOrderState(result.Stt, result.Stt)
			//比较订单状态
			if result.Stt == paymentReceipt.OrderStatus {
				continue
			}
			/*updateReceipt := &repo.BankTransferReceiptDBData{
				OrderState: result.Stt,
			}*/
			updateReceipt := &repo.PaymentReceiptDBData{
				OrderStatus: result.Stt,
				OrderFlowNo: result.HostFlowNo,
			}
			if result.Stt != "" && result.Stt != enum.GuilinBankTransferSuccessResult {
				updateReceipt.RetCode = result.Yhcljg
				updateReceipt.RetMessage = result.BackRem
			}
			if err = s.paymentReceiptRepo.UpdateById(ctx, paymentReceipt.Id, updateReceipt); err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt更新paymentReceiptRepo失败:%v", err))
				continue
			}
		}
	}
	return nil
}

func (s *paymentReceiptService) SyncPaymentReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) (err error) {
	var syncDateRequest stru.SyncDateRequest
	json.Unmarshal(param, &syncDateRequest)
	if syncDateRequest.BeginDate == "" || syncDateRequest.EndDate == "" {
		return s.baseClient.FailTask(ctx, taskId, "BeginDate or EndDate is empty")
	}
	if err := s.HandleSyncPaymentReceipt(ctx, syncDateRequest.BeginDate, syncDateRequest.EndDate, organizationId); err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	return s.baseClient.SuccessTask(ctx, taskId, "")
}

func (s *paymentReceiptService) PaymentReceiptSystemApprove(ctx context.Context, id int64) (err error) {
	PaymentReceiptApplication, err := s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if PaymentReceiptApplication == nil || PaymentReceiptApplication.Id == 0 {
		return handler.HandleNewError("Payment Receipt application not exists")
	}
	var requestMap map[string]int64
	requestMap = make(map[string]int64)
	requestMap["processInstanceId"] = id
	err = s.SystemApprove(ctx, PaymentReceiptApplication.Id, requestMap)
	if err != nil {
		return handler.HandleError(err)
	}
	return nil
}

func (s *paymentReceiptService) PaymentReceiptSystemRefuse(ctx context.Context, id int64) (err error) {
	PaymentReceiptApplication, err := s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
		BaseProcessDBData: repository.BaseProcessDBData{
			ProcessInstanceId: id,
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if PaymentReceiptApplication == nil || PaymentReceiptApplication.Id == 0 {
		return handler.HandleNewError("Reimburse application not exists")
	}
	var requestMap map[string]int64
	requestMap = make(map[string]int64)
	requestMap["processInstanceId"] = id
	err = s.SystemRefuse(ctx, PaymentReceiptApplication.Id, requestMap)
	if err != nil {
		return handler.HandleError(err)
	}
	return nil
}
