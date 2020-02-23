package server

import (
	"html/template"

	"github.com/ian-ross/morse-blinkies/mbaas/model"
	"github.com/ian-ross/morse-blinkies/mbaas/processing"
)

// Notification is a job status notification that should be one of the
// below ...Notification types.
type Notification interface{}

// QueueNotification is an integer queue position notification.
type QueueNotification struct {
	Position int `json:"queue_position"`
}

// URLNotification is a URL string notification for a completed job.
type URLNotification struct {
	URL string `json:"url"`
}

// ErrorNotification is a string error message notification arising
// from a problem during blinky processing.
type ErrorNotification struct {
	ErrorMessage string `json:"error_message"`
}

// MakeQueuer sets up the job queuer.
func MakeQueuer(bm *processing.BlinkyMaker, tmpl *template.Template) *Queuer {
	inCh := make(chan *job)
	workInCh, workOutCh := startWorker(bm, tmpl)
	subCh := startQueue(workInCh, workOutCh, inCh)
	queuer := Queuer{inCh, subCh}
	return &queuer
}

// Submit is the public interface for submitting jobs to the queuer.
func (q *Queuer) Submit(id string, text string, rules *model.Rules) {
	q.inCh <- &job{id, text, rules}
}

// Subscribe is the public interface that allows goroutines to
// indicate interest in the status of a job.
func (q *Queuer) Subscribe(id string) (chan Notification, bool) {
	ch := make(chan Notification)
	reply := make(chan bool)
	q.subCh <- sub{id, ch, reply}
	result := <-reply
	if !result {
		ch = nil
	}
	return ch, result
}

// Queuer is the interface to the job queuer.
type Queuer struct {
	inCh  chan *job
	subCh chan sub
}

// job is what gets stored in the job queue to represent an ongoing
// blinky generation.
type job struct {
	id    string
	text  string
	rules *model.Rules
}

// sub is the information passed in a new subscription.
type sub struct {
	id    string
	ch    chan Notification
	reply chan bool
}

// Start worker goroutine to process jobs one at a time.
func startWorker(bm *processing.BlinkyMaker,
	tmpl *template.Template) (chan *job, chan Notification) {
	inCh := make(chan *job)
	outCh := make(chan Notification)

	go func() {
		for {
			job := <-inCh
			htmlURL, errorMsg, err := bm.Make(job.text, job.rules, tmpl)
			if htmlURL != "" {
				outCh <- URLNotification{htmlURL}
			} else if errorMsg != "" {
				outCh <- ErrorNotification{errorMsg}
			} else if err != nil {
				outCh <- ErrorNotification{err.Error()}
			}
		}
	}()

	return inCh, outCh
}

// Start queue manager goroutine.
func startQueue(workQueue chan *job, workNotify chan Notification,
	inCh chan *job) chan sub {
	jobs := map[string]*job{}               // Currently queued jobs.
	chs := map[string][]chan Notification{} // Interest channels.
	queue := []string{}                     // Job queue.
	current := ""                           // Currently processing job ID.
	subCh := make(chan sub)                 // Subscription channel.

	go func() {
		for {
			select {
			case job := <-inCh:
				// New job submission: add to management structures, add to
				// end of queue, reply with job ID.
				jobs[job.id] = job
				chs[job.id] = []chan Notification{}
				queue = append(queue, job.id)

			case n := <-workNotify:
				// Notification from worker goroutine: pass the notification
				// on to any interested listeners, clear the job from the job
				// management structures, mark that there's currently no job
				// being processed (current == ""), delete job from the front
				// of the queue, inform interested listeners for any queued
				// jobs of changes in queue position.
				for _, ch := range chs[current] {
					ch <- n
					close(ch)
				}
				delete(jobs, current)
				delete(chs, current)
				current = ""
				queue = queue[1:]
				for i, id := range queue {
					for _, ch := range chs[id] {
						ch <- QueueNotification{i}
					}
				}

			case sub := <-subCh:
				// New interested listener subscription: add to the interested
				// list and return whether the job ID actually exists.
				_, ok := chs[sub.id]
				if ok {
					chs[sub.id] = append(chs[sub.id], sub.ch)
				}
				sub.reply <- ok
				for i, id := range queue {
					if id == sub.id {
						sub.ch <- QueueNotification{i}
					}
				}
			}

			// No current job, so let interested listeners know that the job
			// at the front of the queue is now being processed and send the
			// job to the worker.
			if current == "" && len(queue) > 0 {
				current = queue[0]
				workQueue <- jobs[current]
			}
		}
	}()

	return subCh
}
