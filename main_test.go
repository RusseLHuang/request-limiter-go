package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLimitRate(t *testing.T) {
	os.Setenv("LIMIT_DURATION", "3")
	os.Setenv("LIMIT", "60")

	ts := httptest.NewServer(Init())

	for i := 1; i <= 60; i++ {
		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			t.Fatal(err)
		}

		expectedValue := i
		bodyString := string(body)
		bodyNumber, _ := strconv.Atoi(bodyString)

		if expectedValue != bodyNumber {
			t.Fatalf("Expected %d got %d", expectedValue, bodyNumber)
		}
	}

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(body)
	errorExpectedValue := "Error"
	if bodyString == errorExpectedValue {
		t.Fatalf("Expected Error response when limit request exceeded")
	}

	time.Sleep(time.Duration(4 * time.Second))

	res, err = http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	bodyString = string(body)
	bodyNumber, _ := strconv.Atoi(bodyString)
	if bodyNumber != 1 {
		t.Fatalf("Expected request rate refreshed %s", bodyString)
	}

	time.Sleep(time.Duration(3 * time.Second))
}
