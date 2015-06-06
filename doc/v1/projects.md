# 專案

## 建立專案

```
POST /v1/users/:user_id/projects
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`title` | string | 標題。最長為 255。| **必填**
`description` | string | 描述 | **必填**
`is_private` | boolean | 是否為私人專案 | false

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`title` | string | 標題
`description` | string | 描述
`user_id` | uuid | 使用者 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_private` | boolean | 是否為私人專案

## 取得專案

```
GET /v1/projects/:project_id
```

### Response

``` js
{
    "id": "449e2520-52ec-4cc2-b988-f1f92a0ceeaf",
    "title": "Hello",
    "description": "Test",
    "user_id": "5e7a32d2-80c8-452f-8139-5a860522639f",
    "created_at": "2015-05-29T14:58:56Z",
    "updated_at": "2015-05-29T14:58:56Z",
    "is_private": false,
    "owner": {
        "id": "5e7a32d2-80c8-452f-8139-5a860522639f",
        "name": "John",
        "avatar": "https://www.gravatar.com/avatar/144fa42eb34883ecb00cbc3f81a060a1"
    }
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`title` | string | 標題
`description` | string | 描述
`user_id` | uuid | 使用者 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_private` | boolean | 是否為私人專案
`owner` | object | 擁有者

## 更新專案

```
GET /v1/projects/:project_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`title` | string | 標題。最長為 255。
`description` | string | 描述
`is_private` | boolean | 是否為私人專案
`elements` | []uuid | 子元素

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`title` | string | 標題
`description` | string | 描述
`user_id` | uuid | 使用者 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_private` | boolean | 是否為私人專案

## 刪除專案

```
DELETE /v1/projects/:project_id
```

## 取得專案列表

```
GET /v1/users/:user_id/projects
```

### Request

```
/v1/users/:user_id/projects?limit=30&offset=0&order=-created_at
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`limit` | int | 回傳的物件數量。數值可為 1~100。 | 30
`offset` | int | 從第幾項開始。 | 0
`order` | string | 物件的排序方式。用逗號來區分多個排序條件，預設為升冪排序，在欄位名稱前加上負號則為降冪排序，例如：`title,-created_at`。 | `-created_at`

### Response

``` js
{
    "data": [{
        "id": "96ecd5d4-3294-42bd-9cfb-6ede38576d21",
        "title": "Hello",
        "description": "Test",
        "user_id": "5b7758fd-a408-4e80-9b72-3ff2ebcfad94",
        "created_at": "2015-05-12T16:49:09Z",
        "updated_at": "2015-05-12T16:49:09Z",
        "is_private": false,
        "owner": {
            "id": "5b7758fd-a408-4e80-9b72-3ff2ebcfad94",
            "name": "John",
            "avatar": "https://www.gravatar.com/avatar/144fa42eb34883ecb00cbc3f81a060a1"
        }
    }],
    "has_more": false,
    "count": 1,
    "limit": 30,
    "offset": 0
}
```

名稱 | 型別 | 說明
--- | --- | ---
`data` | []object | 資料陣列
`has_more` | boolean | 是否有更多資料
`count` | int | 資料總數
`limit` | int | 回傳的物件數量
`offset` | int | 從第幾項開始