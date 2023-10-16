package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func updateDockerImage(imageName, imageTag string, updateInterval time.Duration) error {
	// 设置 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// 创建一个 context，用于执行 Docker 操作
	ctx := context.Background()

	// 设置定时器，定期执行操作
		ticker := time.NewTicker(updateInterval)
		fmt.Printf("Docker image %s:%s will be updated every %s\n", imageName, imageTag, updateInterval)
		for range ticker.C {
			log.Printf("Updating Docker image %s:%s...\n", imageName, imageTag)
			fmt.Println("Updating Docker image...")
			// 拉取最新的 Docker 镜像，强制重新拉取
			_, err = cli.ImagePull(ctx, imageName+":"+imageTag, types.ImagePullOptions{All: true})
			if err != nil {
				return err
			}

			log.Printf("Docker image %s:%s updated successfully!\n", imageName, imageTag)
			fmt.Println("Docker image updated successfully!")
		}
		return nil
	}
