package main

import "testing"

func TestSafeInt32ToInt8(t *testing.T) {
    tests := []struct{
        name string
        v int32
        want int8
        wantErr bool
    }{
        {"max", 127, 127, false},
        {"min", -128, -128, false},
        {"overflow", 128, 0, true},
        {"underflow", -129, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := SafeInt32ToInt8(tt.v)
            if (err != nil) != tt.wantErr {
                t.Errorf("ошибка: got %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("значение: got %d, want %d", got, tt.want)
            }
        })
    }
}