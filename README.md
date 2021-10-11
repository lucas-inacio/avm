# avm
arduino-cli version manager

## Building
After cloning this repository enter de root directory and type:
```
go build cmd/avm/avm.go
```

## Usage
Typing avm in the terminal will show a help text. A brief list of currently supported commands:
- get - Will download the arduino-cli tool for your platform
- version - Prints installed version
- available - List available releases in the arduino-cli repo
- update - Update your current installation if needed
- install - Install a specific version
