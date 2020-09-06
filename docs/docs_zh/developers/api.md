# API 文档

`Harbor` 相关接口请参考 [Harbor API 文档](https://github.com/goharbor/harbor#api)

## 模型管理

### 上传模型

**Request**

URL: `POST /api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/upload`

| Path | 说明 |
|--|--|
| projectName | string, required, Harbor 项目名 |
| modelName | string, required, 模型名 |
| versionName | string, required, 模型版本名 |

| From Data | Note |
|--|--|
| file | file, required, 要上传的模型文件 |

| header | Note |
|--|--|
| X-Tenant | string, required, 租户ID |
| X-User | string, required, 发起请求的用户ID |

### 下载模型

**Request**

URL: `POST /api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/download`

| Path | 说明 |
|--|--|
| projectName | string, required, Harbor 项目名 |
| modelName | string, required, 模型名 |
| versionName | string, required, 模型版本名 |

| header | Note |
|--|--|
| X-Tenant | string, required, 租户ID |
| X-User | string, required, 发起请求的用户ID |

## 模型解析 / 转换

### 创建 `ModelJob`

**Request**

URL: `POST /api/v1alpha1/namespaces/{namespace}/modeljobs`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |

**Body**:
```
apiVersion:kleveross.io/v1alpha1
kind: ModelJob
metadata:
  name: modeljob-graphdef-extract
  namespace: kube-system
spec:
  # Add fields here
  model: "harbor.caicloud.com/release/graphdef:v1"
  desiredTag: "harbor.caicloud.com/release/graphdef:v1"
  extraction:
    format: "GraphDef"
  conversion:
     mmdnn:
        from: "MXNETParams"
        to: "ONNX"
```

### 列出 `ModelJob`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |

### 获取 `ModelJob`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| modeljobID | string, required, `ModelJob` 名 |

### 删除 `ModelJob`

**Request**

URL: `DELETE /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| modeljobID | string, required, `ModelJob` 名 |

### 获取 `ModelJob` 的事件

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}/events`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| modeljobID | string, required, `ModelJob` 名 |

## 模型服务

### 创建 `Serving`

**Request**

URL: `POST /api/v1alpha1/namespaces/{namespace}/servings`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |

<!-- // TODO: -->
**Body**:
```
```

### 列出 `Serving`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/servings`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |

### 获取 `Serving`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/servings/{servingID}`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| servingID | string, required, `Serving` 名 |

### 删除 `Serving`

**Request**

URL: `DELETE /api/v1alpha1/namespaces/{namespace}/servings/{servingID}`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| servingID | string, required, `ServingJob` 名 |

## 其他 API

### 获取 `Pod` 日志

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/pods/{podID}/logs`

| Path | 说明 |
|--|--|
| namespace | string, required, 命名空间 |
| podID | string, required, `Pod` 名 |

