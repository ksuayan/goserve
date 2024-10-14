import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';

const GlucoseChart = ({ data }) => {
  const svgRef = useRef(); // Create a ref to attach to the SVG element
  const targetMin = 80;
  const targetMax = 130;

  useEffect(() => {
    // Define dimensions for the chart
    const margin = { top: 20, right: 30, bottom: 30, left: 40 };
    const width = 800 - margin.left - margin.right;
    const height = 400 - margin.top - margin.bottom;

    // Clear the SVG on each render
    d3.select(svgRef.current).selectAll('*').remove();

    // Create the SVG element
    const svg = d3
      .select(svgRef.current)
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Parse the time data (assuming device_timestamp is a string)
    const parseTime = d3.timeParse('%Y-%m-%dT%H:%M:%S%Z'); // Modify format based on your data

    // Format the data (parsing the time)
    const formattedData = data.map((d) => ({
      timestamp: parseTime(d.timestamp),
      glucose: d.glucose
    }));

    // Create scales
    const x = d3
      .scaleTime()
      .domain(d3.extent(formattedData, (d) => d.timestamp))
      .range([0, width]);

    const y = d3
      .scaleLinear()
      .domain([0, d3.max(formattedData, (d) => d.glucose)])
      .range([height, 0]);

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

    // Add X axis
    svg
      .append('g')
      .attr('transform', `translate(0,${height})`)
      .call(d3.axisBottom(x).ticks(5));

    // Add Y axis
    svg.append('g').call(d3.axisLeft(y));

    // Add the line
    svg
      .append('path')
      .datum(formattedData)
      .attr('fill', 'none')
      .attr('stroke', 'steelblue')
      .attr('stroke-width', 1.5)
      .attr(
        'd',
        d3
          .line()
          .x((d) => x(d.timestamp))
          .y((d) => y(d.glucose))
      );
  }, [data]); // Rerun the effect when data changes

  return <svg ref={svgRef}></svg>;
};

export default GlucoseChart;
