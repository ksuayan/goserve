import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

const StreamGraph = ({ data }) => {
  const svgRef = useRef();
  const targetMin = 80;
  const targetMax = 130;

  useEffect(() => {
    // Set up SVG dimensions and margins
    const margin = { top: 20, right: 30, bottom: 30, left: 40 };
    const width = 800 - margin.left - margin.right;
    const height = 400 - margin.top - margin.bottom;

    // Clear the SVG if it already contains elements
    d3.select(svgRef.current).selectAll('*').remove();

    // Create the SVG container
    const svg = d3
      .select(svgRef.current)
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Parse the timeslot data
    const parseTime = d3.timeParse('%H:%M:%S');

    data.forEach((d) => {
      d.timeslot = parseTime(d.timeslot);
    });

    // Set up x and y scales
    const x = d3
      .scaleTime()
      .domain(d3.extent(data, (d) => d.timeslot))
      .range([0, width]);

    const y = d3
      .scaleLinear()
      .domain([0, d3.max(data, (d) => d.pct_95)])
      .nice()
      .range([height, 0]);

    // Define the area generators for each percentile
    const area95 = d3
      .area()
      .x((d) => x(d.timeslot))
      .y0(y(0))
      .y1((d) => y(d.pct_95))
      .curve(d3.curveBasis);

    const area05 = d3
      .area()
      .x((d) => x(d.timeslot))
      .y0(y(0))
      .y1((d) => y(d.pct_05))
      .curve(d3.curveBasis);

    const areaAvg = d3
      .area()
      .x((d) => x(d.timeslot))
      .y0(y(0))
      .y1((d) => y(d.average))
      .curve(d3.curveBasis);

    // Add 95th percentile area
    svg
      .append('path')
      .datum(data)
      .attr('fill', 'steelblue')
      .attr('opacity', 0.3)
      .attr('d', area95);

    // Add 5th percentile area
    svg
      .append('path')
      .datum(data)
      .attr('fill', 'white')
      .attr('opacity', 1)
      .attr('d', area05);

    // Add average line
    svg
      .append('path')
      .datum(data)
      .attr('fill', 'none')
      .attr('stroke', 'black')
      .attr('stroke-width', 1)
      .attr('d', areaAvg);

    // Append horizontal lines for target range
    svg
      .append('line')
      .attr('x1', 0)
      .attr('x2', width)
      .attr('y1', y(targetMin))
      .attr('y2', y(targetMin))
      .attr('stroke', 'green')
      .attr('stroke-width', 0.5)
      .attr('stroke-dasharray', '4 2') // Dashed line for targetMin
      .attr('class', 'target-line');

    svg
      .append('line')
      .attr('x1', 0)
      .attr('x2', width)
      .attr('y1', y(targetMax))
      .attr('y2', y(targetMax))
      .attr('stroke', 'green')
      .attr('stroke-width', 0.5)
      .attr('stroke-dasharray', '4 2') // Dashed line for targetMax
      .attr('class', 'target-line');

    // Add labels for target range
    svg
      .append('text')
      .attr('x', width - margin.right)
      .attr('y', y(targetMin) - 5)
      .attr('fill', 'green')
      .style('text-anchor', 'end')
      .text(`${targetMin} mg/dL`);

    svg
      .append('text')
      .attr('x', width - margin.right)
      .attr('y', y(targetMax) - 5)
      .attr('fill', 'green')
      .style('text-anchor', 'end')
      .text(`${targetMax} mg/dL`);

    // Add x-axis
    svg
      .append('g')
      .attr('transform', `translate(0,${height})`)
      .call(
        d3
          .axisBottom(x)
          .ticks(d3.timeMinute.every(600))
          .tickFormat(d3.timeFormat('%H:%M'))
      );

    // Add y-axis
    svg.append('g').call(d3.axisLeft(y));
  }, [data]);

  return <svg ref={svgRef}></svg>;
};

export default StreamGraph;
