import { useEffect, useState } from "react";
import ParkingForm from "./ParkingForm";
import ParkingGetAll from "./ParkingGetAll";

function App() {
  // re-run when refreshTrigger changes
  const [savedData, setSavedData] = useState([]);

  useEffect(() => {
    const data = JSON.parse(localStorage.getItem("vehicleFormData")) || [];
    setSavedData(data);
  }, []);

  const handleSave = (newEntry) => {
    const updated = [...savedData, newEntry];
    setSavedData(updated);
    localStorage.setItem("vehicleFormData", JSON.stringify(updated));
  };
  return (
    <>
      <ParkingForm onSaved={handleSave} />
      <ParkingGetAll savedData={savedData} />
    </>
  );
}

export default App;
