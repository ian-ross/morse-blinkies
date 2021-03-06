package server

import "net/http"

const script = `<svg version="1.1"
     baseProfile="full"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink"
     xmlns:ev="http://www.w3.org/2001/xml-events"
     width="100%" height="100%"
     onload='draw()'>
  <script><![CDATA[
    function draw() {
        var parts = [], mx = 0, mn = 15;
        location.search.substr(1).split('').forEach(function(n) {
            parts.push(n);
        });
        var x1 = 0, y1 = 0, x2 = 0, y2 = 0;
        var last = "X";
        for (var i=0; i<parts.length; i++) {
            if (last !== "X" && parts[i] !== last) {
              var ln = document.createElementNS ("http://www.w3.org/2000/svg", "line");
              ln.setAttribute("x1", x2 + "%");
              ln.setAttribute("x2", x2 + "%");
              ln.setAttribute("y1", "0%");
              ln.setAttribute("y2", "100%");
              ln.setAttribute("stroke", "rgba(0,0,0,0.5)");
              ln.setAttribute("stroke-width", "2");
              document.getElementsByTagName("svg")[0].appendChild(ln);
            }
            var ln = document.createElementNS ("http://www.w3.org/2000/svg", "line");
            if (parts[i] === "1") {
              y1 = 0; y2 = 0;
            } else {
              y1 = 100; y2 = 100;
            }
            x1 = x2;
            x2 = 100 * ((i + 1) / parts.length);
            last = parts[i];
            ln.setAttribute("x1", x1 + "%");
            ln.setAttribute("x2", x2 + "%");
            ln.setAttribute("y1", y1 + "%");
            ln.setAttribute("y2", y2 + "%");
            ln.setAttribute("stroke", "rgba(0,0,0,0.5)");
            ln.setAttribute("stroke-width", "3");
            document.getElementsByTagName("svg")[0].appendChild(ln);
        }
    }
  ]]></script>
</svg>`

func (s *Server) sparkline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(script))
}
