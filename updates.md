reactor最重要一条是不能让conn 读写堵塞event pool
加入 goroutine pool
异步io

Reactor 网络模型
目前 Linux 平台上主流的高性能网络库/框架中，大都采用 Reactor 模式，比如 netty、libevent、libev、ACE，POE(Perl)、Twisted(Python)等。

Reactor 模式本质上指的是使用 I/O 多路复用(I/O multiplexing) + 非阻塞 I/O(non-blocking I/O) 的模式。

通常设置一个主线程负责做 event-loop 事件循环和 I/O 读写，通过 select/poll/epoll_wait 等系统调用监听 I/O 事件，业务逻辑提交给其他工作线程去做。而所谓『非阻塞 I/O』的核心思想是指避免阻塞在 read() 或者 write() 或者其他的 I/O 系统调用上，这样可以最大限度的复用 event-loop 线程，让一个线程能服务于多个 sockets。在 Reactor 模式中，I/O 线程只能阻塞在 I/O multiplexing 函数上（select/poll/epoll_wait）。

Reactor 模式的基本工作流程如下：

Server 端完成在 bind&listen 之后，将 listenfd 注册到 epollfd 中，最后进入 event-loop 事件循环。循环过程中会调用 select/poll/epoll_wait 阻塞等待，若有在 listenfd 上的新连接事件则解除阻塞返回，并调用 socket.accept 接收新连接 connfd，并将 connfd 加入到 epollfd 的 I/O 复用（监听）队列。
当 connfd 上发生可读/可写事件也会解除 select/poll/epoll_wait 的阻塞等待，然后进行 I/O 读写操作，这里读写 I/O 都是非阻塞 I/O，这样才不会阻塞 event-loop 的下一个循环。然而，这样容易割裂业务逻辑，不易理解和维护。
调用 read 读取数据之后进行解码并放入队列中，等待工作线程处理。
工作线程处理完数据之后，返回到 event-loop 线程，由这个线程负责调用 write 把数据写回 client。