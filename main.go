package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	// Carregar config da AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("não foi possível carregar a configuração do SDK, %v", err)
	}

	// Criar client EC2
	svc := ec2.NewFromConfig(cfg)

	// Especificar os detalhes da instância
	runResult, err := svc.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0c55b159cbfafe1f0"), // LEMBRAR DE SUBSTITUIR A ID da AMI correta!!!!
		InstanceType: types.InstanceTypeT2Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})

	if err != nil {
		log.Fatalf("não foi possível criar a instância, %v", err)
	}

	instanceID := *runResult.Instances[0].InstanceId
	fmt.Printf("Instância criada com ID: %s\n", instanceID)

	// add tags à instância
	_, err = svc.CreateTags(context.TODO(), &ec2.CreateTagsInput{
		Resources: []string{instanceID},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("example-instance"),
			},
		},
	})
	if err != nil {
		log.Fatalf("não foi possível adicionar tags à instância, %v", err)
	}

	fmt.Println("Tag adicionada com sucesso à instância")

	// Esperando estado "running" da instância
	err = svc.WaitUntilInstanceRunning(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		log.Fatalf("erro ao esperar pela instância estar em estado 'running', %v", err)
	}

	// Buscar informações da instância
	descResult, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		log.Fatalf("não foi possível descrever a instância, %v", err)
	}

	if len(descResult.Reservations) > 0 && len(descResult.Reservations[0].Instances) > 0 {
		publicIP := *descResult.Reservations[0].Instances[0].PublicIpAddress
		fmt.Printf("IP público da instância: %s\n", publicIP)
	} else {
		log.Println("não foi possível encontrar o IP público da instância")
	}
}
