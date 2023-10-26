package main

import (
	"context"
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
	//fmt.Printf("Docker image %s:%s will be updated every %s\n", imageName, imageTag, updateInterval)
	for range ticker.C {
		log.Printf("Updating Docker image %s:%s...\n", imageName, imageTag)
		//fmt.Println("Updating Docker image...")
		// 拉取最新的 Docker 镜像，强制重新拉取
		_, err = cli.ImagePull(ctx, imageName+":"+imageTag, types.ImagePullOptions{All: true})
		if err != nil {
			return err
		}

		log.Printf("Docker image %s:%s updated successfully!\n", imageName, imageTag)
		//fmt.Println("Docker image updated successfully!")

		// 删除原有 Docker 镜像
		err = deleteDockerImage(imageName, imageTag)
		if err != nil {
			log.Printf("Error deleting old Docker image: %v\n", err)
		}
	}

	return nil
}

func deleteDockerImage(imageName, imageTag string) error {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// 创建一个 context，用于执行 Docker 操作
	ctx := context.Background()

	// 设置删除选项
	removeOptions := types.ImageRemoveOptions{
		Force:         true, // 强制删除
		PruneChildren: true, // 删除关联的子镜像
	}

	// 删除指定的 Docker 镜像
	_, err = cli.ImageRemove(ctx, imageName+":"+imageTag, removeOptions)
	if err != nil {
		return err
	}

	return nil
}
