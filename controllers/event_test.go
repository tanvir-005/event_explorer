package controllers

import (
	"testing"
)

func TestValidateTicketURLAcceptsValidTicketmasterURL(t *testing.T) {
	err := validateTicketURL("https://www.ticketmaster.com/event-123")

	if err != nil {
		t.Fatalf("expected valid Ticketmaster URL, got error: %v", err)
	}
}

func TestValidateTicketURLRejectsHTTP(t *testing.T) {
	err := validateTicketURL("http://www.ticketmaster.com/event-123")

	if err == nil {
		t.Fatal("expected HTTP URL to be rejected")
	}
}

func TestValidateTicketURLRejectsUnapprovedHost(t *testing.T) {
	err := validateTicketURL("https://evil.example.com/event-123")

	if err == nil {
		t.Fatal("expected unapproved host to be rejected")
	}
}

func TestValidateTicketURLRejectsTicketmasterSubdomainAttack(t *testing.T) {
	err := validateTicketURL("https://www.ticketmaster.com.evil.example.com/event-123")

	if err == nil {
		t.Fatal("expected malicious hostname to be rejected")
	}
}

func TestValidateTicketURLRejectsMissingURL(t *testing.T) {
	err := validateTicketURL("")

	if err == nil {
		t.Fatal("expected empty URL to be rejected")
	}
}

func TestValidateTicketURLRejectsMalformedURL(t *testing.T) {
	err := validateTicketURL("://invalid-url")

	if err == nil {
		t.Fatal("expected malformed URL to be rejected")
	}
}
