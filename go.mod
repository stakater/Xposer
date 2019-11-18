module github.com/stakater/Xposer

go 1.13

require (
	github.com/PuerkitoBio/purell v1.0.0
	github.com/PuerkitoBio/urlesc v0.0.0-20160726150825-5bd2802263f2
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/go-restful v1.1.4-0.20170410110728-ff4f55a20633
	github.com/fatih/structs v1.0.0
	github.com/ghodss/yaml v0.0.0-20150909031657-73d445a93680
	github.com/go-openapi/jsonpointer v0.0.0-20160704185906-46af16f9f7b1
	github.com/go-openapi/jsonreference v0.0.0-20160704190145-13c6e3589ad9
	github.com/go-openapi/spec v0.0.0-20170914061247-7abd5745472f
	github.com/go-openapi/swag v0.0.0-20170606142751-f3f9494671f9
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c
	github.com/google/gofuzz v1.0.0
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
	github.com/gregjones/httpcache v0.0.0-20170728041850-787624de3eb7
	github.com/hashicorp/golang-lru v0.5.1
	github.com/howeyc/gopass v0.0.0-20170109162249-bf9dde6d0d2c
	github.com/imdario/mergo v0.3.5
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/json-iterator/go v1.1.7
	github.com/juju/ratelimit v0.0.0-20170523012141-5b9ff8664717
	github.com/mailru/easyjson v0.0.0-20170624190925-2f5df55504eb
	github.com/openshift/api v3.9.1-0.20180801171038-322a19404e37+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/sirupsen/logrus v1.0.5
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f
	golang.org/x/text v0.3.2
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/inf.v0 v0.9.0
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go v6.0.1-0.20180103015815-9389c055a838+incompatible
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
)

replace (
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20190923092516-169848dd8137+incompatible // prebase-1.16
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190923092832-6afefc9bb372 // prebase-1.16
	k8s.io/api => k8s.io/api v0.0.0-20191004120104-195af9ec3521 // release-1.16
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8 // kubernetes-1.16.0
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0
)
