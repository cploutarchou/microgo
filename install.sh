#!/bin/bash

echo "Installing MicroGO binaries..."
wget https://github.com/cploutarchou/MicroGO/releases/download/v1.0.0/microGo
userDir=$USER
echo "Enter your password to install MicroGO binaries using sudo"
sudo chmod +x microGo
sudo mv microGo /usr/local/bin/microGo # move to /usr/local/bin/microGo
echo "MicroGO binaries installed"
echo "Export env"
# shellcheck disable=SC2016
echo 'export PATH=$PATH:/usr/local/bin/microGo' >>/home/"$userDir"/.bashrc
echo "MicroGO binaries exported to PATH"
source /home/"$userDir"/.bashrc
