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
			// Pass the real ctx (not context.Background) so a canceled
			// upstream fails the call instead of silently completing it.
			if _, err := r.Handle(ctx, p); err != nil {
				return err
			}
		}
	}
}
