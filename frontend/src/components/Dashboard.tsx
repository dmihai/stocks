import { useState } from 'react';
import TopGainers from './TopGainers';
import Navbar from './Navbar';

function Dashboard() {
  const [isLive, setLive] = useState(false);

  return (
    <div>
      <Navbar isLive={isLive} setLive={setLive} />
      <div className="container-fluid">
        <div className="row mt-3">
          <div className="col-6">
            <TopGainers isLive={isLive} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
