package seed

import "gorm.io/gorm"

// 存放所有 Seeder
var seeders []Seeder

// 按顺序执行的 Seeder 数组
var orderedSeederNames []string

type SeederFunc func(*gorm.DB)

// Seeder 对应每一个 database/seeders 目录下的 Seeder 文件
type Seeder struct {
	Func SeederFunc
	Name string
}

// Add 注册到 seeders 数组中
func Add(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

// SetRunOrder 设置按顺序执行的 Seeder 数组
func SetRunOrder(names []string) {
	orderedSeederNames = names
}
