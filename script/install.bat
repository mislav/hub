@echo off
CLS
goto checkPrivileges

:: Unfortunately, Windows doesn't have a decent built-in way to append a string to the $PATH.
:: setx is convenient, but it 1) truncates paths longer than 1024 characters, and
:: 2) mucks up the user path with the machine-wide path.
:: This function takes care of these problems by calling Environment.Get/SetEnvironmentVariable
:: via PowerShell, which lacks these issues.
:appendToUserPath
setlocal EnableDelayedExpansion
set "RUNPS=powershell -NoProfile -ExecutionPolicy Bypass -Command" :: Command to start PowerShell.
set "OLDPATHPS=[Environment]::GetEnvironmentVariable('PATH', 'User')" :: PowerShell command to run to get the old $PATH for the current user.

:: Capture the output of %RUNPS% "%OLDPATHPS%" and set it to OLDPATH
for /f "delims=" %%i in ('%RUNPS% "%OLDPATHPS%"') do (
    set "OLDPATH=!OLDPATH!%%i"
)

set "NEWPATH=%OLDPATH%;%1"
:: Set the new $PATH
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
call :appendToUserPath "%HUB_BIN_PATH:"=%"
:DIRECTORY_EXISTS

:: Delete any existing programs
2>NUL del /q %HUB_BIN_PATH%\hub*

1>NUL copy .\bin\hub.exe %HUB_BIN_PATH%\hub.exe

echo hub.exe installed successfully. Press any key to exit
pause > NUL
