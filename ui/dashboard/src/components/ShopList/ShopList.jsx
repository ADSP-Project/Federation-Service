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


const ShopList = ({shop}) => {
  console.log(shop)
  
  const [currentShop, setCurrentShop] = useState({});
  const [partners, setPartners] = useState([]);
  const [message, setMessage] = useState('');


  useEffect(() => {
    const fetchShopData = async () => {
      const res = await fetch(`${import.meta.env.VITE_FEDERATION_SERVICE}/api/v1/shop`);
      const data = await res.json();
      setCurrentShop(data);
    };
    fetchShopData();
  }, []);


  useEffect(() => {
    fetch(`${import.meta.env.VITE_FEDERATION_SERVICE}/api/v1/partners`)
      .then(response => response.json())
      .then(data => setPartners(data));
  }, []);

  const handleAcceptButtonClick = async () => {

    {/* Here correct shop ID needs to be adapted */}
    const partnershipRequest = {
      shopId: currentShop.shopId
  
    };
  
    const res = await fetch(`${import.meta.env.VITE_FEDERATION_SERVICE}/api/v1/partnerships/accept`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(partnershipRequest)
    });
  
    if(res.status === 200) {
      setMessage('Request sent');
    } else {
      setMessage("Failed to send request");
    }
    
  };

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
                <Button onClick={handleAcceptButtonClick}>Approve</Button>
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
