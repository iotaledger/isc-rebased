import { CoinStruct, GetCoinsParams, GetObjectParams, GetOwnedObjectsParams, IotaClient, IotaObjectResponse, PaginatedCoins } from '@iota/iota-sdk/client';
import { paginatedRequest } from './page_reader';
import { consumeAliasOutput } from './alias_output';
import { Ed25519Keypair } from '@iota/iota-sdk/keypairs/ed25519';

const ENDPOINT_URL = 'http://localhost:9000';
const ISC_PACKAGE_ID = '0xb86718fdb1764518084cce16c9024d17c8b82faca232aff9c08245b996d817ac';
const GOVERNOR_ADDRESS = '0x70bc12d8964837afac5978b4e3acc61defe9427e0c975afb1f3663c186e3b1e6';

// random Keypair
const keypair = Ed25519Keypair.deriveKeypair('gospel poem coffee duty cluster plug turkey buffalo aim annual essay mushroom');

//
// Magic, don't touch.
//
const client = new IotaClient({
  url: ENDPOINT_URL,
});

async function main() {
  const objects = await paginatedRequest<IotaObjectResponse, GetOwnedObjectsParams>(x => client.getOwnedObjects(x), {
    owner: GOVERNOR_ADDRESS,
    filter: {
      MatchAll: [
        {
          StructType: '0x107a::alias_output::AliasOutput<0x2::iota::IOTA>',
        },
      ],
    },
    options: {
      showType: true,
      showContent: true,
    },
  });

  const aliasObjects = objects.filter(x => x.data?.type == '0x107a::alias_output::AliasOutput<0x2::iota::IOTA>');

  if (aliasObjects.length != 1) {
    throw new Error(`Invalid amount of Alias objects: ${aliasObjects.length}, expected: 1`);
  }

  const aliasObject = aliasObjects[0];
  const aliasOutputConsumeTX = await consumeAliasOutput(client, ISC_PACKAGE_ID, aliasObject.data?.objectId!);

  const result = await client.signAndExecuteTransaction({
    transaction: aliasOutputConsumeTX,
    signer: keypair,
  });

  console.log(result);
}

main();
