#!/bin/bash
python3 -c "
from PIL import Image
import sys
src = Image.open('appicon.png')
dest = sys.argv[1] if len(sys.argv) > 1 else 'build/bin/ignite.app/Contents/Resources/iconfile.icns'
sizes = [(1024,1024),(512,512),(256,256),(128,128),(64,64),(32,32),(16,16)]
src.save(dest, format='ICNS', append_images=[src.resize(s) for s in sizes[1:]])
" "$@"
touch "${1:-build/bin/ignite.app}"
echo "Icon set on ${1:-build/bin/ignite.app}"