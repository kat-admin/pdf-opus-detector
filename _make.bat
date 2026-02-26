@echo off

go build .

if %errorlevel%==0 goto copyFile

echo.
echo ################################################################
echo ##### Fehler ###################################################
echo ################################################################
echo %errorlevel%
goto endOfScript

:copyFile
if "%COMPUTERNAME%"=="S2S-PC25-17" goto endOfScript

xcopy *.exe Z:\transfer /Y

:endOfScript
