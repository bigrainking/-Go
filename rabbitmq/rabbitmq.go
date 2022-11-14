package rabbitmq

import (
	"Spike-Product-Demo/datamodels"
	"Spike-Product-Demo/services"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// MQ的URL
// URL固定的格式: amqp://账号：密码@rabbitmq服务器地址：端口号/vhost
const MQURL = "amqp://guest:guest@127.0.0.1:5672/rabbitVHost"

// 实例类
type RabbitMQ struct {
	// 连接channel
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// URL
	Mqurl string
	// 隧道锁
	sync.Mutex
}

//RabbitMQ实例创建函数
func NewRabbitMQ(exchange, queueName, key string) *RabbitMQ {
	// 返回一个创建的实例
	rabbitmq := &RabbitMQ{
		Exchange:  exchange,
		QueueName: queueName,
		Key:       key,
		Mqurl:     MQURL}
	// 建立连接
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, ":RabbitMQ 建立连接出错:")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

// 断开连接：类下面的功能function
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理:
// message是给人类看的错误原因
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", err, message) //打印到log
		panic(fmt.Sprintf("%s:%s", err, message))
	}
}

// =================simple=============
// Step1. 创建simple实例:简单实例只有队列
func NewRabbitMQsimple(queueName string) *RabbitMQ {
	// 使用默认exchange而不是没有交换机
	return NewRabbitMQ("", queueName, "")
}

// Step2. 生产者:生产消息
func (r *RabbitMQ) PublishSimple(message string) error {
	// 隧道加锁
	r.Lock()
	defer r.Unlock()
	// 1. 申请队列
	_, err := r.channel.QueueDeclare(
		// 要申请队列的名称
		r.QueueName,
		// durable：bool 持久化； false消息进来会在队列里面，如果服务器重启就没有了
		false,
		// autoDelete bool,消费者自动断开后是否删除消息
		false,
		// 队列是否有排他性exclusive bool,
		false,
		// noWait bool是否阻塞，发送消息后阻塞，等待服务器回应
		false,
		// 额外参数
		nil,
	)
	if err != nil {
		fmt.Printf("%s:%s", err, "生产者申请队列失败")
		return err
	}

	// 2. 生产消息
	r.channel.Publish(
		// exchange string, 交换机
		r.Exchange,
		// key string, routingKey 指定要发送到的queue
		r.QueueName,
		// mandatory bool, true根据交换机androutineKey判断是否有合适的queue，如果没有返回message给生产者
		false,
		// immediate bool, true 队列没有消费者，则返回给生产者
		false,
		// msg amqp.Publishing:要发送的消息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return nil
}

// Step3. 消费者：消费消息
func (r *RabbitMQ) ConsummerSimple(productService services.IProductService, orderService services.IOrderService) {
	// 1. 申请队列:与生产者队列相同，必须是同一个
	_, err := r.channel.QueueDeclare(
		// 要申请队列的名称
		r.QueueName,
		// durable：bool 持久化； false消息进来会在队列里面，如果服务器重启就没有了
		false,
		// autoDelete bool,消费者自动断开后是否删除消息
		false,
		// 队列是否有排他性exclusive bool,
		false,
		// noWait bool是否阻塞，发送消息后阻塞，等待服务器回应
		false,
		// 额外参数
		nil,
	)
	if err != nil {
		fmt.Printf("%s:%s", err, "消费者申请队列失败")
	}
	// 2. 接收消息
	msgs, err := r.channel.Consume(
		// queue string,
		r.QueueName,
		// consumer string, 消费者的名字：区分多个消费者
		"",
		// autoAck bool, 自动回应：true 收到消息立刻回应服务器删除消息，但如果消费失败，无法重新获取消息(消息已经被删除)
		false, //设置手动回应
		// exclusive bool,排他性，队列仅自己可见
		false,
		// noLocal bool, true 表示不能将消息传递给本connection中的另一个消费者
		false, //当当前消费者阻塞，可以传递给其他消费者
		// noWait bool, 消费队列是否阻塞
		false,
		// args amqp.Table额外参数
		nil,
	)
	// 消费者流控
	r.channel.Qos(
		1, //消费者一次能消费的最大容量
		0,
		false, //其他消费者不能消费channel里面的内容
	)
	// 3. 消费消息
	// 阻塞，直到人工主动停止接收消息
	forever := make(chan bool)
	// 协程 一直不断接收消息
	go func() {
		for msg := range msgs {
			// 处理消息:逻辑代码自行添加
			log.Printf("Recive a message:%s", msg.Body)
			// 接收消息
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(msg.Body), message) //接收消息，解析到message里面
			if err != nil {
				log.Fatalln("接收消息出错：Error：", err)
			}
			// 减少productNum
			err = productService.SubProductNum(message.ProductID)
			if err != nil {
				// log.Fatalln("consumer.go商品数量减少出错：Error：", err)：这样写如果出错会阻塞
				fmt.Println("consumer.go商品数量减少出错：Error：", err)
			}
			// 插入订单
			err = orderService.InsertOrderByMessage(message)
			if err != nil {
				// log.Fatalln("consumer.go创建订单出错：Error：", err)
				fmt.Println("consumer.go创建订单出错：Error", err)
			}
			// 通知完成
			//如果为true表示确认所有未确认的消息，
			//为false表示确认当前消息
			msg.Ack(false)
		}
	}()
	log.Printf("[*]Waiting for message， To Exit by ctrl+Enter")
	<-forever //channel中没有信息，一直阻塞，ctrl+c输入一个bool值到channel中，OK可以从channel中读出消息，执行完毕，程序结束
}
