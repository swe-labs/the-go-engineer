# Opslane payment pipeline

## Mission

Understand and apply Opslane payment pipeline in the context of professional Go software engineering.

## Prerequisites

- opslane-07

## Mental Model

The payment pipeline is the application's cash register — it must be absolutely correct (never overcharge, never undercharge), resilient (retry on failure), and auditable (every transaction is logged).

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, Stripe's API uses TLS for transport security. Webhook signatures use HMAC-SHA256 with the webhook secret. The payment pipeline uses database transactions to ensure billing state updates are atomic.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/08-opslane-payment-pipeline
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Calling the payment provider synchronously in the HTTP handler, blocking the user.
- Not idempotent — retrying a failed payment charges the user twice.
- Not validating webhook signatures, allowing fake payment events.

## In Production

Stripe processes billions of dollars in payments through webhook-based pipelines. Opslane implements the same architecture pattern used by production payment systems.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-09`.
