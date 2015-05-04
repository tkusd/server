# 使用者

## 建立使用者

```
POST /v1/users
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 姓名。最大長度 100。 | **必填**
`email` | string | Email | **必填**
`password` | string | 密碼。長度為 6~50。 | **必填**

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動

## 取得使用者

```
GET /v1/users/:user_id
```

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email（僅向本人顯示）
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動（僅向本人顯示）

## 更新使用者

```
PUT /v1/users/:user_id
```

### Request

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 姓名。最大長度 100。
`email` | string | Email
`password` | string | 新密碼。長度為 6~50。
`old_password` | string | 目前密碼。如果要更改密碼的話必填。

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動

## 刪除使用者

```
DELETE /v1/users/:user_id
```