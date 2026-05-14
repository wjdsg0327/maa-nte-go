package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExecuteTaskReturnsErrorWhenServiceCannotStart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/tasks", ExecuteTask)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"task":"Startup"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected service start failure status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestUpdatePipelineRejectsPathLikeName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.PUT("/pipelines/:name", UpdatePipeline)

	req := httptest.NewRequest(http.MethodPut, "/pipelines/..%5Csecret", bytes.NewBufferString(`{"content":{}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid name status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestUpdatePipelineReturnsErrorWhenReloadFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, "resource", "pipeline"), 0755); err != nil {
		t.Fatalf("create pipeline dir: %v", err)
	}
	filePath := filepath.Join(tempDir, "resource", "pipeline", "Test.json")
	original := []byte(`{"Test":{"recognition":"DirectHit"}}`)
	if err := os.WriteFile(
		filePath,
		original,
		0644,
	); err != nil {
		t.Fatalf("write pipeline: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	router := gin.New()
	router.PUT("/pipelines/:name", UpdatePipeline)

	body := bytes.NewBufferString(`{"content":{"Test":{"recognition":"DirectHit","action":"DoNothing"}}}`)
	req := httptest.NewRequest(http.MethodPut, "/pipelines/Test", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected reload failure status %d, got %d with body %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read pipeline after failed reload: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("expected failed reload to keep original file, got %s", string(got))
	}
}

func TestUpdatePipelineRejectsInvalidRoiBeforeWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, "resource", "pipeline"), 0755); err != nil {
		t.Fatalf("create pipeline dir: %v", err)
	}
	filePath := filepath.Join(tempDir, "resource", "pipeline", "Test.json")
	original := []byte(`{"Test":{"recognition":"DirectHit"}}`)
	if err := os.WriteFile(filePath, original, 0644); err != nil {
		t.Fatalf("write pipeline: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	router := gin.New()
	router.PUT("/pipelines/:name", UpdatePipeline)

	body := bytes.NewBufferString(`{"content":{"Test":{"recognition":"OCR","roi":[983,260.169,388]}}}`)
	req := httptest.NewRequest(http.MethodPut, "/pipelines/Test", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid ROI status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("roi")) {
		t.Fatalf("expected ROI validation message, got %s", rec.Body.String())
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read pipeline after invalid ROI: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("expected invalid ROI to keep original file, got %s", string(got))
	}
}
