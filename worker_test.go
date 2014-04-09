package goqless

import (
  "testing"
  "time"
)

var (
	host = "192.168.10.10"
	port = "6379"
)

func TestWorker(t *testing.T) {
	c, err := Dial(host, port)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	
	testNilJobWorker(t, c)
}

func testNilJobWorker(t *testing.T, c *Client) {
	q := c.Queue("goqless_testing_nilJob_queue")
	
	// Add job to queue
	data := struct {
		Str string
	}{
		"a string",
	}
	jid, err := q.Put("", "nilJob", data, -1, -1, []string{}, -1, []string{})
	if err != nil { t.Error(err.Error()) }
	job, err := c.GetJob(jid)
	if err != nil { t.Error(err.Error()) }
	if job.State != "waiting" { t.Error("Expected at least 1 job added to queue") }
	
	// Start Worker
	w, err := NewWorker(host+":"+port, "goqless_testing_nilJob_queue", 1)
	if err != nil { t.Error(err.Error()) }
	w.AddFunc("nilJob", func(j *Job) error { return nil })
	go w.Start()
	
	time.Sleep(1 * time.Second)
	
	// Ensure job was completed
	job, err = c.GetJob(jid)
	if err != nil { t.Error(err.Error()) }
	if job.State != "complete" { t.Error("Expected job to finish") }
}


