layui.use(['table', 'layer'], function() {
    var table = layui.table;
    var layer = layui.layer;
    var $ = layui.jquery;

    // 加载table实例
    table.render({
        elem:"#running_state",  // 绑定到表格ID
        url:"/running/state",   // 获取数据的接口
        page:true,                  // 开启分页
        toolbar:"#toolbar",         // 设置表格工具栏

        // 设置列
        cols:[[
            // 设置表头行号
            {fiele:'id', type:"numbers"},
            // 设置单选框
            //{field:'select', type:"radio"},
            // 设置列
            {field:'protocol', title:'协议', width:80},
            {field:'src-addr', title:'源地址', width:170},
            {field:'src-port', title:'源端口', width:100},
            {field:'dest-addr', title:'目的地址', width:170},
            {field:'dest-port', title:'目的端口', width:100},
            {field:'running-time', title: '运行时间', width:120},
            {field:'forward-flow', title:'转发的流量', width:130},
            {field:'realtime-flow', title:'实时流量', width:120},
            // 设置终止按钮
            {field:'operate', title:'操作' ,toolbar:"#editbar", width:110},
        ]],
    });

    // 监听表格工具栏事件, demo表示待监听容器的lay-filter值
    table.on('toolbar(table)', function(obj) {
        if (obj.event == "LAYTABLE_SEARCH") {
            table.reload('running_state', {
                url:"/running/state_search"
                ,where: {
                      key: $("#key").val()
                }
                ,page: {
                    curr: 1 //重新从第 1 页开始
                }
            });
        }
    });
});