# 元素

## 建立元素

```
POST /v1/projects/:project_id/elements
POST /v1/elements/:element_id/elements
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 名稱 |
`type` | string | 類型（見下方） | **必填**
`attributes` | object | 屬性
`is_visible` | boolean | 元素是否可見 | true

### Response

``` js
{
    "id": "eddc9f25-04fc-4ab1-a060-c3f42b77454d",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "element_id": null,
    "index": 1,
    "name": "Home",
    "type": "screen",
    "created_at": "2015-08-08T09:56:00Z",
    "updated_at": "2015-08-08T09:56:00Z",
    "attributes": {},
    "styles": {},
    "events": [],
    "is_visible": true
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`project_id` | uuid | 專案 ID
`element_id` | uuid | 母元素 ID
`name` | string | 名稱
`type` | string | 類型（見下方）
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`attributes` | object | 屬性
`is_visible` | boolean | 元素是否可見
`index` | int | 索引編號

## 取得元素

```
GET /v1/elements/:element_id
```

### Response

``` js
{
    "id": "eddc9f25-04fc-4ab1-a060-c3f42b77454d",
    "project_id": "a5de8ca0-21b8-477a-b8bd-123dbbdb2d17",
    "element_id": null,
    "index": 1,
    "name": "Home",
    "type": "screen",
    "created_at": "2015-08-08T09:56:00Z",
    "updated_at": "2015-08-08T09:56:00Z",
    "attributes": {},
    "styles": {},
    "events": [],
    "is_visible": true
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`project_id` | uuid | 專案 ID
`element_id` | uuid | 母元素 ID
`name` | string | 名稱
`type` | string | 類型（見下方）
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`attributes` | object | 屬性
`is_visible` | boolean | 元素是否可見
`index` | int | 索引編號

## 更新元素

```
PUT /v1/elements/:element_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 名稱
`type` | string | 類型（見下方）
`attributes` | object | 屬性
`elements` | []uuid | 子元素
`is_visible` | boolean | 元素是否可見

## 刪除元素

```
DELETE /v1/elements/:element_id
```

## 取得元素列表

```
GET /v1/projects/:project_id/elements
GET /v1/elements/:element_id/elements
```

### Request

```
/v1/projects/:project_id/elements?flat&depth=1
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`flat` | boolean | 回傳的元素列表不以階層排列 | false
`depth` | int | 列表的最大深度，0 代表不限制 | 0