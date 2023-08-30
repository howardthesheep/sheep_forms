@ECHO OFF
set BASE_DIR="%AppData%\sheep\sheep-forms"

if "%1"=="build" (
	call :buildApp
) ELSE IF "%1"=="install" (
  	call :buildApp

	if not exist "%BASE_DIR%\NUL" (
		echo Creating program directory for installation...
		mkdir %BASE_DIR%
	)

	echo Copying in program templates...
	xcopy templates\ %BASE_DIR%\ /e

	echo Copying .exe to program directory...
	copy sheep-forms.exe %BASE_DIR%\

	echo Adding program directory to PATH...
	echo Please add %BASE_DIR% to your system environment variables!!!

	exit /b
) ELSE IF "%1"=="run" (
	echo Running sheep-forms.exe...
	start "sheep-forms" CALL sheep-forms.exe
) ELSE (
	echo "Unrecognized argument "%1". Valid args are 'build' 'run' and 'install'"
	exit /b
)

echo Done!
exit /b

:buildApp
	 echo Building .exe with go build...
	 setlocal enabledelayedexpansion enableextensions
     set GO_FILES=
     for %%x in (src\*.go) do set GO_FILES=!GO_FILES! %%x

	 go build -o sheep-forms.exe %GO_FILES%
	 endlocal