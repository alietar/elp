let map = null;
let squareLayers = [];
let pendingBoundsList = [];

function initMap(app) {

  app.ports.initMap.subscribe(function (config) {
    map = L.map("map").setView(
      [config.lat, config.lon],
      config.zoom
    );

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "© OpenStreetMap contributors"
    }).addTo(map);

    // Dessiner les rectangles reçus avant init
    pendingBoundsList.forEach(drawSquare);
    pendingBoundsList = [];
  });

  app.ports.drawSquare.subscribe(function (bounds) {
    if (!map) {
      pendingBoundsList.push(bounds);
      return;
    }

    drawSquare(bounds);
  });
}

function drawSquare(bounds) {
  const rect = L.rectangle(
    [
      [ bounds.southWest[1], bounds.southWest[0] ],
      [ bounds.northEast[1], bounds.northEast[0] ]
    ],
    { color: "blue", weight: 2, fillOpacity: 0.2 }
  ).addTo(map);

  squareLayers.push(rect);
}
