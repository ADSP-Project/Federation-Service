import { Table, Button } from "./Statistics.styles";

const Statistics = ({ shops }) => (
  <Table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Items sold</th>
        <th>Items clicked</th>
      </tr>
    </thead>
    <tbody>
      
      <tr key="0"> 
        <td>Demo-Shop</td>
        <td>
          <Button>View sold items from your shop (+3 today)</Button>
          
        </td>
        <td><Button>View clicked items from your shop (+198 views today)</Button> </td>
      </tr>
    </tbody>
  </Table>
);



export default Statistics;
