# 事件

## 建立事件

```
POST /v1/elements/:element_id/events
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 名稱 |
`action_id` | uuid | 行動 ID | **必填**
`event` | string | 事件 | **必填**

### Response

``` js
{
    "id": "10615bd3-b345-48de-af77-d15b0ae9a6cb",
    "element_id": "5923ae20-ee3d-49f7-a0b3-d6771fcc93d5",
    "action_id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "event": "click",
    "created_at": "2015-08-15T09:28:42Z",
    "updated_at": "2015-08-15T09:28:42Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`element_id` | uuid | 元素 ID
`action_id` | uuid | 行動 ID
`event` | string | 事件
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 取得事件

```
GET /v1/events/:event_id
```

### Response

``` js
{
    "id": "10615bd3-b345-48de-af77-d15b0ae9a6cb",
    "element_id": "5923ae20-ee3d-49f7-a0b3-d6771fcc93d5",
    "action_id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "event": "click",
    "created_at": "2015-08-15T09:28:42Z",
    "updated_at": "2015-08-15T09:28:42Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`element_id` | uuid | 元素 ID
`action_id` | uuid | 行動 ID
`event` | string | 事件
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 更新事件

```
PUT /v1/events/:event_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 名稱
`action_id` | uuid | 行動 ID
`event` | string | 事件

### Response

``` js
{
    "id": "10615bd3-b345-48de-af77-d15b0ae9a6cb",
    "element_id": "5923ae20-ee3d-49f7-a0b3-d6771fcc93d5",
    "action_id": "2a9b2fb1-0a70-430d-8a41-2fa727789bbd",
    "event": "click",
    "created_at": "2015-08-15T09:28:42Z",
    "updated_at": "2015-08-15T09:28:42Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`element_id` | uuid | 元素 ID
`action_id` | uuid | 行動 ID
`event` | string | 事件
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 刪除事件

```
DELETE /v1/events/:event_id
```

## 取得事件列表

```
GET /v1/elements/:element_id/events
```