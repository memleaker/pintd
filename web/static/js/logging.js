layui.use(['table', 'layer'], function() {
    var table = layui.table;
    var layer = layui.layer;
    var $ = layui.jquery;

    // 加载table实例
    table.render({
        elem:"#logging",      // 绑定到表格ID
        url:"/logging/get",   // 获取数据的接口
        page:true,            // 开启分页
        toolbar:"#toolbar",   // 设置表格工具栏

        // 设置列
        cols:[[
            // 设置表头行号
            {fiele:'id', type:"numbers"},
            // 设置复选框
            {field:'select', type:"checkbox"},
            // 设置列
            {field:'id', title:'日志ID', width: 80},
            {field:'time', title:'日志产生时间', width:200},
            {field:'content', title:'日志内容'},
            // 删除按钮
            {field:'operate', title:'操作' ,toolbar:"#editbar", width:110},
        ]],
    });

    // 监听表头工具栏
    table.on('tool(table)', function(obj){
        var data = obj.data;

        // 删除操作
        if (obj.event == "del") {
            layer.confirm('确定删除此条日志吗?', function(index) {
                // 关闭询问弹窗
                layer.close(index);

                // 修改data, 去掉time和content再发给后端，减小JSON大小
                data['time'] = "";
                data['content'] = "";

                // 显示加载中
                layer.load();

                // 发送请求
                $.ajax({
                    url:"/logging/del",
                    type:"post",
                    data:JSON.stringify(data),  //提交的表单数据
        
                    // result 代表服务端返回的JSON, msg和success为JSON里的字段
                    success:function(result) {
                        if (result.success) {
                            layer.closeAll('loading'); // 关闭加载框
                            obj.del(); // 删除条目
                            layer.msg(result.msg, {icon: 6});  //返回数据成功时弹框
                        }
                    },
        
                    // 无返回或处理有报错时弹框
                    error:function(result) {
                        layer.closeAll('loading'); // 关闭加载框

                        // result.status 标志HTTP状态码，无响应时为0
                        if (result.status == 0) {
                            layer.alert('服务器无响应!!!', {icon: 2})
                        }else {
                            layer.msg(result.responseJSON.msg, {icon: 5});
                        }
                    }
                });
            });
        }
    });

    // 监听表格工具栏事件
    table.on('toolbar(table)', function(obj){
        var checkStatus = table.checkStatus(obj.config.id);
        var event = obj.event;
        var arr = checkStatus.data;

        if (event == "DelSelected") {
            if (arr.length == 0) {
                layer.msg('未选中任何数据');
                return
            }

            layer.confirm('确定删除选中的日志吗?', function(index) {
                // 关闭询问弹窗
                layer.close(index);

                // 显示加载中
                layer.load();

                // 保留日志和时间内容会让JSON比较大, 因此将其置位空
                for (let i = 0; i < arr.length; ++i){
                    arr[i]['content'] = "";
                    arr[i]['time'] = "";
                }

                // 发送请求
                $.ajax({
                    url:"/logging/delmore?num="+arr.length,
                    type:"post",
                    data:JSON.stringify(arr),  //提交的表单数据

                    // result 代表服务端返回的JSON, msg和success为JSON里的字段
                    success:function(result) {
                        if (result.success) {
                            layer.closeAll('loading'); // 关闭加载框
                            // 重新加载表格, 因为选中多个时不好删除
                            table.reload('logging', {
                                url:"/logging/get",
                            });
                            layer.msg(result.msg, {icon: 6});  //返回数据成功时弹框
                        }
                    },

                    // 无返回或处理有报错时弹框
                    error:function(result){
                        layer.closeAll('loading'); // 关闭加载框
                        
                        // result.status 标志HTTP状态码，无响应时为0
                        if (result.status == 0) {
                            layer.alert('服务器无响应!!!', {icon: 2})
                        }else {
                            layer.msg(result.responseJSON.msg, {icon: 5});
                        }
                    }
                });
            });
        }
    });
});