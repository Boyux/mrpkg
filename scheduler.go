package mrpkg

import (
	"context"
	"reflect"
	"runtime"
	"sync"
)

var (
	numcpu    = runtime.NumCPU()
	multiple  = 2
	scheduler = &Scheduler{N: numcpu * multiple}
)

func Run(fn any, args ...any) {
	scheduler.Run(fn, args...)
}

func RunTask(task Task) {
	scheduler.RunTask(task)
}

func RunTaskWithPriority(task Task, level Level) {
	scheduler.RunTaskWithPriority(task, level)
}

type Level uint

func (lv Level) Is(target Level) bool {
	return lv == target
}

const (
	LevelHigh Level = iota
	LevelMid
	LevelLow

	LevelTop Level = 0xF690951D
)

type Task interface {
	Run()
}

type PriorityTask interface {
	Task
	Priority() Level
}

type TaskFunc func()

func (f TaskFunc) Run() {
	f()
}

func (f TaskFunc) Priority() Level {
	return LevelLow
}

type Scheduler struct {
	N int

	token Token
	queue Queue[priority]
	mu    sync.Mutex

	ctx       context.Context
	ctxCancel func()
}

func (scheduler *Scheduler) init() {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	if scheduler.token == nil {
		scheduler.token = make(Token)
		if scheduler.queue == nil {
			scheduler.queue = &Heap[priority]{
				inner: make(tinyHeap[priority], 0, Max(scheduler.N, 4)),
			}
		}
		scheduler.ctx, scheduler.ctxCancel = context.WithCancel(context.Background())
		for i := 0; i < Max(scheduler.N, 1); i++ {
			go start(scheduler.ctx, scheduler.token, &scheduler.mu, scheduler.queue)
		}
	}
}

func start(ctx context.Context, token Token, lock sync.Locker, queue Queue[priority]) {
	var task Task

loop:
	for {
		task = nil

		if closed := token.Get(ctx); closed {
			break loop
		}

		lock.Lock()
		if queue.Len() > 0 {
			task = queue.Pop()
		}
		lock.Unlock()

		if task != nil {
			task.Run()
		}
	}
}

type reflectTask struct {
	fn   reflect.Value
	args []reflect.Value
}

func (task *reflectTask) Run() {
	if ty := task.fn.Type(); ty.IsVariadic() {
		task.fn.CallSlice(task.args)
	} else {
		task.fn.Call(task.args)
	}
}

func (scheduler *Scheduler) Run(fn any, args ...any) {
	task := reflectTask{
		fn:   reflect.ValueOf(fn),
		args: Map(args, reflect.ValueOf),
	}

	scheduler.RunTask(&task)
}

func (scheduler *Scheduler) RunTask(task Task) {
	var pTask PriorityTask
	if v, ok := task.(PriorityTask); ok {
		pTask = v
	} else {
		pTask = lowLevel(task)
	}

	if pTask.Priority().Is(LevelTop) {
		go pTask.Run()
		return
	}

	scheduler.init()
	scheduler.mu.Lock()
	scheduler.queue.Push(priority{pTask})
	scheduler.mu.Unlock()
	scheduler.token.Send(scheduler.ctx)
}

func (scheduler *Scheduler) RunTaskWithPriority(task Task, lv Level) {
	scheduler.RunTask(NewPriorityTask(task, lv))
}

func (scheduler *Scheduler) Stop() {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	if scheduler.ctxCancel != nil {
		scheduler.ctxCancel()
		scheduler.ctx = nil
		scheduler.ctxCancel = nil
		if scheduler.token != nil {
			close(scheduler.token)
			scheduler.token = nil
		}
		if scheduler.queue != nil {
			for scheduler.queue.Len() > 0 {
				scheduler.queue.Pop()
			}
			scheduler.queue = nil
		}
	}
}

type Token chan struct{}

func (token Token) Send(ctx context.Context) {
	select {
	case <-ctx.Done():
		break
	case token <- struct{}{}:
		return
	}
}

func (token Token) Get(ctx context.Context) (closed bool) {
	select {
	case <-ctx.Done():
		return true
	case _, ok := <-token:
		return !ok
	}
}

type priority struct {
	PriorityTask
}

func (p priority) Lt(rhs priority) bool {
	return p.Priority() < rhs.Priority()
}

type low struct {
	Task
}

func (l *low) Priority() Level {
	return LevelLow
}

func lowLevel(task Task) PriorityTask {
	return &low{task}
}

func NewPriorityTask(task Task, lv Level) PriorityTask {
	return &taskWithPriority{
		Task: task,
		lv:   lv,
	}
}

type taskWithPriority struct {
	Task
	lv Level
}

func (task *taskWithPriority) Priority() Level {
	return task.lv
}
