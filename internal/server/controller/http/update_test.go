package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		repo repository.Repository
	}
	type args struct {
		wantStatusCode int
		url            string
		method         string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Status OK",
			fields: fields{
				repo: repository.NewMemRepo(),
			},
			args: args{wantStatusCode: 200, url: "/update/gauge/test_metric/100", method: http.MethodPost},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := UpdateHandler{
				repo: tt.fields.repo,
			}
			req := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			uh.ServeHTTP(w, req)
			res := w.Result()
			res.Body.Close()
			if w.Code != tt.args.wantStatusCode {
				t.Errorf("Status code expect: %d got: %d", tt.args.wantStatusCode, w.Code)
			}
		})
	}
}
