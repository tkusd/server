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

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`title` | string | 標題
`description` | string | 描述
`user_id` | uuid | 使用者 ID
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_private` | boolean | 是否為私人專案

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