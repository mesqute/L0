package models

import (
	tests2 "L0/models/tests"
	"reflect"
	"testing"
)

func TestGetFieldValues(t *testing.T) {
	type args struct {
		strct interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "Структура со стандартными типами",
			args:    args{&tests2.A{}},
			want:    tests2.ConvertToInterfaceSlice(0, "", [4]int{}, 0.0),
			wantErr: false,
		},
		{
			name:    "Структура с обрабатываемыми типами",
			args:    args{&tests2.B{}},
			want:    tests2.ConvertToInterfaceSlice(0, "", [1]int{}, 0.0),
			wantErr: false,
		},
		{
			name:    "Структура с необрабатываемыми типами",
			args:    args{&tests2.C{}},
			want:    tests2.ConvertToInterfaceSlice(), // вернет пустой слайс
			wantErr: false,
		},
		{
			name:    "Не структура",
			args:    args{&tests2.D},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldValues(tt.args.strct)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldValues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_Validate(t *testing.T) {
	type fields struct {
		OrderUid          string
		TrackNumber       string
		Entry             string
		Delivery          Delivery
		Payment           Payment
		Items             []Item
		Locale            string
		InternalSignature string
		CustomerId        string
		DeliveryService   string
		Shardkey          string
		SmId              int
		DateCreated       string
		OofShard          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Минимально необходимые поля",
			fields: fields{
				OrderUid:    "1",
				TrackNumber: "1",
				Entry:       "1",
				Delivery:    Delivery{Name: "1", Phone: "1", City: "1", Address: "1", Email: "a@m.c"},
				Payment:     Payment{Transaction: "1", Currency: "1"},
				Items:       []Item{{ChrtId: 1, Status: 1}},
				CustomerId:  "2",
				DateCreated: "2",
			},
			wantErr: false,
		},
		{
			name:    "Пустая структура",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Order{
				OrderUid:          tt.fields.OrderUid,
				TrackNumber:       tt.fields.TrackNumber,
				Entry:             tt.fields.Entry,
				Delivery:          tt.fields.Delivery,
				Payment:           tt.fields.Payment,
				Items:             tt.fields.Items,
				Locale:            tt.fields.Locale,
				InternalSignature: tt.fields.InternalSignature,
				CustomerId:        tt.fields.CustomerId,
				DeliveryService:   tt.fields.DeliveryService,
				Shardkey:          tt.fields.Shardkey,
				SmId:              tt.fields.SmId,
				DateCreated:       tt.fields.DateCreated,
				OofShard:          tt.fields.OofShard,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
