package goutil

import (
	"errors"
	"path/filepath"
	"runtime"
)

func joinPathFromCaller(pathElements ...string) (string, error) {
	// 获取调用者的文件名
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		return "", errors.New("failed to get caller information")
	}

	// 获取文件所在的目录
	callerDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	allPathElements := append([]string{callerDir}, pathElements...)
	dbConfigPath := filepath.Join(allPathElements...)

	// 将路径转换为绝对路径
	absPath, err := filepath.Abs(dbConfigPath)
	if err != nil {
		return "", errors.New("警告: 无法获取绝对路径")
	}

	// 简化路径，比如简化 ..
	simplifiedPath := filepath.Clean(absPath)
	return simplifiedPath, nil
}

/*
JoinPathFromCaller 根据调用文件的路径，构建一个相对于该项目的绝对路径。

参数 pathElements 是相对于调用文件所在目录的路径片段。

例如，在 `/path/to/project/handlers/api.go` 中调用：

	JoinPathFromCaller("..", "config", "app.yaml")

会返回：

	"/path/to/project/config/app.yaml"

注意：该函数使用 runtime.Caller 获取调用者文件路径，因此必须在运行时调用。
*/
func JoinPathFromCaller(pathElements ...string) (string, error) {
	return joinPathFromCaller(pathElements...)
}
