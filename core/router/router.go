package router

type router interface {
	/**
	 * 添加一个处理器
	 */
	Add(url string, restHandler *interface{})
	/**
	 * 根据url获取一个最佳处理器
	 */
	Get(url string) *interface{}
}
