package context

import (
	"github.com/codeready-toolchain/registration-service/pkg/application/service"
	"github.com/codeready-toolchain/registration-service/pkg/kubeclient"
	"k8s.io/client-go/rest"
)

type ServiceContextProducer func() ServiceContext

type ServiceContext interface {
	CRTClient() kubeclient.CRTClient
	Services() service.Services
	Config() *rest.Config
}
