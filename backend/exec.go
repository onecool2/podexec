package backend

import (
	"fmt"
	"io"
	"net/url"

	"k8s.io/kubernetes/pkg/api"
	coreclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
	"k8s.io/kubernetes/pkg/util/interrupt"
	"k8s.io/kubernetes/pkg/util/term"
)

// RemoteExecutor defines the interface accepted by the Exec command - provided for test stubbing
type RemoteExecutor interface {
	Execute(method string, url *url.URL, config *restclient.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool, terminalSizeQueue term.TerminalSizeQueue) error
}

// DefaultRemoteExecutor is the standard implementation of remote command execution
type DefaultRemoteExecutor struct{}

func (*DefaultRemoteExecutor) Execute(method string, url *url.URL, config *restclient.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool, terminalSizeQueue term.TerminalSizeQueue) error {
	// var rConfig restclient.Config
	// CloneValue(config, &rConfig)
	exec, err := remotecommand.NewExecutor(config, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdin:              stdin,
		Stdout:             stdout,
		Stderr:             stderr,
		Tty:                tty,
		TerminalSizeQueue:  terminalSizeQueue,
	})
}

type StreamOptions struct {
	Namespace     string
	PodName       string
	ContainerName string
	Stdin         bool
	TTY           bool
	// minimize unnecessary output
	Quiet bool
	// InterruptParent, if set, is used to handle interrupts while attached
	InterruptParent *interrupt.Handler
	In              io.Reader
	Out             io.Writer
	Err             io.Writer
	Size            *term.Size

	// for testing
	overrideStreams func() (io.ReadCloser, io.Writer, io.Writer)
	isTerminalIn    func(t term.TTY) bool
}

// ExecOptions declare the arguments accepted by the Exec command
type ExecOptions struct {
	StreamOptions

	Command []string

	FullCmdName       string
	SuggestedCmdUsage string

	Executor  RemoteExecutor
	PodClient coreclient.CoreInterface
	Config    *restclient.Config
}

// Complete verifies command line arguments and loads data from the command environment
func (p *ExecOptions) Complete(argsIn []string) error {
	if len(p.PodName) == 0 && (len(argsIn) == 0) {
		return fmt.Errorf("PodName or args is nil")
	}
	if len(p.PodName) != 0 {
		if len(argsIn) < 1 {
			return fmt.Errorf("args is nil")
		}
		p.Command = argsIn
	} else {
		p.PodName = argsIn[0]
		p.Command = argsIn[1:]
		if len(p.Command) < 1 {
			return fmt.Errorf("Command is nil")
		}
	}

	return nil
}

// Validate checks that the provided exec options are specified.
func (p *ExecOptions) Validate() error {
	if len(p.PodName) == 0 {
		return fmt.Errorf("pod name must be specified")
	}
	if len(p.Command) == 0 {
		return fmt.Errorf("you must specify at least one command for the container")
	}
	if p.Out == nil || p.Err == nil {
		return fmt.Errorf("both output and error output must be provided")
	}
	if p.Executor == nil || p.PodClient == nil || p.Config == nil {
		return fmt.Errorf("client, client config, and executor must be provided")
	}
	return nil
}

func (o *StreamOptions) setupTTY() term.TTY {
	t := term.TTY{
		Parent: o.InterruptParent,
		Out:    o.Out,
	}

	if !o.Stdin {
		// need to nil out o.In to make sure we don't create a stream for stdin
		o.In = nil
		o.TTY = false
		return t
	}

	t.In = o.In
	if !o.TTY {
		return t
	}

	// if we get to here, the user wants to attach stdin, wants a TTY, and o.In is a terminal, so we
	// can safely set t.Raw to true
	t.Raw = true

	t.In = o.In
	t.Out = o.Out

	return t
}

// Run executes a validated remote execution against a pod.
func (p *ExecOptions) Run() error {
	pod, err := p.PodClient.Pods(p.Namespace).Get(p.PodName)
	if err != nil {
		return err
	}

	if pod.Status.Phase == api.PodSucceeded || pod.Status.Phase == api.PodFailed {
		return fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}

	containerName := p.ContainerName
	if len(containerName) == 0 {
		if len(pod.Spec.Containers) > 1 {
			usageString := fmt.Sprintf("Defaulting container name to %s.", pod.Spec.Containers[0].Name)
			if len(p.SuggestedCmdUsage) > 0 {
				usageString = fmt.Sprintf("%s\n%s", usageString, p.SuggestedCmdUsage)
			}
			fmt.Fprintf(p.Err, "%s\n", usageString)
		}
		containerName = pod.Spec.Containers[0].Name
	}

	// ensure we can recover the terminal while attached
	t := p.setupTTY()

	var sizeQueue term.TerminalSizeQueue
	if t.Raw {
		if p.Size == nil {
			sizeQueue = t.MonitorSize(t.GetSize())
		} else {
			sizeQueue = t.MonitorSize(p.Size)
		}

		p.Err = nil
	}

	fn := func() error {
		rClient := p.PodClient.RESTClient()
		req := rClient.Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(pod.Namespace).
			SubResource("exec").
			Param("container", containerName)
		req.VersionedParams(&api.PodExecOptions{
			Container: containerName,
			Command:   p.Command,
			Stdin:     p.Stdin,
			Stdout:    p.Out != nil,
			Stderr:    p.Err != nil,
			TTY:       t.Raw,
		}, api.ParameterCodec)
		return p.Executor.Execute("POST", req.URL(), p.Config, p.In, p.Out, p.Err, t.Raw, sizeQueue)
	}

	if err := t.Safe(fn); err != nil {
		return err
	}

	return nil
}
