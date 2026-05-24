package model

import "testing"

func TestIsValidTransition_ValidPaths(t *testing.T) {
	tests := []struct {
		from string
		to   string
	}{
		{StatusReceived, StatusPrimaryScoring},
		{StatusReceived, StatusRejected},
		{StatusPrimaryScoring, StatusSecurityCheck},
		{StatusPrimaryScoring, StatusRejected},
		{StatusPrimaryScoring, StatusRevision},
		{StatusSecurityCheck, StatusCollateralExpertise},
		{StatusCollateralExpertise, StatusLegalCheck},
		{StatusLegalCheck, StatusCreditAnalysis},
		{StatusCreditAnalysis, StatusCreditCommittee},
		{StatusCreditCommittee, StatusApproved},
		{StatusCreditCommittee, StatusRejected},
		{StatusApproved, StatusDocumentation},
		{StatusRevision, StatusPrimaryScoring},
		{StatusRevision, StatusRejected},
		{StatusDocumentation, StatusIssued},
	}

	for _, tc := range tests {
		if !IsValidTransition(tc.from, tc.to) {
			t.Errorf("expected valid: %s → %s", tc.from, tc.to)
		}
	}
}

func TestIsValidTransition_InvalidPaths(t *testing.T) {
	tests := []struct {
		from string
		to   string
	}{
		{StatusReceived, StatusIssued},
		{StatusReceived, StatusApproved},
		{StatusApproved, StatusRejected},
		{StatusIssued, StatusReceived},
		{StatusRejected, StatusApproved},
		{StatusRejected, StatusReceived},
		{StatusDocumentation, StatusRejected},
		{StatusPrimaryScoring, StatusApproved},
	}

	for _, tc := range tests {
		if IsValidTransition(tc.from, tc.to) {
			t.Errorf("expected invalid: %s → %s", tc.from, tc.to)
		}
	}
}

func TestIsValidTransition_UnknownFromStatus(t *testing.T) {
	if IsValidTransition("unknown_status", StatusApproved) {
		t.Error("expected false for unknown from-status")
	}
}

func TestIsValidTransition_TerminalStatuses(t *testing.T) {
	terminals := []string{StatusRejected, StatusIssued}
	targets := []string{StatusReceived, StatusApproved, StatusIssued, StatusRejected}

	for _, term := range terminals {
		for _, to := range targets {
			if IsValidTransition(term, to) {
				t.Errorf("terminal status %s should not transition to %s", term, to)
			}
		}
	}
}
