# API Docs

See `Harbor`'s API [Here](https://github.com/goharbor/harbor#api).

## Model Management

### Upload Model

**Request**

URL: `POST /api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/upload`

| Path | Note |
|--|--|
| projectName | string, required, Harbor project name |
| modelName | string, required, model name |
| versionName | string, required, version name |

| From Data | Note |
|--|--|
| file | file, required, model file |

| header | Note |
|--|--|
| X-Tenant | string, required, tenant name |
| X-User | string, required, user name |

### Download Model

**Request**

URL: `POST /api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/download`

| Path |  Note |
|--|--|
| projectName | string, required, Harbor project name |
| modelName | string, required, model name |
| versionName | string, required, version name |

| header | Note |
|--|--|
| X-Tenant | string, required, tenant name |
| X-User | string, required, user name |

## Model Extraction / Conversion

### Create `ModelJob`

**Request**

URL: `POST /api/v1alpha1/namespaces/{namespace}/modeljobs`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |

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

### List `ModelJob`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |

### Get `ModelJob`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| modeljobID | string, required, `ModelJob` name |

### Delete `ModelJob`

**Request**

URL: `DELETE /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| modeljobID | string, required, `ModelJob` name |

### Get `ModelJob`'s Event

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/modeljobs/{modeljobID}/events`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| modeljobID | string, required, `ModelJob` name |

## Model Serving

### Create `Serving`

**Request**

URL: `POST /api/v1alpha1/namespaces/{namespace}/servings`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |

<!-- // TODO: -->
**Body**:
```
```

### List `Serving`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/servings`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |

### Get `Serving`

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/servings/{servingID}`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| servingID | string, required, `Serving` name |

### Delete `Serving`

**Request**

URL: `DELETE /api/v1alpha1/namespaces/{namespace}/servings/{servingID}`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| servingID | string, required, `ServingJob` name |

## Other APIs

### Get `Pod`'s Log

**Request**

URL: `GET /api/v1alpha1/namespaces/{namespace}/pods/{podID}/logs`

| Path |  Note |
|--|--|
| namespace | string, required, namespace |
| podID | string, required, `Pod` name |

