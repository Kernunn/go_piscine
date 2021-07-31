package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSorted(t *testing.T) {
	cmd := exec.Command("./statistics")
	cmd.Stdin = strings.NewReader("1\n2\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Mean: 1.50\nMedian: 1.50\nMode: 1\nSD: 0.50\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestNotSorted(t *testing.T) {
	cmd := exec.Command("./statistics")
	cmd.Stdin = strings.NewReader("2\n1\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Mean: 1.50\nMedian: 1.50\nMode: 1\nSD: 0.50\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestEmptyLine(t *testing.T) {
	cmd := exec.Command("./statistics")
	cmd.Stdin = strings.NewReader("2\n\n")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error\n")
	}
}

func TestIncorrectInput(t *testing.T) {
	cmd := exec.Command("./statistics")
	cmd.Stdin = strings.NewReader("2\nd\n")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error\n")
	}
}

func TestOutOfRange(t *testing.T) {
	cmd := exec.Command("./statistics")
	cmd.Stdin = strings.NewReader("1000000\n")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error\n")
	}
}

func TestOnlyMean(t *testing.T) {
	cmd := exec.Command("./statistics", "--mean")
	cmd.Stdin = strings.NewReader("2\n1\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Mean: 1.50\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestOnlyMedianAndMode(t *testing.T) {
	cmd := exec.Command("./statistics", "--median", "--mode")
	cmd.Stdin = strings.NewReader("2\n1\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Median: 1.50\nMode: 1\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestMean(t *testing.T) {
	cmd := exec.Command("./statistics", "--mean")
	cmd.Stdin = strings.NewReader("2\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Mean: 2.00\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestMode(t *testing.T) {
	cmd := exec.Command("./statistics", "--mode")
	cmd.Stdin = strings.NewReader("2\n2\n1\n1\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "Mode: 1\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}

func TestSD(t *testing.T) {
	cmd := exec.Command("./statistics", "--sd")
	cmd.Stdin = strings.NewReader("2\n1\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	result := string(out)
	expected := "SD: 0.50\n"
	if result != expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v",
			result, expected)
	}
}
