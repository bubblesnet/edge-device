@echo off
for /f "tokens=*" %%a in ('date /t') do (
		    set TIMESTAMP='%%a
)
for /f "tokens=*" %%a in ('time /t') do (
		    set TIMESTAMP=%TIMESTAMP% %%a'
)

echo %TIMESTAMP: =_%

