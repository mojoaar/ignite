.PHONY: build dev clean test publish

build:
	wails build -clean
	bash scripts/seticon.sh

publish: build
	cd build/bin && hdiutil create -volname Ignite -srcfolder Ignite.app -ov -format UDZO ../Ignite.dmg
	cd build/bin && zip -r ../Ignite.zip Ignite.app
	@echo ""
	@echo "Published: build/Ignite.dmg  build/Ignite.zip"
	@ls -lh build/Ignite.dmg build/Ignite.zip

deploy:
	cp build/Ignite.dmg site/assets/Ignite.dmg
	cp build/Ignite.zip site/assets/Ignite.zip
	@echo "Deployed binaries to site/assets/"

dev:
	wails dev

test:
	go test ./internal/... -v
	cd frontend && pnpm typecheck && pnpm vitest run

clean:
	rm -rf build frontend/dist frontend/node_modules
