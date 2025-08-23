// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package ipinfo

import "testing"

func Test_IPInfoClientIPv4(t *testing.T) {
	client := NewClient()
	data, err := client.GetIPInfo("8.8.8.8")
	if err != nil {
		t.Errorf("Failed to get IP info, got error: %s", err)
	} else if data == nil {
		t.Error("failed to get valid IP info response")
	} else {
		t.Logf("IP Info: %+v", data)
	}
}

func Test_IPInfoClient_IPv6(t *testing.T) {
	client := NewClient()
	data, err := client.GetIPInfo("2001:4860:4860::8888")
	if err != nil {
		t.Errorf("Failed to get IP info, got error: %s", err)
	} else if data == nil {
		t.Error("failed to get valid IP info response")
	} else {
		t.Logf("IP Info: %+v", data)
	}
}

func Test_IPInfoClient_Private(t *testing.T) {
	client := NewClient()
	_, err := client.GetIPInfo("192.168.0.1")
	if err == nil {
		t.Errorf("expected error for invalid IP")
	} else {
		t.Logf("got error: %s, as expected", err)
	}
}

func Test_IPInfoClient_InvalidIP(t *testing.T) {
	client := NewClient()
	_, err := client.GetIPInfo("demo.invalid")
	if err == nil {
		t.Errorf("expected error for invalid IP")
	} else {
		t.Logf("got error: %s, as expected", err)
	}
}
