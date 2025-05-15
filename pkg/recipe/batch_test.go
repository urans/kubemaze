package recipe

import (
	"fmt"
	"log/slog"
	"math/rand"
	"testing"
)

func TestSlowStartBatch(t *testing.T) {
	type args struct {
		count            int
		initialBatchSize int
		fn               func() error
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// {
		// 	"1", args{
		// 		count:            1,
		// 		initialBatchSize: 5,
		// 		fn: func() error {
		// 			slog.Info("slow start batch exec fn")
		// 			return nil
		// 		},
		// 	}, 1, false,
		// },
		// {
		// 	"10", args{
		// 		count:            10,
		// 		initialBatchSize: 5,
		// 		fn: func() error {
		// 			slog.Info("slow start batch exec fn")
		// 			return nil
		// 		},
		// 	}, 10, false,
		// },
		// {
		// 	"20", args{
		// 		count:            20,
		// 		initialBatchSize: 5,
		// 		fn: func() error {
		// 			slog.Info("slow start batch exec fn")
		// 			return nil
		// 		},
		// 	}, 20, false,
		// },
		// {
		// 	"100", args{
		// 		count:            100,
		// 		initialBatchSize: 5,
		// 		fn: func() error {
		// 			slog.Debug("slow start batch exec fn")
		// 			return nil
		// 		},
		// 	}, 100, false,
		// },
		{
			"1000", args{
				count:            1000,
				initialBatchSize: 5,
				fn: func() error {
					slog.Debug("slow start batch exec fn")
					return nil
				},
			}, 1000, false,
		},
		{
			"1000-with-error", args{
				count:            1000,
				initialBatchSize: 5,
				fn: func() error {
					slog.Debug("slow start batch exec fn")

					i := rand.Intn(1000)
					if i%7 == 1 {
						return fmt.Errorf("index %d: error occurs", i)
					}
					return nil
				},
			}, 1000, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SlowStartBatch(tt.args.count, tt.args.initialBatchSize, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlowStartBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlowStartBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
