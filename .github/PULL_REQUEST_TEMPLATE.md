## Description

<!-- Please include a summary of the change and which issue is fixed (if applicable). -->
<!-- Include relevant motivation and context. -->

Fixes # (issue)

## Type of change

- [ ] 🐛 Bug fix (typo, inaccurate comment, broken code)
- [ ] ✨ New feature (new exercise, new section)
- [ ] 📝 Documentation update (README or CONTRIBUTING)
- [ ] ♻️ Refactoring (improving existing code without changing behavior)
- [ ] 🔧 Build system or CI/CD improvement
- [ ] 🧪 Test improvements

## Pre-Submission Checklist

### Code Quality
- [ ] Code formatted with `gofmt` (run `make fmt`)
- [ ] Passes `go vet` checks (run `make vet`)
- [ ] Builds without errors (run `make build`)
- [ ] All tests pass (run `make test`)
- [ ] Race condition check passes (run `make test-race`)

### Documentation
- [ ] File has proper header comment (if applicable)
- [ ] Complex logic is documented with inline comments
- [ ] Exported functions have doc comments
- [ ] No TODO comments without tracking issues

### Tests
- [ ] New code includes tests
- [ ] Tests use table-driven pattern (when appropriate)
- [ ] Coverage is appropriate (aim for > 75%)
- [ ] Tests are named descriptively

### For Lesson Additions
- [ ] Follows the standard file template from [CODE-STANDARDS.md](../../CODE-STANDARDS.md)
- [ ] Includes learning objectives in header
- [ ] Includes "ENGINEERING DEPTH" section explaining production relevance
- [ ] Includes a `_starter/` stub if exercise
- [ ] Curriculum mapping added to `curriculum.json`
- [ ] Navigation footer added (next lesson reference)

### Repository Maintenance
- [ ] `go.mod` and `go.sum` are updated (`go mod tidy`)
- [ ] No unnecessary dependencies added
- [ ] No large files (> 10MB) committed
- [ ] No credentials or sensitive data committed
- [ ] Changelog updated if user-facing change

## Testing Instructions

<!-- Describe how to test this change. Include specific commands if applicable. -->

```bash
# For bug fixes
go test -run TestName ./path/to/package

# For new lessons
go run ./path/to/lesson

# For curriculum changes
go run ./scripts/validate_curriculum.go
```

## Performance Impact

<!-- For changes that might affect performance, include benchmarks. -->

- [ ] No performance impact
- [ ] Improvement: [describe]
- [ ] Regression: [justify]

If applicable, run benchmarks:
```bash
make bench
```

## Reviewer Guidance

<!-- Add any specific things you want the reviewer to look at, concerns, or questions. -->

## Screenshots/Output

<!-- If visual changes, include screenshots. If code output, show terminal output. -->

---

**Note**: Maintainers will verify:
- All CI checks pass
- Code style compliance
- Test coverage
- Curriculum consistency

