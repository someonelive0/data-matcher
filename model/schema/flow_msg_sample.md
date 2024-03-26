

# 流转监测消息样本展示

来自 Kafka 的消息样本，这些消息有流转探针turboflow产生，写入kafka

产生Schema可以在线：https://transform.tools/json-to-json-schema

## 1. Message of HTTP sample

```
{
    "timestamp": "2024-02-19T14:40:44.448550+0800",
    "response_time": 21,
    "flow_id": 989719241942234,
    "event_type": "http",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "10.0.0.1",
    "src_port": 58436,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "11.1.1.1",
    "dest_port": 80,
    "proto": "TCP",
    "tx_id": 98,
    "http": {
        "hostname": "11.1.1.1",
        "url": "/login.jsp",
        "http_user_agent": "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)",
        "http_content_type": "text/html",
        "http_method": "POST",
        "protocol": "HTTP/1.1",
        "status": 200,
        "length": 11293,
        "request_headers": [
            {
                "name": "Host",
                "value": "11.1.1.1"
            },
            {
                "name": "Accept-Charset",
                "value": "iso-8859-1,utf-8;q=0.9,*;q=0.1"
            },
            {
                "name": "Accept-Language",
                "value": "en"
            },
            {
                "name": "Content-Type",
                "value": "application/x-www-form-urlencoded"
            },
            {
                "name": "Connection",
                "value": "Keep-Alive"
            },
            {
                "name": "Cookie",
                "value": "JSESSIONID=0000pcLbdFsfiqBR8tOSi927STB:155keeikh"
            },
            {
                "name": "Content-Length",
                "value": "50"
            },
            {
                "name": "User-Agent",
                "value": "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)"
            },
            {
                "name": "Pragma",
                "value": "no-cache"
            },
            {
                "name": "Accept",
                "value": "image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, image/png, */*"
            }
        ],
        "response_headers": [
            {
                "name": "Date",
                "value": "Fri, 18 Oct 2019 09:26:42 GMT"
            },
            {
                "name": "Server",
                "value": "IBM_HTTP_Server"
            },
            {
                "name": "X-UA-Compatible",
                "value": "IE=edge"
            },
            {
                "name": "Content-Length",
                "value": "11293"
            },
            {
                "name": "Keep-Alive",
                "value": "timeout=10, max=2"
            },
            {
                "name": "Connection",
                "value": "Keep-Alive"
            },
            {
                "name": "Content-Type",
                "value": "text/html;charset=UTF-8"
            },
            {
                "name": "Content-Language",
                "value": "zh-CN"
            }
        ],
        "reqHeaders": "Host: 11.1.1.1\r\nAccept-Charset: iso-8859-1,utf-8;q=0.9,*;q=0.1\r\nAccept-Language: en\r\nContent-Type: application/x-www-form-urlencoded\r\nConnection: Keep-Alive\r\nCookie: JSESSIONID=0000pcLbdFsfiqBR8tOSi927STB:155keeikh\r\nContent-Length: 50\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)\r\nPragma: no-cache\r\nAccept: image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, image/png, */*\r\n\r\n",
        "reqBody": "code=&password=&userId=../../../../../../../../etc",
        "reqLen": 50,
        "respHeaders": "Date: Fri, 18 Oct 2019 09:26:42 GMT\r\nServer: IBM_HTTP_Server\r\nX-UA-Compatible: IE=edge\r\nContent-Length: 11293\r\nKeep-Alive: timeout=10, max=2\r\nConnection: Keep-Alive\r\nContent-Type: text/html;charset=UTF-8\r\nContent-Language: zh-CN\r\n\r\n",
        "respBody": "\r\n<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\r\n\r\n\r\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n<head>\r\n<meta name=\"renderer\" content=\"ie-comp\">\r\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\r\n<title>登录页面</title>\r\n<link rel=\"stylesheet\" type=\"text/css\" href=\"SPESTY/images/css.css\">\r\n<script type=\"text/javascript\" src=\"SPESTY/js/jquery.min.js\"></script>\r\n<link rel=\"stylesheet\" href=\"SPESTY/images/lanrenzhijia.css\">\r\n<script type=\"text/javascript\">\r\n$(document).ready(function(){\r\n\t$(\"#dl\").click(function(){\r\n   $(\"#wrapw\").css('background','url(SPESTY/images/tab_bj2.jpg) top no-repeat');\r\n });\r\n <!--$(\"#zc\").click(function(){\r\n   $(\"#wrapw\").css('background','url(SPESTY/images/tab_bj1.jpg) top no-repeat');\r\n });-->\r\n  });\r\n</script>\r\n</head>\r\n\r\n<body onload=\"refreshImg()\">\r\n<form id=\"forms\" method=\"post\">\r\n<div class=\"top\">\r\n  <div class=\"warp\">\r\n    <div class=\"lefts\">上海青年门户欢迎您！</div>\r\n    <div class=\"rights\"><img src=\"SPESTY/images/biao_02.png\" /><img src=\"SPESTY/images/biao_04.png\" /></div>\r\n  </div>\r\n</div>\r\n<div class=\"dh\">\r\n  <div class=\"warps\">\r\n    <div style=\"float:left\"><img src=\"SPESTY/images/login_04.jpg\" /></div>\r\n    <div style=\"float:left; width:750px\">\r\n      <ul>\r\n        <li><a href=\"#\">首页</a></li>\r\n        <li><a href=\"#\">共青团</a></li>\r\n        <li><a href=\"#\">新闻</a></li>\r\n        <li><a href=\"#\">活动</a></li>\r\n        <li><a href=\"#\">培训</a></li>\r\n        <li><a href=\"#\">工作</a></li>\r\n        <li><a href=\"#\">创业</a></li>\r\n        <li><a href=\"#\">婚恋</a></li>\r\n        <li><a href=\"#\">12355</a></li>\r\n        <li><a href=\"#\">志愿时</a></li>\r\n      </ul>\r\n      <ul>\r\n        <li>Home</li>\r\n        <li>Youths</li>\r\n        <li>News</li>\r\n        <li>Activity</li>\r\n        <li>Training</li>\r\n        <li>Job</li>\r\n        <li>Entrepreneur</li>\r\n        <li>Dating</li>\r\n        <li>12355</li>\r\n        <li>Volunteer</li>\r\n      </ul>\r\n    </div>\r\n  </div>\r\n</div>\r\n<div class=\"content\">\r\n  <div class=\"wrapw\" id=\"wrapw\"  style=\"background:url(SPESTY/images/tab_bj2.jpg) top no-repeat; margin-top:-120px;  \">\r\n    <div class=\"tabs\"> <!--<a href=\"#\" hidefocus=\"true\" class=\"active\" id=\"dl\" >个人登录</a> <a href=\"#\" hidefocus=\"true\" id=\"zc\" >--><a>团务平台登录</a> <a>团员登录请用手机</a></div>\r\n    <div class=\"swiper-container\">\r\n      <div class=\"swiper-wrapper\">\r\n        <!-- <div class=\"swiper-slide\">\r\n          <div class=\"content-slide\">\r\n           <div style=\"margin:0 auto;\">\r\n              <div style=\"float:left; width:300px; margin:3px 45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">姓　　名：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9; height:24px\" name=\"username\"  id=\"username\" value=\"\"/>\r\n                </div>\r\n              </div>\r\n              <div style=\"float:left; width:300px; margin:3px  45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">身份证号：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px\" name=\"idcard\"  id=\"idcard\" value=\"\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">手&nbsp;&nbsp;机&nbsp;&nbsp;号：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px\" name=\"mobile\"  id=\"mobile\" value=\"\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">验&nbsp;&nbsp;证&nbsp;&nbsp;码：</div>\r\n                <div style=\"float:left;\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; width:80px\"  id=\"code1\"  name=\"code1\"/>\r\n                  <iframe src =\"image1.jsp\" width=\"130\" height=\"26\" id=\"ifrname\" name=\"ifrname\" frameborder=\"0\" style=\"vertical-align:middle\"  marginwidth=\"0\" marginheight=\"0\" scrolling=\"no\" > </iframe>\r\n                </div>\r\n                <div style=\"float:left; margin-bottom:40px;\">\r\n                  <h5>账号密码及使用咨询请拨打热线电话</h5>\r\n                  <h5>\"咨询热线：400 821 5854 ，热线受理时间：工作日 09:00-18:00\" </h5>\r\n                  <a style=\"float:left; margin-bottom:-10px;margin-top:-10px;\" href=\"/SPESBX/TJ/importTJMember.jsp\" class=\"textunderline\"><font color=\"red\">团员导入模版</font></a>\r\n                </div>\r\n              </div>\r\n\r\n\r\n              <a style=\"width:120px; float:left; margin-top:-30px;margin-left:118px;background:#FF5B66;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\" onclick=\"doLoginTY()\"> 团员登录 </a> <a href=\"/DispatchAction.do?efFormEname=TYRT06&serviceName=TYApplyCYLService&methodName=intiTYRT06\"  style=\"width:120px; float:left; margin-left:20px;margin-top:-30px;background:#1D92BE;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\"> 团员报到 </a> <a  href=\"/DispatchAction.do?efFormEname=TYRT01&serviceName=TYApplyCYLService&methodName=intiTYRT01\" style=\"text-decoration:underline; float:left; margin-left:18px;margin-top:-10px;color:red;font-size:16px;\">我要入团</a>\r\n\r\n            </div>\r\n          </div>\r\n        </div>-->\r\n        <div class=\"swiper-slide\">\r\n          <div class=\"content-slide\">\r\n            <div style=\"margin:0 auto;\">\r\n              <div style=\"float:left; width:300px; margin:5px 45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">用户名：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9; height:24px\" name=\"userId\"  id=\"userId\" value=\"\"/>\r\n                </div>\r\n              </div>\r\n              <div style=\"float:left; width:300px; margin:5px  45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">密&nbsp;&nbsp;&nbsp;码：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; margin-left:2px;\" id=\"password\" name=\"password\" type=\"password\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">验证码：</div>\r\n                <div style=\"float:left;\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; width:80px\" id=\"code\"  name=\"code\"/>\r\n                  <iframe src=\"image.jsp\" id=\"ifrname1\" name=\"ifrname1\" width=\"130\" height=\"26\" frameborder=\"0\" style=\"vertical-align:middle\"  marginwidth=\"0\" marginheight=\"0\" scrolling=\"no\" > </iframe>\r\n                </div>\r\n              </div>\r\n            </div>\r\n            <div style=\"float:left; margin-bottom:6px; margin-left:45px;\">\r\n              <h5>账号密码及使用咨询请拨打热线电话<br/>\r\n              \"咨询热线：400 821 5854 ，热线受理时间：工作日 09:00-18:00\" </h5>\r\n              <a href=\"/SPESBX/TJ/tydrmb.xls\" class=\"textunderline\"><font color=\"red\">团员导入模版下载</font></a>\r\n            </div>\r\n            <a style=\"width:120px; float:left; margin:5px 100px;background:#FF5B66;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\" onclick=\"doLogin()\">登录 </a></div>\r\n        </div>\r\n      </div>\r\n    </div>\r\n  </div>\r\n</div>\r\n\r\n<div style=\"width:320px;margin:0 auto; padding:20px 0;\">\r\n<div style=\"width:230px;float:left;\">\r\n<a target=\"_blank\" href=\"http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=31009102000020\" ><img src=\"SPESTY/images/batb.png\"/>沪公网安备 31009102000020号</a>\r\n</div>\r\n<div style=\"width:80px;float:right;margin-top:-25px;\">\r\n<script type=\"text/javascript\">document.write(unescape(\"%3Cspan id='_ideConac' %3E%3C/span%3E%3Cscript src='http://dcs.conac.cn/js/02/000/0000/60710135/CA020000000607101350008.js' type='text/javascript'%3E%3C/script%3E\"));</script>\r\n</div>\r\n</div>\r\n<div style=\"bottom: 18px; position: fixed; right: 15px; overflow: hidden\">\r\n  <div style=\"font-size: 13px; text-align: center; margin-bottom: 3px\">“青春上海”微信公众号</div>\r\n  <a href=\"#\" target=\"_parent\" style=\"display: block; text-align: center;\"><img style=\"width: 135px; height: 135px\" src=\"SPESTY/images/erweima_gongzhonghao.jpg\"></a> </div>\r\n\r\n\r\n\r\n\r\n\r\n</form>\r\n</body>\r\n</html>\r\n<script src=\"SPESTY/js/jquery-1.10.1.min.js\"></script>\r\n<script src=\"SPESTY/js/idangerous.swiper.min.js\"></script>\r\n\r\n<script>\r\nvar s= \"\";if(s!=null&&s!=\"\")alert(s);\r\n\r\nvar flag= \"\";\r\nif(flag==\"1\")\r\n{\r\n\tdocument.getElementById(\"userId\").value=\"\";\r\n}\r\nif(flag==\"2\")\r\n{\r\n\tdocument.getElementById(\"username\").value=\"\";\r\n\tdocument.getElementById(\"idcard\").value=\"\";\r\n\r\n}\r\n\r\nvar tabsSwiper = new Swiper('.swiper-container',{\r\n\tspeed:500,\r\n\tonSlideChangeStart: function(){\r\n\t\t$(\".tabs .active\").removeClass('active');\r\n\t\t$(\".tabs a\").eq(tabsSwiper.activeIndex).addClass('active');\r\n\t}\r\n});\r\n\r\n$(\".tabs a\").on('touchstart mousedown',function(e){\r\n\te.preventDefault()\r\n\t$(\".tabs .active\").removeClass('active');\r\n\t$(this).addClass('active');\r\n\ttabsSwiper.swipeTo($(this).index());\r\n});\r\n\r\n\r\n\tfunction doLogin(){\r\n\t\tvar userId = document.getElementById(\"userId\").value;\r\n\t\tvar password = document.getElementById(\"password\").value;\r\n\t\tvar code = document.getElementById(\"code\").value;\r\n\t\tif(userId==\"\"||userId==null){\r\n\t\t\talert(\"用户名不能为空！\");\r\n\t\t\tdocument.getElementById(\"userId\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(password==\"\"||password==null){\r\n\t\t\talert(\"密码不能为空！\");\r\n\t\t\tdocument.getElementById(\"password\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(code==\"\"||code==null){\r\n\t\t\talert(\"验证码不能为空！\");\r\n\t\t\tdocument.getElementById(\"code\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tvar action = \"/login.jsp?flag=1\";\r\n\r\n\t\tdocument.all.forms.action = action;\r\n\r\n\t\tdocument.all.forms.submit();\r\n\t}\r\n\tfunction doLoginTY(){\r\n\t\tvar username = document.getElementById(\"username\").value;\r\n\t\tvar idcard = document.getElementById(\"idcard\").value;\r\n\t\t<!--var mobile = document.getElementById(\"mobile\").value;-->\r\n\t\tvar code1 = document.getElementById(\"code1\").value;\r\n\t\tif(username==\"\"||username==null){\r\n\t\t\talert(\"用户名不能为空！\");\r\n\t\t\tdocument.getElementById(\"username\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(idcard==\"\"||idcard==null){\r\n\t\t\talert(\"身份证号不能为空！\");\r\n\t\t\tdocument.getElementById(\"idcard\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\r\n\t\tif(code1==\"\"||code1==null){\r\n\t\t\talert(\"验证码不能为空！\");\r\n\t\t\tdocument.getElementById(\"code1\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tvar action = \"/login.jsp?flag=2\";\r\n\r\n\t\tdocument.all.forms.action = action;\r\n\r\n\t\tdocument.all.forms.submit();\r\n\t}\r\n\tfunction refreshImg(){\r\n\t\tdocument.getElementById(\"ifrname\").src = \"image1.jsp\";\r\n\t\tdocument.getElementById(\"ifrname1\").src = \"image.jsp\";\r\n\t}\r\n\r\n</script>",
        "respLen": 11293,
        "totalLen": 11992
    },
    "host": "172.0.0.1"
}
```

## 2. Message of DNS sample

```
{
    "timestamp": "2024-01-05T14:53:38.929539+0800",
    "response_time": 0,
    "flow_id": 403514986764205,
    "in_iface": "ens37,ens33",
    "event_type": "dns",
    "pcap_filename": "",
    "src_mac": "00:0c:29:0c:f4:d3",
    "src_ip": "192.168.0.1",
    "src_port": 51436,
    "dst_mac": "00:50:56:f9:ad:98",
    "dest_ip": "114.114.114.114",
    "dest_port": 53,
    "proto": "UDP",
    "dns": {
        "version": 2,
        "type": "answer",
        "id": 2301,
        "flags": "8180",
        "qr": true,
        "rd": true,
        "ra": true,
        "rrname": "www.jetbrains.com",
        "rrtype": "A",
        "rcode": "NOERROR",
        "answers": [
            {
                "rrname": "www.jetbrains.com",
                "rrtype": "A",
                "ttl": 31,
                "rdata": "52.84.251.76"
            },
            {
                "rrname": "www.jetbrains.com",
                "rrtype": "A",
                "ttl": 31,
                "rdata": "52.84.251.110"
            },
            {
                "rrname": "www.jetbrains.com",
                "rrtype": "A",
                "ttl": 31,
                "rdata": "52.84.251.12"
            }
        ],
        "grouped": {
            "A": [
                "52.84.251.76",
                "52.84.251.110",
                "52.84.251.12"
            ]
        }
    },
    "host": "172.0.0.1"
}
```

## 3. Message of SSH sample

```
{
    "timestamp": "2024-01-05T14:55:13.040174+0800",
    "response_time": 0,
    "flow_id": 1411221694815052,
    "in_iface": "ens37,ens33",
    "event_type": "ssh",
    "pcap_filename": "",
    "src_mac": "08:26:ae:31:29:1c",
    "src_ip": "172.0.0.1",
    "src_port": 59112,
    "dst_mac": "9c:06:1b:e1:4b:4c",
    "dest_ip": "172.1.1.1",
    "dest_port": 22,
    "proto": "TCP",
    "ssh": {
        "client": {
            "proto_version": "2.0",
            "software_version": "libssh2_1.9.0_DEV"
        },
        "server": {
            "proto_version": "2.0",
            "software_version": "OpenSSH_8.0"
        }
    },
    "host": "172.1.1.1"
}
```

## 4. Message of RDP sample

```
{
    "timestamp": "2024-01-09T16:36:00.698113+0800",
    "response_time": 0,
    "flow_id": 1112759545657919,
    "in_iface": "ens37,ens33",
    "event_type": "rdp",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "11.1.1.1",
    "src_port": 3389,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "10.0.0.1",
    "dest_port": 60876,
    "proto": "TCP",
    "rdp": {
        "tx_id": 0,
        "event_type": "initial_request",
        "cookie": "msrdp_cve-2019-0708.nbin"
    },
    "host": "172.0.0.1"
}
```

## 5. Message of ElsticSearch sample

```
{
    "timestamp": "2024-02-19T14:25:28.382237+0800",
    "response_time": 0,
    "flow_id": 353909404655123,
    "event_type": "elasticsearch",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "10.0.0.1",
    "src_port": 33842,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "11.1.1.1",
    "dest_port": 9200,
    "proto": "TCP",
    "elasticsearch": {
        "username": "",
        "client_terminal": "",
        "database_name": "",
        "table_name": "",
        "sql": "{\r\n \"action\": \"\",\r\n\"template\": \"\",\r\n\"sql\": \"GET /xmldata?item=All \",\r\n\"tables\": []\r\n}",
        "error_msg": "",
        "result_ret": "{ \r\n \"header\": [\"es\"],\r\n\"data\": [ \r\n[\"<!DOCTYPE HTML PUBLIC \\\"-//IETF//DTD HTML 2.0//EN\\\">\\r\\n<html>\\r\\n<head><title>404 Not Found</title></head>\\r\\n<body bgcolor=\\\"white\\\">\\r\\n<h1>404 Not Found</h1>\\r\\n<p>The requested URL was not found on this server.<hr/>Powered by Tengine</body>\\r\\n</html>\\r\\n\"]]\r\n}",
        "ret_count": 1,
        "exec_timer": 45383,
        "ret_len": 304,
        "error_code": 404,
        "totalLen": 389
    },
    "host": "172.0.0.1"
}
```

## 6. Message of POP3 sample

```
{
    "timestamp": "2024-02-19T14:25:58.260157+0800",
    "response_time": 1,
    "flow_id": 1505334304003198,
    "event_type": "pop3",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "192.168.0.1",
    "src_port": 62026,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "11.1.1.1",
    "dest_port": 110,
    "proto": "TCP",
    "tx_id": 0,
    "pop3": {
        "username": "guan1212@idss-cn.com",
        "command": "RETR 10277",
        "status": "PARSE_DONE",
        "From": "<qian1212@idss-cn.com>",
        "To": [
            "<guan1212@idss-cn.com>",
            "<zhuan1212@idss-cn.com>",
            "<zhang1212@idss-cn.com>",
            "<zhu1212@idss-cn.com>"
        ],
        "Subject": "回复: ****收款",
        "Date": "Fri, 18 Oct 2019 17:16:38 +0800",
        "attachment": [
            {
                "filename": "InsertPic_.jpg",
                "file_256": "c6780b2afa638e78e707c73f0e469d8764d00b7953f9ad0a0a63bc76fe7db9f8"
            },
            {
                "file_256": "ef5e114a579a807d8632c50f8bdf765db09d4e3fa2914372ed8cc1ee2c8cada1"
            },
            {
                "filename": "CatchB1A5(10-18-17-13-43).jpger",
                "file_256": "a87e997f6f726e08a7c0c2986123e54575e9b278a8ff265591c82e396219feb6"
            },
            {
                "file_256": "a10fa55b593c44177455252774dd030198f5f97c7521b1e0bac2be140b7465e0"
            }
        ],
        "email_body": "各位好,\r\n\r\n苏州****的到款是\r\n****万是IDSS2019JC****(****)全款\r\n****万是IDSS2018JC****(****)20%\r\n\r\n\r\n\r\n钱****          销售部门     \r\n \r\n\r\nInformation & Data Security Solutions Co.,Ltd. \r\n上海总部地址：上海市普陀区大渡河路388弄5号华宏商务中心6层\r\n北京总部地址：北京市海淀区东北旺西路8号中关村软件园8号楼一层131室\r\n手机：150********\r\n上海总部电话：(021)62090100\r\nWEB：www.idss-cn.com  \r\n \r\n发件人： guan****@idss-cn.com\r\n发送时间： 2019-10-18 15:02\r\n收件人： qian****; 庄****; 张****; 朱****\r\n主题： ****收款\r\n关于****两笔收款，因为2个项目开票****万金额一致，请****去找客户确认下，回复庄****，谢谢~~\r\n\r\n\r\n\r\n\r\n\r\n\r\n财务部            管****\r\n     \r\n上海观安信息技术股份有限公司\r\nInformation & Data Security Solutions Co.,Ltd. \r\n上海总部地址：上海市普陀区大渡河路388弄5号华宏商务中心6层\r\n北京总部地址：北京市海淀区东北旺西路8号中关村软件园8号楼一层131室\r\n手机：186********\r\n上海总部电话：(021)62090100\r\nWEB：www.idss-cn.com  \r\n",
        "totalLen": 1264
    },
    "host": "172.0.0.1"
}
```

## 7. Message of SMTP sample

```
{
    "timestamp": "2024-02-19T14:40:41.013121+0800",
    "response_time": 15,
    "flow_id": 2077527083938102,
    "event_type": "smtp",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "192.168.0.1",
    "src_port": 62454,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "11.1.1.1",
    "dest_port": 25,
    "proto": "TCP",
    "tx_id": 0,
    "smtp": {
        "helo": "DESKTOP56EGES5",
        "mail_from": "<x1212@idss-cn.com>",
        "From": "<x1212@idss-cn.com>",
        "rcpt_to": [
            "<h1212@idss-cn.com>",
            "<wang1212@idss-cn.com>",
            "<zhangz1212@caict.ac.cn>"
        ],
        "To": [
            "<h1212@idss-cn.com>",
            "<wang1212@idss-cn.com>",
            "<zhang1212@caict.ac.cn>"
        ],
        "status": "BODY_STARTED",
        "cc": [
            "=?utf-8?B?J+eOi+i1nCc=?= <wang1212@idss-cn.com>",
            "=?utf-8?B?J+iDoee7jeWLhyc=?= <h1212@idss-cn.com>"
        ],
        "Subject": "中国****研究院网络安全技术本地化支撑服务（上海）服务合同-验收材料",
        "Date": "Fri, 18 Oct 2019 17:26:16 +0800",
        "email_body": "张****\r\n\r\n \r\n\r\n附件是网络安全技术本地化支撑服务（上海）服务完整验收材料。\r\n\r\n \r\n\r\n《验收单》\r\n\r\n《技术服务质量考核表》\r\n\r\n《风险评估报告》\r\n\r\n《渗透测试报告》\r\n\r\n《重要时期安全保障服务报告》\r\n\r\n《应急响应报告》\r\n\r\n \r\n\r\n验收技术接口人：夏**** 133******** \r\n\r\n \r\n\r\n \r\n\r\n-夏****\r\n\r\n \r\n\r\n \r\n\r\n发件人: 夏**** <x****@idss-cn.com> \r\n发送时间: 2019年10月8日 13:49\r\n收件人: zhang**** <zhang****@caict.ac.cn>; 夏**** <x****@idss-cn.com>\r\n抄送: 王**** <wang****@idss-cn.com>\r\n主题: Fw:中国****研究院网络安全技术本地化支撑服务（上海）服务合同-****W已签\r\n\r\n \r\n\r\n张****,\r\n\r\n \r\n\r\n第一阶段项目款已确认收到。\r\n\r\n附件是对应的验收单，麻烦帮盖章扫描完成财务流程。多谢！\r\n\r\n \r\n\r\n-夏****\r\n\r\n \r\n\r\n \r\n\r\n------------------ Original ------------------\r\n\r\nFrom:  \"夏****\"<x****@idss-cn.com <mailto:x****@idss-cn.com> >;\r\n\r\nDate:  Tue, Sep 24, 2019 07:40 PM\r\n\r\nTo:  \"zhang****\"<zhang****@caict.ac.cn <mailto:zhang****@caict.ac.cn> >; \r\n\r\nSubject:  中国****研究院网络安全技术本地化支撑服务（上海）服务合同-****W已签\r\n\r\n \r\n\r\n张****您好，\r\n\r\n \r\n\r\n中国****研究院网络安全技术本地化支撑服务（上海）服务合同申请第一笔付款。\r\n\r\n \r\n\r\n合同中约定的付款条件如下：\r\n\r\n2、付款措施：\r\n\r\n1）合同签订后，乙方提供给甲方项目实施进度表和等额增值税专用发票后5个工作日内，甲方向乙方支付60%的合同费用，共计：¥****.00元（大写：人民币****元整）。\r\n\r\n \r\n\r\n \r\n\r\n您看项目实施进度表若没问题，我这边申请开票了。\r\n\r\n \r\n\r\n-夏****\r\n\r\n",
        "totalLen": 2027,
        "x_mailer": "Microsoft Outlook 16.0",
        "subject_md5": "6f3f0d47a340d2faa3df20b12c2a912c"
    },
    "host": "172.0.0.1"
}
```
