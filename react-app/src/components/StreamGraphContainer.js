import React, { useState, useEffect } from 'react';
import StreamGraph from './StreamGraph';

const GlucoseDataContainer = ({ fromDate, toDate }) => {
  const [glucoseData, setGlucoseData] = useState([]);

  useEffect(() => {
    // Fetch data from your GetGlucoseData API
    const fetchData = async () => {
      const response = await fetch(`/api/glucose/${fromDate}/${toDate}`);
      const data = await response.json();
      console.log('/api/glucose/...', data);
      setGlucoseData(data);
    };

    fetchData();
  }, [fromDate, toDate]);

  return (
    <div>
      <h2>Avg Glucose, 5th pct &amp; 95th pct</h2>
      <StreamGraph data={glucoseData} />
    </div>
  );
};

export default GlucoseDataContainer;
