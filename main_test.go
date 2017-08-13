package main

import "testing"

func TestAlertTitleMinifyExpectHello(t *testing.T) {
	input := "[hello]"
	v := alertTitleMinify(input)

	expected := "hello"
	if v != "hello" {
		t.Errorf("Result: %s, Expected: %s", v, expected)
	}
}

func TestAlertTitleMinifyExpectHelloWorld(t *testing.T) {
	input := "[hello world]"
	v := alertTitleMinify(input)

	expected := "hello world"
	if v != "hello world" {
		t.Errorf("Result: %s, Expected: %s", v, expected)
	}
}

func TestStringMinifier(t *testing.T) {
	input := "kafka.consumer.totallag over \nburrow_metrics_instance:01\n<= 15000.0 on average during the last 1h."
	v := stringMinifier(input)

	expected := "kafka.consumer.totallag over burrow_metrics_instance:01 <= 15000.0 on average during the last 1h."

	if v != expected {
		t.Errorf("Result: %s, Expected: %s", v, expected)
	}
}
