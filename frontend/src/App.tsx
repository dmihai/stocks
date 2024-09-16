import { useState, useEffect } from 'react';
import { getTopGainers, TopGainer } from './api/api';
import './App.css';

function App() {
  const [topGainers, setTopGainers] = useState<TopGainer[]>([]);

  useEffect(() => {
    const intervalId = setInterval(() => {
      const res = getTopGainers();
      res.then((response) => setTopGainers(response));
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  return (
    <div className="App">
      <table className="table caption-top">
        <caption>Top gainers</caption>
        <thead>
          <tr>
            <th scope="col">Symbol</th>
            <th scope="col">% changed</th>
            <th scope="col">yesterday close</th>
            <th scope="col">yesterday volume</th>
            <th scope="col">current price</th>
            <th scope="col">current volume</th>
          </tr>
        </thead>
        <tbody>
          {topGainers.map((topGainer) => (
            <tr>
              <td>{topGainer.symbol}</td>
              <td className="text-end">{Math.round(topGainer.percentChange)}</td>
              <td className="text-end">{topGainer.yesterday.close}</td>
              <td className="text-end">{topGainer.yesterday.volume}</td>
              <td className="text-end">{topGainer.current.price}</td>
              <td className="text-end">{topGainer.current.volume}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;
