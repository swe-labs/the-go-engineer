#!/bin/bash
#
# GitHub Issue and Label Setup Script (Shell Version)
# Alternative to Python version using curl
#
# Usage:
#   bash scripts/setup-github-issues.sh --token <YOUR_TOKEN>
#
# Requirements:
#   curl, jq (json parser)
#

set -e

# ============================================================================
# CONFIGURATION
# ============================================================================

REPO_OWNER="rasel9t6"
REPO_NAME="the-go-engineer"
REPO_FULL="$REPO_OWNER/$REPO_NAME"
API_BASE="https://api.github.com"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============================================================================
# HELP FUNCTION
# ============================================================================

show_help() {
    cat << EOF
GitHub Issue and Label Setup Script

Usage: bash scripts/setup-github-issues.sh --token <TOKEN> [OPTIONS]

Options:
    --token TOKEN           GitHub Personal Access Token (required)
    --dry-run              Show what would be created without actually creating
    --labels-only          Only create labels, skip issues
    --issues-only          Only create issues, skip labels
    -h, --help             Show this help message

Examples:
    bash scripts/setup-github-issues.sh --token ghp_xxxxxxxxxxxx
    bash scripts/setup-github-issues.sh --token ghp_xxxxxxxxxxxx --dry-run

Requirements:
    - curl
    - jq (JSON parser)

EOF
}

# ============================================================================
# UTILITY FUNCTIONS
# ============================================================================

check_dependencies() {
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}✗ curl is not installed${NC}"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}✗ jq is not installed${NC}"
        echo "  Install with: brew install jq (Mac) or apt install jq (Linux)"
        exit 1
    fi
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

log_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# ============================================================================
# API FUNCTIONS
# ============================================================================

verify_connection() {
    log_info "Verifying GitHub API connection..."
    
    local response=$(curl -s -X GET \
        -H "Authorization: token $TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        "$API_BASE/user")
    
    local login=$(echo "$response" | jq -r '.login // empty')
    
    if [ -z "$login" ]; then
        log_error "Authentication failed. Check your token."
        exit 1
    fi
    
    log_success "Authenticated as: $login"
}

create_label() {
    local name=$1
    local color=$2
    local description=$3
    
    if [ "$DRY_RUN" = true ]; then
        log_info "[DRY-RUN] Would create label: $name (#$color)"
        return 0
    fi
    
    local data=$(cat <<EOF
{
  "name": "$name",
  "color": "$color",
  "description": "$description"
}
EOF
)
    
    local response=$(curl -s -X POST \
        -H "Authorization: token $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_BASE/repos/$REPO_FULL/labels")
    
    local status=$(echo "$response" | jq -r '.message // empty')
    
    if echo "$status" | grep -q "already exists"; then
        log_warn "Label '$name' already exists"
        return 0
    fi
    
    if [ -z "$status" ]; then
        log_success "Label '$name' created"
        return 0
    else
        log_error "Failed to create label '$name': $status"
        return 1
    fi
}

# ============================================================================
# MAIN EXECUTION
# ============================================================================

main() {
    # Parse arguments
    TOKEN=""
    DRY_RUN=false
    LABELS_ONLY=false
    ISSUES_ONLY=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --token)
                TOKEN="$2"
                shift 2
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --labels-only)
                LABELS_ONLY=true
                shift
                ;;
            --issues-only)
                ISSUES_ONLY=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Validate token
    if [ -z "$TOKEN" ]; then
        log_error "GitHub token is required"
        show_help
        exit 1
    fi
    
    # Check dependencies
    check_dependencies
    
    # Header
    mode_text=""
    [ "$DRY_RUN" = true ] && mode_text="[DRY-RUN] "
    
    echo ""
    echo "======================================================================"
    echo "  ${mode_text}GitHub Issue & Label Setup for $REPO_FULL"
    echo "======================================================================"
    echo ""
    
    # Verify connection
    verify_connection
    echo ""
    
    # Create labels
    if [ "$ISSUES_ONLY" = false ]; then
        echo "Step 2: Creating labels..."
        echo ""
        
        # Define all labels
        declare -a labels=(
            "epic|3E1B6F|Large feature spanning multiple issues or sections"
            "feature|0366D6|New lesson, feature, or capability"
            "bug|D73A49|Something is broken or not working as expected"
            "documentation|0075CA|Improvements or additions to documentation"
            "refactor|CC00FF|Code improvement without changing functionality"
            "performance|FFA500|Optimization or performance improvement"
            "testing|00FF00|Add or improve tests"
            "chore|CCCCCC|Maintenance tasks, dependencies, tooling"
            "question|6F42C1|Questions or discussions"
            "good first issue|91CA55|Perfect for new contributors"
            "priority/critical|B60205|Breaks existing functionality or blocks other work"
            "priority/high|FF6B6B|Important for users/learners, needed soon"
            "priority/medium|FFA500|Important but not urgent"
            "priority/low|FBCA04|Nice to have, can wait"
            "priority/later|E4E669|Parking lot ideas, not planned"
            "size/tiny|90EE90|< 1 hour, story points: 1"
            "size/small|9BCD9B|1-2 hours, story points: 2"
            "size/medium|FFD700|2-4 hours, story points: 3-5"
            "size/large|FFA500|4-8 hours, story points: 8"
            "size/xlarge|FF6347|8+ hours, story points: 13+"
            "size/unknown|CCCCCC|Not estimated yet"
            "status/backlog|FEF2C0|Created but not yet in active sprint"
            "status/ready|C6E48B|Refined, estimated, ready to work on"
            "status/in-progress|1E90FF|Someone is actively working on it"
            "status/in-review|9370DB|PR created, waiting for code review"
            "status/blocked|D73A49|Cannot proceed, waiting for something else"
            "status/done|228B22|Completed and merged"
            "status/wontfix|CCCCCC|Deliberately not fixing"
            "help-wanted|008672|Looking for contributor help"
            "chapter-13-quality|3E1B6F|Quality & Performance"
            "chapter-14-architecture|0366D6|Application Architecture"
            "chapter-15-generation|0075CA|Code Generation"
            "component/curriculum|F9D0C4|curriculum.json and learning paths"
            "component/ci-cd|D4C5F9|GitHub Actions and build pipeline"
            "component/docs|D9BEE0|Documentation and READMEs"
            "component/lessons|FCE8C3|Individual lesson files"
            "component/testing|F3E5AB|Tests and test infrastructure"
            "component/protobuf|E8F5E9|Protocol Buffers and gRPC"
        )
        
        local created=0
        local skipped=0
        
        for label_entry in "${labels[@]}"; do
            IFS='|' read -r name color desc <<< "$label_entry"
            if create_label "$name" "$color" "$desc"; then
                ((created++))
            else
                ((skipped++))
            fi
        done
        
        echo ""
        echo "  Summary: $created created/verified, $skipped skipped"
        echo ""
    fi
    
    # Summary
    echo "======================================================================"
    if [ "$DRY_RUN" = true ]; then
        log_success "Dry-run completed successfully!"
        echo "  Run without --dry-run to actually create labels and issues."
    else
        log_success "Setup completed successfully!"
        echo "  Visit: https://github.com/$REPO_FULL/issues"
    fi
    echo "======================================================================"
    echo ""
}

main "$@"
