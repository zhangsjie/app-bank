package repo

import (
	"github.com/google/wire"
	"gitlab.yoyiit.com/youyi/go-core/repository"
	"gitlab.yoyiit.com/youyi/go-core/store"
	"gitlab.yoyiit.com/youyi/go-core/trace"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(trace.NewGormLogger, store.NewReadWriteSeparationDB, NewBankTransferReceiptRepo,
	NewBankTransactionDetailRepo, NewBankTransactionDetailExternalRepo, NewBankCodeRepo, NewBankBusinessPayrollRepo,
	NewBankBusinessPayrollDetailRepo, NewPaymentReceiptRepo)

func NewBankTransferReceiptRepo(db *gorm.DB) BankTransferReceiptRepo {
	return &bankTransferReceiptRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankTransferReceiptDBData{}},
	}
}

func NewBankTransactionDetailRepo(db *gorm.DB) BankTransactionDetailRepo {
	return &bankTransactionDetailRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankTransactionDetailDBData{}},
	}
}

func NewBankTransactionDetailExternalRepo(db *gorm.DB) BankTransactionDetailProcessInstanceRepo {
	return &bankTransactionDetailProcessInstanceRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankTransactionDetailProcessInstanceDBData{}},
	}
}

func NewBankCodeRepo(db *gorm.DB) BankCodeRepo {
	return &bankCodeRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankCodeDBData{}},
	}
}

func NewBankBusinessPayrollRepo(db *gorm.DB) BankBusinessPayrollRepo {
	return &bankBusinessPayrollRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankBusinessPayrollDBData{}},
	}
}

func NewBankBusinessPayrollDetailRepo(db *gorm.DB) BankBusinessPayrollDetailRepo {
	return &bankBusinessPayrollDetailRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: BankBusinessPayrollDetailDBData{}},
	}
}

func NewPaymentReceiptRepo(db *gorm.DB) PaymentReceiptRepo {
	return &paymentReceiptRepo{
		BaseRepo: repository.BaseRepo{Db: db, Model: PaymentReceiptDBData{}},
	}
}
