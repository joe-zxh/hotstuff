@echo off

SET GOOS=linux
SET GOARCH=amd64

set exeSuffix=

echo ����hotstuff...
del .\hotstuffserver%exeSuffix% 
del .\hotstuffclient%exeSuffix% 
go build -o .\hotstuffserver%exeSuffix% .\cmd\hotstuffserver\main.go
go build -o .\hotstuffclient%exeSuffix% .\cmd\hotstuffclient\main.go

