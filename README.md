[![Readme Card](https://github-readme-stats-fast.vercel.app/api/pin/?username=cyclone-github&repo=doughwallet_recovery&theme=gruvbox)](https://github.com/cyclone-github/doughwallet_recovery/)

[![Go Report Card](https://goreportcard.com/badge/github.com/cyclone-github/doughwallet_recovery)](https://goreportcard.com/report/github.com/cyclone-github/doughwallet_recovery)
[![GitHub issues](https://img.shields.io/github/issues/cyclone-github/doughwallet_recovery.svg)](https://github.com/cyclone-github/doughwallet_recovery/issues)
[![License](https://img.shields.io/github/license/cyclone-github/doughwallet_recovery.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/cyclone-github/doughwallet_recovery.svg)](https://github.com/cyclone-github/doughwallet_recovery/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/cyclone-github/doughwallet_recovery.svg)](https://pkg.go.dev/github.com/cyclone-github/doughwallet_recovery)

### Simple CLI tool to recover Dough Wallets

The defunct Dough Wallet iPhone app used non-standard settings which made it impossible to recover or use your Dogecoins without the Dough Wallet app, but since the iPhone app and author's website are long gone, this left many users with no hope of recovering their Dogecoins.

`Enter Dough Wallet Recovery`. The non-standard Dough Wallet settings have been methodically researched, reversed, and reimplemented into this tool which allows users to regain access to their lost Dough Wallet Dogecoins. And for the first time ever, the custom Dough Wallet Dogecoin settings have been publically released (see info below and source code).

Enjoy, 

~ Cyclone 

---

Usage:
```
./doughwallet_recovery
 --------------------------- 
|   Dough Wallet Recovery   |
| github.com/cyclone-github |
 --------------------------- 

Enter your Dough Wallet's 12-word recovery phrase:
abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon

External (Receive) Chain: m/0'/0/n

  m/0'/0/0
    Address:     DNf2JUSPdFBUAwC6NauoCMwzQ5idABnf3J
    Private Key: QTk8fDB45XbHUdesGyPgtA3pJvBmpxNtSEw3WWi958JkApPR4ctV

Internal (Change) Chain: m/0'/1/n

  m/0'/1/0
    Address:     DP4ZdihCRgx88YtfiM8nSW7wkcRHeiCVUk
    Private Key: QVf7tKti8JvGUrm3DchjAB7iqVb7MEmqqXrcTzhmYoXnm7SqRnSH
```
Generate 10 addresses per derivation path:
```
./doughwallet_recovery -count 10
 ---------------------------- 
| Dough Wallet Recovery Tool |
 ---------------------------- 

Enter your Dough Wallet's 12-word recovery phrase:
```

Pipe seed phrase:

```
echo "dough wallet word seed phrase.." | ./doughwallet_recovery
```

---

### Install latest release:
```
go install github.com/cyclone-github/doughwallet_recovery@latest
```
### Install from latest source code (bleeding edge):
```
go install github.com/cyclone-github/doughwallet_recovery@main
```

### Compile from source:
- This assumes you have Go and Git installed
  - `git clone https://github.com/cyclone-github/doughwallet_recovery.git`  # clone repo
  - `cd doughwallet_recovery`                                               # enter project directory
  - `go mod init doughwallet_recovery`                                      # initialize Go module (skips if go.mod exists)
  - `go mod tidy`                                              # download dependencies
  - `go build -ldflags="-s -w" .`                              # compile binary in current directory
  - `go install -ldflags="-s -w" .`                            # compile binary and install to $GOPATH
- Compile from source code how-to:
  - https://github.com/cyclone-github/scripts/blob/main/intro_to_go.txt

---

### Dough Wallet Info:
- Derivation Path = `m/0'/0,1/0`
  - Source: https://github.com/iancoleman/bip39/commit/4062a567f56ed8a6ec8246d034e0aea94a5e554a
- Hardened bit = `0x9e000000`
  - Source: https://github.com/iancoleman/bip39/commit/d98d01a9d00613be9d2042444ac8db629257e73e
- Official Dough Wallet source code
  - https://github.com/peritus/doughwallet
- Official Dough Wallet Web Recovery Toolkit v2
  - https://github.com/peritus/doughwallet-recovery2
