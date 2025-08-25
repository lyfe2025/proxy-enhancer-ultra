package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 从数据库获取的密码哈希
	hashedPassword := "$2a$10$kFoHt8EVdMlg85N.cxmbLuDndS93DmO2FDkV3IZoxJiCrrvgPyYq."
	// 用户输入的密码
	password := "admin123456"
	
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
	} else {
		fmt.Println("密码验证成功!")
	}
	
	// 同时测试生成新的哈希
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成哈希失败: %v\n", err)
	} else {
		fmt.Printf("新生成的哈希: %s\n", string(newHash))
	}
}