package workmanager_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	wm "github.com/riverchu/workmanager"
)

const (
	DummyWorkerA wm.WorkerName = "worker_a"
	DummyWorkerB wm.WorkerName = "worker_b"

	StepA wm.WorkStep = "step_a"
	StepB wm.WorkStep = "step_b"
)

var (
	dummyBuilder wm.WorkerBuilder = func(context.Context, wm.WorkerName, map[string]interface{}) wm.Worker {
		return &dummyWorker{finish: make(chan struct{}, 1)}
	}
	dummyStepRunner wm.StepRunner = func(work wm.Work, target wm.WorkTarget, nexts ...func(wm.WorkTarget)) {
		_, err := work(target, map[wm.WorkerName]wm.WorkerConfig{
			DummyWorkerA: new(wm.DummyConfig),
			DummyWorkerB: new(wm.DummyConfig),
		})
		if err != nil {
			return
		}
		for _, next := range nexts {
			next(&dummyTarget{token: target.Token(), step: StepB})
		}
	}
	dummyStepProcessor wm.StepProcessor = func(results ...wm.WorkTarget) ([]wm.WorkTarget, error) {
		fmt.Printf("got result: %+v\n", results[0])
		return results, nil
	}
)

type dummyWorker struct{ finish chan struct{} }

func (w *dummyWorker) LoadConfig(wm.WorkerConfig) wm.Worker             { return w }
func (w *dummyWorker) WithContext(context.Context) wm.Worker            { return w }
func (w *dummyWorker) GetContext() context.Context                      { return context.Background() }
func (w *dummyWorker) BeforeWork()                                      {}
func (w *dummyWorker) Work(target wm.WorkTarget) (wm.WorkTarget, error) { return target, nil }
func (w *dummyWorker) AfterWork()                                       { w.finish <- struct{}{} }
func (w *dummyWorker) GetResult() wm.WorkTarget                         { return &dummyTarget{} }
func (w *dummyWorker) Finished() <-chan struct{}                        { return w.finish }
func (w *dummyWorker) Terminate() error                                 { return nil }

type dummyTarget struct {
	token string
	step  wm.WorkStep
}

func (t *dummyTarget) Token() string                                 { return t.token }
func (t *dummyTarget) Key() string                                   { return "" }
func (t *dummyTarget) Trans(step wm.WorkStep) (wm.WorkTarget, error) { return t, nil }
func (t *dummyTarget) ToArray() []wm.WorkTarget                      { return nil }
func (t *dummyTarget) Combine(...wm.WorkTarget) wm.WorkTarget        { return t }
func (t *dummyTarget) TTL() int                                      { return 1 }

func Test_Outline(t *testing.T) {
	wm.RegisterWorker(DummyWorkerA, dummyBuilder)
	wm.RegisterWorker(DummyWorkerB, dummyBuilder)

	wm.RegisterStep(StepA, dummyStepRunner, dummyStepProcessor, StepB)
	wm.RegisterStep(StepB, dummyStepRunner, dummyStepProcessor)

	wm.Serve()

	task := wm.NewTask(context.Background())
	wm.AddTask(task)

	err := wm.Recv(StepA, &dummyTarget{token: task.Token(), step: StepA})
	if err != nil {
		t.Errorf("send target fail: %s", err)
	}

	<-time.NewTimer(3 * time.Second).C

	t.Logf("task %+v", wm.GetTask(task.Token()))
}
