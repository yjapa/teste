package handler

import (
	"fmt"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
	sdkStorage "github.com/klever-io/getchain-sdk/storage"
	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/cliFlags"
	"github.com/klever.io/getchain-raw-worker/consumer"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/models"
	"github.com/klever.io/getchain-raw-worker/utils"
	"github.com/klever.io/getchain-raw-worker/worker"
)

func NewHandler(chain common.CoinType, kafkaTopic string) (interfaces.Handler, error) {
	worker, err := worker.NewWorker(chain, kafkaTopic)
	if err != nil {
		return nil, err
	}

	consumerWorker, err := consumer.NewConsumerWorker(chain, kafkaTopic)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     utils.RequireEnv("REDIS_CLIENT_URL"),
		Password: "",
		DB:       0,
	})

	redis := sdkStorage.NewRedis(client, 5*time.Second, 60*time.Second)

	return &Handler{
		Worker:         worker,
		ConsumerWorker: consumerWorker,
		Redis:          redis,
		Coin:           chain,
		KafkaTopic:     kafkaTopic,
	}, nil
}

func (s *Handler) Forward() bool {
	explorer := s.Worker.GetExplorer()
	ticker := time.NewTicker(time.Second * 5)
	done := make(chan bool)
	redisKey := fmt.Sprintf("%s_LAST_BLOCK_SENT", explorer.Type().String())
	for {
		select {
		case <-ticker.C:
			to, err := explorer.GetLastBlock()
			if err != nil {
				return false
			}

			var from *big.Int
			err = s.Redis.Get(redisKey, &from)
			if err == sdkStorage.ErrFailedToGet {
				err = s.Redis.Put(redisKey, to)
				if err != nil {
					return false
				}
				continue
			}

			if from.Cmp(to) == -1 {
				blockBatch, err := explorer.GetBlockBatch(from.Add(from, big.NewInt(1)), to)
				if err != nil {
					return false
				}
				if err := s.Worker.BlockBatchToKafka(blockBatch); err != nil {
					return false
				}

				err = s.Redis.Put(redisKey, to)
				if err != nil {
					return false
				}
			}

		case <-done:
			return false
		}
	}
}

func (s *Handler) Backward() bool {
	return false
}

func (s *Handler) Status() interface{} {
	return utils.ErrNotImplemented
}

func (s *Handler) GetHighestBlock() *big.Int {
	return big.NewInt(0)
}

func (s *Handler) GetLowestBlock() *big.Int {
	return big.NewInt(0)
}

func (s *Handler) GetWorker() interfaces.Worker {
	return s.Worker
}

func (s *Handler) GetConsumerWorker() interfaces.ConsumerWorker {
	return s.ConsumerWorker
}

func (s *Handler) Start(options models.Options) error {

	switch options.Mode {
	case cliFlags.RawWorker:
		if options.Forward {
			ok := s.Forward()
			if !ok {
				return utils.ErrNotImplemented
			}

			return nil
		}
		if options.Backward {
			ok := s.Backward()
			if !ok {
				return utils.ErrNotImplemented
			}

			return nil
		}
	case cliFlags.Listener:
		return s.ConsumerWorker.ConsumeFromKafka(nil)
	case cliFlags.Parser:
		parser, err := GetParser(s.Coin, s.KafkaTopic)
		if err != nil {
			return err
		}

		return s.ConsumerWorker.ConsumeFromKafka(parser.ParseBlock)
	default:
		return utils.ErrNotImplemented
	}

	return utils.ErrNotImplemented
}

func GetParser(chain common.CoinType, kafkaTopic string) (interfaces.Parser, error) {
	parser, ok := parsers[chain]
	if !ok {
		return nil, utils.ErrCoinNotImplemented
	}

	return parser(kafkaTopic), nil
}
