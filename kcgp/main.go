package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/laher/kc/internal"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kcv", flag.ExitOnError)
		verbose = fs.Bool("v", false, "verbose")
		//kubeconfig *string
	)
	/*
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
	*/
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)

	// use the current context in kubeconfig
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kconfig, err := rules.Load()

	if err != nil {
		panic(err.Error())
	}
	overrides := &clientcmd.ConfigOverrides{}
	cconfig := clientcmd.NewDefaultClientConfig(*kconfig, overrides)
	config, err := cconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ns, _, err := cconfig.Namespace()
	pods, err := clientset.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in namespace %s\n", len(pods.Items), ns)

	/*
		namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, ns := range namespaces.Items {
			pods, err := clientset.CoreV1().Pods(ns.Name).List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in namespace %s\n", len(pods.Items), ns.Name)
		}
	*/
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}

		e, err := getpods(context, *verbose, fs.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func getpods(context string, verbose bool, args []string) (int, error) {
	kcArgs := []string{"get", "pod"}
	kcArgs = append(kcArgs, args...)
	cmd := kc.PrepKC(context, kcArgs...)
	return kc.Run(cmd, verbose)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
