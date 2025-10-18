// Example usage of pprof2 with default parameters
package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== pprof2 默认参数示例 ===")
	fmt.Println()

	fmt.Println("1. 基本用法（使用默认参数）:")
	fmt.Println("   go run main.go")
	fmt.Println("   - 默认显示前10个条目")
	fmt.Println("   - 按累积时间排序")
	fmt.Println("   - 隐藏小于0.5%的节点")
	fmt.Println()

	fmt.Println("2. 自定义参数:")
	fmt.Println("   go run main.go -nodecount=20 -cum=false")
	fmt.Println("   - 显示前20个条目")
	fmt.Println("   - 按扁平时间排序")
	fmt.Println()

	fmt.Println("3. 常用命令:")
	fmt.Println("   go run main.go -top                    # 显示top条目")
	fmt.Println("   go run main.go -list=main.main         # 显示main.main函数的源码")
	fmt.Println("   go run main.go -web                    # 启动web界面")
	fmt.Println("   go run main.go -http=:8080             # 在8080端口启动web服务器")
	fmt.Println()

	fmt.Println("4. 过滤选项:")
	fmt.Println("   go run main.go -focus=\"runtime\"        # 只关注runtime相关")
	fmt.Println("   go run main.go -ignore=\"vendor\"        # 忽略vendor目录")
	fmt.Println("   go run main.go -hide=\"test\"            # 隐藏测试相关")
	fmt.Println()

	fmt.Println("5. 输出格式:")
	fmt.Println("   go run main.go -output=profile.txt     # 输出到文件")
	fmt.Println("   go run main.go -granularity=lines      # 按行显示")
	fmt.Println("   go run main.go -call_tree=true         # 显示调用树")
	fmt.Println()

	fmt.Println("6. pprof2 增强功能:")
	fmt.Println("   go run main.go -http=localhost:8080   # 支持HTTP获取profile")
	fmt.Println("   go run main.go -insecure               # 支持不安全的HTTPS连接")
	fmt.Println("   go run main.go -timeout=30s            # 设置超时时间")
	fmt.Println()

	fmt.Println("默认参数配置:")
	fmt.Println("  -cum=true              # 按累积时间排序")
	fmt.Println("  -nodecount=10          # 显示前10个条目")
	fmt.Println("  -nodefraction=0.005    # 隐藏小于0.5%的节点")
	fmt.Println("  -edgefraction=0.001    # 隐藏小于0.1%的边")
	fmt.Println("  -sort=flat             # 默认排序方式")
	fmt.Println("  -granularity=functions # 默认粒度")
}
