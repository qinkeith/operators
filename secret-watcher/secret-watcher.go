package main
import (
  "context"
  "flag"
  "fmt"
  "path/filepath"
  "time"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/util/homedir"

  _ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
  // kubeconfig flag
  var kubeconfig *string
  if home := homedir.HomeDir(); home != "" {
    kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) path to the kubeconfig file")
  } else {
    kubeconfig = flag.String("kubeconfig", "", "path to the kubeconfig file")
  }
  flag.Parse()
  // create config
  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
  if err != nil {
    panic(err.Error())
  }
  // create client set
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }
  // watch for secrets
  for {
    secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
      panic(err.Error())
    }
    fmt.Printf("There are %d secrets in the cluster\n", len(secrets.Items))
    time.Sleep(10 * time.Second)
  }
}
