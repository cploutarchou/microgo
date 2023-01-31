#!/bin/bash

userDir=$USER
echo "Installing MicroGO binaries..."

if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ $(uname -m) == "x86_64" ]]; then
        curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.5/microGo
        mv microGo /Users/"$userDir"/go/bin/microGo
    else
        curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.5/microGo-MacOS-ARM64
        mv microGo-MacOS-ARM64 /Users/"$userDir"/go/bin/microGo
    fi
    
    echo "Enter your password to install MicroGO binaries using sudo"
    chmod +x /Users/"$userDir"/go/bin/microGo
    echo "MicroGO binaries installed"
    echo "Export env"
    
    if grep -q "alias microGo" ~/.zprofile; then
        echo "Already exported"
    else
        source ~/.zprofile
        comm="alias microGo=/Users/$userDir/go/bin/microGo"
        echo $comm >>/Users/"$userDir"/.zprofile
        echo "MicroGO binaries exported to PATH"
        source ~/.zprofile
    fi
else
    echo "Unsupported OS"
fi
