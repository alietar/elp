let map = null;
let squaresGroup = L.featureGroup();
let clickMarker = null;



function initMap(app) {
  // initialise la carte
  app.ports.initMap.subscribe(function (config) {
    map = L.map("map").setView( 
      [config.lat, config.lon],
      config.zoom
    );

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "© OpenStreetMap contributors"
    }).addTo(map);

    map.on('click', function(e){ // gestion du click de l'utilisateur
    var coord = e.latlng
    var lat = coord.lat
    var long = coord.lng
    console.log("lat :" + lat +"long :" + long);

    if (clickMarker) {
        map.removeLayer(clickMarker);
    }

    clickMarker = L.marker([lat, long]).addTo(map); // ajout du marqueur à l'endroit où user clique
    map.setView([lat, long], 13);

    app.ports.click_coord.send({"lat" : lat ,"long" : long});
    

    });
    
    squaresGroup.addTo(map)
  });

  app.ports.clearSquares.subscribe(function () {
    squaresGroup.clearLayers();
  });

  app.ports.addMarker.subscribe(function (data) { //ajout du marqueur si user utilise l'interface
    const lati = data.lat;
    const lon = data.lon;

    if (clickMarker) {
        map.removeLayer(clickMarker);
    }

    clickMarker = L.marker([lati, lon]).addTo(map);
    map.setView([lati, lon], 13);
  });


  app.ports.drawSquare.subscribe(function (bounds) { // construit les rectangles sur la carte
    drawSquare(bounds);
  });

  app.ports.autoView.subscribe(function () { // règle les carrés pour qu'on les voit de loin
    setTimeout(() => {
        if (squaresGroup.getLayers().length > 0 && map) {
            map.fitBounds(squaresGroup.getBounds(), { padding: [50, 50] });
        } else {
            console.warn("Pas de zoom : groupe vide ou carte non prête");
        }
    }, 100);
  });
}

function drawSquare(bounds) {
    L.rectangle(
      [
        [ bounds.southWest[1], bounds.southWest[0] ],
        [ bounds.northEast[1], bounds.northEast[0] ]
      ],
      { color: "blue", weight: 2, fillOpacity: 0.5, stroke: false }
    ).addTo(squaresGroup);
}