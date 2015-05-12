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
`elements` | []uuid | 子元素

#### 更新元素排序

更新元素排序比較複雜，`elements` 中必須包含所有子元素，陣列中可使用字串或是帶有 `id` 和 `elements` 屬性的物件，如下：

``` js
[
    "b91ef654-b81f-4306-87fb-24d27f562b03",
    {
        "id": "c7dfdd95-b3e5-47b1-87cf-24ff50dcc35d",
        "elements": ["96e7f5dd-fbda-4062-8c88-209c4ffb5f9d"]
    }
]
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

## 刪除元素

```
DELETE /v1/elements/:element_id
```

## 取得元素列表

```
GET /v1/projects/:project_id/elements
GET /v1/elements/:element_id/elements
```