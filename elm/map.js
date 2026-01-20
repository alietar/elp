let map = null;
let squaresGroup = L.featureGroup();

function initMap(app) {
  app.ports.initMap.subscribe(function (config) {
    map = L.map("map").setView(
      [config.lat, config.lon],
      config.zoom
    );

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "© OpenStreetMap contributors"
    }).addTo(map);
    
    squaresGroup.addTo(map)
  });

  app.ports.clearSquares.subscribe(function () {
    squaresGroup.clearLayers();
  });

  app.ports.drawSquare.subscribe(function (bounds) {
    drawSquare(bounds);
  });

  app.ports.autoView.subscribe(function () {
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