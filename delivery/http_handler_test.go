package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"lc1-order-api/domain"
	"lc1-order-api/usecase"
)

func init() { gin.SetMode(gin.TestMode) }

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

// do mengirim satu request ke router dan mengembalikan status code serta body.
func do(h http.Handler, method, path, body string) (int, string) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

const orderJSON = `{"buyer_name":"Budi","store_name":"Toko A","item_name":"Kaos","item_qty":2}`

func TestOrderEndpoints(t *testing.T) {
	h := NewRouter(usecase.NewOrderUsecase(newFakeRepo()))
	tests := []struct {
		name         string
		method, path string
		body         string
		want         int
	}{
		{"list kosong", "GET", "/orders", "", 200},
		{"create sukses", "POST", "/orders", orderJSON, 201},
		{"create buyer kosong", "POST", "/orders", `{"buyer_name":"","store_name":"A","item_name":"B","item_qty":1}`, 400},
		{"create qty nol", "POST", "/orders", `{"buyer_name":"A","store_name":"A","item_name":"B","item_qty":0}`, 400},
		{"create qty bukan angka", "POST", "/orders", `{"buyer_name":"A","store_name":"A","item_name":"B","item_qty":"banyak"}`, 400},
		{"create bukan json", "POST", "/orders", `bukan json`, 400},
		{"get by id sukses", "GET", "/orders/1", "", 200},
		{"get by id tidak ada", "GET", "/orders/99", "", 404},
		{"get id bukan angka", "GET", "/orders/abc", "", 404},
		{"put sukses", "PUT", "/orders/1", orderJSON, 200},
		{"put tidak ada", "PUT", "/orders/99", orderJSON, 404},
		{"put input salah", "PUT", "/orders/1", `{"buyer_name":""}`, 400},
		{"delete sukses", "DELETE", "/orders/1", "", 200},
		{"delete tidak ada", "DELETE", "/orders/1", "", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, body := do(h, tt.method, tt.path, tt.body)
			if got != tt.want {
				t.Errorf("status = %d, mau %d, body: %s", got, tt.want, body)
			}
		})
	}
}

// Memastikan bentuk response sesuai kontrak API.
func TestContractShape(t *testing.T) {
	h := NewRouter(usecase.NewOrderUsecase(newFakeRepo()))

	_, body := do(h, "POST", "/orders", orderJSON)
	var created map[string]any
	json.Unmarshal([]byte(body), &created)
	if created["message"] != "success create" || created["shipment"] == nil {
		t.Errorf("response POST harus memuat message 'success create' dan key shipment: %s", body)
	}

	_, body = do(h, "GET", "/orders/99", "")
	if !strings.Contains(body, "order with ID 99 not found") {
		t.Errorf("pesan 404 tidak sesuai kontrak: %s", body)
	}

	_, body = do(h, "GET", "/orders", "")
	if !strings.HasPrefix(body, "[") {
		t.Errorf("GET /orders harus berupa array: %s", body)
	}
}
