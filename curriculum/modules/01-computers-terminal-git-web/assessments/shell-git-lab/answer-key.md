# Shell + Git Lab — Answer Key

## Multiple choice (1 point each)

**1. a) `mkdir projects`**

**2. b) A temporary area where changes are collected before committing**

**3. c) The current state of the working directory and staging area**

**4. a) `git branch feature-login`**

**5. b) 0**

**6. b) Converting domain names to IP addresses**

**7. c) GET**

**8. c) Not Found**

**9. a) Merge creates a new commit with two parents; rebase rewrites history by replaying commits**

**10. b) The output of the first command becomes the input of the second command**

## Short answer (2 points each)

**11.** The three states are:
- **Modified**: A file has been changed in the working directory but not yet staged.
- **Staged**: The modified file has been marked for inclusion in the next commit via `git add`.
- **Committed**: The staged changes have been saved to the local repository via `git commit`.

The transition flow is: Modified → `git add` → Staged → `git commit` → Committed.

*Grading: 1 point for naming all three states, 1 point for describing the transition correctly.*

**12.** The sequence is:
1. The browser checks its DNS cache for `example.com`.
2. If not cached, it queries a DNS resolver (recursively through root, TLD, and authoritative servers) to resolve `example.com` to an IP address.
3. The browser opens a TCP connection to that IP address on port 443 (HTTPS) via a three-way handshake (SYN, SYN-ACK, ACK).
4. A TLS handshake establishes encryption.
5. The browser sends an HTTP GET request for the root path `/`.
6. The server processes the request and sends an HTTP response back.
7. The browser renders the HTML content.

*Grading: 1 point for mentioning DNS resolution, 0.5 for TCP connection, 0.5 for HTTP request/response.*

**13.** To move a commit from the wrong branch to the correct branch:
1. Note the commit hash: `git log --oneline -1`
2. Switch to the correct branch: `git checkout correct-branch`
3. Cherry-pick the commit: `git cherry-pick <commit-hash>`
4. Switch back to the wrong branch: `git checkout wrong-branch`
5. Reset to remove the commit: `git reset HEAD~1`

*Grading: 1 point for cherry-pick approach, 1 point for resetting the wrong branch.*

**14.** An HTTP GET request looks like:
```
GET /index.html HTTP/1.1
Host: example.com
User-Agent: curl/8.0
Accept: */*

```
Components:
- **Method**: GET
- **Path**: /index.html
- **Protocol version**: HTTP/1.1
- **Headers**: Host, User-Agent, Accept (each on a separate line)
- **Empty line**: Separates headers from the body (no body for GET)

*Grading: 1 point for correct format, 1 point for identifying three components correctly.*
