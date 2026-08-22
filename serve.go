package syslogrx

import "context"

func (r *Receiver) ServeContext(ctx context.Context, packets <-chan []byte) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p, ok := <-packets:
			if !ok {
				return nil
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			if _, err := r.Handle(ctx, p); err != nil {
				return err
			}
		}
	}
}
