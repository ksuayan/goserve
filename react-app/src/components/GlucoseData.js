import React, { useEffect, useRef, useState } from 'react';
import * as d3 from 'd3';

const GlucoseData = ({ data }) => {
  const svgRef = useRef(); // Create a ref to attach to the SVG element
  const [tooltip, setTooltip] = useState({
    display: false,
    x: 0,
    y: 0,
    timestamp: '',
    glucose: 0
  });

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

    // Add X axis
    svg
      .append('g')
      .attr('transform', `translate(0,${height})`)
      .call(d3.axisBottom(x).ticks(5));

    // Add Y axis
    svg.append('g').call(d3.axisLeft(y));

    // Append axes to the svg
    svg
      .append('g')
      .attr('transform', `translate(0, ${height - margin.bottom})`)
      .call(x);
    svg.append('g').attr('transform', `translate(${margin.left}, 0)`).call(y);

    // Add the line
    svg
      .append('path')
      .datum(formattedData)
      .attr('fill', 'none')
      .attr('stroke', 'steelblue')
      .attr('stroke-width', 1)
      .attr(
        'd',
        d3
          .line()
          .x((d) => x(d.timestamp))
          .y((d) => y(d.glucose))
      );

    svg
      .selectAll('circle')
      .data(data)
      .enter()
      .append('circle')
      .attr('cx', (d) => x(new Date(d.timestamp)))
      .attr('cy', (d) => y(d.glucose))
      .attr('r', 5)
      .attr('fill', 'orange')
      .attr('opacity', 0.6)
      .on('mouseover', (event, d) => {
        setTooltip({
          display: true,
          x: x(new Date(d.timestamp)),
          y: y(d.glucose),
          timestamp: d.timestamp,
          glucose: d.glucose
        });
      })
      .on('mouseout', () => {
        setTooltip({ display: false, x: 0, y: 0, timestamp: '', glucose: 0 });
      });
  }, [data]); // Rerun the effect when data changes

  return (
    <div style={{ position: 'relative' }}>
      <svg ref={svgRef} width={800} height={400}></svg>

      {tooltip.display && (
        <div
          style={{
            position: 'absolute',
            left: tooltip.x + 10,
            top: tooltip.y + 10,
            backgroundColor: 'rgba(0, 0, 0, 0.7)',
            color: 'white',
            padding: '5px',
            borderRadius: '3px',
            pointerEvents: 'none'
          }}
        >
          <div>
            <strong>Time:</strong>{' '}
            {new Date(tooltip.timestamp).toLocaleString()}
          </div>
          <div>
            <strong>Glucose:</strong> {tooltip.glucose} mg/dL
          </div>
        </div>
      )}
    </div>
  );
};

export default GlucoseData;
