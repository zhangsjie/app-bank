namespace go api

// 最初开始做的时候 BankTransferReceiptData 就是以桂林银行的接口为标准开发的, 所有后续银行的接口字段要统一到这个结构体上
struct BankTransferReceiptData {
    1: i64 id
	2: binary createdAt
	3: binary updatedAt
	4: i64 organizationId
    5: i64 processInstanceId
    6: string originatorUserName
    7: string serialNo
    8: string payAccount
    9: string payAccountName
    10: string recAccount
    11: string recAccountName
    12: double payAmount
    13: string currencyType
    14: string payRem
    15: string pubPriFlag
    16: string recBankType
    17: string recAccountOpenBank
    18: string unionBankNo
    19: string clearBankNo
    20: string rmtType
    21: string transferFlag
    22: double chargeFee
    23: string orderState
    24: string retCode
    25: string retMessage
    26: string orderFlowNo
    27: string processBusinessId
    28: string processComment
    29: string DetailHostFlowNo
    30: i64 CommentUserId
    31: string commentUserName
    32: string processStatus
    33: string recAccountOpenBankFilling
    34: string payAccountOpenBankFilling
    35: string payAccountOpenBank
    36: string title
    37: string payAccountType
}

struct ListBankTransferReceiptRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
    4: list<i64> processInstanceIds
    5: string businessId
    6: string serialNo
    7: string payAccount
    8: string recAccount
    9: string originatorUser
    10: list<double> payAmount
    11: list<string> createTimeArray
    13: string commentUser
    14: string totalStatus
    15: string orderState
    16: double payAmountMin
    17: double payAmountMax
    18: string title
}

struct ListBankTransferReceiptResponse {
    1: list<BankTransferReceiptData> data
    2: i64 count
}

struct BankTransactionDetailData {
    1: i64 id
	2: binary createdAt
	3: binary updatedAt
	4: i64 organizationId
	5: i64 merchantAccountId
	6: string merchantAccount
	7: string merchantAccountName
	8: string type
    9: string cashFlag
    10: double payAmount
    11: double recAmount
    12: string bsnType
    13: string transferDate
    14: string transferTime
    15: string tranChannel
    16: string currencyType
    17: double balance
    18: string orderFlowNo
    19: string hostFlowNo
    20: string vouchersType
    21: string vouchersNo
    22: string summaryNo
    23: string summary
    24: string acctNo
    25: string accountName
    26: string accountOpenNode
    27: string electronicReceiptFile
    28: string processBusinessId
    29: string processTotalStatus
    30: string originatorUser
    31: string operationUser
    32: string operationComment
    33: i64 processInstanceId
    34: string payAccountType
    35: string extField2
    36: string extField3
}

struct ListBankTransactionDetailRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
    4: string type
    5: string businessId
    6: string serialNo
    7: string merchantAccountName
    8: string accountName
    9: string originatorUser
    10: string operationUser
    11: list<double> payAmount
    12: list<string> createTimeArray
    13: string totalStatus
    14: string businessType
    15: double payAmountMin
    16: double payAmountMax
    17: double recAmountMin
    18: double recAmountMax
    19: list<string> transferTimeArray
    20: string payAccountType
    21: string extField2
    22: string extField3
}

struct ListBankTransactionDetailResponse {
    1: list<BankTransactionDetailData> data
    2: i64 count
}

struct BankTransactionRecDetailData {
    1: i64 id
    2: i64 organizationId
	3: string externalId
    4: string processBusinessId
    5: string operationUserId
    6: string operationUserName
    7: string operationComment
    8: string processTotalStatus
    9: i64 ProcessInstanceId
}

struct BankTransactionDetailProcessInstanceData {
    1: i64 bankTransactionDetailId
    2: string externalId
    3: i64 organizationId
}

struct BankCodeData {
    1: i64 id
    2: string bankName
    3: string bankAliasName
    4: string bankCode
    5: string unionBankNo
    6: string clearBankNo
}

struct QueryBankCardInfoResponse {
    1: string cardType
    2: string bank
    3: string key
    4: string messages
    5: bool   validated
    6: string stat
}

struct ListBankCodeRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
    4: string bankName
    5: string bankAliasName
}

struct ListBankCodeResponse {
    1: list<BankCodeData> data
    2: i64 count
}

struct AddBankCodeRequest {
   1: string bankName
   2: string bankAliasName
   3: string unionBankNo
   4: string clearBankNo
}

struct DashboardData {
    1: ChartData dayFlowData
    2: ChartData weekFlowData
    3: ChartData monthFlowData
    4: ChartData weekBalanceData
    5: ChartData monthBalanceData
}

struct ChartData {
    1: list<string> labels
    2: list<ChartDataSet> datasets
}

struct ChartDataSet {
    1: string label
    2: list<double> data
    3: string borderColor
    4: list<string> backgroundColor
    5: list<string> hoverBackgroundColor
}

struct MonthChartDataRequest {
   1: string year
   2: string month
   3: i64 organizationId
}

struct ListBusinessPayrollRequest {
   1: required i32 pageNum
   2: required i32 pageSize
   3: string sort
   4: string pay_account_no
   5: string pay_account_name
   6: string month
   7: i64    status
   8: i64    createUser
   9: list<string> createdTime
}

struct ListBusinessPayrollResponse {
   1: list<BusinessPayrollListVo> data
   2: i64 count
}

struct BusinessPayrollListVo {
    1: i64    id
    2: string name
    3: string payAccountNo
    4: string month
    5: string remark
    6: i64    totalCount
    7: i64    status
    8: string createdUser
    9: string createdAt
    10: string payAccountName
    11: string state
    12: string msg
    13: string submitTime
    14: double totalAmount
}

struct CreateBusinessPayrollRequest {
   1: string payAccountNo
   2: string payAccountName
   3: string month
   4: string remark
   5: i64    count
   6: double totalMoney
   7: i64    createdUserId
   8: binary fileBytes
}

struct ListBusinessPayrollDetailRequest {
   1: required i32 pageNum
   2: required i32 pageSize
   3: string    sort
   4: i64       batchId
   5: string    num
   6: string    month
   7: string    recAccountName
   8: string    recAccountNo
}

struct ListBusinessPayrollDetailResponse {
   1: list<BusinessPayrollDetailListVo> data
   2: i64 count
}

struct BusinessPayrollDetailListVo {
     1: i64     id
     2: i64     batchId
     3: string  recAccountName
     4: string  recAccountNo
     5: double  amount
     6: string  month
     7: string  orderState
     8: string  errorMessage
     9: string  num
     10: string errorCode
     11: string remark
}

struct QueryAccountBalanceRequest {
     1: string  accountNo
     2: i64     organizationId
}

struct QueryAccountBalanceResponse {
     1: bool    success
     2: string  msg
     3: double  balance
}

struct SyncBusinessPayrollResultRequest {
    1: i64 id
}

struct SyncBusinessPayrollResultResponse {
    1: bool     success
    2: string   msg
}

struct ProofData {
    1: i64 id
	2: binary createdAt
	3: binary updatedAt
	4: i64 organizationId
    5: i64 financeAccountSetId
    6: string date
    7: i64 num
    8: double borrowAmount
    9: double loanAmount
    10: list<ProofItemData> items
}

struct ProofItemData {
    1: i64 id
	2: binary createdAt
	3: binary updatedAt
	4: i64 organizationId
    5: i64 proofId
    6: string summary
    7: string subjectFullCodePath
    8: string subjectName
    9: double borrowAmount
    10: double loanAmount
    11: i32 sort
}

struct ListProofRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
}

struct ListProofResponse {
    1: list<ProofData> data
    2: i64 count
}

struct AccountSetPeriodData {
    1: i64 id
	2: binary createdAt
	3: binary updatedAt
	4: i64 organizationId
    5: string year
    6: string month
    7: string status
}

struct ListAccountSetPeriodRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
    4: string status
}

struct ListAccountSetPeriodResponse {
    1: list<AccountSetPeriodData> data
    2: i64 count
}

struct CreateVirtualAccountRequest {
    1: i64 organizationId
    2: VirtualAccountData data
    3: string type
}
struct CreateVirtualAccountResponse {
     1: string  status
     2: string  msg
     3: VirtualAccountData data
}

struct VirtualAccountData {
    1: string virtualAccountName
    2: string virtualAccountNo
    3: double rate
}

struct VirtualAccountBalanceData {
    1: string virtualAccountNo
    2: double virtualBalance
    3: binary updateTime
    4: string bankType
    5: string virtualAccountName
    6: string virtualAccountBankName
}

struct BankVirtualAccountTranscationResponse {
    1: i64 transferReceiptId
    2: string serialNo
    3: string acceptNo
    4: string status
    5: string msg
}

struct PaymentReceiptData {
	1: i64 id
	2: i64 organizationId
	3: binary createdAt
	4: binary updatedAt
	5: i64 processInstanceId
	6: string code
	7: double payAmount
    8: string publicPrivateFlag
    9: string payAccount
    10: string payAccountName
    11: string payAccountBankName
    12: string payAccountType
    13: string receiveAccount
    14: string receiveAccountName
    15: string receiveAccountBankName
    16: string purpose
    17: string unionBankNo
    18: string clearBankNo
    19: string insideOutsideBankType
    20: double chargeFee
    21: string orderStatus
    22: string retCode
    23: string retMessage
    24: string orderFlowNo
    25: bool canWrite
    26: string type
    27: i64 processInstanceItemId
    28: string paymentModeType
    29: i64 applicantId
    30: string applicantName
    31: string fillingDt
    32: i64 departmentId
    33: string departmentName
	34: string attachments
	35: string electronicReceiptFile
    36: string busType
    37: string busOrderNo
    38: string refundSuccess
    39: string receiptOrderNo
    40: string remark
    41: ProcessAddTagItemVO processAddTagItemVO
}

struct ListPaymentReceiptRequest {
    1: required i32 pageNum
    2: required i32 pageSize
    3: string sort
    4:string code
    5:string receiveAccount
    6:string receiveAccountName
    7:string receiveAccountBankName
    8:string purpose
    9:string orderStatus
    10:string orderFlowNo
    11:string unionBankNo
    12:string clearBankNo
    13:string createTimeStart
    14:string createTimeEnd
    15:string type
    16: string processName
    17: string processCodes
    18: string processStatus
    19: string refundSuccess
    20: string processCurrentUserName
    21: i64 processCurrentUserId
    22: string receiptOrderNo
    23: string endTime
    24: string beginTime
    25: string paymentModeType
    26: string applicantName
}

struct ListPaymentReceiptResponse {
    1: list<PaymentReceiptData> data
    2: i64 count
}

struct pinganBankAccountSignatureApplyRequest {
    1: i64 id
	2: i64 organizationId
	3: string type
    4: string account
    5: string accountName
}

struct pinganUserAcctSignatureApplyResponse  {
	1: string  zuID
	2: string  opFlag
	3: string  stt
	4: string  accountNo
}

struct ProcessAddTagItemVO {
	1: i64 tagType
	2: i64 nodeType
	3: list<ProcessAddTagItemUserVO> itemUsers
	4: i64 processInstanceId
}
struct ProcessAddTagItemUserVO {
	1: i64 userId
	2: string nickName
}

service bank {
    ListBankTransferReceiptResponse listBankTransferReceipt(1: ListBankTransferReceiptRequest req)
    BankTransferReceiptData getBankTransferReceipt(1: BankTransferReceiptData req)
    i64 addBankTransferReceipt(1: BankTransferReceiptData req)
    void editBankTransferReceipt(1: BankTransferReceiptData req)
    void deleteBankTransferReceipt(1: i64 req)
    i64 countBankTransferReceipt(1: BankTransferReceiptData req)

    void confirmTransaction(1: BankTransferReceiptData req)
    void handleTransferReceiptResult(1: i64 id)

    ListBankTransactionDetailResponse listBankTransactionDetail(1: ListBankTransactionDetailRequest req)
    BankTransactionDetailData getBankTransactionDetail(1: BankTransactionDetailData req)
    void handleTransactionDetail(1: string beginDate, 2: string endDate, 3: i64 organizationId)
    void createTransactionDetailProcessInstance(1: i64 id)
    void EditBankTransactionDetailExtField(1: BankTransactionDetailData req)

    list<BankTransactionDetailProcessInstanceData> listBankTransactionDetailProcessInstance(1: i64 id)

    BankCodeData GetBankCodeInfo(1: string code)
    QueryBankCardInfoResponse QueryBankCardInfo(1: string cardNo)
    ListBankCodeResponse ListBankCode(1: ListBankCodeRequest req)
    BankCodeData GetBankCode(1: BankCodeData req)
    void AddBankCode(1: AddBankCodeRequest req)
    void EditBankCode(1: BankCodeData req)
    void DeleteBankCode(1: i64 id)

    void HandleSyncTransferReceipt(1: string beginDate, 2: string endDate, 3: i64 organizationId)

    void UpdateBankTransactionRecDetail(1: BankTransactionRecDetailData req)

    void syncTransferReceipt(1: i64 taskId, 2: binary param, 3: i64 organizationId)
    void syncTransactionDetail(1: i64 taskId, 2: binary param, 3: i64 organizationId)

    DashboardData dashboardData(1: i64 organizationId)
    ChartData getCashFlowMonthChartData(1: MonthChartDataRequest req)
    ChartData getBalanceMonthChartData(1: MonthChartDataRequest req)

	QueryAccountBalanceResponse QueryAccountBalance(1: QueryAccountBalanceRequest req)
	void ImportBankBusinessPayrollData(1: i64 taskId, 2: binary param, 3: i64 organizationId)

	ListBusinessPayrollResponse ListBankBusinessPayroll(1: ListBusinessPayrollRequest req)
	ListBusinessPayrollDetailResponse ListBankBusinessPayrollDetail(1: ListBusinessPayrollDetailRequest req)
	SyncBusinessPayrollResultResponse SyncBankBusinessPayrollDetail(1: SyncBusinessPayrollResultRequest req)
	void HandleTransactionDetailReceipt(1: string beginDate, 2: string endDate, 3: i64 organizationId)
	CreateVirtualAccountResponse createVirtualAccount(1: CreateVirtualAccountRequest req)
	void SyncVirtualAccountBalance()
	VirtualAccountBalanceData QueryVirtualAccountBalance(1: i64 organizationId,2:string bankType)
	BankVirtualAccountTranscationResponse spdBankVirtualAccountTranscation(1: i64 organizationId, 2: BankTransferReceiptData req)

    ListPaymentReceiptResponse listPaymentReceipt(1: ListPaymentReceiptRequest req)
    PaymentReceiptData getPaymentReceipt(1: i64 id)
    void addPaymentReceipt(1: PaymentReceiptData req)
    void approvePaymentReceipt(1: i64 id, 2: PaymentReceiptData req)
    void refusePaymentReceipt(1: i64 id, 2: PaymentReceiptData req, 3: string remark)
    void paymentReceiptRun(1: i64 id)
    void transmitPaymentReceipt(1: i64 processInstanceId, 2: i64 transmitUserId)
    void sendBackPaymentApplication(1: i64 id, 2: PaymentReceiptData req, 3: string remark)
    void withDrawPaymentReceipt(1: i64 id, 2: PaymentReceiptData req)
    void commentPaymentReceipt(1: PaymentReceiptData req)
    void addTagPaymentReceipt(1: PaymentReceiptData req)

    void handleSyncPaymentReceipt(1: string beginDate, 2: string endDate, 3: i64 organizationId)
    void syncPaymentReceipt(1: i64 taskId, 2: binary param, 3: i64 organizationId)
    pinganUserAcctSignatureApplyResponse pinganBankAccountSignatureApply(1:pinganBankAccountSignatureApplyRequest req)
    pinganUserAcctSignatureApplyResponse pinganBankAccountSignatureQuery(1:pinganBankAccountSignatureApplyRequest req)
    void systemRefusePaymentReceipt(1: i64 id)
    void systemApprovePaymentReceipt(1: i64 id)

}
