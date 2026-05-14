@echo off
echo Building Vue3 frontend...
cd /d %~dp0
call npm install
call npm run build
echo Build complete! Output in web/dist/
cd ..
