type Props = {
  isLive: boolean;
  setLive: (status: boolean) => void;
};

function LiveSwitch(props: Props) {
  return (
    <div className="form-check form-switch">
      <input
        className="form-check-input"
        type="checkbox"
        role="switch"
        id="liveStatus"
        checked={props.isLive}
        onChange={(e) => props.setLive(e.target.checked)}
      />
      <label className="form-check-label" htmlFor="liveStatus">
        Live
      </label>
    </div>
  );
}

export default LiveSwitch;
