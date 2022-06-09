package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/egnd/go-pipeline"
	"github.com/stretchr/testify/assert"
)

func Test_DecorateDoer(t *testing.T) {
	cases := []string{
		"hello world",
		"",
	}

	for k, phrase := range cases {
		t.Run(fmt.Sprint(k+1), func(t *testing.T) {
			var res string
			decorators := []pipeline.DoerDecorator{}
			if len(phrase) > 0 {
				for i := 0; i < len(phrase)-1; i++ {
					i := i
					decorators = append(decorators, func(next pipeline.Tasker) pipeline.Tasker {
						return func(task pipeline.Task) error {
							res += phrase[i : i+1]
							return next(task)
						}
					})
				}
			}

			pipeline.DecorateDoer(func(_ pipeline.Task) error {
				if len(phrase) > 0 {
					res += phrase[len(phrase)-1:]
				}
				return nil
			}, decorators...)(nil)

			assert.EqualValues(t, phrase, res+"")
		})
	}
}
