package builder

// 允许自定义配置参数的构造器
// 这样比定义一个参数巨多的 DBPool 构造函数要好一点
import (
	"fmt"
	"time"
)

// DB连接池提供了很多配置化的参数。
type DBPool struct {
	dsn             string
	maxOpenConn     int
	maxIdleConn     int
	maxConnLifeTime time.Duration
}

// 给 DB 连接池加一个建造者模式，这样在设置每个配置化参数的时候就可以对参数进行一步检查，
// 避免直接 new 连接池对象，在给每个属性赋值时都加判断，把每个参数的校验内聚到参数自己的建造者步骤里。

// --------------------定义建造者的结构体
type DBPoolBuilder struct {
	DBPool       // 参数检验通过后，逐步将配置赋值
	err    error // 如果某个步骤出现error，后面的步骤直接抛出错误，不做其他处理
}

// --------------------实例化建造者对象
func Builder() *DBPoolBuilder {
	return new(DBPoolBuilder)
}

// --------------------设置每个配置化参数
func (b *DBPoolBuilder) DSN(dsn string) *DBPoolBuilder {
	// 把在外部调用时的错误判断，分散到了每个步骤里。
	if b.err != nil {
		return b
	}
	// 做参数校验
	if dsn == "" {
		b.err = fmt.Errorf("invalid dsn, current is %s", dsn)
	}
	// 配置项赋值
	b.DBPool.dsn = dsn
	return b
}

func (b *DBPoolBuilder) MaxOpenConn(connNum int) *DBPoolBuilder {
	if b.err != nil {
		return b
	}
	if connNum < 1 {
		b.err = fmt.Errorf("invalid MaxOpenConn, current is %d", connNum)
	}

	b.DBPool.maxOpenConn = connNum
	return b
}

func (b *DBPoolBuilder) MaxConnLifeTime(lifeTime time.Duration) *DBPoolBuilder {
	if b.err != nil {
		return b
	}
	if lifeTime < 1*time.Second {
		b.err = fmt.Errorf("connection max life time can not litte than 1 second, current is %v", lifeTime)
	}

	b.DBPool.maxConnLifeTime = lifeTime
	return b
}

// -------------------返回连接池对象
func (b *DBPoolBuilder) Build() (*DBPool, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.DBPool.maxOpenConn < b.DBPool.maxIdleConn {
		return nil, fmt.Errorf("max total(%d) cannot < max idle(%d)", b.DBPool.maxOpenConn, b.DBPool.maxIdleConn)
	}
	return &b.DBPool, nil
}
