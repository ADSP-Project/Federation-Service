import { useEffect, useState } from 'react';
import { Table, Button, RightsTagTrue, RightsTagFalse } from "./ShopList.styles";

const rightsMapping = {
  canEarnCommission: "Earn Commission",
  canShareInventory: "Share Inventory",
  canShareData: "Share Data",
  canCoPromote: "Co-Promote",
  canSell: "Sell",
};

const Rights = ({ rights }) => (
  <div>
    {Object.entries(rights).map(([right, value]) => 
      value 
      ? <RightsTagTrue key={right}>{rightsMapping[right]}</RightsTagTrue> 
      : <RightsTagFalse key={right}>{rightsMapping[right]}</RightsTagFalse>
    )}
  </div>
);

const ShopList = () => {
  const [partners, setPartners] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8091/api/v1/partners')
      .then(response => response.json())
      .then(data => setPartners(data));
  }, []);

  const pendingPartners = partners.filter(partner => partner.requestStatus === 'pending');
  const establishedPartners = partners.filter(partner => partner.requestStatus === 'approved');

  console.log(partners)

  return (
    <>
      <h2>Pending Requests</h2>
      <Table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Rights</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {pendingPartners.map(partner => (
            <tr key={partner.shopId}>
              <td>{partner.shopId}</td>
              <td>{partner.shopName}</td>
              <td><Rights rights={partner.rights} /></td>
              <td>
                <Button>Approve</Button>
                <Button>Deny</Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>

      <h2>Established Partnerships</h2>
      <Table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Rights</th>
          </tr>
        </thead>
        <tbody>
          {establishedPartners.map(partner => (
            <tr key={partner.shopId}>
              <td>{partner.shopId}</td>
              <td>{partner.shopName}</td>
              <td><Rights rights={partner.rights} /></td>
            </tr>
          ))}
        </tbody>
      </Table>
    </>
  );
};

export default ShopList;
