#### architecture-diagram

![这里写图片描述](http://img.blog.csdn.net/20170319141236514?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMzgxMjcxMA==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

#### mini-paas的项目架构说明:
 主要由三层组成：web端，api接口服务端，后台数据存储
 - web层：主要是dashboard组件，主要是用来进行相关功能的操作界面
 - api组件：组要由apiserver，docker-build，registry，monitor，author五大组件组成。apiserver主要做的是关于应用部署相关的api，docker-build主要做的就是实现应用的在线/离线构建成镜像，registry主要是对registry镜像仓库封装的一些api，monitor主要是做监控的组件，author主要是做权限管理的组件。
 - 存储层：mysql和etcd
