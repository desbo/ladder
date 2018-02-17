import * as React from 'react'

import { connect, Dispatch } from 'react-redux'
import { match } from 'react-router'

import { setChartData, setCurrentLadder } from 'actions/actions'

import TitleBar from 'components/TitleBar'
import Chart from 'components/chart/Chart'

import API from 'api'

const mapStateToProps = (state: AppState) => ({
  ladder: state.ladders.current,
  chart: state.chart
})

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLadder: (ladder: Ladder) => dispatch(setCurrentLadder(ladder)),
  setChartData: (data: ChartData) => dispatch(setChartData(data)),
})

type Props = { 
  ladder: Ladder,
  chart: ChartData,
  setChartData: (data: ChartData) => any,
  setLadder: (ladder: Ladder) => any,
  match: match<{id: string}>,
};

class ViewChart extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  componentDidMount() {
    if (!this.props.ladder) {
      this.fetchLadder()
    }

    if (!this.props.chart) {
      this.fetchChart()
    }
  }

  fetchLadder() {
    return API.getLadder(this.props.match.params.id)
      .then(ladder => this.props.setLadder(ladder))
  }

  fetchChart() {
    return API.getChart(this.props.match.params.id)
      .then(data => this.props.setChartData(data))
  }
  
  render() {
    if (this.props.ladder && this.props.chart) {
      return (
        <div>
          <TitleBar ladder={this.props.ladder} />

          <section className="section">
            <div className="container">
              <Chart data={this.props.chart} />
            </div>
          </section>
        </div>
      );
    } else {
      return <div className="container">getting data...</div>
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ViewChart)