# newworld
test
	一个服务更像一座城
	net 便是城市和外面交互的模块
		net.server 是对外开方的城门，可以有多个门，也可以有一个门,有门就需要有管理门的
		net.manager就是对门的管理
		net.conn 就是对各个城门外修的道路，可以有多条道路通向一个城门
		net.client 就是各个路口的收费站
		net.msg 就是每个人来城里办事情的事件，走在那条路上就指定了从哪个城门进城
	
	module 便是城市内的服务模块
		module.manager 就是对各个模块的监管部门，想开商户就得来这里注册
		module.mod 就是各个商户了
			各个商户又分不同的服务方式
			1 线性  	所有来这里请求服务的人，都要排队
			2 非线性	所有来的人都是自助模式，不用排队互不影响
		module.route 模块 就是对各个用户的分流，你是去哪个商户，就指定给分配给相应的商户
		
	data 就是城市内对各个人提供的一个托管仓库，不必每次都带着行李来城里，可以直接托管在这里，
		 每次轻装上阵来城里了，需要什么去这里取
		data.manager 就是对每个人的仓库的管理部门
		{
			由模块自动生成
		}
	
	admin 对整个城市的管控模块
		admin.admin 监管每个 manager， 各个模块的状态都要汇报给admin 
			开启模块，运行模块，关闭模块等控制
			模块问题监控，记录，分析
		各个模块是否安全关闭退出
	requst
	{
		外部请求服务，通过net 筛选绑定服务
		内部请求服务，通过route 筛选绑定服务
		{
			所以一个服务注册需要注册到route 和net
		}
		{
			服务接口也由模块自动生成
		}
		
		各个服务有自己的唯一id  本来是uint 类型 ，但是自动注册不太好做
		golang 的map 以 uint 和string 做key 查找耗时 基本是 1：2
		但是1000以内的服务数量，string 做key 查找亿次耗时1秒多，不算是水桶最矮的那块板
	}