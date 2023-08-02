package main

import (
	"context"

	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/router"
	alabels "github.com/acorn-io/runtime/pkg/labels"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
)

func handlePod(req router.Request, resp router.Response) error {
	pod := req.Object.(*corev1.Pod)

	logrus.WithFields(logrus.Fields{
		"pod.status": pod.Status,
	}).Warnf("handlePod [%s/%s]", pod.Namespace, pod.Name)

	return nil
}

func main() {
	r, err := baaah.DefaultRouter("my-router", scheme.Scheme)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	managedSelector := klabels.SelectorFromSet(map[string]string{
		alabels.AcornManaged: "true",
	})

	r.Type(&corev1.Pod{}).Selector(managedSelector).HandlerFunc(handlePod)

	ctx := context.Background()
	err = r.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	// Background context was canceled. Do whatever cleanup necessary.
}
