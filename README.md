# Go ethereum key broker

The goal of this project is to broke ethereum wallets by randomly generating an existing ethereum wallet.

## Probability

__For now the probability to generate a wallet that exists is barely impossible__

#### Existing keys

Number of possibles public address: 1.5e+48
> 40 hex character = 16 ^ 40 = 1.5e+48

#### Active keys

Number of unique active address: 5e+7
> https://etherscan.io/chart/address

#### Program performance

The current performance of the program is 2Million key per minute (on 12 threads): 2e+6

#### Sum

A simple formula can be apply: `time_to_break_one_key_in_minutes = nb_possible_public_address / nb_unique_active_address / nb_generated_per_minute`

It gives us: 1.5e+48 / 5e+7 / 2e+6 = 1.5e+34 minutes to break one key

To add a little context, the chance to win at the lottery is one of 1.3+7

## Motivation

Why am I doing this project if these probability are so low?

In some years the security of ethereum and other blockchains could be corrupted by the fact that crypto methods used can become obselete and by the nature of a blockchain, there is no way to updates them.

This project aim to keep an eye on blockchain security and the efficiency of crypto methods used.

Currently this project isn't cracking anything but in some year with the growing power of CPU ([Moore's law](https://en.wikipedia.org/wiki/Moore%27s_law)) the probability of generating an existing wallet will become honest and reachable.

Plus it was fun to learn GO :)

## Usage

Use `--help` flag to display the program usage.

#### Keys

In `keys_db` folder there is 5Gb of existing ethereum keys. This is recommended to use all of them to increase the probability of generating an existing wallet (see above). This is what justify the size of the repository and docker image

#### Using GO

Build the project `go build .`

Then run `go-ethereum-key-broker` binary

#### Using docker

In progress...

## Contribute

All kinds of contributions are welcome ! The most basic way to show your support is to star the project, or to raise issues.