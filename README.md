# DevOps开发 golang 测试

1. 72小时内完成
2. fork本仓库
3. 通过kubebuilder或者手动创建golang的项目；完成自定义CRD MyReplicaSet核心功能的开发（功能同kubernetes ReplicaSet）
4. 同时还需要实现ValidatingAdmissionWebhook，为MyReplicaSet提供准入验证
5. 要求不能直接引用kubernetes ReplicaSet模块源码
6. 单元测试覆盖率必须大于80%
7. 必须包含controller的部署helm chart
8. makefile包含完整的编译流程（controller镜像的编译、controller helm chart编译、单元测试）
9. 完成以后通过pull request 提交，并备注面试姓名+联系方式，然后即时联系HR以免超时；

谢谢合作
