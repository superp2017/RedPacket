package activity

import (

	// "fmt"
	// "github.com/astaxie/beego"

	"bytes"

	// "github.com/astaxie/beego"

	// "html"
	"JsLib/JsNet"
	// "bufio"
	"html/template"
)

const tpl = `<!DOCTYPE html>
<html>

	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		<title>{{.Title}}</title>
		<meta charset="utf-8">
		<meta name="apple-touch-fullscreen" content="YES">
		<meta name="format-detection" content="telephone=no">
		<meta name="apple-mobile-web-app-capable" content="yes">
		<meta name="apple-mobile-web-app-status-bar-style" content="black">
		<meta http-equiv="Expires" content="-1">
		<meta http-equiv="pragram" content="no-cache">
		<meta name="viewport" content="width=640, user-scalable=no, target-densitydpi=device-dpi">
		<link rel="stylesheet" type="text/css" href="http://activity.junsie.com/static/css/main.css">
		<link rel="stylesheet" type="text/css" href="http://activity.junsie.com/static/css/myself.css">
        <link rel="stylesheet" type="text/css" href="http://activity.junsie.com/static/css/endpic.css">

        </head>

		<body class="s-bg-ddd pc no-3d" style="-webkit-user-select: none;">
		<section class="u-audio" data-src="http://activity.junsie.com/static/images/music.mp3">
			<p id="coffee_flow" class="btn_audio">
				<strong class="txt_audio z-hide">关闭</strong>
				<span class="css_sprite01 audio_open"></span>

				<div id="song_img" class="" style="display: none;">打开</div>
				<div class="coffee-steam-box" style="height: 100px; width: 44px; left:80px; top: -10px; position: absolute; overflow: hidden; z-index: 0;"></div>
			</p>
		</section>

		<section class="u-arrow">
			<p class="css_sprite01"></p>
		</section>

            <section class="p-ct transformNode-2d" style="height: 1008px;">
                <div class="translate-back" style="height: 1008px;">

                <div class="m-page m-fengye " data-page-type="bigTxt" data-statics="info_list" style="height:70%;">
					<div class="page-con lazy-finish" data-position="50% 50%" data-size="cover" style="background:url(http://activity.junsie.com/static/images/01-01.png); background-size:cover; background-position: 50% 50%;">
						<img src="http://activity.junsie.com/static/images/01-02.png"width="100%" right="100%" class='a-fadeinT page3_1' style='position:absolute; ' >
						<img src="http://activity.junsie.com/static/images/01-03.png"width="100%" right="100%"  class='a-fadeinB page3_3' style='position:absolute; ' >
						<img src="http://activity.junsie.com/static/images/01-04.png"width="100%" right="100%"  class='a-fadeinT page1_3' style='position:absolute;' >
						<img src="http://activity.junsie.com/static/images/01-05.png"width="100%" right="100%"  class='a-fadeinL page1_4' style='position:absolute; top:4.5%' >
						<img src="http://activity.junsie.com/static/images/01-06.png"width="100%" right="100%" class='a-fadeinR page1_4' style='position:absolute;  top:4.5%  ' >	
						<img src="http://activity.junsie.com/static/images/01-07.png"width="100%" right="100%"  class='a-fadeinL page1_5' style='position:absolute; top:4.5%' >
						<img src="http://activity.junsie.com/static/images/01-08.png"width="100%" right="100%" class='a-fadeinR page1_5' style='position:absolute; top:5.5%' >	
						<img src="http://activity.junsie.com/static/images/01-09.png"width="100%" right="100%"  class='a-fadeinL page1_5' style='position:absolute; top:5.5%' >
						<img src="http://activity.junsie.com/static/images/01-10.png"width="100%" right="100%" class='a-fadeinT page3_3' style='position:absolute; top:6.5%' >						
					</div>
				</div>

				<div class="m-page m-bigTxt f-hide" data-page-type="bigTxt" data-statics="info_list" style="height:70%;">
					<div class="page-con lazy-finish" data-position="50% 50%" data-size="cover" style="background:url(http://activity.junsie.com/static/images/02-01.png); background-size: cover; height:1280px; background-position: 50% 50%;">
						<img src="http://activity.junsie.com/static/images/02-02.png" width="100%" right="100%" class='a-rotateinLT page3_1' style='position:absolute;'>
						<img src="http://activity.junsie.com/static/images/02-03.png" width="100%" right="100%" class='a-fadeinB page1_2' style='position:absolute;  '>
						<img src="http://activity.junsie.com/static/images/02-04.png" width="100%" right="100%" class='a-fadeinT page1_2' style='position:absolute;  '>
						<img src="http://activity.junsie.com/static/images/02-06.png" width="100%" right="100%" class='a-bounceinT page1_6' style='position:absolute;  '>
						<img src="http://activity.junsie.com/static/images/02-07.png" width="100%" right="100%" class='a-fadeinB page3_2' style='position:absolute;  '>
						<img src="http://activity.junsie.com/static/images/02-08.png" width="100%" right="100%" class='a-rotateinLT page1_4' style='position:absolute;  '>
						<img src="http://activity.junsie.com/static/images/02-05.png" width="100%" right="100%" class='a-fadeinB page1_4' style='position:absolute;  '>
					</div>
				</div>


				<div class="m-page m-bigTxt f-hide" data-page-type="bigTxt" data-statics="info_list" style="height:70%;">
					<div class="page-con j-txtWrap lazy-finish" data-position="50% 50%" data-size="cover" style="background:url(http://activity.junsie.com/static/images/03-01.png); background-size:cover; background-position: 50% 50%;">
						<img src="http://activity.junsie.com/static/images/03-02.png"width="100%" right="100%" class='a-fadeinT page1_2' style='position:absolute;'>
						<img src="http://activity.junsie.com/static/images/03-03.png"width="100%" right="100%" class='a-fadeinB page1_3' style='position:absolute; '>
						<img src="http://activity.junsie.com/static/images/03-04.png"width="100%" right="100%" class='a-fadeinB page1_4' style='position:absolute; '>
						<img src="http://activity.junsie.com/static/images/03-05.png"width="100%" right="100%" class='a-fadeinR page1_5' style='position:absolute; '>
						<img src="http://activity.junsie.com/static/images/03-06.png"width="100%" right="100%" class='a-fadeinL page1_5' style='position:absolute; '>
						<img src="http://activity.junsie.com/static/images/03-07.png"width="100%" right="100%" class='a-fadeinT page1_5' style='position:absolute; '>
					</div>
				</div>

                </div>
            </section>

			
			{{.Debug}}
            <script src="http://activity.junsie.com/static/js/offline.js" type="text/javascript" ></script>
            <script src="http://activity.junsie.com/static/js/zepto.js" type="text/javascript" charset="utf-8"></script>
            <script src="http://activity.junsie.com/static/js/coffee.js" type="text/javascript" charset="utf-8"></script>
            <script src="http://activity.junsie.com/static/js/zj_main.js" type="text/javascript" charset="utf-8"></script>

        </body>

</html>`

type Para struct {
	Title string
	Debug  string
}

func ActInit() {


	JsNet.Http("/activitypage", activity) //活动
}

//活动入口
func activity(session *JsNet.StSession) {
	t := template.New("activity")
	// t.ParseFiles("activity.html", "music.html")

	// t.ExecuteTemplate(os.Stdout, "music", nil)
	t.Parse(tpl)
	data := Para{
		Title: "君赛",
		Debug:  `<script src="http://activity.junsie.com/static/js/vconsole.min.js" type="text/javascript" ></script>`}

	b := bytes.NewBuffer(make([]byte, 0))

	t.Execute(b, data)

	// session.Env.Ctx.WriteString(string(b.Bytes()))
}
