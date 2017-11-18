# Optimized Group Mining
A web application for convenient device management of multiple XMG mining devices in a local network.

## Why Coin Magi and Raspberry Pi?
With the quick expansion of cryptocurrencies, more and more people are willing to join the community and contribute with processing power to receive their share form the profit. However, mining famous cryptocurrencies such as Bitcoin, Ethereum, Monero, etc. is not for everyone as it requires extensive computational power that is not affordable by everyone. Thus these are being sustained and controlled by a small portion of the population that possesses the necessary equipment. Magi Coin presents a change. It is especially designed to be mined on a CPU thus everyone can afford it, but it also is optimized for the ARM architecture used in the world-wide famous microcomputer Raspberry Pi. Mining Magi Coin, or XMG, on Raspberry Pi is quite beneficial due to its extreme low power consumption compared to a conventional computer and computational capabilities. Thus, rub the dust off your old Raspberry Pi boards at home and join the Magi community now using the group mining utility presented in this repository. You will be able to handle multiple boards at once, control their mining parameters and you would hardly need to SSH multiple times in each of them again.

## Installation
First, download the mining manager with the following and enter its working directory:
```
git clone https://github.com/vanjo9800/GroupMiner && cd GroupMiner/
```

### Configure the miner
Second, install the Proof-of-Work mining software. Currently, the project uses the Magi Coin (XMG), so it itlizies their CPU miner which can be found [here](https://github.com/magi-project/m-cpuminer-v2). It needs to be build locally in order to perform adequatively on the native CPU architecture. For a Raspberry Pi, or another ARM device, you may want to use NEON instructions which can lead to slightly bigger performance. All this can be accomplished via the following command, it will prompt you with choices for the different configurations:
```
# in the GroupMiner/ directory
scripts/setup-miner.sh
```

### Configure the prerequistites
Third, setup the Go environment and install the necessary dependencies of the package. This can be established via the following commands:
```
# in the GroupMiner/ directory
scripts/setup-go.sh
```

### Running the app
After this, you can run the app by executing the generated binary files ```server``` and ```client``` in the bin directory, depending on the mode you would like to use.

## Running configuration
To configure the client, or server, you can use the ```.conf``` files located in the ```config/``` directory.

## Future development
This project is only in its initial phase meaning that much more features are going to be implemented soon. These include:
* Login and Registration Page
* Hashrate staticics for each device
* Graphs with performance results
* Pool management