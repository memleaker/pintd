// 加载form 和 element模块
layui.use(['form', 'jquery', 'element', 'layer'], function(){
    var $ = layui.jquery;
    var form = layui.form;
    var layer = layui.layer;
    var element = layui.element;

    // 点击新增配置
    form.on('submit(indirect_cfg_new)', function(data) {
        $.ajax({
            url:"/indirect/cfg_new",  //提交请求的URL
            type:"post",              //请求方式get,post,put,delete等
            data:JSON.stringify(data.field),          //提交的表单数据

            // result 代表服务端返回的JSON, msg和success为JSON里的字段
            success:function(result) {
                if (result.success) {
                    layer.msg(result.msg)  //返回数据成功时弹框
                }
                else {
                    console.log("failed")
                    layer.msg(result.msg) //返回数据失败时弹框
                }
            },

            // 无返回或处理有报错时弹框
            error:function(result){
                alert("服务器无响应!!!");
            }
        });

        return false;  //阻止表单跳转。如果需要表单跳转，去掉这段即可。
    });
});