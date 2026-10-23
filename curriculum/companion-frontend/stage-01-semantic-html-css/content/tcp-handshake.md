# TCP Handshake & Lifecycle

Transmission Control Protocol (TCP) provides reliable, ordered, and error-checked delivery of a stream of octets between applications running on hosts communicating via an IP network.

## The Three-Way Handshake

Connection establishment requires exchanging three synchronization segments:

1. **SYN**: Client sends `SYN (seq=x)` to server. State: `SYN-SENT`.
2. **SYN-ACK**: Server acknowledges client's sequence and sends its own: `SYN-ACK (seq=y, ack=x+1)`. State: `SYN-RECEIVED`.
3. **ACK**: Client acknowledges server sequence: `ACK (seq=x+1, ack=y+1)`. State: `ESTABLISHED`.

## Connection Teardown (Four-Way Handshake)

Either endpoint can initiate closing:

1. Endpoint A sends `FIN`.
2. Endpoint B responds with `ACK`.
3. Endpoint B sends its own `FIN` once all buffered writes complete.
4. Endpoint A sends `ACK` and enters the `TIME-WAIT` state (typically 2 * MSL) to ensure the final ACK was delivered.
