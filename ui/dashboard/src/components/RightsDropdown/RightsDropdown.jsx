import { RightsTagTrue, RightsTagFalse } from "./RightsDropdown.styles";
import { useState } from "react";

const rightsMapping = {
  canEarnCommission: "Earn Commission",
  canShareInventory: "Share Inventory",
  canShareData: "Share Data",
  canCoPromote: "Co-Promote",
  canSell: "Sell",
};

const RightsChip = ({ option, selected, onClick }) => {
  const displayText = rightsMapping[option];
  return selected ? (
    <RightsTagTrue onClick={onClick}>{displayText}</RightsTagTrue>
  ) : (
    <RightsTagFalse onClick={onClick}>{displayText}</RightsTagFalse>
  );
};
  
const RightsDropdown = ({ options }) => {
  const [selectedRights, setSelectedRights] = useState([]);

  const handleClick = (option) => {
    setSelectedRights(prevState =>
      prevState.includes(option)
        ? prevState.filter(o => o !== option)
        : [...prevState, option]
    );
  };

  return (
    <div>
      {options.map(option => (
        <RightsChip
          key={option}
          option={option}
          selected={selectedRights.includes(option)}
          onClick={() => handleClick(option)}
        />
      ))}
    </div>
  );
};

export default RightsDropdown;
  