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
	"k8s.io/client-go/util/homedir"
)

func main() {
	_ = godotenv.Load()
	client, err := buildKubeClient()
	if err != nil {
		log.Fatal("Failed to load kubeconfig", err.Error())
	}

	r := mux.NewRouter()

	r.HandleFunc("/", api.HealthCheck).Methods("GET")
	r.HandleFunc("/api/v1/health", api.HealthCheck).Methods("GET")

	r.Handle("/api/v1/namespaces", D(client, api.GetNamespaces)).Methods("GET")
	r.Handle("/api/v1/pods/{namespace}", D(client, api.GetPods)).Methods("GET")

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
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
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
