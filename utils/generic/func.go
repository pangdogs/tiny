package generic

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/types"
)

type Func0[R any] func() R

func (f Func0[R]) Exec() (r R) {
	r, _ = f.Call(false, nil)
	return
}

func (f Func0[R]) Invoke() (r R, panicErr error) {
	return f.Call(true, nil)
}

func (f Func0[R]) Call(autoRecover bool, reportError chan error) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(), nil
}

func (f Func0[R]) CastDelegate() DelegateFunc0[R] {
	return []Func0[R]{f}
}

type Func1[A1, R any] func(A1) R

func (f Func1[A1, R]) Exec(a1 A1) (r R) {
	r, _ = f.Call(false, nil, a1)
	return
}

func (f Func1[A1, R]) Invoke(a1 A1) (r R, panicErr error) {
	return f.Call(true, nil, a1)
}

func (f Func1[A1, R]) Call(autoRecover bool, reportError chan error, a1 A1) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1), nil
}

func (f Func1[A1, R]) CastDelegate() DelegateFunc1[A1, R] {
	return []Func1[A1, R]{f}
}

type Func2[A1, A2, R any] func(A1, A2) R

func (f Func2[A1, A2, R]) Exec(a1 A1, a2 A2) (r R) {
	r, _ = f.Call(false, nil, a1, a2)
	return
}

func (f Func2[A1, A2, R]) Invoke(a1 A1, a2 A2) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2)
}

func (f Func2[A1, A2, R]) Call(autoRecover bool, reportError chan error, a1 A1, a2 A2) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2), nil
}

func (f Func2[A1, A2, R]) CastDelegate() DelegateFunc2[A1, A2, R] {
	return []Func2[A1, A2, R]{f}
}

type Func3[A1, A2, A3, R any] func(A1, A2, A3) R

func (f Func3[A1, A2, A3, R]) Exec(a1 A1, a2 A2, a3 A3) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3)
	return
}

func (f Func3[A1, A2, A3, R]) Invoke(a1 A1, a2 A2, a3 A3) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3)
}

func (f Func3[A1, A2, A3, R]) Call(autoRecover bool, reportError chan error, a1 A1, a2 A2, a3 A3) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3), nil
}

func (f Func3[A1, A2, A3, R]) CastDelegate() DelegateFunc3[A1, A2, A3, R] {
	return []Func3[A1, A2, A3, R]{f}
}

type Func4[A1, A2, A3, A4, R any] func(A1, A2, A3, A4) R

func (f Func4[A1, A2, A3, A4, R]) Exec(a1 A1, a2 A2, a3 A3, a4 A4) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4)
	return
}

func (f Func4[A1, A2, A3, A4, R]) Invoke(a1 A1, a2 A2, a3 A3, a4 A4) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4)
}

func (f Func4[A1, A2, A3, A4, R]) Call(autoRecover bool, reportError chan error, a1 A1, a2 A2, a3 A3, a4 A4) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4), nil
}

func (f Func4[A1, A2, A3, A4, R]) CastDelegate() DelegateFunc4[A1, A2, A3, A4, R] {
	return []Func4[A1, A2, A3, A4, R]{f}
}

type Func5[A1, A2, A3, A4, A5, R any] func(
	A1, A2, A3, A4, A5,
) R

func (f Func5[A1, A2, A3, A4, A5, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5)
	return
}

func (f Func5[A1, A2, A3, A4, A5, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5)
}

func (f Func5[A1, A2, A3, A4, A5, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5), nil
}

func (f Func5[A1, A2, A3, A4, A5, R]) CastDelegate() DelegateFunc5[A1, A2, A3, A4, A5, R] {
	return []Func5[A1, A2, A3, A4, A5, R]{f}
}

type Func6[A1, A2, A3, A4, A5, A6, R any] func(
	A1, A2, A3, A4, A5, A6,
) R

func (f Func6[A1, A2, A3, A4, A5, A6, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6)
	return
}

func (f Func6[A1, A2, A3, A4, A5, A6, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6)
}

func (f Func6[A1, A2, A3, A4, A5, A6, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6), nil
}

func (f Func6[A1, A2, A3, A4, A5, A6, R]) CastDelegate() DelegateFunc6[A1, A2, A3, A4, A5, A6, R] {
	return []Func6[A1, A2, A3, A4, A5, A6, R]{f}
}

type Func7[A1, A2, A3, A4, A5, A6, A7, R any] func(
	A1, A2, A3, A4, A5, A6, A7,
) R

func (f Func7[A1, A2, A3, A4, A5, A6, A7, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7)
	return
}

func (f Func7[A1, A2, A3, A4, A5, A6, A7, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7)
}

func (f Func7[A1, A2, A3, A4, A5, A6, A7, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7), nil
}

func (f Func7[A1, A2, A3, A4, A5, A6, A7, R]) CastDelegate() DelegateFunc7[A1, A2, A3, A4, A5, A6, A7, R] {
	return []Func7[A1, A2, A3, A4, A5, A6, A7, R]{f}
}

type Func8[A1, A2, A3, A4, A5, A6, A7, A8, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8,
) R

func (f Func8[A1, A2, A3, A4, A5, A6, A7, A8, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8)
	return
}

func (f Func8[A1, A2, A3, A4, A5, A6, A7, A8, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8)
}

func (f Func8[A1, A2, A3, A4, A5, A6, A7, A8, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8), nil
}

func (f Func8[A1, A2, A3, A4, A5, A6, A7, A8, R]) CastDelegate() DelegateFunc8[A1, A2, A3, A4, A5, A6, A7, A8, R] {
	return []Func8[A1, A2, A3, A4, A5, A6, A7, A8, R]{f}
}

type Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9,
) R

func (f Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	return
}

func (f Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9)
}

func (f Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9), nil
}

func (f Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R]) CastDelegate() DelegateFunc9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R] {
	return []Func9[A1, A2, A3, A4, A5, A6, A7, A8, A9, R]{f}
}

type Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10,
) R

func (f Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	return
}

func (f Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
}

func (f Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10), nil
}

func (f Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R]) CastDelegate() DelegateFunc10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R] {
	return []Func10[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, R]{f}
}

type Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11,
) R

func (f Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
	return
}

func (f Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
}

func (f Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11), nil
}

func (f Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R]) CastDelegate() DelegateFunc11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R] {
	return []Func11[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, R]{f}
}

type Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12,
) R

func (f Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
	return
}

func (f Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
}

func (f Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12), nil
}

func (f Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R]) CastDelegate() DelegateFunc12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R] {
	return []Func12[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, R]{f}
}

type Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13,
) R

func (f Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13)
	return
}

func (f Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13)
}

func (f Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13), nil
}

func (f Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R]) CastDelegate() DelegateFunc13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R] {
	return []Func13[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, R]{f}
}

type Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14,
) R

func (f Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14)
	return
}

func (f Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14)
}

func (f Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14), nil
}

func (f Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R]) CastDelegate() DelegateFunc14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R] {
	return []Func14[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, R]{f}
}

type Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15,
) R

func (f Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
	return
}

func (f Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
}

func (f Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15), nil
}

func (f Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R]) CastDelegate() DelegateFunc15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R] {
	return []Func15[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, R]{f}
}

type Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R any] func(
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16,
) R

func (f Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R]) Exec(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15, a16 A16,
) (r R) {
	r, _ = f.Call(false, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16)
	return
}

func (f Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R]) Invoke(
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15, a16 A16,
) (r R, panicErr error) {
	return f.Call(true, nil, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16)
}

func (f Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R]) Call(
	autoRecover bool, reportError chan error,
	a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9, a10 A10, a11 A11, a12 A12, a13 A13, a14 A14, a15 A15, a16 A16,
) (r R, panicErr error) {
	if f == nil {
		return types.ZeroT[R](), nil
	}

	if autoRecover {
		defer func() {
			if panicErr = types.Panic2Err(recover()); panicErr != nil {
				panicErr = fmt.Errorf("%w: %w", exception.ErrPanicked, panicErr)

				if reportError != nil {
					select {
					case reportError <- exception.PrintStackTrace(panicErr):
					default:
					}
				}
			}
		}()
	}

	return f(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16), nil
}

func (f Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R]) CastDelegate() DelegateFunc16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R] {
	return []Func16[A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12, A13, A14, A15, A16, R]{f}
}
