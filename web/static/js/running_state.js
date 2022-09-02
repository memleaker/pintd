layui.use(['table', 'layer'], function() {
    var table = layui.table;
    var layer = layui.layer;
    var $ = layui.jquery;

    // 为卡片面板赋值
    $.ajax({
        url:"/running/runninginfo",
        type:"get",

        success:function(result) {
            // 协程数量
            var div = document.getElementById("goroutine_num");
            div.textContent += "当前协程数: " + result.goroutine_num;

            // 连接数量
            div = document.getElementById("redirect_conns");
            div.textContent += "当前连接数: " + result.redirect_conns;

            // CPU使用
            div = document.getElementById("cpu_usage");
            div.textContent += "进程占用CPU: " + result.cpu_usage;
            
            // 内存使用
            div = document.getElementById("mem_usage");
            div.textContent += "进程占用内存: " + result.mem_usage;
        },
    });

    // 加载table实例
    table.render({
        elem:"#running_state",  // 绑定到表格ID
        url:"/running/state",   // 获取数据的接口
        page:true,                  // 开启分页
        toolbar:"#toolbar",         // 设置表格工具栏

        // 设置列
        cols:[[
            // 设置表头行号
            {fiele:'num', type:"numbers"},
            // 设置列
            {field:'id', title:'ID', width:80},
            {field:'protocol', title:'协议', width:80},
            {field:'src-addr', title:'源地址', width:170},
            {field:'src-port', title:'源端口', width:100},
            {field:'dest-addr', title:'目的地址', width:170},
            {field:'dest-port', title:'目的端口', width:100},
            {field:'forward-flow', title:'转发的流量', width:130},
            {field:'realtime-flow', title:'实时流量', width:120},
            {field:'create-time', title: '创建时间'},
            // 设置终止按钮
            {field:'operate', title:'操作' ,toolbar:"#editbar", width:105},
        ]],
    });

    // 监听表格工具栏事件
    table.on('toolbar(table)', function(obj) {
        var select = document.getElementById('select_filter');
        var select_index = select.selectedIndex;
        var select_value = select.options[select_index].value;
        var search_key = $("#key").val();

        // 如果为搜索
        if (obj.event == "LAYTABLE_SEARCH") {
            if (search_key == "") {
                layer.msg("请输入搜索关键字!", {icon: 5});
                return
            }

            if (select_value == "") {
                layer.msg("请选择过滤条件!", {icon: 5});
                return
            }

            table.reload('running_state', {
                url:"/running/search?field="+select_value
                ,where: {
                      key: $("#key").val()
                }
                ,page: {
                    curr: 1 //重新从第 1 页开始
                }
            });
        }
    });


    // 监听表头工具栏
    table.on('tool(table)', function(obj){
        var data = obj.data;

        // terminate
        if (obj.event == "terminate") {
            layer.confirm('确定终止此任务吗?', function(index) {
                // 关闭询问弹窗
                layer.close(index);

                // 减小JSON长度
                data['forward-flow'] = "";
                data['realtime-flow'] = "";
                data['create-time'] = "";

                console.log(JSON.stringify(data));

                // 显示加载中
                layer.load();

                // 发送请求
                $.ajax({
                    url:"/running/terminate",
                    type:"post",
                    data:JSON.stringify(data),  //提交的表单数据

                    // result 代表服务端返回的JSON, msg和success为JSON里的字段
                    success:function(result) {
                        if (result.success) {
                            layer.closeAll('loading'); // 关闭加载框
                            obj.del();                 // 删除条目
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
