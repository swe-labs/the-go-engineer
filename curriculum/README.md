# The Go Engineer v3 — Curriculum

> Schema: 3.0.0 | Curriculum: 3.0.0-draft.2 | Status: draft | Updated: 2026-05-17

## Architecture Decision

**Decision:** Break the locked v2 curriculum and rebuild around the Report 1 phase-based zero-to-software-engineer path.

- **v2 Policy:** curriculum.v2.json is frozen legacy import data, not the active source of truth.
- **Core Rule:** Keep good lessons, break the old path, optimize for learner progression, proof, documentation quality, and job readiness.

## Quick Stats

| Metric           | Count |
| ---------------- | ----- |
| Modules          | 18    |
| Core Items       | 329   |
| Elective Items   | 20    |
| Total Items      | 349   |
| Projects         | 25    |
| Assessments      | 18    |
| Cross-References | 878   |

## Module Map

| #   | Module                                                               | Phase            | Status  | Items | Portfolio |
| --- | -------------------------------------------------------------------- | ---------------- | ------- | ----- | --------- |
| 0   | **Orientation**                                                      | orientation      | planned | 10    | No        |
| 1   | **Computers, Terminal, Git, and the Web**                            | foundations      | planned | 15    | No        |
| 2   | **Go Setup and Tooling**                                             | tooling          | planned | 14    | No        |
| 3   | **Programming Fundamentals with Go**                                 | go-core          | planned | 23    | No        |
| 4   | **Functions, Errors, and Data Semantics**                            | go-core          | planned | 18    | No        |
| 5   | **Types, Interfaces, Packages, and Modules**                         | engineering-core | planned | 23    | No        |
| 6   | **Testing, Debugging, and Refactoring**                              | engineering-core | planned | 21    | No        |
| 7   | **CLI, Files, JSON, and Configuration**                              | cli-io           | planned | 21    | No        |
| 8   | **HTTP and REST APIs**                                               | backend          | planned | 20    | Yes       |
| 9   | **SQL and PostgreSQL Persistence**                                   | data             | planned | 26    | Yes       |
| 10  | **Authentication, Authorization, and Security**                      | security         | planned | 24    | Yes       |
| 11  | **Lifecycle, Context, and Concurrency**                              | concurrency      | planned | 28    | No        |
| 12  | **Observability, Performance, and Reliability**                      | reliability      | planned | 20    | Yes       |
| 13  | **Architecture and Systems Thinking**                                | architecture     | planned | 17    | No        |
| 14  | **Docker, CI/CD, Linux, and Deployment**                             | delivery         | planned | 19    | Yes       |
| 15  | **Portfolio, Collaboration, Interview Readiness, and Capstone Prep** | career           | planned | 15    | Yes       |
| 16  | **Advanced Electives**                                               | elective         | planned | 20    | No        |
| 17  | **Flagship Opslane**                                                 | flagship         | planned | 15    | Yes       |

## Repository Structure

```
00-orientation
01-computers-terminal-git-web
02-go-setup-tooling
03-programming-fundamentals
04-functions-errors-data-semantics
05-types-interfaces-packages-modules
06-testing-debugging-refactoring
07-cli-files-json-config
08-http-rest-apis
09-sql-postgres-persistence
10-auth-security
11-context-time-concurrency
12-observability-performance-reliability
13-architecture-systems-thinking
14-docker-cicd-linux-deployment
15-portfolio-interview-capstone
17-flagship-opslane
curriculum/path.core.json
curriculum/path.electives.json
curriculum/projects.json
curriculum/assessments.json
curriculum/crossrefs.json
curriculum/migration.v2-to-v3.json
curriculum/legacy/curriculum.v2.json
docs/
scripts/
internal/tools/curriculum/
```

---

# Core Path

## 0. Orientation

> **Phase:** orientation | **Status:** planned | **Required:** True

**Learning Goal:** Explain the learning system, repository workflow, zero-magic contract, assessments, and job-ready target.

| #   | Type    | Slug                                                | Title                                                      | Minutes |
| --- | ------- | --------------------------------------------------- | ---------------------------------------------------------- | ------- |
| 1   | concept | how-to-use-this-repository                          | **How to use this repository**                             | 45      |
| 2   | concept | what-zero-magic-means                               | **What zero magic means**                                  | 45      |
| 3   | concept | how-lessons-exercises-projects-and-checkpoints-work | **How lessons, exercises, projects, and checkpoints work** | 120     |
| 4   | concept | how-to-run-code-in-this-repository                  | **How to run code in this repository**                     | 45      |
| 5   | concept | how-starter-folders-work                            | **How starter folders work**                               | 45      |
| 6   | concept | how-assessments-work                                | **How assessments work**                                   | 45      |
| 7   | concept | how-to-ask-good-debugging-questions                 | **How to ask good debugging questions**                    | 45      |
| 8   | concept | what-job-ready-means                                | **What job-ready means**                                   | 45      |
| 9   | concept | how-to-use-the-roadmap                              | **How to use the roadmap**                                 | 45      |
| 10  | concept | how-to-build-a-portfolio-from-this-curriculum       | **How to build a portfolio from this curriculum**          | 45      |

## 1. Computers, Terminal, Git, and the Web

> **Phase:** foundations | **Status:** planned | **Required:** True

**Learning Goal:** Make the machine, shell, Git, GitHub, and web basics feel safe before Go syntax.

**Prerequisites:** 0. Orientation

| #   | Type    | Slug                                    | Title                                           | Minutes |
| --- | ------- | --------------------------------------- | ----------------------------------------------- | ------- |
| 1   | concept | what-is-a-program                       | **What is a program?**                          | 45      |
| 2   | concept | source-code-executable-and-process      | **Source code, executable, and process**        | 45      |
| 3   | concept | files-bytes-directories-and-paths       | **Files, bytes, directories, and paths**        | 45      |
| 4   | concept | terminal-basics                         | **Terminal basics**                             | 45      |
| 5   | concept | environment-variables                   | **Environment variables**                       | 45      |
| 6   | concept | exit-codes                              | **Exit codes**                                  | 45      |
| 7   | concept | how-the-os-manages-processes            | **How the OS manages processes**                | 45      |
| 8   | concept | memory-preview-stack-vs-heap            | **Memory preview: stack vs heap**               | 45      |
| 9   | concept | git-mental-model                        | **Git mental model**                            | 45      |
| 10  | concept | git-basics-status-add-commit            | **Git basics: status, add, commit**             | 45      |
| 11  | concept | branching-and-merging                   | **Branching and merging**                       | 45      |
| 12  | concept | github-workflow                         | **GitHub workflow**                             | 45      |
| 13  | concept | pull-requests-and-code-review           | **Pull requests and code review**               | 45      |
| 14  | concept | web-preview-client-server-dns-and-ports | **Web preview: client, server, DNS, and ports** | 45      |
| 15  | concept | http-request-and-response-preview       | **HTTP request and response preview**           | 45      |

## 2. Go Setup and Tooling

> **Phase:** tooling | **Status:** planned | **Required:** True

**Learning Goal:** Make the Go toolchain normal from day one.

**Prerequisites:** 1. Computers, Terminal, Git, and the Web

| #   | Type    | Slug                    | Title                       | Minutes |
| --- | ------- | ----------------------- | --------------------------- | ------- |
| 1   | concept | install-and-verify-go   | **Install and verify Go**   | 45      |
| 2   | concept | hello-world             | **Hello World**             | 45      |
| 3   | concept | go-run                  | **go run**                  | 45      |
| 4   | concept | go-build                | **go build**                | 45      |
| 5   | testing | go-test                 | **go test**                 | 45      |
| 6   | concept | gofmt                   | **gofmt**                   | 45      |
| 7   | concept | go-vet                  | **go vet**                  | 45      |
| 8   | concept | go-doc                  | **go doc**                  | 45      |
| 9   | concept | editor-setup-and-gopls  | **Editor setup and gopls**  | 45      |
| 10  | concept | reading-compiler-errors | **Reading compiler errors** | 45      |
| 11  | concept | reading-runtime-errors  | **Reading runtime errors**  | 45      |
| 12  | testing | reading-test-failures   | **Reading test failures**   | 45      |
| 13  | concept | go-module-root-basics   | **Go module root basics**   | 45      |
| 14  | concept | tooling-checklist       | **Tooling checklist**       | 45      |

## 3. Programming Fundamentals with Go

> **Phase:** go-core | **Status:** planned | **Required:** True

**Learning Goal:** Teach core programming through Go without hidden assumptions.

**Prerequisites:** 2. Go Setup and Tooling

| #   | Type    | Slug                       | Title                          | Minutes |
| --- | ------- | -------------------------- | ------------------------------ | ------- |
| 1   | concept | values-and-expressions     | **Values and expressions**     | 45      |
| 2   | concept | variables                  | **Variables**                  | 45      |
| 3   | concept | types                      | **Types**                      | 45      |
| 4   | concept | zero-values                | **Zero values**                | 45      |
| 5   | concept | type-conversions           | **Type conversions**           | 45      |
| 6   | concept | constants                  | **Constants**                  | 45      |
| 7   | concept | iota                       | **iota**                       | 45      |
| 8   | concept | boolean-logic              | **Boolean logic**              | 45      |
| 9   | concept | if-and-else                | **If and else**                | 45      |
| 10  | concept | switch                     | **Switch**                     | 45      |
| 11  | concept | loops                      | **Loops**                      | 45      |
| 12  | concept | range                      | **Range**                      | 45      |
| 13  | concept | arrays                     | **Arrays**                     | 45      |
| 14  | concept | slices                     | **Slices**                     | 45      |
| 15  | concept | slice-length-and-capacity  | **Slice length and capacity**  | 45      |
| 16  | concept | slice-sharing-and-aliasing | **Slice sharing and aliasing** | 45      |
| 17  | concept | maps                       | **Maps**                       | 45      |
| 18  | concept | comma-ok-idiom             | **Comma-ok idiom**             | 45      |
| 19  | concept | strings                    | **Strings**                    | 45      |
| 20  | concept | bytes                      | **Bytes**                      | 45      |
| 21  | concept | runes-and-unicode          | **Runes and Unicode**          | 45      |
| 22  | concept | pointers-as-addresses      | **Pointers as addresses**      | 45      |
| 23  | concept | beginner-pointer-mistakes  | **Beginner pointer mistakes**  | 45      |

## 4. Functions, Errors, and Data Semantics

> **Phase:** go-core | **Status:** planned | **Required:** True

**Learning Goal:** Teach boundaries, failure, cleanup, mutation, and data behavior.

**Prerequisites:** 3. Programming Fundamentals with Go

| #   | Type    | Slug                                | Title                                   | Minutes |
| --- | ------- | ----------------------------------- | --------------------------------------- | ------- |
| 1   | concept | function-basics                     | **Function basics**                     | 45      |
| 2   | concept | parameters                          | **Parameters**                          | 45      |
| 3   | concept | return-values                       | **Return values**                       | 45      |
| 4   | concept | multiple-return-values              | **Multiple return values**              | 45      |
| 5   | concept | scope                               | **Scope**                               | 45      |
| 6   | concept | call-stack-mental-model             | **Call stack mental model**             | 45      |
| 7   | concept | passing-by-value                    | **Passing by value**                    | 45      |
| 8   | concept | pointer-and-value-mutation-behavior | **Pointer and value mutation behavior** | 45      |
| 9   | concept | errors-as-values                    | **Errors as values**                    | 45      |
| 10  | concept | error-wrapping                      | **Error wrapping**                      | 45      |
| 11  | concept | errors-is                           | **errors.Is**                           | 45      |
| 12  | concept | errors-as                           | **errors.As**                           | 45      |
| 13  | concept | validation                          | **Validation**                          | 45      |
| 14  | concept | orchestration                       | **Orchestration**                       | 45      |
| 15  | concept | defer-mechanics                     | **defer mechanics**                     | 45      |
| 16  | concept | defer-for-cleanup                   | **defer for cleanup**                   | 45      |
| 17  | concept | panic-and-recover                   | **panic and recover**                   | 45      |
| 18  | concept | custom-error-types                  | **Custom error types**                  | 45      |

## 5. Types, Interfaces, Packages, and Modules

> **Phase:** engineering-core | **Status:** planned | **Required:** True

**Learning Goal:** Turn scripts into maintainable Go programs.

**Prerequisites:** 4. Functions, Errors, and Data Semantics

| #   | Type    | Slug                   | Title                      | Minutes |
| --- | ------- | ---------------------- | -------------------------- | ------- |
| 1   | concept | structs                | **Structs**                | 45      |
| 2   | concept | methods                | **Methods**                | 45      |
| 3   | concept | pointer-receivers      | **Pointer receivers**      | 45      |
| 4   | concept | value-receivers        | **Value receivers**        | 45      |
| 5   | concept | receiver-sets          | **Receiver sets**          | 45      |
| 6   | concept | composition            | **Composition**            | 45      |
| 7   | concept | embedding              | **Embedding**              | 45      |
| 8   | concept | interfaces             | **Interfaces**             | 45      |
| 9   | concept | small-interfaces       | **Small interfaces**       | 45      |
| 10  | concept | interface-embedding    | **Interface embedding**    | 45      |
| 11  | concept | stringer               | **Stringer**               | 45      |
| 12  | concept | type-assertions        | **Type assertions**        | 45      |
| 13  | concept | type-switches          | **Type switches**          | 45      |
| 14  | concept | nil-interfaces         | **Nil interfaces**         | 45      |
| 15  | concept | package-names          | **Package names**          | 45      |
| 16  | concept | export-rules           | **Export rules**           | 45      |
| 17  | concept | project-layout         | **Project layout**         | 120     |
| 18  | concept | modules                | **Modules**                | 45      |
| 19  | concept | go-mod                 | **go.mod**                 | 45      |
| 20  | concept | go-sum                 | **go.sum**                 | 45      |
| 21  | concept | dependency-management  | **Dependency management**  | 45      |
| 22  | concept | workspaces             | **Workspaces**             | 45      |
| 23  | concept | documentation-comments | **Documentation comments** | 45      |

## 6. Testing, Debugging, and Refactoring

> **Phase:** engineering-core | **Status:** planned | **Required:** True

**Learning Goal:** Move proof, diagnosis, and safe change before serious backend work.

**Prerequisites:** 5. Types, Interfaces, Packages, and Modules

| #   | Type      | Slug                           | Title                              | Minutes |
| --- | --------- | ------------------------------ | ---------------------------------- | ------- |
| 1   | testing   | why-testing-exists             | **Why testing exists**             | 45      |
| 2   | testing   | unit-testing                   | **Unit testing**                   | 45      |
| 3   | testing   | table-driven-tests             | **Table-driven tests**             | 45      |
| 4   | testing   | subtests                       | **Subtests**                       | 45      |
| 5   | testing   | t-cleanup                      | **t.Cleanup**                      | 45      |
| 6   | testing   | test-fixtures                  | **Test fixtures**                  | 45      |
| 7   | testing   | golden-files                   | **Golden files**                   | 45      |
| 8   | testing   | testable-design-with-io-writer | **Testable design with io.Writer** | 45      |
| 9   | testing   | fakes-before-mocks             | **Fakes before mocks**             | 45      |
| 10  | testing   | mocking-tradeoffs              | **Mocking tradeoffs**              | 45      |
| 11  | debugging | debugging-mindset              | **Debugging mindset**              | 45      |
| 12  | testing   | reading-stack-traces           | **Reading stack traces**           | 45      |
| 13  | debugging | debugging-panics               | **Debugging panics**               | 45      |
| 14  | debugging | delve-basics                   | **Delve basics**                   | 45      |
| 15  | debugging | breakpoints                    | **Breakpoints**                    | 45      |
| 16  | testing   | stepping-through-code          | **Stepping through code**          | 45      |
| 17  | testing   | inspecting-variables           | **Inspecting variables**           | 45      |
| 18  | debugging | debugging-tests                | **Debugging tests**                | 45      |
| 19  | testing   | refactoring-safely             | **Refactoring safely**             | 45      |
| 20  | testing   | fuzz-testing                   | **Fuzz testing**                   | 45      |
| 21  | testing   | race-detector-preview          | **Race detector preview**          | 45      |

## 7. CLI, Files, JSON, and Configuration

> **Phase:** cli-io | **Status:** planned | **Required:** True

**Learning Goal:** Build useful local programs before servers.

**Prerequisites:** 6. Testing, Debugging, and Refactoring

| #   | Type    | Slug                                    | Title                                       | Minutes |
| --- | ------- | --------------------------------------- | ------------------------------------------- | ------- |
| 1   | concept | cli-args                                | **CLI args**                                | 45      |
| 2   | concept | flags                                   | **Flags**                                   | 45      |
| 3   | concept | subcommands                             | **Subcommands**                             | 45      |
| 4   | concept | standard-input-and-standard-output      | **Standard input and standard output**      | 45      |
| 5   | concept | io-reader                               | **io.Reader**                               | 45      |
| 6   | concept | io-writer                               | **io.Writer**                               | 45      |
| 7   | concept | files                                   | **Files**                                   | 45      |
| 8   | concept | paths                                   | **Paths**                                   | 45      |
| 9   | concept | directories                             | **Directories**                             | 45      |
| 10  | concept | temp-files                              | **Temp files**                              | 45      |
| 11  | concept | fs-fs                                   | **fs.FS**                                   | 45      |
| 12  | concept | json-marshal                            | **JSON marshal**                            | 45      |
| 13  | concept | json-unmarshal                          | **JSON unmarshal**                          | 45      |
| 14  | concept | json-encoder                            | **JSON encoder**                            | 45      |
| 15  | concept | json-decoder                            | **JSON decoder**                            | 45      |
| 16  | concept | config-files                            | **Config files**                            | 45      |
| 17  | concept | environment-variables-for-configuration | **Environment variables for configuration** | 45      |
| 18  | concept | config-validation                       | **Config validation**                       | 45      |
| 19  | concept | text-templates                          | **Text templates**                          | 45      |
| 20  | concept | base64-as-encoding-not-encryption       | **Base64 as encoding, not encryption**      | 45      |
| 21  | testing | cli-testability                         | **CLI testability**                         | 45      |

## 8. HTTP and REST APIs

> **Phase:** backend | **Status:** planned | **Required:** True

**Learning Goal:** Build the first serious backend service.

**Prerequisites:** 7. CLI, Files, JSON, and Configuration

| #   | Type    | Slug                          | Title                             | Minutes |
| --- | ------- | ----------------------------- | --------------------------------- | ------- |
| 1   | concept | http-as-a-protocol            | **HTTP as a protocol**            | 45      |
| 2   | concept | net-http-basics               | **net/http basics**               | 45      |
| 3   | concept | handler-lifecycle             | **Handler lifecycle**             | 45      |
| 4   | concept | routing                       | **Routing**                       | 45      |
| 5   | concept | request-parsing               | **Request parsing**               | 45      |
| 6   | concept | input-validation              | **Input validation**              | 45      |
| 7   | concept | response-writing              | **Response writing**              | 45      |
| 8   | concept | api-error-shape               | **API error shape**               | 45      |
| 9   | concept | status-codes                  | **Status codes**                  | 45      |
| 10  | concept | middleware-pattern            | **Middleware pattern**            | 45      |
| 11  | concept | error-middleware              | **Error middleware**              | 45      |
| 12  | testing | handler-testing-with-httptest | **Handler testing with httptest** | 45      |
| 13  | concept | body-size-limits              | **Body size limits**              | 45      |
| 14  | concept | server-timeouts               | **Server timeouts**               | 45      |
| 15  | concept | graceful-shutdown             | **Graceful shutdown**             | 45      |
| 16  | concept | health-and-readiness-probes   | **Health and readiness probes**   | 45      |
| 17  | concept | rest-design-principles        | **REST design principles**        | 45      |
| 18  | concept | api-versioning                | **API versioning**                | 45      |
| 19  | concept | pagination-and-filtering      | **Pagination and filtering**      | 45      |
| 20  | concept | openapi-basics                | **OpenAPI basics**                | 45      |

## 9. SQL and PostgreSQL Persistence

> **Phase:** data | **Status:** planned | **Required:** True

**Learning Goal:** Teach real relational backend work.

**Prerequisites:** 8. HTTP and REST APIs

| #   | Type    | Slug                             | Title                                | Minutes |
| --- | ------- | -------------------------------- | ------------------------------------ | ------- |
| 1   | concept | why-databases-exist              | **Why databases exist**              | 45      |
| 2   | concept | relational-modeling              | **Relational modeling**              | 45      |
| 3   | concept | tables-rows-and-columns          | **Tables, rows, and columns**        | 45      |
| 4   | concept | primary-keys                     | **Primary keys**                     | 45      |
| 5   | concept | foreign-keys                     | **Foreign keys**                     | 45      |
| 6   | concept | constraints                      | **Constraints**                      | 45      |
| 7   | concept | sqlite-as-a-local-learning-tool  | **SQLite as a local learning tool**  | 45      |
| 8   | concept | postgresql-as-production-default | **PostgreSQL as production default** | 45      |
| 9   | concept | docker-compose-for-postgres      | **Docker Compose for Postgres**      | 45      |
| 10  | concept | database-sql                     | **database/sql**                     | 45      |
| 11  | concept | sql-db-as-a-pool                 | **sql.DB as a pool**                 | 45      |
| 12  | concept | insert                           | **INSERT**                           | 45      |
| 13  | concept | select                           | **SELECT**                           | 45      |
| 14  | concept | scanning-rows                    | **Scanning rows**                    | 45      |
| 15  | concept | null-handling                    | **Null handling**                    | 45      |
| 16  | concept | prepared-statements              | **Prepared statements**              | 45      |
| 17  | concept | transactions                     | **Transactions**                     | 45      |
| 18  | concept | rollback-discipline              | **Rollback discipline**              | 45      |
| 19  | concept | migrations                       | **Migrations**                       | 45      |
| 20  | concept | joins                            | **Joins**                            | 45      |
| 21  | concept | indexes                          | **Indexes**                          | 45      |
| 22  | concept | explain                          | **EXPLAIN**                          | 45      |
| 23  | concept | query-timeouts-with-context      | **Query timeouts with context**      | 45      |
| 24  | concept | repository-and-service-seam      | **Repository and service seam**      | 45      |
| 25  | testing | integration-tests                | **Integration tests**                | 45      |
| 26  | concept | sqlc                             | **sqlc**                             | 45      |

## 10. Authentication, Authorization, and Security

> **Phase:** security | **Status:** planned | **Required:** True

**Learning Goal:** Make secure backend behavior explicit.

**Prerequisites:** 9. SQL and PostgreSQL Persistence

| #   | Type     | Slug                            | Title                               | Minutes |
| --- | -------- | ------------------------------- | ----------------------------------- | ------- |
| 1   | security | threat-modeling                 | **Threat modeling**                 | 45      |
| 2   | security | trust-boundaries                | **Trust boundaries**                | 45      |
| 3   | security | input-validation-for-security   | **Input validation for security**   | 45      |
| 4   | security | authentication-vs-authorization | **Authentication vs authorization** | 45      |
| 5   | security | password-hashing                | **Password hashing**                | 45      |
| 6   | security | sessions-and-cookies            | **Sessions and cookies**            | 45      |
| 7   | security | jwt-implementation              | **JWT implementation**              | 45      |
| 8   | security | jwt-risks                       | **JWT risks**                       | 45      |
| 9   | security | oauth2-and-oidc-overview        | **OAuth2 and OIDC overview**        | 45      |
| 10  | security | api-keys                        | **API keys**                        | 45      |
| 11  | security | rbac                            | **RBAC**                            | 45      |
| 12  | security | abac-and-policy-checks          | **ABAC and policy checks**          | 45      |
| 13  | security | tenant-isolation-basics         | **Tenant isolation basics**         | 45      |
| 14  | security | sql-injection-prevention        | **SQL injection prevention**        | 45      |
| 15  | security | xss                             | **XSS**                             | 45      |
| 16  | security | csrf                            | **CSRF**                            | 45      |
| 17  | security | cors                            | **CORS**                            | 45      |
| 18  | security | rate-limiting                   | **Rate limiting**                   | 45      |
| 19  | security | tls-and-https                   | **TLS and HTTPS**                   | 45      |
| 20  | security | secrets-management              | **Secrets management**              | 45      |
| 21  | security | dependency-security             | **Dependency security**             | 45      |
| 22  | security | govulncheck                     | **govulncheck**                     | 45      |
| 23  | security | owasp-for-go-apis               | **OWASP for Go APIs**               | 45      |
| 24  | security | secure-logging                  | **Secure logging**                  | 45      |

## 11. Lifecycle, Context, and Concurrency

> **Phase:** concurrency | **Status:** planned | **Required:** True

**Learning Goal:** Teach safe concurrent systems.

**Prerequisites:** 10. Authentication, Authorization, and Security

| #   | Type    | Slug                        | Title                           | Minutes |
| --- | ------- | --------------------------- | ------------------------------- | ------- |
| 1   | concept | time-basics                 | **Time basics**                 | 45      |
| 2   | concept | durations                   | **Durations**                   | 45      |
| 3   | concept | timers                      | **Timers**                      | 45      |
| 4   | concept | tickers                     | **Tickers**                     | 45      |
| 5   | concept | context-basics              | **Context basics**              | 45      |
| 6   | concept | cancellation                | **Cancellation**                | 45      |
| 7   | concept | timeouts                    | **Timeouts**                    | 45      |
| 8   | concept | deadlines                   | **Deadlines**                   | 45      |
| 9   | concept | context-values-with-caveats | **Context values with caveats** | 45      |
| 10  | concept | why-concurrency-exists      | **Why concurrency exists**      | 45      |
| 11  | concept | goroutines                  | **Goroutines**                  | 45      |
| 12  | concept | waitgroups                  | **WaitGroups**                  | 45      |
| 13  | concept | channels                    | **Channels**                    | 45      |
| 14  | concept | buffered-channels           | **Buffered channels**           | 45      |
| 15  | concept | closing-channels            | **Closing channels**            | 45      |
| 16  | concept | channel-ownership           | **Channel ownership**           | 45      |
| 17  | concept | backpressure                | **Backpressure**                | 45      |
| 18  | concept | pipelines                   | **Pipelines**                   | 45      |
| 19  | concept | mutex-and-rwmutex           | **Mutex and RWMutex**           | 45      |
| 20  | concept | atomics-with-caution        | **Atomics with caution**        | 45      |
| 21  | concept | race-conditions             | **Race conditions**             | 45      |
| 22  | concept | race-detector               | **Race detector**               | 45      |
| 23  | concept | goroutine-leaks             | **Goroutine leaks**             | 45      |
| 24  | concept | deadlocks                   | **Deadlocks**                   | 45      |
| 25  | concept | errgroup                    | **errgroup**                    | 45      |
| 26  | concept | bounded-worker-pools        | **Bounded worker pools**        | 45      |
| 27  | concept | retries-and-backoff         | **Retries and backoff**         | 45      |
| 28  | concept | idempotency-preview         | **Idempotency preview**         | 45      |

## 12. Observability, Performance, and Reliability

> **Phase:** reliability | **Status:** planned | **Required:** True

**Learning Goal:** Make services understandable and improvable in production.

**Prerequisites:** 11. Lifecycle, Context, and Concurrency

| #   | Type       | Slug                         | Title                            | Minutes |
| --- | ---------- | ---------------------------- | -------------------------------- | ------- |
| 1   | operations | why-observability-exists     | **Why observability exists**     | 45      |
| 2   | operations | structured-logging-with-slog | **Structured logging with slog** | 45      |
| 3   | operations | request-scoped-logging       | **Request-scoped logging**       | 45      |
| 4   | operations | correlation-ids              | **Correlation IDs**              | 45      |
| 5   | operations | pii-redaction                | **PII redaction**                | 45      |
| 6   | operations | metrics                      | **Metrics**                      | 45      |
| 7   | operations | prometheus                   | **Prometheus**                   | 45      |
| 8   | operations | tracing                      | **Tracing**                      | 45      |
| 9   | operations | opentelemetry                | **OpenTelemetry**                | 45      |
| 10  | operations | pprof                        | **pprof**                        | 45      |
| 11  | operations | cpu-profiling                | **CPU profiling**                | 45      |
| 12  | operations | memory-profiling             | **Memory profiling**             | 45      |
| 13  | operations | benchmarks                   | **Benchmarks**                   | 45      |
| 14  | operations | escape-analysis              | **Escape analysis**              | 45      |
| 15  | operations | memory-layout                | **Memory layout**                | 45      |
| 16  | operations | caching-basics               | **Caching basics**               | 45      |
| 17  | operations | cache-invalidation           | **Cache invalidation**           | 45      |
| 18  | operations | alerting-mindset             | **Alerting mindset**             | 45      |
| 19  | operations | incident-debugging           | **Incident debugging**           | 45      |
| 20  | operations | reliability-review           | **Reliability review**           | 45      |

## 13. Architecture and Systems Thinking

> **Phase:** architecture | **Status:** planned | **Required:** True

**Learning Goal:** Teach architecture as tradeoff reasoning after learners have built real services.

**Prerequisites:** 12. Observability, Performance, and Reliability

| #   | Type         | Slug                                      | Title                                          | Minutes |
| --- | ------------ | ----------------------------------------- | ---------------------------------------------- | ------- |
| 1   | architecture | why-architecture-exists                   | **Why architecture exists**                    | 45      |
| 2   | architecture | package-boundaries                        | **Package boundaries**                         | 45      |
| 3   | architecture | service-layer                             | **Service layer**                              | 45      |
| 4   | architecture | repository-pattern-deep-dive              | **Repository pattern deep dive**               | 45      |
| 5   | architecture | modular-monolith                          | **Modular monolith**                           | 45      |
| 6   | architecture | hexagonal-architecture                    | **Hexagonal architecture**                     | 45      |
| 7   | architecture | domain-modeling                           | **Domain modeling**                            | 45      |
| 8   | architecture | invariants                                | **Invariants**                                 | 45      |
| 9   | architecture | event-driven-basics                       | **Event-driven basics**                        | 45      |
| 10  | architecture | queues                                    | **Queues**                                     | 45      |
| 11  | architecture | retries-and-idempotent-consumers          | **Retries and idempotent consumers**           | 45      |
| 12  | architecture | payment-workflow-design                   | **Payment workflow design**                    | 45      |
| 13  | architecture | multi-tenancy-architecture                | **Multi-tenancy architecture**                 | 45      |
| 14  | architecture | caching-architecture                      | **Caching architecture**                       | 45      |
| 15  | architecture | when-to-split-services                    | **When to split services**                     | 45      |
| 16  | architecture | microservices-as-a-tradeoff-not-a-default | **Microservices as a tradeoff, not a default** | 45      |
| 17  | architecture | architecture-decision-records             | **Architecture decision records**              | 45      |

## 14. Docker, CI/CD, Linux, and Deployment

> **Phase:** delivery | **Status:** planned | **Required:** True

**Learning Goal:** Make software shippable and repeatable.

**Prerequisites:** 13. Architecture and Systems Thinking

| #   | Type       | Slug                     | Title                        | Minutes |
| --- | ---------- | ------------------------ | ---------------------------- | ------- |
| 1   | operations | linux-processes          | **Linux processes**          | 45      |
| 2   | operations | signals                  | **Signals**                  | 45      |
| 3   | operations | logs                     | **Logs**                     | 45      |
| 4   | operations | file-permissions         | **File permissions**         | 45      |
| 5   | operations | networking-basics        | **Networking basics**        | 45      |
| 6   | operations | docker-basics            | **Docker basics**            | 45      |
| 7   | operations | docker-images-and-layers | **Docker images and layers** | 45      |
| 8   | operations | multi-stage-builds       | **Multi-stage builds**       | 45      |
| 9   | operations | docker-compose           | **Docker Compose**           | 45      |
| 10  | operations | container-health-checks  | **Container health checks**  | 45      |
| 11  | operations | github-actions           | **GitHub Actions**           | 45      |
| 12  | operations | ci-pipeline              | **CI pipeline**              | 45      |
| 13  | operations | release-artifacts        | **Release artifacts**        | 45      |
| 14  | operations | config-in-deployment     | **Config in deployment**     | 45      |
| 15  | operations | secrets-in-deployment    | **Secrets in deployment**    | 45      |
| 16  | operations | vulnerability-scanning   | **Vulnerability scanning**   | 45      |
| 17  | operations | one-deployment-target    | **One deployment target**    | 45      |
| 18  | operations | rollback-basics          | **Rollback basics**          | 45      |
| 19  | operations | deployment-runbook       | **Deployment runbook**       | 45      |

## 15. Portfolio, Collaboration, Interview Readiness, and Capstone Prep

> **Phase:** career | **Status:** planned | **Required:** True

**Learning Goal:** Convert competence into employability.

**Prerequisites:** 14. Docker, CI/CD, Linux, and Deployment

| #   | Type    | Slug                           | Title                              | Minutes |
| --- | ------- | ------------------------------ | ---------------------------------- | ------- |
| 1   | career  | readme-writing                 | **README writing**                 | 45      |
| 2   | career  | api-docs                       | **API docs**                       | 45      |
| 3   | career  | code-review-workflow           | **Code review workflow**           | 45      |
| 4   | career  | pull-request-etiquette         | **Pull request etiquette**         | 45      |
| 5   | career  | portfolio-packaging            | **Portfolio packaging**            | 45      |
| 6   | concept | writing-project-narratives     | **Writing project narratives**     | 120     |
| 7   | concept | resume-project-bullets         | **Resume project bullets**         | 120     |
| 8   | career  | backend-interview-fundamentals | **Backend interview fundamentals** | 45      |
| 9   | career  | go-interview-topics            | **Go interview topics**            | 45      |
| 10  | career  | sql-interview-topics           | **SQL interview topics**           | 45      |
| 11  | career  | debugging-interview-scenarios  | **Debugging interview scenarios**  | 45      |
| 12  | career  | system-design-basics           | **System design basics**           | 45      |
| 13  | career  | architecture-defense           | **Architecture defense**           | 45      |
| 14  | career  | capstone-planning              | **Capstone planning**              | 45      |
| 15  | career  | demo-preparation               | **Demo preparation**               | 45      |

## 17. Flagship Opslane

> **Phase:** flagship | **Status:** planned | **Required:** True

**Learning Goal:** Final graduation proof through a portfolio-grade backend system.

**Prerequisites:** 15. Portfolio, Collaboration, Interview Readiness, and Capstone Prep

| #   | Type     | Slug                                 | Title                                    | Minutes |
| --- | -------- | ------------------------------------ | ---------------------------------------- | ------- |
| 1   | capstone | opslane-foundation-and-configuration | **Opslane foundation and configuration** | 120     |
| 2   | capstone | opslane-database-and-models          | **Opslane database and models**          | 120     |
| 3   | capstone | opslane-authentication               | **Opslane authentication**               | 120     |
| 4   | capstone | opslane-authorization                | **Opslane authorization**                | 120     |
| 5   | capstone | opslane-tenant-isolation             | **Opslane tenant isolation**             | 120     |
| 6   | capstone | opslane-http-api-layer               | **Opslane HTTP API layer**               | 120     |
| 7   | capstone | opslane-order-processing             | **Opslane order processing**             | 120     |
| 8   | capstone | opslane-payment-pipeline             | **Opslane payment pipeline**             | 120     |
| 9   | capstone | opslane-event-bus                    | **Opslane event bus**                    | 120     |
| 10  | capstone | opslane-worker-pools                 | **Opslane worker pools**                 | 120     |
| 11  | capstone | opslane-caching-layer                | **Opslane caching layer**                | 120     |
| 12  | capstone | opslane-observability                | **Opslane observability**                | 120     |
| 13  | capstone | opslane-graceful-shutdown            | **Opslane graceful shutdown**            | 120     |
| 14  | capstone | opslane-deployment                   | **Opslane deployment**                   | 120     |
| 15  | capstone | opslane-portfolio-defense            | **Opslane portfolio defense**            | 120     |

---

# Advanced Electives

## 16. Advanced Electives

> **Phase:** elective | **Status:** planned | **Required:** False

**Learning Goal:** Keep advanced Go and backend topics without blocking the core learner path.

| #   | Type      | Slug                        | Title                           | Minutes |
| --- | --------- | --------------------------- | ------------------------------- | ------- |
| 1   | reference | generics                    | **Generics**                    | 45      |
| 2   | reference | complex-generic-constraints | **Complex generic constraints** | 45      |
| 3   | reference | generic-data-structures     | **Generic data structures**     | 45      |
| 4   | reference | functional-options          | **Functional options**          | 45      |
| 5   | reference | method-values               | **Method values**               | 45      |
| 6   | reference | build-tags                  | **Build tags**                  | 45      |
| 7   | reference | go-generate                 | **go generate**                 | 45      |
| 8   | reference | mockery                     | **Mockery**                     | 45      |
| 9   | reference | protobuf                    | **Protobuf**                    | 45      |
| 10  | reference | grpc-fundamentals           | **gRPC fundamentals**           | 45      |
| 11  | reference | grpc-streaming              | **gRPC streaming**              | 45      |
| 12  | reference | grpc-interceptors           | **gRPC interceptors**           | 45      |
| 13  | reference | rest-vs-grpc-tradeoffs      | **REST vs gRPC tradeoffs**      | 45      |
| 14  | reference | sync-map                    | **sync.Map**                    | 45      |
| 15  | reference | sync-pool                   | **sync.Pool**                   | 45      |
| 16  | reference | cqrs                        | **CQRS**                        | 45      |
| 17  | reference | feature-flags               | **Feature flags**               | 45      |
| 18  | reference | blue-green-deployments      | **Blue/green deployments**      | 45      |
| 19  | reference | microservice-splitting      | **Microservice splitting**      | 45      |
| 20  | reference | deeper-distributed-systems  | **Deeper distributed systems**  | 45      |

---

# Projects

| #   | Title                             | Type              | Stage            | Difficulty   |
| --- | --------------------------------- | ----------------- | ---------------- | ------------ |
| 1   | **Shell + Git Lab**               | mini-project      | foundations      | foundation   |
| 2   | **Tooling Failure Lab**           | mini-project      | tooling          | foundation   |
| 3   | **Pricing Checkout**              | mini-project      | go-core          | core         |
| 4   | **Contact Directory**             | mini-project      | go-core          | core         |
| 5   | **Order Summary Refactor**        | mini-project      | go-core          | core         |
| 6   | **Bank Account**                  | mini-project      | engineering-core | core         |
| 7   | **Payroll Processor**             | mini-project      | engineering-core | core         |
| 8   | **Bug Hunt Pack**                 | mini-project      | engineering-core | core         |
| 9   | **Refactor Under Tests**          | mini-project      | engineering-core | core         |
| 10  | **File Organizer**                | mini-project      | cli-io           | core         |
| 11  | **Log Search**                    | mini-project      | cli-io           | core         |
| 12  | **Config Parser**                 | mini-project      | cli-io           | core         |
| 13  | **Book Catalog REST API**         | portfolio-project | backend          | intermediate |
| 14  | **Task API with Postgres**        | portfolio-project | data             | intermediate |
| 15  | **Secure Authenticated API**      | portfolio-project | security         | intermediate |
| 16  | **Concurrent Downloader**         | mini-project      | concurrency      | intermediate |
| 17  | **URL Health Checker**            | mini-project      | concurrency      | intermediate |
| 18  | **Bounded Worker Pool**           | mini-project      | concurrency      | intermediate |
| 19  | **PII Redactor**                  | portfolio-project | reliability      | intermediate |
| 20  | **Observability Retrofit**        | portfolio-project | reliability      | advanced     |
| 21  | **Profile and Fix Lab**           | portfolio-project | reliability      | advanced     |
| 22  | **Modular Refactor Review**       | mini-project      | architecture     | advanced     |
| 23  | **Dockerized Service with CI/CD** | portfolio-project | delivery         | advanced     |
| 24  | **Portfolio Review Pack**         | portfolio-project | career           | advanced     |
| 25  | **Opslane Flagship**              | flagship          | flagship         | flagship     |

---

# Assessments

| ID                   | Title                                                                                   | Type             | Passing Score |
| -------------------- | --------------------------------------------------------------------------------------- | ---------------- | ------------- |
| assessment-module-00 | Module 00 Assessment — Orientation                                                      | project-review   | 80%           |
| assessment-module-01 | Module 01 Assessment — Computers, Terminal, Git, and the Web                            | project-review   | 80%           |
| assessment-module-02 | Module 02 Assessment — Go Setup and Tooling                                             | project-review   | 80%           |
| assessment-module-03 | Module 03 Assessment — Programming Fundamentals with Go                                 | project-review   | 80%           |
| assessment-module-04 | Module 04 Assessment — Functions, Errors, and Data Semantics                            | project-review   | 80%           |
| assessment-module-05 | Module 05 Assessment — Types, Interfaces, Packages, and Modules                         | project-review   | 80%           |
| assessment-module-06 | Module 06 Assessment — Testing, Debugging, and Refactoring                              | project-review   | 80%           |
| assessment-module-07 | Module 07 Assessment — CLI, Files, JSON, and Configuration                              | project-review   | 80%           |
| assessment-module-08 | Module 08 Assessment — HTTP and REST APIs                                               | project-review   | 80%           |
| assessment-module-09 | Module 09 Assessment — SQL and PostgreSQL Persistence                                   | project-review   | 80%           |
| assessment-module-10 | Module 10 Assessment — Authentication, Authorization, and Security                      | project-review   | 80%           |
| assessment-module-11 | Module 11 Assessment — Lifecycle, Context, and Concurrency                              | project-review   | 80%           |
| assessment-module-12 | Module 12 Assessment — Observability, Performance, and Reliability                      | project-review   | 80%           |
| assessment-module-13 | Module 13 Assessment — Architecture and Systems Thinking                                | project-review   | 80%           |
| assessment-module-14 | Module 14 Assessment — Docker, CI/CD, Linux, and Deployment                             | project-review   | 80%           |
| assessment-module-15 | Module 15 Assessment — Portfolio, Collaboration, Interview Readiness, and Capstone Prep | portfolio-review | 80%           |
| assessment-module-16 | Module 16 Assessment — Advanced Electives                                               | project-review   | 80%           |
| assessment-module-17 | Module 17 Assessment — Flagship Opslane                                                 | capstone-defense | 80%           |

---

# Cross-References

Total cross-references: **878**

### Style Rules

| Relation      | Template                                                     |
| ------------- | ------------------------------------------------------------ |
| builds_on     | Builds on: Module X.Y â€” Lesson Title                       |
| preview_only  | Preview only: Explained fully in Module X.Y â€” Lesson Title |
| reinforced_in | Reinforced in: Project X â€” Project Title                   |
| related       | Related: Module X.Y â€” Lesson Title                         |

### Relation Distribution

| Relation      | Count |
| ------------- | ----- |
| builds_on     | 348   |
| preview_only  | 6     |
| reinforced_in | 524   |

---

# Migration: v2 → v3

**Source Schema:** 2

**Target Schema:** 3.0.0

**Default Policy:** Every legacy item must be mapped to keep, rewrite, split, merge, move-to-elective, replace, or archive. Every migrated README must also receive a legacy_readme_action: reuse, rewrite, split, merge, or archive.

**Frozen Source:** curriculum/legacy/curriculum.v2.json

### High-Confidence Moves

| Legacy Pattern        | Outcome          | Target    | Reason                                                                    |
| --------------------- | ---------------- | --------- | ------------------------------------------------------------------------- |
| API.4-API.9           | move-to-elective | module-16 | Protobuf/gRPC is valuable but should not block the core beginner path.    |
| TI.14-TI.15           | move-to-elective | module-16 | Complex generics and generic data structures are advanced.                |
| TI.12                 | move-to-elective | module-16 | Functional options are useful but not core beginner material.             |
| CP.3                  | move-to-elective | module-16 | sync.Pool is performance-specialized.                                     |
| SY.2 sync.Map portion | move-to-elective | module-16 | sync.Map is specialized and should follow normal map+mutex reasoning.     |
| ARCH.7                | move-to-elective | module-16 | CQRS is advanced architecture.                                            |
| DEPLOY.2              | move-to-elective | module-16 | Blue/green deployment is advanced deployment strategy.                    |
| TE.1-TE.11            | rewrite          | module-06 | Testing and debugging must move before backend work.                      |
| DB.\*                 | rewrite          | module-09 | Database path must become SQL/Postgres-first with migrations and EXPLAIN. |
| SEC.\*                | rewrite          | module-10 | Security needs authn/authz strategy and secure API practice.              |

### Quality Gates

- schema-valid
- graph-valid
- no-orphan-core-items
- all-prerequisites-exist
- no-circular-dependencies
- item-orders-sequential
- phase-sequence-valid
- no-null-arrays
- crossref-targets-exist
- crossref-targets-in-same-bundle
- module-prerequisites-exist
- item-prerequisites-exist
- project-reinforces-items-exist
- assessment-targets-exist
- project-assessment-exists
- no-duplicate-id
- no-duplicate-slug
- no-unused-crossref-item
- item-has-zero-magic
- item-has-proof
- item-has-content-contract
- item-has-verification
- item-has-files
- item-files-path-valid
- assessment-retake-policy-present

### Validation Commands

| Command          | Script                                              |
| ---------------- | --------------------------------------------------- |
| validate_all     | go run ./internal/tools/curriculum validate-all     |
| validate_graph   | go run ./internal/tools/curriculum validate-graph   |
| validate_lessons | go run ./internal/tools/curriculum validate-lessons |
| validate_readmes | go run ./internal/tools/curriculum validate-readmes |
| validate_schema  | go run ./internal/tools/curriculum validate-schema  |

---

_Generated from 3.0.0-draft.2 | Schema 3.0.0 | 2026-05-17_
