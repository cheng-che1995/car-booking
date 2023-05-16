# Car Booking System
## 專案目標：
~~為了促進家庭和諧~~

為了方便協調我跟我爸用車日期及時間而開發的後端API。

## 功能：
* **使用者資訊管理**：提供使用者資訊新增、刪除、查詢、修改的功能。
* **車輛資訊管理**：提供車輛資訊新增、刪除、查詢、修改的功能。
* **預約功能**：允許使用者選擇車輛和時間進行預約。
* **預約記錄管理**：允許使用者刪除自己建立的預約記錄。
* **預約記錄歷史查詢**：允許使用者查詢過去車輛的使用紀錄。

## 技術堆疊：
* 採用[Echo](https://github.com/labstack/echo) 網頁框架，並導入 JWT 驗證方法，提供安全的 API 服務。
* 使用基於HTTP方法的RESTful API實作，提供標準的API介面。
* 使用[Hoppscotch](https://docs.hoppscotch.io/)執行API測試，確保 API 的正確性和可靠性。
* ***尚未更新*** ~~使用Golang原生數據庫[BoltDB](https://github.com/boltdb/bolt) 實現輕量級的鍵值儲存，確保高效且可靠的數據操作。~~
* 導入[MySQL資料庫](https://github.com/go-sql-driver/mysql) 實現較複雜的數據操作，例如 CRUD（新增、讀取、更新、刪除）等。

## 如何開始：
### 環境要求
* Golang v1.19
* MySQL v8.0.31

如果尚未安裝上述環境，請參考以下官方網站進行安裝：
* [Golang官方網站](https://go.dev/doc/install)
* [MySQL官方網站](https://dev.mysql.com/downloads/installer/)進行安裝。



* 下載專案
在終端機中下載此專案至本地端：
```shell
git clone https://github.com/diverwil1995/car-booking.git
```
* 建立資料庫
在終端機中使用MYSQL命令列手動建立名為"testdb"的Database，命令如下：

```mysql
create database testdb;
```
> "請自行修改"testdb"為所需的資料庫名稱，並同時更新mysql_repository.go檔案中的環境變數。"

* 運行程式
執行以下命令，將應用程式運行"localhost:1323"上：

```golang
go run .
```
> 若需要查看更詳細的 API 文件，請點擊[這裡](https://)。該文件包含所有可用的 API 端點、請求和回應的範例，以及相關的參數和資源說明。請參考該文件以便深入了解並開始使用專案的 API 功能。


## 貢獻者名單：
* Felix Werth [@qwwqe](https://github.com/qwwqe)