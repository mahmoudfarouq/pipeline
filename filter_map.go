package pipeline

type Option[I InOut] struct {
	Some *I
}

func (o Option[I]) IsSome() bool {
	return o.Some != nil
}

func FilterMap[I Input, O Output](transformer func(I) Option[O]) Step[I, O] {
	//return func(in <-chan I) <-chan O {
	//	out := make(chan O, 4)
	//
	//	dirtyChannel := Workers(transformer, 1, in)
	//	go func() {
	//		defer close(out)
	//
	//		for o := range dirtyChannel {
	//			if o.IsSome() {
	//				out <- *o.Some
	//			}
	//		}
	//	}()
	//
	//	return out
	//}

	return Compose(
		Map(transformer),
		Compose(
			Filter(Option[O].IsSome),
			Map(func(o Option[O]) O {
				return *o.Some
			}),
		),
	)
}
