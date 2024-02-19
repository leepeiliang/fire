package globals

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func task() {
	fmt.Println("I am a running job.")
}

func TestScheduleNextRunFromNow(t *testing.T) {
	now := time.Now()

	sched := gocron.NewScheduler()
	sched.ChangeLoc(time.UTC)

	job := sched.Every(10).Hour().From(NextTick())
	job.Do(task)

	next := job.NextScheduledTime()
	nextRounded := time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, time.UTC)

	expected := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC).Add(5 * time.Second)

	assert.Exactly(t, expected, nextRounded)

	fmt.Println(expected)
	fmt.Println(nextRounded)
}
