package service

import (
	"context"
	"echo-crud/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	// ErrNilSupplier occurs when a nil supplier is passed.
	ErrNilSupplier = errors.New("supplier is nil")
)

// SupplierService responsible for any flow related to supplier.
// It also implements SupplierService.
type SupplierService struct {
	supplierRepo SupplierRepository
}

func NewSupplierService(supplierRepo SupplierRepository) *SupplierService {
	return &SupplierService{
		supplierRepo: supplierRepo,
	}
}

type SupplierUseCase interface {
	Create(ctx context.Context, supplier *entity.Supplier) error
	GetListSupplier(ctx context.Context, limit, offset string) ([]*entity.Supplier, error)
	GetDetailSupplier(ctx context.Context, ID uuid.UUID) (*entity.Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *entity.Supplier) error
	DeleteSupplier(ctx context.Context, ID uuid.UUID) error
}

type SupplierRepository interface {
	Insert(ctx context.Context, supplier *entity.Supplier) error
	GetListSupplier(ctx context.Context, limit, offset string) ([]*entity.Supplier, error)
	GetDetailSupplier(ctx context.Context, ID uuid.UUID) (*entity.Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *entity.Supplier) error
	DeleteSupplier(ctx context.Context, ID uuid.UUID) error
}

func (svc SupplierService) Create(ctx context.Context, supplier *entity.Supplier) error {
	// Checking nil file
	if supplier == nil {
		return ErrNilSupplier
	}

	// Generate id if nil
	if supplier.ID == uuid.Nil {
		supplier.ID = uuid.New()
	}

	if err := svc.supplierRepo.Insert(ctx, supplier); err != nil {
		return errors.Wrap(err, "[SupplierService-Create]")
	}
	return nil
}

func (svc SupplierService) GetListSupplier(ctx context.Context, limit, offset string) ([]*entity.Supplier, error) {
	supplier, err := svc.supplierRepo.GetListSupplier(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SupplierService-Create]")
	}
	return supplier, nil
}

func (svc SupplierService) GetDetailSupplier(ctx context.Context, ID uuid.UUID) (*entity.Supplier, error) {
	supplier, err := svc.supplierRepo.GetDetailSupplier(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "[SupplierService-Create]")
	}
	return supplier, nil
}

func (svc SupplierService) UpdateSupplier(ctx context.Context, supplier *entity.Supplier) error {
	// Checking nil supplier
	if supplier == nil {
		return ErrNilSupplier
	}

	// Generate id if nil
	if supplier.ID == uuid.Nil {
		supplier.ID = uuid.New()
	}

	if err := svc.supplierRepo.UpdateSupplier(ctx, supplier); err != nil {
		return errors.Wrap(err, "[SupplierService-Create]")
	}
	return nil
}

func (svc SupplierService) DeleteSupplier(ctx context.Context, ID uuid.UUID) error {
	err := svc.supplierRepo.DeleteSupplier(ctx, ID)
	if err != nil {
		return errors.Wrap(err, "[SupplierService-Create]")
	}
	return nil
}
