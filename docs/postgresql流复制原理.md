# 简介

1、跨机器部署
2、容器运行环境
3、支持为现有的主库配置流复制从库: 备份(获取旧数据)+流复制(同步新数据)
3、全新的环境: 直接配置流复制(没有旧数据无需备份)

主从连接需要通过ssh进行处理

！！！从库完全仅用于备份、或者用于测试环境，不考虑主从切换

## 参考文档

[PostgreSQL流复制之一：原理+环境搭建（转发+整理）](https://www.cnblogs.com/yickel/p/11161706.html)

# 原理

postgres的主从复制分为两种，物理复制(流复制)以及逻辑复制。

物理复制用于复制整个数据库实例，通过wal日志来进行数据同步，分为异步和同步两种，支持从库为只读模式

![](./assets/%E6%B5%81%E5%A4%8D%E5%88%B6%E5%8E%9F%E7%90%86.png)

事务commit后，日志在主库写入wal日志，还需要根据配置的日志同步级别，等待从库反馈的接收结果。主库通过日志传输进程将日志块传给从库，从库接收进程收到日志开始回放，最终保证主从数据一致性。

## 同步级别

下面的级别从上至下降低

- remote_apply：事务commit或rollback时，等待其redo在primary、以及同步standby(s)已持久化，并且其redo在同步
standby(s)已apply。
- on：事务commit或rollback时，等待其redo在primary、以及同步standby(s)已持久化。
- remote_write：事务commit或rollback时，等待其redo在primary已持久化; 其redo在同步standby(s)已调用write接口(写到 OS, 但是还没有调用持久化接口如fsync)。
- local：事务commit或rollback时，等待其redo在primary已持久化;
- off：事务commit或rollback时，等待其redo在primary已写入wal buffer，不需要等待其持久化;

## wal日志

> 用于避免数据库中的事务一提交就刷新到数据库磁盘中，使用wal日志机制来减少磁盘性能影响并保持数据的持久性

以record为单位首先记录到wal日志中，在checkpoint时才对数据进行刷
盘(background writer会定时刷脏数据，但最终还是都由checkpoint确认都刷盘成功)。

- wal日志位置：$PGDATA/pg_wal(pg10之前叫pg_xlog)
- wal日志文件命名规则：`000000010000000100000092`, 前8位timeline，中8位logid，后8位logseg(前六位都是零，后两位是LSN的前两位，LSN可以用于计算在日志文件中的偏移量)

wal日志是二进制文件，可以使用pg_waldump这个工具来将其转换成可读的文件，wal日志中保存的是增量的sql语句

- rmgr: 这条记录所属的资源管理器
- len: wal记录的长度
- tx: 事务号
- lsn: 本条wal记录的lsn
- prev: 上条wal记录的lsn
- desc
- blkref

## 复制槽

在归档模式下wal日志，归档完成后会自动清理

提供了复制槽来避免主库在所有的备库收到 WAL 日志之前不会移除它们

创建复制槽，配置备库时使用`primary_slot_name = 'node_a_slot'`

## 免密登录

备份和还原过程中所用的archive_command和restore_command命令都以postgres用户运行

针对postgres用户实现ssh无密码登录。

## ^pg9.6

9.6开始，wal_level没有hot_standby（热备模式） 

## ^pg10

> 资源准备: 主数据库、从数据库、专门用于复制的数据库账号并配置权限
> 这些配置文件修改后都需要重新启动应用才会生效

*主库配置文件*

```conf
# postgresql.conf

# 异步模式
listen_address = '*' # 用于指定能够接收哪些ip的连接请求
wal_level = replica # 设置流复制模式至少设置为replica # minimal, replica, or logical
archive_mode = on 
archive_command = 'cp %p /data/postgresql/archive/%f '
max_wal_senders= 10  # 最大WAL发送进程数 要大于从库个数
wal_keep_segments=1024 # pg_wal目录下保留WAL日志的个数
hot_standby = on # 设置为ON后从库为只读模式

# 同步模式
synchronous_commit = remote_write、on、remote_apply
synchronous_standby_names = 'standby2'
```

*权限配置文件*

创建用户 

```shell
ceate user replica with replication login   password 'replication'; 

alter user replica with password 'replication'; 
```

在主库配置接受流复制的连接，修改pg_hba.conf文件，添加另一个备库的信息

```conf
# pg_hba.conf
host  replication  replica 127.0.0.1/32  md5
```

*从库配置文件*

```conf
# recovery.conf

# 异步模式
recovery_target_timeline = 'latest'
standby_mode = on
primary_conninfo = 'host=192.168.7.180 port=1921 user=bill password=xxx'

# 同步模式
primary_conninfo = 'host=192.168.7.180 port=1921 user=bill password=xxx application_name=standby2'
```

## ^pg12

PostgreSQL12中流复制有了一些改变,把recovery.conf的内容全部移入postgresql.conf，配置恢复、archive based standby、stream based standby，都在postgresql.conf中。postgresql.conf以及对应的两个signal文件来表示进入recovery 模式或standby模式。

### 恢复模式



```conf
# stream恢复模式配置  
#primary_conninfo = ''  
或  
# archvie恢复模式配置  
#restore_command = ''  
  
hot_standby = on   
  
# 配置是否跨时间线  
#recovery_target_timeline = 'latest'  
  
# 配置恢复目标，例如  
# 立即（达到一致性即停止恢复）、时间、XID、restore point name, LSN.  
#recovery_target = ''           # 'immediate' to end recovery as soon as a  
                                # consistent state is reached  
                                # (change requires restart)  
#recovery_target_name = ''      # the named restore point to which recovery will proceed  
                                # (change requires restart)  
#recovery_target_time = ''      # the time stamp up to which recovery will proceed  
                                # (change requires restart)  
#recovery_target_xid = ''       # the transaction ID up to which recovery will proceed  
                                # (change requires restart)  
#recovery_target_lsn = ''       # the WAL LSN up to which recovery will proceed  
                                # (change requires restart)  
#recovery_target_inclusive = on # Specifies whether to stop:  
                                # just after the specified recovery target (on)  
                                # just before the recovery target (off)  
                                # (change requires restart)  
  
# 恢复目标到达后，暂停恢复、激活、停库  
#recovery_target_action = 'pause'  # 'pause', 'promote', 'shutdown'  
                                   # (change requires restart) 
```

### standby模式

设置`hot_standby`为`on`

```conf
# stream恢复模式配置  
#primary_conninfo = ''  
或   
# archvie恢复模式配置  
#restore_command = ''  
  
hot_standby = on   
  
# 配置是否跨时间线  
#recovery_target_timeline = 'latest'  
```


### 基础备份

```shell
pg_basebackup -h 10.10.22.151 -p 5432 -U replica -W -R -Fp -Xs -Pv -D /data/postgresql-12/data/
```

## ^pg14

没有wal_keep_segments