package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type jobReader struct {
	path  string
	cache map[string]Job
	sync.Mutex
}

func newJobsReader(p string) *jobReader {
	os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0666)

	return &jobReader{
		path:  p,
		cache: map[string]Job{},
	}
}

func (r *jobReader) Read() map[string]Job {
	r.Lock()
	defer r.Unlock()
	if len(r.cache) > 0 {
		return r.cache
	}

	file, err := ioutil.ReadFile(r.path)
	if err != nil {
		panic(fmt.Sprintf("Error opening file '%s'", r.path))
	}

	var jobs map[string]Job
	json.Unmarshal(file, &jobs)
	if len(jobs) == 0 {
		jobs = map[string]Job{}
	}

	r.cache = jobs
	return jobs
}

func (r *jobReader) Write(jobs map[string]Job) {
	file, err := json.MarshalIndent(jobs, "", " ")
	if err != nil {
		panic(fmt.Sprintf("Error Marshaling jobs: %v", jobs))
	}

	err = ioutil.WriteFile(r.path, file, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing to file '%s'", r.path))
	}

	r.Lock()
	r.cache = jobs
	r.Unlock()
}

// JobState - allows persisting job data
var JobState = newJobsReader("./data/Jobs.json")