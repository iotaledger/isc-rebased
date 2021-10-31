@echo off
cd %1
if not exist ts\%1\lib.ts goto :xit
echo Building %1
schema -ts %2
call asc ts/%1/lib.ts --binaryFile ts/pkg/%1_ts.wasm --textFile ts/pkg/%1_ts.wat
:xit
cd ..
