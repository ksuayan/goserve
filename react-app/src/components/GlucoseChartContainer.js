import React, { useState, useEffect } from 'react';
import GlucoseChart from './GlucoseChart';

const GlucoseChartContainer = ({ fromDate, toDate }) => {
  const [glucoseData, setGlucoseData] = useState([]);

  useEffect(() => {
    // Fetch data from your GetGlucoseData API
    const fetchData = async () => {
      const response = await fetch(`/api/raw/${fromDate}/${toDate}`);
      const data = await response.json();
      console.log('/raw/...', data);
      setGlucoseData(data);
    };

    fetchData();
  }, [fromDate, toDate]);

  return (
    <div>
      <h2>Raw Chart</h2>
      <GlucoseChart data={glucoseData} />
    </div>
  );
};

export default GlucoseChartContainer;
