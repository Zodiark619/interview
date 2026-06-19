import { useState } from "react";

function ParkingForm({ onSaved }) {
  const [formData, setFormData] = useState({
    name: "",
    vehicleNumber: "",
    duration: 0,
  });
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (
      !formData.name.trim() ||
      !formData.vehicleNumber.trim() ||
      formData.duration <= 0
    ) {
      alert("All fields must have a value");
      return;
    }

    // Get existing array from localStorage
    const existing = JSON.parse(localStorage.getItem("vehicleFormData")) || [];

    // Add new entry
    existing.push(formData);

    // Save back to localStorage
    localStorage.setItem("vehicleFormData", JSON.stringify(existing));

    alert("Form saved to localStorage!");
    console.log("Saved:", formData);
    if (onSaved) {
      onSaved(formData);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Name:
        <input
          type="text"
          name="name"
          value={formData.name}
          onChange={handleChange}
        />
      </label>
      <br />

      <label>
        Vehicle Number:
        <input
          type="text"
          name="vehicleNumber"
          value={formData.vehicleNumber}
          onChange={handleChange}
        />
      </label>
      <br />

      <label>
        Duration (hours):
        <input
          type="number"
          name="duration"
          value={formData.duration}
          onChange={handleChange}
          min="1"
        />
      </label>
      <br />

      <button type="submit">Submit</button>
    </form>
  );
}

export default ParkingForm;
