package storage

type Option func(s *SalesStorage)

func WithIndexGranularity(size int64) Option {
	return func(s *SalesStorage) {
		s.indexGranularity = size
	}
}
