#!/bin/bash
echo "MicroGO Installation"
echo " ____    ____    _                                ______      ___    "
echo "|_   \  /   _|  (_)                             .' ___  |   .'   \`.  "
echo "  |   \/   |    __    .---.   _ .--.    .--.   / .'   \_|  /  .-.  \\ "
echo "  | |\  /| |   [  |  / /'\`\\] [ \`/'\`\\ ] / .'\`\\ \ | |   ____  | |   | | "
echo " _| |_\/_| |_   | |  | \\__.   | |     | \\__. | \\ \`\.___]  | \\  \`-'  / "
echo "|_____||_____| [___] '.___.' [___]     '.__.'   \`._____.'   \`.___.'  "
echo ""
echo "v1.0.6"
echo "Author: Christos Ploutarchou"
echo "Installing MicroGO binaries..."
userDir=$USER

spinner() {
  local pid=$1
  local delay=0.75
  # shellcheck disable=SC1003
  local spinstr='|/-\'
  # shellcheck disable=SC2143
  while [ "$(ps a | awk '{print $1}' | grep "$pid")" ]; do
    local temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    local spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
}

if [[ "$OSTYPE" == "darwin"* ]]; then
  if [[ $(uname -m) == "x86_64" ]]; then
    curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.6/microGo-MacOS-x86_64 &
    spinner $!
    mv microGo /Users/"$userDir"/go/bin/microGo
  else
    curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.6/microGo-MacOS-ARM64 &
    spinner $!
    mv microGo-MacOS-ARM64 /Users/"$userDir"/go/bin/microGo
  fi

  echo "Enter your password to install MicroGO binaries using sudo"
  chmod +x /Users/"$userDir"/go/bin/microGo
  echo "MicroGO binaries installed"
  echo "Export env"

  if grep -q "alias microGo" ~/.zprofile; then
    echo "Already exported"
  else
    # shellcheck disable=SC1090
    source ~/.zprofile
    comm="alias microGo=/Users/$userDir/go/bin/microGo"
    echo "$comm" >>/Users/"$userDir"/.zprofile
    echo "MicroGO binaries exported to PATH"
    # shellcheck disable=SC1090
    source ~/.zprofile
  fi

elif [[ "$OSTYPE" == "linux-gnu" ]]; then
  if [[ $(uname -m) == "x86_64" ]]; then
    curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.6/microGo-Linux-x86_64 &
    spinner $!
    mv microGo-Linux-x86_64 ~/go/bin/microGo
  else
    curl -LJO https://github.com/cploutarchou/MicroGO/releases/download/v1.0.6/microGo-Linux-ARM64 &
    spinner $!
    mv microGo-Linux-ARM64 ~/go/bin/microGo
  fi

  chmod +x ~/go/bin/microGo
  echo "MicroGO binaries installed"
  echo "Export env"

  if grep -q "alias microGo" ~/.bashrc; then
    echo "Already exported"
  else
    # shellcheck disable=SC1090
    source ~/.bashrc
    comm="alias microGo=~/go/bin/microGo"
    echo "$comm" >>~/.bashrc
    echo "MicroGO binaries exported to PATH"
    # shellcheck disable=SC1090
    source ~/.bashrc
  fi
else
  echo "Unsupported OS"
fi
