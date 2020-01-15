package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/KscSDK/ksc-sdk-go/service/kcm"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv2"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/slb"
	"github.com/KscSDK/ksc-sdk-go/service/sqlserver"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
)

type KsyunClient struct {
	region        string
	eipconn       *eip.Eip
	slbconn       *slb.Slb
	vpcconn       *vpc.Vpc
	kecconn       *kec.Kec
	sqlserverconn *sqlserver.Sqlserver
	krdsconn      *krds.Krds
	kcmconn       *kcm.Kcm
	sksconn       *sks.Sks
	kcsv1conn     *kcsv1.Kcsv1
	kcsv2conn     *kcsv2.Kcsv2
	epcconn       *epc.Epc
	ebsconn       *ebs.Ebs
	mongodbconn   *mongodb.Mongodb
	ks3conn       *s3.S3
}
