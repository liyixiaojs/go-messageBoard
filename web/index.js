function ajax(options) {
  options=options||{};
  options.type=(options.type||'GET').toUpperCase();
  options.dataType=options.dataType||'json';
  params=formatParams(options.data);

  //创建-第一步
  var xhr;
  //非IE6
  if(window.XMLHttpRequest){
      xhr=new XMLHttpRequest();
  }

  //接收-第三步
  xhr.onreadystatechange=function(){
      if(xhr.readyState==4){
          var status=xhr.status;
          if(status>=200&&status<300){
              options.success&&options.success(xhr.responseText,xhr.responseXML);
          }else{
              options.error&&options.error(status);
          }
      }
  }

  //连接和发送-第二步
  if(options.type=='GET'){
      xhr.open('GET',options.url+'?'+params,true);
      xhr.send(null);
  }else if(options.type=='POST'){
      xhr.open('POST',options.url,true);
      //设置表单提交时的内容类型
      // xhr.setRequestHeader("Content-Type", "application/json;charset=utf-8");
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(JSON.stringify(params));
  }
}

//格式化参数
function formatParams(data){
  // var arr=[];
  // for(var name in data){
  //     arr.push(encodeURIComponent(name)+'='+encodeURIComponent(data[name]));
  // }
  // arr.push(('v=' + Math.random()).replace('.',''));
  // return arr.join('&');
  return data
}

function deleteItem(id) {
  ajax({
    url:'/delete',
    type:'POST',
    dataType:'json',
    data: {id: id, name: userList.value},
    success:function(result){
      loadData()
    },
    error:function(status){
      console.error('获取列表失败')
    }
  })
}

function contentTemplate(id, name, content, time) {
  time = new Date(time).format("yyyy-MM-dd hh:mm:ss")
  if (userList.value === name) {
    return `
      <div class="item owner">
      <div class="user-info">
        <a href="javascript:;" onclick="deleteItem('${id}')">删除</a>
        <div class="time">
          ${time}
        </div>
        <div class="name">
          ${name}
        </div>
      </div>
      <div class="item-content">
        ${content}
      </div>
    </div>
    `
  }
  return `
    <div class="item">
      <div class="user-info">
        <div class="name">
          ${name}
        </div>
        <div class="time">
          ${time}
        </div>
      </div>
      <div class="item-content">
        ${content}
      </div>
    </div>
  `
}

function load(data) {
  let html = ''
  data.forEach(item => {
    html += contentTemplate(item.id, item.name, item.content, item.time) + '\n'
  })
  content.innerHTML = html
}

function loadData() {
  ajax({
    url:'/queryList',
    type:'POST',
    dataType:'json',
    data: {},
    success:function(result){
      load(JSON.parse(result).data || [])
    },
    error:function(status){
      console.error('获取列表失败')
    }
  })

  // ajax({
  //   url:'/error',
  //   type:'POST',
  //   dataType:'json',
  //   data: {},
  //   success:function(result){
  //     load(JSON.parse(result).data || [])
  //   },
  //   error:function(status){
  //     console.error('获取列表失败')
  //   }
  // })
};

(() => {
  Date.prototype.format = function(fmt) { 
    var o = { 
      "M+" : this.getMonth()+1,                 //月份 
      "d+" : this.getDate(),                    //日 
      "h+" : this.getHours(),                   //小时 
      "m+" : this.getMinutes(),                 //分 
      "s+" : this.getSeconds(),                 //秒 
      "q+" : Math.floor((this.getMonth()+3)/3), //季度 
      "S"  : this.getMilliseconds()             //毫秒 
    }; 
    if(/(y+)/.test(fmt)) {
            fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length)); 
    }
    for(var k in o) {
      if(new RegExp("("+ k +")").test(fmt)){
        fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
      }
    }
    return fmt; 
  }
})();

(() => {
  let button = window.document.querySelector('#submitButton')
  button.addEventListener('click', e => {
     ajax({
      url:'/save',
      type:'POST',
      dataType:'json',
      data: {content:contentText.value, name: userList.value},
      success:function(result){
        if (JSON.parse(result).success) {
          loadData()
        }
        console.log(JSON.parse(result).success)
      },
      error:function(status){
        console.error('保存失败')
      }
    });
  })
  userList.addEventListener('change', e => {
    loadData()
  })
  loadData()
})();