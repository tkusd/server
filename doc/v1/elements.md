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
`type` | string | 類別 | **必填**
`attributes` | object | 屬性

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`project_id` | uuid | 專案 ID
`element_id` | uuid | 母元素 ID
`name` | string | 名稱
`type` | string | 類別
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`attributes` | object | 屬性

## 取得元素

```
GET /v1/elements/:element_id
```

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`project_id` | uuid | 專案 ID
`element_id` | uuid | 母元素 ID
`name` | string | 名稱
`type` | string | 類別
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`attributes` | object | 屬性

## 更新元素

```
PUT /v1/elements/:element_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 名稱
`type` | string | 類別
`attributes` | object | 屬性
`parent_id` | uuid | 母元素 ID
`elements` | []uuid | 子元素 ID 陣列

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`project_id` | uuid | 專案 ID
`element_id` | uuid | 母元素 ID
`name` | string | 名稱
`type` | string | 類別
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`attributes` | object | 屬性

## 刪除元素

```
DELETE /v1/elements/:element_id
```

## 取得元素列表

```
GET /v1/projects/:project_id/elements
GET /v1/elements/:element_id/elements
```