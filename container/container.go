package container

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"taskio"
)

// todo: output should be protected data
func Run(binary io.Reader, output io.Writer) {
	cmd := exec.Command("/proc/self/exe")
	cmd.Stdin = os.Stdin
	cmd.Stdout = output
	cmd.Stderr = os.Stderr

	r, w, err := os.Pipe() // todo: protect binary with encryption so that child process decrypts and runs it
	taskio.Must(err)
	cmd.ExtraFiles = []*os.File{r}

	go func() {
		_, err = io.Copy(w, binary)
		taskio.Must(err)
		taskio.Must(w.Close())
	}()

	cmd.SysProcAttr = &syscall.SysProcAttr{ // todo: isolate users, rootless container
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	taskio.Must(cmd.Run())
}

func RunChild() {
	pwd, err := os.Getwd()
	taskio.Must(err)

	path := filepath.Join(pwd, "taskio")
	err = os.MkdirAll(path, 0770)

	taskPath := filepath.Join(path, "exec.task")

	cmd := exec.Command(taskPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	taskio.Must(err)

	taskio.Must(syscall.Sethostname([]byte("taskio")))
	taskio.Must(os.Chdir("/"))
	taskio.Must(syscall.Mount("proc", "proc", "proc", 0, ""))

	taskio.Must(syscall.Mount(path, path, "tmpfs", 0, "size=512M"))
	taskio.Must(os.Chdir(path))

	taskioFile, err := os.OpenFile(taskPath, os.O_CREATE|os.O_RDWR, 0770)
	taskio.Must(err)

	binary := os.NewFile(uintptr(3), "pipe")
	_, err = io.Copy(taskioFile, binary)
	taskio.Must(err)

	_ = taskioFile.Close()

	taskio.Must(cmd.Run())

	taskio.Must(os.Chdir("/"))
	taskio.Must(syscall.Unmount("proc", 0))
	taskio.Must(syscall.Unmount(path, 0))
}
