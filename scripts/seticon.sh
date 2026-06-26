#!/bin/bash
set -euo pipefail
APP="${1:-build/bin/ignite.app}"
sips -s format icns appicon.png --out "$APP/Contents/Resources/iconfile.icns" >/dev/null 2>&1
touch "$APP"
echo "Icon set on $APP"
