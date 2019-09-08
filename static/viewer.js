let host = document.currentScript.getAttribute('host');
window.addEventListener("load", function(evt) {
    document.getElementById("spinner").style.display = 'none';
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var mouse_moved = false;

    var print = function(message) {
        var d = document.createElement("li");
		d.setAttribute("class","list-group-item")
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
		}
        h = document.getElementById("canvas").getAttribute("height");
        w = document.getElementById("canvas").getAttribute("width");
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
                if (feedback.action == "loaded"){
                    document.getElementById("spinner").style.display = 'none';
                }
                if (feedback.action == "loading"){
                        document.getElementById("spinner").style.display = 'block';
                }
				print(evt.data)
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

	document.getElementById("canvas").onmousemove = function(evt){
			if (!ws) {return false;}
            mouse_moved = true;
			var rect = evt.target.getBoundingClientRect();
			var x = (evt.clientX - rect.left); 
			var y = (evt.clientY - rect.top); 
			ws.send(`{"x":${x},"y":${y}, "cmd":""}`);
			return false;	
	}

	document.getElementById("canvas").onwheel = function(evt){
		evt.preventDefault();
			if (!ws) {return false;}
			ws.send(`{"x":${evt.deltaX},"y":${evt.deltaY}, "cmd":"zoom"}`);
			return false;	
	}

	document.getElementById("canvas").oncontextmenu = function(evt){
		evt.preventDefault();
        return false;
    }

	document.getElementById("canvas").onmousedown = function(evt){
		if (!ws) {return false;}
        evt.preventDefault();
        mouse_moved = false;
		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
		var y = (evt.clientY - rect.top); 
		ws.send(`{"x":${x},"y":${y}, "cmd":"mousedown", "val":"${evt.button}"}`);
		return false;	
	}

    document.getElementById("canvas").onmouseup = function(evt){
		evt.preventDefault();
		if (!ws) {return false;}
		var rect = evt.target.getBoundingClientRect();
		var x = (evt.clientX - rect.left); 
		var y = (evt.clientY - rect.top); 
		ws.send(`{"x":${x},"y":${y}, "cmd":"mouseup", "val":"${evt.button}", "moved":${mouse_moved}}`);
		return false;	
	}

    document.getElementById("canvas").onkeydown = function(e) {
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keydown", "val":"${e.keyCode}"}`);
        return false;   
    }

    document.getElementById("canvas").onkeyup = function(e) {
        e = e || window.event;
        if (!ws) {return false;}
        ws.send(`{"cmd":"keyup", "val":"${e.keyCode}"}`);
        return false;   
    }

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {return false;}
		print(`Set Field of View: ${input.value}`);
        ws.send(`{"cmd":"fov", "val":"${input.value}"}`);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {return false;}
		ws.send(`{"cmd":"close"}`);
        ws.close();
        return false;
	};
});
