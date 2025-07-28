# 个人博客系统后端

基于 Go 语言、Gin 框架和 GORM 实现的博客系统后端

## 功能特性
- 用户认证（注册/登录）
- 文章管理（CRUD）
- 评论功能
- JWT 认证授权
- 请求日志记录

## 运行环境
- Go 1.18+
- MySQL 5.7+ 或 SQLite
- Gin 框架
- GORM

## 安装步骤
1. 克隆仓库：


## API 文档
| 端点 | 方法 | 描述 | 认证 |
|------|------|------|------|
| /auth/register | POST | 用户注册 | 否 |
| /auth/login | POST | 用户登录 | 否 |
| /posts | GET | 获取文章列表 | 否 |
| /posts/:id | GET | 获取单篇文章 | 否 |
| /posts | POST | 创建文章 | JWT |
| /posts/:id | PUT | 更新文章 | JWT + 作者 |
| /posts/:id | DELETE | 删除文章 | JWT + 作者 |
| /comments/post/:post_id | GET | 获取文章评论 | 否 |
| /comments | POST | 创建评论 | JWT |

## 测试说明
1. 使用 Postman 导入测试集合
2. 依次测试注册、登录、文章操作等接口
3. 验证权限控制（非作者无法修改/删除文章）

作业中遇到的几个问题
1:在开始作业之前，务必设计好项目结构，mvc框架，每个目录，每个go内容是什么，配置文件放config.go,env 存放 DSN JWT密钥 端口号等信息
2:提前规划路由，认证路由，文章路由，评论路由，需要注意的是，部分文章路由不需要JWT认证也可使用，例如查看文章等
还有些路由虽然经过了jwt 认证，但是需要二次调用中间件，确保jwt 认证的用户 和数据库中存储的用户一致
关于jwt认证，如何用postman模拟测试 搞了一阵子，参考的教程是https://www.cnblogs.com/simawenbo/p/13685836.html
3:ShouldBind直接绑定对象的时候，要求某些字段不能为空，比如绑定Post，但是传输过来的json没有User，这种情况就建议用 interface直接提取请求中的key value
4:gorm Updates Post当UserID为空 ，底层会改写为inser 导致更新失败，最后改写为 Updates 只更新固定字段解决
5:用户秘密 加密后存放在数据库中，但是查询的时候，不能通过针对秘密加密，然后用加密后的作为查询条件，因为每次加密后的结果都不一样
，需要将数据库存放的秘密取出，调用CompareHashAndPassword，进行对比判断