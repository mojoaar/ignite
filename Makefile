.PHONY: build dev clean test

build:
	wails build -clean
	bash scripts/seticon.sh

dev:
	wails dev

test:
	go test ./internal/... -v
	cd frontend && pnpm typecheck && pnpm vitest run

clean:
	rm -rf build frontend/dist frontend/node_modules
