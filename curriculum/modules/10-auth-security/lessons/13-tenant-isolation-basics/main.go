package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type tenantKeyType string

const tenantKey tenantKeyType = "tenant_id"

type Document struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type TenantStore struct {
	mu        sync.RWMutex
	documents []Document
}

func NewTenantStore() *TenantStore {
	return &TenantStore{}
}

func (ts *TenantStore) AddDocument(ctx context.Context, title, content string) (*Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 16)
	rand.Read(b)
	doc := Document{
		ID: hex.EncodeToString(b), TenantID: tenantID,
		Title: title, Content: content, CreatedAt: time.Now(),
	}
	ts.mu.Lock()
	ts.documents = append(ts.documents, doc)
	ts.mu.Unlock()
	return &doc, nil
}

func (ts *TenantStore) GetDocuments(ctx context.Context) ([]Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	var result []Document
	for _, d := range ts.documents {
		if d.TenantID == tenantID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (ts *TenantStore) GetDocument(ctx context.Context, docID string) (*Document, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	for _, d := range ts.documents {
		if d.ID == docID && d.TenantID == tenantID {
			return &d, nil
		}
	}
	return nil, errors.New("document not found")
}

func (ts *TenantStore) GetDocumentsNoFilter(ctx context.Context) ([]Document, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	out := make([]Document, len(ts.documents))
	copy(out, ts.documents)
	return out, nil
}

func getTenantID(ctx context.Context) (string, error) {
	tid, ok := ctx.Value(tenantKey).(string)
	if !ok || tid == "" {
		return "", errors.New("tenant not found in context")
	}
	return tid, nil
}

func CreateDocument(ctx context.Context, title, content string) (*Document, error) {
	return NewTenantStore().AddDocument(ctx, title, content)
}

func GetDocuments(ctx context.Context) ([]Document, error) {
	return NewTenantStore().GetDocuments(ctx)
}

func main() {
	store := NewTenantStore()

	ctx1 := context.WithValue(context.Background(), tenantKey, "acme-corp")
	ctx2 := context.WithValue(context.Background(), tenantKey, "globex-inc")

	store.AddDocument(ctx1, "Acme Strategy 2026", "Confidential Acme plans")
	store.AddDocument(ctx1, "Acme Product Roadmap", "Q3-Q4 features")
	store.AddDocument(ctx2, "Globex Financial Report", "Q2 earnings")

	fmt.Println("=== Correct tenant isolation ===")
	docs1, _ := store.GetDocuments(ctx1)
	fmt.Printf("Acme-Corp documents (%d):\n", len(docs1))
	for _, d := range docs1 {
		fmt.Printf("  - %s\n", d.Title)
	}

	docs2, _ := store.GetDocuments(ctx2)
	fmt.Printf("Globex-Inc documents (%d):\n", len(docs2))
	for _, d := range docs2 {
		fmt.Printf("  - %s\n", d.Title)
	}

	fmt.Println("\n=== Cross-tenant leak (bug without tenant filter) ===")
	all, _ := store.GetDocumentsNoFilter(context.Background())
	fmt.Printf("All documents without filter (%d items):\n", len(all))
	for _, d := range all {
		fmt.Printf("  - %s (tenant: %s)\n", d.Title, d.TenantID)
	}
	fmt.Println("SECURITY ISSUE: Any tenant can access any document!")

	// HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/api/documents", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			http.Error(w, `{"error":"tenant not identified"}`, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), tenantKey, tenantID)
		switch r.Method {
		case "GET":
			docs, err := store.GetDocuments(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(docs)
		case "POST":
			var req struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			doc, err := store.AddDocument(ctx, req.Title, req.Content)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(doc)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("\nTenant isolation server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
