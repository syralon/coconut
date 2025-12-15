package proto

//go:generate protoc --go_out=paths=source_relative:. --coconut-errors_out=paths=source_relative:. ./syralon/coconut/errors/*.proto
//go:generate protoc --go_out=paths=source_relative:. ./syralon/coconut/field/*.proto
