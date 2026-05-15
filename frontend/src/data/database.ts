export interface User {
  id: string
  username: string
  password: string  // plaintext for mock only
  email: string
  firstName: string
  lastName: string
}

export interface Channel {
  id: string
  name: string
  type: 'public' | 'dm' | 'group'
  memberIds: string[]
}

export interface Message {
  id: string
  channelId: string
  authorId: string
  content: string
  timestamp: number
}

export interface Friendship {
  id: string
  userId1: string
  userId2: string
}

// Five generic users
export const users: User[] = [
  {
    id: 'u1',
    username: 'steven',
    password: 'password123',
    email: 'Steven_Miller@mymail.com',
    firstName: 'Steven',
    lastName: 'Miller',
  },
  {
    id: 'u2',
    username: 'jessica',
    password: 'password123',
    email: 'Jessica_Williams@mymail.com',
    firstName: 'Jessica',
    lastName: 'Williams',
  },
  {
    id: 'u3',
    username: 'michael',
    password: 'password123',
    email: 'Michael_Brown@mymail.com',
    firstName: 'Michael',
    lastName: 'Brown',
  },
  {
    id: 'u4',
    username: 'sarah',
    password: 'password123',
    email: 'Sarah_Davis@mymail.com',
    firstName: 'Sarah',
    lastName: 'Davis',
  },
  {
    id: 'u5',
    username: 'david',
    password: 'password123',
    email: 'David_Wilson@mymail.com',
    firstName: 'David',
    lastName: 'Wilson',
  },
]

// Additional users that are not friends with anyone (for search/requests)
export const otherUsers: User[] = [
  {
    id: 'u6',
    username: 'amanda',
    password: 'password123',
    email: 'Amanda_Chen@mymail.com',
    firstName: 'Amanda',
    lastName: 'Chen',
  },
  {
    id: 'u7',
    username: 'robert',
    password: 'password123',
    email: 'Robert_Taylor@mymail.com',
    firstName: 'Robert',
    lastName: 'Taylor',
  },
  {
    id: 'u8',
    username: 'lisa',
    password: 'password123',
    email: 'Lisa_Anderson@mymail.com',
    firstName: 'Lisa',
    lastName: 'Anderson',
  },
]

// All users combined
export const allUsers = [...users, ...otherUsers]

// Channels the current user may want to discover/join
export const discoverableChannels: Channel[] = [
  {
    id: 'c4',
    name: '#engineering',
    type: 'public',
    memberIds: [],
  },
  {
    id: 'c5',
    name: '#random',
    type: 'public',
    memberIds: [],
  },
  {
    id: 'c6',
    name: '#design',
    type: 'public',
    memberIds: [],
  },
]

// Steven (u1) is friends with every other user
export const friendships: Friendship[] = [
  { id: 'f1', userId1: 'u1', userId2: 'u2' },
  { id: 'f2', userId1: 'u1', userId2: 'u3' },
  { id: 'f3', userId1: 'u1', userId2: 'u4' },
  { id: 'f4', userId1: 'u1', userId2: 'u5' },
]

// Three sample public channels (all users are members)
export const channels: Channel[] = [
  {
    id: 'c1',
    name: '#general',
    type: 'public',
    memberIds: ['u1', 'u2', 'u3', 'u4', 'u5'],
  },
  {
    id: 'c2',
    name: '#news',
    type: 'public',
    memberIds: ['u1', 'u2', 'u3', 'u4', 'u5'],
  },
  {
    id: 'c3',
    name: '#technical_support',
    type: 'public',
    memberIds: ['u1', 'u2', 'u3', 'u4', 'u5'],
  },
  // Steven's DM with Jessica (u2)
  {
    id: 'dm1',
    name: 'Steven & Jessica',
    type: 'dm',
    memberIds: ['u1', 'u2'],
  },
  // Group chat: Steven, Jessica, Sarah
  {
    id: 'g1',
    name: 'Steven, Jessica & Sarah',
    type: 'group',
    memberIds: ['u1', 'u2', 'u4'],
  },
]

export const getFriendsForUser = (userId: string): User[] => {
  const friendIds = friendships
    .filter(f => f.userId1 === userId || f.userId2 === userId)
    .map(f => (f.userId1 === userId ? f.userId2 : f.userId1))
  return users.filter(u => friendIds.includes(u.id))
}

// Creates a new dm or group channel if one with the exact same members doesn't already exist.
// Returns the existing or newly created channel.
export const createDmOrGroupChannel = (
  currentUserId: string,
  otherMemberIds: string[]
): Channel => {
  const memberIds = [currentUserId, ...otherMemberIds].sort()
  const existing = channels.find(ch => {
    if (ch.type !== 'dm' && ch.type !== 'group') return false
    const ids = [...ch.memberIds].sort()
    return ids.length === memberIds.length && ids.every((id, i) => id === memberIds[i])
  })
  if (existing) return existing

  const isGroup = otherMemberIds.length > 1
  const id = isGroup ? `g${Date.now()}` : `dm${Date.now()}`
  const otherNames = otherMemberIds
    .map(id => users.find(u => u.id === id))
    .filter(Boolean) as User[]
  const name = isGroup
    ? `${otherNames.map(u => `${u.firstName} ${u.lastName}`).join(', ')}`
    : `${otherNames[0]?.firstName ?? 'Unknown'} ${otherNames[0]?.lastName ?? ''}`

  const channel: Channel = {
    id,
    name,
    type: isGroup ? 'group' : 'dm',
    memberIds: [currentUserId, ...otherMemberIds],
  }
  channels.push(channel)
  return channel
}

export const messages: Message[] = [
  // #general messages
  {
    id: 'm1',
    channelId: 'c1',
    authorId: 'u1',
    content: 'Hey everyone, welcome to #general!',
    timestamp: Date.now() - 3600000 * 24, // 1 day ago
  },
  {
    id: 'm2',
    channelId: 'c1',
    authorId: 'u2',
    content: 'Thanks Steven! Excited to be here.',
    timestamp: Date.now() - 3600000 * 23,
  },
  {
    id: 'm3',
    channelId: 'c1',
    authorId: 'u3',
    content: 'What are we working on today?',
    timestamp: Date.now() - 3600000 * 22,
  },
  {
    id: 'm4',
    channelId: 'c1',
    authorId: 'u1',
    content: 'Just setting up the new chat app. Let me know if you find any bugs!',
    timestamp: Date.now() - 3600000 * 20,
  },

  // #news messages
  {
    id: 'm5',
    channelId: 'c2',
    authorId: 'u1',
    content: 'New update rolled out this morning.',
    timestamp: Date.now() - 3600000 * 12,
  },
  {
    id: 'm6',
    channelId: 'c2',
    authorId: 'u4',
    content: 'What are the main changes?',
    timestamp: Date.now() - 3600000 * 11,
  },
  {
    id: 'm7',
    channelId: 'c2',
    authorId: 'u1',
    content: 'Performance improvements and bug fixes mainly.',
    timestamp: Date.now() - 3600000 * 10,
  },

  // #technical_support messages
  {
    id: 'm8',
    channelId: 'c3',
    authorId: 'u5',
    content: 'Having trouble connecting to the API.',
    timestamp: Date.now() - 3600000 * 6,
  },
  {
    id: 'm9',
    channelId: 'c3',
    authorId: 'u1',
    content: 'What error are you seeing?',
    timestamp: Date.now() - 3600000 * 5,
  },
  {
    id: 'm10',
    channelId: 'c3',
    authorId: 'u5',
    content: 'Getting a 502 on the websocket handshake.',
    timestamp: Date.now() - 3600000 * 4,
  },
  {
    id: 'm11',
    channelId: 'c3',
    authorId: 'u3',
    content: 'Check if the gateway service is running.',
    timestamp: Date.now() - 3600000 * 3,
  },

  // Steven & Jessica DM
  {
    id: 'm12',
    channelId: 'dm1',
    authorId: 'u1',
    content: 'Hey Jessica, are we still on for lunch tomorrow?',
    timestamp: Date.now() - 3600000 * 2,
  },
  {
    id: 'm13',
    channelId: 'dm1',
    authorId: 'u2',
    content: 'Yes! Let me check my calendar and get back to you.',
    timestamp: Date.now() - 3600000 * 1.5,
  },
  {
    id: 'm14',
    channelId: 'dm1',
    authorId: 'u1',
    content: 'Sure, no rush. How is the project going btw?',
    timestamp: Date.now() - 3600000 * 1,
  },
  {
    id: 'm15',
    channelId: 'dm1',
    authorId: 'u2',
    content: 'Pretty well! Almost done with the frontend mockup.',
    timestamp: Date.now() - 3600000 * 0.5,
  },

  // Group chat: Steven, Jessica, Sarah
  {
    id: 'm16',
    channelId: 'g1',
    authorId: 'u1',
    content: 'Hey you two, just setting up this group chat for the project!',
    timestamp: Date.now() - 3600000 * 8,
  },
  {
    id: 'm17',
    channelId: 'g1',
    authorId: 'u2',
    content: 'Great idea! We needed a space to coordinate.',
    timestamp: Date.now() - 3600000 * 7.5,
  },
  {
    id: 'm18',
    channelId: 'g1',
    authorId: 'u4',
    content: 'Agreed. When is our first standup?',
    timestamp: Date.now() - 3600000 * 7,
  },
  {
    id: 'm19',
    channelId: 'g1',
    authorId: 'u1',
    content: 'How about Monday at 10am?',
    timestamp: Date.now() - 3600000 * 6,
  },
  {
    id: 'm20',
    channelId: 'g1',
    authorId: 'u2',
    content: 'Works for me!',
    timestamp: Date.now() - 3600000 * 5,
  },
]
