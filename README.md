## NewChainBlind

This is a commandline client for the NewChain.
It's designed to be easy-to-use. It contains the following:
* Anonymous transaction

## QuickStart

### Download from releases

Binary archives are published at https://release.cloud.diynova.com/newton/NewChainBlind/.

### Building the source

To get from gitlab via `go get`, this will get source and install dependens(cobra, viper, logrus).

#### Windows

install command

```bash
go get gitlab.newtonproject.org/yangchenzhong/NewChainBlind
```

run NewCommander

```bash
%GOPATH%/bin/NewChainBlind.exe
```

#### Linux or Mac

install:

```bash
git config --global url."git@gitlab.newtonproject.org:".insteadOf "https://gitlab.newtonproject.org/"
go get gitlab.newtonproject.org/yangchenzhong/NewChainBlind
```
run NewCommander

```bash
$GOPATH/bin/NewChainBlind
```

### Usage

#### Help

Use command `NewChainBlind help` to display the usage.

```bash
Usage:
  NewChainBlind [flags]
  NewChainBlind [command]

Available Commands:
  account     Manage NewChain accounts
  balance     Get balance of address
  blind       Blind 1NEW for address
  deposit     Use the TxHash to deposit NEW(decimal not counted)
  faucet      Get free money for address
  help        Help about any command
  info        Show the info of the bank
  init        Initialize config file and bank
  sign        Sign blinded file for address of cash 1NEW
  unblind     Unblind 1NEW for address
  verify      Verify unblinder sign file and send 1NEW to address
  version     Get version of NewCommander CLI

Flags:
  -c, --config path            The path to config file (default "./config.toml")
  -h, --help                   help for NewChainBlind
  -i, --rpcURL url             NewChain json rpc or ipc url (default "https://rpc1.newchain.newtonproject.org")
  -w, --walletPath directory   Wallet storage directory (default "./wallet/")

Use "NewChainBlind [command] --help" for more information about a command.
```

#### Use config.toml

You can use a configuration file to simplify the command line parameters.

One available configuration file `config.toml` is as follows:


```conf
rpcurl = "https://rpc1.newchain.newtonproject.org"
unit = "NEW"
usedhash = ["0xe5ed16d511bb6f511c04fdc19be38b47d3494844e94db04f36fe2c7d4348b75b"]
usedunblinderhash = ["0xb18f64a9cc9ce8082675b65c174041d5fe2fb7b4b2d50084961c280ebe896d29"]
walletpath = "./wallet/"

[balances]
  0x2a8996ebb0314717dfdcd879685a9246649d7bc1 = 1
  0x511eef866c847c78b3ff67b581064007166575b1 = 22
  0x97549e368acafdcae786bb93d98379f1d1561a29 = 53

[bank]
  address = "0x873054eAcB22516E1dBC966C9aE338eef40FE15c"
  cash = "1"
  rsapemprivatekeyfile = "./key.pem"
  rsapempublickeyfile = "./pubkey.pem"
```

#### Initialize config file

```bash
# Initialize config file
NewChainBlind init
```

Just press Enter to use the default configuration, and it's best to create a new user.

```bash
Initialize config file
Enter file in which to save (./config.toml):
Enter the wallet storage directory (./wallet/):
Enter NewChain json rpc or ipc url (https://rpc1.newchain.newtonproject.org):
Enter path of Private Key RSA Pem file(./key.pem):
Enter path of Public Key RSA Pem file(./pubkey.pem):
Create a new bank account or not: [Y/n]
Your new account is locked with a password. Please give a password. Do not forget this password.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
0x873054eAcB22516E1dBC966C9aE338eef40FE15c
Your configuration has been saved in  ./config.toml
```

#### Info

```bash
# Get the info of the bank
NewChainBlind info
```

### Deposit (by Bank)

Use `NewCommander pay` to send 10 NEW to the address of bank, 
then use `NewChainBlind deposit <txHash>` to add the balance

```bash
# list all accounts of the walletPath
NewChainBlind deposit 0xfea3844616766cea49c21447d6e5a8c4521192b820c567bff81678615966987e
```

### Blind (By User1)

```bash
# blind 1NEW, only support 1NEW
NewChainBlind blind pubkey.pem
```

the output of the cmd is as follow:

```txt
Data:  20190511140547.data
Blinded:  20190511140547.blinded
Unblinder:  20190511140547.unblinder
Send blinded file to bank
```

Send 20190511140547.blinded to bank

### Sign (By Bank)

```bash
# Sign the blinded file
NewChainBlind sign 20190511140547.blinded 0x97549E368AcaFdCAE786BB93D98379f1D1561a29
```

the output is as follow:
```txt
Sub 1 NEW for 0x97549E368AcaFdCAE786BB93D98379f1D1561a29
Current balance of 0x97549E368AcaFdCAE786BB93D98379f1D1561a29 is 59 NEW
Sign:  20190511140547.blinded.sig
Send signed blinded file back to user1
```

then send the 20190511140547.blinded.sig back to the user

### Unblind (By User1)
```bash
# Unblind with unblinder
NewChainBlind unblind pubkey.pem 20190511140547.blinded.sig 20190511140547.unblinder
```

the output is as follow:
````txt
UnblinderSig:  20190511140547.unblinder.sig
Send unblinder sig file and cash data to others
````

### Pay
1. when `user1` pay 1NEW to `user2`, just send the file `20190511140547.unblinder.sig` and `20190511140547.data`
2. when `user2` get the file, then send both of the file `20190511140547.unblinder.sig` and `20190511140547.data` to bank to verify

### Verify (By Bank)
```bash
# Verify unblinder and data, then add 1NEW to user2
NewChainBlind verify 20190511140547.unblinder.sig 20190511140547.data 0x2a8996eBb0314717dfdCd879685A9246649D7BC1
```

if ok, then the bank pay 1 NEW to `user2`