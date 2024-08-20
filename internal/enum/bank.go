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
const IcbcHost = "bankConfig.icbc.host"
const IcbcCorpNo = "bankConfig.icbc.corpNo"       //一级合作方编号
const IcbcAccCompNo = "bankConfig.icbc.acccompno" //二级合作方编号
const IcbcAccountNo = "bankConfig.icbc.accountNo" //主账号
const IcbcPrivateKey = "bankConfig.icbc.privateKey"
const IcbcAccDetailURL = "/api/mybank/account/accountdetailservice/adsaccountdtlqry/V1"      //工行单账户流水提取url
const IcbcAdsAgreementGryURL = "/api/mybank/account/accountdetailservice/adsagreementqry/V1" //协议批量查询
const IcbcAdsReceiptAryURL = "/api/mybank/account/accountdetailservice/adsreceiptqry/V1"     //准实时回单查询

const IcbcTempFilePath = "tempFile/icbc"
const GetBankTransactionReceipt = "getBankTransactionReceipt:"

const IcbcSftPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAsI1N3bI2WhbjzFnmHhqN7Pdkeh8TICF9OPFjfUC0hoMdURlX2XsM
zuoakQAfVsf9KvFgwC4KW5f++bn47t28uqDyfBCq4puDVK0UFUQJ0/jptIjd3A5ghktUDc
Pyk7X0E9X791DoapKHOLJTSiZZgD+srl1bjpIAqJ9CsrfjS6fjoCExlpq0SsKbNFhchj0S
E8RgBjfM+lKu+DKeYBJIikvvsdMBMq6dOf2eRkvNrhvSpn8nANd9YKnTG25uDaTijkmOOh
RvScgmujLxf4sqvFazSxseDzhd2Q7kBGEX/gFugDKns69dMaM2PAEoFqQE4GqO3rHwI3YL
ao8KQM5r/wAAA9CPq3/4j6t/+AAAAAdzc2gtcnNhAAABAQCwjU3dsjZaFuPMWeYeGo3s92
R6HxMgIX048WN9QLSGgx1RGVfZewzO6hqRAB9Wx/0q8WDALgpbl/75ufju3by6oPJ8EKri
m4NUrRQVRAnT+Om0iN3cDmCGS1QNw/KTtfQT1fv3UOhqkoc4slNKJlmAP6yuXVuOkgCon0
Kyt+NLp+OgITGWmrRKwps0WFyGPRITxGAGN8z6Uq74Mp5gEkiKS++x0wEyrp05/Z5GS82u
G9KmfycA131gqdMbbm4NpOKOSY46FG9JyCa6MvF/iyq8VrNLGx4POF3ZDuQEYRf+AW6AMq
ezr10xozY8ASgWpATgao7esfAjdgtqjwpAzmv/AAAAAwEAAQAAAQA5AkfEcIlQadfA4r6F
tfliLThKnsIkO+wdeQSxKzWfwbKzv0U4up0WK03MyIdWFFnRhgPBypwZm2j/5mdValBIyz
PBj/g+GA0+SG0VuNSbl+KPIyrQpevRMX3AvCcWP0jDJvOnln6V+x6i1iJC7UM1QFpYK1kn
HkoMKPD2mJ5SjShxvzZjvRicz3FCo1jiq5HTKELHUndoVL9fPCKxos+FFvVDRdzt6eIfLA
IcHLqrs4iXbim7hxB+pp9UrYtrxQBVbVeAE/hzLh++td1uIMk9+vc0xNbxlRZqCneEKVzW
ulqenT+dJf2QCuew4a+ApuDLxcqeT967tpZKwuLiGsPxAAAAgQCiCZey5abMlkd6infc1S
HNtO1LXN2UxtolVAViK+Z388NVK7KRRnDTk7ONiOYozxnlGmSeDNMFhbl3tXjUm7xPffU5
2k9wzumvShmb16ISSW5lu7avaOCMY+bUluETZ04e7UeNEtXlPPrGymmMjaoXRPqvPHCYPX
Q1UnaIkuAcBgAAAIEA3SaO4k8Oe7bmpB+rkRTk05vBgrCJ6HFf7jz5nmUpaEepOFmrzhNr
20mR3JD09Odhfd++wY7hdG0ayGhU73bHC0lvpCw70l1AM6qQkzeHZi/6t/qZe66UBoRWXH
AvpLizh0s6wyAn0Jl5I065LzpU46MXwPJsZvBXnsi636Km3B0AAACBAMxfl0rEkOLh6M+9
ZhSbK9E5AhdVX7ZKPBbjMpbbzUGekXYdGvqnwVd+4B2lMQWH0fcoLEO/CNVrxEsvl8mfyg
ZF99TXutj8CFyWKMpI6EQ2XsYdSVZtLJqBQHuYWA1jK3uHz63ePw3v0I7mrwLTFiJ9YgST
niWlJ6Runo5iFZXLAAAAFDAwMTAzMTg3NkAwMDEwMzE4NzZUAQIDBAUG
-----END OPENSSH PRIVATE KEY-----`

const BankReceiptSyncLogKey = "bankReceiptDownloadSyncLog"
const BankErrorLogKey = "bankErrorLog"
