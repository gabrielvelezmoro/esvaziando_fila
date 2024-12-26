package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	amqpURL   = "amqp://guest:guest@localhost:5672/" // URL de conexão com RabbitMQ
	queueName = "minha_fila"                      // Nome da fila a ser esvaziada
)

func main() {
	// Conectar ao RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Abrir um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir um canal: %s", err)
	}
	defer ch.Close()

	// Verificar a fila
	queue, err := ch.QueueDeclarePassive(
		queueName, // Nome da fila
		true,      // Durável
		false,     // Auto-delete
		false,     // Exclusivo
		false,     // No-wait
		nil,       // Argumentos
	)
	if err != nil {
		log.Fatalf("Erro ao acessar a fila '%s': %s", queueName, err)
	}

	if queue.Messages == 0 {
		fmt.Printf("A fila '%s' já está vazia.\n", queueName)
		return
	}

	fmt.Printf("Esvaziando a fila '%s' com %d mensagens...\n", queueName, queue.Messages)

	// Consumir mensagens
	msgs, err := ch.Consume(
		queueName, // Nome da fila
		"",        // Consumer
		true,      // Auto-ack
		false,     // Exclusivo
		false,     // Local
		false,     // No-wait
		nil,       // Argumentos
	)
	if err != nil {
		log.Fatalf("Erro ao consumir mensagens: %s", err)
	}

	// Controlar o consumo e finalizar quando a fila for esvaziada
	done := make(chan bool)
	go func() {
		counter := 0
		timer := time.NewTimer(1 * time.Second) // Timer para detectar inatividade

		for {
			select {
			case <-msgs:
				counter++
				timer.Reset(1 * time.Second) // Reinicia o timer a cada mensagem consumida
			case <-timer.C:
				// Nenhuma mensagem foi recebida por 1 segundo
				done <- true
				return
			}
		}
	}()

	<-done
	fmt.Printf("Fila '%s' esvaziada com sucesso!\n", queueName)
}
