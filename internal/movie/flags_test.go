package movie

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestFromFlags(t *testing.T) {
	testMovie := func(t *testing.T, path string) {
		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		movie, err := FromFlags(flags, path)
		if !assert.NoError(t, err) {
			return
		}

		bodyPad, err := flags.GetIntSlice(PadFlag)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, movie.BodyStyle.GetPaddingTop(), bodyPad[0])
		assert.Equal(t, movie.BodyStyle.GetPaddingRight(), bodyPad[1])
		assert.Equal(t, movie.BodyStyle.GetPaddingBottom(), bodyPad[2])
		assert.Equal(t, movie.BodyStyle.GetPaddingLeft(), bodyPad[1])

		progressPad, err := flags.GetIntSlice(ProgressPadFlag)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, movie.ProgressStyle.GetPaddingTop(), progressPad[0])
		assert.Equal(t, movie.ProgressStyle.GetPaddingRight(), progressPad[1])
		assert.Equal(t, movie.ProgressStyle.GetPaddingBottom(), progressPad[2])
		assert.Equal(t, movie.ProgressStyle.GetPaddingLeft(), progressPad[1])
	}

	t.Run("default embedded", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "")
	})

	t.Run("short_intro embedded", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "short_intro")
	})

	t.Run("short_intro file", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "../../movies/short_intro.txt")
	})

	t.Run("invalid speed", func(t *testing.T) {
		t.Parallel()

		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		if err := flags.Set(SpeedFlag, "-1"); !assert.NoError(t, err) {
			return
		}

		if _, err := FromFlags(flags, ""); !assert.Error(t, err) {
			return
		}
	})
}
