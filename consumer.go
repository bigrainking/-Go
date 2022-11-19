package main

// func main() {

// 	// =============== 1.链接数据库 ==============
// 	// 因为要单独部署，因此需要链接数据库

// 	db, err := common.NewMysqlConn()
// 	if err != nil {
// 		log.Fatalln("consumer 链接数据库错误， Error：", err)
// 	}

// 	// =================2. 创建simpleConsummer： 创建order product的操作实例===========
// 	productRepository := repository.NewProductManager("product", db)
// 	productService := services.NewIPoductSeviceManager(productRepository)
// 	orderRepository := repository.NewOrderManagerRepo("spikeSystem.order", db)
// 	orderService := services.NewOrderServiceManager(orderRepository)
// 	rabbitConsumSimple := rabbitmq.NewRabbitMQsimple("spikeProduct") //名字与生产端对齐
// 	// 消费消息
// 	rabbitConsumSimple.ConsummerSimple(productService, orderService)
// }
