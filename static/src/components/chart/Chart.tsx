import * as React from 'react';
import { Ref } from 'react';

import * as palette from 'google-palette';

import { Chart } from 'chart.js';

type Props = { 
  data: ChartData
}

// reduce to 1 rating per player per day
const summarise = (data: ChartData): ChartData =>
  Object.keys(data).reduce((cd, name) => {
    const sorted = data[name]
      .sort((a, b) => new Date(a.x).getTime() - new Date(b.x).getTime())
      .reverse()

    cd[name] = sorted.reduce((ps: Array<Point>, point: Point): Array<Point> => {
      const dayAlreadySet = ps.some(p => new Date(p.x).toDateString === new Date(point.x).toDateString)
      return dayAlreadySet ? ps : ps.concat(point)
    }, []);

    return cd;
  }, data);

export default class ChartComponent extends React.Component<Props> {
  canvasElement: null | HTMLCanvasElement

  constructor(props: Props) {
    super(props)
  }

  componentDidMount() {
    const playerNames = Object.keys(this.props.data);
    const colours = palette('rainbow', playerNames.length);
    const summarised = summarise(this.props.data);

    new Chart(this.canvasElement, {
      type: 'line',
      options: {
        scales: {
          xAxes: [{
            type: 'time'
          }]
        }
      },
      data: {
        datasets: playerNames.map((name, i) => {
          const data = summarised[name];

          return {
            label: name,
            borderColor: `#${colours[i]}`,
            fill: false,
            data: data.map(d => ({
              x: new Date(d.x),
              y: d.y
            }))
          }
        })
      }
    })
  }

  render() {
    return (
      <div style={{'position': 'relative'}}>
        <canvas ref={e => this.canvasElement = e} id="chart"></canvas>
      </div>
    )
  }
}