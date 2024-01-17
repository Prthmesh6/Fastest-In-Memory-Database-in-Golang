package main

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func getReader(input string) *bufio.Reader {
	reader := bufio.NewReader(strings.NewReader(input))
	return reader
}

func TestResp_readLine(t *testing.T) {
	input := "Hello\r\n"
	// reader := bufio.NewReader(strings.NewReader(input))
	type fields struct {
		reader *bufio.Reader
	}
	tests := []struct {
		name     string
		fields   fields
		wantLine []byte
		wantN    int
		wantErr  bool
	}{
		{
			name: "test1",
			fields: fields{
				reader: getReader(input),
			},
			wantLine: []byte("Hello"),
			wantErr:  false,
			wantN:    7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resp{
				reader: tt.fields.reader,
			}
			gotLine, gotN, err := r.readLine()
			fmt.Println(gotLine)
			fmt.Println(gotN)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resp.readLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLine, tt.wantLine) {
				t.Errorf("Resp.readLine() gotLine = %v, want %v", gotLine, tt.wantLine)
			}
			if gotN != tt.wantN {
				t.Errorf("Resp.readLine() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestResp_readInteger(t *testing.T) {
	type fields struct {
		reader *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantX   int
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testReadInteger",
			fields: fields{
				reader: getReader("369\r\n"),
			},
			wantX:   369,
			wantN:   5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resp{
				reader: tt.fields.reader,
			}
			gotX, gotN, err := r.readInteger()
			fmt.Printf("gotX = %v, gotN = %v, err = %v", gotX, gotN, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resp.readInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotX != tt.wantX {
				t.Errorf("Resp.readInteger() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotN != tt.wantN {
				t.Errorf("Resp.readInteger() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
