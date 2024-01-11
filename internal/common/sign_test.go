package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHMACSignature(t *testing.T) {
	type args struct {
		key     string
		message []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Normal",
			args{
				key:     "test",
				message: []byte("test"),
			},
			"88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateHMACSignature(tt.args.key, tt.args.message)
			assert.NoError(t, err)
			t.Log(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVerifyHMACSignature(t *testing.T) {
	type args struct {
		sign string
		key  string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Normal",
			args{
				sign: "88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7",
				key:  "test",
				data: []byte("test"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyHMACSignature(tt.args.sign, tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("VerifyHMACSignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
