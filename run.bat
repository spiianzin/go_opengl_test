set filename="cells.go"
set arg=%1 

if "%arg%"=="1" set filename="simple_triangle.go"
if "%arg%"=="2" set filename="square.go"
if "%arg%"=="3" set filename="cells.go" 

pushd game_of_life
go run %filename%
popd