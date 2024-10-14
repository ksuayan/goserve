import { useState, useEffect } from 'react';
import './App.css';
import GlucoseChartContainer from './components/GlucoseChartContainer';
import StreamGraphContainer from './components/StreamGraphContainer';

function App() {
  const [fromDate, setFromDate] = useState(null);
  const [toDate, setToDate] = useState(null);

  useEffect(() => {
    // Extract query parameters from URL hash
    const hash = window.location.hash;
    const params = new URLSearchParams(hash.split('?')[1]); // Extract after `#?`

    const from = params.get('fromDate');
    const to = params.get('toDate');

    // Set state
    // Set state if valid
    if (from && !isNaN(new Date(from))) {
      setFromDate(from);
    }
    if (to && !isNaN(new Date(to))) {
      setToDate(to);
    }
  }, []);

  return (
    <div className='App'>
      <header className='App-header'>
        <h1>Glucose Visualization</h1>
        <p>with MySQL, Go Lang, ReactJS and D3.js.</p>
      </header>
      <main>
        {fromDate && toDate ? (
          <>
            <GlucoseChartContainer fromDate={fromDate} toDate={toDate} />
            <StreamGraphContainer fromDate={fromDate} toDate={toDate} />
          </>
        ) : (
          <p>Please provide a valid fromDate and toDate in the URL.</p>
        )}
      </main>
      <footer className='App-footer'>
        <p>
          CGM data loaded into MySQL. Experiment with GoRM and Gorilla MUX as
          framework for web server.
        </p>
      </footer>
    </div>
  );
}

export default App;
