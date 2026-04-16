# Web — HTTP Client

The naive client demonstrates the problem; the refactor lesson demonstrates the solution. Strictly sequential.

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| HC.1 | [basic GET](./1-get-posts) | http.Get · resp.Body.Close · DefaultClient no-timeout problem | 🟢 entry |
| HC.2 | [refactor for testability](./2-refactor-for-testability) | HTTPClient interface · IoTClient struct · NewIoTClient DI | HC.1 |
