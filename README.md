# estool
简介：一个简单的一个从文件导入到es的工具

### 用法
1. 用tool 生成数据  例如 ./tool --count=10000000
2. 用esImport导入数据, 如果有需要在配置文件esImport.yaml更改配置 例如 ./esImport --conf=esImport.yaml
3. 用esQuery查询数据,当前只提供一个时间查询 例如  ./esQuery --esUrl= --esUser=elastic --esPasswd= --indexName=log4es --size=10

### 编译
进主目录,执行bash build.sh,生成文件在build/bin中


### 不足
1.退避
2.换行符号支持更多种
3.中途失败只从失败前的数据开始导入（因为es的批量导入不太好记住上次失败行数,这里要做准确需要花点时间）
4.可以自定义buffer,节约内存,不然gc压力很大
5.注释少，测试覆盖不足
