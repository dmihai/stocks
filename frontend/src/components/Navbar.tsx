import LiveSwitch from './LiveSwitch';

type Props = {
  isLive: boolean;
  setLive: (status: boolean) => void;
};

function Navbar(props: Props) {
  return (
    <nav className="navbar sticky-top bg-body-tertiary">
      <div className="container-fluid">
        <span className="navbar-brand mb-0 h1">Stocks Scanner</span>
        <LiveSwitch isLive={props.isLive} setLive={props.setLive} />
      </div>
    </nav>
  );
}

export default Navbar;
