# Tekton CI Golang SDK
[Tekton](https://tekton.dev/) 是一个云原生的 CI 平台，本项目提供 Tekton 的 Golang SDK。

## 用法
引入包
```go
import (
  tekton "github.com/hongyuxuan/tekton-sdk-go"
  "github.com/hongyuxuan/tekton-sdk-go/core/option"
)
```

有两种初始化tekton连接的方式

- 通过 kubeconfig：
```go
client := tekton.NewClient(
  option.WithKubeconfig("./kubeconfig"),
  option.WithSecretPrefix("default-token"),
  option.WithDebug(true),
)
```
SDK 通过 `kubernetes.io/service-account-token` 类型的 secrets 绑定的 token 作为 `Bearer Token` 向 Kubernetes 里的 Tekton CRD 资源发送请求，因此需通过 `WithSecretPrefix("default-token")` 指定该 secrets 的前缀。

- 程序直接运行在 Kubernetes 中，此时无需 `kubeconfig`，直接初始化：
```go
client := tekton.NewClient(
  option.WithSecretPrefix("default-token"),
  option.WithDebug(true),
)
```

## 示例
列出所有 `Pipeline`
```go
res, err := client.Pipeline(namespace).List(context.TODO(), metav1.ListOptions{
  LabelSelector: "app=testpipeline",
  Limit:         3,
})
if err != nil {
  panic(err)
}
fmt.Println(res)
```
更多示例详见test。