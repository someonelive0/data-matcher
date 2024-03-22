package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

func TestFlowHttp(t *testing.T) {
	// s := `{}`

	flowHttp := &FlowHttp{}
	err := json.Unmarshal([]byte(s), flowHttp)
	if err != nil {
		t.Fatalf("TestMstHttp failed %s", err)
	}
	t.Logf("%d, %#v", len(s), *flowHttp)
}

func TestFlowHttpSchema1(t *testing.T) {
	b, err := os.ReadFile("schema/flow_http_schema.json")
	if err != nil {
		panic(err.Error())
	}
	schemaLoader := gojsonschema.NewBytesLoader(b)
	documentLoader := gojsonschema.NewStringLoader(s)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		t.Fail()
	}
}

func TestFlowHttpSchema2(t *testing.T) {
	sch, err := jsonschema.Compile("schema/flow_http_schema.json")
	if err != nil {
		log.Fatalf("%#v", err)
	}

	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		fmt.Printf("The document is not valid. see errors :\n")
		log.Fatal(err)
	}

	if err = sch.ValidateInterface(v); err != nil {
		fmt.Printf("The document is not valid. see errors :\n")
		log.Printf("%#v", err)
		schemaerr := err.(*jsonschema.ValidationError)
		fmt.Printf(" invalid Message %s\n\n", schemaerr.Message)

		for i, e := range schemaerr.Causes { // e is *jsonschema.ValidationError
			fmt.Printf(" %d - %#v\n", i, e.Message)
		}
		t.Fail()
	}
}

const s string = `
{
    "timestamp": "2024-02-19T14:40:44.448550+0800",
    "response_time": 21,
    "flow_id": 989719241942234,
    "event_type": "http",
    "pcap_filename": "",
    "src_mac": "e8:65:49:4b:46:55",
    "src_ip": "10.10.20.113",
    "src_port": 58436,
    "dst_mac": "00:27:0d:ff:5b:e1",
    "dest_ip": "101.231.211.176",
    "dest_port": 80,
    "proto": "TCP",
    "tx_id": 98,
    "http": {
        "hostname": "101.231.211.176",
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
                "value": "101.231.211.176"
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
        "reqHeaders": "Host: 101.231.211.176\r\nAccept-Charset: iso-8859-1,utf-8;q=0.9,*;q=0.1\r\nAccept-Language: en\r\nContent-Type: application/x-www-form-urlencoded\r\nConnection: Keep-Alive\r\nCookie: JSESSIONID=0000pcLbdFsfiqBR8tOSi927STB:155keeikh\r\nContent-Length: 50\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)\r\nPragma: no-cache\r\nAccept: image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, image/png, */*\r\n\r\n",
        "reqBody": "code=&password=&userId=../../../../../../../../etc",
        "reqLen": 50,
        "respHeaders": "Date: Fri, 18 Oct 2019 09:26:42 GMT\r\nServer: IBM_HTTP_Server\r\nX-UA-Compatible: IE=edge\r\nContent-Length: 11293\r\nKeep-Alive: timeout=10, max=2\r\nConnection: Keep-Alive\r\nContent-Type: text/html;charset=UTF-8\r\nContent-Language: zh-CN\r\n\r\n",
        "respBody": "\r\n<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\r\n\r\n\r\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n<head>\r\n<meta name=\"renderer\" content=\"ie-comp\">\r\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\r\n<title>登录页面</title>\r\n<link rel=\"stylesheet\" type=\"text/css\" href=\"SPESTY/images/css.css\">\r\n<script type=\"text/javascript\" src=\"SPESTY/js/jquery.min.js\"></script>\r\n<link rel=\"stylesheet\" href=\"SPESTY/images/lanrenzhijia.css\">\r\n<script type=\"text/javascript\">\r\n$(document).ready(function(){\r\n\t$(\"#dl\").click(function(){\r\n   $(\"#wrapw\").css('background','url(SPESTY/images/tab_bj2.jpg) top no-repeat');\r\n });\r\n <!--$(\"#zc\").click(function(){\r\n   $(\"#wrapw\").css('background','url(SPESTY/images/tab_bj1.jpg) top no-repeat');\r\n });-->\r\n  });\r\n</script>\r\n</head>\r\n\r\n<body onload=\"refreshImg()\">\r\n<form id=\"forms\" method=\"post\">\r\n<div class=\"top\">\r\n  <div class=\"warp\">\r\n    <div class=\"lefts\">上海青年门户欢迎您！</div>\r\n    <div class=\"rights\"><img src=\"SPESTY/images/biao_02.png\" /><img src=\"SPESTY/images/biao_04.png\" /></div>\r\n  </div>\r\n</div>\r\n<div class=\"dh\">\r\n  <div class=\"warps\">\r\n    <div style=\"float:left\"><img src=\"SPESTY/images/login_04.jpg\" /></div>\r\n    <div style=\"float:left; width:750px\">\r\n      <ul>\r\n        <li><a href=\"#\">首页</a></li>\r\n        <li><a href=\"#\">共青团</a></li>\r\n        <li><a href=\"#\">新闻</a></li>\r\n        <li><a href=\"#\">活动</a></li>\r\n        <li><a href=\"#\">培训</a></li>\r\n        <li><a href=\"#\">工作</a></li>\r\n        <li><a href=\"#\">创业</a></li>\r\n        <li><a href=\"#\">婚恋</a></li>\r\n        <li><a href=\"#\">12355</a></li>\r\n        <li><a href=\"#\">志愿时</a></li>\r\n      </ul>\r\n      <ul>\r\n        <li>Home</li>\r\n        <li>Youths</li>\r\n        <li>News</li>\r\n        <li>Activity</li>\r\n        <li>Training</li>\r\n        <li>Job</li>\r\n        <li>Entrepreneur</li>\r\n        <li>Dating</li>\r\n        <li>12355</li>\r\n        <li>Volunteer</li>\r\n      </ul>\r\n    </div>\r\n  </div>\r\n</div>\r\n<div class=\"content\">\r\n  <div class=\"wrapw\" id=\"wrapw\"  style=\"background:url(SPESTY/images/tab_bj2.jpg) top no-repeat; margin-top:-120px;  \">\r\n    <div class=\"tabs\"> <!--<a href=\"#\" hidefocus=\"true\" class=\"active\" id=\"dl\" >个人登录</a> <a href=\"#\" hidefocus=\"true\" id=\"zc\" >--><a>团务平台登录</a> <a>团员登录请用手机</a></div>\r\n    <div class=\"swiper-container\">\r\n      <div class=\"swiper-wrapper\">\r\n        <!-- <div class=\"swiper-slide\">\r\n          <div class=\"content-slide\">\r\n           <div style=\"margin:0 auto;\">\r\n              <div style=\"float:left; width:300px; margin:3px 45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">姓　　名：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9; height:24px\" name=\"username\"  id=\"username\" value=\"\"/>\r\n                </div>\r\n              </div>\r\n              <div style=\"float:left; width:300px; margin:3px  45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">身份证号：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px\" name=\"idcard\"  id=\"idcard\" value=\"\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">手&nbsp;&nbsp;机&nbsp;&nbsp;号：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px\" name=\"mobile\"  id=\"mobile\" value=\"\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">验&nbsp;&nbsp;证&nbsp;&nbsp;码：</div>\r\n                <div style=\"float:left;\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; width:80px\"  id=\"code1\"  name=\"code1\"/>\r\n                  <iframe src =\"image1.jsp\" width=\"130\" height=\"26\" id=\"ifrname\" name=\"ifrname\" frameborder=\"0\" style=\"vertical-align:middle\"  marginwidth=\"0\" marginheight=\"0\" scrolling=\"no\" > </iframe>\r\n                </div>\r\n                <div style=\"float:left; margin-bottom:40px;\">\r\n                  <h5>账号密码及使用咨询请拨打热线电话</h5>\r\n                  <h5>\"咨询热线：400 821 5854 ，热线受理时间：工作日 09:00-18:00\" </h5>\r\n                  <a style=\"float:left; margin-bottom:-10px;margin-top:-10px;\" href=\"/SPESBX/TJ/importTJMember.jsp\" class=\"textunderline\"><font color=\"red\">团员导入模版</font></a>\r\n                </div>\r\n              </div>\r\n\r\n\r\n              <a style=\"width:120px; float:left; margin-top:-30px;margin-left:118px;background:#FF5B66;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\" onclick=\"doLoginTY()\"> 团员登录 </a> <a href=\"/DispatchAction.do?efFormEname=TYRT06&serviceName=TYApplyCYLService&methodName=intiTYRT06\"  style=\"width:120px; float:left; margin-left:20px;margin-top:-30px;background:#1D92BE;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\"> 团员报到 </a> <a  href=\"/DispatchAction.do?efFormEname=TYRT01&serviceName=TYApplyCYLService&methodName=intiTYRT01\" style=\"text-decoration:underline; float:left; margin-left:18px;margin-top:-10px;color:red;font-size:16px;\">我要入团</a>\r\n\r\n            </div>\r\n          </div>\r\n        </div>-->\r\n        <div class=\"swiper-slide\">\r\n          <div class=\"content-slide\">\r\n            <div style=\"margin:0 auto;\">\r\n              <div style=\"float:left; width:300px; margin:5px 45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">用户名：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9; height:24px\" name=\"userId\"  id=\"userId\" value=\"\"/>\r\n                </div>\r\n              </div>\r\n              <div style=\"float:left; width:300px; margin:5px  45px;\">\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">密&nbsp;&nbsp;&nbsp;码：</div>\r\n                <div style=\"float:left\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; margin-left:2px;\" id=\"password\" name=\"password\" type=\"password\"/>\r\n                </div>\r\n                <br />\r\n                <br />\r\n                <div style=\"float:left; font-size:14px; font-family:'微软雅黑'; line-height:24px\">验证码：</div>\r\n                <div style=\"float:left;\">\r\n                  <input style=\"border:1px solid #c8c9c9;height:24px; width:80px\" id=\"code\"  name=\"code\"/>\r\n                  <iframe src=\"image.jsp\" id=\"ifrname1\" name=\"ifrname1\" width=\"130\" height=\"26\" frameborder=\"0\" style=\"vertical-align:middle\"  marginwidth=\"0\" marginheight=\"0\" scrolling=\"no\" > </iframe>\r\n                </div>\r\n              </div>\r\n            </div>\r\n            <div style=\"float:left; margin-bottom:6px; margin-left:45px;\">\r\n              <h5>账号密码及使用咨询请拨打热线电话<br/>\r\n              \"咨询热线：400 821 5854 ，热线受理时间：工作日 09:00-18:00\" </h5>\r\n              <a href=\"/SPESBX/TJ/tydrmb.xls\" class=\"textunderline\"><font color=\"red\">团员导入模版下载</font></a>\r\n            </div>\r\n            <a style=\"width:120px; float:left; margin:5px 100px;background:#FF5B66;height:50px;display:block;color:#fff;font-size:20px;line-height:50px;text-align:center;cursor:pointer;\" onclick=\"doLogin()\">登录 </a></div>\r\n        </div>\r\n      </div>\r\n    </div>\r\n  </div>\r\n</div>\r\n\r\n<div style=\"width:320px;margin:0 auto; padding:20px 0;\">\r\n<div style=\"width:230px;float:left;\">\r\n<a target=\"_blank\" href=\"http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=31009102000020\" ><img src=\"SPESTY/images/batb.png\"/>沪公网安备 31009102000020号</a>\r\n</div>\r\n<div style=\"width:80px;float:right;margin-top:-25px;\">\r\n<script type=\"text/javascript\">document.write(unescape(\"%3Cspan id='_ideConac' %3E%3C/span%3E%3Cscript src='http://dcs.conac.cn/js/02/000/0000/60710135/CA020000000607101350008.js' type='text/javascript'%3E%3C/script%3E\"));</script>\r\n</div>\r\n</div>\r\n<div style=\"bottom: 18px; position: fixed; right: 15px; overflow: hidden\">\r\n  <div style=\"font-size: 13px; text-align: center; margin-bottom: 3px\">“青春上海”微信公众号</div>\r\n  <a href=\"#\" target=\"_parent\" style=\"display: block; text-align: center;\"><img style=\"width: 135px; height: 135px\" src=\"SPESTY/images/erweima_gongzhonghao.jpg\"></a> </div>\r\n\r\n\r\n\r\n\r\n\r\n</form>\r\n</body>\r\n</html>\r\n<script src=\"SPESTY/js/jquery-1.10.1.min.js\"></script>\r\n<script src=\"SPESTY/js/idangerous.swiper.min.js\"></script>\r\n\r\n<script>\r\nvar s= \"\";if(s!=null&&s!=\"\")alert(s);\r\n\r\nvar flag= \"\";\r\nif(flag==\"1\")\r\n{\r\n\tdocument.getElementById(\"userId\").value=\"\";\r\n}\r\nif(flag==\"2\")\r\n{\r\n\tdocument.getElementById(\"username\").value=\"\";\r\n\tdocument.getElementById(\"idcard\").value=\"\";\r\n\r\n}\r\n\r\nvar tabsSwiper = new Swiper('.swiper-container',{\r\n\tspeed:500,\r\n\tonSlideChangeStart: function(){\r\n\t\t$(\".tabs .active\").removeClass('active');\r\n\t\t$(\".tabs a\").eq(tabsSwiper.activeIndex).addClass('active');\r\n\t}\r\n});\r\n\r\n$(\".tabs a\").on('touchstart mousedown',function(e){\r\n\te.preventDefault()\r\n\t$(\".tabs .active\").removeClass('active');\r\n\t$(this).addClass('active');\r\n\ttabsSwiper.swipeTo($(this).index());\r\n});\r\n\r\n\r\n\tfunction doLogin(){\r\n\t\tvar userId = document.getElementById(\"userId\").value;\r\n\t\tvar password = document.getElementById(\"password\").value;\r\n\t\tvar code = document.getElementById(\"code\").value;\r\n\t\tif(userId==\"\"||userId==null){\r\n\t\t\talert(\"用户名不能为空！\");\r\n\t\t\tdocument.getElementById(\"userId\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(password==\"\"||password==null){\r\n\t\t\talert(\"密码不能为空！\");\r\n\t\t\tdocument.getElementById(\"password\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(code==\"\"||code==null){\r\n\t\t\talert(\"验证码不能为空！\");\r\n\t\t\tdocument.getElementById(\"code\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tvar action = \"/login.jsp?flag=1\";\r\n\r\n\t\tdocument.all.forms.action = action;\r\n\r\n\t\tdocument.all.forms.submit();\r\n\t}\r\n\tfunction doLoginTY(){\r\n\t\tvar username = document.getElementById(\"username\").value;\r\n\t\tvar idcard = document.getElementById(\"idcard\").value;\r\n\t\t<!--var mobile = document.getElementById(\"mobile\").value;-->\r\n\t\tvar code1 = document.getElementById(\"code1\").value;\r\n\t\tif(username==\"\"||username==null){\r\n\t\t\talert(\"用户名不能为空！\");\r\n\t\t\tdocument.getElementById(\"username\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tif(idcard==\"\"||idcard==null){\r\n\t\t\talert(\"身份证号不能为空！\");\r\n\t\t\tdocument.getElementById(\"idcard\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\r\n\t\tif(code1==\"\"||code1==null){\r\n\t\t\talert(\"验证码不能为空！\");\r\n\t\t\tdocument.getElementById(\"code1\").focus();\r\n\t\t\treturn;\r\n\t\t}\r\n\t\tvar action = \"/login.jsp?flag=2\";\r\n\r\n\t\tdocument.all.forms.action = action;\r\n\r\n\t\tdocument.all.forms.submit();\r\n\t}\r\n\tfunction refreshImg(){\r\n\t\tdocument.getElementById(\"ifrname\").src = \"image1.jsp\";\r\n\t\tdocument.getElementById(\"ifrname1\").src = \"image.jsp\";\r\n\t}\r\n\r\n</script>",
        "respLen": 11293,
        "totalLen": 11992
    },
    "host": "172.16.0.230"
}
`
