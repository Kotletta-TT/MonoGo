package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"reflect"
	"testing"
)

func TestNewCollector(t *testing.T) {
	type args struct {
		repo      *repository.Repository
		wantPanic bool
	}
	tests := []struct {
		name string
		args args
		want Collector
	}{
		{
			name: "Empty repository",
			args: args{
				repo:      nil,
				wantPanic: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.args.wantPanic {
					t.Errorf("Not expected panic but got")
				}
			}()
			collect := NewCollector(tt.args.repo)
			if tt.args.wantPanic {
				return
			}
			if !reflect.DeepEqual(collect, tt.want) {
				t.Errorf("NewCollector() = %v, want %v", collect, tt.want)
			}
		})
	}
}
