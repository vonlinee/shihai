package utils

import (
	"encoding/json"
	"testing"
)

func TestClaimsMarshalUserIDAsString(t *testing.T) {
	claims := Claims{
		UserID:   706545104525070300,
		Username: "admin",
	}

	data, err := json.Marshal(claims)

	if err != nil {
		t.Fatalf("Marshal error = %v, want nil", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal marshaled claims error = %v, want nil", err)
	}
	if body["userId"] != "706545104525070300" {
		t.Fatalf("userId = %#v, want string ID", body["userId"])
	}
}

func TestClaimsUnmarshalUserIDFromNumber(t *testing.T) {
	var claims Claims

	err := json.Unmarshal([]byte(`{"userId":706545104525070300,"username":"admin"}`), &claims)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	if claims.UserID != 706545104525070300 {
		t.Fatalf("UserID = %d, want 706545104525070300", claims.UserID)
	}
}

func TestClaimsUnmarshalUserIDFromString(t *testing.T) {
	var claims Claims

	err := json.Unmarshal([]byte(`{"userId":"706545104525070300","username":"admin"}`), &claims)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	if claims.UserID != 706545104525070300 {
		t.Fatalf("UserID = %d, want 706545104525070300", claims.UserID)
	}
}
