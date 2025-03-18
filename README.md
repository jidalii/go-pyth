# Pyth: High-Concurrency Message Aggregation Platform

## Descripion

Pyth is a high-concurrent message aggregation platform.

- **Tech stack**: Go, go-zero, MySQL, Kafka, Redis, ants, asynq

## Features

### Rating Limiting

## Todo

### pyth-support

- [x] Encapsulation of message queue consumer and producer.

### pyth-handler

#### 1) pre-handler services

- [x] Deduplication service:
  - [x] Service: Content-based + Frequency-based
  - [x] Logic: Simple + TokenBucket + SlidingWindow
- [x] Shielding service
- [x] Discard service
- [ ] Connect each services

#### 2) handlers

- [ ] SMS handler
- [ ] Email handler
- [ ] Twitter handler

### pyth-cron

## Thanks

- austin-go: [https://github.com/rxrddd/austin-go](https://github.com/rxrddd/austin-go)
- austin: [https://github.com/ZhongFuCheng3y/austin](https://github.com/ZhongFuCheng3y/austin)