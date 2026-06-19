function ParkingGetAll({ savedData }) {
  //const savedData = JSON.parse(localStorage.getItem("vehicleFormData")) || [];

  return (
    <ul>
      {savedData.map((item, index) => (
        <li key={index}>
          {item.name} - {item.vehicleNumber} - {item.duration} hours
        </li>
      ))}
    </ul>
  );
}

export default ParkingGetAll;
