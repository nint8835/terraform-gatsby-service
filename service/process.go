package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	cli, err := client.NewEnvClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stopTimeout := 300

	container, err := cli.ContainerCreate(
		context.TODO(),
		&container.Config{
			Image:       "worker",
			StopTimeout: &stopTimeout,
			Env:         []string{fmt.Sprintf("TERRAFORM_SOURCE=%s", body.Code)},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		fmt.Sprintf("terraform-gatsby-service-%s", uuid.New().String()),
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

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	_, err = stdcopy.StdCopy(stdout, stderr, logReader)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"container": container,
		"status":    status,
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
	})

}
