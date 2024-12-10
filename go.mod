module github.com/kubevirt/kubevirt-api-lifecycle-automation

go 1.22.0

require (
	k8s.io/apimachinery v0.31.2
	k8s.io/code-generator v0.30.4
)

require (
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	kubevirt.io/containerized-data-importer-api v1.57.0-alpha1

)

require (
	github.com/ghodss/yaml v1.0.0
	github.com/golang/mock v1.6.0
	github.com/onsi/ginkgo/v2 v2.19.0
	github.com/onsi/gomega v1.33.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.68.0
	k8s.io/client-go v8.0.0+incompatible
	kubevirt.io/application-aware-quota v1.2.3
	kubevirt.io/qe-tools v0.1.8
	sigs.k8s.io/controller-runtime v0.16.3
)

require (
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/pprof v0.0.0-20240525223248-4bfdf5a9a2af // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/moby/spdystream v0.4.0 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sync v0.7.0 // indirect
	k8s.io/gengo/v2 v2.0.0-20240228010128-51d4e06bde70 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.0
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/api v0.0.0 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/oauth2 v0.18.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.30.4
	k8s.io/apiextensions-apiserver v0.30.4 // indirect
	k8s.io/klog/v2 v2.130.1
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8
	kubevirt.io/api v1.2.0
	kubevirt.io/controller-lifecycle-operator-sdk v0.2.6
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace k8s.io/component-helpers => k8s.io/component-helpers v0.30.4

replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20230406152840-ce21e3fe5da2
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20230324103026-3f1513df25e0
	github.com/openshift/library-go => github.com/mhenriks/library-go v0.0.0-20230310153733-63d38b55bd5a
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	k8s.io/api => k8s.io/api v0.30.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.30.4
	k8s.io/apiserver => k8s.io/apiserver v0.30.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.30.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.30.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.30.4
	k8s.io/component-base => k8s.io/component-base v0.30.4
	k8s.io/cri-api => k8s.io/cri-api v0.30.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.30.4
	k8s.io/endpointslice => k8s.io/staging/src/k8s.io/endpointslice v0.30.4
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.30.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.30.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.30.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.30.4
	k8s.io/kubectl => k8s.io/kubectl v0.30.4
	k8s.io/kubelet => k8s.io/kubelet v0.30.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.30.4
	k8s.io/metrics => k8s.io/metrics v0.30.4
	k8s.io/node-api => k8s.io/node-api v0.30.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.30.4
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.30.4
	k8s.io/sample-controller => k8s.io/sample-controller v0.30.4
	kubevirt.io/qe-tools => kubevirt.io/qe-tools v0.1.8
)

replace k8s.io/controller-manager => k8s.io/controller-manager v0.30.4

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.30.4

replace k8s.io/kms => k8s.io/kms v0.30.4

replace k8s.io/mount-utils => k8s.io/mount-utils v0.30.4

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.30.4

replace kubevirt.io/controller-lifecycle-operator-sdk/api => kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90

replace k8s.io/client-go => k8s.io/client-go v0.30.4
