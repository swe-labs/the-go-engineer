#!/usr/bin/env python3
"""
GitHub Issue and Label Setup Script
Automatically creates labels and issues for The Go Engineer repository

Usage:
    python3 scripts/setup-github-issues.py --token <YOUR_TOKEN> [--dry-run]

Requirements:
    pip3 install requests pyyaml

Author: The Go Engineer Maintainers
License: MIT
"""

import argparse
import json
import os
import sys
from typing import Dict, List, Optional, Tuple
import requests
from datetime import datetime

# ============================================================================
# CONFIGURATION
# ============================================================================

REPO_OWNER = "rasel9t6"
REPO_NAME = "the-go-engineer"
REPO_FULL = f"{REPO_OWNER}/{REPO_NAME}"
API_BASE = "https://api.github.com"

# Color codes (hex without #)
COLORS = {
    "epic": "3E1B6F",
    "feature": "0366D6",
    "bug": "D73A49",
    "documentation": "0075CA",
    "refactor": "CC00FF",
    "performance": "FFA500",
    "testing": "00FF00",
    "chore": "CCCCCC",
    "question": "6F42C1",
    "good_first": "91CA55",
    "critical": "B60205",
    "high": "FF6B6B",
    "medium": "FFA500",
    "low": "FBCA04",
    "later": "E4E669",
    "tiny": "90EE90",
    "small": "9BCD9B",
    "medium_size": "FFD700",
    "large": "FFA500",
    "xlarge": "FF6347",
    "unknown": "CCCCCC",
    "backlog": "FEF2C0",
    "ready": "C6E48B",
    "in_progress": "1E90FF",
    "in_review": "9370DB",
    "blocked": "D73A49",
    "done": "228B22",
    "wontfix": "CCCCCC",
    "help_wanted": "008672",
}

# ============================================================================
# LABELS CONFIGURATION
# ============================================================================

LABELS = [
    # LEVEL 1: Issue Types
    ("epic", COLORS["epic"], "Large feature spanning multiple issues or sections"),
    ("feature", COLORS["feature"], "New lesson, feature, or capability"),
    ("bug", COLORS["bug"], "Something is broken or not working as expected"),
    ("documentation", COLORS["documentation"], "Improvements or additions to documentation"),
    ("refactor", COLORS["refactor"], "Code improvement without changing functionality"),
    ("performance", COLORS["performance"], "Optimization or performance improvement"),
    ("testing", COLORS["testing"], "Add or improve tests"),
    ("chore", COLORS["chore"], "Maintenance tasks, dependencies, tooling"),
    ("question", COLORS["question"], "Questions or discussions"),
    ("good first issue", COLORS["good_first"], "Perfect for new contributors"),
    
    # LEVEL 2: Priority
    ("priority/critical", COLORS["critical"], "Breaks existing functionality or blocks other work"),
    ("priority/high", COLORS["high"], "Important for users/learners, needed soon"),
    ("priority/medium", COLORS["medium"], "Important but not urgent"),
    ("priority/low", COLORS["low"], "Nice to have, can wait"),
    ("priority/later", COLORS["later"], "Parking lot ideas, not planned"),
    
    # LEVEL 3: Effort/Size
    ("size/tiny", COLORS["tiny"], "< 1 hour, story points: 1"),
    ("size/small", COLORS["small"], "1-2 hours, story points: 2"),
    ("size/medium", COLORS["medium_size"], "2-4 hours, story points: 3-5"),
    ("size/large", COLORS["large"], "4-8 hours, story points: 8"),
    ("size/xlarge", COLORS["xlarge"], "8+ hours, story points: 13+"),
    ("size/unknown", COLORS["unknown"], "Not estimated yet"),
    
    # LEVEL 4: Status
    ("status/backlog", COLORS["backlog"], "Created but not yet in active sprint"),
    ("status/ready", COLORS["ready"], "Refined, estimated, ready to work on"),
    ("status/in-progress", COLORS["in_progress"], "Someone is actively working on it"),
    ("status/in-review", COLORS["in_review"], "PR created, waiting for code review"),
    ("status/blocked", COLORS["blocked"], "Cannot proceed, waiting for something else"),
    ("status/done", COLORS["done"], "Completed and merged"),
    ("status/wontfix", COLORS["wontfix"], "Deliberately not fixing"),
    ("help-wanted", COLORS["help_wanted"], "Looking for contributor help"),
    
    # LEVEL 5: Chapters
    ("chapter-13-quality", "3E1B6F", "Quality & Performance"),
    ("chapter-14-architecture", "0366D6", "Application Architecture"),
    ("chapter-15-generation", "0075CA", "Code Generation"),
    
    # LEVEL 6: Components
    ("component/curriculum", "F9D0C4", "curriculum.json and learning paths"),
    ("component/ci-cd", "D4C5F9", "GitHub Actions and build pipeline"),
    ("component/docs", "D9BEE0", "Documentation and READMEs"),
    ("component/lessons", "FCE8C3", "Individual lesson files"),
    ("component/testing", "F3E5AB", "Tests and test infrastructure"),
    ("component/protobuf", "E8F5E9", "Protocol Buffers and gRPC"),
]

# ============================================================================
# ISSUES CONFIGURATION
# ============================================================================

ISSUES = [
    {
        "title": "[EPIC] Complete Chapters 13-15 Implementation",
        "body": """## Overview
Complete implementation of Chapters 13-15 with all lessons, documentation, and curriculum updates.

## Sub-Tasks (Issues)
- [ ] #1: [DOCS] Documentation Updates
- [ ] #2: [FEATURE] PR.3 Memory Profile Lesson
- [ ] #3: [FEATURE] GR.4 Streaming RPC Lesson
- [ ] #4: [FEATURE] GR.5 Interceptors Lesson
- [ ] #5: [FEATURE] CG.2-4 Code Generation Lessons
- [ ] #6: [TASK] Update curriculum.json

## Definition of Done
- [ ] All 6 sub-issues completed and merged
- [ ] All CI checks pass
- [ ] Documentation updated (README, ROADMAP, CHANGELOG)
- [ ] curriculum.json validated with validation script
- [ ] No broken internal links
- [ ] All new lessons follow curriculum style guide
- [ ] Section-level READMEs created for chapters 13-15

## Timeline
- **Start**: Today
- **Target Completion**: 3-5 days
- **Total Effort**: 20-25 hours

## Success Criteria
✓ All lessons compile and run without errors
✓ All follow curriculum style (comments, structure, footer)
✓ Learning paths documented in section READMEs
✓ All changes merged to main branch
✓ Curriculum completion: 15/15 chapters done

## Related
Blocks all other issues in this epic.
""",
        "labels": ["epic", "priority/high", "size/xlarge", "status/ready", "chapter-13-quality", "chapter-14-architecture", "chapter-15-generation", "component/curriculum"],
    },
    {
        "title": "[DOCS] Documentation Updates for Chapters 13-15",
        "body": """## Overview
Update all documentation files to reflect completion of Chapters 13-15.

## Tasks
- [ ] Update ROADMAP.md (mark chapters 13-15 as ✅)
- [ ] Update CHANGELOG.md (add 2026-04-15 entry)
- [ ] Update README.md (add new lessons to projects table)
- [ ] Create 13-quality-and-performance/README.md
- [ ] Create 14-application-architecture/README.md
- [ ] Create 15-code-generation/README.md
- [ ] Verify all internal links and paths
- [ ] Pass CI: make build, make test, make lint

## Success Criteria
✓ All section READMEs created with learning paths
✓ ROADMAP and CHANGELOG updated with new information
✓ README projects table includes all new lessons (PR.3, GR.4, GR.5, CG.2-4)
✓ All links are valid (no broken references)
✓ CI pipeline passes completely

## Related
Depends on: Epic #1
Blocks: Issues #2, #3, #4, #5, #6 (other improvements reference docs)
""",
        "labels": ["documentation", "priority/high", "size/large", "status/ready", "chapter-13-quality", "chapter-14-architecture", "chapter-15-generation", "component/docs"],
    },
    {
        "title": "[FEATURE] Add PR.3: Memory Profile Lesson",
        "body": """## Overview
Add memory profiling lesson to Chapter 13 (Profiling subsection).

## File to Create
`13-quality-and-performance/profiling/2-memory-profile/main.go`

## Content Requirements
- [ ] Header with copyright, section/level info
- [ ] WHAT YOU'LL LEARN section
- [ ] ENGINEERING DEPTH section explaining:
  - Heap vs allocs profiles
  - sync.Pool optimization strategy
  - Memory leak detection workflow
- [ ] Two implementations:
  - slowDecode (excessive allocations)
  - fastDecode (pool-based optimization)
- [ ] printMemStats() helper function
- [ ] Main function demonstrating:
  - Heap profile collection
  - Allocs profile collection
  - Before/after statistics
  - Comparison workflow
- [ ] KEY TAKEAWAY section with key concepts
- [ ] NEXT UP footer (→ PR.4 or similar)

## Testing
- [ ] `go run ./13-quality-and-performance/profiling/2-memory-profile` runs without error
- [ ] `go build ./13-quality-and-performance/profiling/2-memory-profile` compiles
- [ ] Profiles generated: mem_heap.prof, mem_allocs.prof
- [ ] `gofmt` passes (use `make fmt`)
- [ ] `go vet` passes (use `make vet`)
- [ ] `make build` and `make test` pass

## Success Criteria
✓ File compiles and runs successfully
✓ Follows exact curriculum style
✓ Generates profile files for inspection
✓ Demonstrates sync.Pool optimization
✓ All Go best practices followed

## Related
Depends on: #1 (documentation updates)
Part of: Epic #1
""",
        "labels": ["feature", "priority/high", "size/medium", "status/ready", "chapter-13-quality", "component/lessons"],
    },
    {
        "title": "[FEATURE] Add GR.4: Streaming RPC Lesson",
        "body": """## Overview
Add streaming RPC lesson to Chapter 14 (gRPC subsection).

## Files to Create
- `14-application-architecture/grpc/2-streaming/proto/stream.proto`
- `14-application-architecture/grpc/2-streaming/server/main.go`
- `14-application-architecture/grpc/2-streaming/client/main.go`

## Proto Content (stream.proto)
- [ ] Header comments explaining streaming patterns
- [ ] Service definition with 4 RPC types:
  - Add (unary, for reference)
  - CountUp (server streaming)
  - Sum (client streaming)
  - Exchange (bidirectional streaming)
- [ ] Message definitions for each RPC
- [ ] go_package option with correct path

## Server Content (server/main.go)
- [ ] Header with copyright, section info, engineering depth
- [ ] Server struct implementing CalculatorServer
- [ ] Add method (unary RPC)
- [ ] CountUp method (server streaming with context checking)
- [ ] Sum method (accumulating client stream with io.EOF)
- [ ] Exchange method (bidirectional stream)
- [ ] Main function starting server on :50051
- [ ] Comprehensive comments explaining patterns
- [ ] KEY TAKEAWAY section

## Client Content (client/main.go)
- [ ] Header with copyright, section info
- [ ] Demonstrate unary call (Add)
- [ ] Demonstrate server streaming (CountUp with io.EOF)
- [ ] Demonstrate client streaming (Sum with CloseAndRecv)
- [ ] Demonstrate bidirectional (Exchange with CloseSend)
- [ ] KEY TAKEAWAY section
- [ ] Clear instructions for running with server

## Testing
- [ ] `protoc --go_out=. --go-grpc_out=. stream.proto` generates code
- [ ] Generated .pb.go files created
- [ ] Server compiles and runs on :50051
- [ ] Client connects and all patterns work
- [ ] Server and client tested together
- [ ] Follows curriculum style exactly

## Success Criteria
✓ Protobuf compiles successfully
✓ Server and client work together
✓ All 3 streaming patterns demonstrated
✓ io.EOF and context handling shown
✓ Follows curriculum style guide

## Related
Depends on: #1 (documentation updates)
Blocks: #4 (GR.5 interceptors uses proto definitions)
Part of: Epic #1
""",
        "labels": ["feature", "priority/high", "size/large", "status/ready", "chapter-14-architecture", "component/lessons", "component/protobuf"],
    },
    {
        "title": "[FEATURE] Add GR.5: Interceptors Lesson",
        "body": """## Overview
Add gRPC interceptors lesson to Chapter 14 (gRPC subsection).

## File to Create
`14-application-architecture/grpc/3-interceptors/main.go`

## Content Requirements

### Server-Side Interceptors
- [ ] authUnaryInterceptor (checks authorization header)
- [ ] loggingUnaryInterceptor (logs request/response + latency)
- [ ] Metrics struct with unaryInterceptor (call count tracking)
- [ ] authStreamInterceptor (checks auth for streams)

### Server Implementation
- [ ] Server struct implementing CalculatorServer
- [ ] Add method (unary RPC)
- [ ] CountUp method (streaming, simulated with sleep)
- [ ] Server startup with chained interceptors

### Client-Side Interceptors
- [ ] authClientInterceptor (adds authorization header)
- [ ] loggingClientInterceptor (logs calls and latency)

### Main Function
- [ ] Start server in goroutine
- [ ] Create client with interceptors
- [ ] Test successful call (valid token)
- [ ] Test failed call (invalid token with proper error)
- [ ] Print metrics summary
- [ ] KEY TAKEAWAY section

## Teaching Points
- [ ] Explain interceptor pattern (middleware)
- [ ] Show metadata usage (HTTP/2 headers)
- [ ] Demonstrate status codes (Unauthenticated, PermissionDenied)
- [ ] Show chaining with grpc.ChainUnaryInterceptors()

## Testing
- [ ] `go build ./14-application-architecture/grpc/3-interceptors` compiles
- [ ] Server and client run without error
- [ ] Auth rejection works and shows error
- [ ] Metrics print correctly at end
- [ ] Follows curriculum style exactly
- [ ] All CI checks pass

## Success Criteria
✓ File compiles successfully
✓ Server and client work together
✓ Auth validation demonstrated
✓ Metrics collection shown
✓ Follows curriculum style guide

## Related
Depends on: #1 (docs), #3 (GR.4 proto definitions)
Part of: Epic #1
""",
        "labels": ["feature", "priority/high", "size/medium", "status/ready", "chapter-14-architecture", "component/lessons"],
    },
    {
        "title": "[FEATURE] Add CG.2-4: Code Generation Lessons",
        "body": """## Overview
Add three code generation lessons to Chapter 15:
- CG.2: Mockery (automatic mock generation)
- CG.3: Stringer (enum String() generation)
- CG.4: sqlc (type-safe SQL generation)

## Files to Create
1. `15-code-generation/29-mockery/main.go`
2. `15-code-generation/30-stringer/main.go`
3. `15-code-generation/31-sqlc/main.go`

## CG.2: Mockery
- [ ] Header with copyright, section info
- [ ] UserRepository interface with //go:generate mockery directive
- [ ] RealUserRepository implementation
- [ ] UserService using DI with interface
- [ ] Methods: GetUserEmail, UpdateUserEmail
- [ ] Comments explaining interface-based testing
- [ ] Example test code (commented out)
- [ ] KEY TAKEAWAY: Interface design, DI, mock assertions
- [ ] Setup instructions for running go generate

## CG.3: Stringer
- [ ] Header with copyright, section info
- [ ] Multiple enum types with //go:generate stringer directives:
  - Status (Unknown, Pending, Completed, Failed)
  - Priority (Low, Normal, High, Critical)
  - Color (Red, Green, Blue, Yellow, Purple)
- [ ] Task struct using these enums
- [ ] Main demonstrating:
  - Printing enums (shows String() output)
  - Switch statements on enums
  - Formatting in templates
  - Manual parseStatus() helper
- [ ] Generated code example (commented)
- [ ] KEY TAKEAWAY section

## CG.4: sqlc
- [ ] Header explaining sqlc workflow
- [ ] User struct (example of generated type)
- [ ] Query functions (examples of generated code)
- [ ] Main demonstrating usage patterns
- [ ] Comments about setup (sqlc.yaml, schema.sql, queries.sql)
- [ ] Full setup instructions
- [ ] KEY TAKEAWAY section

## Testing
- [ ] All three files compile without errors
- [ ] CG.2 demonstrates interface pattern clearly
- [ ] CG.3 shows enum usage patterns with String()
- [ ] CG.4 explains sqlc workflow and setup
- [ ] All follow curriculum style guide
- [ ] //go:generate directives included
- [ ] Instructions clear for running tools
- [ ] All CI checks pass

## Success Criteria
✓ All three files compile and run (or explain why)
✓ Pattern demonstrations are clear
✓ Teaching flow is logical
✓ Follows curriculum style exactly

## Related
Depends on: #1 (documentation updates)
Part of: Epic #1
""",
        "labels": ["feature", "priority/high", "size/large", "status/ready", "chapter-15-generation", "component/lessons"],
    },
    {
        "title": "[TASK] Update curriculum.json with New Lessons",
        "body": """## Overview
Update curriculum.json to include all new lessons from Chapters 13-15.

## Lessons to Add

### Chapter 13 Additions
- [ ] PR.3: memory-profile
  - Concept: "Heap allocations · sync.Pool · comparative analysis"
  - Requires: PR.1
  - Path: 13-quality-and-performance/profiling/2-memory-profile
  - is_entry: false
  - is_exercise: false

### Chapter 14 Additions
- [ ] GR.4: streaming-rpc
  - Concept: "Server/client/bidirectional streaming · io.EOF detection"
  - Requires: GR.1, GR.2
  - Path: 14-application-architecture/grpc/2-streaming
  - is_entry: false
  - is_exercise: false

- [ ] GR.5: interceptors
  - Concept: "Auth, logging, metrics middleware · metadata headers"
  - Requires: GR.2
  - Path: 14-application-architecture/grpc/3-interceptors
  - is_entry: false
  - is_exercise: false

### Chapter 15 Additions
- [ ] CG.2: mockery-mocking
  - Concept: "Interface mocks · mockery tool · dependency injection"
  - Requires: CG.1
  - Path: 15-code-generation/29-mockery
  - is_entry: false
  - is_exercise: false

- [ ] CG.3: stringer-enum
  - Concept: "Enum String() generation · diagnostics · readability"
  - Requires: CG.1
  - Path: 15-code-generation/30-stringer
  - is_entry: false
  - is_exercise: false

- [ ] CG.4: sqlc-type-safety
  - Concept: "Type-safe SQL · schema-driven generation · zero reflection"
  - Requires: CG.1
  - Path: 15-code-generation/31-sqlc
  - is_entry: false
  - is_exercise: false

## Validation
- [ ] Run: `go run scripts/validate_curriculum.go`
- [ ] Verify no ID collisions
- [ ] Verify all paths exist
- [ ] Verify all prerequisites reference existing lessons
- [ ] All JSON syntax valid (no commas misplaced)

## Testing
- [ ] curriculum.json validates successfully
- [ ] All lessons have correct path, concept, prerequisites
- [ ] IDs follow convention (PR.3, GR.4, GR.5, CG.2, CG.3, CG.4)
- [ ] Section headers updated with new lesson counts
- [ ] `make build` passes

## Success Criteria
✓ curriculum.json is valid JSON
✓ All lessons have correct structure
✓ Validation script passes
✓ No path or ID errors

## Related
Depends on: #1, #2, #3, #4, #5 (all feature branches must be complete)
Part of: Epic #1
""",
        "labels": ["chore", "priority/high", "size/small", "status/ready", "chapter-13-quality", "chapter-14-architecture", "chapter-15-generation", "component/curriculum"],
    },
]

# ============================================================================
# HELPER FUNCTIONS
# ============================================================================

class GitHubAPI:
    """GitHub API client"""
    
    def __init__(self, token: str, dry_run: bool = False):
        self.token = token
        self.dry_run = dry_run
        self.headers = {
            "Authorization": f"token {token}",
            "Accept": "application/vnd.github.v3+json",
            "User-Agent": "The-Go-Engineer-Setup-Script",
        }
        self.session = requests.Session()
        self.session.headers.update(self.headers)
    
    def create_label(self, name: str, color: str, description: str) -> Tuple[bool, str]:
        """Create a label in the repository"""
        url = f"{API_BASE}/repos/{REPO_FULL}/labels"
        data = {
            "name": name,
            "color": color,
            "description": description,
        }
        
        if self.dry_run:
            print(f"  [DRY-RUN] Would create label: {name} (#{color})")
            return True, f"Label '{name}' (dry-run)"
        
        try:
            response = self.session.post(url, json=data)
            
            if response.status_code == 201:
                return True, f"✓ Label '{name}' created"
            elif response.status_code == 422:
                # Label already exists
                return True, f"⚠ Label '{name}' already exists"
            else:
                error = response.json().get("message", "Unknown error")
                return False, f"✗ Failed to create label '{name}': {error}"
        except Exception as e:
            return False, f"✗ Error creating label '{name}': {str(e)}"
    
    def create_issue(self, title: str, body: str, labels: List[str]) -> Tuple[bool, str, Optional[int]]:
        """Create an issue in the repository"""
        url = f"{API_BASE}/repos/{REPO_FULL}/issues"
        data = {
            "title": title,
            "body": body,
            "labels": labels,
        }
        
        if self.dry_run:
            print(f"  [DRY-RUN] Would create issue: {title}")
            return True, f"Issue '{title}' (dry-run)", None
        
        try:
            response = self.session.post(url, json=data)
            
            if response.status_code == 201:
                issue_num = response.json()["number"]
                return True, f"✓ Issue '{title}' created (#{issue_num})", issue_num
            else:
                error = response.json().get("message", "Unknown error")
                return False, f"✗ Failed to create issue: {error}", None
        except Exception as e:
            return False, f"✗ Error creating issue: {str(e)}", None
    
    def verify_connection(self) -> bool:
        """Verify API connection and authentication"""
        url = f"{API_BASE}/user"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                user = response.json()["login"]
                print(f"✓ Authenticated as: {user}")
                return True
            else:
                print(f"✗ Authentication failed: {response.status_code}")
                print(f"  Message: {response.json().get('message', 'Unknown error')}")
                return False
        except Exception as e:
            print(f"✗ Connection error: {str(e)}")
            return False

# ============================================================================
# MAIN EXECUTION
# ============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Setup GitHub labels and issues for The Go Engineer",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 scripts/setup-github-issues.py --token ghp_xxxxxxxxxxxx
  python3 scripts/setup-github-issues.py --token ghp_xxxxxxxxxxxx --dry-run
        """
    )
    
    parser.add_argument(
        "--token",
        required=True,
        help="GitHub Personal Access Token (from https://github.com/settings/tokens)"
    )
    
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would be created without actually creating it"
    )
    
    parser.add_argument(
        "--labels-only",
        action="store_true",
        help="Only create labels, skip issues"
    )
    
    parser.add_argument(
        "--issues-only",
        action="store_true",
        help="Only create issues, skip labels"
    )
    
    args = parser.parse_args()
    
    # Initialize API client
    api = GitHubAPI(args.token, dry_run=args.dry_run)
    
    # Header
    mode = "[DRY-RUN] " if args.dry_run else ""
    print(f"\n{'='*70}")
    print(f"  {mode}GitHub Issue & Label Setup for {REPO_FULL}")
    print(f"{'='*70}\n")
    
    # Verify connection
    print("Step 1: Verifying GitHub API connection...")
    if not api.verify_connection():
        print("\n✗ Failed to authenticate. Check your token and try again.")
        sys.exit(1)
    print()
    
    # Create labels
    if not args.issues_only:
        print("Step 2: Creating labels...")
        print(f"  Total labels to create: {len(LABELS)}\n")
        
        created = 0
        skipped = 0
        
        for name, color, description in LABELS:
            success, message = api.create_label(name, color, description)
            print(f"  {message}")
            if success:
                created += 1
            else:
                skipped += 1
        
        print(f"\n  Summary: {created} created/verified, {skipped} skipped\n")
    
    # Create issues
    if not args.labels_only:
        print("Step 3: Creating issues...")
        print(f"  Total issues to create: {len(ISSUES)}\n")
        
        issue_numbers = {}
        for i, issue_data in enumerate(ISSUES, 1):
            title = issue_data["title"]
            body = issue_data["body"]
            labels = issue_data["labels"]
            
            print(f"  Creating issue #{i}: {title}")
            success, message, issue_num = api.create_issue(title, body, labels)
            print(f"    {message}")
            
            if issue_num:
                issue_numbers[title] = issue_num
        
        print(f"\n  Created {len(issue_numbers)} issue(s)\n")
    
    # Summary
    print("="*70)
    if args.dry_run:
        print("  ✓ Dry-run completed successfully!")
        print("  Run without --dry-run to actually create labels and issues.")
    else:
        print("  ✓ Setup completed successfully!")
        print(f"  Visit: https://github.com/{REPO_FULL}/issues")
    print("="*70)
    print()

if __name__ == "__main__":
    main()
