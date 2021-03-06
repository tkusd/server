# 使用者

## 建立使用者

```
POST /v1/users
```

### Request

``` js
{
  "name": "abc",
  "email": "abc@example.com",
  "password": "123456",
  "language": "en"
}
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`name` | string | 姓名。最大長度 100。 | **必填**
`email` | string | Email | **必填**
`password` | string | 密碼。長度為 6~50。 | **必填**
`language` | string | 語言。使用 [IETF 語言標籤]，最大長度 35。 | `en`

### Response

``` js
{
  "id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "name": "abc",
  "email": "abc@example.com",
  "avatar": "https://www.gravatar.com/avatar/b28d5fe8da784e36235a487c03a47353",
  "created_at": "2015-05-08T05:04:35Z",
  "updated_at": "2015-05-08T05:04:35Z",
  "is_activated": false,
  "language": "en"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動
`language` | string | 語言

## 取得使用者

```
GET /v1/users/:user_id
```

### Response

包含隱私資訊：

``` js
{
  "id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "name": "abc",
  "email": "abc@example.com",
  "avatar": "https://www.gravatar.com/avatar/b28d5fe8da784e36235a487c03a47353",
  "created_at": "2015-05-08T05:04:35Z",
  "updated_at": "2015-05-08T05:04:35Z",
  "is_activated": false,
  "language": "en"
}
```

不包含隱私資訊：

``` js
{
  "avatar": "https://www.gravatar.com/avatar/b28d5fe8da784e36235a487c03a47353",
  "created_at": "2015-05-08T05:04:35Z",
  "id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "name": "abc",
  "updated_at": "2015-05-08T05:04:35Z",
  "language": "en"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email（僅向本人顯示）
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動（僅向本人顯示）
`language` | string | 語言

## 更新使用者

```
PUT /v1/users/:user_id
```

### Request

``` js
{
  "name": "abc",
  "email": "abc@example.com",
  "password": "234567",
  "old_password": "123456",
  "language": "en"
}
```

參數 | 型別 | 說明
--- | --- | ---
`name` | string | 姓名。最大長度 100。
`email` | string | Email
`password` | string | 新密碼。長度為 6~50。
`old_password` | string | 目前密碼。如果要更改密碼的話必填。
`language` | string | 語言。使用 [IETF 語言標籤]，最大長度 35。 | `en`

### Response

``` js
{
  "id": "cfb4955e-ebdf-4e5b-88f3-6f919dd58907",
  "name": "abc",
  "email": "abc@example.com",
  "avatar": "https://www.gravatar.com/avatar/b28d5fe8da784e36235a487c03a47353",
  "created_at": "2015-05-08T05:04:35Z",
  "updated_at": "2015-05-08T05:04:35Z",
  "is_activated": false,
  "language": "en"
}
```

名稱 | 型別 | 說明
--- | --- | ---
`id` | uuid | ID
`name` | string | 姓名
`email` | string | Email
`avatar` | string | 大頭貼
`created_at` | date | 建立日期
`updated_at` | date | 更新日期
`is_activated` | boolean | 使用者是否已啟動
`language` | string | 語言

## 刪除使用者

```
DELETE /v1/users/:user_id
```

## 啟用使用者

```
POST /v1/activation/:activation_token
```

## 申請重設密碼

```
POST /v1/passwords/reset
```

### Request

``` js
{
  "email": "abc@example.com"
}
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`email` | string | Email | **必填**

## 重設密碼

```
POST /v1/passwords/reset/:password_reset_token
```

``` js
{
  "password": "123456"
}
```

參數 | 型別 | 說明 | 預設值
--- | --- | --- | ---
`password` | string | 密碼。長度為 6~50。 | **必填**

[IETF 語言標籤]: https://en.wikipedia.org/wiki/IETF_language_tag