import cassandra from 'cassandra-driver'

const TEST_KEYSPACE = 'my_chat_app_test'

const TRUNCATE_TABLES = [
  'users_by_id',
  'users_by_username',
  'friendships',
  'friend_requests',
  'channels',
  'channels_by_user',
  'members_by_channel',
  'messages_by_channel',
]

export default async function globalTeardown() {
  console.log('\n[global-teardown] Truncating test database...')
  const client = new cassandra.Client({
    contactPoints: ['127.0.0.1'],
    localDataCenter: 'datacenter1',
  })

  await client.connect()
  await client.execute(`USE ${TEST_KEYSPACE}`)

  for (const table of TRUNCATE_TABLES) {
    await client.execute(`TRUNCATE ${table}`)
  }

  await client.shutdown()
  console.log('[global-teardown] Done.\n')
}
