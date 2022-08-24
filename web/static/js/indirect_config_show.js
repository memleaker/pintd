layui.use(['laypage', 'table'], function() {
    var laypage = layui.laypage;
    var table = layui.table;

    // 加载laypage实例
    laypage.render({
        elem:"page",  // 绑定到容器的ID
        count:100,    // 条目数量，从服务器获取
        limit:15,     // 每页显示数量
        groups:5,     // 一次显示的页码
    });

    // 加载table实例
    table.render({
        elem:"#indirect_cfg_show",  // 绑定到表格ID
        url:"/indirect/cfg_show",   // 获取数据的接口
        page:true,                  // 开启分页

        // 设置列
        cols:[[
            {field:'protocol', title:'协议', width:60},
            {field:'listen-addr', title:'监听地址', width:120},
            {field:'listen-port', title:'监听端口', width:90},
            {field:'dest-addr', title:'目的地址', width:120},
            {field:'dest-port', title:'目的端口', width:90},
            {field:'acl', title:'访问控制', width:90},
            {field:'admit-addr', title:'白名单', width:120},
            {field:'deny-addr', title:'黑名单', width:120},
            {field:'max-conns', title:'最大连接数', width:100},
            {field:'memo', title:'备注'},
        ]],
    });
});