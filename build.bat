@echo off

REM 更改为自己的目标目录,执行时会自动复制过去,或者注释掉此句以关闭自动复制功能 

SET DevDir=D:\Program Files (x86)\MiraiProject\MiraiOk-M4\data\MiraiNative\plugins

echo Setting proxy
SET GOPROXY=https://goproxy.cn

echo Checking go installation...
go version > nul
IF ERRORLEVEL 1 (
	echo Please install go first...
	goto RETURN
)

echo Checking gcc installation...
gcc --version > nul
IF ERRORLEVEL 1 (
	echo Please install gcc first...
	goto RETURN
)

echo Checking cqcfg installation...
cqcfg -v
IF ERRORLEVEL 1 (
	echo Install cqcfg...
	go get github.com/Tnze/CoolQ-Golang-SDK/tools/cqcfg
	IF ERRORLEVEL 1 (
		echo Install cqcfg fail
		goto RETURN
	)
)

echo Generating app.json ...
go generate
IF ERRORLEVEL 1 (
	echo Generate app.json fail
	goto RETURN
)
echo.

echo Setting env vars..
SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386

echo Building app.dll ...
go build -ldflags "-s -w" -buildmode=c-shared -o app.dev.dll
IF ERRORLEVEL 1 (pause) ELSE (echo Build success!)

if defined DevDir (
    echo Copy app.dll and app.json ...
    for %%f in (app.dev.dll) do COPY %%f "%DevDir%\%%f" > nul
    IF ERRORLEVEL 1 pause
)

exit /B

:RETURN
pause
exit /B
