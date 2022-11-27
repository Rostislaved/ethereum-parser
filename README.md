# Ethereum Parser

Daemon application. Listens to new blocks from the ethereum blockchain.


If parser encounters a transaction with subscribed address it stores this transaction to storage.

## API:

### 1. GetCurrentBlockNumber:  
```
curl --request GET \
--url http://localhost:8080/get-current-block
```

### 2. GetStorageInfo:  
```
curl --request GET \
--url http://localhost:8080/get-storage-info
```

### 3. Subscribe address:  
#### Address 1 (a lot of txes):
```
curl --request GET \
--url 'http://localhost:8080/subscribe?address=0x06450dee7fd2fb8e39061434babcfc05599a6fb8'
```
#### Address 2 (few txes):  
```
curl --request GET \
--url 'http://localhost:8080/subscribe?address=0xb8feffac830c45b4cd210ecdaab9d11995d338ee'
```

### 4. GetTXes:  
#### Address 1 (a lot of txes):
```
curl --request GET \
--url 'http://localhost:8080/get-transactions?address=0x06450dee7fd2fb8e39061434babcfc05599a6fb8'
```

#### Address 2 (few txes):
```
curl --request GET \
--url 'http://localhost:8080/get-transactions?address=0xb8feffac830c45b4cd210ecdaab9d11995d338ee'
```