# ServiceComputing
## selpg程序设计

作业要求[传送门](https://pmlpml.github.io/ServiceComputingOnCloud/ex-cli-basic)

-----
## 程序介绍
selpg程序即 *select page* 程序，可以在文件内选择部分页面或页内的行输出 

## 命令说明
* 使用`-s number`指定选择的开始页，必须大于0
* 使用`-e number`指定选择的结束页，必须大于开始页
* 使用`-l number`指定每一页包含的行数，如不填写则使用默认值72
* 使用`-f`在页面中通过`\f`换行符作为页面分割的符号（使用`-f`将忽略`-l number`指令）
* 使用`-d destination`指定输出文件所在位置

## 使用示范
1. 首先进入文件所在位置使用命令`go build selpg.go`编译文件
2. 从测试文档test.txt中读取第一页的前10行内容   

    `./selpg -s 1 -e 1 -l 10 < test.txt`
3. 将测试文档test.txt中的第一页输入到另一空白文档out.txt中

    `./selpg -s 1 -e 1 -f < test.txt >out.txt`

