// Package usecase berisi logika bisnis: validasi input dan alur tiap fitur.
package usecase

import (
	"fmt"
	"strings"

	"lc1-order-api/domain"
)

type orderUsecase struct {
	repo domain.OrderRepository
}

// NewOrderUsecase membuat usecase order dengan repository yang diberikan.
func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{repo: repo}
}

// validate merapikan spasi lalu memeriksa semua field wajib.
func validate(in *domain.OrderInput) error {
	in.BuyerName = strings.TrimSpace(in.BuyerName)
	in.StoreName = strings.TrimSpace(in.StoreName)
	in.ItemName = strings.TrimSpace(in.ItemName)

	switch {
	case in.BuyerName == "":
		return fmt.Errorf("%w: buyer_name is required", domain.ErrInvalidInput)
	case in.StoreName == "":
		return fmt.Errorf("%w: store_name is required", domain.ErrInvalidInput)
	case in.ItemName == "":
		return fmt.Errorf("%w: item_name is required", domain.ErrInvalidInput)
	case in.ItemQty <= 0:
		return fmt.Errorf("%w: item_qty must be greater than 0", domain.ErrInvalidInput)
	}
	return nil
}

func (u *orderUsecase) Create(in domain.OrderInput) (*domain.Order, error) {
	if err := validate(&in); err != nil {
		return nil, err
	}
	o := &domain.Order{BuyerName: in.BuyerName, StoreName: in.StoreName, ItemName: in.ItemName, ItemQty: in.ItemQty}
	if err := u.repo.Create(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (u *orderUsecase) GetByID(id int64) (*domain.Order, error) {
	return u.repo.GetByID(id)
}

func (u *orderUsecase) GetAll() ([]domain.Order, error) {
	return u.repo.GetAll()
}

func (u *orderUsecase) Update(id int64, in domain.OrderInput) (*domain.Order, error) {
	if err := validate(&in); err != nil {
		return nil, err
	}
	// Cek dulu apakah order ada (MySQL melaporkan 0 baris berubah bila isinya sama persis).
	o, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	o.BuyerName, o.StoreName, o.ItemName, o.ItemQty = in.BuyerName, in.StoreName, in.ItemName, in.ItemQty
	if err := u.repo.Update(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (u *orderUsecase) Delete(id int64) error {
	// Cek dulu apakah order ada supaya bisa mengembalikan 404 bila tidak ditemukan.
	if _, err := u.repo.GetByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}
