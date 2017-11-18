#!/bin/bash

mkdir -p log/
touch log/build.log
echo "If any errors occur during the setup, report them via the log located at log/build.log"
echo "Downloading miner..."
git clone https://github.com/magi-project/m-cpuminer-v2 tools/m-cpuminer-v2/ >> log/build.log
echo "Installing build tools."
echo "Please select an operating system:"
echo "1) Ubuntu, or other DEB-based"
echo "2) Arch Linux"
echo "3) Other Linux distribution (you will need to install the necessary build tools yourself; these are libcurl, libgmp, openssl, and jansson)"
echo -n "Select (1/2/None):" && read n
case $n in
    1) sudo apt-get install automake autoconf pkg-config libcurl4-openssl-dev libjansson-dev libssl-dev libgmp-dev make g++ git libgmp-dev;;
    2) sudo pacman -S base-devel;;
esac
echo "Building the miner..."
cd tools/m-cpuminer-v2/
./autogen.sh > ../../log/build.log
echo -n "Would you like NEON support (recommended for Raspberry Pi and other ARM devices)? (y/N)" && read c
echo "Configuring..."
case $c in
    "y") ./configure CFLAGS="-O3 -mfpu=neon" CXXFLAGS="-O3" > ../../log/build.log;;
    *) ./configure CFLAGS="-O3" CXXFLAGS="-O3" > /dev/null > ../../log/build.log;;
esac
echo "BUilding..."
make > ../../log/build.log
cd ../../
echo "Successfully finished setting up the miner!"
