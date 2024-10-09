import { useState, useEffect } from 'react';
import { getTopGainers, TopGainer } from '../api/api';
import LiveSwitch from './LiveSwitch';

function Dashboard() {
  const [topGainers, setTopGainers] = useState<TopGainer[]>([]);
  const [isLive, setLive] = useState(false);

  useEffect(() => {
    if (isLive) {
      const intervalId = setInterval(async () => {
        const response = await getTopGainers();
        setTopGainers(response);
      }, 1000);

      return () => clearInterval(intervalId);
    }
  }, [isLive]);

  return (
    <div>
      <LiveSwitch isLive={isLive} setLive={setLive} />
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
            <th scope="col">last updated</th>
          </tr>
        </thead>
        <tbody>
          {topGainers.map((topGainer) => (
            <tr>
              <td>{topGainer.symbol}</td>
              <td className="text-end">{Math.round(topGainer.percentChanged)}</td>
              <td className="text-end">{topGainer.yesterday.close}</td>
              <td className="text-end">{topGainer.yesterday.volume}</td>
              <td className="text-end">{topGainer.current.price}</td>
              <td className="text-end">{topGainer.current.volume}</td>
              <td className="text-end">{topGainer.lastUpdated}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Dashboard;
