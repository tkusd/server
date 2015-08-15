# 行動

## 建立行動

```
POST /v1/projects/:project_id/actions
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 名稱 |
`action` | string | 行動 | **必填**
`data` | multipart | 資料 | **必填**

### Response

``` js
{
    "id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "name": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "action": "alert",
    "data": {},
    "created_at": "2015-08-15T09:21:37Z",
    "updated_at": "2015-08-15T09:21:37Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`project_id` | uuid | 專案 ID
`action` | string | 行動
`data` | object | 資料
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 取得行動

```
GET /v1/actions/:action_id
```

### Response

``` js
{
    "id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "name": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "action": "alert",
    "data": {},
    "created_at": "2015-08-15T09:21:37Z",
    "updated_at": "2015-08-15T09:21:37Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`project_id` | uuid | 專案 ID
`action` | string | 行動
`data` | object | 資料
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 更新行動

```
PUT /v1/actions/:action_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 名稱
`action` | string | 行動
`data` | multipart | 資料

### Response

``` js
{
    "id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "name": "",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "action": "alert",
    "data": {},
    "created_at": "2015-08-15T09:21:37Z",
    "updated_at": "2015-08-15T09:21:37Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 名稱
`project_id` | uuid | 專案 ID
`action` | string | 行動
`data` | object | 資料
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 刪除行動

```
DELETE /v1/actions/:action_id
```

## 取得行動列表

```
GET /v1/projects/:project_id/actions
```