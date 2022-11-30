@echo off
for /f "tokens=*" %%a in ('git rev-parse HEAD') do (
    set GITHASH=%%a
)
echo %GITHASH%
