package devdash // import "github.com/flicaflow/devdash"

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	updateInputCh chan *updateData
	updates       chan map[string]string
	running       bool = false
)

type updateData struct {
	Id   string
	Html string
}

func Start() {
	if running {
		panic("DevDash already started")
	}
	running = true

	// buffered and therefore none-blocking update Ch
	updateInputCh = make(chan *updateData, 1000)
	updates = make(chan map[string]string)

	go updateManager()

	http.HandleFunc("/", index)
	http.HandleFunc("/update", update)
	go http.ListenAndServe(":27272", nil)
}

// the update manager collates updates from the input channel
func updateManager() {
	data := make(map[string]string)
	hasData := false
	for {
		if !hasData {
			d := <-updateInputCh
			data[d.Id] = d.Html
			hasData = true

		} else {
			select {
			case updates <- data:
				// clear data
				data = make(map[string]string)
				hasData = false
			case d := <-updateInputCh:
				data[d.Id] = d.Html
			}
		}
	}
}

// Message is just a demo dashboard
func Message(id string, text string) {
	updateInputCh <- &updateData{
		"message:" + id,
		"<h2>" + text + "</h2>",
	}
}

func GraphValues(id string, valueName string, value float64) {

}

func update(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
		// client closed connection
	case ud := <-updates:
		// got an update
		json.NewEncoder(w).Encode(ud)

	}
}

func index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, `<!doctype html>
	<html>
	<head>
	<title>DevDash</title>
	
	<style>

	article {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
	}

	section {
		border: 1px solid #888;
		margin: 5px;
		padding: 5px;
	}

	</style>

	</head>
	<body>
		<nav>
			<h1>DevDash</h1>
		</nav>
		<article id="content">

		</article>

		<script>
		function update() {
			fetch("/update")
				.then(a => a.json() )
				.then(data => {
					for(let id in data) {
						let html = data[id];
						let el = document.getElementById(id);
						if(el === null) {
							el = document.createElement("section");
							el.setAttribute("id", id);

							let c = document.getElementById("content");
							c.appendChild(el);
						}

						el.innerHTML = html;

					}
					window.requestAnimationFrame(update);
				})
				.catch(r => {
					console.log("Error reason: ", r);
					setTimeout(update, 5000);
				});
		}
		update();
		</script>
	</body>
	</html>`)
}
