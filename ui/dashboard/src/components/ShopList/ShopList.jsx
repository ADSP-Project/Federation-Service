import { Table, Button, DeniedButton } from "./ShopList.styles";

const ShopList = ({ shops }) => (
  <Table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {shops.map(shop => (
        <tr key={shop.id}>
          <td>{shop.name}</td>
          <td>
            <Button>Approve</Button>
            <Button>Deny</Button>
          </td>
        </tr>
      ))}
      <tr key="0"> 
        <td>Demo-Shop</td>
        <td>
          <Button>Allow item sale</Button> 
          <DeniedButton>Allow commission</DeniedButton>
          <DeniedButton>Allow data sharing</DeniedButton>
        </td>
      </tr>
    </tbody>
  </Table>
);



export default ShopList;
