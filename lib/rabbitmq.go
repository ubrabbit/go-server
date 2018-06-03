package lib

//https://segmentfault.com/a/1190000010516906
//http://www.damonyi.cc/%E5%9F%BA%E4%BA%8Egolang%E5%AE%9E%E7%8E%B0%E7%9A%84rabbitmq-%E8%BF%9E%E6%8E%A5%E6%B1%A0/

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

import (
	rabbitmq "github.com/streadway/amqp"

	. "github.com/ubrabbit/go-server/common"
	. "github.com/ubrabbit/go-server/config"
)

type RabbitMsgReceiver struct {
	Name     string
	Receiver func(b []byte) (bool, error)
}

type RabbitMQSession struct {
	sync.Mutex

	Name         string
	Exchange     string
	ExchangeType string
	Conn         *rabbitmq.Connection
	Channel      *rabbitmq.Channel
	Queue        *rabbitmq.Queue

	waitPushConfirm chan rabbitmq.Confirmation
	OnPushConfirm   Functor

	OnReceiver  func([]byte, bool, error) (bool, bool)
	MsgReceiver map[string]*RabbitMsgReceiver

	connectRetry int
	shutdown     bool
	closed       bool
	err          error
	errNotify    string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitRabbitMQ() *RabbitMQConfig {
	cfg := GetRabbitMQConfig()
	if cfg == nil {
		fmt.Println("rabbitmq conf is nil, use default setting")
		cfg = new(RabbitMQConfig)
		cfg.Account = "guest"
		cfg.Password = "guest"
		cfg.Host = "127.0.0.1"
		cfg.HostName = ""
		cfg.Port = 5672
	}
	return cfg
}

func ConnRabbitMQ() (*rabbitmq.Connection, error) {
	cfg := InitRabbitMQ()
	account := cfg.Account
	password := cfg.Password
	host := cfg.Host
	hostname := cfg.HostName
	port := cfg.Port

	// 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
	server := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", account, password, host, port, hostname)
	conn, err := rabbitmq.Dial(server)
	//failOnError(err, "Failed to connect to RabbitMQ")
	return conn, err
}

func NewRabbitMQSession(q string) *RabbitMQSession {
	session := new(RabbitMQSession)
	session.Name = q
	session.Exchange = q
	session.ExchangeType = "fanout"

	session.Conn = nil
	session.Channel = nil
	session.Queue = nil

	session.connectRetry = 0
	session.shutdown = false
	session.closed = true
	session.err = nil
	session.errNotify = ""

	session.waitPushConfirm = nil
	session.OnPushConfirm = nil

	session.OnReceiver = nil
	session.MsgReceiver = make(map[string]*RabbitMsgReceiver)
	return session
}

func (session *RabbitMQSession) setError(err error, desc string) {
	session.closed = true
	session.err = err
	session.errNotify = desc
}

func (session *RabbitMQSession) SetConnectRetry(i int) {
	session.connectRetry = i
}

func (session *RabbitMQSession) SetExchange(name string, t string) {
	last, last_t := session.Exchange, session.ExchangeType
	session.Exchange = name
	session.ExchangeType = name
	if last != name || last_t != t {
		fmt.Println("reconnect by exchange changed")
		session.reconnectSession()
	}
}

func (session *RabbitMQSession) GetExchange() (string, string) {
	return session.Exchange, session.ExchangeType
}

func (session *RabbitMQSession) Connect() error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Connect Err:  ", err)
		}
	}()

	conn, err := ConnRabbitMQ()
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	ex, extype := session.GetExchange()
	if err = channel.ExchangeDeclare(
		ex,     // name of the exchange
		extype, // type
		true,   // durable
		false,  // delete when complete
		false,  // internal
		false,  // noWait
		nil,    // arguments
	); err != nil {
		return err
	}

	//channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
	queue, err := channel.QueueDeclare(
		session.Name, // name
		true,         // durable    持久化标识
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	session.closed = false
	session.err = nil
	session.errNotify = ""
	session.Conn = conn
	session.Channel = channel
	session.Queue = &queue

	session.waitPushConfirm = channel.NotifyPublish(make(chan rabbitmq.Confirmation, 1))
	if session.OnPushConfirm != nil {
		session.SetPushConfirm(session.OnPushConfirm)
	}
	return nil
}

func (session *RabbitMQSession) reconnectSession() (bool, error) {
	if session.IsShutdown() {
		return false, nil
	}
	session.Lock()
	defer session.Unlock()

	session.Close()
	retry := 0
	for {
		if session.connectRetry > 0 {
			if retry >= session.connectRetry {
				break
			}
			retry++
		}
		time.Sleep(time.Second * 1)
		fmt.Println("try connect rabbitmq")
		err := session.Connect()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("success connect rabbitmq")
		return true, nil
	}
	return false, errors.New("fail reconnect rabbitmq")
}

func (session *RabbitMQSession) Ping() error {
	if session.Conn == nil || session.Channel == nil {
		return rabbitmq.ErrClosed
	}

	channel := session.Channel
	err := channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
	if err != nil {
		//fmt.Println("11111 ", err)
		return err
	}

	msgContent := "ping.ping"
	err = channel.Publish("ping.ping", "ping.ping", false, false, rabbitmq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})
	if err != nil {
		//fmt.Println("22222 ", err)
		return err
	}

	err = channel.ExchangeDelete("ping.ping", false, false)
	//fmt.Println("33333 ", err)
	return err
}

func (session *RabbitMQSession) AddReceiver(name string, f func(b []byte) (bool, error)) bool {
	session.Lock()
	defer session.Unlock()

	o_recv := new(RabbitMsgReceiver)
	o_recv.Name = name
	o_recv.Receiver = f

	session.OnReceiver = nil
	session.MsgReceiver[name] = o_recv
	return true
}

func (session *RabbitMQSession) RemoveReceiver(name string) bool {
	session.Lock()
	defer session.Unlock()

	_, ok := session.MsgReceiver[name]
	if ok {
		delete(session.MsgReceiver, name)
		return true
	}
	return false
}

func (session *RabbitMQSession) ExecuteReceiver(body []byte) (bool, error) {
	session.Lock()
	defer func() {
		session.Unlock()
		err := recover()
		if err != nil {
			fmt.Println("ExecuteReceiver Error: ", err)
		}
	}()
	for name := range session.MsgReceiver {
		recever := session.MsgReceiver[name]
		_, err := recever.Receiver(body)
		if err != nil {
			fmt.Println("Receiver Error: ", recever.Name, err)
		}
	}
	return true, nil
}

func (session *RabbitMQSession) RegisterOnReceiver(f func([]byte, bool, error) (bool, bool)) {
	session.Lock()
	defer session.Unlock()

	session.OnReceiver = f
}

func (session *RabbitMQSession) publishMsg(data []byte, save bool) error {
	if session.Channel == nil {
		return fmt.Errorf("publishMsg Error by channel is nil")
	}

	exchange, exchangeType := session.GetExchange()
	if err := session.Channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	mod := rabbitmq.Transient //temp save
	if save {
		mod = rabbitmq.Persistent // 持久化标记
	}
	err := session.Channel.Publish(
		exchange,           // exchange
		session.Queue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		rabbitmq.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: mod,
		})
	return err
}

func (session *RabbitMQSession) confirmPushMsg() {
	log.Printf("waiting for confirmMsg")
	for {
		confirmed, ok := <-session.waitPushConfirm
		if !ok {
			log.Printf("break confirmMsg by waitPushConfirm closed")
			break
		}
		session.OnPushConfirm.Call(confirmed.DeliveryTag, confirmed.Ack)
	}
}

func (session *RabbitMQSession) SetPushConfirm(f Functor) {
	if session.Channel == nil {
		session.OnPushConfirm = f
		return
	}
	log.Printf("enabling publishing confirms.")
	if err := session.Channel.Confirm(false); err != nil {
		fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		return
	}
	if f == nil {
		session.OnPushConfirm = nil
		return
	}
	session.OnPushConfirm = f
	go session.confirmPushMsg()
}

func (session *RabbitMQSession) PushMsg(data []byte) bool {
	for {
		err := session.publishMsg(data, true)
		if err != nil {
			fmt.Println("pushMsg Error: ", err)
			succ, err := session.reconnectSession()
			if succ {
				continue
			}
			msg := JoinString("", "Failed to push message: ", string(data))
			session.setError(err, msg)
			if !session.IsShutdown() {
				log.Fatalf("%s reason: %s", msg, err.Error())
			}
			break
		}
		return true
	}
	return false
}

func (session *RabbitMQSession) PushMsgNoSave(data []byte) bool {
	for {
		err := session.publishMsg(data, false)
		if err != nil {
			fmt.Println("pushMsg Error: ", err)
			succ, err := session.reconnectSession()
			if succ {
				continue
			}
			msg := JoinString("", "Failed to push message: ", string(data))
			session.setError(err, msg)
			break
		}
		return true
	}
	return false
}

func (session *RabbitMQSession) Close() error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if session.Conn == nil {
		return nil
	}
	// will close() the deliveries channel
	if err := session.Channel.Cancel(session.Name, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	if err := session.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	session.closed = true
	return nil
}

func (session *RabbitMQSession) IsClosed() bool {
	if session.IsShutdown() {
		return true
	}
	if session.closed {
		return true
	}

	err := session.Ping()
	if err != nil {
		session.setError(err, "Ping Error")
		session.Close()
		return true
	}
	return false
}

func (session *RabbitMQSession) Shutdown() {
	session.Lock()
	defer session.Unlock()
	session.shutdown = true
	session.Close()
}

func (session *RabbitMQSession) IsShutdown() bool {
	return session.shutdown
}

func (session *RabbitMQSession) ConsumeMsg() {
	for {
		if session.IsShutdown() {
			fmt.Println("consume finished by shutdown")
			break
		}
		ok, _ := session.reconnectSession()
		if !ok {
			fmt.Println("consume finished by connect failure")
			break
		}

		ex, _ := session.GetExchange()
		if err := session.Channel.QueueBind(
			session.Name, // name of the queue
			"",           // bindingKey
			ex,           // sourceExchange
			false,        // noWait
			nil,          // arguments
		); err != nil {
			msg := "Failed to Queue Bind"
			fmt.Println(msg)
			session.setError(fmt.Errorf("Queue Bind: %s", err), msg)
			continue
		}

		ch_msgs, err := session.Channel.Consume(
			session.Queue.Name, // queue
			session.Queue.Name, // consumerTag,
			false,              // auto-ack
			false,              // exclusive
			false,              // no-local
			false,              // no-wait
			nil,                // args
		)
		if err != nil {
			msg := "Failed to create Consume"
			fmt.Println(msg)
			session.setError(err, msg)
			continue
		}
		fmt.Println(">>>>>>>>>>> start consumer")
		for msg := range ch_msgs {
			log.Printf(
				"got %dB delivery: [%v] %q",
				len(msg.Body),
				msg.DeliveryTag,
				msg.Body,
			)

			succ, err := session.ExecuteReceiver(msg.Body)
			ack := true
			requeue := false
			if session.OnReceiver != nil {
				ack, requeue = session.OnReceiver(msg.Body, succ, err)
			}
			if ack {
				// 确认收到本条消息, multiple必须为false
				msg.Ack(false)
			} else {
				msg.Nack(false, requeue)
			}
		}
		fmt.Println(">>>>>>>>>>> finished consumer")
	}
}
