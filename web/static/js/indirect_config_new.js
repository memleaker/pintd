// 加载form 和 element模块
layui.use(['form', 'jquery', 'element', 'layer'], function(){
    var $ = layui.jquery;
    var form = layui.form;
    var layer = layui.layer;
    var element = layui.element;

    // 点击新增配置
    form.on('submit(indirect_cfg_new)', function(data) {
        // 显示加载中
        layer.load();

        // 提交数据到后端
        $.ajax({
            url:"/indirect/cfg_new",  //提交请求的URL
            type:"post",              //请求方式get,post,put,delete等
            data:JSON.stringify(data.field),          //提交的表单数据
            dateType:"json",

            // 调用success回调时，result 代表服务端返回的JSON, msg和success为JSON里的字段
            // success 为返回码200系列时的逻辑
            success:function(result) {
                if (result.success) {
                    layer.closeAll('loading'); // 关闭加载框
                    layer.msg(result.msg, {icon: 6});  //返回数据成功时弹框
                }
            },

            // 调用error 回调时result则不是服务端返回的JSON，而是要更复杂一些
            // error 为无返回或返回码不为200系列以及其它错误逻辑
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

        return false;  //阻止表单跳转。如果需要表单跳转，去掉这段即可。
    });
});