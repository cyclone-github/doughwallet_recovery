# doughwallet_recovery

### Simple CLI tool to recover Dough Wallets

Usage:
```
./doughwallet_recovery
 ---------------------------- 
| Dough Wallet Recovery Tool |
 ---------------------------- 

Enter your Dough Wallet's 12-word recovery phrase:
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

### Dough Wallet Info:
- Derivation Path = `m/0'/0,1/0`
- Hardened bit = `0x9e000000`
- Official Dough Wallet source code
  - https://github.com/peritus/doughwallet
- Official Dough Wallet Recover Tool v2
  - https://github.com/peritus/doughwallet-recovery2