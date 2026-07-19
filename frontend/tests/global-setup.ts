import cassandra from 'cassandra-driver'

const TEST_KEYSPACE = 'my_chat_app_test'
const BACKEND_URL = 'http://localhost:8080'

const SCHEMA_STATEMENTS = [
  `CREATE TABLE IF NOT EXISTS users_by_id (
    user_id uuid PRIMARY KEY,
    created_at timestamp,
    email text,
    first_name text,
    last_name text,
    password_hash text,
    username text
  )`,
  `CREATE TABLE IF NOT EXISTS users_by_username (
    username_lower text PRIMARY KEY,
    user_id uuid,
    username text
  )`,
  `CREATE TABLE IF NOT EXISTS friendships (
    user_id uuid,
    friend_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, friend_id)
  )`,
  `CREATE TABLE IF NOT EXISTS friend_requests (
    recipient_id uuid,
    sender_id uuid,
    status text,
    created_at timestamp,
    PRIMARY KEY (recipient_id, sender_id)
  )`,
  `CREATE TABLE IF NOT EXISTS channels (
    channel_id uuid PRIMARY KEY,
    name text,
    type text,
    created_by uuid,
    created_at timestamp
  )`,
  `CREATE TABLE IF NOT EXISTS channels_by_user (
    user_id uuid,
    channel_id uuid,
    channel_name text,
    channel_type text,
    joined_at timestamp,
    PRIMARY KEY (user_id, channel_id)
  )`,
  `CREATE TABLE IF NOT EXISTS members_by_channel (
    channel_id uuid,
    user_id uuid,
    role text,
    joined_at timestamp,
    PRIMARY KEY (channel_id, user_id)
  )`,
  `CREATE TABLE IF NOT EXISTS messages_by_channel (
    channel_id uuid,
    bucket int,
    created_at timestamp,
    message_id uuid,
    author_id uuid,
    content text,
    PRIMARY KEY ((channel_id, bucket), created_at, message_id)
  ) WITH CLUSTERING ORDER BY (created_at DESC, message_id ASC)`,
]

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

async function registerUser(username: string, firstName: string, lastName: string): Promise<{ user_id: string }> {
  const resp = await fetch(`${BACKEND_URL}/api/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username,
      password: 'password123',
      email: `${username}@test.com`,
      first_name: firstName,
      last_name: lastName,
    }),
  })
  if (!resp.ok) {
    const text = await resp.text()
    throw new Error(`Failed to register ${username}: ${resp.status} ${text}`)
  }
  const data = await resp.json()
  return data.user
}

async function sendFriendRequest(senderCookie: string, recipientId: string) {
  await fetch(`${BACKEND_URL}/api/friend-requests/${recipientId}`, {
    method: 'POST',
    headers: { Cookie: senderCookie },
  })
}

async function acceptFriendRequest(recipientCookie: string, senderId: string) {
  await fetch(`${BACKEND_URL}/api/friend-requests/${senderId}`, {
    method: 'PUT',
    headers: { Cookie: recipientCookie },
  })
}

function extractCookie(setCookieHeader: string | null): string {
  if (!setCookieHeader) return ''
  const match = setCookieHeader.match(/auth_token=([^;]+)/)
  return match ? `auth_token=${match[1]}` : ''
}

async function loginAndGetCookie(username: string): Promise<string> {
  const resp = await fetch(`${BACKEND_URL}/api/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password: 'password123' }),
  })
  return extractCookie(resp.headers.get('set-cookie'))
}

export default async function globalSetup() {
  console.log('\n[global-setup] Checking if backend is running...')
  try {
    const resp = await fetch(`${BACKEND_URL}/api/me`)
    if (resp.status !== 401) {
      throw new Error(`Expected 401 from /api/me, got ${resp.status}`)
    }
  } catch (err) {
    throw new Error(
      `Backend is not running at ${BACKEND_URL}. Start it with:\n` +
      `  CASSANDRA_KEYSPACE=my_chat_app_test go run cmd/monolithic/main.go\n` +
      `Error: ${err}`
    )
  }

  console.log('[global-setup] Connecting to Cassandra...')
  const client = new cassandra.Client({
    contactPoints: ['127.0.0.1'],
    localDataCenter: 'datacenter1',
  })

  await client.connect()

  await client.execute(`CREATE KEYSPACE IF NOT EXISTS ${TEST_KEYSPACE} WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`)
  await client.execute(`USE ${TEST_KEYSPACE}`)

  for (const stmt of SCHEMA_STATEMENTS) {
    await client.execute(stmt)
  }

  console.log('[global-setup] Truncating all tables...')
  for (const table of TRUNCATE_TABLES) {
    await client.execute(`TRUNCATE ${table}`)
  }

  await client.shutdown()

  console.log('[global-setup] Registering users via API...')
  const steven = await registerUser('steven', 'Steven', 'Miller')
  const jessica = await registerUser('jessica', 'Jessica', 'Williams')
  await registerUser('bob', 'Bob', 'Johnson')

  console.log('[global-setup] Creating friendship between Steven and Jessica...')
  const stevenCookie = await loginAndGetCookie('steven')
  const jessicaCookie = await loginAndGetCookie('jessica')

  await sendFriendRequest(stevenCookie, jessica.user_id)
  await acceptFriendRequest(jessicaCookie, steven.user_id)

  console.log('[global-setup] Done. Test database is ready.\n')
}
