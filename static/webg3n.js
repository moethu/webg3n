
let host = document.currentScript.getAttribute('host');

window.addEventListener("load", function(evt) {
    var canvas = document.getElementById("canvas");
    var spinner = document.getElementById("spinner");
    var selection_ui = document.getElementById("selection");
    var ws;
    var mouse_moved = false;
    var selection = [];

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
                        selection = []
                        selection_ui.innerHTML = "No selection"
                    } else {
                        selection = [feedback.value]
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
        mouse_moved = false;
		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
		var y = (evt.clientY - rect.top); 
		ws.send(`{"x":${x},"y":${y}, "cmd":"mousedown", "val":"${evt.button}"}`);
		return false;	
	}

    canvas.onmouseup = function(evt){
		evt.preventDefault();
		if (!ws) {return false;}
		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
		var y = (evt.clientY - rect.top); 
		ws.send(`{"x":${x},"y":${y}, "cmd":"mouseup", "val":"${evt.button}", "moved":${mouse_moved}}`);
		return false;	
	}

    canvas.onkeydown = function(e) {
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keydown", "val":"${e.keyCode}"}`);
        return false;   
    }

    canvas.onkeyup = function(e) {
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keyup", "val":"${e.keyCode}"}`);
        return false;   
    }

    document.getElementById("cmd_parallel").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"fov", "val":"10"}`);
        return false;
    };

    document.getElementById("cmd_perspective").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"fov", "val":"60"}`);
        return false;
    };

    document.getElementById("cmd_focus").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"focus"}`);
        return false;
    };

    document.getElementById("cmd_viewtop").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"viewtop"}`);
        return false;
    };

    document.getElementById("cmd_unhideall").onclick = function(evt) {
        if (!ws) {return false;}
        ws.send(`{"cmd":"unhide", "val":""}`);
        return false;
    };

    document.getElementById("cmd_hide").onclick = function(evt) {
        if (!ws) {return false;}
        if (selection) {
        ws.send(`{"cmd":"hide", "val":"${selection[0]}"}`);
        }
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {return false;}
		ws.send(`{"cmd":"close"}`);
        ws.close();
        return false;
	};
});
