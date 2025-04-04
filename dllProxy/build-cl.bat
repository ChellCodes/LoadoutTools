@echo off
@REM echo dir = %cd%

:: compiler choice
set VCVARS="C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Auxiliary\Build\vcvars32.bat"

:: builds dir for out executables/libraries
if not exist builds mkdir builds
if not exist builds\int mkdir builds\int

@REM MTd: static crt debug
@REM Zi: debug info
@REM WX: warnings as errors
@REM W4: warning level 4
@REM FC: full path in diagnostics
@REM EHsc: warning C4530: C++ exception handler used, but unwind semantics are not enabled
set COMPILER_FLAGS=^
                   /W4 ^
                   /Zi ^
                   /FC ^
                   /EHsc ^
                   /nologo

set INCLUDES=/I"imgui" ^
            /I"C:\Windows\System32" ^
            /I "%WindowsSdkDir%Include\um" ^
            /I "%WindowsSdkDir%Include\shared"

set LINKER_FLAGS=^
                 /incremental:no ^
                 /nologo

set LINK_LIBRARIES=^
             VERSION.lib ^
             user32.lib

:: build
echo.
echo --------------------------- STARTING BUILD ---------------------------
echo.


:: devenv
if not defined DevEnvDir (
    call %VCVARS% || (echo "vcvars64.bat not found" & goto build_failed)
)

:: build
echo building console.exe...
:: cl console.c /MT %COMPILER_FLAGS% /Fd"builds\\int\\" /Fo"builds\\int\\" /link /subsystem:console %LINKER_FLAGS% %LINK_LIBRARIES% /OUT:"builds\\console.exe" /PDB:"builds\\int\\"

set IMGUISOURCE=imgui\backends\imgui_impl_dx11.cpp imgui\backends\imgui_impl_win32.cpp imgui\imgui*.cpp
echo building version.dll...
cl dll.cpp %IMGUISOURCE% %INCLUDES% %COMPILER_FLAGS% /Fd"builds\\int\\" /Fo"builds\\int\\" /link %LINKER_FLAGS% %LINK_LIBRARIES% /DLL /DEF:"exports.def" /OUT:"builds\\version.dll"

if ERRORLEVEL 1 goto build_failed

if "%~1"=="run" (
    copy /Y "builds\version.dll" "N:\SteamLibrary\steamapps\common\Loadout\"
    start "" "steam://rungameid/208090"
)
:: build success
echo.
echo --------------------------- BUILD COMPLETE ---------------------------
echo.
exit /b 0


:: :(
:build_failed
echo.
echo --------------------------- BUILD FAILED -----------------------------
echo.
exit /b 1
