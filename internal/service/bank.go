package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	bankEnum "gitlab.yoyiit.com/youyi/app-bank/internal/enum"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk"
	sdkStru "gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service/stru"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	baseApi "gitlab.yoyiit.com/youyi/app-base/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-base/kitex_gen/api/base"
	dingtalkApi "gitlab.yoyiit.com/youyi/app-dingtalk/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-dingtalk/kitex_gen/api/dingtalk"
	"gitlab.yoyiit.com/youyi/app-finance/kitex_gen/api/finance"
	flexApi "gitlab.yoyiit.com/youyi/app-flex/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-flex/kitex_gen/api/flex"
	oaApi "gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api/oa"
	"gitlab.yoyiit.com/youyi/go-common/enum"
	"gitlab.yoyiit.com/youyi/go-common/kafka"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/store"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BankService interface {
	ListBankTransferReceipt(ctx context.Context, req *api.ListBankTransferReceiptRequest) (*api.ListBankTransferReceiptResponse, error)
	GetBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (*api.BankTransferReceiptData, error)
	AddBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (int64, error)
	EditBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) error
	DeleteBankTransferReceipt(ctx context.Context, id int64) error
	CountBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (int64, error)

	ConfirmTransaction(ctx context.Context, req *api.BankTransferReceiptData) error
	GuilinBankTranscation(ctx context.Context, req *api.BankTransferReceiptData) error
	SPDBankTranscation(ctx context.Context, req *api.BankTransferReceiptData) error

	HandleTransferReceiptResult(ctx context.Context, id int64) error

	ListBankTransactionDetail(context.Context, *api.ListBankTransactionDetailRequest) (*api.ListBankTransactionDetailResponse, error)
	SimpleListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest) (*api.ListBankTransactionDetailResponse, error)
	GetBankTransactionDetail(context.Context, *api.BankTransactionDetailData) (*api.BankTransactionDetailData, error)
	SimpleGetBankTransactionDetail(context.Context, *api.BankTransactionDetailData) (*api.BankTransactionDetailData, error)
	CreateTransactionDetailProcessInstance(ctx context.Context, id int64) error
	EditBankTransactionDetailExtField(ctx context.Context, req *api.BankTransactionDetailData) error

	HandleTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	HandleGuilinBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error
	HandleSPDBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error
	HandlePinganBankTransactionDetail(ctx context.Context, bankType string, date string, date2 string, id int64) error

	HandleMinShengBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error
	ListBankTransactionDetailProcessInstance(ctx context.Context, id int64) ([]*api.BankTransactionDetailProcessInstanceData, error)

	GetBankCodeInfo(ctx context.Context, code string) (*api.BankCodeData, error)
	QueryBankCardInfo(ctx context.Context, cardNo string) (*api.QueryBankCardInfoResponse, error)
	ListBankCode(ctx context.Context, req *api.ListBankCodeRequest) (*api.ListBankCodeResponse, error)
	GetBankCode(ctx context.Context, req *api.BankCodeData) (*api.BankCodeData, error)
	AddBankCode(ctx context.Context, req *api.AddBankCodeRequest) error
	EditBankCode(ctx context.Context, req *api.BankCodeData) error
	DeleteBankCode(ctx context.Context, id int64) error

	HandleSyncTransferReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	HandleGuilinBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	HandleSPDBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	HandlePinganBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	UpdateBankTransactionRecDetail(context.Context, *api.BankTransactionRecDetailData) error

	SyncTransferReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) error
	SyncTransactionDetail(ctx context.Context, taskId int64, param []byte, organizationId int64) error

	DashboardData(ctx context.Context, organizationId int64) (*api.DashboardData, error)
	GetCashFlowMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (*api.ChartData, error)
	GetBalanceMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (*api.ChartData, error)

	QueryAccountBalance(ctx context.Context, req *api.QueryAccountBalanceRequest) (*api.QueryAccountBalanceResponse, error)
	QuerySPDBankAccountBalance(ctx context.Context, accountNo string, organizationBankConfig *baseApi.OrganizationBankConfigData) (*api.QueryAccountBalanceResponse, error)
	QueryGuilinBankAccountBalance(ctx context.Context, accountNo string, organizationBankConfig *baseApi.OrganizationBankConfigData) (*api.QueryAccountBalanceResponse, error)
	ImportBankBusinessPayrollData(ctx context.Context, taskId int64, param []byte, organizationId int64) error

	ListBankBusinessPayroll(ctx context.Context, req *api.ListBusinessPayrollRequest) (*api.ListBusinessPayrollResponse, error)
	ListBankBusinessPayrollDetail(ctx context.Context, req *api.ListBusinessPayrollDetailRequest) (*api.ListBusinessPayrollDetailResponse, error)
	SyncBankBusinessPayrollDetail(ctx context.Context, req *api.SyncBusinessPayrollResultRequest) (*api.SyncBusinessPayrollResultResponse, error)
	HandleTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	HandleSPDTransactionDetailReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	/*	HandlePinganTransactionDetailReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error
	 */
	CreateVirtualAccount(ctx context.Context, req *api.CreateVirtualAccountRequest) (*api.CreateVirtualAccountResponse, error)
	SyncVirtualAccountBalance(ctx context.Context) (err error)
	PinganSyncVirtualAccountBalance(ctx context.Context) (err error)
	QueryVirtualAccountBalance(ctx context.Context, organizationId int64, bankType string) (*api.VirtualAccountBalanceData, error)
	VirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error)
	PinganBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error)
	SPDBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error)
	PinganBankAccountSignatureApply(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (*api.PinganUserAcctSignatureApplyResponse, error)
	PinganBankAccountSignatureQuery(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (*api.PinganUserAcctSignatureApplyResponse, error)
	PinganBankVirtualSubAcctBalanceAdjust(ctx context.Context, id int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error)
	IcbcBankAccountSignatureQuery(ctx context.Context, req *api.IcbcBankAccountSignatureRequest) (*api.IcbcBankAccountSignatureQueryResponse, error)
	MinShengBankAccountSignatureApply(ctx context.Context, req *api.MinShengBankAccountSignatureRequest) (string, error)
	MinShengBankAccountSignatureQuery(ctx context.Context, req *api.MinShengBankAccountSignatureRequest) (*api.MinShengBankAccountSignatureQueryResponse, error)
	IcbcBankListTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	SyncIcbcBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	SyncMinshengBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error
	GetBankTransactionReceipt(ctx context.Context, bankTransactionDetailId int64) error
	SyncBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, bankType string) error
	PinganBankTransaction(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankAccountTranscationResponse, error)
}

type bankService struct {
	baseClient                               base.Client
	guilinBankSDK                            sdk.GuilinBankSDK
	spdBankSDK                               sdk.SPDBankSDK
	pinganBankSDK                            sdk.PinganBankSDK
	bankTransferReceiptRepo                  repo.BankTransferReceiptRepo
	kafkaProducer                            *store.KafkaProducer
	dingtalkClient                           dingtalk.Client
	bankTransactionDetailRepo                repo.BankTransactionDetailRepo
	bankTransactionDetailProcessInstanceRepo repo.BankTransactionDetailProcessInstanceRepo
	ossConfig                                *store.OSSConfig
	bankCodeRepo                             repo.BankCodeRepo
	businessPayrollRepo                      repo.BankBusinessPayrollRepo
	businessPayrollDetailRepo                repo.BankBusinessPayrollDetailRepo
	oaClient                                 oa.Client
	paymentReceiptRepo                       repo.PaymentReceiptRepo
	pdfToImageService                        PdfToImageService
	financeClient                            finance.Client
	icbcBank                                 sdk.IcbcBankSDK
	minShengBank                             sdk.MinShengSDK
	redisClient                              *store.RedisClient
	flexClient                               flex.Client
}

func (s *bankService) SyncBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64, bankType string) error {
	switch bankType {
	case enum.PinganBankType:
		s.SyncIcbcBankTransactionReceipt(ctx, beginDate, endDate, organizationId)
	case enum.MinShengBankType:
		s.SyncMinshengBankTransactionReceipt(ctx, beginDate, endDate, organizationId)

	}
	return nil
}
func (s *bankService) SyncIcbcBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	//先更新本地回单
	_, err := s.icbcBank.IcbcReceiptFileDownload(ctx)
	if err != nil {
		return handler.HandleError(err)
	}
	zap.L().Info(fmt.Sprintf("sSyncIcbcBankTransactionReceipt开始查询icbc没有回单的流水"))
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId:       organizationId,
		Type:                 enum.IcbcBankType,
		SignatureApplyStatus: "0",
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, bankAccount := range bankAccounts {
		//找到没有回单附件的流水
		datas, _, err := s.bankTransactionDetailRepo.List(ctx, "-1", 0, 1000, &repo.BankTransactionDetailDBDataParam{
			BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
				BaseDBData:        repository.BaseDBData{},
				MerchantAccountId: bankAccount.Id,
				TransferDate:      beginDate,
			},
			IsElectronicReceiptFileNull: false,
		})
		if err != nil {
			return err
		}
		zap.L().Info(fmt.Sprintf("==IcbcBankAccountListTransactionDetail%v", datas))
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.icbcBank.ListTransactionDetail__error_info%s", err.Error()))
			continue
		}
		if *datas != nil {
			for _, data := range *datas {
				//在临时回单文件中寻找匹配的回单文件,
				dir, err := ioutil.ReadDir(bankEnum.IcbcTempFilePath)
				if err != nil {
					return handler.HandleError(err)
				}
				for _, fileInfo := range dir {
					fileName := fileInfo.Name()
					if strings.Contains(fileName, bankAccount.Account) && strings.Contains(fileName, data.HostFlowNo) && strings.HasSuffix(fileName, ".pdf") {
						f, err := os.ReadFile(path.Join(bankEnum.IcbcTempFilePath, fileName))
						if err != nil {
							zap.L().Error(fmt.Sprintf("SyncIcbcBankTransactionReceipt读取icbc回单失败: %v\n", err.Error()))
						}
						electronicReceiptFile, err := store.UploadOSSFileBytes("pdf", ".pdf", f, s.ossConfig, false)
						if err != nil {
							zap.L().Error(fmt.Sprintf("SyncIcbcBankTransactionReceipt上传icbc电子凭证到OSS失败: %v\n", err.Error()))
						}
						err = s.bankTransactionDetailRepo.UpdateById(ctx, data.Id, &repo.BankTransactionDetailDBData{
							ElectronicReceiptFile: electronicReceiptFile,
						})
						if err != nil {
							zap.L().Error(fmt.Sprintf("更新icbc电子凭证失败: %v\n", err.Error()))
						}
						//todo 更新单据的明细回单
						break
					}
				}

			}
		}
	}
	return nil
}
func (s *bankService) SyncMinshengBankTransactionReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	//先更新本地回单
	zap.L().Info(fmt.Sprintf("sSyncMingshengBankTransactionReceipt开始查询minsheng没有回单的流水"))
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId: organizationId,
		Type:           enum.MinShengBankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, bankAccount := range bankAccounts {
		//找到没有回单附件的流水
		datas, _, err := s.bankTransactionDetailRepo.List(ctx, "-1", 0, 1000, &repo.BankTransactionDetailDBDataParam{
			BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
				BaseDBData:        repository.BaseDBData{},
				MerchantAccountId: bankAccount.Id,
				TransferDate:      beginDate,
			},
			IsElectronicReceiptFileNull: false,
		})

		zap.L().Info(fmt.Sprintf("==minshengBankAccountListTransactionDetail%v", datas))
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.ListTransactionDetail__error_info%s", err.Error()))
			continue
		}
		if *datas != nil {
			for _, data := range *datas {
				//下载回单并上传到oss==>通过jsdk服务
				s.minshengDownloadReceiptAndUpdatepPaymentReceipt(ctx, &data, bankAccount)
			}
		}
	}
	return nil
}

func (s *bankService) GetBankTransactionReceipt(ctx context.Context, bankTransactionDetailId int64) error {
	//将该id放入redis中,
	zap.L().Info(fmt.Sprintf("==GetBankTransactionReceipt开始执行回单下载任务::id=%d", bankTransactionDetailId))
	lockKey := bankEnum.GetBankTransactionReceipt + strconv.FormatInt(bankTransactionDetailId, 10)
	// 设置过期时间为20分钟
	lock, err := s.redisClient.SetNX(ctx, lockKey, time.Now(), 30*time.Minute).Result()
	if err != nil {
		zap.L().Info(fmt.Sprintf("==GetIBankTransactionReceipt设置redis锁失败: %+v", err))
		return errors.New(err.Error())
	}
	if lock {
		bankTransactionDetail, _ := s.bankTransactionDetailRepo.Get(ctx, &repo.BankTransactionDetailDBData{
			BaseDBData: repository.BaseDBData{
				BaseCommonDBData: repository.BaseCommonDBData{Id: bankTransactionDetailId},
			},
		})
		account, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
			OrganizationId: bankTransactionDetail.OrganizationId,
			Type:           bankTransactionDetail.PayAccountType,
		})
		if err != nil {
			return handler.HandleError(err)
		}
		if account.Type == enum.IcbcBankType {
			go func() {
				query, err := s.icbcBank.IcbcReceiptNoQuery(ctx, account.Account, account.ZuId, bankTransactionDetail.HostFlowNo, "")
				if err != nil {
					zap.L().Info(fmt.Sprintf("==GetBankTransactionReceipt查询回单orderid失败%s,err=%+v", query.OrderId, err))
				}
				zap.L().Info(fmt.Sprintf("==GetBankTransactionReceipt%s", query.OrderId))
				//等待20分钟之后继续执行
				time.Sleep(30 * time.Minute)
				zap.L().Info(fmt.Sprintf("==GetBankTransactionReceipt开始执行回单下载任务%s", query.OrderId))
				err = s.icbcBank.IcbcReceiptFileDownloadByOrderId(ctx, query.OrderId)
				if err != nil {
					zap.L().Info(fmt.Sprintf("==GetBankTransactionReceipt下载回单orderid失败%s,err=%+v", query.OrderId, err))
				}
				//寻找对应hostflow
				//在临时回单文件中寻找匹配的回单文件,
				dir, _ := ioutil.ReadDir(bankEnum.IcbcTempFilePath)
				for _, fileInfo := range dir {
					fileName := fileInfo.Name()
					if strings.Contains(fileName, bankTransactionDetail.HostFlowNo) && strings.HasSuffix(fileName, ".pdf") {
						f, err := os.ReadFile(path.Join(bankEnum.IcbcTempFilePath, fileName))
						if err != nil {
							zap.L().Error(fmt.Sprintf("SyncIcbcBankTransactionReceipt读取icbc回单失败: %v\n", err.Error()))
						}
						electronicReceiptFile, err := store.UploadOSSFileBytes("pdf", ".pdf", f, s.ossConfig, false)
						if err != nil {
							zap.L().Error(fmt.Sprintf("SyncIcbcBankTransactionReceipt上传icbc电子凭证到OSS失败: %v\n", err.Error()))
						}
						err = s.bankTransactionDetailRepo.UpdateById(ctx, bankTransactionDetailId, &repo.BankTransactionDetailDBData{
							ElectronicReceiptFile: electronicReceiptFile,
						})
						if err != nil {
							zap.L().Error(fmt.Sprintf("GetIcbcBankTransactionReceipt更新icbc电子凭证失败: %v\n", err.Error()))
						}
						break
					}
				}

			}()

		} else if account.Type == enum.SPDBankType {
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: bankTransactionDetail.OrganizationId,
				Type:           enum.SPDBankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			f, err := s.spdBankSDK.DownloadTransactionDetailElectronicReceipt(ctx, account.Account, bankTransactionDetail.TransferDate, bankTransactionDetail.TransferDate,
				bankTransactionDetail.OrderFlowNo, bankTransactionDetail.ExtField1, organizationBankConfig.Host, organizationBankConfig.SignHost,
				organizationBankConfig.FileHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			if err != nil {
				zap.L().Info(fmt.Sprintf("s.bankService.GetBankTransactionReceipt 下载浦发电子凭证失败: %v\n", err.Error()))
			} else {
				str := fmt.Sprintf("bankAccount.Account=%s,newDbData.TransferDate=%s,newDbData.TransferDate=%s,newDbData.OrderFlowNo=%s,newDbData.ExtField1=%s,organizationBankConfig.Host=%s,organizationBankConfig.SignHost=%s,organizationBankConfig.FileHost=%s,organizationBankConfig.BankCustomerId=%s,organizationBankConfig.BankUserId=%s",
					account.Account, bankTransactionDetail.TransferDate, bankTransactionDetail.TransferDate,
					bankTransactionDetail.OrderFlowNo, bankTransactionDetail.ExtField1, organizationBankConfig.Host, organizationBankConfig.SignHost,
					organizationBankConfig.FileHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
				value := fmt.Sprintf("%sid=%d单独下载浦发电子凭证成功参数=%s", util.FormatDateTime(time.Now()), bankTransactionDetail.Id, str)
				s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
			}
			var electronicReceiptFile string
			if f != nil && len(f) > 0 {
				electronicReceiptFile, err = store.UploadOSSFileBytes("pdf", ".pdf", f, s.ossConfig, false)
				if err != nil {
					zap.L().Error(fmt.Sprintf("上传浦发电子凭证到OSS失败: %v\n", err.Error()))
				}
				err := s.bankTransactionDetailRepo.UpdateById(ctx, bankTransactionDetail.Id, &repo.BankTransactionDetailDBData{
					ElectronicReceiptFile: electronicReceiptFile,
				})
				if err != nil {
					zap.L().Error(fmt.Sprintf("更新浦发电子凭证失败: %v\n", err.Error()))
				}
				//更新单据的明细ID
				if err = s.updateRelevanceElectronicDocument(ctx, bankTransactionDetail.Summary, bankTransactionDetail.HostFlowNo, electronicReceiptFile, enum.SPDBankType); err != nil {
					zap.L().Error(fmt.Sprintf("更新浦发更新单据的回单失败: %v\n", err.Error()))
				}
				zap.L().Info("GetBankTransactionReceipt 处理浦发电子凭证成功")
			}
		} else if account.Type == enum.MinShengBankType {
			go func() {
				s.minshengDownloadReceiptAndUpdatepPaymentReceipt(ctx, bankTransactionDetail, account)

				receiptResponse, err := s.minShengBank.GetTransactionDetailElectronicReceipt(ctx, account.Account, bankTransactionDetail.OrderFlowNo, bankTransactionDetail.TransferDate, account.OpenId)
				if err != nil {
					zap.L().Info(fmt.Sprintf("GetBankTransactionReceipt 处理民生电子凭证失败,err=%+v", err))
				}
				if receiptResponse["return_code"] != "0000" {
					zap.L().Info(fmt.Sprintf("GetBankTransactionReceipt 处理民生电子凭证失败查询回单结果:%+v", receiptResponse))
				}
				electronicReceiptFile := receiptResponse["response_busi"].(string)
				err = s.bankTransactionDetailRepo.UpdateById(ctx, bankTransactionDetail.Id, &repo.BankTransactionDetailDBData{
					ElectronicReceiptFile: electronicReceiptFile,
				})
				if err != nil {
					zap.L().Error(fmt.Sprintf("GetBankTransactionReceipt-更新民生电子凭证失败: %v\n", err.Error()))
				}
				//更新单据的电子凭证
				summary := bankTransactionDetail.Summary
				//可,
				if summary != "" && strings.Index(summary, "[") >= 0 && strings.Index(summary, "]") >= 0 {
					serialNo := summary[strings.Index(summary, "[")+1 : strings.Index(summary, "]")]
					paymentReceipt, err := s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
						Code: serialNo,
					})
					if err != nil {
						zap.L().Error(fmt.Sprintf("GetBankTransactionReceipt-更新民生单据凭证失败: %v\n", err.Error()))
					}
					if paymentReceipt != nil {
						paymentReceipt.ElectronicDocument = electronicReceiptFile
					}
					s.paymentReceiptRepo.UpdateById(ctx, bankTransactionDetail.Id, paymentReceipt)
				}
			}()

		}
	} else {
		zap.L().Info(fmt.Sprintf("==GetIBankTransactionReceipt申请回单重复,此任务不在执行::id=%d", bankTransactionDetailId))
		ttl, _ := s.redisClient.TTL(ctx, lockKey).Result()
		sc := int(ttl.Seconds())
		return errors.New(fmt.Sprintf("该流水已申请电子回单,请等待%d分钟%d秒后刷新页面查看结果", sc/60, sc%60))
	}

	return nil
}

func (s *bankService) minshengDownloadReceiptAndUpdatepPaymentReceipt(ctx context.Context, bankTransactionDetail *repo.BankTransactionDetailDBData, account *baseApi.OrganizationBankAccountData) {
	receiptResponse, err := s.minShengBank.GetTransactionDetailElectronicReceipt(ctx, account.Account, bankTransactionDetail.OrderFlowNo, bankTransactionDetail.TransferDate, account.OpenId)
	if err != nil {
		zap.L().Info(fmt.Sprintf("GetBankTransactionReceipt 处理民生电子凭证失败,err=%+v", err))
	}
	if receiptResponse["return_code"] != "0000" {
		zap.L().Info(fmt.Sprintf("GetBankTransactionReceipt 处理民生电子凭证失败查询回单结果:%+v", receiptResponse))
	}
	electronicReceiptFile := receiptResponse["response_busi"].(string)
	err = s.bankTransactionDetailRepo.UpdateById(ctx, bankTransactionDetail.Id, &repo.BankTransactionDetailDBData{
		ElectronicReceiptFile: electronicReceiptFile,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("GetBankTransactionReceipt-更新民生电子凭证失败: %v\n", err.Error()))
	}
	//更新单据的电子凭证
	summary := bankTransactionDetail.Summary
	//可,
	if summary != "" && strings.Index(summary, "[") >= 0 && strings.Index(summary, "]") >= 0 {
		serialNo := summary[strings.Index(summary, "[")+1 : strings.Index(summary, "]")]
		paymentReceipt, err := s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
			Code: serialNo,
		})
		if err != nil {
			zap.L().Error(fmt.Sprintf("GetBankTransactionReceipt-更新民生单据凭证失败: %v\n", err.Error()))
		}
		if paymentReceipt != nil {
			paymentReceipt.ElectronicDocument = electronicReceiptFile
		}
		s.paymentReceiptRepo.UpdateById(ctx, bankTransactionDetail.Id, paymentReceipt)
	}
}

func (s *bankService) IcbcBankListTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId:       organizationId,
		Type:                 enum.IcbcBankType,
		SignatureApplyStatus: "0",
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, bankAccount := range bankAccounts {
		datas, err := s.icbcBank.ListTransactionDetail(ctx, bankAccount.Account, beginDate, endDate, bankAccount.ZuId)
		if err != nil {
			return err
		}
		zap.L().Info(fmt.Sprintf("==IcbcBankAccountListTransactionDetail%v", datas))
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.icbcBank.ListTransactionDetail__error_info%s", err.Error()))
			continue
		}
		if datas != nil {
			var addDatas []repo.BankTransactionDetailDBData
			for _, data := range datas {
				count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
					BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						MerchantAccountId: bankAccount.Id,
						HostFlowNo:        strconv.FormatInt(data.SerialNo, 10), //流水号
					},
				})
				if err != nil {
					continue
				}
				// 保存交易明细
				if count == 0 {
					//  D借，出账；C贷，入账
					//1-借，2-贷，借是出账，贷是入账-icbc
					payAmount := 0.00
					recAmount := 0.00
					TransactionType := ""
					tranAmount := data.Amount
					if err != nil {
						continue
					}
					if data.DrcrF == 2 {
						recAmount = tranAmount / 100
						TransactionType = enum.GuilinBankTransactionDetailRecType

					} else if data.DrcrF == 1 {
						payAmount = tranAmount / 100
						TransactionType = enum.GuilinBankTransactionDetailPayType

					}
					acctBalance := data.Balance / 100
					if err != nil {
						continue
					}
					transferDate := strings.Replace(data.BusiDate, "-", "", -1)
					transferTime := strings.Replace(data.BusiTime, ".", "", -1)
					transactionDetailDBData := repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						Type:                TransactionType,
						MerchantAccountId:   bankAccount.Id,
						MerchantAccountName: bankAccount.AccountName,
						CashFlag:            "0", // 写死0: 现钞
						PayAmount:           payAmount,
						RecAmount:           recAmount,
						BsnType:             "TR", // 写死TR: 转账
						TransferDate:        transferDate,
						TransferTime:        transferDate + transferTime,
						TranChannel:         "",
						CurrencyType:        "CNY",
						Balance:             acctBalance,
						//OrderFlowNo:         data.BussSeqNo,
						OrderFlowNo:        strconv.FormatInt(data.TrxCode, 10),
						HostFlowNo:         strconv.FormatInt(data.SerialNo, 10),
						VouchersType:       strconv.Itoa(data.VouhType),
						VouchersNo:         strconv.Itoa(data.VouhNo),
						SummaryNo:          data.FinFo,
						Summary:            data.Summary,
						AcctNo:             data.RecipAcc,
						AccountName:        data.RecipNam,
						AccountOpenNode:    data.RecipBna,
						ProcessTotalStatus: enum.ProcessInstanceTotalStatusRunning,
						PayAccountType:     enum.IcbcBankType,
					}

					//查询回单
					//f, err := s.pinganBankSDK.UploadTransactionDetailElectronic(ctx, request, bankAccount.ZuId)
					//if err != nil {
					//	zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证失败: %v\n", err.Error()))
					//}
					//if f != "" {
					//	transactionDetailDBData.ElectronicReceiptFile = f
					//}
					addDatas = append(addDatas, transactionDetailDBData)
					//
					////更新单据的回单
					//if err = s.updateRelevanceElectronicDocument(ctx, "", data.HostTrace, f, enum.PinganBankType); err != nil {
					//	return handler.HandleError(err)
					//}
				}
			}
			if len(addDatas) > 0 {

				ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDatas)
				if err != nil {
					return handler.HandleError(err)
				}
				for _, id := range ids {
					s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
						Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
						Type:     kafka.DingtalkType,
						Id:       id,
					})
				}
			}
		}
	}
	return nil
}

func (s *bankService) IcbcBankAccountSignatureQuery(ctx context.Context, req *api.IcbcBankAccountSignatureRequest) (*api.IcbcBankAccountSignatureQueryResponse, error) {
	bankAccount, _ := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{Id: req.Id})

	agreement, err := s.icbcBank.QueryAgreeNo(ctx, bankAccount.ZuId, bankAccount.Account)

	if err != nil {
		return nil, err
	}

	status := "1"
	agreeNo := strconv.FormatInt(agreement.AgreeNo, 10)
	if agreement.Status == "1" {
		status = "0"
	}
	s.baseClient.EditOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id:                   req.Id,
		ZuId:                 agreeNo,
		SignatureApplyStatus: status,
		Remark:               "",
	})
	return &api.IcbcBankAccountSignatureQueryResponse{
		Signatureapplystatus: status,
		ZuId:                 agreeNo,
		Remark:               bankAccount.Remark,
	}, nil
}

func (s *bankService) GetBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (*api.BankTransferReceiptData, error) {
	dbData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if "" != dbData.DetailHostFlowNo && dbData.PayAccountType == enum.GuilinBankType {
		detailData, err := s.bankTransactionDetailRepo.Get(ctx, &repo.BankTransactionDetailDBData{
			HostFlowNo: dbData.DetailHostFlowNo,
		})
		if err == nil && nil != detailData {
			dbData.ElectronicReceiptFile = detailData.ElectronicReceiptFile
		}
	}
	return stru.ConvertBankTransferReceiptData(*dbData), nil
}

func (s *bankService) ListBankTransferReceipt(ctx context.Context, req *api.ListBankTransferReceiptRequest) (*api.ListBankTransferReceiptResponse, error) {
	dbData, count, err := s.bankTransferReceiptRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankTransferReceiptDBDataParam{
		ProcessInstanceIds: req.ProcessInstanceIds,
		SerialNo:           req.SerialNo,
		BusinessId:         req.BusinessId,
		PayAccount:         req.PayAccount,
		RecAccount:         req.RecAccount,
		PayAmount:          req.PayAmount,
		OriginatorUser:     req.OriginatorUser,
		CommentUser:        req.CommentUser,
		CreateTime:         req.CreateTimeArray,
		TotalStatus:        req.TotalStatus,
		OrderState:         req.OrderState,
		PayAmountMin:       req.PayAmountMin,
		PayAmountMax:       req.PayAmountMax,
		BankTransferReceiptDBData: repo.BankTransferReceiptDBData{
			Title: req.Title,
		},
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	data := make([]*api.BankTransferReceiptData, len(*dbData))
	for i, v := range *dbData {
		data[i] = stru.ConvertBankTransferReceiptData(v)
	}
	return &api.ListBankTransferReceiptResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) AddBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (int64, error) {
	id, err := s.bankTransferReceiptRepo.Add(ctx, stru.ConvertBankTransferReceiptDBData(*req))
	if err != nil {
		return 0, handler.HandleError(err)
	}
	return id, nil
}

func (s *bankService) EditBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) error {
	data := stru.ConvertBankTransferReceiptDBData(*req)
	currentUserId, err := util.GetMetaInfoCurrentUserId(ctx)
	if err != nil {
		return handler.HandleError(err)
	}
	currentUserName, err := util.GetMetaInfoCurrentUserName(ctx)
	if err != nil {
		return handler.HandleError(err)
	}
	data.CommentUserId = currentUserId
	data.CommentUserName = currentUserName
	err = s.bankTransferReceiptRepo.UpdateById(ctx, req.Id, data)
	if err != nil {
		return handler.HandleError(err)
	}
	s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
		Business: kafka.ProcessFinanceTransferResultBusiness,
		Type:     kafka.DingtalkType,
		Id:       req.Id,
	})
	return nil
}

func (s *bankService) DeleteBankTransferReceipt(ctx context.Context, id int64) (err error) {
	return s.bankTransferReceiptRepo.DeleteById(ctx, id)
}

func (s *bankService) CountBankTransferReceipt(ctx context.Context, req *api.BankTransferReceiptData) (int64, error) {
	return s.bankTransferReceiptRepo.Count(ctx, &repo.BankTransferReceiptDBDataParam{
		BankTransferReceiptDBData: *stru.ConvertBankTransferReceiptDBData(*req),
	})
}

func (s *bankService) ConfirmTransaction(ctx context.Context, req *api.BankTransferReceiptData) error {
	// 根据 payAccountType 获取银行类型 0: 桂林银行 1: 浦发银行,2:平安银行
	payAccountType := req.PayAccountType
	if payAccountType == enum.GuilinBankType {
		return s.GuilinBankTranscation(ctx, req)
	} else if payAccountType == enum.SPDBankType {
		return s.SPDBankTranscation(ctx, req)
	} else if payAccountType == enum.PinganBankType {
		return s.PinganBankTranscation(ctx, req)
	}
	return nil
}

// GuilinBankTranscation
//
//	@Description: 桂林银行转账接口
//	@receiver s
//	@param ctx
//	@param req
//	@return error
func (s *bankService) GuilinBankTranscation(ctx context.Context, req *api.BankTransferReceiptData) error {
	dbData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})

	// 0 行内 1 行外
	recBankType := "0"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		recBankType = "1"
	}

	if err != nil {
		return handler.HandleError(err)
	}
	if req.PayAccount == "" {
		return errors.New("PayAccount must not empty")
	}
	if req.PayAccountName == "" {
		return errors.New("PayAccountName must not empty")
	}
	if dbData.RecAccount == "" {
		return errors.New("RecAccount must not empty")
	}
	if dbData.RecAccountName == "" {
		return errors.New("RecAccountName must not empty")
	}
	if dbData.PayAmount <= 0 {
		return errors.New("PayAmount must greater than 0")
	}
	if dbData.PubPriFlag == "" {
		return errors.New("PubPriFlag must not empty")
	}
	if dbData.CurrencyType == "" {
		return errors.New("CurrencyType must not empty")
	}
	if recBankType == "" {
		return errors.New("RecBankType must not empty")
	}
	if recBankType == "1" {
		if req.UnionBankNo == "" {
			return errors.New("UnionBankNo must not empty")
		}
		if req.ClearBankNo == "" {
			return errors.New("ClearBankNo must not empty")
		}
		if req.RecAccountOpenBank == "" {
			return errors.New("RecAccountOpenBank must not empty")
		}
		if dbData.RmtType == "" {
			return errors.New("RmtType must not empty")
		}
	}

	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: dbData.OrganizationId,
		Type:           enum.GuilinBankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	payRem := fmt.Sprintf("%s[%s]", dbData.PayRem, dbData.SerialNo)
	var bankTransferResponse sdkStru.BankTransferResponse
	switch recBankType {
	case "0":
		res, err := s.guilinBankSDK.IntrabankTransfer(ctx, dbData.SerialNo, sdkStru.IntrabankTransferRequest{
			PayAccount:     req.PayAccount,
			PayAccountName: req.PayAccountName,
			RecAccount:     dbData.RecAccount,
			RecAccountName: dbData.RecAccountName,
			PayAmount:      dbData.PayAmount,
			PayRem:         payRem,
			PubPriFlag:     dbData.PubPriFlag,
			CurrencyType:   dbData.CurrencyType,
			RecBankType:    recBankType,
		}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			return handler.HandleError(err)
		}
		if res != nil {
			bankTransferResponse = *res
		}
	case "1":
		res, err := s.guilinBankSDK.OutofbankTransfer(ctx, dbData.SerialNo, sdkStru.OutofbankTransferRequest{
			PayAccount:         req.PayAccount,
			PayAccountName:     req.PayAccountName,
			RecAccount:         dbData.RecAccount,
			RecAccountName:     dbData.RecAccountName,
			PayAmount:          dbData.PayAmount,
			PayRem:             payRem,
			PubPriFlag:         dbData.PubPriFlag,
			CurrencyType:       dbData.CurrencyType,
			RecBankType:        recBankType,
			TransferFlag:       dbData.TransferFlag,
			UnionBankNo:        req.UnionBankNo,
			ClearBankNo:        req.ClearBankNo,
			RecAccountOpenBank: req.RecAccountOpenBank,
			RmtType:            dbData.RmtType,
		}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			return handler.HandleError(err)
		}
		if res != nil {
			bankTransferResponse = *res
		}
	}
	if bankTransferResponse.Head.RetCode != enum.GuilinBankSuccessRetCode {
		bankTransferResponse.Body.OrderState = enum.GuilinBankTransferFailResult
	}

	if err = s.EditBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		Id:                 dbData.Id,
		PayAccount:         req.PayAccount,
		PayAccountName:     req.PayAccountName,
		ChargeFee:          bankTransferResponse.Body.ChargeFee,
		OrderState:         bankTransferResponse.Body.OrderState,
		RetCode:            bankTransferResponse.Head.RetCode,
		RetMessage:         bankTransferResponse.Head.RetMessage,
		OrderFlowNo:        bankTransferResponse.Head.OrderFlowNo,
		ProcessComment:     req.ProcessComment,
		PayRem:             payRem,
		CommentUserName:    req.CommentUserName,
		RecBankType:        recBankType,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.GuilinBankType,
	}); err != nil {
		return handler.HandleError(err)
	}
	return nil
}

// SPDBankTranscation
//
//	@Description: 浦发银行转账接口
//	@receiver s
//	@param ctx
//	@param req
//	@return error
func (s *bankService) SPDBankTranscation(ctx context.Context, req *api.BankTransferReceiptData) error {
	dbData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: dbData.OrganizationId,
		Type:           enum.SPDBankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}

	if req.PayAccount == "" {
		return errors.New("PayAccount must not empty")
	}
	if req.PayAccountName == "" {
		return errors.New("PayAccountName must not empty")
	}
	if dbData.RecAccount == "" {
		return errors.New("RecAccount must not empty")
	}
	if dbData.RecAccountName == "" {
		return errors.New("RecAccountName must not empty")
	}
	if dbData.PayAmount <= 0 {
		return errors.New("PayAmount must greater than 0")
	}
	if dbData.PubPriFlag == "" {
		return errors.New("PubPriFlag must not empty")
	}
	if dbData.CurrencyType == "" {
		return errors.New("CurrencyType must not empty")
	}
	// 0 行内 1 行外
	sysFlag := "0"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		sysFlag = "1"
	}
	if sysFlag == "" {
		return errors.New("SysFlag must not empty")
	}
	// 收款行名称 (当SysFlag=1即跨行转帐时必须输入)
	payeeBankName := ""
	// 同城异地标志
	remitLocation := "0"
	// 收款行速选标志 (1-速选, 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效。如果希望跨行汇款自动处理，请务必填写此项。
	payeeBankSelectFlag := ""
	// 收款行行号（人行现代支付系统行号）如果速选标志为1，请务必填写此项。 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效
	payeeBankNo := ""
	if sysFlag == "1" {
		if req.UnionBankNo == "" {
			return errors.New("UnionBankNo must not empty")
		}
		if req.ClearBankNo == "" {
			return errors.New("ClearBankNo must not empty")
		}
		if req.RecAccountOpenBank == "" {
			return errors.New("RecAccountOpenBank must not empty")
		}
		// 由于分行同城已经取消，跨行支付时，同城异地标志建议固定送1异地
		remitLocation = "1"
		payeeBankName = req.RecAccountOpenBank
		payeeBankSelectFlag = "1"
		payeeBankNo = req.UnionBankNo
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	dbData.SerialNo, _ = util.SonyflakeID()
	// 多次支付会有,例如: 转账测试[123][222][456], 这里要截取出前面的转账测试
	if strings.Contains(dbData.PayRem, "[") {
		dbData.PayRem = dbData.PayRem[:strings.Index(dbData.PayRem, "[")]
	}
	note := fmt.Sprintf("%s[%s]", dbData.PayRem, dbData.SerialNo)
	request := sdkStru.SPDBankTransferRequest{
		AuthMasterID:        organizationBankConfig.BankCustomerId,
		ElecChequeNo:        dbData.SerialNo,
		AcctNo:              req.PayAccount,
		AcctName:            req.PayAccountName,
		PayeeAcctNo:         dbData.RecAccount,
		PayeeName:           dbData.RecAccountName,
		PayeeBankName:       payeeBankName,
		Amount:              dbData.PayAmount,
		SysFlag:             sysFlag,
		RemitLocation:       remitLocation,
		Note:                note,
		PayeeBankSelectFlag: payeeBankSelectFlag,
		PayeeBankNo:         payeeBankNo,
	}
	bankTransferResponse, err := s.spdBankSDK.BankTransfer(ctx, dbData.SerialNo, request,
		organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return handler.HandleError(err)
	}
	// 处理订单状态
	orderState := bankTransferResponse.TransStatus
	if bankTransferResponse.TransStatus != "" {
		orderState = enum.GetOrderState(bankTransferResponse.TransStatus, orderState)
		bankTransferResponse.ReturnCode = bankTransferResponse.TransStatus
		bankTransferResponse.ReturnMsg = enum.GetOrderMsg(bankTransferResponse.TransStatus, bankTransferResponse.ReturnMsg)
	}

	if err = s.EditBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		Id:                 dbData.Id,
		SerialNo:           dbData.SerialNo,
		PayAccount:         req.PayAccount,
		PayAccountName:     req.PayAccountName,
		ChargeFee:          0.00, // 目前没看到浦发返回的手续费字段
		OrderState:         orderState,
		RetCode:            bankTransferResponse.ReturnCode,
		RetMessage:         bankTransferResponse.ReturnMsg,
		OrderFlowNo:        bankTransferResponse.AcceptNo,
		ProcessComment:     req.ProcessComment,
		PayRem:             note,
		CommentUserName:    req.CommentUserName,
		RecBankType:        sysFlag,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.SPDBankType,
	}); err != nil {
		return handler.HandleError(err)
	}
	return nil
}

func (s *bankService) PinganBankTranscation(ctx context.Context, req *api.BankTransferReceiptData) error {
	dbData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}

	if req.PayAccount == "" {
		return errors.New("PayAccount must not empty")
	}
	if req.PayAccountName == "" {
		return errors.New("PayAccountName must not empty")
	}
	if dbData.RecAccount == "" {
		return errors.New("RecAccount must not empty")
	}
	if dbData.RecAccountName == "" {
		return errors.New("RecAccountName must not empty")
	}
	if dbData.PayAmount <= 0 {
		return errors.New("PayAmount must greater than 0")
	}
	if dbData.PubPriFlag == "" {
		return errors.New("PubPriFlag must not empty")
	}
	if dbData.CurrencyType == "" {
		return errors.New("CurrencyType must not empty")
	}
	// 0 行内 1 行外  1：行内转账，0：跨行转账
	sysFlag := "1"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		sysFlag = "0"
	}
	if sysFlag == "" {
		return errors.New("SysFlag must not empty")
	}
	// 收款行名称 (当SysFlag=1即跨行转帐时必须输入)
	payeeBankName := ""
	// 同城异地标志
	remitLocation := "1"
	// 收款行行号（人行现代支付系统行号）如果速选标志为1，请务必填写此项。 当本行/他行标志SysFlag为“1”（他行）、同城异地标志remitLocation为“1”（异地）时才能生效
	payeeBankNo := ""
	if sysFlag == "0" {
		if req.UnionBankNo == "" {
			return errors.New("UnionBankNo must not empty")
		}
		if req.ClearBankNo == "" {
			return errors.New("ClearBankNo must not empty")
		}
		if req.RecAccountOpenBank == "" {
			return errors.New("RecAccountOpenBank must not empty")
		}
		payeeBankName = req.RecAccountOpenBank
		payeeBankNo = req.UnionBankNo
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	dbData.SerialNo, _ = util.SonyflakeID()
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Account: req.PayAccount,
		Type:    enum.PinganBankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	note := fmt.Sprintf("%s[%s]", dbData.PayRem, dbData.SerialNo)
	uuid, _ := util.SonyflakeID()
	request := sdkStru.PingAnBankTransferRequest{
		MrchCode:       config.GetString(bankEnum.PinganMrchCode, ""),
		CnsmrSeqNo:     uuid,
		OutAcctNo:      req.PayAccount,
		OutAcctName:    req.PayAccountName,
		ThirdVoucher:   dbData.SerialNo,
		CcyCode:        "RMB",
		InAcctNo:       dbData.RecAccount,
		InAcctName:     dbData.RecAccountName,
		InAcctBankName: payeeBankName,
		TranAmount:     strconv.FormatFloat(dbData.PayAmount, 'f', -1, 64),
		UnionFlag:      sysFlag,
		AddrFlag:       remitLocation,
		UseEx:          note,
		InAcctBankNode: payeeBankNo,
	}
	transfer, err := s.pinganBankSDK.BankTransfer(ctx, request, bankAccount.ZuId)

	zap.L().Info(fmt.Sprintf("s.pinganBankSDK.BankTransfer.transfer.HostFlowNo银行流水号3:%+v", transfer.HostFlowNo))
	if err != nil {
		return handler.HandleError(err)
	}
	// 处理订单状态
	chargeF, err := strconv.ParseFloat(transfer.Fee1, 64)
	if err != nil {
		return handler.HandleError(err)
	}
	if err = s.EditBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		Id:                 dbData.Id,
		SerialNo:           dbData.SerialNo,
		PayAccount:         req.PayAccount,
		PayAccountName:     req.PayAccountName,
		ChargeFee:          chargeF,
		OrderState:         enum.GetOrderState(transfer.Stt, transfer.Stt),
		RetCode:            transfer.Stt,
		RetMessage:         transfer.HostTxDate,
		OrderFlowNo:        transfer.HostFlowNo,
		ProcessComment:     req.ProcessComment,
		PayRem:             note,
		CommentUserName:    req.CommentUserName,
		RecBankType:        sysFlag,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.PinganBankType,
		DetailHostFlowNo:   transfer.FrontLogNo,
	}); err != nil {
		zap.L().Info(fmt.Sprintf("EditBankTransferReceipt-error:%v", err))
		return handler.HandleError(err)
	}
	return nil
}

func (s *bankService) HandleTransferReceiptResult(ctx context.Context, id int64) error {
	bankTransferDBData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: id,
			},
		},
	})
	if err != nil {
		return handler.HandleError(err)
	}
	organizationUnionConfig, err := s.baseClient.GetOrganizationUnionConfig(ctx, &baseApi.OrganizationUnionConfigData{
		OrganizationId: bankTransferDBData.OrganizationId,
		Type:           enum.DingtalkType,
	})
	if err != nil {
		return err
	}
	processInstance, err := s.dingtalkClient.GetProcessInstance(ctx, &dingtalkApi.ProcessInstanceData{
		Id: bankTransferDBData.ProcessInstanceId,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	var statusBg string
	var statusValue string
	var totalStatus string
	switch bankTransferDBData.OrderState {
	case enum.GuilinBankTransferSuccessResult:
		statusBg = "67c23a"
		statusValue = "交易成功"
		totalStatus = enum.ProcessInstanceTotalStatusSuccess
	case enum.GuilinBankTransferFailResult:
		statusBg = "e6a23c"
		statusValue = "交易失败"
		totalStatus = enum.ProcessInstanceTotalStatusFail
	case enum.GuilinBankTransferRefuseResult:
		statusBg = "f56c6c"
		statusValue = "拒绝"
		totalStatus = enum.ProcessInstanceTotalStatusRefuse
	case enum.GuilinBankTransferRejectResult:
		statusBg = "ff9900"
		statusValue = "落地拒绝"
		totalStatus = enum.ProcessInstanceTotalStatusRefuse
	case enum.GuilinBankTransferCreditResult:
		statusBg = "ff9900"
		statusValue = "有贷户落地"
		totalStatus = enum.ProcessInstanceTotalStatusRefuse
	case enum.GuilinBankTransferRevokeResult:
		statusBg = "cc6633"
		statusValue = "交易被撤销"
		totalStatus = enum.ProcessInstanceTotalStatusRefuse
	case enum.GuilinBankTransferDeleteResult:
		statusBg = "3399cc"
		statusValue = "交易作废"
		totalStatus = enum.ProcessInstanceTotalStatusRefuse
	}

	if err = s.dingtalkClient.EditProcessInstance(ctx, &dingtalkApi.ProcessInstanceData{
		Id:                   processInstance.Id,
		TransferReceiptState: bankTransferDBData.OrderState,
		TotalStatus:          totalStatus,
	}); err != nil {
		return handler.HandleError(err)
	}

	//更新付款单据的审批实例总状态
	if err = s.bankTransferReceiptRepo.UpdateById(ctx, id, &repo.BankTransferReceiptDBData{
		ProcessStatus: totalStatus,
	}); err != nil {
		return handler.HandleError(err)
	}
	if statusValue != "" {
		if err = s.dingtalkClient.AddProcessComment(ctx, &dingtalkApi.ProcessCommentRequest{
			ProcessInstanceId: processInstance.ExternalId,
			Text:              "付款单据" + statusValue + "，单据号：" + bankTransferDBData.SerialNo,
			UserId:            organizationUnionConfig.AdminUnionUserId,
		}, organizationUnionConfig.AppKey, organizationUnionConfig.AppSecret); err != nil {
			return handler.HandleError(err)
		}
		return s.dingtalkClient.SendOACorpConversation(ctx, bankTransferDBData.ProcessInstanceId,
			&dingtalkApi.ProcessOAMessageRequest{
				Url:         fmt.Sprintf("https://aflow.dingtalk.com/dingtalk/mobile/homepage.htm?corpid=%s&dd_share=false&showmenu=false&dd_progress=false&back=native&procInstId=%s&wfrom=onebox&dinghash=approval&dtaction=os&dd_from=onebox#approval", processInstance.CorpId, processInstance.ExternalId),
				Title:       "付款单据结果：" + statusValue,
				StatusBg:    statusBg,
				StatusValue: "结果: " + statusValue,
				Items: []*dingtalkApi.ProcessOAMessageItem{
					{Key: "流水号: ", Value: bankTransferDBData.SerialNo},
					{Key: "审批实例发起人: ", Value: bankTransferDBData.OriginatorUserName},
					{Key: "付款账号: ", Value: bankTransferDBData.PayAccount},
					{Key: "付款方户名: ", Value: bankTransferDBData.PayAccountName},
					{Key: "收款账号: ", Value: bankTransferDBData.RecAccount},
					{Key: "收款方户名: ", Value: bankTransferDBData.RecAccountName},
					{Key: "交易金额（元）: ", Value: strconv.FormatFloat(bankTransferDBData.PayAmount, 'f', 2, 64)},
				},
			}, organizationUnionConfig.AppKey, organizationUnionConfig.AppSecret, organizationUnionConfig.UnionAgentId)
	}
	return nil
}

func (s *bankService) ListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest) (*api.ListBankTransactionDetailResponse, error) {
	dbData, count, err := s.bankTransactionDetailRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankTransactionDetailDBDataParam{
		BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
			Type:                req.Type,
			ProcessBusinessId:   req.BusinessId,
			OriginatorUserName:  req.OriginatorUser,
			OperationUserName:   req.OperationUser,
			ProcessTotalStatus:  req.TotalStatus,
			HostFlowNo:          req.SerialNo,
			MerchantAccountName: req.MerchantAccountName,
			AccountName:         req.AccountName,
			BsnType:             req.BusinessType,
			PayAccountType:      req.PayAccountType,
			ExtField3:           req.ExtField3,
			MerchantAccountId:   req.MerchantAccountId,
		},
		PayAmountMin:      req.PayAmountMin,
		PayAmountMax:      req.PayAmountMax,
		RecAmountMin:      req.RecAmountMin,
		RecAmountMax:      req.RecAmountMax,
		TransferTimeArray: req.TransferTimeArray,
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	data := make([]*api.BankTransactionDetailData, len(*dbData))
	for i, v := range *dbData {
		data[i] = stru.ConvertBankTransactionDetailData(v)
	}
	return &api.ListBankTransactionDetailResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) SimpleListBankTransactionDetail(ctx context.Context, req *api.ListBankTransactionDetailRequest) (*api.ListBankTransactionDetailResponse, error) {
	dbData, count, err := s.bankTransactionDetailRepo.SimpleList(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankTransactionDetailDBDataParam{
		BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: req.OrganizationId,
			},
			Type: req.Type,
		},
		TransferTimeArray: req.TransferTimeArray,
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	data := make([]*api.BankTransactionDetailData, len(*dbData))
	for i, v := range *dbData {
		data[i] = stru.ConvertBankTransactionDetailData(v)
	}
	return &api.ListBankTransactionDetailResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) GetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData) (*api.BankTransactionDetailData, error) {
	dbData, err := s.bankTransactionDetailRepo.Get(ctx, &repo.BankTransactionDetailDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})

	if err != nil {
		return nil, handler.HandleError(err)
	}
	merchantAccountData, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id: dbData.MerchantAccountId,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	mAccountNo := merchantAccountData.Account
	if dbData.MerchantAccountId == 0 && dbData.PayAccountType == "2" {
		mAccountNo = config.GetString(bankEnum.PinganIntelligenceAccountNo, "")
	}
	if merchantAccountData == nil && dbData.PayAccountType == "2" && dbData.MerchantAccountId != 0 {
		r, err := s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
			Id: dbData.MerchantAccountId})
		if err != nil {
			return nil, handler.HandleError(err)
		}
		mAccountNo = r.VirtualAccountNo
	}

	return stru.ConvertBankTransactionDetailDataAndMerchantAccount(*dbData, mAccountNo), nil
}

func (s *bankService) SimpleGetBankTransactionDetail(ctx context.Context, req *api.BankTransactionDetailData) (*api.BankTransactionDetailData, error) {
	dbData, err := s.bankTransactionDetailRepo.SimpleGet(ctx, &repo.BankTransactionDetailDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
		ElectronicReceiptFile: req.ElectronicReceiptFile,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	merchantAccountData, err := s.baseClient.SimpleGetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id: dbData.MerchantAccountId,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	bankTransactionDetailData := stru.ConvertBankTransactionDetailDataAndMerchantAccount(*dbData, merchantAccountData.Account)
	bankTransactionDetailData.MerchantAccountOpenName = merchantAccountData.OpenBankName
	return bankTransactionDetailData, nil
}

func (s *bankService) EditBankTransactionDetailExtField(ctx context.Context, req *api.BankTransactionDetailData) error {
	if req.Id == 0 {
		return nil
	}
	return s.bankTransactionDetailRepo.UpdateSelectedFieldsById(ctx, req.Id, &repo.BankTransactionDetailDBData{
		ExtField3: req.ExtField3,
		ExtField2: req.ExtField2,
	}, []string{"ExtField2", "ExtField3"})
}

func (s *bankService) CreateTransactionDetailProcessInstance(ctx context.Context, id int64) error {
	zap.L().Info(fmt.Sprintf("==CreateTransactionDetailProcessInstance%v", id))
	dbData, err := s.bankTransactionDetailRepo.Get(ctx, &repo.BankTransactionDetailDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: id,
			},
		},
	})
	zap.L().Info(fmt.Sprintf("==bankTransactionDetailRepo%v", dbData))
	if err != nil || dbData == nil {
		return handler.HandleError(err)
	}
	if dbData.Type != enum.GuilinBankTransactionDetailRecType {
		return nil
	}
	processMappings, err := s.baseClient.ListOrganizationUnionProcessConfig(ctx, &baseApi.ListOrganizationUnionProcessConfigRequest{
		OrganizationId: dbData.OrganizationId,
		Type:           "0",
		TemplateType:   enum.GuilinBankCollectionNoticeType,
	})
	zap.L().Info(fmt.Sprintf("==ListOrganizationUnionProcessConfig%v", processMappings))
	if err != nil || processMappings == nil {
		return handler.HandleError(err)
	}
	organizationUnionConfig, err := s.baseClient.GetOrganizationUnionConfig(ctx, &baseApi.OrganizationUnionConfigData{
		OrganizationId: dbData.OrganizationId,
		Type:           "0",
	})
	zap.L().Info(fmt.Sprintf("==GetOrganizationUnionConfig%v", organizationUnionConfig))
	if err != nil || organizationUnionConfig == nil {
		return handler.HandleError(err)
	}
	for _, processMapping := range processMappings {
		configMap := make(map[string]string)
		if err = json.Unmarshal(processMapping.Config, &configMap); err != nil {
			return handler.HandleError(err)
		}
		var cashFlag string
		if dbData.CashFlag == "0" {
			cashFlag = enum.GuilinBankTransactionDetailCashFlag
		} else if dbData.CashFlag == "1" {
			cashFlag = enum.GuilinBankTransactionDetailTransferFlag
		}
		var attachmentValueBytes []byte
		unionId := organizationUnionConfig.FileManagerUnionId
		externalUserId := organizationUnionConfig.AdminUnionUserId
		spaceId := organizationUnionConfig.FileSpaceId
		instanceId := dbData.ProcessInstanceId
		if dbData.ElectronicReceiptFile != "" {
			fileUrl, err := store.GetOSSPrivateFile(dbData.ElectronicReceiptFile, s.ossConfig)
			if err != nil {
				return handler.HandleError(err)
			}
			zap.L().Info(fmt.Sprintf("==s.dingtalkClient.UploadDingtalkFile.args=unionId:%s,AppKey:%s,AppSecret%s,spaceId:%s,instanceId:%s\n",
				unionId, organizationUnionConfig.AppKey, organizationUnionConfig.AppSecret, spaceId, string(instanceId)))
			attachment, err := s.dingtalkClient.UploadDingtalkFile(ctx, unionId, fileUrl, "电子回单.pdf", organizationUnionConfig.AppKey, organizationUnionConfig.AppSecret, spaceId, instanceId)
			if err != nil {
				return handler.HandleError(err)
			}
			attachmentValueBytes, err = json.Marshal([]interface{}{attachment})
			if err != nil {
				return handler.HandleError(err)
			}
		}
		transferTime, _ := util.ParseDateTimeMMddHHmmss2(dbData.TransferTime)
		externalId, err := s.dingtalkClient.CreateProcessInstance(ctx, &dingtalkApi.CreateProcessInstanceRequest{
			ProcessCode: processMapping.ProcessCode,
			UserId:      externalUserId,
			FormComponentValues: []*dingtalkApi.ProcessInstanceFormComponentRequest{
				{Name: configMap["transferTime"], Value: util.FormatDateTime(transferTime)},
				{Name: configMap["otherSideAccount"], Value: dbData.AcctNo},
				{Name: configMap["otherSideAccountName"], Value: dbData.AccountName},
				{Name: configMap["cashFlag"], Value: cashFlag},
				{Name: configMap["payAmount"], Value: strconv.FormatFloat(dbData.RecAmount, 'f', 2, 64)},
				{Name: configMap["electronicReceipt"], Value: string(attachmentValueBytes)},
			},
		}, organizationUnionConfig.AppKey, organizationUnionConfig.AppSecret)
		if err != nil {
			return handler.HandleError(err)
		}
		_, err = s.bankTransactionDetailProcessInstanceRepo.Add(ctx, &repo.BankTransactionDetailProcessInstanceDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: dbData.OrganizationId,
			},
			BankTransactionDetailId: dbData.Id,
			ExternalId:              externalId,
		})
		if err != nil {
			return handler.HandleError(err)
		}
	}
	return handler.HandleError(err)
}

func (s *bankService) HandleTransactionDetail(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	//查询交易明细处理存到交易明细表里面去
	s.HandleGuilinBankTransactionDetail(ctx, enum.GuilinBankType, beginDate, endDate, organizationId)
	s.HandleSPDBankTransactionDetail(ctx, enum.SPDBankType, beginDate, endDate, organizationId)
	//1.平安银行的当日和历史交易需要分开查询,2.结束日期为当前日期的,不能查出当日期,3.查询昨日单据需要等到6点之后
	//a.先查询当天的,这个交易结束就可以查询,把开始时间和结束时间设置成一样
	s.HandlePinganBankTransactionDetail(ctx, enum.PinganBankType, endDate, endDate, organizationId)
	s.HandleMinShengBankTransactionDetail(ctx, enum.MinShengBankType, beginDate, endDate, organizationId)
	now := time.Now()
	nowHour := now.Hour()
	if nowHour >= 6 {
		//在查询一次昨天的
		//银行要求 同一个账户重新查询第一页限制1分钟间隔,这里等待1分钟之后重新查询
		time.Sleep(time.Duration(60) * time.Second)
		oneDayDuration, _ := time.ParseDuration("-24h")
		yesterday := now.Add(oneDayDuration)
		yesterdayString := util.FormatTimeyyyyMMdd(yesterday)
		s.HandlePinganBankTransactionDetail(ctx, enum.PinganBankType, beginDate, yesterdayString, organizationId)
	}
	//若endDate==今天,那么要分成两部分来进行
	s.HandlePinganBankVirtualTransactionDetail(ctx, enum.PinganBankType, beginDate, endDate, organizationId)
	if endDate == util.FormatTimeyyyyMMdd(now) {
		s.HandlePinganBankVirtualTransactionDetail(ctx, enum.PinganBankType, endDate, endDate, organizationId)
	}
	//查询工商银行
	s.IcbcBankListTransactionDetail(ctx, beginDate, endDate, organizationId)
	return nil
}

func (s *bankService) HandleGuilinBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error {
	merchantAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if merchantAccounts != nil && len(merchantAccounts) > 0 {
		for _, merchantAccount := range merchantAccounts {
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: merchantAccount.OrganizationId,
				Type:           bankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			//organizationUnionConfig, err := s.baseClient.GetOrganizationUnionConfig(ctx, &baseApi.OrganizationUnionConfigData{
			//	OrganizationId: merchantAccount.OrganizationId,
			//	Type:           "0",
			//})
			if err != nil {
				return err
			}
			datas, err := s.guilinBankSDK.ListTransactionDetail(ctx, merchantAccount.Account, beginDate, endDate, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			if err != nil {
				zap.L().Error(fmt.Sprintf("s.guilinBankSDK.ListTransactionDetail__error_info%s", err.Error()))
				return handler.HandleError(err)
			}
			zap.L().Info(fmt.Sprintf("s.guilinBankSDK.ListTransactionDetail_info:%v", datas))
			if datas != nil {
				var addDatas []repo.BankTransactionDetailDBData
				for _, data := range datas {
					count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
						BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
							BaseDBData: repository.BaseDBData{
								OrganizationId: merchantAccount.OrganizationId,
							},
							MerchantAccountId: merchantAccount.Id,
							OrderFlowNo:       data.OrderFlowNo,
							HostFlowNo:        data.HostFlowNo,
						},
					})
					if err != nil {
						return handler.HandleError(err)
					}
					if count == 0 {
						f, err := s.guilinBankSDK.GetTransactionDetailElectronicReceipt(ctx, data.OrderFlowNo, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
						if err != nil {
							return err
						}
						var electronicReceiptFile string
						if len(f) > 0 {
							electronicReceiptFile, err = store.UploadOSSFileBytes("pdf", ".pdf", f, s.ossConfig, false)
							if err != nil {
								return err
							}
						}

						transactionDetailDBData := repo.BankTransactionDetailDBData{
							BaseDBData: repository.BaseDBData{
								OrganizationId: merchantAccount.OrganizationId,
							},
							MerchantAccountId:     merchantAccount.Id,
							MerchantAccountName:   merchantAccount.AccountName,
							CashFlag:              data.CashFlag,
							PayAmount:             data.PayAmount,
							RecAmount:             data.RecAmount,
							BsnType:               data.BsnType,
							TransferDate:          data.TransferDate,
							TransferTime:          data.TransferTime,
							TranChannel:           data.TranChannel,
							CurrencyType:          data.CurrencyType,
							Balance:               data.Balance,
							OrderFlowNo:           data.OrderFlowNo,
							HostFlowNo:            data.HostFlowNo,
							VouchersType:          data.VouchersType,
							VouchersNo:            data.VouchersNo,
							SummaryNo:             data.SummaryNo,
							Summary:               data.Summary,
							AcctNo:                data.AcctNo,
							AccountName:           data.AccountName,
							AccountOpenNode:       data.AccountOpenNode,
							ElectronicReceiptFile: electronicReceiptFile,
							ProcessTotalStatus:    enum.ProcessInstanceTotalStatusRunning,
							PayAccountType:        enum.GuilinBankType,
						}
						if transactionDetailDBData.PayAmount < 0 {
							transactionDetailDBData.RecAmount = -transactionDetailDBData.PayAmount
							transactionDetailDBData.PayAmount = 0
						}
						if data.RecAmount > 0 {
							//查询发起人
							//userData, err := s.dingtalkClient.GetUserDataByUnionId(ctx, organizationUnionConfig.NotificationUnionId)
							//if err != nil || userData == nil {
							//	return err
							//}
							//
							//transactionDetailDBData.OriginatorUserId = userData.UserId
							//transactionDetailDBData.OriginatorUserName = userData.Name
							transactionDetailDBData.Type = enum.GuilinBankTransactionDetailRecType
							// 扩展字段3：用来标识-收款确认单同步状态（0-待同步 1-同步中 2-同步成功 3-同步失败）
							transactionDetailDBData.ExtField3 = "0"
						} else {
							transactionDetailDBData.Type = enum.GuilinBankTransactionDetailPayType
						}
						addDatas = append(addDatas, transactionDetailDBData)

						if err != nil {
							return handler.HandleError(err)
						}

						//更新单据的明细ID
						if err = s.updateRelevanceElectronicDocument(ctx, data.Summary, data.HostFlowNo, electronicReceiptFile, enum.GuilinBankType); err != nil {
							return handler.HandleError(err)
						}
					}
				}
				if len(addDatas) > 0 {
					ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDatas)
					if err != nil {
						return handler.HandleError(err)
					}
					for _, id := range ids {
						s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
							Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
							Type:     kafka.DingtalkType,
							Id:       id,
						})
					}
				}
			}
			//now := time.Now()
			//if err = s.guilinBankMerchantAccountRepo.UpdateById(ctx, merchantAccount.Id, &repo.GuilinBankMerchantAccountDBData{
			//	SyncTime: &now,
			//}, 0); err != nil {
			//	return err
			//}
		}
	}
	return nil
}

func (s *bankService) HandleMinShengBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error {
	merchantAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if merchantAccounts != nil && len(merchantAccounts) > 0 {
		for _, merchantAccount := range merchantAccounts {
			result, err := s.minShengBank.ListTransactionDetail(ctx, merchantAccount.Account, beginDate, endDate, "1", "200", merchantAccount.OpenId)
			if err != nil {
				zap.L().Error(fmt.Sprintf("s.minShengBankSDK.ListTransactionDetail__error_info%s", err.Error()))
				return handler.HandleError(err)
			}
			zap.L().Info(fmt.Sprintf("s.minShengBankSDK.ListTransactionDetail_info:%v", result))
			if result["return_code"] != "0000" {
				zap.L().Info(fmt.Sprintf("s.minShengBankSDK.ListTransactionDetail_info查询转账结果:%v", result))
				continue
			}
			responseBusi := result["response_busi"].(string)
			var busiMap map[string]string
			err = json.Unmarshal([]byte(responseBusi), &busiMap)
			if err != nil {
				zap.L().Info(fmt.Sprintf("s.minShengBankSDK.ListTransactionDetail_info转换response_busi异常:%v", err))
				continue
			}
			var minShengTransactionDetails []stru.MinShengTransactionDetailResponse
			err = json.Unmarshal([]byte(busiMap["result_list"]), &minShengTransactionDetails)
			if err != nil {
				zap.L().Info(fmt.Sprintf("s.minShengBankSDK.ListTransactionDetail_info转换result_list异常:%v", err))
				continue
			}
			if minShengTransactionDetails != nil && len(minShengTransactionDetails) > 0 {
				var addDatas []repo.BankTransactionDetailDBData
				for _, data := range minShengTransactionDetails {
					count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
						BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
							BaseDBData: repository.BaseDBData{
								OrganizationId: merchantAccount.OrganizationId,
							},
							MerchantAccountId: merchantAccount.Id,
							OrderFlowNo:       data.TransSeqNo,
							HostFlowNo:        data.TransSeqNo,
						},
					})
					if err != nil {
						return handler.HandleError(err)
					}
					if count == 0 {
						receiptResponse, err := s.minShengBank.GetTransactionDetailElectronicReceipt(ctx, data.AcctNo, data.TransSeqNo, data.EnterAcctDate, merchantAccount.OpenId)
						if err != nil {
							return err
						}
						if receiptResponse["return_code"] != "0000" {
							zap.L().Info(fmt.Sprintf("s.minShengBankSDK.GetTransactionDetailElectronicReceipt查询回单结果:%v", receiptResponse))
							continue
						}
						electronicReceiptFile := receiptResponse["response_busi"].(string)

						amount, _ := strconv.ParseFloat(data.Amount, 64)
						balance, _ := strconv.ParseFloat(data.Balance, 64)
						tranChannel := ""
						if data.DcFlag == "1" {
							tranChannel = "PAY"
						} else if data.DcFlag == "2" {
							tranChannel = "REC"
						}
						transactionDetailDBData := repo.BankTransactionDetailDBData{
							BaseDBData: repository.BaseDBData{
								OrganizationId: merchantAccount.OrganizationId,
							},
							MerchantAccountId:     merchantAccount.Id,
							MerchantAccountName:   merchantAccount.AccountName,
							PayAmount:             amount,
							RecAmount:             amount,
							TransferDate:          data.EnterAcctDate,
							TransferTime:          data.Timestamp,
							TranChannel:           tranChannel,
							CurrencyType:          data.Currency,
							Balance:               balance,
							OrderFlowNo:           data.TransSeqNo,
							HostFlowNo:            data.TransSeqNo,
							Summary:               data.Explain,
							AcctNo:                data.CpAcctNo,
							AccountName:           data.CpAcctName,
							AccountOpenNode:       data.CpBankName,
							ElectronicReceiptFile: electronicReceiptFile,
							ProcessTotalStatus:    enum.ProcessInstanceTotalStatusRunning,
							PayAccountType:        enum.MinShengBankType,
						}
						if transactionDetailDBData.PayAmount < 0 {
							transactionDetailDBData.RecAmount = -transactionDetailDBData.PayAmount
							transactionDetailDBData.PayAmount = 0
						}
						//if data.RecAmount > 0 {
						//	transactionDetailDBData.Type = enum.GuilinBankTransactionDetailRecType
						//	// 扩展字段3：用来标识-收款确认单同步状态（0-待同步 1-同步中 2-同步成功 3-同步失败）
						//	transactionDetailDBData.ExtField3 = "0"
						//} else {
						//	transactionDetailDBData.Type = enum.GuilinBankTransactionDetailPayType
						//}
						addDatas = append(addDatas, transactionDetailDBData)

						if err != nil {
							return handler.HandleError(err)
						}

						//更新单据的明细ID
						if err = s.updateRelevanceElectronicDocument(ctx, data.Explain, data.TransSeqNo, electronicReceiptFile, enum.MinShengBankType); err != nil {
							return handler.HandleError(err)
						}
					}
				}
				if len(addDatas) > 0 {
					ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDatas)
					if err != nil {
						return handler.HandleError(err)
					}
					for _, id := range ids {
						s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
							Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
							Type:     kafka.DingtalkType,
							Id:       id,
						})
					}
				}
			}
			//now := time.Now()
			//if err = s.guilinBankMerchantAccountRepo.UpdateById(ctx, merchantAccount.Id, &repo.GuilinBankMerchantAccountDBData{
			//	SyncTime: &now,
			//}, 0); err != nil {
			//	return err
			//}
		}
	}
	return nil
}

// bankType 银行类型:0,桂林,1:浦发,2:平安,4:民生
func (s *bankService) updateRelevanceElectronicDocument(ctx context.Context, summary, hostFlowNo, electronicDocument, bankType string) error {
	var paymentReceipt *repo.PaymentReceiptDBData
	electronicDocumentPng := electronicDocument
	var err error
	switch bankType {
	case enum.GuilinBankType:
		if summary != "" && strings.Index(summary, "[") >= 0 && strings.Index(summary, "]") >= 0 {
			serialNo := summary[strings.Index(summary, "[")+1 : strings.Index(summary, "]")]
			transferReceiptData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
				SerialNo: serialNo,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			if transferReceiptData != nil && transferReceiptData.Id != 0 {
				if err = s.bankTransferReceiptRepo.UpdateById(ctx, transferReceiptData.Id, &repo.BankTransferReceiptDBData{
					DetailHostFlowNo: hostFlowNo,
				}); err != nil {
					return handler.HandleError(err)
				}
			}
			paymentReceipt, err = s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
				Code: serialNo,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			// 上传png
			//electronicDocumentPng, err = s.pdfToImageService.UploadPdfToImageJsdk(ctx, electronicDocument)
			/*if err != nil {
				return handler.HandleError(err)
			}*/
		}
	case enum.SPDBankType:
		//浦发只用更新回单即可,
		if summary != "" && strings.Index(summary, "[") >= 0 && strings.Index(summary, "]") >= 0 {
			serialNo := summary[strings.Index(summary, "[")+1 : strings.Index(summary, "]")]
			paymentReceipt, err = s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
				Code: serialNo,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			// 上传png
			/*electronicDocumentPng, err = s.pdfToImageService.UploadPdfToImageJsdk(ctx, electronicDocument)
			if err != nil {
				return handler.HandleError(err)
			}*/
		}
	case enum.PinganBankType:
		paymentReceipt, err = s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
			OrderFlowNo: hostFlowNo,
		})
	case enum.MinShengBankType:
		paymentReceipt, err = s.paymentReceiptRepo.GetWithoutPermission(ctx, &repo.PaymentReceiptDBData{
			OrderFlowNo: hostFlowNo,
		})
	default:

	}

	if paymentReceipt != nil && paymentReceipt.Id != 0 {
		// 付款申请单
		isProcessSuccess := false
		if paymentReceipt.Type == "1" || paymentReceipt.Type == "" {
			paymentApplication, err := s.oaClient.GetPaymentApplicationByProcessInstanceId(ctx, paymentReceipt.ProcessInstanceId)
			if err != nil {
				return handler.HandleError(err)
			}
			if paymentApplication != nil && paymentApplication.Id != 0 && electronicDocument != "" {
				if paymentApplication.OrderStatus == enum.GuilinBankTransferHandling {
					paymentApplication.OrderStatus = enum.GuilinBankTransferSuccessResult
					isProcessSuccess = true
				}
				err := s.oaClient.EditPaymentApplicationWithoutPermission(ctx, &oaApi.PaymentApplicationData{
					Id:                    paymentApplication.Id,
					ElectronicDocument:    electronicDocument,
					ElectronicDocumentPng: electronicDocumentPng,
					OrderStatus:           paymentApplication.OrderStatus,
				})
				if err != nil {
					return handler.HandleError(err)
				}
			}
		} else if paymentReceipt.Type == "2" {
			// 报销申请单
			reimburseApplication, err := s.oaClient.GetReimburseApplicationByProcessInstanceId(ctx, paymentReceipt.ProcessInstanceId)
			if err != nil {
				return handler.HandleError(err)
			}
			if reimburseApplication != nil && reimburseApplication.Id != 0 && electronicDocument != "" {
				if reimburseApplication.OrderStatus == enum.GuilinBankTransferHandling {
					reimburseApplication.OrderStatus = enum.GuilinBankTransferSuccessResult
					isProcessSuccess = true
				}
				err := handler.HandleError(s.oaClient.EditReimburseApplicationWithoutPermission(ctx, &oaApi.ReimburseApplicationData{
					Id:                    reimburseApplication.Id,
					ElectronicDocument:    electronicDocument,
					ElectronicDocumentPng: electronicDocumentPng,
					OrderStatus:           reimburseApplication.OrderStatus,
				}))
				if err != nil {
					return handler.HandleError(err)
				}
			}
		}

		// 更新付款单据表回单和状态
		if electronicDocument != "" {
			if paymentReceipt.OrderStatus == enum.GuilinBankTransferHandling {
				paymentReceipt.OrderStatus = enum.GuilinBankTransferSuccessResult
			}
			err := s.paymentReceiptRepo.UpdateByIdWithoutPermission(ctx, paymentReceipt.Id, &repo.PaymentReceiptDBData{
				ElectronicDocument:    electronicDocument,
				ElectronicDocumentPng: electronicDocumentPng,
				OrderStatus:           paymentReceipt.OrderStatus,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			if isProcessSuccess {
				_, err := s.baseClient.SuccessProcessInstance(ctx, paymentReceipt.ProcessInstanceId)
				if err != nil {
					return handler.HandleError(err)
				}
			}
		}
	}
	return nil
}

func (s *bankService) HandleSPDBankTransactionDetail(ctx context.Context, bankType, beginDate, endDate string, organizationId int64) error {
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, bankAccount := range bankAccounts {
		organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
			OrganizationId: bankAccount.OrganizationId,
			Type:           bankType,
		})
		if err != nil {
			return handler.HandleError(err)
		}
		if err != nil {
			return err
		}
		datas, err := s.spdBankSDK.ListTransactionDetail(ctx, bankAccount.Account, beginDate, endDate, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.spdBankSDK.ListTransactionDetail__error_info%s", err.Error()))
			return handler.HandleError(err)
		}
		// 将对象转换为JSON
		jsonData, _ := json.Marshal(datas)
		zap.L().Info("s.spdBankSDK.ListTransactionDetail_info", zap.String("details", string(jsonData)))
		if datas != nil {
			var addDatas []repo.BankTransactionDetailDBData
			for _, data := range datas {
				count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
					BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						MerchantAccountId: bankAccount.Id,
						OrderFlowNo:       data.TellerJnlNo, // 浦发柜员流水号
						HostFlowNo:        data.TellerJnlNo,
					},
				})
				if err != nil {
					return handler.HandleError(err)
				}

				// 保存交易明细
				if count == 0 {
					transtime := fmt.Sprintf("%06s", data.TransTime)
					// 日期8位 时间6位 脏数据校验
					newTransTime := data.TransDate + transtime[0:2] + ":" + transtime[2:4] + ":" + transtime[4:6]

					// 浦发借贷标记 0-借/收 1-贷/付
					payAmount := 0.00
					recAmount := 0.00
					if data.DebitFlag == "1" {
						recAmount = data.TransAmount
					} else if data.DebitFlag == "0" {
						payAmount = data.TransAmount
					}

					transactionDetailDBData := repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						MerchantAccountId:   bankAccount.Id,
						MerchantAccountName: bankAccount.AccountName,
						CashFlag:            "0", // 写死0: 现钞
						PayAmount:           payAmount,
						RecAmount:           recAmount,
						BsnType:             "TR", // 写死TR: 转账
						TransferDate:        data.TransDate,
						TransferTime:        newTransTime,
						TranChannel:         "",
						CurrencyType:        "CNY",
						Balance:             data.AcctBalance,
						OrderFlowNo:         data.TellerJnlNo,
						HostFlowNo:          data.TellerJnlNo,
						VouchersType:        "",
						VouchersNo:          "",
						SummaryNo:           data.SummaryCode,
						Summary:             data.Remark,
						AcctNo:              data.SubAccount,
						AccountName:         data.SubAcctName,
						AccountOpenNode:     data.OppositeBankName,
						//ElectronicReceiptFile: electronicReceiptFile,
						ProcessTotalStatus: enum.ProcessInstanceTotalStatusRunning,
						PayAccountType:     enum.SPDBankType,
						ExtField1:          data.SummonsNumber,
					}
					if recAmount > 0 {
						transactionDetailDBData.Type = enum.GuilinBankTransactionDetailRecType
						// 扩展字段3：用来标识-收款确认单同步状态（0-待同步 1-同步中 2-同步成功 3-同步失败）
						transactionDetailDBData.ExtField3 = "0"
					} else {
						transactionDetailDBData.Type = enum.GuilinBankTransactionDetailPayType
					}

					addDatas = append(addDatas, transactionDetailDBData)

					//更新单据的明细ID
					if data.Remark != "" && strings.Index(data.Remark, "[") >= 0 && strings.Index(data.Remark, "]") >= 0 {
						serialNo := data.Remark[strings.Index(data.Remark, "[")+1 : strings.Index(data.Remark, "]")]
						transferReceiptData, err := s.bankTransferReceiptRepo.Get(ctx, &repo.BankTransferReceiptDBData{
							SerialNo: serialNo,
						})
						if err == nil && transferReceiptData != nil {
							s.bankTransferReceiptRepo.UpdateById(ctx, transferReceiptData.Id, &repo.BankTransferReceiptDBData{
								DetailHostFlowNo: data.TellerJnlNo,
							})
						}

					}
				}
			}
			if len(addDatas) > 0 {
				ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDatas)
				if err != nil {
					return handler.HandleError(err)
				}
				for _, id := range ids {
					s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
						Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
						Type:     kafka.DingtalkType,
						Id:       id,
					})
				}
			}
		}
		//now := time.Now()
		//if err = s.guilinBankMerchantAccountRepo.UpdateById(ctx, merchantAccount.Id, &repo.GuilinBankMerchantAccountDBData{
		//	SyncTime: &now,
		//}, 0); err != nil {
		//	return err
		//}
	}
	return nil
}
func (s *bankService) HandlePinganBankTransactionDetail(ctx context.Context, bankType string, beginDate string, endDate string, organizationId int64) error {
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		OrganizationId:       organizationId,
		Type:                 bankType,
		SignatureApplyStatus: "0",
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, bankAccount := range bankAccounts {

		datas, err := s.pinganBankSDK.ListTransactionDetail(ctx, bankAccount.Account, beginDate, endDate, bankAccount.ZuId)
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.pinganBankSDK.ListTransactionDetail__error_info%s", err.Error()))
			continue
		}
		zap.L().Info(fmt.Sprintf("s.pinganBankSDK.ListTransactionDetail_info:%+v", datas))
		if datas != nil {
			var addDatas []repo.BankTransactionDetailDBData
			feeMap := make(map[string]string)
			var rechargeList []*flexApi.FlexVirtualAccountRechargeFlowItem
			for _, data := range datas {
				count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
					BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						MerchantAccountId: bankAccount.Id,
						//OrderFlowNo:       data.BussSeqNo, // 业务流水号
						HostFlowNo: data.HostTrace, //主机流水号
					},
				})
				if err != nil {
					continue
				}
				// 保存交易明细
				if count == 0 {
					//如果本条数据AbstractStr = "FEE"就说明这笔明细是手续费,将其放入map中,
					//最后放入对应的HostTrace的数据中
					if data.AbstractStr == "FEE" {
						feeMap[data.HostTrace] = data.TranAmount
						continue
					}
					//  D借，出账；C贷，入账
					payAmount := 0.00
					recAmount := 0.00
					TransactionType := ""
					tranAmount, err := strconv.ParseFloat(data.TranAmount, 64)
					if err != nil {
						continue
					}

					//对方银行
					oppAccountNo := ""
					oppAccountName := ""
					oppAccountBankName := ""

					if data.DcFlag == "C" {
						recAmount = tranAmount
						TransactionType = enum.GuilinBankTransactionDetailRecType
						oppAccountNo = data.OutAcctNo
						oppAccountName = data.OutAcctName
						oppAccountBankName = data.OutBankNo
					} else if data.DcFlag == "D" {
						payAmount = tranAmount
						TransactionType = enum.GuilinBankTransactionDetailPayType
						oppAccountNo = data.InAcctNo
						oppAccountName = data.InAcctName
						oppAccountBankName = data.InBankNo
					}
					acctBalance, err := strconv.ParseFloat(data.AcctBalance, 64)
					if err != nil {
						continue
					}
					transactionDetailDBData := repo.BankTransactionDetailDBData{
						BaseDBData: repository.BaseDBData{
							OrganizationId: bankAccount.OrganizationId,
						},
						Type:                TransactionType,
						MerchantAccountId:   bankAccount.Id,
						MerchantAccountName: bankAccount.AccountName,
						CashFlag:            "0", // 写死0: 现钞
						PayAmount:           payAmount,
						RecAmount:           recAmount,
						BsnType:             "TR", // 写死TR: 转账
						TransferDate:        data.AcctDate,
						TransferTime:        data.TxTime,
						TranChannel:         "",
						CurrencyType:        "CNY",
						Balance:             acctBalance,
						//OrderFlowNo:         data.BussSeqNo,
						OrderFlowNo:        data.HostTrace,
						HostFlowNo:         data.HostTrace,
						VouchersType:       "",
						VouchersNo:         "",
						SummaryNo:          data.AbstractStr,
						Summary:            data.AbstractStrDesc,
						AcctNo:             oppAccountNo,
						AccountName:        oppAccountName,
						AccountOpenNode:    oppAccountBankName,
						ProcessTotalStatus: enum.ProcessInstanceTotalStatusRunning,
						PayAccountType:     enum.PinganBankType,
						ExtField1:          data.TranFee,
					}

					if strings.Contains(data.AbstractStrDesc, "服务费") && bankAccount.AccountName == config.GetString(bankEnum.PinganIntelligenceAccountName, "") && data.DcFlag == "C" {
						rechargeList = append(rechargeList, &flexApi.FlexVirtualAccountRechargeFlowItem{
							RecAccountNo:      bankAccount.Account,
							RecAccountName:    bankAccount.AccountName,
							PayAccountNo:      oppAccountNo,
							PayAccountName:    oppAccountName,
							RechargeAmount:    tranAmount,
							Summary:           data.AbstractStrDesc,
							OrderStatus:       enum.ProcessInstanceTotalStatusRunning,
							ResMessage:        "",
							ResCode:           "",
							TransferReceiptId: 0,
							OrderFlowNo:       data.HostTrace,
						})
					}
					if data.DcFlag == "C" {
						// 扩展字段3：用来标识-收款确认单同步状态（0-待同步 1-同步中 2-同步成功 3-同步失败）
						transactionDetailDBData.ExtField3 = "0"
					}

					//封装请求body
					serialNo, _ := util.SonyflakeID()

					request := sdkStru.PinganSameDayHistoryReceiptDataQueryRequest{
						MrchCode:         config.GetString(bankEnum.PinganMrchCode, ""),
						CnsmrSeqNo:       serialNo,
						OutAccNo:         data.OutAcctNo,
						AccountBeginDate: data.HostDate,
						AccountEndDate:   data.HostDate,
						HostFlow:         data.HostTrace,
					}
					//查询回单并且转成png到oss
					f, err := s.pinganBankSDK.UploadTransactionDetailElectronic(ctx, request, bankAccount.ZuId)
					if err != nil {
						zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证失败: %v\n", err.Error()))
					}
					if f != "" {
						transactionDetailDBData.ElectronicReceiptFile = f
					}
					addDatas = append(addDatas, transactionDetailDBData)

					//更新单据的回单
					if err = s.updateRelevanceElectronicDocument(ctx, "", data.HostTrace, f, enum.PinganBankType); err != nil {
						return handler.HandleError(err)
					}
				}
			}
			if len(rechargeList) > 0 {
				s.flexClient.SyncVirtualAccountRechargeFlow(ctx, &flexApi.FlexVirtualAccountRechargeFlowData{
					Data: rechargeList,
				})
			}
			if len(addDatas) > 0 {
				var addDetails []repo.BankTransactionDetailDBData
				for _, data := range addDatas {
					if fee, ok := feeMap[data.HostFlowNo]; ok {
						data.ExtField1 = fee
					}
					addDetails = append(addDetails, data)
				}
				ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDetails)
				if err != nil {
					return handler.HandleError(err)
				}
				for _, id := range ids {
					s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
						Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
						Type:     kafka.DingtalkType,
						Id:       id,
					})
				}
			}
		}
	}
	return nil
}
func (s *bankService) HandlePinganBankVirtualTransactionDetail(ctx context.Context, bankType string, beginDate string, endDate string, organizationId int64) error {
	virtualBankAccounts, err := s.baseClient.ListOrganizationBankVirtualAccountData(ctx, &baseApi.ListOrganizationBankVirtualAccountRequest{
		//OrganizationId: organizationId,
		Type: bankType,
	})
	//使用主账号去查所有的流水,然后根据摘要中写的去判断数属于哪个子账号
	bankAccountNo := config.GetString(bankEnum.PinganIntelligenceAccountNo, "")
	datas, err := s.pinganBankSDK.ListVirtualTransactionDetail(ctx, bankAccountNo, beginDate, endDate)
	if err != nil {
		zap.L().Error(fmt.Sprintf("s.pinganBankSDK.HandlePinganBankVirtualTransactionDetail%s", err.Error()))
	}
	//zap.L().Info(fmt.Sprintf("s.pinganBankSDK.HandlePinganBankVirtualTransactionDetail:%+v", datas))
	if datas != nil {
		var addDatas []repo.BankTransactionDetailDBData
		feeMap := make(map[string]string)
		for _, data := range datas {
			//根据 "Purpose": "代30210294284702[鑫旷世碧园] 付款 （服务费CZ1692263918387",
			//寻找子账号,然后根据填入正确的组织id和accountId
			currentOrganizationId := int64(1)
			bankAccountId := int64(0)
			bankAccountName := config.GetString(bankEnum.PinganIntelligenceAccountName, "")
			subBankAccountNo := bankAccountNo

			re := regexp.MustCompile(`\d+`)
			match := re.FindStringSubmatch(data.Purpose)
			if len(match) > 0 {
				subBankAccountNo = match[0]
				for _, subAccount := range virtualBankAccounts {
					if subAccount.VirtualAccountNo == subBankAccountNo {
						//currentOrganizationId = subAccount.OrganizationId
						bankAccountId = subAccount.Id
						bankAccountName = subAccount.VirtualAccountName
						break
					}
				}
			}

			count, err := s.bankTransactionDetailRepo.Count(ctx, &repo.BankTransactionDetailDBDataParam{
				BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
					BaseDBData: repository.BaseDBData{
						OrganizationId: currentOrganizationId,
					},
					MerchantAccountId: bankAccountId,
					//OrderFlowNo:       data.BussSeqNo, // 业务流水号
					HostFlowNo: data.HostTrace, //主机流水号
				},
			})
			if err != nil {
				continue
			}
			// 保存交易明细
			if count == 0 {
				//如果本条数据AbstractStr = "FEE"就说明这笔明细是手续费,将其放入map中,
				//最后放入对应的HostTrace的数据中
				if data.AbstractStr == "FEE" {
					feeMap[data.HostTrace] = data.TranAmount
					continue
				}
				//  D借，出账；C贷，入账
				payAmount := 0.00
				recAmount := 0.00
				TransactionType := ""
				tranAmount, err := strconv.ParseFloat(data.TranAmount, 64)
				//对方银行
				oppAccountNo := ""
				oppAccountName := ""
				oppAccountBankName := ""
				if err != nil {
					continue
				}

				if data.DcFlag == "C" { //收钱
					recAmount = tranAmount
					TransactionType = enum.GuilinBankTransactionDetailRecType
					oppAccountNo = data.OutAcctNo
					oppAccountName = data.OutAcctName
					oppAccountBankName = data.OutBankNo
				} else if data.DcFlag == "D" { //出钱
					payAmount = tranAmount
					TransactionType = enum.GuilinBankTransactionDetailPayType
					oppAccountNo = data.InAcctNo
					oppAccountName = data.InAcctName
					oppAccountBankName = data.InBankName
				}
				acctBalance, err := strconv.ParseFloat(data.AcctBalance, 64)
				if err != nil {
					continue
				}
				transactionDetailDBData := repo.BankTransactionDetailDBData{
					BaseDBData: repository.BaseDBData{
						OrganizationId: currentOrganizationId,
					},
					Type:                TransactionType,
					MerchantAccountId:   bankAccountId,
					MerchantAccountName: bankAccountName,
					CashFlag:            "0", // 写死0: 现钞
					PayAmount:           payAmount,
					RecAmount:           recAmount,
					BsnType:             "TR", // 写死TR: 转账
					TransferDate:        data.AcctDate,
					TransferTime:        data.AcctDate + data.TxTime,
					TranChannel:         "",
					CurrencyType:        "CNY",
					Balance:             acctBalance,
					//OrderFlowNo:         data.BussSeqNo,
					OrderFlowNo:        data.HostTrace,
					HostFlowNo:         data.HostTrace,
					VouchersType:       "",
					VouchersNo:         "",
					SummaryNo:          data.AbstractStr + data.AbstractStrDesc,
					Summary:            data.Purpose,
					AcctNo:             oppAccountNo,
					AccountName:        oppAccountName,
					AccountOpenNode:    oppAccountBankName,
					ProcessTotalStatus: enum.ProcessInstanceTotalStatusRunning,
					PayAccountType:     enum.PinganBankType,
					ExtField1:          data.TranFee,
				}

				if data.DcFlag == "C" {
					// 扩展字段3：用来标识-收款确认单同步状态（0-待同步 1-同步中 2-同步成功 3-同步失败）
					transactionDetailDBData.ExtField3 = "0"
				}

				//封装请求body
				serialNo, _ := util.SonyflakeID()

				request := sdkStru.PinganSameDayHistoryReceiptDataQueryRequest{
					MrchCode:         config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
					CnsmrSeqNo:       serialNo,
					OutAccNo:         bankAccountNo,
					AccountBeginDate: data.HostDate,
					AccountEndDate:   data.HostDate,
					HostFlow:         data.HostTrace,
				}
				//查询回单并且转成png到oss
				f, err := s.pinganBankSDK.UploadVirtualTransactionDetailElectronic(ctx, request)
				if err != nil {
					zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证失败: %v\n", err.Error()))
				}
				if f != "" {
					transactionDetailDBData.ElectronicReceiptFile = f
					if TransactionType == enum.GuilinBankTransactionDetailPayType {
						//同步付款表里面的回单
						s.bankTransferReceiptRepo.UpdateElectronicReceiptFile(ctx, data.HostTrace, enum.PinganBankType, f)
					}
				}
				addDatas = append(addDatas, transactionDetailDBData)
			}
		}
		if len(addDatas) > 0 {
			var addDetails []repo.BankTransactionDetailDBData
			for _, data := range addDatas {
				if fee, ok := feeMap[data.HostFlowNo]; ok {
					data.ExtField1 = fee
				}
				addDetails = append(addDetails, data)
			}
			ids, err := s.bankTransactionDetailRepo.BatchAdd(ctx, &addDetails)
			if err != nil {
				return handler.HandleError(err)
			}
			for _, id := range ids {
				s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
					Business: kafka.ProcessFinanceTransactionDetailProcessInstanceBusiness,
					Type:     kafka.DingtalkType,
					Id:       id,
				})
			}
		}
	}
	return nil
}

func (s *bankService) ListBankTransactionDetailProcessInstance(ctx context.Context, id int64) ([]*api.BankTransactionDetailProcessInstanceData, error) {
	dbData, _, err := s.bankTransactionDetailProcessInstanceRepo.List(ctx, "created_at", 0, 0, &repo.BankTransactionDetailProcessInstanceDBData{
		BankTransactionDetailId: id,
	})
	if err != nil || dbData == nil {
		return nil, err
	}
	data := make([]*api.BankTransactionDetailProcessInstanceData, len(*dbData))
	for i, v := range *dbData {
		data[i] = &api.BankTransactionDetailProcessInstanceData{
			BankTransactionDetailId: v.BankTransactionDetailId,
			ExternalId:              v.ExternalId,
		}
	}
	return data, nil
}

func (s *bankService) GetBankCodeInfo(ctx context.Context, code string) (*api.BankCodeData, error) {
	bankCodeData, err := s.bankCodeRepo.Get(ctx, &repo.BankCodeDBData{
		BankCode: code,
	})
	if err != nil || bankCodeData == nil {
		return nil, err
	}
	return &api.BankCodeData{
		BankName:      bankCodeData.BankName,
		BankAliasName: bankCodeData.BankAliasName,
		BankCode:      bankCodeData.BankCode,
		UnionBankNo:   bankCodeData.UnionBankNo,
		ClearBankNo:   bankCodeData.ClearBankNo,
	}, err
}

func (s *bankService) QueryBankCardInfo(ctx context.Context, cardNo string) (*api.QueryBankCardInfoResponse, error) {

	//TODO 测试账号
	if "660500051188900010" == cardNo || "000674790600182" == cardNo || "660000015930800015" == cardNo {
		return &api.QueryBankCardInfoResponse{
			Bank:      "GLBANK",
			Stat:      "ok",
			Validated: true,
		}, nil
	}
	if "63010078801300002727" == cardNo {
		return &api.QueryBankCardInfoResponse{
			Bank:      "SPDB",
			Stat:      "ok",
			Validated: true,
		}, nil
	}

	if "" == cardNo {
		return nil, errors.New("查询银行卡开户行失败: 参数为空")
	}

	reqUrl := "https://ccdcapi.alipay.com/validateAndCacheCardInfo.json?_input_charset=utf-8&cardNo=" + cardNo + "&cardBinCheck=true"
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(robots, &result); err != nil {
		return nil, err
	}
	if "ok" != result["stat"].(string) || !result["validated"].(bool) || "" == result["bank"].(string) {
		return nil, errors.New("查询银行卡开户行失败: " + cardNo)
	}

	return &api.QueryBankCardInfoResponse{
		Bank:      result["bank"].(string),
		Stat:      result["stat"].(string),
		Validated: true,
	}, nil
}

func (s *bankService) ListBankCode(ctx context.Context, req *api.ListBankCodeRequest) (*api.ListBankCodeResponse, error) {
	dbData, count, err := s.bankCodeRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankCodeDBData{
		BankName:      req.BankName,
		BankAliasName: req.BankAliasName,
	})
	if err != nil || dbData == nil {
		return nil, handler.HandleError(err)
	}
	data := make([]*api.BankCodeData, len(*dbData))
	for i, v := range *dbData {
		data[i] = stru.ConvertBankCodeData(v)
	}
	return &api.ListBankCodeResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) GetBankCode(ctx context.Context, req *api.BankCodeData) (*api.BankCodeData, error) {
	data, err := s.bankCodeRepo.Get(ctx, stru.ConvertBankCodeDBData(*req))
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return stru.ConvertBankCodeData(*data), nil
}

func (s *bankService) AddBankCode(ctx context.Context, req *api.AddBankCodeRequest) error {
	_, err := s.bankCodeRepo.Add(ctx, stru.ConvertAddBankCodeDBData(*req))
	return err
}

func (s *bankService) EditBankCode(ctx context.Context, req *api.BankCodeData) error {
	return s.bankCodeRepo.UpdateById(ctx, req.Id, stru.ConvertBankCodeDBData(*req))
}

func (s *bankService) DeleteBankCode(ctx context.Context, id int64) error {
	return s.bankCodeRepo.DeleteById(ctx, id)
}

func (s *bankService) HandleSyncTransferReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	s.HandleGuilinBankSyncTransferReceipt(ctx, enum.GuilinBankType, beginDate, endDate, organizationId)
	s.HandleSPDBankSyncTransferReceipt(ctx, enum.SPDBankType, beginDate, endDate, organizationId)
	s.HandlePinganBankSyncTransferReceipt(ctx, enum.PinganBankType, beginDate, endDate, organizationId)
	return nil
}

func (s *bankService) HandleGuilinBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("HandleGuilinBankSyncTransferReceipt 查询桂林所有未完成的付款单据")
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRejectResult}
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
	transferReceiptList, _, err := s.bankTransferReceiptRepo.List(ctx, "", 0, 0, &repo.BankTransferReceiptDBDataParam{
		ExcludeOrderStates:    excludeOrderStates,
		IsOrderFlowNoNotEmpty: true,
		BankTransferReceiptDBData: repo.BankTransferReceiptDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: organizationId,
			},
			ProcessStatus: enum.ProcessInstanceTotalStatusRunning,
		},
		CreateTime:     createTimeParam,
		PayAccountType: bankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if transferReceiptList != nil && len(*transferReceiptList) > 0 {
		for _, transferReceipt := range *transferReceiptList {
			//查询转账结果
			isIntrabankTransfer := transferReceipt.RecBankType == "0"

			beginDate2 := stru.FormatDayTime(transferReceipt.CreatedAt)
			endDate2 := stru.FormatDayTime(time.Now())
			if beginDate2 == endDate {
				endDate2 = stru.FormatDayTimeStamp(time.Now().Unix() + 86400)
			}
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: transferReceipt.OrganizationId,
				Type:           bankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			transferResults, err := s.guilinBankSDK.QueryTransferResult(ctx, transferReceipt.PayAccount, transferReceipt.RecAccount, transferReceipt.OrderFlowNo, isIntrabankTransfer, beginDate2, endDate2, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			zap.L().Info(fmt.Sprintf("guilinBankSDK.QueryTransferResult查询转账结果:%v", transferResults))
			if err != nil {
				return handler.HandleError(err)
			}
			if transferResults != nil && len(transferResults) > 0 {
				for _, result := range transferResults {
					//比较订单状态
					if result.OrderState == transferReceipt.OrderState && result.ErrorCode == "" {
						continue
					}
					//更新单据
					processStatus := enum.ProcessInstanceTotalStatusFinish
					if result.OrderState == enum.GuilinBankTransferSuccessResult {
						processStatus = enum.ProcessInstanceTotalStatusSuccess
					}
					updateReceipt := &repo.BankTransferReceiptDBData{
						OrderState:    result.OrderState,
						ProcessStatus: processStatus,
					}
					if result.ErrorCode != "" && result.ErrorCode != "000000" {
						updateReceipt.RetCode = result.ErrorCode
						updateReceipt.RetMessage = result.ErrorMessage
					}
					if err = s.bankTransferReceiptRepo.UpdateById(ctx, transferReceipt.Id, updateReceipt); err != nil {
						return handler.HandleError(err)
					}
					//推送钉钉成功消息
					s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
						Business: kafka.ProcessFinanceTransferResultBusiness,
						Type:     kafka.DingtalkType,
						Id:       transferReceipt.Id,
					})
					//更新关联审批任务的总状态
					if err = s.dingtalkClient.UpdateProcessInstanceStatus(ctx, &dingtalkApi.ProcessInstanceData{
						TransferReceiptState: result.OrderState,
						TotalStatus:          processStatus,
					}, transferReceipt.ProcessBusinessId); err != nil {
						return handler.HandleError(err)
					}
				}
			}
		}
	}
	return nil
}
func (s *bankService) HandleSPDBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("HandleSPDBankSyncTransferReceipt 查询浦发所有未完成的付款单据")
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRejectResult}
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
	transferReceiptList, _, err := s.bankTransferReceiptRepo.List(ctx, "", 0, 0, &repo.BankTransferReceiptDBDataParam{
		ExcludeOrderStates:    excludeOrderStates,
		IsOrderFlowNoNotEmpty: true,
		BankTransferReceiptDBData: repo.BankTransferReceiptDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: organizationId,
			},
			ProcessStatus: enum.ProcessInstanceTotalStatusRunning,
		},
		CreateTime:     createTimeParam,
		PayAccountType: bankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if transferReceiptList != nil && len(*transferReceiptList) > 0 {
		for _, transferReceipt := range *transferReceiptList {
			beginDate2 := stru.FormatDayTime(transferReceipt.CreatedAt)
			endDate2 := stru.FormatDayTime(time.Now())
			if beginDate2 == endDate {
				endDate2 = stru.FormatDayTimeStamp(time.Now().Unix() + 86400)
			}
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: transferReceipt.OrganizationId,
				Type:           bankType,
			})
			if err != nil {
				return handler.HandleError(err)
			}
			bankTransferResultResponseData, err := s.spdBankSDK.QueryTransferResult(ctx, transferReceipt.PayAccount, transferReceipt.RecAccount, transferReceipt.OrderFlowNo, transferReceipt.SerialNo, beginDate2, endDate2, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
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
					if result.TransStatus == transferReceipt.OrderState {
						continue
					}
					//更新单据
					processStatus := enum.ProcessInstanceTotalStatusFinish
					if result.TransStatus == enum.GuilinBankTransferSuccessResult {
						processStatus = enum.ProcessInstanceTotalStatusSuccess
					}
					updateReceipt := &repo.BankTransferReceiptDBData{
						OrderState:    result.TransStatus,
						ProcessStatus: processStatus,
					}
					if bankTransferResultResponseData.ReturnCode != "" && bankTransferResultResponseData.ReturnCode != enum.SPDBankSuccessRetCode {
						updateReceipt.RetCode = bankTransferResultResponseData.ReturnCode
						updateReceipt.RetMessage = bankTransferResultResponseData.ReturnMsg
					}
					if err = s.bankTransferReceiptRepo.UpdateById(ctx, transferReceipt.Id, updateReceipt); err != nil {
						return handler.HandleError(err)
					}
					//推送钉钉成功消息
					s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
						Business: kafka.ProcessFinanceTransferResultBusiness,
						Type:     kafka.DingtalkType,
						Id:       transferReceipt.Id,
					})
					//更新关联审批任务的总状态
					if err = s.dingtalkClient.UpdateProcessInstanceStatus(ctx, &dingtalkApi.ProcessInstanceData{
						TransferReceiptState: result.TransStatus,
						TotalStatus:          processStatus,
					}, transferReceipt.ProcessBusinessId); err != nil {
						return handler.HandleError(err)
					}
				}
			}
		}
	}
	return nil
}
func (s *bankService) HandlePinganBankSyncTransferReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	//查询所有未完成的付款单据
	zap.L().Info("HandlePinganBankSyncTransferReceipt 查询平安所有未完成的付款单据") //因为平安中间审审核环节可以最长10天,所以对于未完成的付款单据不能按照时间范围来选择
	excludeOrderStates := []string{enum.GuilinBankTransferSuccessResult, enum.GuilinBankTransferFailResult, enum.GuilinBankTransferRevokeResult, enum.GuilinBankTransferDeleteResult, enum.GuilinBankTransferRejectResult}
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
	transferReceiptList, _, err := s.bankTransferReceiptRepo.List(ctx, "", 0, 0, &repo.BankTransferReceiptDBDataParam{
		ExcludeOrderStates: excludeOrderStates,
		BankTransferReceiptDBData: repo.BankTransferReceiptDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: organizationId,
			},
			ProcessStatus: enum.ProcessInstanceTotalStatusRunning,
		},
		CreateTime:     createTimeParam,
		PayAccountType: bankType,
	})
	if err != nil {
		zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt查询付款单据失败:%v", err))
		return handler.HandleError(err)
	}
	if transferReceiptList != nil && len(*transferReceiptList) > 0 {
		zap.L().Info(fmt.Sprintf("pinganBankSDK.QueryTransferResult当前批次处理平安%v条数据,data=%v", len(*transferReceiptList), transferReceiptList))
		for _, transferReceipt := range *transferReceiptList {

			bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
				OrganizationId: transferReceipt.OrganizationId,
				Type:           bankType,
				Account:        transferReceipt.PayAccount,
			})
			if err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt查询账号详情失败:%v", err))
				continue
			}

			result, err := s.pinganBankSDK.QueryTransferResult(ctx, transferReceipt.SerialNo, bankAccount.ZuId, bankAccount.Account)
			zap.L().Info(fmt.Sprintf("pinganBankSDK.QueryTransferResult查询转账结果:%v", result))
			if err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt查询转账结果失败:%v", err))
				continue
			}

			// 处理订单状态(默认失败)
			//orderState := enum.GuilinBankTransferFailResult
			// 获取映射的订单状态
			result.Stt = enum.GetOrderState(result.Stt, result.Stt)
			zap.L().Info(fmt.Sprintf("enum.GetOrderState获取平安的状态:%v", result.Stt))
			//比较订单状态
			if result.Stt == transferReceipt.OrderState {
				zap.L().Info(fmt.Sprintf("transferReceipt.OrderState当前交易状态没有更新result.Stt:%v,transferReceipt.OrderState:%v", result.Stt, transferReceipt.OrderState))
				continue
			}
			chargeF, _ := strconv.ParseFloat(result.Fee, 64)
			//更新单据
			processStatus := enum.ProcessInstanceTotalStatusRunning
			if result.Stt == enum.GuilinBankTransferSuccessResult {
				processStatus = enum.ProcessInstanceTotalStatusSuccess
			}
			if result.Stt == enum.GuilinBankTransferFailResult {
				processStatus = enum.ProcessInstanceTotalStatusFinish
			}
			updateReceipt := &repo.BankTransferReceiptDBData{
				OrderState:       result.Stt,
				ProcessStatus:    processStatus,
				RetCode:          result.Yhcljg,
				RetMessage:       result.BackRem,
				DetailHostFlowNo: result.FrontLogNo,
				OrderFlowNo:      result.HostFlowNo,
				ChargeFee:        chargeF,
			}

			if err = s.bankTransferReceiptRepo.UpdateById(ctx, transferReceipt.Id, updateReceipt); err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt更新付款详情失败:%v", err))
				continue
			}
			//推送钉钉成功消息
			s.kafkaProducer.Send(kafka.BankTopic, kafka.TypeMessage{
				Business: kafka.ProcessFinanceTransferResultBusiness,
				Type:     kafka.DingtalkType,
				Id:       transferReceipt.Id})
			//更新关联审批任务的总状态
			if err = s.dingtalkClient.UpdateProcessInstanceStatus(ctx, &dingtalkApi.ProcessInstanceData{
				TransferReceiptState: result.Stt,
				TotalStatus:          processStatus}, transferReceipt.ProcessBusinessId); err != nil {
				zap.L().Info(fmt.Sprintf("HandlePinganBankSyncTransferReceipt关联审批任务的总状态:%v", err))
				continue
			}
		}
	}
	return nil
}

func (s *bankService) UpdateBankTransactionRecDetail(ctx context.Context, req *api.BankTransactionRecDetailData) error {
	dbData, err := s.bankTransactionDetailProcessInstanceRepo.Get(ctx, &repo.BankTransactionDetailProcessInstanceDBData{
		ExternalId: req.ExternalId,
	})
	if err != nil || dbData == nil {
		return err
	}
	detailId := dbData.BankTransactionDetailId
	if err = s.bankTransactionDetailRepo.UpdateById(ctx, detailId, &repo.BankTransactionDetailDBData{
		ProcessBusinessId:  req.ProcessBusinessId,
		ProcessInstanceId:  req.ProcessInstanceId,
		OperationUserId:    req.OperationUserId,
		OperationUserName:  req.OperationUserName,
		OperationComment:   req.OperationComment,
		ProcessTotalStatus: req.ProcessTotalStatus,
	}); err != nil {
		return err
	}
	return nil
}

func (s *bankService) SyncTransferReceipt(ctx context.Context, taskId int64, param []byte, organizationId int64) error {
	var syncDateRequest stru.SyncDateRequest
	json.Unmarshal(param, &syncDateRequest)
	if syncDateRequest.BeginDate == "" || syncDateRequest.EndDate == "" {
		return s.baseClient.FailTask(ctx, taskId, "BeginDate or EndDate is empty")
	}
	if err := s.HandleSyncTransferReceipt(ctx, syncDateRequest.BeginDate, syncDateRequest.EndDate, organizationId); err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	return s.baseClient.SuccessTask(ctx, taskId, "")
}

func (s *bankService) SyncTransactionDetail(ctx context.Context, taskId int64, param []byte, organizationId int64) error {
	var syncDateRequest stru.SyncDateRequest
	json.Unmarshal(param, &syncDateRequest)
	if syncDateRequest.BeginDate == "" || syncDateRequest.EndDate == "" {
		return s.baseClient.FailTask(ctx, taskId, "BeginDate or EndDate is empty")
	}
	if err := s.HandleTransactionDetail(ctx, syncDateRequest.BeginDate, syncDateRequest.EndDate, organizationId); err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	return s.baseClient.SuccessTask(ctx, taskId, "")
}

func (s *bankService) DashboardData(ctx context.Context, organizationId int64) (*api.DashboardData, error) {
	var dashboardData api.DashboardData
	//近15日
	halfMonthRange := stru.FormatDayTimeCondition(-15)
	halfMonthDayArray, halfMonthPayAmountArray, halfMonthRecAmountArray, err := s.bankTransactionDetailRepo.CashFlowCount(halfMonthRange, organizationId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	halfMonthTimeLabels := stru.FormatWeekChartTimeLabelData(-15)
	dashboardData.DayFlowData = &api.ChartData{
		Labels:   halfMonthTimeLabels,
		Datasets: stru.FormatLineChartMultipleDataSetArray(halfMonthTimeLabels, *halfMonthDayArray, *halfMonthPayAmountArray, *halfMonthRecAmountArray),
	}

	//近7日
	weekTimeRange := stru.FormatDayTimeCondition(-7)
	dayArray, payAmountArray, recAmountArray, err := s.bankTransactionDetailRepo.CashFlowCount(weekTimeRange, organizationId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	weekTimeLabels := stru.FormatWeekChartTimeLabelData(-7)
	dashboardData.WeekFlowData = &api.ChartData{
		Labels:   weekTimeLabels,
		Datasets: stru.FormatLineChartMultipleDataSetArray(weekTimeLabels, *dayArray, *payAmountArray, *recAmountArray),
	}

	weekTimeRange2 := stru.FormatDayTimeCondition(-14) //第一天的余额数据可能为空 多预查一周
	weekBalanceDataMap, err := s.bankTransactionDetailRepo.BalanceFlowCountMap(weekTimeRange2, organizationId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	weekTimeLabels2 := stru.FormatWeekChartTimeLabelData(-14)
	weekBalanceLabelArray, weekBalanceDataArray := stru.HandleBalanceData(weekTimeLabels2, weekBalanceDataMap, 7)
	dashboardData.WeekBalanceData = &api.ChartData{
		Labels:   weekTimeLabels,
		Datasets: stru.FormatLineChartSingleDataSetArray(weekTimeLabels, *weekBalanceLabelArray, *weekBalanceDataArray),
	}

	//当月
	year, month, _ := time.Now().Date()
	monthTimeRange := stru.FormatMonthTimeCondition(year, month)
	monthDayArray, monthPayAmountArray, monthRecAmountArray, err := s.bankTransactionDetailRepo.CashFlowCount(monthTimeRange, organizationId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	monthTimeLabels := stru.FormatMonthChartTimeLabelData(year, month)
	dashboardData.MonthFlowData = &api.ChartData{
		Labels:   monthTimeLabels,
		Datasets: stru.FormatLineChartMultipleDataSetArray(monthTimeLabels, *monthDayArray, *monthPayAmountArray, *monthRecAmountArray),
	}

	monthTimeRange2 := stru.FormatMonthTimeConditionMore(year, month, -7) //第一天的余额数据可能为空 多预查一周
	monthBalanceDataMap, err := s.bankTransactionDetailRepo.BalanceFlowCountMap(monthTimeRange2, organizationId)
	monthTimeLabels2, dayCount := stru.FormatMonthChartTimeLabelDataMore(year, month, -7)
	monthBalanceLabelArray, monthBalanceDataArray := stru.HandleBalanceData(monthTimeLabels2, monthBalanceDataMap, dayCount)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	dashboardData.MonthBalanceData = &api.ChartData{
		Labels:   monthTimeLabels,
		Datasets: stru.FormatLineChartSingleDataSetArray(monthTimeLabels, *monthBalanceLabelArray, *monthBalanceDataArray),
	}

	return &dashboardData, nil
}

func (s *bankService) GetCashFlowMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (*api.ChartData, error) {
	month, err := strconv.Atoi(req.Month)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	year, err := strconv.Atoi(req.Year)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	monthTime := stru.FormatMonthTime(month)
	monthTimeRange := stru.FormatMonthTimeCondition(year, monthTime)
	monthDayArray, monthPayAmountArray, monthRecAmountArray, err := s.bankTransactionDetailRepo.CashFlowCount(monthTimeRange, req.OrganizationId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	monthTimeLabels := stru.FormatMonthChartTimeLabelData(year, monthTime)
	return &api.ChartData{
		Labels:   monthTimeLabels,
		Datasets: stru.FormatLineChartMultipleDataSetArray(monthTimeLabels, *monthDayArray, *monthPayAmountArray, *monthRecAmountArray),
	}, nil
}

func (s *bankService) GetBalanceMonthChartData(ctx context.Context, req *api.MonthChartDataRequest) (*api.ChartData, error) {
	month, err := strconv.Atoi(req.Month)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	year, err := strconv.Atoi(req.Year)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	monthTime := stru.FormatMonthTime(month)
	monthTimeRange2 := stru.FormatMonthTimeConditionMore(year, monthTime, -7) //第一天的余额数据可能为空 多预查一周
	monthBalanceDataMap, err := s.bankTransactionDetailRepo.BalanceFlowCountMap(monthTimeRange2, req.OrganizationId)
	monthTimeLabels2, dayCount := stru.FormatMonthChartTimeLabelDataMore(year, monthTime, -7)
	monthBalanceLabelArray, monthBalanceDataArray := stru.HandleBalanceData(monthTimeLabels2, monthBalanceDataMap, dayCount)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	monthTimeLabels := stru.FormatMonthChartTimeLabelData(year, monthTime)
	return &api.ChartData{
		Labels:   monthTimeLabels,
		Datasets: stru.FormatLineChartSingleDataSetArray(monthTimeLabels, *monthBalanceLabelArray, *monthBalanceDataArray),
	}, nil
}

func (s *bankService) QueryAccountBalance(ctx context.Context, req *api.QueryAccountBalanceRequest) (*api.QueryAccountBalanceResponse, error) {
	// 先根据accountNo 查询是哪个银行类型的, 然后再调用不同银行的接口
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		OrganizationId: req.OrganizationId,
		Account:        req.AccountNo,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}

	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: req.OrganizationId,
		Type:           bankAccount.Type,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}

	if bankAccount.Type == enum.GuilinBankType {
		return s.QueryGuilinBankAccountBalance(ctx, req.AccountNo, organizationBankConfig)
	} else if bankAccount.Type == enum.SPDBankType {
		return s.QuerySPDBankAccountBalance(ctx, req.AccountNo, organizationBankConfig)
	} else if bankAccount.Type == enum.PinganBankType {
		return s.QueryPinganBankAccountBalance(ctx, req.AccountNo, bankAccount.ZuId)
	}
	return nil, nil
}

func (s *bankService) QueryGuilinBankAccountBalance(ctx context.Context, accountNo string, organizationBankConfig *baseApi.OrganizationBankConfigData) (*api.QueryAccountBalanceResponse, error) {
	res, err := s.guilinBankSDK.QueryAccountBalance(ctx, accountNo, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if res.Body.AccountNo == "" {
		return &api.QueryAccountBalanceResponse{
			Success: false,
			Msg:     "未找到该账号信息",
		}, nil
	}
	return &api.QueryAccountBalanceResponse{
		Success: true,
		Balance: res.Body.Balance,
	}, nil
}

func (s *bankService) QuerySPDBankAccountBalance(ctx context.Context, accountNo string, organizationBankConfig *baseApi.OrganizationBankConfigData) (*api.QueryAccountBalanceResponse, error) {
	res, err := s.spdBankSDK.QueryAccountBalance(ctx, accountNo, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if (res == nil || len(res.Items) <= 0) || (res != nil && len(res.Items) > 0 && res.Items[0].AcctNo == "") {
		return &api.QueryAccountBalanceResponse{
			Success: false,
			Msg:     "未找到该账号信息",
		}, nil
	}
	// res返回是一个list 取第一个
	return &api.QueryAccountBalanceResponse{
		Success: true,
		Balance: res.Items[0].Balance,
	}, nil
}
func (s *bankService) QueryPinganBankAccountBalance(ctx context.Context, accountNo string, zuId string) (*api.QueryAccountBalanceResponse, error) {
	res, err := s.pinganBankSDK.QueryAccountBalance(ctx, accountNo,
		config.GetString(bankEnum.PinganJsdkUrl, ""),
		config.GetString(bankEnum.PinganMrchCode, ""), zuId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if res == nil || res.Account == "" {
		return &api.QueryAccountBalanceResponse{
			Success: false,
			Msg:     "未找到该账号信息",
		}, nil
	}
	acctBalance, err := strconv.ParseFloat(res.AcctBalance, 64)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return &api.QueryAccountBalanceResponse{
		Success: true,
		Balance: acctBalance,
	}, nil
}

func (s *bankService) ImportBankBusinessPayrollData(ctx context.Context, taskId int64, param []byte, organizationId int64) error {
	var importDateRequest stru.ImportPayrollDataRequest
	json.Unmarshal(param, &importDateRequest)
	if importDateRequest.PayAccountNo == "" || importDateRequest.Month == "" {
		return s.baseClient.FailTask(ctx, taskId, "Param is empty")
	}
	var fileData [][]string
	err := json.Unmarshal(importDateRequest.FileData, &fileData)
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: organizationId,
		Type:           "0",
	})
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	//sdk 上送批次
	uploadPayInfoRes, err := s.guilinBankSDK.UploadBatchTransferPayInfo(ctx, sdkStru.UploadBatchTransferPayInfoRequest{
		TotalNumber:    strconv.FormatInt(importDateRequest.TotalCount, 10),
		TotalAmount:    importDateRequest.TotalAmount,
		PayAccount:     importDateRequest.PayAccountNo,
		PayAccountName: importDateRequest.PayAccountName,
	}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	if uploadPayInfoRes.Head.RetMessage != "" || uploadPayInfoRes.Body.BatchNo == "" {
		return s.baseClient.FailTask(ctx, taskId, uploadPayInfoRes.Head.RetMessage)
	}
	batchNo := uploadPayInfoRes.Body.BatchNo
	//入库 批次
	createUser, _ := util.GetMetaInfoCurrentUserName(ctx)
	payrollId, err := s.businessPayrollRepo.Add(ctx, &repo.BankBusinessPayrollDBData{
		BaseDBData: repository.BaseDBData{
			OrganizationId: organizationId,
		},
		Name:            fmt.Sprintf("%s月代发工资", importDateRequest.Month),
		PayAccountName:  importDateRequest.PayAccountName,
		PayAccountNo:    importDateRequest.PayAccountNo,
		Month:           importDateRequest.Month,
		Remark:          importDateRequest.Remark,
		Count:           importDateRequest.TotalCount,
		TotalMoney:      importDateRequest.TotalAmount,
		Status:          1,
		CreatedUserName: createUser,
		FileUrl:         importDateRequest.FileUrl,
		BatchNo:         batchNo,
	})
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	//sdk 上送交易信息
	var recAccountData []sdkStru.UploadBatchTransferRecInfoList
	var payrollDetailData []repo.BankBusinessPayrollDetailDBData
	for i := 1; i < len(fileData); i++ {
		accountName := fileData[i][0]
		accountNo := fileData[i][1]
		amount, pErr := strconv.ParseFloat(fileData[i][2], 64)
		if pErr != nil {
			return s.baseClient.FailTask(ctx, taskId, err.Error())
		}
		payUse := fileData[i][3]
		recNum := fileData[i][4]
		//serialNo := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		serialNo, _ := util.SonyflakeID()
		recAccountData = append(recAccountData, sdkStru.UploadBatchTransferRecInfoList{
			SerialNo:       serialNo,
			RecAccount:     accountNo,
			RecAccountName: accountName,
			PayAmount:      amount,
			PayUse:         payUse,
			BusinessCode:   "020104",
			CurrencyType:   "CNY",
			RmtType:        "0",
			RecBankType:    "0",
			PubPriFlag:     "1",
		})
		payrollDetailData = append(payrollDetailData, repo.BankBusinessPayrollDetailDBData{
			BaseDBData: repository.BaseDBData{
				OrganizationId: organizationId,
			},
			SerialNo:       serialNo,
			BatchId:        payrollId,
			RecAccountName: accountName,
			RecAccountNo:   accountNo,
			Amount:         amount,
			Month:          importDateRequest.Month,
			Num:            recNum,
			Remark:         payUse,
		})
	}
	uploadRecInfoRes, err := s.guilinBankSDK.UploadBatchTransferRecInfo(ctx, batchNo, recAccountData, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	if uploadRecInfoRes.Head.RetMessage != "" {
		return s.baseClient.FailTask(ctx, taskId, uploadRecInfoRes.Head.RetMessage)
	}
	//入库 批次明细
	err = s.businessPayrollDetailRepo.Transaction(func(tx *gorm.DB) error {
		_, err = s.businessPayrollDetailRepo.BatchAddWithTX(ctx, tx, &payrollDetailData)
		if err != nil {
			return handler.HandleError(err)
		}
		return nil
	})
	if err != nil {
		return s.baseClient.FailTask(ctx, taskId, err.Error())
	}
	//更新批次订单状态
	err = s.businessPayrollRepo.UpdateById(ctx, payrollId, &repo.BankBusinessPayrollDBData{
		State:           uploadRecInfoRes.Body.OrderState,
		Msg:             uploadRecInfoRes.Body.ErrorMessage,
		OrderSubmitTime: uploadRecInfoRes.Body.OrderSubmitTime,
	})
	if err != nil {
		return err
	}
	return s.baseClient.SuccessTask(ctx, taskId, "")
}

func (s *bankService) ListBankBusinessPayroll(ctx context.Context, req *api.ListBusinessPayrollRequest) (*api.ListBusinessPayrollResponse, error) {
	dbData, count, err := s.businessPayrollRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankBusinessPayrollDBDataParam{
		CreateTimeRange: req.CreatedTime,
		BankBusinessPayrollDBData: repo.BankBusinessPayrollDBData{
			PayAccountNo:   req.PayAccountNo,
			PayAccountName: req.PayAccountName,
			Month:          req.Month,
			Status:         req.Status,
			BaseDBData: repository.BaseDBData{
				BaseCommonDBData: repository.BaseCommonDBData{
					CreatedUserId: req.CreateUser,
				},
			},
		},
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	var data []*api.BusinessPayrollListVo
	if count > 0 && dbData != nil {
		data = make([]*api.BusinessPayrollListVo, len(*dbData))
		for i, v := range *dbData {
			createdAt := v.CreatedAt.Format("2006-01-02 15:04:05")
			data[i] = &api.BusinessPayrollListVo{
				Id:             v.Id,
				Name:           v.Name,
				PayAccountNo:   v.PayAccountNo,
				PayAccountName: v.PayAccountName,
				Month:          v.Month,
				Remark:         v.Remark,
				TotalCount:     v.Count,
				TotalAmount:    v.TotalMoney,
				Status:         v.Status,
				CreatedUser:    v.CreatedUserName,
				CreatedAt:      createdAt,
				State:          v.State,
				Msg:            v.Msg,
				SubmitTime:     v.OrderSubmitTime,
			}
		}
	}

	return &api.ListBusinessPayrollResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) ListBankBusinessPayrollDetail(ctx context.Context, req *api.ListBusinessPayrollDetailRequest) (*api.ListBusinessPayrollDetailResponse, error) {
	dbData, count, err := s.businessPayrollDetailRepo.List(ctx, req.Sort, req.PageNum, req.PageSize, &repo.BankBusinessPayrollDetailDBData{
		BatchId:        req.BatchId,
		Num:            req.Num,
		Month:          req.Month,
		RecAccountName: req.RecAccountName,
		RecAccountNo:   req.RecAccountNo,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	var data []*api.BusinessPayrollDetailListVo
	if count > 0 && dbData != nil {
		data = make([]*api.BusinessPayrollDetailListVo, len(*dbData))
		for i, v := range *dbData {
			data[i] = &api.BusinessPayrollDetailListVo{
				Id:             v.Id,
				Num:            v.Num,
				RecAccountName: v.RecAccountName,
				RecAccountNo:   v.RecAccountNo,
				Amount:         v.Amount,
				Month:          v.Month,
				OrderState:     v.OrderState,
				ErrorMessage:   v.ErrorMessage,
				ErrorCode:      v.ErrorCode,
				BatchId:        v.BatchId,
				Remark:         v.Remark,
			}
		}
	}
	return &api.ListBusinessPayrollDetailResponse{
		Data:  data,
		Count: count,
	}, nil
}

func (s *bankService) SyncBankBusinessPayrollDetail(ctx context.Context, req *api.SyncBusinessPayrollResultRequest) (*api.SyncBusinessPayrollResultResponse, error) {
	dbData, err := s.businessPayrollRepo.Get(ctx, &repo.BankBusinessPayrollDBData{
		BaseDBData: repository.BaseDBData{
			BaseCommonDBData: repository.BaseCommonDBData{
				Id: req.Id,
			},
		},
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	//判断订单状态 90：交易成功 99：交易失败
	if dbData.State == "90" || dbData.State == "99" {
		return &api.SyncBusinessPayrollResultResponse{
			Success: false,
			Msg:     "批量转账订单已结束",
		}, nil
	}
	detailDataList, count, err := s.businessPayrollDetailRepo.List(ctx, "", 0, 0, &repo.BankBusinessPayrollDetailDBData{
		BatchId: req.Id,
		Month:   dbData.Month,
	})
	if count <= 0 {
		return &api.SyncBusinessPayrollResultResponse{
			Success: false,
			Msg:     "该批次下无批量转账明细",
		}, nil
	}
	transferMap := make(map[string]repo.BankBusinessPayrollDetailDBData)
	for _, item := range *detailDataList {
		transferMap[item.SerialNo] = item
	}

	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: dbData.OrganizationId,
		Type:           "0",
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	//sdk 上送批次
	orderSubmitTime := dbData.CreatedAt.Format("20060102")
	if dbData.OrderSubmitTime != "" {
		orderSubmitTime = dbData.OrderSubmitTime[0:8]
	}

	batchTransferRes, err := s.guilinBankSDK.QueryBatchTransferResult(ctx, sdkStru.QueryBatchTransferResultRequestBodyData{
		SearchPayAccount: dbData.PayAccountNo,
		BeginDate:        orderSubmitTime,
		EndDate:          orderSubmitTime,
		BatchNo:          dbData.BatchNo,
	}, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	batchOrderState := "90"
	totalSuccessCount := 0
	if len(batchTransferRes) > 0 {
		for _, transfer := range batchTransferRes {
			item, ok := transferMap[transfer.OrderFlowNo]
			if !ok {
				continue
			}
			if transfer.OrderState == "90" {
				totalSuccessCount++
			}
			if item.OrderState == "90" || item.OrderState == "99" {
				continue
			}
			err = s.businessPayrollDetailRepo.UpdateById(ctx, item.Id, &repo.BankBusinessPayrollDetailDBData{
				OrderState:   transfer.OrderState,
				ErrorCode:    transfer.ErrorCode,
				ErrorMessage: transfer.ErrorMessage,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	//更新转账批次的订单状态
	if totalSuccessCount != len(batchTransferRes) {
		batchOrderState = fmt.Sprintf("交易结束(%d/%d)", totalSuccessCount, dbData.Count)
	}
	err = s.businessPayrollRepo.UpdateById(ctx, req.Id, &repo.BankBusinessPayrollDBData{
		State: batchOrderState,
	})
	if err != nil {
		return nil, err
	}
	return &api.SyncBusinessPayrollResultResponse{
		Success: true,
	}, nil
}

// HandleTransactionDetailReceipt
//
//	@Description: 桂林银行下载电子凭证可以直接下载, 浦发银行需要申请,然后过几分钟才能下载, 这里只处理了浦发
//	@receiver s
//	@param ctx
//	@param beginDate
//	@param endDate
//	@param organizationId
//	@return error
func (s *bankService) HandleTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {
	value := beginDate + "-" + endDate + ":" + strconv.FormatInt(organizationId, 10) + "==" + util.FormatDateTime(time.Now())
	s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
	s.HandleSPDTransactionDetailReceipt(ctx, enum.SPDBankType, beginDate, endDate, organizationId)

	s.HandlePinganTransactionDetailReceipt(ctx, beginDate, endDate, organizationId)
	return nil
}

func (s *bankService) HandleSPDTransactionDetailReceipt(ctx context.Context, bankType, beginDate string, endDate string, organizationId int64) error {
	// 查找所有 bank_transaction_detail 表 pay_account_type 是浦发 且 electronic_receipt_file 字段为空的数据
	value := util.FormatDateTime(time.Now()) + "开始处理数据"
	s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)

	dbDatas, count, err := s.bankTransactionDetailRepo.List(ctx, "", 0, 0, &repo.BankTransactionDetailDBDataParam{
		BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
			PayAccountType: bankType,
			//Type:           "pay",
		},
		IsElectronicReceiptFileNull: true,
		IsAccountNoNull:             false,
		TransferTimeArray:           []string{beginDate, endDate},
	})
	if err != nil || dbDatas == nil {
		value = fmt.Sprintf("%s=当前spd-s.bankTransactionDetailRepo.List方法出错或查询结果为空-err=%+v,count=%d", util.FormatDateTime(time.Now()), err, count)
		s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
		return handler.HandleError(err)
	}
	zap.L().Info(fmt.Sprintf("s.bankService.HandleSPDTransactionDetailReceipt_count:%v", count))
	value = fmt.Sprintf("%s=当前spd需要处理:%d条数据", util.FormatDateTime(time.Now()), count)
	s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
	for _, dbData := range *dbDatas {
		go func() {
			newDbData := dbData
			value = fmt.Sprintf("%s id=%d当前处理的交易明细量:%+v", util.FormatDateTime(time.Now()), dbData.Id, dbData)
			s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
			//查询银行账号
			bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
				OrganizationId: newDbData.OrganizationId,
				Type:           enum.SPDBankType,
			})
			value = fmt.Sprintf("%s=当前spd-SimpleGetOrganizationBankAccount方法查询结果:%+v", util.FormatDateTime(time.Now()), bankAccount)
			s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
			if err != nil {
				value = fmt.Sprintf("%s=当前spd-SimpleGetOrganizationBankAccount方法出错:%+v", util.FormatDateTime(time.Now()), err)
				s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				return
			}
			if bankAccount == nil {
				value = fmt.Sprintf("%s=当前spd-SimpleGetOrganizationBankAccount方法bankAccount为空", util.FormatDateTime(time.Now()))
				s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				return
			}
			organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
				OrganizationId: organizationId,
				Type:           bankType,
			})
			value = fmt.Sprintf("%s=当前spd-GetOrganizationBankConfig方法查询结果:%+v", util.FormatDateTime(time.Now()), organizationBankConfig)
			s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
			if err != nil {
				value = fmt.Sprintf("%s=当前spd-GetOrganizationBankConfig方法出错:%+v", util.FormatDateTime(time.Now()), err)
				s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				return
			}

			// 更新凭证地址
			// dbData.OrderFlowNo 是 浦发的柜员流水号: 对应交易明细的 tellerJnlNo
			// dbData.ExtField1 是 浦发的传票组内序号: 对应交易明细的 summonsNumber

			//f, err := s.spdBankSDK.DownloadTransactionDetailElectronicReceipt(ctx, bankAccount.Account, beginDate, endDate, newDbData.OrderFlowNo, newDbData.ExtField1, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.FileHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			f, err := s.spdBankSDK.DownloadTransactionDetailElectronicReceipt(ctx, bankAccount.Account, newDbData.TransferDate, newDbData.TransferDate, newDbData.OrderFlowNo, newDbData.ExtField1, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.FileHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
			if err != nil {
				str := fmt.Sprintf("bankAccount.Account=%s,newDbData.TransferDate=%s,newDbData.TransferDate=%s,newDbData.OrderFlowNo=%s,newDbData.ExtField1=%s,organizationBankConfig.Host=%s,organizationBankConfig.SignHost=%s,organizationBankConfig.FileHost=%s,organizationBankConfig.BankCustomerId=%s,organizationBankConfig.BankUserId=%s",
					bankAccount.Account, newDbData.TransferDate, newDbData.TransferDate, newDbData.OrderFlowNo, newDbData.ExtField1, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.FileHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
				value = fmt.Sprintf("%sid=%d下载浦发电子凭证失败:%+v,参数=%s", util.FormatDateTime(time.Now()), dbData.Id, err, str)
				s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				zap.L().Info(fmt.Sprintf("s.bankService.spdBankSDK 下载浦发电子凭证失败: %v\n", err.Error()))
				return
			}
			value = fmt.Sprintf("%sid=%d 下载浦发电子凭证成功,OrderFlowNo=%s", util.FormatDateTime(time.Now()), newDbData.Id, newDbData.OrderFlowNo)
			s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)

			var electronicReceiptFile string
			if f != nil && len(f) > 0 {
				electronicReceiptFile, err = store.UploadOSSFileBytes("pdf", ".pdf", f, s.ossConfig, false)
				if err != nil {
					value = fmt.Sprintf("%sid=%d上传浦发电子凭证失败:%+v", util.FormatDateTime(time.Now()), newDbData.Id, err)
					s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
					zap.L().Error(fmt.Sprintf("上传浦发电子凭证到OSS失败: %v\n", err.Error()))
					return
				} else {
					value = fmt.Sprintf("%sid=%d上传浦发电子凭证成功:%s", util.FormatDateTime(time.Now()), newDbData.Id, electronicReceiptFile)
					s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				}
				err := s.bankTransactionDetailRepo.UpdateById(ctx, newDbData.Id, &repo.BankTransactionDetailDBData{
					ElectronicReceiptFile: electronicReceiptFile,
				})
				if err != nil {
					value = fmt.Sprintf("%sid=%d更新浦发电子凭证失败:%+v", util.FormatDateTime(time.Now()), newDbData.Id, err)
					s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
					zap.L().Error(fmt.Sprintf("更新浦发电子凭证失败: %v\n", err.Error()))
					return
				} else {
					value = fmt.Sprintf("%sid=%d更新浦发电子凭证成功", util.FormatDateTime(time.Now()), newDbData.Id)
					s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
				}

				//更新单据的明细ID
				if err = s.updateRelevanceElectronicDocument(ctx, newDbData.Summary, newDbData.HostFlowNo, electronicReceiptFile, enum.SPDBankType); err != nil {
					value = fmt.Sprintf("%s更新浦发单据电子凭证失败:%+v", util.FormatDateTime(time.Now()), err)
					s.setRedisLog(ctx, bankEnum.BankReceiptSyncLogKey, value)
					zap.L().Error(fmt.Sprintf("更新浦发更新单据的回单失败: %v\n", err.Error()))
				}

				zap.L().Info("HandleSPDTransactionDetailReceipt 处理浦发电子凭证成功")
			}
		}()
	}
	return nil
}

func (s *bankService) HandlePinganTransactionDetailReceipt(ctx context.Context, beginDate string, endDate string, organizationId int64) error {

	dbDatas, count, err := s.bankTransactionDetailRepo.List(ctx, "", 0, 0, &repo.BankTransactionDetailDBDataParam{
		BankTransactionDetailDBData: repo.BankTransactionDetailDBData{
			PayAccountType: enum.PinganBankType,
		},
		IsElectronicReceiptFileNull: true,
		TransferTimeArray:           []string{beginDate, endDate},
	})
	if err != nil || dbDatas == nil {
		return handler.HandleError(err)
	}
	zap.L().Info(fmt.Sprintf("s.bankService.本次下载平安电子回单数量_count:%v", count))

	for _, dbData := range *dbDatas {
		go func() {
			newDbData := dbData
			//获取付款账户地址
			bankAccount, _ := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
				Id: newDbData.MerchantAccountId,
			})
			// 更新凭证地址
			//封装请求body
			serialNo, _ := util.SonyflakeID()
			outAccountNo := ""
			if newDbData.Type == "pay" {
				outAccountNo = bankAccount.Account
			} else {
				outAccountNo = newDbData.AcctNo
			}

			request := sdkStru.PinganSameDayHistoryReceiptDataQueryRequest{
				MrchCode:         config.GetString("bankConfig.pingAn.accountKeeper.mrchCode", ""),
				CnsmrSeqNo:       serialNo,
				OutAccNo:         outAccountNo,
				AccountBeginDate: newDbData.TransferDate,
				AccountEndDate:   newDbData.TransferDate,
				HostFlow:         newDbData.HostFlowNo,
			}
			zap.L().Info(fmt.Sprintf("PinganSameDayHistoryReceiptDataQueryRequest 下载pingan电子凭证参数: %v\n", request))
			f, err := s.pinganBankSDK.UploadTransactionDetailElectronic(ctx, request, bankAccount.ZuId)
			if err != nil {
				zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证失败: %v\n", err.Error()))
			}
			if f == "" {
				zap.L().Info(fmt.Sprintf("s.bankService.pinganBankSDK 下载pingan电子凭证失败,没有可供下载的zip文件"))
			} else {
				err = s.bankTransactionDetailRepo.UpdateById(ctx, newDbData.Id, &repo.BankTransactionDetailDBData{
					ElectronicReceiptFile: f,
				})
				if err != nil {
					zap.L().Error(fmt.Sprintf("更新平安电子凭证失败: %v\n", err.Error()))
				}

				//更新单据的回单
				if err = s.updateRelevanceElectronicDocument(ctx, "", newDbData.OrderFlowNo, f, enum.PinganBankType); err != nil {
					zap.L().Error(fmt.Sprintf("更新平安电子更新单据的回单失败: %v\n", err.Error()))
				}
				zap.L().Info("HandlePinganTransactionDetailReceipt 处理平安电子凭证成功")
			}
		}()
	}
	return nil
}

// CreateVirtualAccount
//
//	@Description: 创建虚账户
//	@receiver s
//	@param ctx
//	@param req
//	@return *api.CreateVirtualAccountResponse
//	@return error
func (s *bankService) CreateVirtualAccount(ctx context.Context, req *api.CreateVirtualAccountRequest) (*api.CreateVirtualAccountResponse, error) {
	organizationId := req.OrganizationId
	bankType := req.Type
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if bankAccount == nil {
		return nil, nil
	}

	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if organizationBankConfig == nil {
		return nil, nil
	}

	if bankType == enum.SPDBankType {
		// 从配置中心获取虚账户
		virtualAccountParentNo := config.GetString("virtual.account.parent.no", "")
		if virtualAccountParentNo == "" {
			return nil, handler.HandleError(errors.New("未配置虚账户实账号"))
		}
		virtualAccountParentName := config.GetString("virtual.account.parent.name", "")
		if virtualAccountParentName == "" {
			return nil, handler.HandleError(errors.New("未配置虚账户实账号名称"))
		}
		// 封装masterName
		var createAccountRequestItems []sdkStru.SPDCreateVirtualAccountRequestItem
		if req.Data != nil {
			createAccountRequestItems = make([]sdkStru.SPDCreateVirtualAccountRequestItem, 1)
			createAccountRequestItems[0] = sdkStru.SPDCreateVirtualAccountRequestItem{
				VirtualAccountName: fmt.Sprintf("%s（%s）", virtualAccountParentName, req.Data.VirtualAccountName),
				//VirtualAccountNo:   req.Data.VirtualAccountNo,
				Rate: req.Data.Rate,
			}
		}

		accountData, err := s.spdBankSDK.CreateVirtualAccount(ctx, createAccountRequestItems, virtualAccountParentNo, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		// 校验accountData
		if accountData == nil || (accountData.Items != nil && len(accountData.Items) > 0 && accountData.Items[0].VirtualAccountNo == "") {
			zap.L().Error(fmt.Sprintf("spdBankSDK.CreateVirtualAccount 返回虚账户数据为空"))
		}

		var data *api.VirtualAccountData
		if accountData.Items != nil && len(accountData.Items) > 0 {
			data = &api.VirtualAccountData{
				VirtualAccountName: accountData.Items[0].VirtualAccountName,
				VirtualAccountNo:   accountData.Items[0].VirtualAccountNo,
				Rate:               accountData.Items[0].Rate,
			}
		}
		// 保存到虚账户表
		err = s.baseClient.AddOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
			OrganizationId:     organizationId,
			Type:               bankType,
			AccountNo:          virtualAccountParentNo,
			VirtualAccountNo:   data.VirtualAccountNo,
			VirtualAccountName: data.VirtualAccountName,
		})
		if err != nil {
			zap.L().Error(fmt.Sprintf("CreateVirtualAccount error: %v\n", err.Error()))
			return nil, err
		}

		return &api.CreateVirtualAccountResponse{
			Status: accountData.BusinessStatus,
			Msg:    accountData.Explain,
			Data:   data,
		}, nil
	} else if bankType == "2" {

		accountDatas, err := s.baseClient.ListOrganizationBankVirtualAccountData(ctx, &baseApi.ListOrganizationBankVirtualAccountRequest{
			Type: bankType,
			//AccountNo:      bankAccount.Account,
			//VirtualAccountNo: virtualAccountNo,
		})
		if err != nil {
			return nil, handler.HandleError(err)
		}
		//生成随机6位数
	NEWSEQNUMBER:
		subAccountSeqNumber := util.Createcaptcha()
		//判断这个随机数是不是被用过,
		for _, accountData := range accountDatas {
			SeqNumber := accountData.VirtualAccountNo[len(accountData.VirtualAccountNo)-6:]
			if SeqNumber == subAccountSeqNumber {
				goto NEWSEQNUMBER
			}
		}

		// 封装masterName
		seqNoId, _ := util.SonyflakeID()
		if req.Data != nil {
			var createAccountRequest = sdkStru.PinganCreateVirtualAccountRequest{
				MrchCode:       config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
				CnsmrSeqNo:     seqNoId,
				MainAccount:    config.GetString(bankEnum.PinganIntelligenceAccountNo, ""), //智能账号,
				SubAccountSeq:  subAccountSeqNumber,                                        //清分台账编码序号
				SubAccountName: fmt.Sprintf("%s[%s]", bankAccount.AccountName, req.Data.VirtualAccountName),
				OpFlag:         "A",
			}
			accountData, err := s.pinganBankSDK.CreateVirtualAccount(ctx, createAccountRequest)
			if err != nil {
				return nil, handler.HandleError(err)
			}
			// 校验accountData
			if accountData == nil {
				zap.L().Error(fmt.Sprintf("pinganBankSDK.CreateVirtualAccount 返回虚账户数据为空"))
			}
			zap.L().Info(fmt.Sprintf("==CreateVirtualAccountResult%v", accountData))
			var data *api.VirtualAccountData
			rate, err := strconv.ParseFloat(accountData.Rate, 64)
			if err != nil {
				return nil, handler.HandleError(err)
			}
			data = &api.VirtualAccountData{
				VirtualAccountName: accountData.SubAccountName,
				VirtualAccountNo:   accountData.SubAccountNo,
				Rate:               rate,
			}
			// 保存到虚账户表
			err = s.baseClient.AddOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
				OrganizationId:     organizationId,
				Type:               bankType,
				AccountNo:          bankAccount.Account,
				VirtualAccountNo:   data.VirtualAccountNo,
				VirtualAccountName: data.VirtualAccountName,
			})
			if err != nil {
				zap.L().Error(fmt.Sprintf("CreateVirtualAccount error: %v\n", err.Error()))
				return nil, err
			}

			return &api.CreateVirtualAccountResponse{
				Status: "200",
				Msg:    "添加成功",
				Data:   data,
			}, nil
		}
	}
	return nil, nil
}

// QueryVirtualAccountBalance
//
//	@Description: 查询虚账户余额
//	@receiver s
//	@param ctx
//	@param accountNo
//	@param virtualAccountNo
//	@param organizationId
//	@return *api.VirtualAccountBalanceData
//	@return error
func (s *bankService) QueryVirtualAccountBalance(ctx context.Context, organizationId int64, bankType string) (*api.VirtualAccountBalanceData, error) {
	bankTypeStr := ""
	var resultVirtualAccountBalance float64
	if bankType == enum.SPDBankType {
		bankTypeStr = "浦发银行"
	} else if bankType == enum.GuilinBankType {
		bankTypeStr = "桂林银行"
	} else if bankType == enum.PinganBankType {
		bankTypeStr = "平安银行"
	}

	virtualAccount, err := s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	virtualAccountBankName := ""
	switch bankType {
	case enum.SPDBankType:
		// 从配置中心获取虚账户
		virtualAccountParentNo := config.GetString("virtual.account.parent.no", "")
		if virtualAccountParentNo == "" {
			return nil, handler.HandleError(errors.New("未配置虚账户实账号"))
		}
		virtualAccountParentName := config.GetString("virtual.account.parent.name", "")
		if virtualAccountParentName == "" {
			return nil, handler.HandleError(errors.New("未配置虚账户实账号名称"))
		}

		if virtualAccount == nil {
			zap.L().Error(fmt.Sprintf("本地数据库->主账户: %s 未查询到虚账户: %s 的数据\n", virtualAccountParentNo, virtualAccount.VirtualAccountNo))
			return nil, nil
		}

		organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
			OrganizationId: organizationId,
			Type:           bankType,
		})
		if err != nil {
			return nil, handler.HandleError(errors.New(fmt.Sprintf("查询虚账户余额获取组织配置异常: %d,银行类型: %s", organizationId, bankTypeStr)))
		}
		//virtualAccountParentNo = "63160078801600000054"
		//virtualAccountParentName = "广西呦亿灵动科技有限公司"
		//virtualAccountNos := []string{"0154802949371"}
		virtualAccountNos := []string{virtualAccount.VirtualAccountNo}
		// sdk 查询浦发虚账号金额
		balanceDatas, err := s.spdBankSDK.QueryVirtualAccountBalance(ctx, virtualAccountParentNo, virtualAccountNos, -1, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
		if err != nil {
			return nil, handler.HandleError(errors.New(fmt.Sprintf("spdBankSDK.QueryVirtualAccountBalance error: %v\n", err.Error())))
		}
		if balanceDatas == nil || len(balanceDatas) == 0 {
			return nil, handler.HandleError(errors.New(fmt.Sprintf("浦发接口->主账户: %s 未查询到虚账户: %s 的数据\n", virtualAccountParentNo, virtualAccount.VirtualAccountNo)))
		}
		// 默认给上数据库的金额
		virtualAccountBalance := virtualAccount.VirtualAccountBalance
		// 循环遍历金额
		for _, v := range balanceDatas {
			if v.VirtualAcctNo == virtualAccount.VirtualAccountNo {
				virtualAccountBalance = v.AccountBalance
				break
			}
		}
		// 更新数据库
		virtualAccount.VirtualAccountBalance = virtualAccountBalance
		resultVirtualAccountBalance = virtualAccountBalance
		err = s.baseClient.EditOrganizationBankVirtualAccount(ctx, virtualAccount)
		if err != nil {
			zap.L().Error(fmt.Sprintf("母账号:%s 更新虚账号: %s 金额失败\n", virtualAccountParentNo, virtualAccount.VirtualAccountNo))
			return nil, handler.HandleError(errors.New(fmt.Sprintf("浦发接口->主账户: %s 未查询到虚账户: %s 的数据\n", virtualAccountParentNo, virtualAccount.VirtualAccountNo)))
		}
		break
	case enum.PinganBankType:
		virtualAccountBankName = config.GetString("bankConfig.pingAn.intelligence.accountBankName", "")

		seqNo, err := util.SonyflakeID()
		if err != nil {
			return nil, handler.HandleError(err)
		}
		var createAccountRequest = sdkStru.PinganQueryVirtualAccountBalanceRequest{
			MrchCode:        config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
			CnsmrSeqNo:      seqNo,
			MainAccount:     config.GetString(bankEnum.PinganIntelligenceAccountNo, ""),
			ReqSubAccountNo: virtualAccount.VirtualAccountNo,
		}
		// sdk 查询Pingan虚账号金额
		balanceData, err := s.pinganBankSDK.QueryVirtualAccount(ctx, createAccountRequest)
		//更新金额
		acctBalance, err := strconv.ParseFloat(balanceData.SubAccBalance, 64)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		virtualAccount.VirtualAccountBalance = acctBalance
		resultVirtualAccountBalance = acctBalance
		//更新金额
		err = s.baseClient.EditOrganizationBankVirtualAccount(ctx, virtualAccount)
		if err != nil {
			zap.L().Error(fmt.Sprintf("更新虚账号: %s 金额失败\n", virtualAccount.VirtualAccountNo))
			return nil, handler.HandleError(err)
		}
		break
	default:
		break
	}

	virtualAccount, err = s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
		OrganizationId: organizationId,
		Type:           bankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}

	return &api.VirtualAccountBalanceData{
		VirtualAccountNo:       virtualAccount.VirtualAccountNo,
		VirtualBalance:         resultVirtualAccountBalance,
		UpdateTime:             virtualAccount.UpdatedAt,
		BankType:               bankTypeStr,
		VirtualAccountName:     config.GetString(bankEnum.PinganIntelligenceAccountName, ""),
		VirtualAccountBankName: virtualAccountBankName,
	}, nil
}

// SyncVirtualAccountBalance
//
//	@Description: 同步浦发虚账户余额
//	@receiver s
//	@param ctx
//	@param organizationId
//	@return err
func (s *bankService) SyncVirtualAccountBalance(ctx context.Context) (err error) {
	// 从配置中心获取虚账户
	virtualAccountParentNo := config.GetString("virtual.account.parent.no", "")
	if virtualAccountParentNo == "" {
		return handler.HandleError(errors.New("未配置虚账户实账号"))
	}
	virtualAccountParentName := config.GetString("virtual.account.parent.name", "")
	if virtualAccountParentName == "" {
		return handler.HandleError(errors.New("未配置虚账户实账号名称"))
	}

	// 查询所有数据库浦发母账号下的虚账号列表
	virtualAccountDatas, err := s.baseClient.ListOrganizationBankVirtualAccountData(ctx, &baseApi.ListOrganizationBankVirtualAccountRequest{
		OrganizationId: 0,
		Type:           enum.SPDBankType,
		AccountNo:      virtualAccountParentNo,
	})
	if err != nil {
		return handler.HandleError(errors.New(fmt.Sprintf("s.baseClient.ListOrganizationBankVirtualAccountData error: %v\n", err.Error())))
	}
	if virtualAccountDatas == nil || len(virtualAccountDatas) == 0 {
		return handler.HandleError(errors.New(fmt.Sprintf("母账号: %s s.baseClient.ListOrganizationBankVirtualAccountData 为空\n", virtualAccountParentNo)))
	}
	organizationId := virtualAccountDatas[0].OrganizationId
	organizationBankConfig, err := s.baseClient.GetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: organizationId,
		Type:           enum.SPDBankType,
	})
	virtualAccountNo := []string{}
	for _, v := range virtualAccountDatas {
		virtualAccountNo = append(virtualAccountNo, v.VirtualAccountNo)
	}
	// sdk 查询浦发虚账号金额
	balanceDatas, err := s.spdBankSDK.QueryVirtualAccountBalance(ctx, virtualAccountParentNo, virtualAccountNo, -1, organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return handler.HandleError(errors.New(fmt.Sprintf("QueryVirtualAccountBalance error: %v\n", err.Error())))
	}
	if balanceDatas == nil || len(balanceDatas) == 0 {
		return handler.HandleError(errors.New(fmt.Sprintf("母账号:%s  s.spdBankSDK.QueryVirtualAccountBalance 为空\n", virtualAccountParentNo)))
	}
	// 虚账号金额 map
	balanceMap := make(map[string]float64, len(balanceDatas))
	for _, v := range balanceDatas {
		balanceMap[v.VirtualAcctNo] = v.AccountBalance
	}

	// 循环更新金额
	for _, v := range virtualAccountDatas {
		v.VirtualAccountBalance = balanceMap[v.VirtualAccountNo]
		err := s.baseClient.EditOrganizationBankVirtualAccount(ctx, v)
		if err != nil {
			zap.L().Error(fmt.Sprintf("母账号:%s 更新虚账号: %s 金额失败\n", virtualAccountParentNo, v.VirtualAccountNo))
			continue
		}
	}
	return nil
}
func (s *bankService) PinganSyncVirtualAccountBalance(ctx context.Context) (err error) {
	// 查询所有平安账户
	bankAccounts, err := s.baseClient.ListOrganizationBankAccount(ctx, &baseApi.ListOrganizationBankAccountRequest{
		Type: enum.PinganBankType,
	})
	if err != nil {
		return handler.HandleError(err)
	}
	if bankAccounts == nil || len(bankAccounts) <= 0 {
		return nil
	}
	for _, v := range bankAccounts {
		accountNo := v.Account
		organizationId := v.OrganizationId

		// 查询所有数据库平安母账号下的虚账号列表
		virtualAccountDatas, err := s.baseClient.ListOrganizationBankVirtualAccountData(ctx, &baseApi.ListOrganizationBankVirtualAccountRequest{
			OrganizationId: organizationId,
			Type:           enum.PinganBankType,
			AccountNo:      accountNo,
		})
		if err != nil {
			zap.L().Error(fmt.Sprintf("s.baseClient.ListOrganizationBankVirtualAccountData error: %v\n", err.Error()))
			continue
		}
		if virtualAccountDatas == nil || len(virtualAccountDatas) == 0 {
			zap.L().Info(fmt.Sprintf("母账号: %s s.baseClient.ListOrganizationBankVirtualAccountData 为空\n", accountNo))
			continue
		}

		for _, virtualAccount := range virtualAccountDatas {
			seqNo, err := util.SonyflakeID()
			if err != nil {
				return handler.HandleError(err)
			}
			var createAccountRequest = sdkStru.PinganQueryVirtualAccountBalanceRequest{
				MrchCode:        config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
				CnsmrSeqNo:      seqNo,
				MainAccount:     config.GetString(bankEnum.PinganIntelligenceAccountNo, ""),
				ReqSubAccountNo: virtualAccount.VirtualAccountNo,
			}
			// sdk 查询Pingan虚账号金额
			balanceData, err := s.pinganBankSDK.QueryVirtualAccount(ctx, createAccountRequest)
			if err != nil {
				zap.L().Error(fmt.Sprintf("QueryVirtualAccountBalance error: %v\n", err.Error()))
				continue
			}
			//更新金额
			acctBalance, err := strconv.ParseFloat(balanceData.SubAccBalance, 64)
			if err != nil {
				return handler.HandleError(err)
			}
			virtualAccount.VirtualAccountBalance = acctBalance
			//更新金额
			err = s.baseClient.EditOrganizationBankVirtualAccount(ctx, virtualAccount)
			if err != nil {
				zap.L().Error(fmt.Sprintf("母账号:%s 更新虚账号: %s 金额失败\n", accountNo, virtualAccount.VirtualAccountNo))
				continue
			}
		}
	}
	return nil
}

func (s *bankService) VirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error) {
	if req.PayAccountType == "2" {
		return s.PinganBankVirtualAccountTranscation(ctx, organizationId, req)
	} else if req.PayAccountType == "1" {
		return s.SPDBankVirtualAccountTranscation(ctx, organizationId, req)
	} else if req.PayAccountType == "21" {
		return s.PinganBankVirtualSubAcctBalanceAdjust(ctx, organizationId, req)
	}
	return nil, nil
}

// SPDBankVirtualAccountTranscation
//
//	@Description: 浦发银行虚账户转账接口
//	@receiver s
//	@param ctx
//	@param req
//	@return error
func (s *bankService) SPDBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error) {
	currentUserName, err := util.GetMetaInfoCurrentUserName(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	virtualAccount, err := s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
		OrganizationId: organizationId,
		Type:           enum.SPDBankType,
		//AccountNo:      bankAccount.Account,
		//VirtualAccountNo: virtualAccountNo,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("未查询到组织下: %d的虚账号\n", organizationId))
		return nil, handler.HandleError(err)
	}
	if req.RecAccount == "" {
		return nil, errors.New("RecAccount must not empty")
	}
	if req.RecAccountName == "" {
		return nil, errors.New("RecAccountName must not empty")
	}
	if req.PayAmount <= 0 {
		return nil, errors.New("PayAmount must greater than 0")
	}
	if virtualAccount.VirtualAccountNo == "" {
		return nil, errors.New("VirtualAccountNo must not empty")
	}
	if virtualAccount.VirtualAccountName == "" {
		return nil, errors.New("VirtualAccountName must not empty")
	}
	/*organizationBankAccount, err := s.baseClient.SimpleGetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		OrganizationId: organizationId,
		Type:           enum.SPDBankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}*/
	// 从配置中心获取虚账户
	virtualAccountParentNo := config.GetString("virtual.account.parent.no", "")
	if virtualAccountParentNo == "" {
		return nil, handler.HandleError(errors.New("未配置虚账户实账号"))
	}
	virtualAccountParentName := config.GetString("virtual.account.parent.name", "")
	if virtualAccountParentName == "" {
		return nil, handler.HandleError(errors.New("未配置虚账户实账号名称"))
	}

	organizationBankConfig, err := s.baseClient.SimpleGetOrganizationBankConfig(ctx, &baseApi.OrganizationBankConfigData{
		OrganizationId: organizationId,
		Type:           enum.SPDBankType,
	})
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 0 行内 1 行外
	ownItBankFlag := "0"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		ownItBankFlag = "1"
	}
	// 同城异地标志
	remitLocation := "0"
	// 收款行名称 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankName := ""
	// 收款行地址 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankAddress := ""
	if ownItBankFlag == "1" {
		//if req.RecAccountOpenBankFilling == "" {
		//	return transferReceiptId, errors.New("payeeBankAddress must not empty")
		//}
		if req.RecAccountOpenBank == "" {
			return nil, errors.New("payeeBankName must not empty")
		}
		// 由于分行同城已经取消，跨行支付时，同城异地标志建议固定送1异地
		remitLocation = "1"
		payeeBankAddress = req.RecAccountOpenBank // 先改成 RecAccountOpenBank, 之前是 RecAccountOpenBankFilling
		payeeBankName = req.RecAccountOpenBank
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	serialNo, _ := util.SonyflakeID()
	note := req.PayRem // 给灵活用工的附言, 不用拼接字符串
	request := sdkStru.SPDBankVirtualAccountTransferRequest{
		ElectronNumber:   serialNo,
		AcctNo:           virtualAccountParentNo,
		PayerVirAcctNo:   virtualAccount.VirtualAccountNo, // 根据组织去查虚账号
		PayerName:        virtualAccountParentName,        // 根据组织去查虚账号的实账户户名
		PayeeAcctNo:      req.RecAccount,
		PayeeAcctName:    req.RecAccountName,
		PayeeBankName:    payeeBankName,
		PayeeBankAddress: payeeBankAddress,
		TransAmount:      req.PayAmount,
		OwnItBankFlag:    ownItBankFlag,
		RemitLocation:    remitLocation,
		Note:             note,
	}
	bankTransferResponse, err := s.spdBankSDK.BankVirtualAccountTransfer(ctx, serialNo, request,
		organizationBankConfig.Host, organizationBankConfig.SignHost, organizationBankConfig.BankCustomerId, organizationBankConfig.BankUserId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 处理订单状态
	orderState := bankTransferResponse.BackhostStatus
	if bankTransferResponse.BackhostStatus != "" {
		orderState = enum.GetOrderState(bankTransferResponse.BackhostStatus, orderState)
		bankTransferResponse.ReturnCode = bankTransferResponse.BackhostStatus
		bankTransferResponse.ReturnMsg = enum.GetOrderMsg(bankTransferResponse.BackhostStatus, bankTransferResponse.ReturnMsg)
	}

	payRem := fmt.Sprintf("%s[%s]", req.PayRem, serialNo)
	transferReceiptId := int64(0)
	if transferReceiptId, err = s.AddBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		OrganizationId: organizationId,
		//ProcessInstanceId:  processInstanceDBData.Id,
		OriginatorUserName:        currentUserName,
		SerialNo:                  serialNo,
		RecAccount:                req.RecAccount,
		RecAccountName:            req.RecAccountName,
		RecAccountOpenBankFilling: payeeBankAddress,
		PayAmount:                 req.PayAmount,
		CurrencyType:              enum.CurrencyTypeCNY,
		PayRem:                    payRem,
		//PubPriFlag:                pubPriFlag,
		//RmtType:                   rmtType,
		OrderState: orderState,
		//ProcessBusinessId:         processInstanceDBData.BusinessId,
		//ProcessStatus:             enum.ProcessInstanceTotalStatusRunning,
		Title:              fmt.Sprintf("%s%s", currentUserName, "发起的付款单据"),
		PayAccount:         virtualAccount.VirtualAccountNo,
		PayAccountName:     virtualAccount.VirtualAccountName,
		ChargeFee:          0.00,
		RetCode:            bankTransferResponse.ReturnCode,
		RetMessage:         bankTransferResponse.ReturnMsg,
		OrderFlowNo:        bankTransferResponse.JnlSeqNo,
		ProcessComment:     req.ProcessComment,
		CommentUserName:    req.CommentUserName,
		RecBankType:        ownItBankFlag,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.SPDBankType,
	}); err != nil {
		return nil, handler.HandleError(err)
	}
	return &api.BankVirtualAccountTranscationResponse{
		TransferReceiptId: transferReceiptId,
		SerialNo:          serialNo,
		AcceptNo:          bankTransferResponse.JnlSeqNo,
		Status:            orderState,
		Msg:               bankTransferResponse.ReturnMsg,
	}, nil
}
func (s *bankService) PinganBankVirtualAccountTranscation(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error) {

	currentUserName, err := util.GetMetaInfoCurrentUserName(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	virtualAccount, err := s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
		OrganizationId: organizationId,
		Type:           enum.PinganBankType,
		//AccountNo:      bankAccount.Account,
		//VirtualAccountNo: virtualAccountNo,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("未查询到组织下: %d的虚账号\n", organizationId))
		return nil, handler.HandleError(err)
	}
	if req.RecAccount == "" {
		return nil, errors.New("RecAccount must not empty")
	}
	if req.RecAccountName == "" {
		return nil, errors.New("RecAccountName must not empty")
	}
	if req.PayAmount <= 0 {
		return nil, errors.New("PayAmount must greater than 0")
	}
	if virtualAccount.VirtualAccountNo == "" {
		return nil, errors.New("VirtualAccountNo must not empty")
	}
	if virtualAccount.VirtualAccountName == "" {
		return nil, errors.New("VirtualAccountName must not empty")
	}

	// 1 行内 0 行外
	ownItBankFlag := "1"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		ownItBankFlag = "0"
	}
	// 同城异地标志
	remitLocation := "1"
	// 收款行名称 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankName := ""
	// 收款行地址 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankAddress := ""
	if ownItBankFlag == "0" {
		//if req.RecAccountOpenBankFilling == "" {
		//	return transferReceiptId, errors.New("payeeBankAddress must not empty")
		//}
		if req.RecAccountOpenBank == "" {
			return nil, errors.New("payeeBankName must not empty")
		}
		// 1”—同城   “2”—异地；若无法区分，可默认送1-同城。
		remitLocation = "1"
		payeeBankAddress = req.RecAccountOpenBank // 先改成 RecAccountOpenBank, 之前是 RecAccountOpenBankFilling
		payeeBankName = req.RecAccountOpenBank
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	serialNo, _ := util.SonyflakeID()
	serialNo2, _ := util.SonyflakeID()
	serialNo = bankEnum.PinganFlexPrefix + serialNo
	useE := req.PayRem // 给灵活用工的附言, 不用拼接字符串
	request := sdkStru.PingAnBankTransferRequest{
		MrchCode:           config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
		CnsmrSeqNo:         serialNo2,
		ThirdVoucher:       serialNo,
		CcyCode:            "RMB",
		OutAcctNo:          virtualAccount.VirtualAccountNo,
		OutAcctName:        config.GetString(bankEnum.PinganIntelligenceAccountName, ""),
		InAcctNo:           req.RecAccount,
		InAcctName:         req.RecAccountName,
		InAcctBankName:     payeeBankName,
		InAcctProvinceCode: payeeBankAddress,
		InAcctBankNode:     req.UnionBankNo,
		TranAmount:         strconv.FormatFloat(req.PayAmount, 'f', -1, 64),
		UseEx:              useE,
		UnionFlag:          ownItBankFlag,
		AddrFlag:           remitLocation,
	}
	bankTransferResponse, err := s.pinganBankSDK.BankVirtualTransfer(ctx, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 处理订单状态
	orderState := bankTransferResponse.Stt
	if bankTransferResponse.Stt != "" {
		orderState = enum.GetOrderState(bankTransferResponse.Stt, orderState)
	}

	payRem := fmt.Sprintf("%s[%s]", req.PayRem, serialNo)
	transferReceiptId := int64(0)
	if transferReceiptId, err = s.AddBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		OrganizationId: organizationId,
		//ProcessInstanceId:  processInstanceDBData.Id,
		OriginatorUserName:        currentUserName,
		SerialNo:                  serialNo,
		RecAccount:                req.RecAccount,
		RecAccountName:            req.RecAccountName,
		RecAccountOpenBankFilling: payeeBankAddress,
		PayAmount:                 req.PayAmount,
		CurrencyType:              enum.CurrencyTypeCNY,
		PayRem:                    payRem,
		//PubPriFlag:                pubPriFlag,
		//RmtType:                   rmtType,
		OrderState: orderState,
		//ProcessBusinessId:         processInstanceDBData.BusinessId,
		ProcessStatus:      enum.ProcessInstanceTotalStatusRunning,
		Title:              fmt.Sprintf("%s%s", currentUserName, "发起的付款单据"),
		PayAccount:         virtualAccount.VirtualAccountNo,
		PayAccountName:     virtualAccount.VirtualAccountName,
		ChargeFee:          0.00,
		RetCode:            bankTransferResponse.Stt,
		RetMessage:         "",
		OrderFlowNo:        bankTransferResponse.HostFlowNo,
		ProcessComment:     req.ProcessComment,
		CommentUserName:    req.CommentUserName,
		RecBankType:        ownItBankFlag,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.PinganBankType,
		DetailHostFlowNo:   bankTransferResponse.HostFlowNo,
	}); err != nil {
		return nil, handler.HandleError(err)
	}
	return &api.BankVirtualAccountTranscationResponse{
		TransferReceiptId: transferReceiptId,
		SerialNo:          serialNo,
		AcceptNo:          bankTransferResponse.HostFlowNo,
		Status:            orderState,
		Msg:               "",
	}, nil
}
func (s *bankService) PinganBankTransaction(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankAccountTranscationResponse, error) {

	currentUserName, err := util.GetMetaInfoCurrentUserName(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	account, err := s.baseClient.SimpleGetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		OrganizationId: organizationId,
		Type:           enum.PinganBankType,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("未查询到组织下: %d的账号\n", organizationId))
		return nil, handler.HandleError(err)
	}
	if req.RecAccount == "" {
		return nil, errors.New("RecAccount must not empty")
	}
	if req.RecAccountName == "" {
		return nil, errors.New("RecAccountName must not empty")
	}
	if req.PayAmount <= 0 {
		return nil, errors.New("PayAmount must greater than 0")
	}
	if account.Account == "" {
		return nil, errors.New("account must not empty")
	}
	if account.AccountName == "" {
		return nil, errors.New("outAccountName must not empty")
	}

	// 1 行内 0 行外
	ownItBankFlag := "1"
	if req.PayAccountOpenBank != req.RecAccountOpenBank {
		ownItBankFlag = "0"
	}
	// 同城异地标志
	remitLocation := "1"
	// 收款行名称 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankName := ""
	// 收款行地址 (当 ownItBankFlag =1即跨行转帐时必须输入)
	payeeBankAddress := ""
	if ownItBankFlag == "0" {
		//if req.RecAccountOpenBankFilling == "" {
		//	return transferReceiptId, errors.New("payeeBankAddress must not empty")
		//}
		if req.RecAccountOpenBank == "" {
			return nil, errors.New("payeeBankName must not empty")
		}
		// 1”—同城   “2”—异地；若无法区分，可默认送1-同城。
		remitLocation = "1"
		payeeBankAddress = req.RecAccountOpenBank // 先改成 RecAccountOpenBank, 之前是 RecAccountOpenBankFilling
		payeeBankName = req.RecAccountOpenBank
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	serialNo, _ := util.SonyflakeID()
	serialNo2, _ := util.SonyflakeID()
	serialNo = bankEnum.PinganFlexPrefix + serialNo
	useE := req.PayRem // 给灵活用工的附言, 不用拼接字符串
	request := sdkStru.PingAnBankTransferRequest{
		MrchCode:           config.GetString(bankEnum.PinganPlatformAccount, ""),
		CnsmrSeqNo:         serialNo2,
		ThirdVoucher:       serialNo,
		CcyCode:            "RMB",
		OutAcctNo:          account.Account,
		OutAcctName:        account.AccountName,
		InAcctNo:           req.RecAccount,
		InAcctName:         req.RecAccountName,
		InAcctBankName:     payeeBankName,
		InAcctProvinceCode: payeeBankAddress,
		InAcctBankNode:     req.UnionBankNo,
		TranAmount:         strconv.FormatFloat(req.PayAmount, 'f', -1, 64),
		UseEx:              useE,
		UnionFlag:          ownItBankFlag,
		AddrFlag:           remitLocation,
	}

	//todo
	bankTransferResponse := &sdkStru.PingAnBankTransferResponse{
		OutAcctName: request.OutAcctName,
		HostFlowNo:  "orderFlowNo2",
		Stt:         "20",
	}
	//bankTransferResponse, err := s.pinganBankSDK.BankTransfer(ctx, request, account.ZuId)
	//if err != nil {
	//	return nil, handler.HandleError(err)
	//}
	// 处理订单状态
	orderState := bankTransferResponse.Stt
	if bankTransferResponse.Stt != "" {
		orderState = enum.GetOrderState(bankTransferResponse.Stt, orderState)
	}

	payRem := fmt.Sprintf("%s[%s]", req.PayRem, serialNo)
	transferReceiptId := int64(0)
	if transferReceiptId, err = s.AddBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		OrganizationId: organizationId,
		//ProcessInstanceId:  processInstanceDBData.Id,
		OriginatorUserName:        currentUserName,
		SerialNo:                  serialNo,
		RecAccount:                req.RecAccount,
		RecAccountName:            req.RecAccountName,
		RecAccountOpenBankFilling: payeeBankAddress,
		PayAmount:                 req.PayAmount,
		CurrencyType:              enum.CurrencyTypeCNY,
		PayRem:                    payRem,
		//PubPriFlag:                pubPriFlag,
		//RmtType:                   rmtType,
		OrderState: orderState,
		//ProcessBusinessId:         processInstanceDBData.BusinessId,
		ProcessStatus:      enum.ProcessInstanceTotalStatusRunning,
		Title:              fmt.Sprintf("%s%s", currentUserName, "发起的付款单据"),
		PayAccount:         account.Account,
		PayAccountName:     account.AccountName,
		ChargeFee:          0.00,
		RetCode:            bankTransferResponse.Stt,
		RetMessage:         "",
		OrderFlowNo:        bankTransferResponse.HostFlowNo,
		ProcessComment:     req.ProcessComment,
		CommentUserName:    req.CommentUserName,
		RecBankType:        ownItBankFlag,
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.PinganBankType,
		DetailHostFlowNo:   bankTransferResponse.HostFlowNo,
	}); err != nil {
		return nil, handler.HandleError(err)
	}
	return &api.BankAccountTranscationResponse{
		TransferReceiptId: transferReceiptId,
		SerialNo:          serialNo,
		AcceptNo:          bankTransferResponse.HostFlowNo,
		Status:            orderState,
		Msg:               "",
	}, nil
}

func (s *bankService) PinganBankAccountSignatureApply(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (*api.PinganUserAcctSignatureApplyResponse, error) {
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id: req.Id,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("PinganBankAccountSignatureApply error: %v\n", err.Error()))
	}
	response, err := s.pinganBankSDK.UserAcctSignatureApply(ctx, bankAccount.Account, bankAccount.AccountName,
		config.GetString(bankEnum.PinganJsdkUrl, ""),
		config.GetString(bankEnum.PinganMrchCode, ""))
	if err != nil {
		zap.L().Error(fmt.Sprintf("PinganBankAccountSignatureApply error: %v\n", err.Error()))
		return nil, handler.HandleError(err)
	}
	return &api.PinganUserAcctSignatureApplyResponse{
		ZuID:      response.ZuID,
		OpFlag:    response.OpFlag,
		Stt:       response.Stt,
		AccountNo: response.AccountNo,
	}, nil
}
func (s *bankService) PinganBankAccountSignatureQuery(ctx context.Context, req *api.PinganBankAccountSignatureApplyRequest) (*api.PinganUserAcctSignatureApplyResponse, error) {
	bankAccount, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id: req.Id,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("PinganBankAccountSignatureApply error: %v\n", err.Error()))
	}

	response, err := s.pinganBankSDK.UserAcctSignatureQuery(ctx, bankAccount.Account, bankAccount.AccountName, config.GetString(bankEnum.PinganJsdkUrl, ""), config.GetString(bankEnum.PinganMrchCode, ""), bankAccount.ZuId)
	if err != nil {
		zap.L().Error(fmt.Sprintf("PinganBankAccountSignatureApply error: %v\n", err.Error()))
		return nil, handler.HandleError(err)
	}
	s.baseClient.EditOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
		Id:                   req.Id,
		SyncTime:             nil,
		SignatureApplyStatus: response.Stt,
		Remark:               response.Remark,
	})
	return &api.PinganUserAcctSignatureApplyResponse{
		ZuID:      response.ZuID,
		OpFlag:    response.OpFlag,
		Stt:       response.Stt,
		AccountNo: response.AccountNo,
	}, nil
}

func (s *bankService) PinganBankVirtualSubAcctBalanceAdjust(ctx context.Context, organizationId int64, req *api.BankTransferReceiptData) (*api.BankVirtualAccountTranscationResponse, error) {
	currentUserName, err := util.GetMetaInfoCurrentUserName(ctx)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	virtualAccount, err := s.baseClient.SimpleGetOrganizationBankVirtualAccount(ctx, &baseApi.OrganizationBankVirtualAccountData{
		OrganizationId: organizationId,
		Type:           enum.PinganBankType,
		//AccountNo:      bankAccount.Account,
		//VirtualAccountNo: virtualAccountNo,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("未查询到组织下: %d的虚账号\n", organizationId))
		return nil, handler.HandleError(err)
	}
	if req.RecAccount == "" {
		return nil, errors.New("RecAccount must not empty")
	}
	if req.RecAccountName == "" {
		return nil, errors.New("RecAccountName must not empty")
	}
	if req.PayAmount <= 0 {
		return nil, errors.New("PayAmount must greater than 0")
	}
	if virtualAccount.VirtualAccountNo == "" {
		return nil, errors.New("VirtualAccountNo must not empty")
	}
	if virtualAccount.VirtualAccountName == "" {
		return nil, errors.New("VirtualAccountName must not empty")
	}
	// 和桂林不同, serialNo一天内唯一, 重新生成
	serialNo, _ := util.SonyflakeID()
	serialNo2, _ := util.SonyflakeID()
	serialNo = bankEnum.PinganFlexSubPrefix + serialNo
	useE := req.PayRem // 给灵活用工的附言, 不用拼接字符串
	request := sdkStru.PinganSubAcctBalanceAdjustRuest{
		MrchCode:          config.GetString(bankEnum.PinganIntelligenceMrchCode, ""),
		CnsmrSeqNo:        serialNo2,
		ThirdVoucher:      serialNo,
		MainAccount:       config.GetString(bankEnum.PinganIntelligenceAccountNo, ""),
		MainAccountName:   config.GetString(bankEnum.PinganIntelligenceAccountName, ""),
		OutSubAccount:     virtualAccount.VirtualAccountNo,
		OutSubAccountName: virtualAccount.VirtualAccountName,
		TranAmount:        strconv.FormatFloat(req.PayAmount, 'f', -1, 64),
		InSubAccNo:        req.RecAccount,
		InSubAccName:      req.RecAccountName,
		UseEx:             useE,
	}

	bankTransferResponse, err := s.pinganBankSDK.SubAcctBalanceAdjust(ctx, request)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	// 处理订单状态
	orderState := "51"

	payRem := fmt.Sprintf("%s[%s]", req.PayRem, serialNo)
	transferReceiptId := int64(0)
	if transferReceiptId, err = s.AddBankTransferReceipt(ctx, &api.BankTransferReceiptData{
		OrganizationId: organizationId,
		//ProcessInstanceId:  processInstanceDBData.Id,
		OriginatorUserName: currentUserName,
		SerialNo:           serialNo,
		RecAccount:         req.RecAccount,
		RecAccountName:     req.RecAccountName,
		//RecAccountOpenBankFilling: payeeBankAddress,
		PayAmount:    req.PayAmount,
		CurrencyType: enum.CurrencyTypeCNY,
		PayRem:       payRem,
		//PubPriFlag:                pubPriFlag,
		//RmtType:                   rmtType,
		OrderState: orderState,
		//ProcessBusinessId:         processInstanceDBData.BusinessId,
		ProcessStatus:  enum.ProcessInstanceTotalStatusRunning,
		Title:          fmt.Sprintf("%s%s", currentUserName, "发起的付款单据"),
		PayAccount:     virtualAccount.VirtualAccountNo,
		PayAccountName: virtualAccount.VirtualAccountName,
		ChargeFee:      0.00,
		//RetCode:            bankTransferResponse.Stt,
		RetMessage:         "",
		OrderFlowNo:        bankTransferResponse.FrontFlowNo,
		ProcessComment:     req.ProcessComment,
		CommentUserName:    req.CommentUserName,
		RecBankType:        "1",
		PayAccountOpenBank: req.PayAccountOpenBank,
		RecAccountOpenBank: req.RecAccountOpenBank,
		UnionBankNo:        req.UnionBankNo,
		ClearBankNo:        req.ClearBankNo,
		PayAccountType:     enum.PinganBankType,
		//DetailHostFlowNo:   bankTransferResponse.HostFlowNo,
	}); err != nil {
		return nil, handler.HandleError(err)
	}
	return &api.BankVirtualAccountTranscationResponse{
		TransferReceiptId: transferReceiptId,
		SerialNo:          serialNo,
		AcceptNo:          bankTransferResponse.FrontFlowNo,
		Status:            orderState,
		Msg:               "",
	}, nil
}

func (s *bankService) MinShengBankAccountSignatureApply(ctx context.Context, req *api.MinShengBankAccountSignatureRequest) (string, error) {
	if req.Id == 0 {
		return "", handler.HandleNewError("id不能为空")
	}
	account, err := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{Id: req.Id})
	if err != nil {
		return "", handler.HandleError(err)
	}
	if account == nil || account.Id == 0 {
		return "", handler.HandleNewError("账号不存在")
	}

	minShengEnterpriseIdCode := config.GetString(bankEnum.MinShengEnterpriseIdCode, "")

	//生成请求流水号
	msgId, _ := util.SonyflakeID()
	result, err := s.minShengBank.AuthRequest(ctx, minShengEnterpriseIdCode, account.Account, msgId)
	if err != nil {
		return "", handler.HandleError(err)
	}
	if result["return_code"] == "0000" {
		err = s.baseClient.EditOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
			Id:                   req.Id,
			SignatureApplyStatus: "1", // 受理成功
			ZuId:                 msgId,
			Remark:               result["return_code"].(string) + result["return_msg"].(string),
			//PaymentMode:          "1",
		})
		if err != nil {
			return "", handler.HandleError(err)
		}
	} else {
		return "", handler.HandleNewError("授权失败：" + result["return_msg"].(string))
	}
	return "", nil
}

func (s *bankService) MinShengBankAccountSignatureQuery(ctx context.Context, req *api.MinShengBankAccountSignatureRequest) (*api.MinShengBankAccountSignatureQueryResponse, error) {
	bankAccount, _ := s.baseClient.GetOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{Id: req.Id})
	if bankAccount == nil || bankAccount.Id == 0 {
		return nil, handler.HandleNewError("账户不存在")
	}
	if bankAccount.SignatureApplyStatus == "0" { // 授权成功
		return &api.MinShengBankAccountSignatureQueryResponse{
			Signatureapplystatus: bankAccount.SignatureApplyStatus,
			StartTime:            bankAccount.StartTime,
			EndTime:              bankAccount.EndTime,
		}, nil
	}
	//生成请求流水号
	msgId, _ := util.SonyflakeID()
	response, err := s.minShengBank.QueryAuthStatus(ctx, bankAccount.ZuId, msgId)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if response["return_code"].(string) == "0000" { // 请求成功
		responseBusi := response["response_busi"].(string)
		var busiMap map[string]string
		err = json.Unmarshal([]byte(responseBusi), &busiMap)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		if busiMap["status"] == "1" { // 授权成功
			err = s.baseClient.EditOrganizationBankAccount(ctx, &baseApi.OrganizationBankAccountData{
				Id:                   req.Id,
				SignatureApplyStatus: "0", //授权成功
				OpenId:               busiMap["open_id"],
				StartTime:            busiMap["start_time"],
				EndTime:              busiMap["end_time"],
			})
			if err != nil {
				return nil, handler.HandleError(err)
			}
		}
		return &api.MinShengBankAccountSignatureQueryResponse{
			Signatureapplystatus: busiMap["status"],
			StartTime:            busiMap["start_time"],
			EndTime:              busiMap["end_time"],
		}, nil
	} else {
		return nil, handler.HandleNewError("请求失败：" + response["return_msg"].(string))
	}
}

func (s *bankService) setRedisLog(ctx context.Context, key, log string) {
	//获取key的长度,超过范围就删掉一个
	len, _ := s.redisClient.LLen(ctx, key).Result()
	if len >= 500 {
		s.redisClient.RPop(ctx, key)
		s.redisClient.LPush(ctx, key, log)
	} else {
		s.redisClient.LPush(ctx, key, log)
	}
}
