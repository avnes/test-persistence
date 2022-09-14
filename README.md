# Test persistence

## Build and deploy

```bash
make build
kubectl apply --kustomize k8s
```

## Create files

```bash
curl -k -X 'POST' \
  'https://rancher-lb/persistence/api/v1/files' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'number_of_files=5' \
  -F 'size_in_bytes=10485760'
```

## List files

```bash
curl -k -X 'GET' \
  'https://rancher-lb/persistence/api/v1/files' \
  -H 'accept: application/json'
```

## Count files

```bash
curl -k -X 'GET' \
  'https://rancher-lb/persistence/api/v1/files/count' \
  -H 'accept: application/json'
```

## Delete files

```bash
curl -k -X 'DELETE' \
  'https://rancher-lb/persistence/api/v1/files' \
  -H 'accept: application/json'
```

## Create random data

```bash
curl -k -X 'POST' \
  'https://rancher-lb/persistence/api/v1/database' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'number_of_inserts=10'
```

## List random data

```bash
curl -k -X 'GET' \
  'https://rancher-lb/persistence/api/v1/database' \
  -H 'accept: application/json'
```

## Count random data

```bash
curl -k -X 'GET' \
  'https://rancher-lb/persistence/api/v1/database/count' \
  -H 'accept: application/json'
```

## Delete random data

```bash
curl -k -X 'DELETE' \
  'https://rancher-lb/persistence/api/v1/database' \
  -H 'accept: application/json'
```
