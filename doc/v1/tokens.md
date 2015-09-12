# 驗證

## 建立 Token

```
POST /v1/tokens
```

### Request

``` js
{
  "email": "abc@example.com",
  "password": "123456"
}
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`email` | string | Email | **必填**
`password` | string | 密碼 | **必填**

### Response

``` js
{
  "id": "9354bbb1-2cfd-4808-8a73-e3b03f432cf9",
  "user_id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "secret": "cl7aZacFjkd5aJF7AU3UZU/cfNTTOMIAbyPPM4ws/zA=",
  "created_at": "2015-05-08T05:08:06Z",
  "updated_at": "2015-05-08T05:08:06Z"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`user_id` | uuid | 使用者 ID
`secret` | string | 密鑰，Base 64 格式
`created_at` | date | 建立日期
`updated_at` | date | 更新日期

## 使用 Token

大部份的 API 都需要使用 Token 進行驗證，在 Token 建立成功後，把密鑰放在 Header 的 `Authorization` 欄位即可使用。

```
Authorization: Bearer <secret>
```

## 刪除 Token

```
DELETE /v1/tokens/:token_id
```