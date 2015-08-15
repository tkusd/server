# 資源

## 建立資源

```
POST /v1/projects/:project_id/assets
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 名稱 |
`description` | string | 描述 |
`data` | multipart | 資料 | **必填**

### Response

``` js
{
    "id": "e6cfe403-20a1-4667-b6eb-de76fb9e6266",
    "name": "140114-0001.png",
    "description": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "created_at": "2015-08-09T02:23:19Z",
    "updated_at": "2015-08-09T02:23:19Z",
    "size": 106840,
    "type": "image/png",
    "slug": "2bdb7769-be94-47f1-8621-d7ae9e4d440b.png",
    "width": 570,
    "height": 451,
    "hash": "0d3cb384ecd3b445110278e1c4028058e1aa27fd"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`description` | string | 描述
`project_id` | uuid | 專案 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`size` | int | 檔案大小 (Bytes)
`type` | string | MIME 類型
`slug` | string | 網址的檔案名稱
`width` | int | 圖片寬度
`height` | int | 圖片高度
`hash` | string | SHA1 Hash

## 取得資源

```
GET /v1/assets/:asset_id
```

### Response

``` js
{
    "id": "e6cfe403-20a1-4667-b6eb-de76fb9e6266",
    "name": "140114-0001.png",
    "description": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "created_at": "2015-08-09T02:23:19Z",
    "updated_at": "2015-08-09T02:23:19Z",
    "size": 106840,
    "type": "image/png",
    "slug": "2bdb7769-be94-47f1-8621-d7ae9e4d440b.png",
    "width": 570,
    "height": 451,
    "hash": "0d3cb384ecd3b445110278e1c4028058e1aa27fd"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`description` | string | 描述
`project_id` | uuid | 專案 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`size` | int | 檔案大小 (Bytes)
`type` | string | MIME 類型
`slug` | string | 網址的檔案名稱
`width` | int | 圖片寬度
`height` | int | 圖片高度
`hash` | string | SHA1 Hash

### 圖片網址

Slug 和資源 ID 沒有關聯性，只是在上傳圖片時隨機產生的 UUID 加上副檔名。

```
/uploads/assets/:slug
```

## 更新使用者

```
PUT /v1/assets/:asset_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 名稱
`description` | string | 描述
`data` | multipart | 資料

### Response

``` js
{
    "id": "e6cfe403-20a1-4667-b6eb-de76fb9e6266",
    "name": "140114-0001.png",
    "description": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "created_at": "2015-08-09T02:23:19Z",
    "updated_at": "2015-08-09T02:23:19Z",
    "size": 106840,
    "type": "image/png",
    "slug": "2bdb7769-be94-47f1-8621-d7ae9e4d440b.png",
    "width": 570,
    "height": 451,
    "hash": "0d3cb384ecd3b445110278e1c4028058e1aa27fd"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`description` | string | 描述
`project_id` | uuid | 專案 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`size` | int | 檔案大小 (Bytes)
`type` | string | MIME 類型
`slug` | string | 網址的檔案名稱
`width` | int | 圖片寬度
`height` | int | 圖片高度
`hash` | string | SHA1 Hash

## 刪除資源

```
DELETE /v1/assets/:asset_id
```

## 取得資源列表

```
GET /v1/projects/:project_id/assets
```