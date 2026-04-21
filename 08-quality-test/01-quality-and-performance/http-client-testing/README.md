# Quality — HTTP Client Testing

Manual mock establishes the pattern (interface → fake struct). Function-injection mock extends the fake with a GetFunc field. Table-driven mock applies §13 testing patterns to the function mock. Testify/mock replaces the hand-written mock with a library.

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| HM.1 | [manual mock](./3-manual-mock) | MockHTTPClient struct · io.NopCloser · injected interface | 🟢 entry |
| HM.2 | [function-injection mock](./4-manual-mock-2) | GetFunc field · GetCalls spy · per-test dynamic behaviour | HM.1 |
| HM.3 | [table-driven mock](./5-manual-mock-table-driven) | []struct test matrix · sub-tests per case · errContains check | HM.1, HM.2 |
| HM.4 | [testify/mock](./6-testify-mock) | mock.Mock embed · .On/.Return · AssertNumberOfCalls | HM.1, HM.2, HM.3 |
