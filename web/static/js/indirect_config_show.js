layui.use(['table', 'layer'], function() {
    var table = layui.table;
    var layer = layui.layer;
    var $ = layui.jquery;

    // 加载table实例
    table.render({
        elem:"#indirect_cfg_show",  // 绑定到表格ID
        url:"/indirect/cfg_show",   // 获取数据的接口
        page:true,                  // 开启分页
        toolbar:"#toolbar",         // 设置表格工具栏

        // 设置列
        cols:[[
            // 设置表头行号
            {fiele:'id', type:"numbers"},
            // 设置单选框
            {field:'select', type:"radio"},
            // 设置列
            {field:'protocol', title:'协议', width:60},
            {field:'listen-addr', title:'监听地址', width:170},
            {field:'listen-port', title:'监听端口', width:100},
            {field:'dest-addr', title:'目的地址', width:170},
            {field:'dest-port', title:'目的端口', width:100},
            // 除了五元组开启表格可编辑
            {field:'acl', title:'访问控制', width:90, edit:'text'},
            {field:'admit-addr', title:'白名单', width:120, edit:'text'},
            {field:'deny-addr', title:'黑名单', width:120, edit:'text'},
            {field:'max-conns', title:'最大连接数', width:100, edit:'text'},
            {field:'memo', title:'备注', edit:'text'},
            // 设置编辑删除按钮
            {field:'operate', title:'操作' ,toolbar:"#editbar", width:70},
        ]],
    });

    /* 监听表格工具栏事件, demo表示待监听容器的lay-filter值
     * 语法 table.on('toolbar(demo)', function(obj){});
    */
    table.on('toolbar(table)', function(obj){
        var checkStatus = table.checkStatus(obj.config.id);
        var event = obj.event;
        var dataobj = checkStatus.data;
        var arr = dataobj[0];

        // 1. 不处理右侧几个按钮事件，让其自动处理
        if (event == "LAYTABLE_COLS" || event == "LAYTABLE_EXPORT"
            || event == "LAYTABLE_PRINT") {
            return
        }

        // 2. 处理自定义事件
        if (!arr) {
            layer.msg('未选中任何数据');
            return
        }

        switch (event) {
            case "getConnInfo":
                // 获取五元组
                layer.alert(arr['protocol']    + " " + arr['listen-addr'] + " " + 
                            arr['listen-port'] + " " + arr['dest-addr']   + " " + 
                            arr['dest-port']);
                break;
            case "getWhiteList":
                // 获取白名单
                layer.alert(arr['admit-addr']);
                break;
            case "getBlackList":
                // 获取黑名单
                layer.alert(arr['deny-addr']);
                break;
            case "getMemo":
                // 获取备注
                layer.alert(arr['memo']);
                break;
        }
    });

    // 监听表头工具栏
    table.on('tool(table)', function(obj){
        var data = obj.data;

        // 删除操作
        if (obj.event == "del") {
            layer.confirm('确定删除此配置吗? 删除配置不会影响已经建立的连接', function(index) {
                // 关闭询问弹窗
                layer.close(index);

                // 显示加载中
                layer.load();

                // 发送请求
                $.ajax({
                    url:"/indirect/cfg_del",
                    type:"post",
                    data:JSON.stringify(data.field),  //提交的表单数据

                    // result 代表服务端返回的JSON, msg和success为JSON里的字段
                    success:function(result) {
                        if (result.success) {
                            layer.closeAll('loading'); // 关闭加载框
                            obj.del(); // 删除条目
                            layer.msg(result.msg, {icon: 6});  //返回数据成功时弹框
                        }
                        else {
                            layer.closeAll('loading'); // 关闭加载框
                            layer.msg(result.msg, {icon: 5}); //返回数据失败时弹框
                        }
                    },
        
                    // 无返回或处理有报错时弹框
                    error:function(result){
                        layer.closeAll('loading'); // 关闭加载框
                        layer.alert('服务器无响应!!!', {icon: 2})
                    }
                });
            });
        }
    });

    // 监听单元格编辑事件
    table.on('edit(table)', function(obj){
        var field = obj.field; // 修改的字段名
        var data = obj.data; // 全部数据
        var old=$(this).prev().text(); //字段的旧值
        var success = false;

        // 加载中
        layer.load();
       
        // 发送请求
        $.ajax({
            url:"/indirect/cfg_edit?field="+field,
            type:"post",
            data:JSON.stringify(data.field),  //提交的表单数据

            // result 代表服务端返回的JSON, msg和success为JSON里的字段
            success:function(result) {
                if (result.success) {
                    success = true;
                    layer.closeAll('loading'); // 关闭加载框
                    layer.msg(result.msg, {icon: 6});  //返回数据成功时弹框
                }
                else {
                    layer.closeAll('loading'); // 关闭加载框
                    layer.msg(result.msg, {icon: 5}); //返回数据失败时弹框
                }
            },

            // 无返回或处理有报错时弹框
            error:function(result){
                layer.closeAll('loading'); // 关闭加载框
                layer.alert('编辑表格失败，服务器无响应!!!', {icon: 2})
            }
        });

        if (!success) {
            // 修改失败，界面上修改回旧值, 不能在ajax内部调用此代码
            $(this).val(old);
        }
    });
});