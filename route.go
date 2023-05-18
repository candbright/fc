package frouter

type routes[Data any] struct {
	asyncRoutes []route[Data]
	syncRoutes  []route[Data]
}

type route[Data any] struct {
	f        func(instance Data) error
	rollback func(instance Data)
}

func (r *routes[Data]) Add(f func(instance Data) error, rollback func(instance Data)) {
	if r.syncRoutes == nil {
		r.syncRoutes = make([]route[Data], 0)
	}
	r.syncRoutes = append(r.syncRoutes, route[Data]{f: f, rollback: rollback})
}

func (r *routes[Data]) AddAsync(f func(instance Data) error, rollback func(instance Data)) {
	if r.asyncRoutes == nil {
		r.asyncRoutes = make([]route[Data], 0)
	}
	r.asyncRoutes = append(r.asyncRoutes, route[Data]{f: f, rollback: rollback})
}

func (r *routes[Data]) ExecuteAll(instance Data) error {
	var err error
	var rollbacks []func(instance Data)
	defer func() {
		if rollbacks != nil {
			for i := len(rollbacks) - 1; i >= 0; i-- {
				if rollbacks[i] != nil {
					rollbacks[i](instance)
				}
			}
		}
	}()
	for _, f := range r.syncRoutes {
		if f.f == nil {
			continue
		}
		err = f.f(instance)
		if err != nil {
			rollbacks = append(rollbacks, f.rollback)
			return err
		}
		if rollbacks == nil {
			rollbacks = make([]func(instance Data), 0)
		}
		rollbacks = append(rollbacks, f.rollback)
	}
	return nil
}

func (r *routes[Data]) ExecuteAsyncAll(instance Data) error {
	var rollbacks []func(instance Data)
	defer func() {
		if rollbacks != nil {
			for _, rb := range rollbacks {
				if rb != nil {
					rb(instance)
				}
			}
		}
	}()
	for _, f := range r.asyncRoutes {
		if f.f == nil {
			continue
		}
		err := f.f(instance)
		if err != nil {
			rollbacks = append(rollbacks, f.rollback)
		}
	}
	return nil
}
