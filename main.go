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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	loadEnv()

	client, err := buildKubeClient()
	if err != nil {
		log.Fatal("Failed to load kubeconfig", err.Error())
	}

	r := mux.NewRouter()

	r.HandleFunc("/", api.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/health", api.HealthCheck).Methods(http.MethodGet)

	r.Handle("/api/v1/nodes", D(client, api.GetNodes)).Methods(http.MethodGet)
	r.Handle("/api/v1/namespaces", D(client, api.GetNamespaces)).Methods(http.MethodGet)
	r.Handle("/api/v1/namespaces/status/{namespace}", D(client, api.GetNamespaceStatus)).Methods(http.MethodGet)
	r.Handle("/api/v1/pods/{namespace}", D(client, api.GetPods)).Methods(http.MethodGet)

	r.Handle("/api/v1/namespaces/{namespace}", D(client, api.CreateNamespace)).Methods(http.MethodPost)

	r.HandleFunc("/terraform/cloudflare/{name}", api.CreateCloudflareTerraformResource).Methods(http.MethodPost)
	r.HandleFunc("/terraform/cloudflare/{name}", api.DeleteCloudflareTerraformResource).Methods(http.MethodDelete)

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

func D(client *kubernetes.Clientset, fn func(http.ResponseWriter, *http.Request, *kubernetes.Clientset)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, client)
	})
}

func buildKubeClient() (*kubernetes.Clientset, error) {
	var kubeconfig *string

	// homedir
	/*
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
	*/

	// Getting the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Pointing kubeconfig to the file in the current directory
	kubeconfig = flag.String("kubeconfig", filepath.Join(cwd, "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	log.Println("K3s HOST_ADDR:", config.Host)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
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
