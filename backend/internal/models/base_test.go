package models

import (
	"encoding/json"
	"testing"
)

func TestBaseModelMarshalIDAsNumber(t *testing.T) {
	model := BaseModel{
		ID:        706545104525070300,
		CreatedBy: 706545104525070301,
		UpdatedBy: 706545104525070302,
	}

	data, err := json.Marshal(model)

	if err != nil {
		t.Fatalf("Marshal error = %v, want nil", err)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal marshaled model error = %v, want nil", err)
	}
	if string(body["id"]) != "706545104525070300" {
		t.Fatalf("id = %s, want numeric ID", body["id"])
	}
	if string(body["createdBy"]) != "706545104525070301" {
		t.Fatalf("createdBy = %s, want numeric ID", body["createdBy"])
	}
	if string(body["updatedBy"]) != "706545104525070302" {
		t.Fatalf("updatedBy = %s, want numeric ID", body["updatedBy"])
	}
}
