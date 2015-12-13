@echo off
CLS
goto checkPrivileges

:appendToPath
set "RUNPS=powershell -NoProfile -ExecutionPolicy Bypass -Command"
set "OLDPATHPS=[Environment]::GetEnvironmentVariable('PATH', 'User')"
for /f "delims=" %%i in ('%RUNPS% "%OLDPATHPS%"') do (
    set OLDPATH=%%i
)
set NEWPATH=%OLDPATH%;%1
%RUNPS% "[Environment]::SetEnvironmentVariable('PATH', '%NEWPATH%', 'User')"
goto :eof

:checkPrivileges
NET FILE 1>NUL 2>NUL
if '%errorlevel%' == '0' ( goto gotPrivileges ) else ( goto getPrivileges )

:getPrivileges
if '%1'=='ELEV' (shift & goto gotPrivileges)
echo.
echo **************************************
echo Installing GitHub CLI as Administrator
echo **************************************

setlocal DisableDelayedExpansion
set "batchPath=%~0"
setlocal EnableDelayedExpansion
echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\OEgetPrivileges.vbs"
echo UAC.ShellExecute "!batchPath!", "ELEV", "", "runas", 1 >> "%temp%\OEgetPrivileges.vbs"
"%temp%\OEgetPrivileges.vbs"
exit /B

:gotPrivileges

setlocal & cd /d %~dp0

set HUB_BIN_PATH="%LOCALAPPDATA%\GitHubCLI\bin"
IF EXIST %HUB_BIN_PATH% GOTO DIRECTORY_EXISTS
mkdir %HUB_BIN_PATH%
set "path=%PATH%;%HUB_BIN_PATH:"=%"
call :apppendToPath %HUB_BIN_PATH:"=%
:DIRECTORY_EXISTS

:: Delete any existing programs
2>NUL del /q %HUB_BIN_PATH%\hub*

1>NUL copy .\bin\hub.exe %HUB_BIN_PATH%\hub.exe

echo hub.exe installed successfully. Press any key to exit
pause > NUL
