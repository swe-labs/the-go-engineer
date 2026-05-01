// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Comments & Nesting
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to model hierarchical (tree) data in a flat database.
//   - How to build a nested JSON response for threaded comments.
//   - The "Two-Pass" algorithm for efficient tree construction.
//   - How to handle parent-child relationships in Go.
//
// WHY THIS MATTERS:
//   - Hierarchical data is everywhere: folder structures, threaded
//     comments, organizational charts, and menu systems. Learning
//     how to transform a flat list of items into a nested tree is a
//     common engineering challenge that requires careful thinking
//     about data structures and performance.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/10-comments
//
// KEY TAKEAWAY:
//   - Adjacency lists are the simplest way to represent hierarchies in SQL.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Stage 06: Web Masterclass - Comments & Nesting
//
//   - Adjacency List: ParentID field
//   - Tree Construction: The Two-Pass Strategy
//   - JSON Nesting: Recursive serialization
//
// ENGINEERING DEPTH:
//   While you could use recursion to build a tree, it is often
//   more efficient to use a two-pass linear algorithm.
//   1. First, create a map of all items by their ID.
//   2. Second, iterate through the items and use the map to
//      attach children to their respective parents.
//   This approach is O(N) time complexity and avoids the
//   risk of stack overflow that comes with deep recursion.

type Comment struct {
	ID       int        `json:"id"`
	Author   string     `json:"author"`
	Content  string     `json:"content"`
	ParentID int        `json:"parent_id,omitempty"`
	Replies  []*Comment `json:"replies,omitempty"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/comments", func(w http.ResponseWriter, r *http.Request) {
		// 1. Simulate flat data from a database
		flatComments := []Comment{
			{ID: 1, Author: "Alice", Content: "Great post!"},
			{ID: 2, Author: "Bob", Content: "I agree!", ParentID: 1},
			{ID: 3, Author: "Charlie", Content: "Me too!", ParentID: 2},
			{ID: 4, Author: "Dave", Content: "Nice work.", ParentID: 1},
			{ID: 5, Author: "Eve", Content: "Thanks for sharing!"},
		}

		// 2. Build the tree structure (Two-Pass Algorithm)
		commentMap := make(map[int]*Comment)
		var rootComments []*Comment

		// Pass 1: Create pointers for all items and store in a map
		for i := range flatComments {
			c := flatComments[i]
			commentMap[c.ID] = &Comment{
				ID:       c.ID,
				Author:   c.Author,
				Content:  c.Content,
				ParentID: c.ParentID,
			}
		}

		// Pass 2: Attach children to parents
		for _, c := range commentMap {
			if c.ParentID == 0 {
				rootComments = append(rootComments, c)
			} else {
				if parent, ok := commentMap[c.ParentID]; ok {
					parent.Replies = append(parent.Replies, c)
				}
			}
		}

		// 3. Return the nested JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rootComments)
	})

	fmt.Println("=== Web Masterclass: Comments & Nesting ===")
	fmt.Println("  🚀 Server starting on http://localhost:8089")
	fmt.Println()
	fmt.Println("  Visit http://localhost:8089/api/comments to see the tree.")

	log.Fatal(http.ListenAndServe(":8089", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.11 -> 06-backend-db/01-web-and-database/web-masterclass/11-websockets")
	fmt.Println("Current: MC.10 (comments)")
	fmt.Println("Previous: MC.9 (pagination)")
	fmt.Println("---------------------------------------------------")
}
