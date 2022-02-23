package models

import (
	"reflect"
	"sync"
	"testing"
)

func TestCache_Get(t *testing.T) {
	type fields struct {
		mtx  sync.RWMutex
		data map[string]Order
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Order
		wantErr bool
	}{
		{
			name: "Нормальное поведение",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: map[string]Order{"1": {}},
			},
			args:    args{id: "1"},
			want:    Order{},
			wantErr: false,
		},
		{
			name: "Неверный id",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: map[string]Order{"1": {}},
			},
			args:    args{id: "2"},
			want:    Order{},
			wantErr: true,
		},
		{
			name: "Неициализированный кеш",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: nil,
			},
			args:    args{id: "1"},
			want:    Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mtx:  tt.fields.mtx,
				data: tt.fields.data,
			}
			got, err := c.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Insert(t *testing.T) {
	type fields struct {
		mtx  sync.RWMutex
		data map[string]Order
	}
	type args struct {
		data Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Нормальное поведение",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: map[string]Order{"1": {}},
			},
			args: args{
				data: Order{OrderUid: "2"},
			},
			wantErr: false,
		},
		{
			name: "Добавление уже существующего",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: map[string]Order{"1": {}},
			},
			args: args{
				data: Order{OrderUid: "1"},
			},
			wantErr: true,
		},
		{
			name: "Неициализированный кеш",
			fields: fields{
				mtx:  sync.RWMutex{},
				data: nil,
			},
			args: args{
				data: Order{OrderUid: "1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mtx:  tt.fields.mtx,
				data: tt.fields.data,
			}
			if err := c.Insert(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
