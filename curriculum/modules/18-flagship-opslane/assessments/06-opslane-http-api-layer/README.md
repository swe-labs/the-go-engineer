# Opslane HTTP API layer

## Mission

Understand and apply Opslane HTTP API layer in the context of professional Go software engineering.

## Prerequisites

- opslane-05

## Mental Model

The HTTP API layer is the application's front door — it receives requests, ensures they're well-formed, routes them to the right internal department, and returns a clean response. A well-designed front door makes visitors feel welcome and guided.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, the router matches the request path and method using a radix tree. middleware.Chain composes middleware into a single handler. Request body decoding uses json.NewDecoder with a limited io.LimitReader to prevent memory exhaustion.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/06-opslane-http-api-layer
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Putting all business logic directly in HTTP handlers instead of a service layer.
- Not validating request bodies before passing them to the service layer.
- Returning stack traces or internal error messages in API responses.

## In Production

RESTful APIs power the modern web — every major service (GitHub API, Stripe API, Slack API) exposes functionality through well-designed HTTP APIs. Opslane follows the same conventions.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-07`.
