package deploynatsuserrule

import (
	"bytes"
	"fmt"

	"github.com/datasance/iofog-go-sdk/v3/pkg/client"
	"github.com/datasance/potctl/internal/config"
	"github.com/datasance/potctl/internal/execute"
	clientutil "github.com/datasance/potctl/internal/util/client"
	"github.com/datasance/potctl/pkg/util"
	"gopkg.in/yaml.v2"
)

type Options struct {
	Namespace string
	Yaml      []byte
	FullYAML  []byte
	Name      string
}

type executor struct {
	namespace string
	name      string
	fullYAML  []byte
}

func (exe *executor) GetName() string {
	return exe.name
}

func (exe *executor) Execute() error {
	util.SpinStart(fmt.Sprintf("Deploying NATS user rule %s", exe.GetName()))
	clt, err := clientutil.NewControllerClient(exe.namespace)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(exe.fullYAML)
	_, err = clt.UpdateNatsUserRuleFromYAML(exe.name, reader)
	if err == nil {
		return nil
	}
	if _, ok := err.(*client.NotFoundError); !ok {
		return err
	}

	reader = bytes.NewReader(exe.fullYAML)
	_, err = clt.CreateNatsUserRuleFromYAML(reader)
	return err
}

func NewExecutor(opt Options) (execute.Executor, error) {
	ns, err := config.GetNamespace(opt.Namespace)
	if err != nil {
		return nil, err
	}
	controlPlane, err := ns.GetControlPlane()
	if err != nil {
		return nil, err
	}
	if len(controlPlane.GetControllers()) == 0 {
		return nil, util.NewInputError("This namespace does not have a Controller. You must first deploy a Controller before deploying NATS user rules")
	}

	// Validate YAML shape using strict decode against expected top-level fields.
	var validateDoc struct {
		Kind     config.Kind            `yaml:"kind"`
		Metadata config.HeaderMetadata  `yaml:"metadata"`
		Spec     map[string]interface{} `yaml:"spec"`
	}
	if err = yaml.UnmarshalStrict(opt.FullYAML, &validateDoc); err != nil {
		return nil, util.NewUnmarshalError(err.Error())
	}

	return &executor{
		namespace: opt.Namespace,
		name:      opt.Name,
		fullYAML:  opt.FullYAML,
	}, nil
}
