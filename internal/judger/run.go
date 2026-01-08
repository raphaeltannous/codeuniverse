package judger

// func (g *golangJudge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
// 	runWorkspace, err := os.MkdirTemp("", "run-go-*")
// 	if err != nil {
// 		return err
// 	}
// 	defer os.RemoveAll(runWorkspace)
// 	g.logger.Debug("creating run workspace", "runWorkspace", runWorkspace)

// 	l, ok := SupportedLanguages["golang"]
// 	if !ok {
// 		return errors.New("unsupported language")
// 	}
// 	problemTestDir := filepath.Join(ProblemsDataDir, "problems", problemSlug, l.internalSlug)
// 	srcDir := os.DirFS(problemTestDir)

// 	err = os.CopyFS(runWorkspace, srcDir)
// 	if err != nil {
// 		return err
// 	}

// 	utilsDir := os.DirFS(filepath.Join(ProblemsDataDir, "utils", l.internalSlug))
// 	err = os.CopyFS(filepath.Join(runWorkspace, "utils"), utilsDir)
// 	if err != nil {
// 		return err
// 	}

// 	solutionsDir := os.DirFS(filepath.Join(ProblemsDataDir, "solutions", l.internalSlug, problemSlug))
// 	err = os.CopyFS(filepath.Join(runWorkspace, "solutions"), solutionsDir)
// 	if err != nil {
// 		return err
// 	}

// 	srcTestFile, err := os.Open(filepath.Join(ProblemsDataDir, "problems", problemSlug, "example-tests.txt"))
// 	if err != nil {
// 		return err
// 	}

// 	dstTestFile, err := os.Create(filepath.Join(runWorkspace, "tests.txt"))
// 	if err != nil {
// 		return nil
// 	}

// 	_, err = io.Copy(dstTestFile, srcTestFile)
// 	if err != nil {
// 		return nil
// 	}

// 	srcTestFile.Close()
// 	dstTestFile.Close()

// 	userSubmittedCode := filepath.Join(runWorkspace, "main.go")
// 	if run.Code[:12] != "package main" {
// 		run.Code = "package main\n" + run.Code
// 	}
// 	if err := os.WriteFile(userSubmittedCode, []byte(run.Code), 0644); err != nil {
// 		return err
// 	}

// 	modFile := filepath.Join(runWorkspace, "go.mod")
// 	if err := os.WriteFile(modFile, g.getModContent(), 0644); err != nil {
// 		return err
// 	}
// 	g.logger.Debug("finished settinp up workspace")

// 	resp, err := g.cli.ContainerCreate(
// 		ctx,
// 		&container.Config{
// 			Image:           SupportedLanguages[run.Language].containerImage,
// 			Cmd:             []string{"go", "test", "-timeout", "7s", "."},
// 			NetworkDisabled: true,
// 			WorkingDir:      "/app",
// 		},
// 		&container.HostConfig{
// 			AutoRemove: false,
// 			Binds:      []string{runWorkspace + ":/app:rw"},
// 		},
// 		nil,
// 		nil,
// 		"",
// 	)
// 	defer g.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

// 	if err := g.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
// 		panic(err)
// 	}
// 	g.logger.Debug("started container")

// 	statusCh, errCh := g.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

// 	select {
// 	case status := <-statusCh:
// 		g.logger.Warn("status", "status", status)
// 		if status.StatusCode == 0 {
// 			run.IsAccepted = true
// 			run.Status = "ACCEPTED"
// 		} else {
// 			run.IsAccepted = false
// 			run.Status = "FAILED"
// 		}
// 	case <-errCh:
// 		run.IsAccepted = false
// 		run.Status = "ERROR"
// 	case <-ctx.Done():
// 		run.IsAccepted = false
// 		run.Status = "TIME LIMIT EXCEEDED"
// 	}
// 	g.logger.Debug("container is finished")

// 	out, err := g.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
// 	if err != nil {
// 		panic(err)
// 	}

// 	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

// 	return nil
// }
