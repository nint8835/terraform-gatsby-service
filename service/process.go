package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	Code string `json:"code" binding:"required,max=5000"`
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

func _Error(c *gin.Context, err error, wrapperMessage string) {
	wrappedError := fmt.Errorf("%s: %w", wrapperMessage, err)
	var returnedError error
	if gin.Mode() == gin.DebugMode {
		returnedError = wrappedError
	} else {
		returnedError = errors.New("an internal error occurred processing your request")
	}
	fmt.Fprintln(os.Stderr, wrappedError.Error())
	c.JSON(http.StatusInternalServerError, gin.H{"error": returnedError.Error()})
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

	stopTimeout := 5

	container, err := cli.ContainerCreate(
		context.TODO(),
		&container.Config{
			Image:           "nint8835/terraform-gatsby-service-worker",
			StopTimeout:     &stopTimeout,
			Env:             []string{fmt.Sprintf("TERRAFORM_SOURCE=%s", body.Code)},
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory: 128 * 1024 * 1024,
			},
		},
		&network.NetworkingConfig{},
		fmt.Sprintf("terraform-gatsby-service-%s", uuid.New().String()),
	)
	if err != nil {
		_Error(c, err, "an error occurred while creating a worker container")
		return
	}

	defer _CleanupContainer(container.ID)

	err = cli.ContainerStart(
		context.TODO(),
		container.ID,
		types.ContainerStartOptions{},
	)

	if err != nil {
		_Error(c, err, "an error occurred while starting a worker container")
		return
	}

	_, err = cli.ContainerWait(context.TODO(), container.ID)

	if err != nil {
		_Error(c, err, "an error occurred while waiting for a worker container to terminate")
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
		_Error(c, err, "an error occurred while reading logs from a worker container")
		return
	}

	defer logReader.Close()

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	_, err = stdcopy.StdCopy(stdout, stderr, logReader)

	if err != nil {
		_Error(c, err, "an error occurred while processing logs from a worker container")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	})

}
