package model

import "testing"

func TestRolePriority(t *testing.T) {
	cases := []struct {
		role string
		want int
	}{
		{RoleAdmin, 4},
		{RoleExpert, 3},
		{RoleEmployee, 2},
		{RoleBorrower, 1},
		{RolePublic, 0},
		{"unknown", 0},
		{"", 0},
	}

	for _, tc := range cases {
		got := RolePriority(tc.role)
		if got != tc.want {
			t.Errorf("RolePriority(%q) = %d, want %d", tc.role, got, tc.want)
		}
	}
}

func TestHighestRole_SingleRole(t *testing.T) {
	cases := []struct {
		roles []string
		want  string
	}{
		{[]string{"admin"}, "admin"},
		{[]string{"expert"}, "expert"},
		{[]string{"employee"}, "employee"},
		{[]string{"borrower"}, "borrower"},
		{[]string{"public"}, "public"},
	}

	for _, tc := range cases {
		got := HighestRole(tc.roles)
		if got != tc.want {
			t.Errorf("HighestRole(%v) = %q, want %q", tc.roles, got, tc.want)
		}
	}
}

func TestHighestRole_MultipleRoles(t *testing.T) {
	// Highest priority wins
	got := HighestRole([]string{"borrower", "employee", "admin"})
	if got != "admin" {
		t.Errorf("want admin, got %q", got)
	}

	got = HighestRole([]string{"employee", "expert"})
	if got != "expert" {
		t.Errorf("want expert, got %q", got)
	}

	got = HighestRole([]string{"borrower", "employee"})
	if got != "employee" {
		t.Errorf("want employee, got %q", got)
	}
}

func TestHighestRole_EmptySlice(t *testing.T) {
	got := HighestRole([]string{})
	if got != "public" {
		t.Errorf("want public for empty slice, got %q", got)
	}
}

func TestHighestRole_UnknownRoles(t *testing.T) {
	got := HighestRole([]string{"unknown", "garbage"})
	if got != "public" {
		t.Errorf("want public for unknown roles, got %q", got)
	}
}

func TestHighestRole_MixedKnownUnknown(t *testing.T) {
	got := HighestRole([]string{"unknown", "employee"})
	if got != "employee" {
		t.Errorf("want employee, got %q", got)
	}
}
