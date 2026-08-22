package syslogrx

import "context"

func (r *Receiver) ServeContext(ctx context.Context, packets <-chan []byte) error {
	_ = ctx
	for p := range packets {
		if _, err := r.Handle(context.Background(), p); err != nil {
			return err
		}
	}
	return nil
}
