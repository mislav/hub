@echo off
CLS

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
1>NUL setx PATH "%PATH%" /M
:DIRECTORY_EXISTS

:: Delete any existing programs
2>NUL del /q %HUB_BIN_PATH%\hub*

1>NUL copy .\bin\hub.exe %HUB_BIN_PATH%\hub.exe

echo hub.exe installed successfully. Press any key to exit
pause > NUL
