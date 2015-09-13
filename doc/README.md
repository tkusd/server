# 文件

## 網址

- Base URL: <http://tkusd.zespia.tw>
- API URL: <http://tkusd.zespia.tw/v1>

## 資料

資料一概以 JSON 傳輸，無論是向伺服器請求或是伺服器回傳的資料格式皆為 JSON。

## 目錄

- [使用者](v1/users.md)
- [驗證](v1/tokens.md)
- [專案](v1/projects.md)
- [元素](v1/elements.md)
- [資源](v1/assets.md)
- [事件](v1/events.md)

## JSON-P

在 URL 中加入 `callback` 參數，伺服器會回傳 JSON-P，`meta` 是回應標頭，`data` 是回應內容，如下：

``` js
GET http://tkusd.zespia.tw?callback=foo

foo({
  "meta": {
    "status": 200
  },
  "data": {
    // data
  }
})
```

## 錯誤

當發生錯誤時，你可以使用 `error` 欄位來判斷是否發生錯誤。

``` js
{
  "error": 1000,
  "message": "Unknown error",
  "field": "name"
}
```

欄位 | 說明
--- | ---
`error` | 錯誤代碼（見下方）
`message` | 錯誤訊息
`field` | 欄位（不一定有）

### 1000: 通用錯誤

- 1000: 未知
- 1001: 伺服器錯誤
- 1002: 找不到
- 1003: 存取次數超過限制

### 1100: 資料驗證錯誤

- 1100: 必填欄位
- 1101: 內容類型（Content-Type）錯誤
- 1102: JSON 解析錯誤
- 1103: 欄位類型錯誤
- 1104: Email 格式錯誤
- 1105: 字串長度錯誤
- 1106: URL 格式錯誤
- 1108: UUID 格式錯誤

### 1200: 資源錯誤

- 1200: 找不到使用者
- 1201: 找不到 Token
- 1202: 找不到專案
- 1203: 找不到元素
- 1204: 找不到資源
- 1206: 找不到事件

### 1300: 資料錯誤

- 1300: 密碼錯誤
- 1301: Email 已使用
- 1302: 需要 Token
- 1303: Token 已失效
- 1304: 使用者沒有權限存取資源
- 1307: 元素不被專案擁有
- 1309: 密碼重設密鑰錯誤
- 1310: 密碼重設密鑰已過期（6 小時）
- 1311: 使用者已被啟用
- 1312: 使用者啟用密鑰錯誤