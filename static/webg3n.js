let host = document.currentScript.getAttribute('host');

$(document).ready(function () {
    console.log("ready!");
    fitToContainer(canvas)
});

function fitToContainer(canvas) {
    // Make it visually fill the positioned parent
    canvas.style.width = '100%';
    canvas.style.height = '100%';
    // ...then set the internal size to match
    canvas.width = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;
}

window.addEventListener("load", function (evt) {
    var canvas = document.getElementById("canvas");
    var spinner = document.getElementById("spinner");
    var selection_ui = document.getElementById("selection");
    var ws;
    var mouse_moved = false;
    var prev_x = undefined;
    var prev_y = undefined;

    spinner.style.display = 'none';

    var print = function (message) {
        console.log(message);
    };

    document.getElementById("open").onclick = function (evt) {
        if (ws) {
            return false;
        }

        h = $('#canvas').height();
        w = $('#canvas').width();
        

        var jmuxer = new JMuxer({
            node: 'player',
            mode: 'video',
            flushingTime: 1,
            fps: 20,
            debug: false
         });

        ws = new WebSocket(`${host}?h=${h}&w=${w}`);
        ws.binaryType = 'arraybuffer';
            ws.onmessage = function (evt) {
                jmuxer.feed({video: new Uint8Array(evt.data)});               
            }
            ws.onopen = function (evt) {
                print("Connected to Server");
            }
            ws.onclose = function (evt) {
                print("Closed Connection");
                ws = null;
            }
            ws.onerror = function (evt) {
                print("Error: " + evt.data);
            }

        return false;
    };

    canvas.onmousemove = function (evt) {
        if (!ws) {
            return false;
        }
        mouse_moved = true;
        var rect = evt.target.getBoundingClientRect();
        var x = (evt.clientX - rect.left);
        var y = (evt.clientY - rect.top);
        ws.send(`{"x":${x},"y":${y}, "cmd":""}`);
        return false;
    }

    canvas.onwheel = function (evt) {
        evt.preventDefault();
        if (!ws) {
            return false;
        }
        ws.send(`{"x":${evt.deltaX},"y":${evt.deltaY}, "cmd":"Zoom"}`);
        return false;
    }

    canvas.oncontextmenu = function (evt) {
        evt.preventDefault();
        return false;
    }

    canvas.onmousedown = function (evt) {
        if (!ws) {
            return false;
        }
        evt.preventDefault();

        var rect = evt.target.getBoundingClientRect();
        var x = (evt.clientX - rect.left);
        var y = (evt.clientY - rect.top);
        checkMouseMoved(x, y);
        ws.send(`{"x":${x},"y":${y}, "cmd":"Mousedown", "val":"${evt.button}", "moved":${mouse_moved}}`);
        prev_x = x;
        prev_y = y;
        return false;
    }

    function checkMouseMoved(x, y) {
        if (prev_x && prev_y) {
            if (Math.abs(prev_x - x) < 1 && Math.abs(prev_y - y) < 1) {
                mouse_moved = false;
            } else {
                mouse_moved = true;
            }
        }
    }

    canvas.onmouseup = function (evt) {
        evt.preventDefault();
        if (!ws) {
            return false;
        }
        var rect = evt.target.getBoundingClientRect();
        var x = (evt.clientX - rect.left);
        var y = (evt.clientY - rect.top);
        checkMouseMoved(x, y);
        ws.send(`{"x":${x},"y":${y}, "cmd":"Mouseup", "val":"${evt.button}", "moved":${mouse_moved}, "ctrl":${evt.ctrlKey}}`);

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

    this.document.onkeydown = function (e) {
        e.preventDefault();
        e = e || window.event;
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Keydown", "val":"${e.keyCode}"}`);
        return false;
    }

    this.document.onkeyup = function (e) {
        e.preventDefault();
        e = e || window.event;
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Keyup", "val":"${e.keyCode}"}`);
        return false;
    }

    document.getElementById("cmd_parallel").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Fov", "val":"15"}`);
        return false;
    };

    document.getElementById("cmd_perspective").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Fov", "val":"65"}`);
        return false;
    };

    document.getElementById("cmd_zoomextent").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Zoomextent"}`);
        return false;
    };

    document.getElementById("cmd_focus").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Focus"}`);
        return false;
    };

    document.getElementById("cmd_viewtop").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"top"}`);
        return false;
    };

    document.getElementById("cmd_viewbottom").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"bottom"}`);
        return false;
    };

    document.getElementById("cmd_viewleft").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"left"}`);
        return false;
    };

    document.getElementById("cmd_viewright").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"right"}`);
        return false;
    };

    document.getElementById("cmd_viewrear").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"rear"}`);
        return false;
    };

    document.getElementById("cmd_viewfront").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"View", "val":"front"}`);
        return false;
    };

    document.getElementById("cmd_unhideall").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Unhide", "val":""}`);
        return false;
    };

    document.getElementById("cmd_hide").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Hide"}`);
        return false;
    };

    document.getElementById("cmd_debug").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Debugmode"}`);
        return false;
    };

    document.getElementById("cmd_qlow").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Quality", "val":"2"}`);
        return false;
    };

    document.getElementById("cmd_qmid").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Quality", "val":"1"}`);
        return false;
    };

    document.getElementById("cmd_encodepng").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Encoder", "val":"png"}`);
        return false;
    };

    document.getElementById("cmd_encodejpeg").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Encoder", "val":"jpeg"}`);
        return false;
    };

    document.getElementById("cmd_encodelibjpeg").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Encoder", "val":"libjpeg"}`);
        return false;
    };

    document.getElementById("cmd_qhigh").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Quality", "val":"0"}`);
        return false;
    };

    document.getElementById("cmd_imagesettings").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        let contrast = document.getElementById("img_contrast").value
        let saturation = document.getElementById("img_saturation").value
        let brightness = document.getElementById("img_brightness").value
        let blur = document.getElementById("img_blur").value
        let pixel = document.getElementById("img_pixel").value
        ws.send(`{"cmd":"Imagesettings", "val":"${brightness}:${contrast}:${saturation}:${blur}:${pixel}"}`);
        return false;
    };

    document.getElementById("cmd_resetimagesettings").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        document.getElementById("img_contrast").value = 0
        document.getElementById("img_saturation").value = 0
        document.getElementById("img_brightness").value = 0
        document.getElementById("img_blur").value = 0
        document.getElementById("img_pixel").value = 1.0
        ws.send(`{"cmd":"Imagesettings", "val":"0:0:0:0:1"}`);
        return false;
    };

    document.getElementById("cmd_imageinvert").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Invert"}`);
        return false;
    };

    document.getElementById("close").onclick = function (evt) {
        if (!ws) {
            return false;
        }
        ws.send(`{"cmd":"Close"}`);
        ws.close();
        return false;
    };


    $("#context-menu").on("click", function () {
        $("#context-menu").removeClass("show").hide();
    });

    $("#context-menu a").on("click", function () {
        $(this).parent().removeClass("show").hide();
    });
});