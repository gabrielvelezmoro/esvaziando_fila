package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// URL de conexão com o RabbitMQ
	amqpURL := "amqp://guest:guest@localhost:5672/"
	queueName := "minha_fila"

	// Conectando ao servidor RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Abrindo um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir um canal: %s", err)
	}
	defer ch.Close()

	// Verificando se a fila existe
	queue, err := ch.QueueDeclarePassive(
		queueName, // nome da fila
		true,      // durável
		false,     // auto-delete
		false,     // exclusivo
		false,     // sem-wait
		nil,       // argumentos
	)
	if err != nil {
		log.Fatalf("Falha ao declarar fila: %s", err)
	}

	fmt.Printf("Esvaziando a fila '%s' com %d mensagens...\n", queueName, queue.Messages)

	// Consumindo as mensagens e descartando-as
	msgs, err := ch.Consume(
		queueName, // nome da fila
		"",        // consumer
		true,      // auto-ack
		false,     // exclusivo
		false,     // sem-local
		false,     // no-wait
		nil,       // argumentos
	)
	if err != nil {
		log.Fatalf("Falha ao consumir mensagens: %s", err)
	}

	// Usar um canal para controle de fluxo
	done := make(chan bool)

	go func() {
		for range msgs {
				// Apenas consumir e descartar as mensagens
		}
		done <- true
	}()

	// Aguardar o término do consumo
	<-done

	fmt.Println("Fila esvaziada com sucesso!")
}
