import React, { useState, useEffect } from 'react';
import GlucoseData from './GlucoseData';

const GlucoseDataContainer = ({ fromDate, toDate }) => {
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
      <h2>Raw Data</h2>
      <GlucoseData data={glucoseData} />
    </div>
  );
};

export default GlucoseDataContainer;
