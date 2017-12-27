package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/laher/kc/internal"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	log.SetFlags(0)
	var (
		fs            = flag.NewFlagSet("kcgp", flag.ExitOnError)
		verbose       = fs.Bool("v", false, "verbose")
		labelSelector = fs.String("l", "", "label selector")
		format        = fs.String("f", "{{.Name}}{{t}}{{.Status.Phase}}{{nl}}", "Representation of each pod")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of kcgp:\n")
		fmt.Fprintf(os.Stderr, " kcgp [contexts] [options]\n")
		fmt.Fprintf(os.Stderr, " [contexts] is a comma-delimited list of contexts, as defined in your kubernetes config.\n")
		fmt.Fprintf(os.Stderr, " [options] defined as follows:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, " For field list see API docs for your version.\n e.g. https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#pod-v1-core\n")
	}
	fs.Parse(args)

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kconfig, err := rules.Load()
	if err != nil {

		log.Printf("Error loading kubectl config: [%s]", err)
		os.Exit(1)
	}

	if len(contexts) == 1 && contexts[0] == "" {
		contexts[0] = kconfig.CurrentContext
	}

	funcMap := template.FuncMap{
		"nl":        func() string { return "\n" },
		"t":         func() string { return "\t" },
		"tableflip": func() string { return "(╯°□°）╯︵ ┻━┻" },
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(*format)
	if err != nil {
		log.Printf("Error parsing template [%s]: [%v]", *format, err)
		os.Exit(1)
	}

	for _, context := range contexts {
		c, ok := kconfig.Contexts[context]
		if !ok {
			log.Printf("Error: context [%s] does not exist", context)
			os.Exit(1)
		}
		if len(contexts) > 1 || *verbose {
			log.Printf("Context: %s, Namespace: %s", context, c.Namespace)
		}

		e, err := getpods(kconfig, context, c.Namespace, *verbose, *labelSelector, tmpl, fs.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func getClientset(kconfig *api.Config, context string) (*kubernetes.Clientset, error) {
	overrides := &clientcmd.ConfigOverrides{CurrentContext: context}
	cconfig := clientcmd.NewDefaultClientConfig(*kconfig, overrides)
	config, err := cconfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getpods(kconfig *api.Config, context string, namespace string, verbose bool, labelSelector string, tmpl *template.Template, args []string) (int, error) {
	clientset, err := getClientset(kconfig, context)
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return 1, err
	}
	if verbose {
		log.Printf("There are %d pods in namespace %s", len(pods.Items), namespace)
	}
	for _, pod := range pods.Items {
		err = tmpl.Execute(os.Stderr, pod)
		if err != nil {
			return 1, err
		}
	}
	return 0, nil
}
