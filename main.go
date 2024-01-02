package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/api"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	loadEnv()

	kubeconfig, err := kubeConfigFlag()
	if err != nil {
		log.Fatal("Failed to load kubeconfig", err.Error())
	}

	client, dynamicClient, err := buildKubeClient(kubeconfig)
	if err != nil {
		log.Fatal("Failed to load kubeconfig", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", api.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/health", api.HealthCheck).Methods(http.MethodGet)

	r.Handle("/api/v1/nodes", C(client, api.GetNodes)).Methods(http.MethodGet)
	r.Handle("/api/v1/namespaces", C(client, api.GetNamespaces)).Methods(http.MethodGet)
	r.Handle("/api/v1/namespaces/status/{namespace}", C(client, api.GetNamespaceStatus)).Methods(http.MethodGet)
	r.Handle("/api/v1/pods/{namespace}", C(client, api.GetPods)).Methods(http.MethodGet)

	r.Handle("/api/v1/namespaces", C(client, api.CreateNamespace)).Methods(http.MethodPost)

	r.Handle("/api/cert-manager/v1/certificates/{namespace}", DC(dynamicClient, client, api.GetCertificates)).Methods(http.MethodGet)
	// r.HandleFunc("/api/cert-manager/v1/certificates", api.CreateCertificate).Methods(http.MethodPost)

	r.HandleFunc("/terraform/cloudflare", api.CreateCloudflareTerraformResource).Methods(http.MethodPost)
	r.HandleFunc("/terraform/cloudflare", api.DeleteCloudflareTerraformResource).Methods(http.MethodDelete)

	server := &http.Server{
		Addr:         os.Getenv("LISTEN_ADDR"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.Default().Handler(r), // Pass our instance of gorilla/mux in.
	}

	log.Println("Server is running on", os.Getenv("LISTEN_ADDR"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func C(client *kubernetes.Clientset, fn func(http.ResponseWriter, *http.Request, *kubernetes.Clientset)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, client)
	})
}

func DC(dynamicClient *dynamic.DynamicClient, client *kubernetes.Clientset, fn func(http.ResponseWriter, *http.Request, *dynamic.DynamicClient, *kubernetes.Clientset)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, dynamicClient, client)
	})
}

func kubeConfigFlag() (string, error) {
	// Getting the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Pointing kubeconfig to the file in the current directory
	kubeconfigFlag := flag.String("kubeconfig", filepath.Join(cwd, "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	return *kubeconfigFlag, nil
}

func buildKubeClient(kubeconfig string) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return client, nil, err
	}

	log.Println("KUBECLIENT: K3s HOST_ADDR:", config.Host)
	log.Println("DYNAMIC_KUBE_CLIENT: K3s HOST_ADDR:", config.Host)
	return client, dynamicClient, nil
}

func loadEnv() {
	_ = godotenv.Load()
	checkList := []string{"LISTEN_ADDR", "GIT_ACCESS_TOKEN"}

	for _, v := range checkList {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set\n", v)
		}
	}
}
