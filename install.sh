#!/bin/bash
userDir=$USER
echo "Installing MicroGO binaries..."
curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.5/microGo
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if grep -q "alias microGo" ~/.zprofile; then
        echo "Already exported"
        sudo mv microGo /usr/local/bin/microGo
        chmod +x /usr/local/bin/microGo
    else
        sudo mv microGo /usr/local/bin/microGo
        chmod +x /usr/local/bin/microGo
        echo "export PATH=$PATH:/usr/local/bin/microGo" >>/home/"$userDir"/.bashrc
        echo "MicroGO binaries exported to PATH"
        source /home/"$userDir"/.bashrc
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    mv microGo /Users/"$userDir"/go/bin/microGo
    echo "Enter your password to install MicroGO binaries using sudo"
    chmod +x /Users/"$userDir"/go/bin/microGo # move to /usr/local/bin/microGo
    echo "MicroGO binaries installed"
    echo "Export env"
    if grep -q "alias microGo" ~/.zprofile; then
        echo "Already exported"
    else
        source ~/.zprofile
        comm="alias microGo=/Users/$userDir/go/bin/microGo"
        echo $comm >>/Users/"$userDir"/.zprofile # add to .zprofile
        echo "MicroGO binaries exported to PATH"
        source ~/.zprofile
    fi
else
    echo "Unsupported OS"
fi
