import { Tile, TileHeader, TileBody, TileFooter, JoinButton } from './ShopTile.styles';

const ShopTile = ({ shop }) => (
  <Tile>
    <TileHeader>{shop.name}</TileHeader>
    <TileBody>{shop.webhookURL}</TileBody>
    <TileFooter>
      <JoinButton>Join</JoinButton>
    </TileFooter>
  </Tile>
);

export default ShopTile;
