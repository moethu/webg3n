
let host = document.currentScript.getAttribute('host');

window.addEventListener("load", function(evt) {
    var canvas = document.getElementById("canvas");
    var spinner = document.getElementById("spinner");
    var selection_ui = document.getElementById("selection");
    var ws;
    var mouse_moved = false;
    var prev_x = undefined;
    var prev_y = undefined;

    spinner.style.display = 'none';

    var print = function(message) {
        console.log(message);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
		}
        h = canvas.getAttribute("height");
        w = canvas.getAttribute("width");
		ws = new WebSocket(`${host}?h=${h}&w=${w}`);
		
        ws.onopen = function(evt) {
            print("Connected to Server");
        }
        ws.onclose = function(evt) {
			print("Closed Connection");
            ws = null;
        }
        ws.onmessage = function(evt) {
			if (evt.data.startsWith('{') && evt.data.endsWith('}')){
                var feedback = JSON.parse(evt.data);
                print(evt.data)
                if (feedback.action == "loaded"){
                    spinner.style.display = 'none';
                }
                if (feedback.action == "loading"){
                    spinner.style.display = 'block';
                }
                if (feedback.action == "selected") {
                    if (feedback.value == "") { 
                        selection_ui.innerHTML = "No selection"
                    } else {
                        selection_ui.innerHTML = `Selected Node <span class="badge badge-secondary">${feedback.value}</span>`;
                    }
                }
			}else{
				var ctx = document.getElementById('canvas').getContext('2d');
				var img = new Image();
				img.onload = function() {ctx.drawImage(img, 0, 0);};
				img.src = 'data:image/jpeg;base64,'+evt.data;
			}
        }
        ws.onerror = function(evt) {
            print("Error: " + evt.data);
        }
        return false;
    };

	canvas.onmousemove = function(evt){
			if (!ws) {return false;}
            mouse_moved = true;
			var rect = evt.target.getBoundingClientRect();
			var x = (evt.clientX - rect.left); 
			var y = (evt.clientY - rect.top); 
			ws.send(`{"x":${x},"y":${y}, "cmd":""}`);
			return false;	
	}

	canvas.onwheel = function(evt){
		evt.preventDefault();
			if (!ws) {return false;}
			ws.send(`{"x":${evt.deltaX},"y":${evt.deltaY}, "cmd":"zoom"}`);
			return false;	
	}

	canvas.oncontextmenu = function(evt){
		evt.preventDefault();
        return false;
    }

	canvas.onmousedown = function(evt){
		if (!ws) {return false;}
        evt.preventDefault();

		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
        var y = (evt.clientY - rect.top); 
        if (prev_x && prev_y) {
            console.log(Math.abs(prev_x - x), Math.abs(prev_y == y) )
            if (Math.abs(prev_x - x) < 1 && Math.abs(prev_y == y) < 1) {
                mouse_moved = false;
            }else {
                mouse_moved = true;
            }
        }
        ws.send(`{"x":${x},"y":${y}, "cmd":"mousedown", "val":"${evt.button}", "moved":${mouse_moved}}`);
        prev_x = x;
        prev_y = y;
		return false;	
	}

    canvas.onmouseup = function(evt){
		evt.preventDefault();
		if (!ws) {return false;}
		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
		var y = (evt.clientY - rect.top); 
        ws.send(`{"x":${x},"y":${y}, "cmd":"mouseup", "val":"${evt.button}", "moved":${mouse_moved}, "ctrl":${evt.ctrlKey}}`);

        // open context menu if mouse hasn't been moved
        if (evt.button == 2 && !mouse_moved) {
              var top = evt.clientY;
              var left = evt.clientX;
              $("#context-menu").css({
                display: "block",
                top: top,
                left: left
              }).addClass("show");
        }
		return false;	
	}

    this.document.onkeydown = function(e) {
        e.preventDefault();
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keydown", "val":"${e.keyCode}"}`);
        return false;   
    }

    this.document.onkeyup = function(e) {
        e.preventDefault();
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keyup", "val":"${e.keyCode}"}`);
        return false;   
    }

    document.getElementById("cmd_parallel").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"fov", "val":"5"}`);
        return false;
    };

    document.getElementById("cmd_perspective").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"fov", "val":"65"}`);
        return false;
    };

    document.getElementById("cmd_zoomextent").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"zoomextent"}`);
        return false;
    };

    document.getElementById("cmd_focus").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"focus"}`);
        return false;
    };

    document.getElementById("cmd_viewtop").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"top"}`);
        return false;
    };

    document.getElementById("cmd_viewbottom").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"bottom"}`);
        return false;
    };

    document.getElementById("cmd_viewleft").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"left"}`);
        return false;
    };

    document.getElementById("cmd_viewright").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"right"}`);
        return false;
    };

    document.getElementById("cmd_viewrear").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"rear"}`);
        return false;
    };

    document.getElementById("cmd_viewfront").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"view", "val":"front"}`);
        return false;
    };

    document.getElementById("cmd_unhideall").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"unhide", "val":""}`);
        return false;
    };

    document.getElementById("cmd_hide").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"hide"}`);
        return false;
    };

    document.getElementById("cmd_qlow").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"quality", "val":"20"}`);
        return false;
    };

    document.getElementById("cmd_qmid").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"quality", "val":"60"}`);
        return false;
    };

    document.getElementById("cmd_qhigh").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"quality", "val":"90"}`);
        return false;
    };

    document.getElementById("cmd_imagesettings").onclick = function(evt) {
        if (!ws) {return false;}
        let contrast = document.getElementById("img_contrast").value
        let saturation = document.getElementById("img_saturation").value
        let brightness = document.getElementById("img_brightness").value
        let blur = document.getElementById("img_blur").value
        ws.send(`{"cmd":"imagesettings", "val":"${brightness}:${contrast}:${saturation}:${blur}"}`);
        return false;
    };

    document.getElementById("cmd_resetimagesettings").onclick = function(evt) {
        if (!ws) {return false;}
        document.getElementById("img_contrast").value = 0
        document.getElementById("img_saturation").value = 0
        document.getElementById("img_brightness").value = 0
        document.getElementById("img_blur").value = 0
        ws.send(`{"cmd":"imagesettings", "val":"0:0:0:0"}`);
        return false;
    };

    document.getElementById("cmd_imageinvert").onclick = function(evt) {
        if (!ws) {return false;}
		ws.send(`{"cmd":"invert"}`);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {return false;}
		ws.send(`{"cmd":"close"}`);
        ws.close();
        return false;
    };
    

      $("#context-menu").on("click", function() {
        $("#context-menu").removeClass("show").hide();
      });
      
      $("#context-menu a").on("click", function() {
        $(this).parent().removeClass("show").hide();
      });
});
