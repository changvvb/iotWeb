$(function () {
        "use strict";
        var url = "http://127.0.0.1:8000/monitor/humi_json/";
        $.getJSON(url,{}, function(res){
            console.log(res);
            /* Transformando o dicionário em lista.
               Com o comando map eu coloco uma lista dentro da outra,
               necessário para este tipo de gráfico. */
            var data = res.humidity.map(function (v) {
                return [v.dia, v.humi]
            }
          );


        Highcharts.setOptions({
            global: {
                useUTC: false
            }
        });


         console.log(data);

         $('#humi-chart').highcharts({
             chart: {
                 type: 'spline',
                 animation: Highcharts.svg,
                 marginRight: 10,
                 events: {
                   load: function() {
                     var series = this.series[0];

                     setInterval(function(){
                         console.log(data);
                         var x = data[0].via,

                             y = data[0].temp;

                         series.addPoint([x, y], true, true);
                     },1000);
                   }
                 }
             },
             title: {
                 text: '湿度实时监测'
             },

             xAxis: {

                 type: 'datetime',//时间轴要加上这个type，默认是linear
                 tickPixelInterval: 150
             },
             yAxis: {

                 title: {
                     text: '相对湿度（RH）'
                 },
                 plotLines: [{
                   value:0,
                   width: 1,
                   color: '#808080'
                 }]

             },
             tooltip: {
                 formatter: function () {
                     return '<b>' + this.series.name + '</b><br/>' +
                         Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                         Highcharts.numberFormat(this.y, 2);
                 }
             },
             legend: {
                 enabled: false
             },

             exporting: {
                   enabled: false
               },

            series: [{
                name: 'Random data',
                data: /*(function () {
                    // generate an array of random data

                    var data = [],
                        time = (new Date()).getTime(),
                        i;

                    for (i = -19; i <= 0; i += 1) {
                        data.push({
                            x: time + i * 1000,
                            y: Math.random()*(19-18)+18
                        });
                    }
                    return data;
                }())*/data
            }]
        });
    });
});
