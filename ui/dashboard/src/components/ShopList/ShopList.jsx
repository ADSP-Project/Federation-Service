import { Table, Button } from "./ShopList.styles";

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
    </tbody>
  </Table>
);



export default ShopList;
