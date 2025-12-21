# Quick start

`go run ./cmd/init-database`
`go run ./cmd/example`


# Framework

- api:
    The server implementations. It will be grpc service implementation in usually.
- domain.biz:
    The define of business. Some complex logic or reusable logic can be wrapped into here.
- domain.repository:
    The abstract define of data layer.
- infrastructure.data:
    The data layer implementation. It will be database access operate in usually.
- infrastructure.dependency:
    The third party tools, such as database, redis... and so on.
- server:
    The entrances.