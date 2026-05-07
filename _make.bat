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
if "%COMPUTERNAME%"=="S2S-PC25-17" goto copyToTimbo

:copyToTimbo

xcopy *.exe W:\timbo\downloads\temp /Y

echo.
echo https://tim.s2s.gmbh/timbo/downloads/temp/pdf-opus-detector.exe
echo.

goto endOfScript

:copyToTransfer

xcopy *.exe Z:\transfer /Y

:endOfScript
