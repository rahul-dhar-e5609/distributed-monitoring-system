function createChart(parentNode, data) {
  var node = $(
  '<div class="col-sm-12 panel panel-default">' +
  '  <div class="panel-heading row">' +
  '    <div class="col-sm-2 text-right">Name</div>' +
  '    <div class="col-sm-10">' +
  '      <b>' + data.name + '</b>' +
  '    </div>' +
  '    <div class="col-sm-2 text-right">Serial No</div>' +
  '  <div class="col-sm-10">' +
  '      <b>' + data.serialNo + '</b>' +
  '    </div>' +
  '    <div class="col-sm-2 text-right">SafeRange</div>' +
  '    <div class="col-sm-10">' +
  '      <b>' +
  '       <span>'
          + data.minSafeValue +
  '       </span> - ' +
  '       <span>' +
          + data.maxSafeValue +
'         </span>' +
  '        <span>' + data.unitType + '</span>' +
  '      </b>' +
  '    </div>' +
  '  </div>' +
  '  <div class="panel-body ' + data.name + '"' +
  '    style="width:100%;height:200px;" >' +
  '  </div>' +
  '</div>');

  $(parentNode).append(node);

  var avg = (data.maxSafeValue + data.minSafeValue) / 2;
  var range = (data.maxSafeValue - data.minSafeValue);

  $('.' + data.name).CanvasJSChart({
    axisY: {
      minimum: avg - range * 0.6,
      maximum: avg + range * 0.6
    },
    axisX: {
      labelFormatter: function(e) {
        return e.value.toTimeString().substr(0,8);
      }
    },
    data: [
      {
        type: "line",
        dataPoints: []
      },
      {
          type: 'rangeArea',
          dataPoints: [],
          color: 'green',
          fillOpacity: 0.2
      }
    ]
  });
  var chart = $('.' + data.name).CanvasJSChart();
  var pts = chart.options.data[0].dataPoints;
  var range = chart.options.data[1].dataPoints;
  setInterval(function() {
    pts.push({x: new Date(), y: Math.random() * (15.3 - 15.1) + 15.1});
    while (pts.length > 20) {
      pts.shift();
    }
    range[0] = {x: pts[0].x, y: [data.minSafeValue, data.maxSafeValue]};
    range[1] = {x: pts[pts.length-1].x, y:[data.minSafeValue, data.maxSafeValue]};
    chart.render();
  }, 500);
}
