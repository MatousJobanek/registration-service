package workspace

import (
	context2 "context"
	"fmt"
	"strings"

	"github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetHome(cfg *rest.Config, token string) (string, error) {
	// Create a Bearer string by appending string access token
	cl, err := getClient(cfg, fmt.Sprintf("%s/clusters/root/", cfg.Host), token)
	if err != nil {
		return "", err
	}

	homeWorkspace := &v1beta1.Workspace{}
	if err := cl.Get(context2.Background(), types.NamespacedName{Name: "~"}, homeWorkspace); err != nil {
		return "", err
	}
	fmt.Println("home workspace", *homeWorkspace)
	return homeWorkspace.Status.URL, nil
}

func CreateAppStudio(cfg *rest.Config, token string) (string, error) {
	// Create a Bearer string by appending string access token
	home, err := GetHome(cfg, token)
	if err != nil {
		return "", err
	}
	home = strings.ReplaceAll(home, "/clusters/", "/services/workspaces/") + "/personal"
	fmt.Println("for home", home)
	cl, err := getClient(cfg, home, token)
	if err != nil {
		return "", err
	}

	apptudioWorkspace := &v1beta1.Workspace{}
	err = cl.Get(context2.TODO(), types.NamespacedName{Name: "appstudio"}, apptudioWorkspace)
	if err == nil || !errors.IsNotFound(err) {
		return "", err
	}

	apptudioWorkspace = &v1beta1.Workspace{
		ObjectMeta: v1.ObjectMeta{
			Name: "appstudio",
		},
		Spec: v1beta1.WorkspaceSpec{
			Type: v1alpha1.ClusterWorkspaceTypeReference{
				Name: "appstudio",
				Path: "root:plane:usersignup",
			},
		},
	}

	if err := cl.Create(context2.Background(), apptudioWorkspace); err != nil {
		return "", err
	}
	fmt.Println("appstudio workspace", *apptudioWorkspace)
	return apptudioWorkspace.Status.URL, nil
}

func getClient(cfg *rest.Config, url, token string) (client.Client, error) {
	config := rest.CopyConfig(cfg)
	config.Host = url
	config.BearerToken = token
	config.BearerTokenFile = ""
	scheme := runtime.NewScheme()
	if err := v1beta1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	fmt.Println(config)

	cl, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return cl, err
}
