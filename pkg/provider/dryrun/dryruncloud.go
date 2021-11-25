package dryrun

import (
	log "github.com/sirupsen/logrus"
	prvd "k8s.io/cloud-provider-alibaba-cloud/pkg/provider"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/base"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/ecs"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/pvtz"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/slb"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/vpc"
)

func NewDryRunCloud() prvd.Provider {
	auth, err := base.NewClientMgr()
	if err != nil {
		log.Warnf("initialize alibaba cloud client auth: %s", err.Error())
	}
	if auth == nil {
		panic("auth should not be nil")
	}
	err = auth.Start(base.RefreshToken)
	if err != nil {
		log.Warnf("refresh token: %s", err.Error())
	}

	cloud := &alibaba.AlibabaCloud{
		IMetaData:    auth.Meta,
		EcsProvider:  ecs.NewEcsProvider(auth),
		SLBProvider:  slb.NewLBProvider(auth),
		PVTZProvider: pvtz.NewPVTZProvider(auth),
		VPCProvider:  vpc.NewVPCProvider(auth),
	}

	return &DryRunCloud{
		IMetaData:  auth.Meta,
		DryRunECS:  NewDryRunECS(auth, cloud.EcsProvider),
		DryRunPVTZ: NewDryRunPVTZ(auth, cloud.PVTZProvider),
		DryRunVPC:  NewDryRunVPC(auth, cloud.VPCProvider),
		DryRunSLB:  NewDryRunSLB(auth, cloud.SLBProvider),
	}
}

var _ prvd.Provider = &DryRunCloud{}

type DryRunCloud struct {
	*DryRunECS
	*DryRunPVTZ
	*DryRunVPC
	*DryRunSLB
	prvd.IMetaData
}
