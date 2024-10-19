package main

import (
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	data := "HELLO\r\n"
	reader := NewResp(strings.NewReader(data))

	line, _, err := reader.readLine()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expected := "HELLO"
	if string(line) != expected {
		t.Fatalf("Expected %s, but got %s", expected, string(line))
	}
}

func TestReadInteger(t *testing.T) {
	data := "100\r\n"
	reader := NewResp(strings.NewReader(data))

	num, _, err := reader.readInteger()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expected := int64(100)
	if num != expected {
		t.Fatalf("Expected %d, but got %d", expected, num)
	}
}

func TestReadBulk(t *testing.T) {
	data := "$5\r\nHELLO\r\n"
	reader := NewResp(strings.NewReader(data))

	val, err := reader.Read()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if val.typ != "bulk" {
		t.Fatalf("Expected type bulk, but got %s", val.typ)
	}

	expected := "HELLO"
	if val.bulk != expected {
		t.Fatalf("Expected %s, but got %s", expected, val.bulk)
	}
}

func TestReadArray(t *testing.T) {
	data := "2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	reader := NewResp(strings.NewReader(data))
	value, err := reader.readArray()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(value.array) != 2 {
		t.Fatalf("expected array length 2, got %d", len(value.array))
	}
	if value.array[0].bulk != "foo" || value.array[1].bulk != "bar" {
		t.Fatalf("expected [foo bar], got [%s %s]", value.array[0].bulk, value.array[1].bulk)
	}
}

func TestSetAndGet(t *testing.T) {
	setCmd := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	reader := NewResp(strings.NewReader(setCmd))
	value, err := reader.Read()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := set(value.array[1:])
	if result.str != "OK" {
		t.Fatalf("expected OK, got %s", result.str)
	}

	getCmd := "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	reader = NewResp(strings.NewReader(getCmd))
	value, err = reader.Read()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result = get(value.array[1:])
	if result.bulk != "value" {
		t.Fatalf("expected value, got %s", result.bulk)
	}
}

func TestMarshalBulk(t *testing.T) {
	val := Value{typ: "bulk", bulk: "HELLO"}

	result := string(val.Marshal())
	expected := "$5\r\nHELLO\r\n"

	if result != expected {
		t.Fatalf("Expected %s, but got %s", expected, result)
	}
}

func TestMarshalArray(t *testing.T) {
	val := Value{typ: "array", array: []Value{
		{typ: "bulk", bulk: "HELLO"},
		{typ: "bulk", bulk: "WORLD"},
	}}

	result := string(val.Marshal())
	expected := "*2\r\n$5\r\nHELLO\r\n$5\r\nWORLD\r\n"

	if result != expected {
		t.Fatalf("Expected %s, but got %s", expected, result)
	}
}

func TestWriterWrite(t *testing.T) {
	val := Value{typ: "bulk", bulk: "HELLO"}
	var buf strings.Builder
	writer := NewWriter(&buf)

	err := writer.Write(val)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expected := "$5\r\nHELLO\r\n"
	if buf.String() != expected {
		t.Fatalf("Expected %s, but got %s", expected, buf.String())
	}
}
