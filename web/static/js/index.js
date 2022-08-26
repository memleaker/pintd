layui.use('element', function(){
    // 导入jquery element模块
    // tab的切换功能，切换事件监听等，需要依赖element模块
    var $ = layui.jquery
    ,element = layui.element;
    
    // 定义对象active,以及函数tabmange, js可通过active[tabmange]来调用函数
    var active = {
      tabmanage: function() {
        var htmlurl  = $(this).attr('data-url');
        var mytitle  = $(this).attr('mytitle');
        var myid     = $(this).attr('myid');
        var arrayobj = new Array();

        // 先判断是否已经有此tab, 从.layui-tab-title类找创建tab时自动生成的li标签
        // 找到其id放进数组里.
        $(".layui-tab-title").find('li').each(function() {
                var y = $(this).attr("lay-id");
                arrayobj.push(y);
        });

        // 判断用户点击的id是否在数组里
        if (arrayobj.indexOf(myid) >= 0) {
            //标签存在,切换到当前点击的页面
            //changeFrameHeight();
            element.tabChange('demo', myid);
        } else {
            // 标签不存在，创建标签
            // 为让iframe 自适应高度，加了段js
            // 注意, iframe的id必须是不同的，因此我设置成了myid
            element.tabAdd('demo', {
              title:mytitle
              ,content: '<iframe frameborder="no" onload="changeFrameHeight('+myid+')" id = '+myid+' \
                        style="width:100%;hight:100%;" scrolling="no" src='+htmlurl+'></iframe> \
                        <script> \
                        function changeFrameHeight(myid) { \
                          var ifm = document.getElementById(myid); \
                          ifm.height = document.documentElement.clientHeight; \
                        } </script>'
              ,id: myid
            })
            // 切换到当前点击的页面
            element.tabChange('demo', myid);
        }
      }
    };

    // 监听click, 点击后调用函数创建或切换
    $(".leftdaohang").click(function() {
      var type="tabmanage";
      var othis = $(this);
      active[type] ? active[type].call(this, othis) : '';
    });


    // 主动调用click, 让其创建mainpage界面，即主页界面
    // jquery 语法
    $('#mainpage').click();
});
