package args

// FactoryArgs contains the required parameters to generate all namespaced resources
type FactoryArgs struct {
	MachineTypeGlob        string `required:"true" split_words:"true"`
	TargetNamespace        string `required:"true" split_words:"true"`
	RestartRequired        string `required:"true" split_words:"true"`
	OperatorVersion        string `required:"true" split_words:"true"`
	KubevirtApiLifecycleAutomationImage       string `required:"true" split_words:"true"`
	DeployClusterResources string `required:"true" split_words:"true"`
	DeployPrometheusRule   string `required:"true" split_words:"true"`
	Verbosity              string `required:"true"`
	LabelSelector          string `required:"true"`
	PullPolicy             string `required:"true" split_words:"true"`
	Namespace              string
}
