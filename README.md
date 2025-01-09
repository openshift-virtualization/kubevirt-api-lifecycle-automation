# KubeVirt API Lifecycle Automation
This repository contains the code for the KubeVirt API 
lifecycle automation image. The automation is designed 
to handle machine type transitions, automating the 
replacement of machine types with the newest available 
version.

### Functionality
The automation job patches the machine types of VMs, and 
the changes take effect after the next VM restart. 
This is particularly useful when VMs are using outdated 
machine types after a RHEL upgrade.

### Environment Variables
The job exposes several environment variables for customization:

- `MACHINE_TYPE_GLOB` (required): A regex pattern to match the machine types that need to be patched. For example, \*rhel7\* targets RHEL 7 machine types.
- `NAMESPACE` (optional): Specifies the target namespace. If not provided, the job will apply to all namespaces.
- `RESTART_REQUIRED` (optional): A boolean flag indicating whether to restart VMs after the machine type patch is applied. The default value is false.
- `LABEL_SELECTOR` (optional): Allows filtering VMs by labels. If not specified, the job will apply to all VMs. Example usage: `label1 in (value1,value2)`.

### How It Works
The job searches for VMs based on the specified machine 
type regex and optionally filters by namespace and label 
selectors. It patches the machine type of the matching VMs. 
The updates will only be applied after the VM is restarted, 
unless `RESTART_REQUIRED` is set to true, in which case the 
VMs are restarted immediately after the patch.

### Try it
The development tools include a version of kubectl that can be 
used to communicate with the cluster. A wrapper script to 
interact with the cluster can be invoked using ./cluster-up/kubectl.sh.

### Deploy locally
```console
$ mkdir $GOPATH/src/kubevirt.io && cd $GOPATH/src/kubevirt.io
$ git clone git@github.com:openshift-virtualization/kubevirt-api-lifecycle-automation.git && cd kubevirt-api-lifecycle-automation
$ make cluster-up
$ make cluster-sync
$ ./cluster-up/kubectl.sh .....
```

### Deploy on a cluster
```console
$ export KUBECONFIG=</path/to/kubeconfig>
$ make cluster-sync
```

This will create a suspended job that patches the machine types 
of the VMs in the cluster once it is activated.
Modify the environment variables and activate the suspended job by running:
```console
$ ./cluster-up/kubectl.sh edit job kubevirt-api-lifecycle-automation -nkubevirt-api-lifecycle-automation
```