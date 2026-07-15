package dto

import (
	"encoding/json"
	"testing"
)

func TestBatchDeleteRequestUnmarshalIDsFromStrings(t *testing.T) {
	var req BatchDeleteRequest

	err := json.Unmarshal([]byte(`{"ids":["706545104525070300","706545104525070301"]}`), &req)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	assertIDListEqual(t, req.IDs, []uint64{706545104525070300, 706545104525070301})
}

func TestBatchDeleteRequestUnmarshalIDsFromNumbers(t *testing.T) {
	var req BatchDeleteRequest

	err := json.Unmarshal([]byte(`{"ids":[1,2]}`), &req)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	assertIDListEqual(t, req.IDs, []uint64{1, 2})
}

func TestPoemResponseMarshalIDAsString(t *testing.T) {
	resp := PoemResponse{
		ID:        706545104525070300,
		AuthorID:  706545104525070301,
		DynastyID: 706545104525070302,
		Author: AuthorResponse{
			ID: 706545104525070301,
		},
		Dynasty: DynastyResponse{
			ID: 706545104525070302,
		},
	}

	data, err := json.Marshal(resp)

	if err != nil {
		t.Fatalf("Marshal error = %v, want nil", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal marshaled response error = %v, want nil", err)
	}
	if body["id"] != "706545104525070300" {
		t.Fatalf("id = %#v, want string ID", body["id"])
	}
	if body["authorId"] != "706545104525070301" {
		t.Fatalf("authorId = %#v, want string ID", body["authorId"])
	}
	if body["dynastyId"] != "706545104525070302" {
		t.Fatalf("dynastyId = %#v, want string ID", body["dynastyId"])
	}
}

func TestPoemAnnotationResponseMarshalIDAsString(t *testing.T) {
	resp := PoemAnnotationResponse{
		ID:     706545104525070300,
		PoemID: 706545104525070301,
	}

	data, err := json.Marshal(resp)

	if err != nil {
		t.Fatalf("Marshal error = %v, want nil", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal marshaled response error = %v, want nil", err)
	}
	if body["id"] != "706545104525070300" {
		t.Fatalf("id = %#v, want string ID", body["id"])
	}
	if body["poemId"] != "706545104525070301" {
		t.Fatalf("poemId = %#v, want string ID", body["poemId"])
	}
}

func TestPoemCreateRequestMarshalIDAsNumber(t *testing.T) {
	req := PoemCreateRequest{
		Title:     "title",
		Content:   []string{"content"},
		AuthorID:  706545104525070300,
		DynastyID: 706545104525070301,
	}

	data, err := json.Marshal(req)

	if err != nil {
		t.Fatalf("Marshal error = %v, want nil", err)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal marshaled request error = %v, want nil", err)
	}
	if string(body["authorId"]) != "706545104525070300" {
		t.Fatalf("authorId = %s, want numeric ID", body["authorId"])
	}
	if string(body["dynastyId"]) != "706545104525070301" {
		t.Fatalf("dynastyId = %s, want numeric ID", body["dynastyId"])
	}
}

func TestPoemCreateRequestUnmarshalIDFromStrings(t *testing.T) {
	var req PoemCreateRequest

	err := json.Unmarshal([]byte(`{"title":"title","content":["content"],"authorId":"706545104525070300","dynastyId":"706545104525070301"}`), &req)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	if uint64(req.AuthorID) != 706545104525070300 {
		t.Fatalf("AuthorID = %d, want 706545104525070300", req.AuthorID)
	}
	if uint64(req.DynastyID) != 706545104525070301 {
		t.Fatalf("DynastyID = %d, want 706545104525070301", req.DynastyID)
	}
}

func TestAdminCreateUserRequestUnmarshalRoleIdsFromStrings(t *testing.T) {
	var req AdminCreateUserRequest

	err := json.Unmarshal([]byte(`{"username":"admin","password":"password123","roleIds":["706545104525070300"]}`), &req)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	assertIDListEqual(t, req.RoleIds, []uint64{706545104525070300})
}

func TestAssignRoleRequestUnmarshalRoleIdsFromStrings(t *testing.T) {
	var req AssignRoleRequest

	err := json.Unmarshal([]byte(`{"roleIds":["706545104525070300","706545104525070301"]}`), &req)

	if err != nil {
		t.Fatalf("Unmarshal error = %v, want nil", err)
	}
	assertIDListEqual(t, req.RoleIDs, []uint64{706545104525070300, 706545104525070301})
}

func assertIDListEqual(t *testing.T, got IDList, want []uint64) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d; got %v", i, got[i], want[i], got)
		}
	}
}
