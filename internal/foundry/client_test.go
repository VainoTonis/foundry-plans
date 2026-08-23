package foundry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreatePlanCreatesPlanAndSteps(t *testing.T) {
	var requests []struct {
		path string
		body map[string]interface{}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			requests = append(requests, struct {
				path string
				body map[string]interface{}
			}{r.URL.Path, body})
		}

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/plans":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":7,"title":"Plan","summary":"Summary","status":"pending","repositories":[{"id":11,"name":"one"},{"id":12,"name":"two"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/plans/7/steps":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/plans/7":
			_, _ = w.Write([]byte(`{"id":7,"title":"Plan","summary":"Summary","status":"pending","repositories":[{"id":11,"name":"one"},{"id":12,"name":"two"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/plans/7/steps":
			_, _ = w.Write([]byte(`[{"id":21,"plan_id":7,"position":1,"text":"first","status":"pending"},{"id":22,"plan_id":7,"position":2,"text":"second","status":"pending","parallel_group":3}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	parallelGroup := 3
	plan, err := NewClient(server.URL).CreatePlan(
		[]int64{11, 12}, "Plan", "Summary",
		[]CreateStepInput{{Text: "first"}, {Text: "second", ParallelGroup: &parallelGroup}},
	)
	if err != nil {
		t.Fatalf("CreatePlan() error = %v", err)
	}
	if got := plan.Repositories; !reflect.DeepEqual(got, []Repository{{ID: 11, Name: "one"}, {ID: 12, Name: "two"}}) {
		t.Fatalf("repositories = %#v", got)
	}
	if len(plan.Steps) != 2 || plan.Steps[1].ParallelGroup == nil || *plan.Steps[1].ParallelGroup != 3 {
		t.Fatalf("steps = %#v", plan.Steps)
	}

	if len(requests) != 3 {
		t.Fatalf("POST request count = %d, want 3", len(requests))
	}
	if got := requests[0].body["repository_ids"]; !reflect.DeepEqual(got, []interface{}{float64(11), float64(12)}) {
		t.Fatalf("repository_ids = %#v", got)
	}
	if _, exists := requests[0].body["repo_name"]; exists {
		t.Fatal("create payload contains obsolete repo_name")
	}
	if got := requests[1].body; !reflect.DeepEqual(got, map[string]interface{}{"position": float64(1), "status": "pending", "text": "first"}) {
		t.Fatalf("first step payload = %#v", got)
	}
	if got := requests[2].body["parallel_group"]; got != float64(3) {
		t.Fatalf("parallel_group = %#v", got)
	}
}

func TestCreatePlanRejectsInvalidRepositoryIDs(t *testing.T) {
	client := NewClient("http://unused")
	for _, ids := range [][]int64{nil, {}, {0}, {-1}} {
		if _, err := client.CreatePlan(ids, "Plan", "", nil); err == nil {
			t.Fatalf("CreatePlan(%v) error = nil", ids)
		}
	}
}

func TestCreatePlanAPIFailures(t *testing.T) {
	t.Run("plan", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "no plan", http.StatusBadRequest)
		}))
		defer server.Close()

		if _, err := NewClient(server.URL).CreatePlan([]int64{1}, "Plan", "", nil); err == nil {
			t.Fatal("CreatePlan() error = nil")
		}
	})

	t.Run("step", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/plans" {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":7,"repositories":[{"id":1,"name":"one"}]}`))
				return
			}
			http.Error(w, "no step", http.StatusInternalServerError)
		}))
		defer server.Close()

		if _, err := NewClient(server.URL).CreatePlan([]int64{1}, "Plan", "", []CreateStepInput{{Text: "step"}}); err == nil {
			t.Fatal("CreatePlan() error = nil")
		}
	})
}
