package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/go-chi/chi/v5"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockProgramSvc struct {
	listActiveFn func(context.Context) ([]*model.LoanProgram, error)
	listAllFn    func(context.Context) ([]*model.LoanProgram, error)
	getByIDFn    func(context.Context, string) (*model.LoanProgram, error)
	createFn     func(context.Context, repository.CreateProgramInput) (*model.LoanProgram, error)
	updateFn     func(context.Context, string, repository.UpdateProgramInput) (*model.LoanProgram, error)
	deactivateFn func(context.Context, string) error
}

func (m *mockProgramSvc) ListActive(ctx context.Context) ([]*model.LoanProgram, error) {
	if m.listActiveFn != nil {
		return m.listActiveFn(ctx)
	}
	return nil, nil
}
func (m *mockProgramSvc) ListAll(ctx context.Context) ([]*model.LoanProgram, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx)
	}
	return nil, nil
}
func (m *mockProgramSvc) GetByID(ctx context.Context, id string) (*model.LoanProgram, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockProgramSvc) Create(ctx context.Context, inp repository.CreateProgramInput) (*model.LoanProgram, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockProgramSvc) Update(ctx context.Context, id string, inp repository.UpdateProgramInput) (*model.LoanProgram, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inp)
	}
	return nil, nil
}
func (m *mockProgramSvc) Deactivate(ctx context.Context, id string) error {
	if m.deactivateFn != nil {
		return m.deactivateFn(ctx, id)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func stubProg() *model.LoanProgram {
	return &model.LoanProgram{
		ID: "prog-1", Name: "Агро 2025", Rate: 7.5,
		MinAmount: 100_000, MaxAmount: 5_000_000,
		MinTermMonths: 6, MaxTermMonths: 60,
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func withParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func decodeResponse(t *testing.T, body *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(body).Decode(&m); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return m
}

// ── ListPrograms ──────────────────────────────────────────────────────────────

func TestProgramHandler_ListPrograms_Active(t *testing.T) {
	prog := stubProg()
	h := NewProgramHandler(&mockProgramSvc{
		listActiveFn: func(_ context.Context) ([]*model.LoanProgram, error) {
			return []*model.LoanProgram{prog}, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs", nil)
	w := httptest.NewRecorder()
	h.ListPrograms(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
	body := decodeResponse(t, w.Body)
	if body["data"] == nil {
		t.Error("response missing 'data' field")
	}
}

func TestProgramHandler_ListPrograms_All(t *testing.T) {
	var calledAll bool
	h := NewProgramHandler(&mockProgramSvc{
		listAllFn: func(_ context.Context) ([]*model.LoanProgram, error) {
			calledAll = true
			return []*model.LoanProgram{stubProg()}, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs?all=true", nil)
	w := httptest.NewRecorder()
	h.ListPrograms(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
	if !calledAll {
		t.Error("expected ListAll to be called when ?all=true")
	}
}

func TestProgramHandler_ListPrograms_ServiceError(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		listActiveFn: func(_ context.Context) ([]*model.LoanProgram, error) {
			return nil, errors.New("db error")
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs", nil)
	w := httptest.NewRecorder()
	h.ListPrograms(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── GetProgram ────────────────────────────────────────────────────────────────

func TestProgramHandler_GetProgram_Found(t *testing.T) {
	prog := stubProg()
	h := NewProgramHandler(&mockProgramSvc{
		getByIDFn: func(_ context.Context, id string) (*model.LoanProgram, error) {
			if id == prog.ID {
				return prog, nil
			}
			return nil, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs/prog-1", nil)
	r = withParam(r, "id", prog.ID)
	w := httptest.NewRecorder()
	h.GetProgram(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestProgramHandler_GetProgram_NotFound(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, nil },
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs/missing", nil)
	r = withParam(r, "id", "missing")
	w := httptest.NewRecorder()
	h.GetProgram(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestProgramHandler_GetProgram_ServiceError(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) {
			return nil, errors.New("db error")
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/programs/id", nil)
	r = withParam(r, "id", "id")
	w := httptest.NewRecorder()
	h.GetProgram(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── CreateProgram ─────────────────────────────────────────────────────────────

func TestProgramHandler_CreateProgram_Success(t *testing.T) {
	prog := stubProg()
	h := NewProgramHandler(&mockProgramSvc{
		createFn: func(_ context.Context, _ repository.CreateProgramInput) (*model.LoanProgram, error) {
			return prog, nil
		},
	})

	body, _ := json.Marshal(repository.CreateProgramInput{
		Name: "Агро 2025", Rate: 7.5, MinAmount: 100_000, MaxAmount: 5_000_000,
		MinTermMonths: 6, MaxTermMonths: 60,
	})
	r := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProgram(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("want 201, got %d", w.Code)
	}
}

func TestProgramHandler_CreateProgram_BadJSON(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{})

	r := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader([]byte(`not json`)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProgram(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestProgramHandler_CreateProgram_ValidationError(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		createFn: func(_ context.Context, _ repository.CreateProgramInput) (*model.LoanProgram, error) {
			return nil, errors.New("name is required")
		},
	})

	body, _ := json.Marshal(repository.CreateProgramInput{})
	r := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProgram(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
	resp := decodeResponse(t, w.Body)
	if resp["error"] == nil {
		t.Error("response missing 'error' field")
	}
}

// ── UpdateProgram ─────────────────────────────────────────────────────────────

func TestProgramHandler_UpdateProgram_Success(t *testing.T) {
	prog := stubProg()
	h := NewProgramHandler(&mockProgramSvc{
		updateFn: func(_ context.Context, _ string, _ repository.UpdateProgramInput) (*model.LoanProgram, error) {
			return prog, nil
		},
	})

	rate := 9.0
	body, _ := json.Marshal(repository.UpdateProgramInput{Rate: &rate})
	r := httptest.NewRequest(http.MethodPut, "/api/v1/programs/prog-1", bytes.NewReader(body))
	r = withParam(r, "id", "prog-1")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateProgram(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestProgramHandler_UpdateProgram_NotFound(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		updateFn: func(_ context.Context, _ string, _ repository.UpdateProgramInput) (*model.LoanProgram, error) {
			return nil, nil
		},
	})

	body, _ := json.Marshal(repository.UpdateProgramInput{})
	r := httptest.NewRequest(http.MethodPut, "/api/v1/programs/missing", bytes.NewReader(body))
	r = withParam(r, "id", "missing")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateProgram(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestProgramHandler_UpdateProgram_BadJSON(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{})

	r := httptest.NewRequest(http.MethodPut, "/api/v1/programs/id", bytes.NewReader([]byte(`bad`)))
	r = withParam(r, "id", "id")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateProgram(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestProgramHandler_UpdateProgram_ServiceError(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		updateFn: func(_ context.Context, _ string, _ repository.UpdateProgramInput) (*model.LoanProgram, error) {
			return nil, errors.New("db error")
		},
	})

	body, _ := json.Marshal(repository.UpdateProgramInput{})
	r := httptest.NewRequest(http.MethodPut, "/api/v1/programs/id", bytes.NewReader(body))
	r = withParam(r, "id", "id")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateProgram(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── DeactivateProgram ─────────────────────────────────────────────────────────

func TestProgramHandler_DeactivateProgram_Success(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		deactivateFn: func(_ context.Context, _ string) error { return nil },
	})

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/programs/prog-1", nil)
	r = withParam(r, "id", "prog-1")
	w := httptest.NewRecorder()
	h.DeactivateProgram(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestProgramHandler_DeactivateProgram_NotFound(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		deactivateFn: func(_ context.Context, _ string) error { return errors.New("not found") },
	})

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/programs/missing", nil)
	r = withParam(r, "id", "missing")
	w := httptest.NewRecorder()
	h.DeactivateProgram(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestProgramHandler_DeactivateProgram_ServiceError(t *testing.T) {
	h := NewProgramHandler(&mockProgramSvc{
		deactivateFn: func(_ context.Context, _ string) error { return errors.New("db failure") },
	})

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/programs/id", nil)
	r = withParam(r, "id", "id")
	w := httptest.NewRecorder()
	h.DeactivateProgram(w, r)

	// "db failure" != "not found" → 500
	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}
