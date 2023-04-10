package main

import (
	"context"

	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/apply"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	r, err := baaah.DefaultRouter("my-router", scheme.Scheme)
	if err != nil {
		panic(err)
	}

	r.Type(new(appsv1.Deployment)).HandlerFunc(handleDeployment)
	r.HandleFunc(new(appsv1.Deployment), annotateDeployments)

	ctx := context.Background()
	err = r.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	// Background context was canceled. Do whatever cleanup necessary.
}

func annotateDeployments(req router.Request, resp router.Response) error {
	// Act on the deployment, which can be retrieved via req.Object
	// If you need to create new objects based on the deployment, use resp.Objects()
	// A req.Client and req.Ctx are provided for instances when you need to directly interact with the Kubernetes API.
	deployment := req.Object.(*appsv1.Deployment)
	if deployment.Annotations["foo"] == "baz" {
		return nil
	}
	logrus.Warnf("annotateDeployment [%s/%s]: %s", deployment.Namespace, deployment.Name, deployment.ResourceVersion)

	deployment.Annotations["foo"] = "baz"

	// logrus.Warnf("%s/%s.metadata.annotations: %s", deployment.Namespace, deployment.Name, deployment.Annotations)
	// logrus.Warnf("%s/%s.metadata.annotations: %s", deployment.Namespace, deployment.Name, deployment.Annotations)

	return apply.Ensure(req.Ctx, req.Client, deployment)
}

func handleDeployment(req router.Request, resp router.Response) error {
	// Act on the deployment, which can be retrieved via req.Object
	// If you need to create new objects based on the deployment, use resp.Objects()
	// A req.Client and req.Ctx are provided for instances when you need to directly interact with the Kubernetes API.
	// logrus.Warnf("handleDeployment [%s]", req.Object.GetName())
	deployment := req.Object.(*appsv1.Deployment)
	// logrus.Warnf("handleDeployment [%s/%s]: %s", deployment.Namespace, deployment.Name, deployment.ResourceVersion)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name + "-key",
			Namespace: deployment.Namespace,
			Annotations: map[string]string{
				"foo": "bar",
			},
		},
		Data: map[string][]byte{
			"token": []byte("secret-key"),
		},
	}

	// logrus.Warnf("%s/%s.metadata.annotations: %s", deployment.Namespace, deployment.Name, deployment.Annotations)
	// logrus.Warnf("%s/%s.metadata.annotations: %s", deployment.Namespace, deployment.Name, deployment.Annotations)

	resp.Objects(secret)

	return nil
}
