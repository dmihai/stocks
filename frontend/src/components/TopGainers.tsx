import { useState, useEffect } from 'react';
import { getTopGainers, TopGainer } from '../api/api';
import Number from './Number';

type Props = {
  isLive: boolean;
};

function TopGainers(props: Props) {
  const [topGainers, setTopGainers] = useState<TopGainer[]>([]);

  useEffect(() => {
    if (props.isLive) {
      const intervalId = setInterval(async () => {
        const response = await getTopGainers();
        setTopGainers(response);
      }, 1000);

      return () => clearInterval(intervalId);
    }
  }, [props.isLive]);

  return (
    <div className="card border-primary-subtle">
      <h5 className="card-header bg-primary-subtle text-primary-emphasis">Top gainers</h5>
      <table className="table data-table table-sm lh-1 mb-1">
        <thead>
          <tr>
            <th scope="col">Symbol</th>
            <th scope="col" className="text-end">
              % changed
            </th>
            <th scope="col" className="text-end">
              yesterday close
            </th>
            <th scope="col" className="text-end">
              yesterday volume
            </th>
            <th scope="col" className="text-end">
              current price
            </th>
            <th scope="col" className="text-end">
              current volume
            </th>
            <th scope="col" className="text-end">
              last updated
            </th>
          </tr>
        </thead>
        <tbody className="table-group-divider">
          {topGainers.map((topGainer) => (
            <tr key={topGainer.symbol}>
              <td>{topGainer.symbol}</td>
              <td className="text-end">
                <Number value={topGainer.percentChanged} />
              </td>
              <td className="text-end">
                <Number value={topGainer.yesterday.close} precision={2} />
              </td>
              <td className="text-end">
                <Number value={topGainer.yesterday.volume} />
              </td>
              <td className="text-end">
                <Number value={topGainer.current.price} precision={2} />
              </td>
              <td className="text-end">
                <Number value={topGainer.current.volume} />
              </td>
              <td className="text-end">{topGainer.lastUpdated}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default TopGainers;
