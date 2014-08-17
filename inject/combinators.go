package inject

func Constant(x Any) func() Any {
	return func() Any {
		return x
	}
}
