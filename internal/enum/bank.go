package enum

const PinganPlatformAccount = "bankConfig.pingAn.accountKeeper.accountNo"
const PinganMrchCode = "bankConfig.pingAn.accountKeeper.mrchCode"
const PinganJsdkUrl = "bankConfig.pingAn.jsdkUrl"
const PinganJsdkFileUrl = "bankConfig.pingAn.jsdkFileUrl"
const PinganIntelligenceMrchCode = "bankConfig.pingAn.intelligence.mrchCode"
const PinganIntelligenceAccountNo = "bankConfig.pingAn.intelligence.accountNo"
const PinganIntelligenceAccountName = "bankConfig.pingAn.intelligence.accountName"

const PinganFlexPrefix = "LH-"     //平安子账户转账前缀
const PinganFlexSubPrefix = "LH-S" //平安

const IcbcAppId = "bankConfig.icbc.appId"
const IcbcPrivateKey = "bankConfig.icbc.privateKey"
const IcbcAccDetailURL = "/mybank/account/accountdetailservice/adsaccountdtlqry/V1"          //工行单账户流水提取url
const IcbcBatchAccDetailURL = "/mybank/account/accountdetailservice/adsbatchaccdtlqry/V1"    //工行批量账户流水提取url
const IcbcAdsAgrConfirmSynURL = "/mybank/account/accountdetailservice/adsagrconfirmsyn/V1"   //协议待确认信息同步
const IcbcAdsAgreementGryURL = "/mybank/account/accountdetailservice/adsagreementqry/V1"     //协议批量查询
const IcbcAdsAgrSigUiURL = "/mybank/account/accountdetailservice/adsagrsignui/V1"            //协议签订页面API
const IcbcAdsReceiptAryURL = "/mybank/account/accountdetailservice/adsreceiptqry/V1"         //准实时回单查询
const IcbcAdsPartNerSynURL = "/mybank/account/accountdetailservice/adspartnersyn/V1"         //代账公司信息同步
const IcbcAdsPartNerGryURL = "/mybank/account/accountdetailservice/adspartnerqry/V1"         //代账公司信息查询
const IcbcAdsApplyResultGryURL = "/mybank/account/accountdetailservice/adsapplyresultqry/V1" //批量申请结果查询
const IcbcAdsAgreementPushURL = "/mybank/account/accountdetailservice/adsagreementpush/V1"   //协议变化通知（推送接口）
