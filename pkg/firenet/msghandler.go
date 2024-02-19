package firenet

import (
	"fire/pkg/fire"
	"k8s.io/klog/v2"
	"strconv"

	"fire/pkg/fireface"
)

// MsgHandle -
type MsgHandle struct {
	Apis             map[uint32]fireface.IRouter //存放每个MsgID 所对应的处理方法的map属性
	WorkerPoolSize   uint16                      //业务工作Worker池的数量
	MaxWorkerTaskLen uint16                      //业务工作Worker对应负责的任务队列最大任务存储数量
	TaskRequestQueue []chan fireface.IRequest    //Worker负责取任务的消息队列
}

// NewMsgHandle 创建MsgHandle
func NewMsgHandle(workerPoolSize, maxWorkerTaskLen uint16) *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]fireface.IRouter),
		WorkerPoolSize: workerPoolSize,

		MaxWorkerTaskLen: maxWorkerTaskLen,
		//一个worker对应一个queue
		TaskRequestQueue: make([]chan fireface.IRequest, workerPoolSize),
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request fireface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则
	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)

	//将请求消息发送给任务队列
	mh.TaskRequestQueue[workerID] <- request

}

// DoMsgRequestHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgRequestHandler(request fireface.IRequest) {
	klog.Infof("DoMsgRequestHandler = %d", request.GetMsgID())
	handler, ok := mh.Apis[fire.DefaultFireMsgID]
	if !ok {
		klog.Errorf("api msgID = ", fire.DefaultFireMsgID, " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.RequestHandle(request)
	handler.PostHandle(request)
}

// DoMsgDataHandler 马上以非阻塞方式处理数据
func (mh *MsgHandle) DoMsgDataHandler(request fireface.IRequest) {
	klog.Infof("DoMsgDataHandler = %d", request.GetMsgID())
	handler, ok := mh.Apis[fire.DefaultFireMsgID]
	if !ok {
		klog.Errorf("api msgID = ", fire.DefaultFireMsgID, " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.DataHandle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router fireface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	klog.Infof("Add api msgID = %d", msgID)
}

// StartOneRequestWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneRequestWorker(workerID int, taskQueue chan fireface.IRequest) {
	klog.Infof("Request Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgRequestHandler(request)
			mh.DoMsgDataHandler(request)
		}
	}
}

// StartWorkerPool 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskRequestQueue[i] = make(chan fireface.IRequest, mh.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneRequestWorker(i, mh.TaskRequestQueue[i])
	}
}
