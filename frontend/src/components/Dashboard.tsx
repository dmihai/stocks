import { useState } from 'react';
import LiveSwitch from './LiveSwitch';
import TopGainers from './TopGainers';

function Dashboard() {
  const [isLive, setLive] = useState(false);

  return (
    <div>
      <LiveSwitch isLive={isLive} setLive={setLive} />
      <TopGainers isLive={isLive} />
    </div>
  );
}

export default Dashboard;
