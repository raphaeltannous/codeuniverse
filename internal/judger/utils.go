package judger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func writeToWorkspace(workspace string, filename string, content string, logger *slog.Logger) (*models.RunResult, error) {
	filePath := filepath.Join(workspace, filename)
	if err := os.WriteFile(filePath, []byte(content), 0640); err != nil {
		logger.Error("failed to write file", "filename", filename, "content", content, "err", err)
		return &models.RunResult{Status: models.StatusInternalServerError}, err
	}

	return nil, nil
}

func copyRuntimeFilesToWorkspace(workspace string, files []string) (*models.RunResult, error) {
	for _, file := range files {
		src, err := os.Open(file)
		if err != nil {
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		dstPath := filepath.Join(workspace, filepath.Base(file))
		dst, err := os.Create(dstPath)
		if err != nil {
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		_, err = io.Copy(dst, src)
		if err != nil {
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}
	}

	return nil, nil
}

func writeTestcasesToWorkspace(workspace string, testcases []*models.ProblemTestcase) (*models.RunResult, error) {
	testcasesPath := filepath.Join(workspace, "testcases.json")
	testcasesFile, err := os.Create(testcasesPath)
	if err != nil {
		return &models.RunResult{Status: models.StatusInternalServerError}, err
	}
	defer testcasesFile.Close()

	encoder := json.NewEncoder(testcasesFile)
	if err := encoder.Encode(testcases); err != nil {
		return &models.RunResult{Status: models.StatusInternalServerError}, err
	}

	return nil, nil
}

func watchContainerMemory(statsCtx context.Context, cli *client.Client, containerId string, peakMemory *uint64) {
	stats, err := cli.ContainerStats(statsCtx, containerId, true)
	if err != nil {
		return
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)
	for {
		var v container.StatsResponse
		if err := decoder.Decode(&v); err != nil {
			return
		}
		if v.MemoryStats.Usage > *peakMemory {
			*peakMemory = v.MemoryStats.Usage
		}
	}
}
