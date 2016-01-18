#ndb-go#

ndb(node database) 是一种轻量级嵌入式的文档型数据库, ndb-go是ndb数据库的Go语言实现 

- 可扩展的数据格式（Flexible Data Model）
- 使用表达式查询语言(Expressive Query Language)
- 采用程序嵌入的方式（Embedded Program）

下载路径：ndb-1.1.zip

## ndb语法规则 ##

**ndb 语法格式：**
> command:path->query:value

command是操作指令，包括select、one、update、delete、insert、delete<br/>
path是检索路径，格式为：node->node，通过层级关系进行指定检索路径<br/>
item是检索项，value是检索值

**查询列表(select)**
> select:path->item1:value1 && item2:value2

检索路径(path)和检索值(value) 支持正则表达式，数值区域，startwith 和 endwith：

- 正则表达式使用 / 进行分割，例如 /[aA].*/
- 数值区域仅对整数区域有效，使用 [min, max] 描述，例如 [12,17]
- startwith使用^开头，例如^hello
- endwith使用$结尾，例如hello$

select的返回值是一个列表，其中包括全部符合查询路径和条件的节点

**查询单个(one)**
> one:path->item1:value1 && item2:value2

one的使用方法与select一致，当出现多个结果时，仅返回第一个结果

**修改(update)**
> update:node->node->item:value !! item1=value1,item2=value2

item : value用于检索对应的node，根据后面的值对节点中的数据进行修改

**删除(delete)**
> delete:path->item : value !! [item1, item2] / block

如果使用item 标注，则仅删除item条目，多个item使用","分割
如果使用block标注，则删除item的父node项目

**新增(insert)**
> insert:path && item1=new value1,item2=new value

在节点中新建item，如果父路径不存在，则自动建立父路径


## ndb使用样例 ##

举例说明（数据样例）：

    root {
    	parent {
			name : green
    		child [{
    				name: jim
    				age: 20
    				sex: male
    				}, {
    				name: lily
    				age: 17
    				sex: female
    				}, {
    				name: tom
    				age: 28
    				sex: male
    			}]
			nephew {
				name: lucy
				age: 12
				sex: female
			}
    	}
		
    }


**one**

获得sex为male的第一个数据

	one:root->parent->child->sex:male

**select**

获得name为jim和tom的child

	select:root->parent->child->name:/.*m/

获得age在15-25之间的child

	select:root->parent->child->age:[15,25]  

获得sex以fe开头的child

	select:root->parent->child->sex:^fe
 
获得name为m结尾的child  

	select:root->parent->child->name:m$  

获得sex为male，age在15-25之间的child

	select:root->parent->child->sex:male && age:[15,25] 

获得root->parent下所有的child  

	select:root->parent->child   

获得root->parent中sex为female的child和nephew  

	select:root->parent->:/child|nephew/->sex:female 

**delete**

删除name为jim的sex和age属性：

	delete:root->parent->child->name:jim !! [sex, age]  

删除name为jim的所有属性，其中block为可选参数：

	delete:root->parent->child->name:jim !! block 

**update**

修改name为jim的age值，并增加address值
	
	update:root->parent->child->name:jim !! age=21, address=China 

**insert**

增加name为bill的child

	insert:root->parent->child !! name=bill, sex=male, age=31


## 在go中使用ndb ##
	package main
	
	import (
		"github.com/node-db/ndb-go/ndb"
		"fmt"
	)
	
	func main() {
		node, _ := ndb.Read("d:/example.ndb")
		result, found, err := ndb.Execute(node, "select:root->parent->child")
		if found && err == nil {
			children, ok := result.([]*ndb.Node)
			if ok {
				for _, child := range children {
					fmt.Println(child.GetName() + ":" + child.GetValue("name")[0])
				}
			}
		}
		
		query := "select:root->parent->:/child|nephew/->sex:female >> print"
		ndb.Execute(node, query)
	}
		
## 命令行中使用ndb ##

查询ndb中节点信息

	ndb_cli -f example.ndb -q "select:root->parent->child->name:/.*m/"
	
输出help信息
	
	ndb_cli -h
	
