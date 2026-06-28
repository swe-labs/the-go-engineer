#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="${1:-project}"
PASS=0
FAIL=0
TOTAL=14

check() {
    local name="$1"
    if shift; "$@"; then
        PASS=$((PASS + 1))
        printf "  PASS  %s\n" "$name"
    else
        PASS=$((PASS - 1))
        FAIL=$((FAIL + 1))
        printf "  FAIL  %s\n" "$name"
    fi
}

printf "\nShell + Git Lab Verification\n\n"

check "project directory exists" test -d "$PROJECT_DIR"
check "src subdirectory exists" test -d "$PROJECT_DIR/src"
check "docs subdirectory exists" test -d "$PROJECT_DIR/docs"
check "tests subdirectory exists" test -d "$PROJECT_DIR/tests"
check "README.md exists" test -f "$PROJECT_DIR/README.md"
check "src/main.go exists" test -f "$PROJECT_DIR/src/main.go"
check "docs/design.md exists" test -f "$PROJECT_DIR/docs/design.md"

# Git checks
pushd "$PROJECT_DIR" > /dev/null
check ".git directory exists" test -d ".git"
check "at least 3 commits exist" bash -c '[[ $(git log --oneline 2>/dev/null | wc -l) -ge 3 ]]'
check "feature branch merged" bash -c 'git branch --merged 2>/dev/null | grep -q "feature"'
check "git log shows commit history" bash -c '[[ $(git log --oneline --graph 2>/dev/null | wc -l) -ge 3 ]]'
check "README.md is tracked" bash -c 'git ls-files README.md 2>/dev/null | grep -q .'
check "no unresolved merge conflicts" bash -c '! git grep -l ">>>>>>>" -- "*.md" "*.go" 2>/dev/null | grep -q .'
popd > /dev/null

printf "\nResults: %d/%d passed, %d failed\n" "$PASS" "$TOTAL" "$FAIL"
if [ "$FAIL" -eq 0 ]; then
    printf "All checks passed!\n"
else
    printf "Some checks failed. Review the output above.\n"
fi
