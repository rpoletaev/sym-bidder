package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	MULTI   = "MULTI"
	INCR    = "INCR"
	EXPIRE  = "EXPIRE"
	EXEC    = "EXEC"
	HINCRBY = "HINCRBY"
)

// ifaItem структура для хранения информации о последнем обновлении ключа
// и текущем номере серии
type ifaItem struct {
	lastUpdate      time.Time
	curentSeriesNum int64
}

type service struct {
	*Config
	mu          sync.Mutex
	lastUpdates map[string]*ifaItem
	pool        *redis.Pool
}

func CreateService(cfg *Config) *service {

	address := cfg.RedisHost + ":" + cfg.RedisPort

	return &service{
		Config:      cfg,
		lastUpdates: map[string]*ifaItem{},
		pool: redis.NewPool(func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		}, cfg.MaxIddleCons),
	}
}

// обрабатывает входные данные
func (s *service) Process(key, stat string) (pos int64, err error) {
	go s.saveStat(stat)

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// если ключа еще не было, то установим и вернем 0
	item, ok := s.lastUpdates[key]
	if !ok {
		s.lastUpdates[key] = &ifaItem{
			lastUpdate:      now,
			curentSeriesNum: 0,
		}

		return 0, nil
	}

	afterLastUpdate := time.Since(item.lastUpdate)

	// если ключ пришел в течение секунды со времени предыдущего обновления
	// то прочто вернем сохраненный номер серии
	if afterLastUpdate <= 1*time.Second {
		return item.curentSeriesNum, nil
	}

	con := s.pool.Get()
	defer con.Close()

	ifaKey := "ifa" + key

	// если время с последнего обновления больше секунды, но меньше интервала установки новой сессии
	// обновим ttl в redis для данного ключа, обновим время последнего обновления и вернем текущий номер серии
	if afterLastUpdate < s.IFASeriesInterval*time.Second {
		_, err = con.Do(EXPIRE, ifaKey, s.IFACounterTTL)
		item.lastUpdate = now
		return item.curentSeriesNum, err
	}

	// иначе обновим ttl, обновим счетчик, обновим время последнего апдейта
	con.Send(MULTI)
	con.Send(EXPIRE, ifaKey, s.IFACounterTTL)
	con.Send(INCR, ifaKey)
	res, err := redis.Values(con.Do(EXEC))

	if err != nil {
		return 0, err
	}

	item.curentSeriesNum = res[1].(int64)
	item.lastUpdate = now

	return item.curentSeriesNum, err
}

// saveStat сохранить в редис статистику
func (s *service) saveStat(statKey string) {
	con := s.pool.Get()
	defer con.Close()

	con.Do(HINCRBY, "stat", statKey, 1)
}

// Stats получает статистику из redis и возвращает, при успешном получении,
// иначе ошибка
func (s *service) Stats() ([]Stat, error) {
	con := s.pool.Get()
	defer con.Close()

	result, err := redis.Values(con.Do("HGETALL", "stat"))
	if err != nil {
		return nil, err
	}

	resLength := len(result)
	stats := make([]Stat, resLength/2, resLength/2)
	statsIndex := 0

	for i, val := range result {
		fmt.Printf("%d %v\n", i, val)
	}
	for i := 1; i < resLength; i += 2 {

		rawKey := result[i-1].([]byte)
		statItem := stats[statsIndex]

		if err := json.Unmarshal(rawKey, &statItem); err != nil {
			return nil, err
		}

		fmt.Printf("%d %v\n", i, result[i])
		// statItem.Count, _ = result[i]
	}

	return stats, nil
}
