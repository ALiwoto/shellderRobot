@echo off
git pull && powershell -command "Stop-service -Force -name "ShellderRobot" -ErrorAction SilentlyContinue; go build; Start-service -name "ShellderRobot""
:: Hail Hydra
