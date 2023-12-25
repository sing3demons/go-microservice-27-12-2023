package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/sing3demons/product/logger"
)

func NewConsumerGroup(servers, topics []string, groupID string, logger *logger.Logger) (ConsumerHandler, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	version, _ := sarama.ParseKafkaVersion("1.0.0")
	config.Version = version

	consumer, err := sarama.NewConsumerGroup(servers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumerHandler{consumer, topics, logger}, nil
}

type ConsumerHandler interface {
	StartConsumer(handler sarama.ConsumerGroupHandler)
}

type consumerHandler struct {
	consumer sarama.ConsumerGroup
	topics   []string
	logger   *logger.Logger
}

func (h *consumerHandler) StartConsumer(handler sarama.ConsumerGroupHandler) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	h.logger.Info("Starting consumer...")
	go func() {
		defer wg.Done()
		for {
			if err := h.consumer.Consume(ctx, h.topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Printf("error from consumer: %v\n", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// Handle graceful shutdown
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Println("Received termination signal. Initiating shutdown...")
		cancel()
	case <-ctx.Done():
		fmt.Println("terminating: context cancelled")

	}
	// Wait for the consumer to finish processing
	wg.Wait()
}

type Event struct {
	Header map[string]any `json:"header"`
	Body   any            `json:"body"`
}
