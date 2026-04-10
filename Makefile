## Makefile for dice project: testing and benchmarking helpers

.PHONY: test bench bench-race bench-mem fmt

TESTPKG=./...

test:
	@echo "Running unit tests"
	go test $(TESTPKG)

bench:
	@echo 'Running benchmarks (default)'
	@echo 'Use make bench BENCH="BenchmarkName" to run a specific benchmark.'
	go test -bench=. -benchmem $(TESTPKG)

bench-race:
	@echo "Running benchmarks with race detector (note: slower)"
	go test -race -bench=. -benchmem $(TESTPKG)

bench-mem:
	@echo "Running benchmarks and writing cpu/alloc profiles to ./profiles/"
	mkdir -p profiles
	go test -bench=. -benchmem -cpuprofile profiles/cpu.prof -memprofile profiles/mem.prof $(TESTPKG)

fmt:
	gofmt -w .
