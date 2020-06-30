# Net
	service{
		in 服务调用请求消息内容
		out 返回结果内容
	}
	
	client{
		服务器间通信连接
		一个监听端口可以建立多个链接
		ServerA<--->ServerB conn1
		ServerA<--->ServerB conn2
		这样conn1 和conn2 功能是一样的，算是负载均衡
		需要进行负载均衡管理
		每个conn 都会注册一个mod 模块
		类型一样的mod 会有相同的 funclist
		调用服务选择哪个mod 需要route 去处理这样就会稍微麻烦一些
		如果让 funclist 均衡的分配到几个 conn 那样就各走各的，相当于分区管理，
		而且conn 数量是由链接发起方主动申请建立的，数量控制可以通过消息来确认
		也可以走互相建立链接申请，每个模块建立单独的 conn ,这样会有多个链接，
		两个服务器间的连接数 相当于两个服务器的func_mod 的总和
		
	}