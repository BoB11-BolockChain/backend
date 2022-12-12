#!/bin/bash

echo "Deploy PDxF Backend Server"

Front_REPOSITORY="https://github.com/BoB11-BolockChain/frontend.git"
Back_REPOSITORY="https://github.com/BoB11-BolockChain/backend.git"

echo "apt update"
apt update

echo "Install Golang..."
if which go > /dev/null; then
  echo "Golang already installed."
else
  rm -rf /usr/local/go
  latest=$(curl https://go.dev/VERSION?m=text)
  wget "https://dl.google.com/go/$latest.linux-amd64.tar.gz"
  tar -C /usr/local -xzf $latest.linux-amd64.tar.gz
  printf "\nPATH=\$PATH:/usr/local/go/bin" >> home.profile / etc.profile
fi

echo "Install golang dependencies..."
go get ./...

echo "Install sqlite3 database..."
if which sqlite3 > /dev/null; then
  echo "sqlite3 already installed."
else
  apt install sqlite3
fi

echo "Import database schema..."
sqlite3 schema

echo "Install virtual environment..."
echo "Install Docker..."
if which docker > /dev/null; then
  echo "Docker already installed."
else
  apt install docker
fi

echo "Install libvirt..."
if which virsh > /dev/null; then
  echo "virsh already installed."
else
  sudo apt install qemu-kvm libvirt-daemon-system

fi

echo "Install qemu..."
echo "Install kvm..."
echo "Install libvirt..."
echo "Install novnc..."

echo "Done"