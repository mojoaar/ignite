.PHONY: build dev clean test publish

build:
	wails build -clean
	bash scripts/seticon.sh

publish: build
	cd build/bin && hdiutil create -volname Ignite -srcfolder ignite.app -ov -format UDZO ../Ignite.dmg
	cd build/bin && zip -r ../Ignite.zip ignite.app
	@echo ""
	@echo "Published: build/Ignite.dmg  build/Ignite.zip"
	@ls -lh build/Ignite.dmg build/Ignite.zip

dev:
	wails dev

test:
	go test ./internal/... -v
	cd frontend && pnpm typecheck && pnpm vitest run

clean:
	rm -rf build frontend/dist frontend/node_modules
