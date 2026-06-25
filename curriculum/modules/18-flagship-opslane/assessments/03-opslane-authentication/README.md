# Opslane authentication

## Mission

Understand and apply Opslane authentication in the context of professional Go software engineering.

## Prerequisites

- opslane-02

## Mental Model

Authentication is the application's bouncer — it checks IDs at the door before letting anyone in. A good bouncer checks every ID, every time, and knows how to spot fakes.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, JWT is a base64-encoded JSON payload with a cryptographic signature. The signature uses HMAC-SHA256 (symmetric) or RSA/ECDSA (asymmetric). bcrypt includes a random salt and is intentionally slow to thwart brute force attacks.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/03-opslane-authentication
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Storing passwords in plain text instead of hashing them.
- Using a fixed API key that never rotates.
- Implementing custom authentication instead of using established standards.

## In Production

JWT-based authentication is used by Kubernetes, Docker, GitHub API, and most modern web services. Opslane implements the same industry-standard patterns with Go.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-04`.
