// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ============================================================================
// Stage 06: Comments
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Comment data model with parent-child relationships
//   - In-memory comment storage (production: use database)
//   - Nested comment retrieval
//   - Thread-safe operations with sync.RWMutex
//
// ENGINEERING DEPTH:
//   Storing hierarchical tree data (Reddit threads, comments) in a flat relational
//   system is famously difficult. We are using the "Adjacency List" pattern:
//   each node simply points to its `ParentID`. When we retrieve the flat list from
//   the database, we perform a 2-Pass algorithm in Go (O(N) time complexity). First
//   pass Maps all items into memory by their ID. The Second pass iterates through
//   all items, checking their `ParentID`, and appends them as a child slice to
//   their respective parent's memory address. This instantly constructs an n-ary
//   tree out of a flat array with zero recursion overhead!
//
// RUN: go run ./06-backend-db/01-web-and-database/web-masterclass/10-comments
// ============================================================================

// Comment represents a single comment on a post.
type Comment struct {
	ID        int        `json:"id"`
	PostID    int        `json:"post_id"`
	ParentID  *int       `json:"parent_id,omitempty"` // nil = top-level comment
	Author    string     `json:"author"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	Replies   []*Comment `json:"replies,omitempty"` // Nested replies
}

// CommentStore is a thread-safe in-memory comment store.
type CommentStore struct {
	mu       sync.RWMutex
	comments []*Comment
	nextID   int
}

func NewCommentStore() *CommentStore {
	return &CommentStore{nextID: 1}
}

// Add creates a new comment.
func (s *CommentStore) Add(postID int, parentID *int, author, content string) *Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := &Comment{
		ID:        s.nextID,
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.comments = append(s.comments, c)
	return c
}

// GetByPost returns all comments for a post, organized as a tree.
func (s *CommentStore) GetByPost(postID int) []*Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Build a map of comment ID → comment for quick lookup
	commentMap := make(map[int]*Comment)
	var topLevel []*Comment

	// First pass: collect all comments for this post and create copies
	for _, c := range s.comments {
		if c.PostID == postID {
			copy := &Comment{
				ID:        c.ID,
				PostID:    c.PostID,
				ParentID:  c.ParentID,
				Author:    c.Author,
				Content:   c.Content,
				CreatedAt: c.CreatedAt,
			}
			commentMap[copy.ID] = copy
		}
	}

	// Second pass: build the tree structure
	for _, c := range commentMap {
		if c.ParentID == nil {
			topLevel = append(topLevel, c)
		} else {
			if parent, ok := commentMap[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, c)
			}
		}
	}

	return topLevel
}

func main() {
	store := NewCommentStore()

	// Seed some example comments
	store.Add(1, nil, "Alice", "Great post!")
	parentID := 1
	store.Add(1, &parentID, "Bob", "Thanks Alice!")
	store.Add(1, &parentID, "Charlie", "I agree with Alice")
	store.Add(1, nil, "Dave", "Very informative")

	mux := http.NewServeMux()

	// List comments for a post (nested tree structure)
	mux.HandleFunc("GET /api/posts/{postID}/comments", func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		comments := store.GetByPost(postID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"post_id":  postID,
			"comments": comments,
		})
	})

	// Add a comment
	mux.HandleFunc("POST /api/posts/{postID}/comments", func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var input struct {
			Author   string `json:"author"`
			Content  string `json:"content"`
			ParentID *int   `json:"parent_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		comment := store.Add(postID, input.ParentID, input.Author, input.Content)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	})

	fmt.Println("Comments API: http://localhost:8080/api/posts/1/comments")
	log.Fatal(http.ListenAndServe(":8080", mux))
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: HC.1 basic GET")
	fmt.Println("   Current: WM.10 (comments)")
	fmt.Println("---------------------------------------------------")
}
