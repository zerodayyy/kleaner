package main

import (
	"context"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
)

var args struct {
	Debug      bool   `arg:"-d" help:"Enable debug logs"`
	Namespaces string `arg:"-n" help:"List of namespaces to clean pods in, separated by comma"`
}

func main() {
	arg.MustParse(&args)
	log := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &prefixed.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	if args.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	log.Info("Kleaner v1.0.0 is starting...")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error("Unable to connect to Kubernetes")
		log.Debug(err.Error())
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Unable to connect to Kubernetes")
		log.Debug(err.Error())
		os.Exit(1)
	}

	var namespaces []string
	if ns := strings.Split(args.Namespaces, ","); len(ns) == 0 {
		namespaces = []string{metav1.NamespaceAll}
	} else {
		namespaces = ns
	}

	for _, ns := range namespaces {
		if ns != metav1.NamespaceAll {
			log.Info("Cleaning namespace", ns)
		} else {
			log.Info("Cleaning all namespaces")
		}

		pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Error("Unable to list pods")
			log.Debug(err.Error())
			os.Exit(1)
		}

		for _, pod := range pods.Items {
			strings.Contains(pod.Status.Reason, "Evicted")
			log.Info("Pod", pod.Name, "will be deleted")
			err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				log.Error("Unable to delete pod", pod.Name)
				log.Debug(err.Error())
				os.Exit(1)
			}
			log.Info("Pod", pod.Name, "deleted")
		}
	}
}
