package prisma

import (
	"sketch-blend-document-module/db"
	"sync"
)

var (
	once   sync.Once
	client *db.PrismaClient
)

func GetPrismaClient() *db.PrismaClient {
	once.Do(func() {
		client = db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			panic(err)
		}
	})

	return client
}

func Disconnect() {
	if client == nil {
		return
	}

	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}
