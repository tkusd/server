# 驗證

## 建立 Token

```
POST /v1/tokens
```

### Request

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`email` | string | Email | **必填**
`password` | string | 密碼 | **必填**

### Response

名稱 | 型別 | 說明
--- | --- | ---
`id` | hash | ID
`user_id` | uuid | 使用者 ID
`created_at` | date | 建立日期

## 使用 Token

大部份的 API 都需要使用 Token 進行驗證，在 Token 建立成功後，把 Token 放在 Header 的 `Authorization` 欄位即可使用。

```
Authorization: Bearer <token>
```

## 刪除 Token

```
DELETE /v1/tokens/:token_id
```