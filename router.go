package frouter

import "time"

const (
	Init   = "init"
	Add    = "add"
	Update = "update"
	Delete = "delete"
)

type Router[Data any] struct {
	asyncSignal      bool
	asyncStore       []asyncData[Data]
	asyncChan        chan asyncData[Data]
	preSyncRoutes    *routes[Data]
	afterSyncRoutes  *routes[Data]
	preAsyncRoutes   *routes[Data]
	afterAsyncRoutes *routes[Data]
	asyncRoutesMap   map[string]*routes[Data]
	syncRoutesMap    map[string]*routes[Data]
}

type asyncData[Data any] struct {
	asyncRoutes *routes[Data]
	instance    Data
}

func New[Data any]() *Router[Data] {
	return &Router[Data]{}
}

func (s *Router[Data]) start() {
	go func() {
		for {
			as := <-s.asyncChan
			if !s.asyncSignal {
				s.store(as)
				continue
			}
			s.load()
			err := s.preAsync(as.instance)
			if err != nil {
				continue
			}
			err = as.asyncRoutes.ExecuteAsyncAll(as.instance)
			if err != nil {
				continue
			}
			err = s.afterAsync(as.instance)
			if err != nil {
				continue
			}
		}
	}()
}

func (s *Router[Data]) Start() {
	if s.asyncSignal == false {
		s.asyncSignal = true
	}
}

func (s *Router[Data]) Stop() {
	if s.asyncSignal == true {
		s.asyncSignal = false
	}
}

func (s *Router[Data]) store(as asyncData[Data]) {
	if s.asyncStore == nil {
		s.asyncStore = make([]asyncData[Data], 0)
	}
	s.asyncStore = append(s.asyncStore, as)
}

func (s *Router[Data]) load() {
	if s.asyncStore != nil {
		for _, as := range s.asyncStore {
			err := s.preAsync(as.instance)
			if err != nil {
				continue
			}
			err = as.asyncRoutes.ExecuteAsyncAll(as.instance)
			if err != nil {
				continue
			}
			err = s.afterAsync(as.instance)
			if err != nil {
				continue
			}
		}
	}
}

func (s *Router[Data]) SetDelay(delay time.Duration) {
	s.RegisterPreSync(func(instance Data) error {
		time.Sleep(delay)
		return nil
	})
	s.RegisterPreAsync(func(instance Data) error {
		time.Sleep(delay)
		return nil
	})
}

func (s *Router[Data]) RegisterAsync(key string, asyncFunc func(instance Data) error, rollback func(instance Data)) {
	if s.asyncRoutesMap == nil {
		s.asyncRoutesMap = make(map[string]*routes[Data])
		s.asyncChan = make(chan asyncData[Data])
		s.asyncSignal = true
		s.start()
	}
	if s.asyncRoutesMap[key] == nil {
		s.asyncRoutesMap[key] = &routes[Data]{}
	}
	s.asyncRoutesMap[key].AddAsync(asyncFunc, rollback)
}

func (s *Router[Data]) Async(key string, instance Data) {
	if s.asyncRoutesMap == nil || s.asyncRoutesMap[key] == nil {
		return
	}
	s.asyncChan <- asyncData[Data]{s.asyncRoutesMap[key], instance}
}

func (s *Router[Data]) RegisterSync(key string, syncFunc func(instance Data) error, rollback func(instance Data)) {
	if s.syncRoutesMap == nil {
		s.syncRoutesMap = make(map[string]*routes[Data])
	}
	if s.syncRoutesMap[key] == nil {
		s.syncRoutesMap[key] = &routes[Data]{}
	}
	s.syncRoutesMap[key].Add(syncFunc, rollback)
}

func (s *Router[Data]) Sync(key string, instance Data) error {
	if s.syncRoutesMap == nil || s.syncRoutesMap[key] == nil {
		return nil
	}
	err := s.preSync(instance)
	if err != nil {
		return err
	}
	err = s.syncRoutesMap[key].ExecuteAll(instance)
	if err != nil {
		return err
	}
	err = s.afterSync(instance)
	if err != nil {
		return err
	}
	return nil
}

func (s *Router[Data]) RegisterPreSync(preFunc func(instance Data) error) {
	if s.preSyncRoutes == nil {
		s.preSyncRoutes = &routes[Data]{}
	}
	s.preSyncRoutes.Add(preFunc, nil)
}

func (s *Router[Data]) preSync(instance Data) error {
	if s.preSyncRoutes == nil {
		return nil
	}
	err := s.preSyncRoutes.ExecuteAll(instance)
	if err != nil {
		return err
	}
	return nil
}

func (s *Router[Data]) RegisterAfterSync(afterFunc func(instance Data) error) {
	if s.afterSyncRoutes == nil {
		s.afterSyncRoutes = &routes[Data]{}
	}
	s.afterSyncRoutes.Add(afterFunc, nil)
}

func (s *Router[Data]) afterSync(instance Data) error {
	if s.afterSyncRoutes == nil {
		return nil
	}
	err := s.afterSyncRoutes.ExecuteAll(instance)
	if err != nil {
		return err
	}
	return nil
}

func (s *Router[Data]) RegisterPreAsync(preFunc func(instance Data) error) {
	if s.preAsyncRoutes == nil {
		s.preAsyncRoutes = &routes[Data]{}
	}
	s.preAsyncRoutes.Add(preFunc, nil)
}

func (s *Router[Data]) preAsync(instance Data) error {
	if s.preAsyncRoutes == nil {
		return nil
	}
	err := s.preAsyncRoutes.ExecuteAll(instance)
	if err != nil {
		return err
	}
	return nil
}

func (s *Router[Data]) RegisterAfterAsync(afterFunc func(instance Data) error) {
	if s.afterAsyncRoutes == nil {
		s.afterAsyncRoutes = &routes[Data]{}
	}
	s.afterAsyncRoutes.Add(afterFunc, nil)
}

func (s *Router[Data]) afterAsync(instance Data) error {
	if s.afterAsyncRoutes == nil {
		return nil
	}
	err := s.afterAsyncRoutes.ExecuteAll(instance)
	if err != nil {
		return err
	}
	return nil
}
