package usecase

import (
	"errors"
	"testing"

	"lc1-order-api/domain"
)

// fakeRepo adalah repository palsu (data di memori) agar test tidak butuh MySQL.
type fakeRepo struct {
	data   map[int64]domain.Order
	nextID int64
}

func newFakeRepo() *fakeRepo { return &fakeRepo{data: map[int64]domain.Order{}, nextID: 1} }

func (f *fakeRepo) Create(o *domain.Order) error {
	o.ID, o.CreatedAt = f.nextID, "2026-11-16"
	f.nextID++
	f.data[o.ID] = *o
	return nil
}
func (f *fakeRepo) GetByID(id int64) (*domain.Order, error) {
	o, ok := f.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &o, nil
}
func (f *fakeRepo) GetAll() ([]domain.Order, error) {
	list := []domain.Order{}
	for id := int64(1); id < f.nextID; id++ {
		if o, ok := f.data[id]; ok {
			list = append(list, o)
		}
	}
	return list, nil
}
func (f *fakeRepo) Update(o *domain.Order) error { f.data[o.ID] = *o; return nil }
func (f *fakeRepo) Delete(id int64) error        { delete(f.data, id); return nil }

func TestCreateValidation(t *testing.T) {
	uc := NewOrderUsecase(newFakeRepo())
	tests := []struct {
		name    string
		in      domain.OrderInput
		wantErr bool
	}{
		{"valid", domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemName: "Kaos", ItemQty: 2}, false},
		{"buyer kosong", domain.OrderInput{StoreName: "Toko A", ItemName: "Kaos", ItemQty: 2}, true},
		{"store kosong", domain.OrderInput{BuyerName: "Budi", ItemName: "Kaos", ItemQty: 2}, true},
		{"item kosong", domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemQty: 2}, true},
		{"qty nol", domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemName: "Kaos", ItemQty: 0}, true},
		{"qty negatif", domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemName: "Kaos", ItemQty: -1}, true},
		{"hanya spasi", domain.OrderInput{BuyerName: "   ", StoreName: "Toko A", ItemName: "Kaos", ItemQty: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Create(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, domain.ErrInvalidInput) {
				t.Errorf("error harus ErrInvalidInput, dapat: %v", err)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	uc := NewOrderUsecase(newFakeRepo())
	ok := domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemName: "Kaos", ItemQty: 1}
	if _, err := uc.GetByID(99); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByID: mau ErrNotFound, dapat %v", err)
	}
	if _, err := uc.Update(99, ok); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Update: mau ErrNotFound, dapat %v", err)
	}
	if err := uc.Delete(99); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Delete: mau ErrNotFound, dapat %v", err)
	}
}

func TestUpdateAndDelete(t *testing.T) {
	uc := NewOrderUsecase(newFakeRepo())
	o, _ := uc.Create(domain.OrderInput{BuyerName: "Budi", StoreName: "Toko A", ItemName: "Kaos", ItemQty: 1})
	u, err := uc.Update(o.ID, domain.OrderInput{BuyerName: "Budi S", StoreName: "Toko B", ItemName: "Topi", ItemQty: 5})
	if err != nil || u.ItemQty != 5 || u.StoreName != "Toko B" {
		t.Fatalf("update gagal: %+v, %v", u, err)
	}
	if err := uc.Delete(o.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := uc.GetByID(o.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Error("setelah delete, data harus tidak ditemukan")
	}
}
