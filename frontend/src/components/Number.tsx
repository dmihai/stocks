type Props = {
  value: number;
  precision?: number;
};

function Number(props: Props) {
  const formatted = new Intl.NumberFormat('en', {
    minimumFractionDigits: props.precision ? props.precision : 0,
    maximumFractionDigits: props.precision ? props.precision : 0,
  }).format(props.value);

  return <span>{formatted}</span>;
}

export default Number;
