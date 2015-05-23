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
  "id": "2o-j6R88UfFHQuHqiRA8rZQZDnc_-9SlJF3RICNxFag=",
  "user_id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "created_at": "2015-05-08T05:08:06Z"
}
```

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

## 更新  Token

```
PUT /v1/tokens/:token_id
```

## 刪除 Token

```
DELETE /v1/tokens/:token_id
```