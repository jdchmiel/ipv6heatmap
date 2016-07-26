var map = L.map('map').setView([-37.87, 175.475], 12);
var tiles = L.tileLayer('https://api.mapbox.com/styles/v1/jdchmiel/cir2dz7hl001ebknkrpyooepv/tiles/256/{z}/{x}/{y}?access_token=pk.eyJ1IjoiamRjaG1pZWwiLCJhIjoiY2lyMmR1d3g1MDMxa2ZxbThweXN5MW81ZiJ9.vCLpDDTnien_h_G8SAOhQw', {
    attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',
    id : "jdchmiel",
    accessToken: "pk.eyJ1IjoiamRjaG1pZWwiLCJhIjoiY2lyMmR3Y3poMDMxdmZsbThwbzM4d2VzZCJ9.D97ij5EYdsNGnbrIFiavzw"
}).addTo(map);

    function loadPoints() {
        var addressPoints;
        var xmlhttp = new XMLHttpRequest();
        var bounds = map.getBounds();
        var min = bounds.getSouthWest().wrap();
        var max = bounds.getNorthEast().wrap();

        xmlhttp.onreadystatechange = function() {
            if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
                if (xmlhttp.status == 200) {
                    points = JSON.parse(xmlhttp.responseText);
                    points = points.map(function (p) { return [p['Lat'], p['Lon']]; });

                    var heat = L.heatLayer(points).addTo(map);
                }
                else if (xmlhttp.status == 400) {
                    alert('There was an error 400');
                }
                else {
                    alert('something else other than 200 was returned');
                }
            }
        };

        xmlhttp.open("GET", "/api/coordinates?minLat="+min.lat+"&minLon="+min.lng+"&maxLat="+max.lat+"&maxLon="+max.lng , true);
        xmlhttp.send();
    }
map.on("moveend", onMapMove)


function onMapMove(e){
    loadPoints();
}