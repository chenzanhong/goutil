package sshx

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// CreateSSHConn 创建 SSH 连接到指定地址
func CreateSSHConn(addr, user, auth string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(auth), // 注意：生产环境建议使用私钥
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境应使用安全的 HostKey 验证
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器 %s: %v", addr, err)
	}

	return client, nil
}

// UploadFile 上传本地文件到远程服务器
//
// 参数:
//   - localPath: 本地文件路径
//   - remotePath: 远程目标路径（包括文件名）
//   - serverAddr: 服务器地址，如 "192.168.1.100:22"
//   - user: SSH 用户名
//   - auth: SSH 密码
//
// 返回:
//   - error: 成功返回 nil，失败返回具体错误
func UploadFile(localPath, remotePath, serverAddr, user, auth string) error {
	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息（用于创建远程文件）
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取本地文件信息失败: %v", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("本地路径是目录，不能上传: %s", localPath)
	}

	// 创建 SSH 连接
	sshClient, err := CreateSSHConn(serverAddr, user, auth)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 复制文件内容
	_, err = io.Copy(remoteFile, file)
	if err != nil {
		return fmt.Errorf("上传文件内容失败: %v", err)
	}

	return nil
}

// DownloadFile 从远程服务器下载文件到本地
//
// 参数:
//   - remotePath: 远程文件路径
//   - localPath: 本地保存路径（包括文件名）
//   - serverAddr: 服务器地址，如 "192.168.1.100:22"
//   - user: SSH 用户名
//   - auth: SSH 密码
//
// 返回:
//   - error: 成功返回 nil，失败返回具体错误
func DownloadFile(remotePath, localPath, serverAddr, user, auth string) error {
	// 创建 SSH 连接
	sshClient, err := CreateSSHConn(serverAddr, user, auth)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer remoteFile.Close()

	// 获取远程文件信息，检查是否为目录
	fileInfo, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("获取远程文件信息失败: %v", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("远程路径是目录，无法下载: %s", remotePath)
	}

	// 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer localFile.Close()

	// 复制文件内容
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("下载文件内容失败: %v", err)
	}

	return nil
}
