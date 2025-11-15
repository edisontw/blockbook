# Blockbook + PepePow 指南

此文件說明如何在本倉庫的 PepePow 整合上完成後端部署、建置 Blockbook，以及測試同步流程。所有命令與範例皆以 Debian 11/12 AMD64 為基準，若使用其他發行版請依需求調整。

## 1. 環境需求

- Debian 11/12 或相容 Linux，具備 `sudo` 權限
- Go 1.23（可由 [go.dev](https://go.dev/dl/) 下載）與 `git`
- RocksDB 與 ZeroMQ 執行期：`sudo apt install build-essential libzmq3-dev zlib1g-dev libsnappy-dev`
- PePe-core 後端：v2.8.1.1 以上，下載自 <https://github.com/MattF42/PePe-core/releases>
- 建議硬體：8 核 CPU、16 GB RAM、500 GB SSD（完整節點同步仍以鏈狀況為準）

> Blockbook 本身可使用 `make` 目標在 Docker 內建置，但 PepePow 的後端節點務必在宿主 OS 原生執行，以便提供 RPC 與 ZMQ 服務給 Blockbook。

## 2. 部署 PepePow 後端

1. 建立目錄與系統使用者
   ```bash
   sudo useradd --system --home /var/lib/pepepow --shell /usr/sbin/nologin pepepow
   sudo mkdir -p /opt/pepepow/bin /var/lib/pepepow/backend
   sudo chown -R pepepow:pepepow /opt/pepepow /var/lib/pepepow
   ```

2. 下載並驗證二進位檔（以官方 ARM 套件為例，若需 AMD64 請改用自行編譯的檔案）
   ```bash
   cd /opt/pepepow
   sudo -u pepepow curl -L -o pepepow.tgz \
     https://github.com/MattF42/PePe-core/releases/download/v2.8.1.1/PEPEPOW-v2.8.1.1-40b8862-release-aarch64-linux-gnu.tgz
   echo "b9e975faef36f1ce07b19ab3beb8eace499d1d54f258a98face91812c4700449  pepepow.tgz" | sha256sum -c -
   sudo -u pepepow tar -C /opt/pepepow/bin --strip-components=1 -xzf pepepow.tgz
   ```

3. 建立 `pepepow.conf`（對應 `configs/coins/pepepow.json` 中的預設值）
   ```bash
   sudo tee /opt/pepepow/pepepow.conf >/dev/null <<'EOF'
   daemon=1
   server=1
   mainnet=1
   txindex=1
   rpcuser=rpc
   rpcpassword=rpc
   rpcallowip=127.0.0.1
   rpcbind=127.0.0.1
   rpcport=8093
   zmqpubhashtx=tcp://127.0.0.1:38393
   zmqpubhashblock=tcp://127.0.0.1:38393
   whitelist=127.0.0.1
   maxconnections=64
   EOF
   sudo chown pepepow:pepepow /opt/pepepow/pepepow.conf
   ```

4. 啟動守護行程（一次性測試）
   ```bash
   sudo -u pepepow /opt/pepepow/bin/pepepowd \
     -conf=/opt/pepepow/pepepow.conf \
     -datadir=/var/lib/pepepow/backend
   ```
   確認 `pepepow-cli -conf=/opt/pepepow/pepepow.conf getblockchaininfo` 可正常回應後，再寫入 systemd 服務。

## 3. 建置 Blockbook（PepePow）

1. 取得程式碼並建置二進位
   ```bash
   git clone https://github.com/trezor/blockbook.git
   cd blockbook
   make build
   # 產物：build/bin/blockbook
   ```

2. 產生 `blockchaincfg.json`（也可直接使用 `configs/coins/pepepow.json`）
   ```bash
   ./contrib/scripts/build-blockchaincfg.sh pepepow
   # 會將 configs/coins/pepepow.json 轉為 build/blockchaincfg.json
   ```

3. 建立資料與日誌目錄
   ```bash
   sudo useradd --system --home /var/lib/blockbook --shell /usr/sbin/nologin blockbook || true
   sudo mkdir -p /var/lib/blockbook/pepepow /var/log/blockbook
   sudo chown -R blockbook:blockbook /var/lib/blockbook /var/log/blockbook
   ```

4. 啟動 Blockbook
   ```bash
   sudo -u blockbook ./build/bin/blockbook \
     -sync \
     -coin=PEPEPOW \
     -blockchaincfg=build/blockchaincfg.json \
     -datadir=/var/lib/blockbook/pepepow \
     -certfile=server/testcert \
     -logfile=/var/log/blockbook/pepepow.log \
     -workers=4 \
     -cache=256
   ```
   - `-sync` 會在啟動時執行區塊導入；移除即可啟用 websocket/API 服務。
   - `-workers` 與 `-cache` 可依硬體調整。若記憶體不足，可加上 `-workers=1 -dbcache=0`。
   - 內外部介面預設來自 `configs/coins/pepepow.json`，可藉由 `-public` / `-internal` 參數覆寫。

5. 驗證
   - HTTP：`curl http://127.0.0.1:9193/api` 應返回 Blockbook API 版本。
   - WebSocket：`wscat -c ws://127.0.0.1:9193/websocket`，訂閱 `mempool-subscribe` 可看到交易推播。
   - 指標：`curl http://127.0.0.1:9093/metrics` 取得 Prometheus 指標。

## 4. systemd 範例

`/etc/systemd/system/pepepowd.service`
```ini
[Unit]
Description=PepePow daemon
After=network.target

[Service]
User=pepepow
Group=pepepow
ExecStart=/opt/pepepow/bin/pepepowd -conf=/opt/pepepow/pepepow.conf -datadir=/var/lib/pepepow/backend
ExecStop=/opt/pepepow/bin/pepepow-cli -conf=/opt/pepepow/pepepow.conf stop
Restart=on-failure
TimeoutStopSec=180

[Install]
WantedBy=multi-user.target
```

`/etc/systemd/system/blockbook-pepepow.service`
```ini
[Unit]
Description=Blockbook PepePow indexer
After=network-online.target pepepowd.service
Wants=pepepowd.service

[Service]
User=blockbook
Group=blockbook
WorkingDirectory=/opt/blockbook
ExecStart=/opt/blockbook/build/bin/blockbook \
  -sync \
  -coin=PEPEPOW \
  -blockchaincfg=/opt/blockbook/build/blockchaincfg.json \
  -datadir=/var/lib/blockbook/pepepow \
  -logfile=/var/log/blockbook/pepepow.log
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

啟用服務：
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pepepowd.service blockbook-pepepow.service
```

## 5. 疑難排解

| 症狀 | 可能原因與處理方式 |
| --- | --- |
| `blockbook` 啟動後立即退出並顯示 `database is in inconsistent state` | 初次同步被中止。刪除 `/var/lib/blockbook/pepepow/db` 後重新執行 `-sync`；同步時建議 `-workers=1` 以降低記憶體需求。 |
| API 返回空值或地址無法辨識 | 確認使用的地址前綴為 `P...`（主網）或 `y...`（測試網），Blockbook 會根據 `pepepowparser.go` 中的魔術數進行驗證。 |
| `zmq` 連線錯誤 | 檢查 `pepepow.conf` 是否啟用了 `zmqpubhashtx`/`zmqpubhashblock`，以及防火牆是否放行 `38393`。 |
| `pepepowd` RPC 認證失敗 | Blockbook 與 backend 必須共用 `rpcuser`/`rpcpassword`；若改成其他值，請同步更新 `configs/coins/pepepow.json` 或 `build/blockchaincfg.json`。 |

## 6. 更新與版本控制

- PepePow 的 `slip44` 設為 `5`，BIP32/44 導出路徑請使用 `m/44'/5'/...`。
- 當 PePe-core 發布新版本，更新流程：
  1. 調整 `configs/coins/pepepow.json` 中的 `version`、`binary_url` 與 `verification_source`。
  2. 重新執行 `./contrib/scripts/build-blockchaincfg.sh pepepow` 以生成新的 `blockchaincfg.json`。
  3. 重新建置並重新啟動 Blockbook。

有任何疑問可於 Issues 提出，並附上 `blockbook` 與 `pepepowd` 的日誌片段，便於協助排查。
