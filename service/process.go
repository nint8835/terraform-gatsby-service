package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type _ProcessBody struct {
	Code string `json:"code" binding:"required"`
}

func _CleanupContainer(containerID string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	err = cli.ContainerRemove(context.TODO(), containerID, types.ContainerRemoveOptions{})

	if err != nil {
		panic(err)
	}
}

func _ProcessPost(c *gin.Context) {
	var body _ProcessBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(body)

	cli, err := client.NewEnvClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stopTimeout := 3

	container, err := cli.ContainerCreate(
		context.TODO(),
		&container.Config{
			Image:       "hashicorp/terraform:light",
			StopTimeout: &stopTimeout,
			Entrypoint:  strslice.StrSlice{"terraform", "version", "-no-color"},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		"terraform-gatsby-service-test",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer _CleanupContainer(container.ID)

	err = cli.ContainerStart(
		context.TODO(),
		container.ID,
		types.ContainerStartOptions{},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := cli.ContainerWait(context.TODO(), container.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logReader, err := cli.ContainerLogs(
		context.TODO(),
		container.ID,
		types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer logReader.Close()

	logBuffer := new(strings.Builder)
	_, err = io.Copy(logBuffer, logReader)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logs := logBuffer.String()
	fmt.Println(logs)

	c.JSON(http.StatusOK, gin.H{
		"container": container,
		"status":    status,
		"logs":      logs,
	})

}
