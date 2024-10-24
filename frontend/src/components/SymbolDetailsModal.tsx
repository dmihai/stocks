import { useEffect, useState } from 'react';

import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import { getSymbolDetails, SymbolDetails } from '../api/api';

type Props = {
  symbol: string;
  handleClose: () => void;
};

function SymbolDetailsModal(props: Props) {
  const [details, setDetails] = useState<SymbolDetails>({});

  useEffect(() => {
    setDetails({});
    if (props.symbol !== '') {
      const fetchData = async () => {
        const details = await getSymbolDetails(props.symbol);
        setDetails(details);
      };

      fetchData().catch(console.error);
    }
  }, [props.symbol]);

  return (
    <Modal show={!!props.symbol} onHide={props.handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>{props.symbol}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <dl className="row">
          <dt className="col-sm-3">Symbol</dt>
          <dd className="col-sm-9">{details.symbol}</dd>
          <dt className="col-sm-3">Name</dt>
          <dd className="col-sm-9">{details.name}</dd>
          <dt className="col-sm-3">Industry</dt>
          <dd className="col-sm-9">{details.industry}</dd>
          <dt className="col-sm-3">Sector</dt>
          <dd className="col-sm-9">{details.sector}</dd>
          <dt className="col-sm-3">IPO date</dt>
          <dd className="col-sm-9">{details.ipoDate}</dd>
          <dt className="col-sm-3">Shares outstanding</dt>
          <dd className="col-sm-9">{details.sharesOutstanding}</dd>
        </dl>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={props.handleClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default SymbolDetailsModal;
