// package common

// import (
// 	"errors"
// 	"fmt"
// 	"hash/crc32"
// 	"sort"
// 	"strconv"
// 	"sync"
// )

// // 指定切片:实现接口
// // 使用sort.Sort需要实现接口，实现接口对应的三个function
// type units []uint32

// // 返回切片长度
// func (x units) Len() int {
// 	return len(x)
// }

// // 对比两个数字的大小
// func (x units) Less(i, j int) bool {
// 	return x[i] < x[j]
// }
// func (x units) Swap(i, j int) {
// 	x[i], x[j] = x[j], x[i]
// }

// var errEmpty = errors.New("Hash 环没有数据")

// // - 创建hash环
// type Consistent struct {
// 	// 环
// 	circle map[uint32]string //存储服务器key值and服务器信息
// 	// 虚拟节点个数
// 	VirtualNode int
// 	// 排序后的环的key对应的列表
// 	sortedHashes units
// 	// 读写锁:增删环上节点时不能读写
// 	sync.RWMutex
// }

// func NewConsistent() *Consistent {
// 	return &Consistent{
// 		circle:      make(map[uint32]string), //初始化一个环
// 		VirtualNode: 20,
// 	}
// }

// // - 【AddKey】向环中插入服务器节点and虚拟节点
// func (c *Consistent) Add(element string) {
// 	// 读写锁
// 	c.Lock()
// 	defer c.Unlock()
// 	c.add(element)
// }

// //   - 【generateKey】每个节点生成副本Key
// func (c *Consistent) generateKey(element string, index int) string {
// 	return element + strconv.Itoa(index)
// }

// //   - 【add】将节点添加到环上
// func (c *Consistent) add(element string) {
// 	// 1. 生成虚拟节点的副本 2. 副本添加到环上 3. 每个虚拟节点都指向服务器
// 	for i := 0; i < c.VirtualNode; i++ {
// 		c.circle[c.hashKey(c.generateKey(element, i))] = element
// 	}
// 	// 更新circle上新的hash值的排序
// 	c.updateSortedHashes()
// }

// // - 【Get】获取内容在哪台服务器上
// // name:要查找的对象的信息
// // element：存储name的服务器信息
// func (c *Consistent) Get(name string) (string, error) {
// 	// 查找时，添加读锁
// 	c.RLock()
// 	defer c.Unlock()
// 	if len(c.circle) == 0 {
// 		return "", errEmpty
// 	}
// 	// 计算hash值，返回对应的服务器
// 	fmt.Println("查找服务器：", name)
// 	key := c.hashKey(name) //在环上找到自己的位置
// 	index := c.seachKey(key)
// 	return c.circle[c.sortedHashes[index]], nil
// }

// //   - 【hashKey】计算内容对应的hashKey
// func (c *Consistent) hashKey(element string) uint32 {
// 	// 如果element不够容量
// 	if len(element) < 64 {
// 		var scratch [64]byte
// 		copy(scratch[:], element) //扩容
// 		return crc32.ChecksumIEEE(scratch[:len(element)])
// 	}
// 	return crc32.ChecksumIEEE([]byte(element))
// }

// //   - 【searchKey】在hash环上查找对应的节点，顺时针查找最近的节点
// // 传入对应的hashKey值, 返回在环上满足条件的节点，节点对应的sortedHashes的下标
// func (c *Consistent) seachKey(hashKey uint32) int {
// 	f := func(i int) bool {
// 		return c.sortedHashes[i] > hashKey
// 	}
// 	index := sort.Search(len(c.sortedHashes), f)
// 	// 如果index没有找到：超出了最大长度
// 	if index >= len(c.sortedHashes) {
// 		return 0
// 	}
// 	return index
// }

// // - 【Remove】删除服务器节点
// func (c *Consistent) Remove(element string) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.remove(element)
// }

// //   - 【remove】删除服务器节点
// // element：服务器节点信息ip
// func (c *Consistent) remove(element string) {
// 	// 删除节点:删除所有副本
// 	for i := 0; i < c.VirtualNode; i++ {
// 		delete(c.circle, c.hashKey(c.generateKey(element, i)))
// 	}

// 	// 更新排序hash
// 	c.updateSortedHashes()
// }

// // - 【updateSortedHashes】更新排sortedHashes，用于二分查找(sort.Search)
// func (c *Consistent) updateSortedHashes() {
// 	// 拷贝空串
// 	hashes := c.sortedHashes[:0]
// 	// 核对容量是否过大,过大则重置:切片容量 是 hash环上节点数*虚拟节点个数的4倍
// 	if cap(c.sortedHashes) > len(c.circle)*c.VirtualNode*4 {
// 		hashes = nil
// 	}
// 	// 将环上circle现有的节点添加到sortedHashes中
// 	for key := range c.circle {
// 		hashes = append(hashes, key)
// 	}
// 	// 将hashes排序
// 	sort.Sort(hashes)
// 	// 排序后更新到c.sortedHases
// 	c.sortedHashes = hashes
// }
package common

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//声明新切片类型
type units []uint32

//返回切片长度
func (x units) Len() int {
	return len(x)
}

//比对两个数大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

//切片中两个值的交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//当hash环上没有数据时，提示错误
var errEmpty = errors.New("Hash 环没有数据")

//1. ===================创建结构体保存一致性hash信息======================
type Consistent struct {
	//hash环，key为哈希值，值存放节点的信息 ： 这个节点是hash环上的每个节点？？
	circle map[uint32]string
	//已经排序的节点hash切片 ： 方便在hash环上查找图片对应的服务器
	sortedHashes units
	//虚拟节点个数，用来增加hash的平衡性：每台服务器都有VirtualNode个虚拟节点
	VirtualNode int
	//map 读写锁 ： 防止超卖
	sync.Mutex
}

//创建一致性hash算法结构体，设置默认节点数量
func NewConsistent() *Consistent {
	return &Consistent{
		//初始化变量
		circle: make(map[uint32]string),
		//设置虚拟节点个数
		VirtualNode: 20,
	}
}

//2. ==================自动生成服务器key值：根据服务器信息element生成服务器的Key============
func (c *Consistent) generateKey(element string, index int) string {
	//副本key生成逻辑
	return element + strconv.Itoa(index)
}

//3. ==================向环中添加服务器and虚拟节点===================
// 环上有N多个虚拟节点，但每个虚拟节点对应的服务器都是我们的实体节点
func (c *Consistent) add(element string) {
	//循环虚拟节点，设置副本
	// 传进来实体节点，依次生成他的虚拟节点， 虚拟节点对应的服务器都是实体节点
	for i := 0; i < c.VirtualNode; i++ {
		//根据生成的节点添加到hash环中
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}

// 外部使用添加节点，需要加锁
//向hash环中添加服务器：添加服务器
// element：IP
func (c *Consistent) Add(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.add(element)
}

//4. ===================获取hash位置,以便查找图片===============================
//传入key信息，找到key在hash环上的位置
func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组中
		copy(srcatch[:], key)
		//通过这个函数标准来计算hash值 ： 就想 key % 机器数量 = ？ 一样
		//使用IEEE 多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//更新排序，方便查找
func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0] //已经排序好的空切片，复位一下切片
	//判断切片容量，是否过大，如果过大则重置为空
	// 切片容量 是 hash环上节点数*虚拟节点个数的4倍
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil // 切片容量过大，重置
	}

	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	//对所有节点hash值进行排序，
	//方便之后进行二分查找
	sort.Sort(hashes)
	//重新赋值
	c.sortedHashes = hashes

}

//5. ===================删除节点============================
//删除节点：并且删除所有的副本
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

//删除一个节点
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用"二分查找"算法来搜索指定切片满足条件的最小值
	// 指定切片
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

//6.===================根据数据标示获取最近的服务器节点信息============
func (c *Consistent) Get(name string) (string, error) {
	fmt.Println("开始执行consistent获取最近的服务器节点信息：")
	// //添加锁
	// c.Lock()
	// //解锁
	// defer c.Unlock()
	//如果为零则返回错误
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//计算hash值
	// fmt.Println("获取用户所在服务器uid：", name)
	key := c.hashkey(name) //图片
	i := c.search(key)     //在hash环上根据数据标志，找到最近的服务器
	return c.circle[c.sortedHashes[i]], nil
}
