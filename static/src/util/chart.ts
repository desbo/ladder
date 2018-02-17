const sortPoints = (points: Array<Point>): Array<Point> =>
  points.sort((a, b) => new Date(a.x).getTime() - new Date(b.x).getTime())

// reduce to 1 rating per player per day
export const summarise = (data: ChartData): ChartData =>
  Object.keys(data).reduce((cd, name) => {
    const sorted = sortPoints(data[name]);
      
    cd[name] = sorted.reverse().reduce((ps: Array<Point>, point: Point): Array<Point> => {
      const dayAlreadySet = ps.some(p => new Date(p.x).toDateString() === new Date(point.x).toDateString())
      return dayAlreadySet ? ps : ps.concat(point)
    }, []).reverse();

    return cd;
  }, data);