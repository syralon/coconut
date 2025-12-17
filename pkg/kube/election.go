package kube

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/syralon/coconut/mesh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type Elector struct {
	id            string
	config        string
	leaseDuration time.Duration
	renewDeadline time.Duration
	retryPeriod   time.Duration

	state atomic.Bool

	cancel context.CancelFunc
}

type ElectorOption func(*Elector)

func WithConfig(config string) ElectorOption {
	return func(e *Elector) {
		e.config = config
	}
}

func WithLeaseDuration(duration time.Duration) ElectorOption {
	return func(e *Elector) {
		e.leaseDuration = duration
	}
}

func WithRenewDeadline(deadline time.Duration) ElectorOption {
	return func(e *Elector) {
		e.renewDeadline = deadline
	}
}

func WithRetryPeriod(retryPeriod time.Duration) ElectorOption {
	return func(e *Elector) {
		e.retryPeriod = retryPeriod
	}
}

func (e *Elector) Start(ctx context.Context, name string, runner mesh.ElectionRunner) error {
	ctx, e.cancel = context.WithCancel(ctx)

	config, err := buildKubeConfig(e.config)
	if err != nil {
		return err
	}
	namespace, err := getNamespace()
	if err != nil {
		return err
	}
	client := clientset.NewForConfigOrDie(config)
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Client:     client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{Identity: e.id},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   e.leaseDuration,
		RenewDeadline:   e.renewDeadline,
		RetryPeriod:     e.retryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				e.state.Store(true)
				runner.Run(ctx)
			},
			OnStoppedLeading: func() {
				if e.state.CompareAndSwap(true, false) {
					runner.Close()
				}
			},
			OnNewLeader: func(_ string) {},
		},
	})
	return nil
}

func (e *Elector) Close() {
	e.cancel()
}

func NewElector(options ...ElectorOption) *Elector {
	e := &Elector{id: uuid.New().String()}
	for _, opt := range options {
		opt(e)
	}
	return e
}
